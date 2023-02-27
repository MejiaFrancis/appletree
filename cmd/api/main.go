// Filename : cmd/api/main.go
package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
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
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
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
	flag.StringVar(&cfg.db.dsn, "db_dsn", os.Getenv("APPLETREE_DB_DSN"), "PostgreSQL DSN")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")
	// flag.StringVar(&cfg.db.dsn, "dsn", os.Getenv("$MELONS_DB_DSN"), "PostgreSQL DSN") // this line is used to  allow communication between DB and API
	flag.Parse()
	// Create a logger (customised)
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	// Create the connection pool
	db, err := openDB(cfg)
	if err != nil {
		logger.Fatal(err)
	}
    // close the database connection pool
	defer db.Close()
	// Log the successful connection pool
	logger.Println("database connection pool establishes")
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
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,      // inactive connections
		ReadTimeout:  10 * time.Second, // time to read response body or header
		WriteTimeout: 30 * time.Second, // time to write response body or header
	}
	// Start the customer server
	logger.Printf("starting %s server on %s", cfg.env, srv.Addr)
	err = srv.ListenAndServe()
	logger.Fatal(err)
}

// The openDB() function returns a  *sql.DB connection pool
func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("pgx", cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)
	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)
	// Create a context with a 5-seconnd timeout deadline
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
    // ping the database
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}
