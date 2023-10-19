package server

import (
	"context"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/mfajri11/family-catering-micro-service/auth-service/config"
	"github.com/mfajri11/family-catering-micro-service/auth-service/pkg/log"
)

var defaultShutdownTimeout time.Duration = 5 * time.Second

type Server struct {
	httpSrv *http.Server
	// grpcSrv         *grpc.Server
	errChan         chan error
	sigChan         chan os.Signal
	shutdownTimeout time.Duration
	httpListener    net.Listener
}

// New return Server struct
func New(address string, handler http.Handler, opts ...Option) *Server {
	httpsrv := &http.Server{
		Addr:    address,
		Handler: handler,
	}
	srv := &Server{
		httpSrv:         httpsrv,
		errChan:         make(chan error, 1),
		shutdownTimeout: defaultShutdownTimeout,
	}

	for _, opt := range opts {
		opt(srv)
	}

	return srv
}

// Notify return a channel which will be closed when a server got an error from listening
func (s *Server) Notify() <-chan os.Signal {
	return s.sigChan
}

func (s *Server) Err() <-chan error {
	return s.errChan
}

// TODO: implement seamless restart

// Start will start the server and will notify the error channel if error occur
func (s *Server) Serve() error {
	log.Info("server.Server.Start: %s listen at %s\n", config.Cfg.App.Name, s.httpSrv.Addr)
	l, err := getOrNewListener(s.httpSrv.Addr)
	if err != nil {
		return err
	}

	s.httpListener = l

	err = s.httpSrv.Serve(l)
	// err := s.httpSrv.ListenAndServe()
	if err != nil {
		// TODO: wrap error
		s.errChan <- err
		// sender of the channel should be responsible for closing the channel
		close(s.errChan)
		return err
	}

	return nil
}

// Shutdown shutdown the server gracefully and wait until the server time out is passed and return the error if there is an error while performing the shutdown
func (s *Server) Shutdown() error {
	log.Info("server.Server.Shutdown: %s starting to shutdown gracefully", config.Cfg.App.Name)

	ctx, cancelFunc := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancelFunc()

	s.httpSrv.SetKeepAlivesEnabled(false)
	err := s.httpSrv.Shutdown(ctx)
	if err != nil {
		return err // TODO: wrap error
	}

	return nil
}

// Restart seamlessly restart the server by make a fork child of current process and then exit itself
// ! need to confirm does forking a child process copy the exact source file or it goes the way like re build does?
func (s *Server) Restart() error {
	log.Info("server.Server.Restart: %s starting seamless restart", config.Cfg.App.Name)
	var (
		ln  net.Listener
		err error
	)

	ln = s.httpListener

	if ln == nil {
		ln, err = getOrNewListener(s.httpSrv.Addr)
	}

	if err != nil {
		// TODO: wrap error
		return err
	}
	p, err := forkChild(s.httpSrv.Addr, ln)
	if err != nil {
		// TODO: wrap error
		return err
	}

	log.Info("server.Server.Restart: fork child successfully created with pid: %d", p.Pid)
	return s.Shutdown()
}
