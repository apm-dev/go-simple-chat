package rest

import (
	"fmt"
	"net/http"

	"apm.dev/go-simple-chat/src/pkg/logger"
	ctrl "apm.dev/go-simple-chat/src/presentation/rest/controllers"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type RestServer struct {
	server *http.Server
	quit   chan bool
	auth   *ctrl.AuthController
}

func NewServer(auth *ctrl.AuthController) *RestServer {
	return &RestServer{
		server: &http.Server{},
		quit:   make(chan bool, 1),
		auth:   auth,
	}
}

func (s *RestServer) Start(addr string) {

	fmt.Println("Starting RestServer on", addr)

	r := gin.Default()

	s.registerRoutes(r)

	s.server.Addr = addr
	s.server.Handler = r

	// start http server on different goroutine
	go func() {
		if err := s.server.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				fmt.Println("RestServer stopped under request")
			} else {
				logger.Error(errors.Wrap(err, "RestServer stopped unexpectedly"))
				panic(err)
			}
		}
	}()

	// listen to quit channel to close the server
	go func() {
		<-s.quit
		if err := s.server.Close(); err != nil {
			logger.Error(errors.Wrap(err, "failed to stop RestServer"))
		}
	}()
}

func (s *RestServer) Stop() {
	fmt.Println("Stopping RestServer on", s.server.Addr)
	s.quit <- true
}
