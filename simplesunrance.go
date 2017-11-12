package simplesurance

import "net/http"

type (
	CounterService interface {
		TotalRequests() int
		NewRequest() int
		Start()
		Stop()
	}

	Router interface {
		Run(port string) error
		Stop() error
		Get(path string, handler http.HandlerFunc)
	}

	HTTPHandler interface {
		RegisterRoutes(Router)
	}
)
