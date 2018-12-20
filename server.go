package graceful

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/golang/glog"
	"github.com/pressly/chi"
)

const (
	contextTimeoutHeader = "context-timeout"
)

func init() {
	http.DefaultTransport.(*http.Transport).MaxIdleConns = 1000
	http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = 1000
}

type Server struct {
	mux          *chi.Mux
	serv         http.Server
	errc         chan error
	shutdownFlag int32
}

func NewServer() *Server {
	s := &Server{
		mux:  chi.NewMux(),
		errc: make(chan error, 1),
	}
	// s.restoreContextTimeoutMiddleware()
	s.serv.Handler = s.mux
	return s
}

func (s *Server) Use(m ...func(http.Handler) http.Handler) {
	s.mux.Use(m...)
}

func (s *Server) HandleHostname() {
	h, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	s.mux.Get("/hostname", func(w http.ResponseWriter, req *http.Request) {
		_, err := fmt.Fprintf(w, "%s", h)
		if err != nil {
			log.Printf("response error: %v", err)
		}
	})
}

func (s *Server) Handle(p string, h http.Handler) {
	s.mux.Handle(p, h)
}

func (s *Server) Mount(dir string, h http.Handler) {
	if !strings.HasSuffix(dir, "/") {
		dir += "/"
	}
	s.mux.Mount(dir, h)
}

func (s *Server) ListenAndServe(addr string) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return s.Serve(l)
}

func (s *Server) Serve(l net.Listener) error {
	err := s.serv.Serve(l)
	s.errc <- err
	return err
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.mux.ServeHTTP(w, req)
}

func (s *Server) Stop() error {
	if err := s.serv.Close(); err != nil {
		return err
	}
	return <-s.errc
}

func (s *Server) Shutdown(ctx context.Context, d time.Duration) error {
	atomic.StoreInt32(&s.shutdownFlag, 1)
	time.Sleep(d)

	if err := s.serv.Shutdown(ctx); err != nil {
		return err
	}

	return <-s.errc
}

func (s *Server) UseShutdownMiddleware(code int, body []byte) {
	m := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if atomic.LoadInt32(&s.shutdownFlag) == 0 {
				next.ServeHTTP(w, r)
				return
			}

			w.WriteHeader(code)
			if _, err := w.Write(body); err != nil {
				glog.Errorf("write shutdown response error: %v", err)
			}
		}
		return http.HandlerFunc(fn)
	}
	s.mux.Use(m)
}

// restoreContextTimeoutMiddleware write timeout from http-header to context
func (s *Server) restoreContextTimeoutMiddleware() {
	m := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			str := r.Header.Get(contextTimeoutHeader)
			timeout, err := strconv.ParseInt(str, 16, 64)
			if err == nil {
				ctx, cancel := context.WithTimeout(r.Context(), time.Duration(timeout))
				defer cancel()
				r = r.WithContext(ctx)
				r.Header.Del(contextTimeoutHeader)
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
	s.mux.Use(m)
}
