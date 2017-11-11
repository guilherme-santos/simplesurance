package main

import (
	"log"
	"os"
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

	counterHandler := http.NewCounterHandler(counterService)
	router := http.NewRouter(counterHandler)

	err := router.Run(apiPort)
	if err != nil {
		log.Fatal("Cannot run webserver: ", err.Error())
	}
}
