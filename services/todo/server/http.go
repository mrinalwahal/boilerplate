package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mrinalwahal/boilerplate/services/todo/service"
)

type HTTPServer struct {

	//	Port
	port string

	//	Router
	router http.Handler

	//	Service
	service service.Service
}

func NewHTTPServer(config *NewHTTPServerConfig) Server {

	router := config.Router
	if router == nil {

		r := chi.NewRouter()

		// A good base middleware stack
		r.Use(middleware.RequestID)
		r.Use(middleware.RealIP)
		r.Use(middleware.Logger)
		r.Use(middleware.Recoverer)

		// Set a timeout value on the request context (ctx), that will signal
		// through ctx.Done() that the request has timed out and further
		// processing should be stopped.
		r.Use(middleware.Timeout(60 * time.Second))

		r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
		router = r
	}

	return &HTTPServer{
		router:  router,
		service: config.Service,
	}
}

type NewHTTPServerConfig struct {

	//	Port
	Port string

	//	Router
	Router http.Handler

	//	Service
	Service service.Service
}

func (s *HTTPServer) Serve() {
	http.ListenAndServe(fmt.Sprintf(":%s", s.port), s.router)
}
