// Filename : cmd/api/main.go
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// A global variable to hold the application version number
// The application version number
const version = "1.0.0"

// We want to customize our server using command line flags.
// Therefore, we will store the settings in a config struct
// The configuration settings
type config struct {
	port int
	env  string // development, staging, production, etc.
}

// We will use dependency injection. We will create a type that will hold
// the dependencies from our HTTP handlers, helpers, and middleware
// Dependency Injection
type application struct {
	config config
	logger *log.Logger
}

func main() {
	var cfg config
	// Get the commandline arguments and store them in the config struct
	// read in the flags that are needed to populate our config
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development | staging | production")
	flag.Parse()
	// Create a logger (customised)
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	// Create an instance of our application struct
	// dependencies available
	app := &application{
		config: cfg,
		logger: logger,
	}
	// Create our new servemux (multiplexer)
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/healthcheck", app.healthcheckHandler) //register our routes
	// Create our HTTP(customised) server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      mux,
		IdleTimeout:  time.Minute,      // inactive connections
		ReadTimeout:  10 * time.Second, // time to read response body or header
		WriteTimeout: 30 * time.Second, // time to write response body or header
	}
	// Start the customer server
	logger.Printf("starting %s server on %s", cfg.env, srv.Addr)
	err := srv.ListenAndServe()
	logger.Fatal(err)

}
