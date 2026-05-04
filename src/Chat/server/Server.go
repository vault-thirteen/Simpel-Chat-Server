package server

import (
	"context"
	"errors"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	jrm1 "github.com/vault-thirteen/JSON-RPC-M1"

	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/settings"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/helper"
)

type Server struct {
	settings           *settings.ChatServerSettings
	rp                 *jrm1.Processor
	criticalErrorsChan *chan error
	listenDsn          string
	httpServer         *http.Server
	startTime          time.Time
	isRunning          *atomic.Bool
	stopTime           time.Time
	listenError        error
}

func NewServer(
	settings *settings.ChatServerSettings,
	criticalErrorsChan *chan error,
	rp *jrm1.Processor,
	serverStartTimeTS time.Time,
) (srv *Server, err error) {
	srv = &Server{
		settings:           settings,
		rp:                 rp,
		criticalErrorsChan: criticalErrorsChan,
		listenDsn:          net.JoinHostPort(settings.HostName, strconv.Itoa(int(settings.PortNumber))),
		startTime:          serverStartTimeTS,
		isRunning:          new(atomic.Bool),
	}

	srv.httpServer = &http.Server{
		Addr:    srv.listenDsn,
		Handler: http.Handler(http.HandlerFunc(srv.router)),
	}

	go srv.run()

	return srv, nil
}

func (s *Server) router(rw http.ResponseWriter, req *http.Request) {
	left, right, ok := strings.Cut(req.URL.Path, helper.UrlPathSeparator)
	if !ok {
		s.httpRespond_BadRequest(rw)
		return
	}

	if len(left) != 0 {
		s.httpRespond_NotFound(rw)
		return
	}

	switch right {
	case helper.UrlPath_Api:
		s.rp.ServeHTTP(rw, req)
		return

	default:
		s.httpRespond_NotFound(rw)
		return
	}
}

func (s *Server) httpRespond_BadRequest(rw http.ResponseWriter) {
	rw.WriteHeader(http.StatusBadRequest)
}
func (s *Server) httpRespond_NotFound(rw http.ResponseWriter) {
	rw.WriteHeader(http.StatusNotFound)
}

// Async.
func (s *Server) run() {
	s.isRunning.Store(true)

	defer func() {
		s.stopTime = time.Now().UTC()
		s.isRunning.Store(false)
	}()

	// If either the server crashes or it is stopped manually, we get here. So,
	// the function returns in any case. So, it is safe to leave this function
	// as it is without any WaitGroup guards.
	s.listenError = s.httpServer.ListenAndServeTLS(s.settings.CertFile, s.settings.KeyFile)

	// Report about a critical error to the parent object.
	if !errors.Is(s.listenError, http.ErrServerClosed) {
		*s.criticalErrorsChan <- s.listenError
	}
}

func (s *Server) IsRunning() bool {
	return s.isRunning.Load()
}

func (s *Server) Stop() (err error) {
	ctx, cf := context.WithTimeout(context.Background(), time.Minute)
	defer cf()

	err = s.httpServer.Shutdown(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) GetStartTime() time.Time { return s.startTime }
