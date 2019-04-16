package taodb

import (
	"github.com/markusleevip/taodb/leveldb"
	"github.com/markusleevip/taodb/resp"
	"net"
	"strings"
	"sync"
	"time"
)

type Server struct{
	config	*Config
	info	*ServerInfo
	cmds 	map[string] interface{}
	db   	*DB
	mu		sync.RWMutex
}

// NewServer creates a new server instance
func NewServer(config *Config,dbPath string) *Server{
	if config == nil{
		config = new(Config)
	}
	var db DB = leveldb.NewDB(dbPath)
	return &Server{
		config :config,
		info: newServerInfo(),
		cmds: make(map[string]interface{}),
		db: &db,}
}


// Info returns the server info registry
func (srv *Server) Info() *ServerInfo { return srv.info }


// Handle registers a handler for a command.
func (srv *Server) Handle(name string, h Handler) {
	srv.mu.Lock()
	srv.cmds[strings.ToLower(name)] = h
	srv.mu.Unlock()
}

// HandleFunc registers a handler func for a command.
func (srv *Server) HandleFunc(name string, fn HandlerFunc) {
	srv.Handle(name, fn)
}


// HandleStream registers a handler for a streaming command.
func (srv *Server) HandleStream(name string, h StreamHandler) {
	srv.mu.Lock()
	srv.cmds[strings.ToLower(name)] = h
	srv.mu.Unlock()
}


// HandleStreamFunc registers a handler func for a command
func (srv *Server) HandleStreamFunc(name string, fn StreamHandlerFunc) {
	srv.HandleStream(name, fn)
}


// Serve accepts incoming connections on a listener, creating a
// new service goroutine for each.
func (srv *Server) Serve(lis net.Listener) error {
	for {
		cn, err:= lis.Accept()
		if err != nil{
			return err
		}

		if ka := srv.config.TCPKeepAlive; ka >0{
			if tc, ok := cn.(*net.TCPConn); ok {
				tc.SetKeepAlive(true)
				tc.SetKeepAlivePeriod(ka)
			}
		}

		go srv.serveClient(newClient(cn))

	}
}

// Starts a new session, serving client
func (srv *Server) serveClient(c *Client) {
	// Release client on exit
	defer c.release()

	// Register client
	srv.info.Register(c)
	defer srv.info.deregister(c.id)

	// Create perform callback
	perform := func(name string) error {
		return srv.perform(c,name)
	}
	for !c.closed{
		if d:= srv.config.Timeout; d>0{
			c.cn.SetDeadline(time.Now().Add(d))
		}

		// perform pipeline
		if err := c.pipeline(perform) ; err != nil{
			c.wr.AppendError("ERR "+err.Error())

			if !resp.IsProtocolError(err) {
				c.wr.Flush()
				return
			}
		}

		// flush buffer, return on errors
		if err := c.wr.Flush(); err != nil{
			return
		}
	}
}

func (srv *Server) perform(c *Client, name string) (err error) {
	norm := strings.ToLower(name)

	// find handler
	srv.mu.RLock()
	h, ok := srv.cmds[norm]
	srv.mu.RUnlock()

	if !ok {
		c.wr.AppendError(UnknownCommand(name))
		c.rd.SkipCmd()
		return
	}

	// register call
	srv.info.command(c.id, norm)

	switch handler := h.(type) {
	case Handler:
		if c.cmd, err = c.readCmd(c.cmd); err != nil{
			return
		}
		handler.ServeRedeo(c.wr, c.cmd)

	case StreamHandler:
		if c.scmd, err = c.streamCmd(c.scmd) ; err != nil{
			return
		}
		defer c.scmd.Discard()

		handler.ServeRedeoStream(c.wr, c.scmd)
	}

	// flush when buffer is large enough
	if n := c.wr.Buffered(); n > resp.MaxBufferSize/2 {
		err = c.wr.Flush()
		if err !=nil{
			return err
		}
	}
	return
}