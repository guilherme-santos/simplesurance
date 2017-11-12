package file

import (
	"encoding/gob"
	"io"
	"log"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

type CounterService struct {
	requests     []time.Time
	dbfile       *os.File
	mutex        sync.Mutex
	wokerStarted uint32
	closeChan    chan bool
}

var (
	WindowSize = 60 * time.Second
	FlushDisk  = 30 * time.Second
)

func NewCounterService(filename string) *CounterService {
	service := &CounterService{
		requests:  make([]time.Time, 0, 60),
		closeChan: make(chan bool),
	}

	var err error
	service.dbfile, err = os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Fatalf("Cannot open \"%s\": %s\n", filename, err)
	}

	service.Load()

	return service
}

func Un(f func()) {
	f()
}

func Lock(x sync.Locker) func() {
	x.Lock()
	return func() { x.Unlock() }
}

func (s *CounterService) TotalRequests() int {
	defer Un(Lock(&s.mutex))
	return len(s.requests)
}

func (s *CounterService) NewRequest() int {
	req := time.Now()
	s.mutex.Lock()
	s.requests = append(s.requests, req)
	s.mutex.Unlock()
	log.Println("New request:", req.Format(time.RFC3339Nano))

	return s.TotalRequests()
}

func (s *CounterService) Load() {
	log.Println("Loading data from disk...")

	enc := gob.NewDecoder(s.dbfile)
	s.mutex.Lock()
	err := enc.Decode(&s.requests)
	s.mutex.Unlock()

	if err != nil && err != io.EOF {
		log.Println("Error decoding data:", err)
		return
	}
}

func (s *CounterService) Flush() {
	log.Println("Flushing data to disk...")
	// Clear file before write
	s.dbfile.Truncate(0)

	enc := gob.NewEncoder(s.dbfile)
	s.mutex.Lock()
	err := enc.Encode(s.requests)
	s.mutex.Unlock()

	if err != nil {
		log.Println("Error encoding data:", err)
		return
	}

	err = s.dbfile.Sync()
	if err != nil {
		log.Println("Error flushing:", err)
		return
	}
}

func (s *CounterService) Start() {
	// Ensure to run just once
	if atomic.LoadUint32(&s.wokerStarted) == 1 {
		return
	}

	atomic.StoreUint32(&s.wokerStarted, 1)
	log.Println("Starting CounterService worker...")

	go func() {
		tickerToSync := time.NewTicker(FlushDisk).C
		tickerToRemoveExpired := time.NewTicker(1 * time.Second).C

		for {
			select {
			case <-s.closeChan:
				close(s.closeChan)
				return
			case <-tickerToSync:
				go s.Flush()
			case <-tickerToRemoveExpired:
				if s.TotalRequests() == 0 {
					continue
				}

				expired := time.Now().Add(-1 * WindowSize)
				s.mutex.Lock()
				for _, req := range s.requests {
					// First request that is not expired, all others for sure don't too
					// because older requests are first in the array (queue FIFO)
					if !req.Before(expired) {
						break
					}

					s.requests = s.requests[1:]
					log.Println("Request removed:", req.Format(time.RFC3339Nano), "missing", len(s.requests))
				}
				s.mutex.Unlock()
			}
		}
	}()
}

func (s *CounterService) Stop() {
	// Stop worker just if it's running
	if atomic.LoadUint32(&s.wokerStarted) == 1 {
		defer atomic.StoreUint32(&s.wokerStarted, 0)
		log.Println("Stopping CounterService worker...")
		s.closeChan <- true
		<-s.closeChan
		log.Println("Worker stopped!")
	}

	// Flush and close dbfile
	s.Flush()

	err := s.dbfile.Close()
	if err != nil {
		log.Println("Error closing dbfile:", err)
	}
}
