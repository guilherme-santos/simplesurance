package http_test

import (
	"io/ioutil"
	"log"
	gohttp "net/http"
	"strings"
	"testing"
	"time"

	"github.com/guilherme-santos/simplesurance"
	"github.com/guilherme-santos/simplesurance/http"
)

type testHandler struct {
	Endpoint             string
	Response             string
	RegisterRoutesCalled bool
}

func (h *testHandler) RegisterRoutes(router simplesurance.Router) {
	h.RegisterRoutesCalled = true
	router.Get(h.Endpoint, h.handleGet)
}

func (h *testHandler) handleGet(w gohttp.ResponseWriter, r *gohttp.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(h.Response))
}

func init() {
	// Disable log
	log.SetOutput(ioutil.Discard)
}

var routerPort = "8282"

func runServer(router *http.Router, errChan chan error) {
	go func() {
		err := router.Run(routerPort)
		errChan <- err
	}()

	// Just to guarantee that server is running
	time.Sleep(1 * time.Second)
}

func TestNewRouter_CallRegisterRoutesForAllHandlers(t *testing.T) {
	testHandlerOne := &testHandler{
		Endpoint: "/test-one",
	}
	testHandlerTwo := &testHandler{
		Endpoint: "/test-two",
	}
	http.NewRouter(testHandlerOne, testHandlerTwo)

	if !testHandlerOne.RegisterRoutesCalled {
		t.Error("It was expected testHandlerOne.RegisterRoutes be called but don't")
	}
	if !testHandlerTwo.RegisterRoutesCalled {
		t.Error("It was expected testHandlerTwo.RegisterRoutes be called but don't")
	}
}

func TestRunAndStop(t *testing.T) {
	r := http.NewRouter()
	errChan := make(chan error)

	runServer(r, errChan)

	url := "http://localhost:" + routerPort
	_, err := gohttp.Get(url)
	if err != nil {
		t.Errorf("Error making request to %s: %s", url, err)
		return
	}

	err = r.Stop()
	if err != nil {
		t.Error("Error stoping router:", err)
		return
	}

	select {
	case err = <-errChan:
		if err != gohttp.ErrServerClosed {
			t.Error("It was expected http.ErrServerClosed but reveived:", err)
		}
	case <-time.After(time.Second * 3):
		t.Error("Server took so long to shutdown")
	}
}

func TestRegistryHandler(t *testing.T) {
	handler := &testHandler{
		Endpoint: "/test",
		Response: "my test handler",
	}

	r := http.NewRouter(handler)
	errChan := make(chan error)

	runServer(r, errChan)

	url := "http://localhost:" + routerPort + handler.Endpoint
	resp, err := gohttp.Get(url)
	if err != nil {
		t.Errorf("Error making request to %s: %s", url, err)
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	bodyStr := string(body)
	if !strings.EqualFold(handler.Response, bodyStr) {
		t.Errorf("It was expected '%s' as response but received: '%s'", handler.Response, strings.TrimSpace(bodyStr))
	}

	err = r.Stop()
	if err != nil {
		t.Error("Error stoping router:", err)
	}
}
