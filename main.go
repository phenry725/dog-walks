package main

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

var (
	isReady  bool
	initOnce sync.Once
	port     string
)

const ()

func init() {
	var set bool
	var err error
	initOnce.Do(func() {
		envFileNameArg := os.Args[1]
		err = godotenv.Load(envFileNameArg)
		if err != nil {
			log.Fatalf("Error loading env file named: %v", envFileNameArg)
		}

		port, set = os.LookupEnv("port")
		if !set {
			port = "80"
			log.Printf("No port set, defaulting to %v", port)
		}

		isReady = true
	})
}

func main() {

	log.Println("hello")
}
