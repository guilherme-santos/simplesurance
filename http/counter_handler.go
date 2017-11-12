package http

import (
	"fmt"
	"net/http"

	"github.com/guilherme-santos/simplesurance"
)

type CounterHandler struct {
	counterService simplesurance.CounterService
}

func NewCounterHandler(counterService simplesurance.CounterService) *CounterHandler {
	return &CounterHandler{
		counterService: counterService,
	}
}

func (s *CounterHandler) RegisterRoutes(router simplesurance.Router) {
	router.Get("/", s.handleGet)
}

func (s *CounterHandler) handleGet(w http.ResponseWriter, r *http.Request) {
	totalRequests := s.counterService.NewRequest()
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(fmt.Sprintf("%d\n", totalRequests)))
}
