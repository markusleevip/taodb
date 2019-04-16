// Package pool is a generic, high-performance pool for net.Conn
// objects.
package pool

import (
	"net"
	"sync"
	"sync/atomic"
	"time"
)

// Factory must returns new connections
type Factory func() (net.Conn, error)

// Options can tweak Pool configuration
type Options struct {
	// InitialSize creates a number of connection on pool initialization
	// Default: 0
	InitialSize int

	// MaxCap sets the maximum pool capacity. Will be automatically adjusted when InitialSize
	// is larger.
	// Default: 10
	MaxCap int

	// IdleTimeout timeout after which connections are reaped and
	// automatically removed from the pool.
	// Default: 0 (= never)
	IdleTimeout time.Duration

	// ReapInterval determines the frequency of reap cycles
	// Default: 1 minute
	ReapInterval time.Duration
}

func (o *Options) norm() Options {
	x := *o
	if x.ReapInterval <= 0 {
		x.ReapInterval = time.Minute
	}
	if x.MaxCap <= 0 {
		x.MaxCap = 10
	}
	if x.MaxCap < x.InitialSize {
		x.MaxCap = x.InitialSize
	}
	return x
}

type none struct{}

// Pool contains a number of connections
type Pool struct {
	conns   []member
	opt     Options
	factory Factory

	dying, dead chan none

	avail  uint32
	closed int32

	mu sync.Mutex
}

// New creates a pool with an initial number of connection and a maximum cap
func New(opt *Options, factory Factory) (*Pool, error) {
	if opt == nil {
		opt = new(Options)
	}

	p := &Pool{
		conns:   make([]member, 0, opt.MaxCap),
		factory: factory,
		opt:     opt.norm(),
		dying:   make(chan none),
		dead:    make(chan none),
	}

	for i := 0; i < opt.InitialSize; i++ {
		cn, err := factory()
		if err != nil {
			_ = p.close()
			return nil, err
		}
		p.Put(cn)
	}

	go p.loop()
	return p, nil
}

// Len returns the number of available connections in the pool
func (s *Pool) Len() int { return int(atomic.LoadUint32(&s.avail)) }

// Get returns a connection from the pool or creates a new one
func (s *Pool) Get() (net.Conn, error) {
	if cn := s.pop(); cn != nil {
		return cn, nil
	}

	return s.factory()
}

// Put adds/returns a connection to the pool
func (s *Pool) Put(cn net.Conn) bool {
	if s.Len() >= s.opt.MaxCap || atomic.LoadInt32(&s.closed) == 1 {
		_ = cn.Close()
		return false
	}

	m := member{cn: cn, lastAccess: time.Now()}
	s.mu.Lock()
	s.conns = append(s.conns, m)
	atomic.StoreUint32(&s.avail, uint32(len(s.conns)))
	s.mu.Unlock()

	return true
}

// Close closes all connections and the pool
func (s *Pool) Close() error {
	if !atomic.CompareAndSwapInt32(&s.closed, 0, 1) {
		return nil
	}

	close(s.dying)
	<-s.dead
	return s.close()
}

func (s *Pool) pop() net.Conn {
	s.mu.Lock()

	pos := len(s.conns) - 1
	if pos < 0 {
		s.mu.Unlock()
		return nil
	}

	m := s.conns[pos]
	s.conns = s.conns[:pos]
	atomic.StoreUint32(&s.avail, uint32(len(s.conns)))
	s.mu.Unlock()

	return m.cn

}

func (s *Pool) close() (err error) {
	for {
		cn := s.pop()
		if cn == nil {
			break
		}
		if e := cn.Close(); e != nil {
			err = e
		}
	}
	return err
}

func (s *Pool) reap() {
	timeout := s.opt.IdleTimeout
	if timeout <= 0 {
		return
	}

	cutoff := time.Now().Add(-timeout)

	s.mu.Lock()
	if sz := len(s.conns); sz != 0 {
		if m := s.conns[0]; m.lastAccess.Before(cutoff) {
			defer m.cn.Close()

			copy(s.conns, s.conns[1:])
			s.conns = s.conns[:sz-1]
			atomic.StoreUint32(&s.avail, uint32(len(s.conns)))
		}
	}
	s.mu.Unlock()
}

func (s *Pool) loop() {
	defer close(s.dead)

	ticker := time.NewTicker(s.opt.ReapInterval)
	defer ticker.Stop()

	for {
		select {
		case <-s.dying:
			return
		case <-ticker.C:
			s.reap()
		}
	}
}

type member struct {
	cn         net.Conn
	lastAccess time.Time
}
