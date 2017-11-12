package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/guilherme-santos/simplesurance/file"
	"github.com/guilherme-santos/simplesurance/http"
)

func getEnvVar(envvar string) string {
	value := os.Getenv(envvar)
	if value == "" {
		log.Fatalf("Environment variable \"%s\" is empty or missing\n", envvar)
	}

	return value
}

func main() {
	apiPort := getEnvVar("SIMPLEINSURANCE_API_PORT")
	filename := getEnvVar("SIMPLEINSURANCE_API_COUNTER_FILENAME")

	// Using file implementation
	counterService := file.NewCounterService(filename)
	// But you could also use:
	// counterService := mysql.NewCounterService(filename)
	// counterService := sqlite.NewCounterService(filename)
	// counterService := mongodb.NewCounterService(filename)

	// Run worker
	counterService.Start()
	defer counterService.Stop()

	counterHandler := http.NewCounterHandler(counterService)
	router := http.NewRouter(counterHandler)

	go func() {
		err := router.Run(apiPort)
		if err != nil {
			log.Fatalln("Cannot run webserver:", err)
		}
	}()

	// If you CONTROL+C defer functions will not be called, to fix that
	// I'm dealing here with this signals to gracefully shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	sig := <-sigChan

	log.Printf("Shutting down, %v signal received\n", sig)
}
