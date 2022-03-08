package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"github.com/3n0ugh/BasedWeb/internal/data"
	"github.com/3n0ugh/BasedWeb/internal/jsonlog"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type config struct {
	port int
	db   struct {
		dsn string
	}
}

type application struct {
	config config
	logger *jsonlog.Logger
	model  data.Model
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 8080, "API server port")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "", "PostgreSQL DSN")
	flag.Parse()

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	db, err := openDB(cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}
	defer func() {
		if err := db.Close(); err != nil {
			panic(err)
		}
	}()

	logger.PrintInfo("database connection pool established", nil)

	app := &application{
		config: cfg,
		logger: logger,
		model:  data.NewModel(db),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.port),
		Handler: app.routes(),
	}
	logger.PrintFatal(srv.ListenAndServe(), nil)
}

func openDB(cfg config) (*sql.DB, error) {

	// Create an empty connection pool
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	// Create a context with a 5-second timeout deadline
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// If the connection couldn't be established successfully
	// within the 5-second deadline, then this will return an error
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
