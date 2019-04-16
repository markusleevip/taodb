package pool_test

import (
	"github.com/markusleevip/taodb/pool"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Pool", func() {
	var (
		subject *pool.Pool
		server  *httptest.Server
		factory pool.Factory
	)

	BeforeEach(func() {
		server, factory = mockServer()

		var err error
		subject, err = pool.New(&pool.Options{
			InitialSize: 3,
			MaxCap:      5,
		}, factory)
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		Expect(subject.Close()).To(Succeed())
		server.Close()
	})

	It("should check len", func() {
		Expect(subject.Len()).To(Equal(3))
	})

	It("should get/put", func() {
		cn1, err := subject.Get()
		Expect(err).NotTo(HaveOccurred())
		Expect(subject.Len()).To(Equal(2))

		cn2, err := subject.Get()
		Expect(err).NotTo(HaveOccurred())
		Expect(subject.Len()).To(Equal(1))

		cn3, err := subject.Get()
		Expect(err).NotTo(HaveOccurred())
		Expect(subject.Len()).To(Equal(0))

		cn4, err := subject.Get()
		Expect(err).NotTo(HaveOccurred())
		Expect(subject.Len()).To(Equal(0))

		cn5, err := subject.Get()
		Expect(err).NotTo(HaveOccurred())

		cn6, err := subject.Get()
		Expect(err).NotTo(HaveOccurred())

		Expect(subject.Put(cn1)).To(BeTrue())
		Expect(subject.Put(cn2)).To(BeTrue())
		Expect(subject.Put(cn3)).To(BeTrue())
		Expect(subject.Put(cn4)).To(BeTrue())
		Expect(subject.Len()).To(Equal(4))

		Expect(subject.Put(cn5)).To(BeTrue())
		Expect(subject.Len()).To(Equal(5))

		Expect(subject.Put(cn6)).To(BeFalse())
		Expect(subject.Len()).To(Equal(5))

		_, err = cn6.Write([]byte("x"))
		Expect(err).To(HaveOccurred())
		Expect(err.(*net.OpError).Err).To(MatchError("use of closed network connection"))

		last, err := subject.Get()
		Expect(err).NotTo(HaveOccurred())
		Expect(last).NotTo(BeNil())
		Expect(last).NotTo(Equal(cn1))
		Expect(last).NotTo(Equal(cn2))
		Expect(last).NotTo(Equal(cn3))
		Expect(last).NotTo(Equal(cn4))
		Expect(last).To(Equal(cn5))
	})

	It("should be thread-safe", func() {
		n := 10000
		if testing.Short() {
			n = 100
		}

		wg := new(sync.WaitGroup)
		for c := 0; c < 10; c++ {
			wg.Add(1)
			go func() {
				defer GinkgoRecover()
				defer wg.Done()

				for i := 0; i < n; i++ {
					cn, err := subject.Get()
					Expect(err).NotTo(HaveOccurred())

					subject.Put(cn)
				}
			}()
		}
		wg.Wait()

		Expect(subject.Len()).To(Equal(5))
	})

	It("should reap idle connections", func() {
		p, err := pool.New(&pool.Options{
			InitialSize:  3,
			MaxCap:       5,
			ReapInterval: 200 * time.Millisecond,
			IdleTimeout:  300 * time.Millisecond,
		}, factory)
		Expect(err).NotTo(HaveOccurred())
		defer p.Close()

		Expect(p.Len()).To(Equal(3))
		Eventually(func() (int, error) {
			cn, err := p.Get()
			if err != nil {
				return 0, err
			}
			if err != nil {
				return 0, err
			}

			p.Put(cn)
			return p.Len(), nil
		}).Should(Equal(1))
	})

})

// --------------------------------------------------------------------

func mockServer() (*httptest.Server, pool.Factory) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	factory := func() (net.Conn, error) {
		return net.Dial("tcp", strings.Replace(server.URL, "http://", "", -1))
	}
	return server, factory
}

// --------------------------------------------------------------------

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "pool")
}

func BenchmarkPool(b *testing.B) {
	srv, factory := mockServer()
	defer srv.Close()

	p, err := pool.New(nil, factory)
	if err != nil {
		b.Fatal(err)
	}
	defer p.Close()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			cn, err := p.Get()
			if err != nil {
				b.Fatal(err)
			}
			p.Put(cn)
		}
	})
}
