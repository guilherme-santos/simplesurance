package file_test

import (
	"encoding/gob"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/guilherme-santos/simplesurance/file"
)

func init() {
	// Disable log
	log.SetOutput(ioutil.Discard)
}

func createdFile(filename *string) func() {
	tmpfile, _ := ioutil.TempFile("", "simplesurance-api")
	*filename = tmpfile.Name()
	return func() {
		os.Remove(tmpfile.Name())
	}
}

func remove(f func()) {
	f()
}

func TestNewCounterService_CanLoadFile(t *testing.T) {
	var filename string
	defer remove(createdFile(&filename))

	// write requests to file
	requests := []time.Time{time.Now(), time.Now(), time.Now(), time.Now(), time.Now()}

	f, _ := os.OpenFile(filename, os.O_WRONLY, 0644)
	enc := gob.NewEncoder(f)
	err := enc.Encode(requests)
	if err != nil {
		t.Error("Cannot encode list of request:", err)
	}

	counter := file.NewCounterService(filename)

	if c := counter.TotalRequests(); c != 5 {
		t.Error("It was expected 5 but get", c)
		return
	}
}

func TestTotalRequests(t *testing.T) {
	var filename string
	defer remove(createdFile(&filename))

	counter := file.NewCounterService(filename)

	if c := counter.TotalRequests(); c != 0 {
		t.Error("It was expected 0 but get", c)
		return
	}

	counter.NewRequest()

	if c := counter.TotalRequests(); c != 1 {
		t.Error("It was expected 1 but get", c)
		return
	}
}

func TestNewRequest_ReturnNextValue(t *testing.T) {
	var filename string
	defer remove(createdFile(&filename))

	counter := file.NewCounterService(filename)

	if c := counter.TotalRequests(); c != 0 {
		t.Error("It was expected 0 but get", c)
		return
	}

	if c := counter.NewRequest(); c != 1 {
		t.Error("It was expected 1 but get", c)
		return
	}
}

func TestNewRequest_InRaceCondition(t *testing.T) {
	var filename string
	defer remove(createdFile(&filename))

	counter := file.NewCounterService(filename)

	var wg sync.WaitGroup
	wg.Add(1000)
	for i := 0; i < 1000; i++ {
		go func() {
			counter.TotalRequests()
			counter.NewRequest()
			wg.Done()
		}()
	}

	wg.Wait()
}

func TestStart_DecreaseValueAfterWindowSize(t *testing.T) {
	file.WindowSize = 2 * time.Second

	var filename string
	defer remove(createdFile(&filename))

	counter := file.NewCounterService(filename)
	counter.Start()

	counter.NewRequest()
	counter.NewRequest()
	if c := counter.TotalRequests(); c != 2 {
		t.Error("It was expected 2 but get", c)
		return
	}

	time.Sleep(1 * time.Second)
	counter.NewRequest()
	if c := counter.TotalRequests(); c != 3 {
		t.Error("It was expected 3 but get", c)
		return
	}
	time.Sleep(1050 * time.Millisecond)
	if c := counter.TotalRequests(); c != 1 {
		t.Error("It was expected 1 but get", c)
		return
	}

	counter.Stop()
}

func TestFlush(t *testing.T) {
	var filename string
	defer remove(createdFile(&filename))

	counter := file.NewCounterService(filename)
	counter.NewRequest()
	counter.NewRequest()
	counter.NewRequest()
	counter.NewRequest()
	counter.NewRequest()
	counter.Stop()

	var requests []time.Time

	f, _ := os.Open(filename)
	enc := gob.NewDecoder(f)
	err := enc.Decode(&requests)
	if err != nil {
		t.Error("Cannot decode list of request:", err)
	}

	if len(requests) != 5 {
		t.Error("It was expected 5 but get", len(requests))
	}
}

func TestStop(t *testing.T) {
	var filename string
	defer remove(createdFile(&filename))

	counter := file.NewCounterService(filename)
	counter.Start()
	time.Sleep(1 * time.Second)

	wait := make(chan struct{})
	go func() {
		counter.Stop()
		close(wait)
	}()

	select {
	case <-wait:
	case <-time.After(time.Second * 3):
		t.Error("Worker took so long to stop")
	}
}
