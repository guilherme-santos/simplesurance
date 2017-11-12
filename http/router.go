package http

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/guilherme-santos/simplesurance"
)

type Router struct {
	server *http.Server
}

func NewRouter(handlers ...simplesurance.HTTPHandler) *Router {
	router := &Router{}
	for _, h := range handlers {
		h.RegisterRoutes(router)
	}

	return router
}

func (r *Router) Run(port string) error {
	log.Println("Running server on port", port)

	r.server = &http.Server{Addr: ":" + port}
	return r.server.ListenAndServe()
}

func (r *Router) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return r.server.Shutdown(ctx)
}

func (r *Router) Get(path string, handler http.HandlerFunc) {
	http.HandleFunc(path, handler)
}
