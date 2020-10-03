package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/agkee/go-boilerplate.git/internal/db"
	"github.com/agkee/go-boilerplate.git/internal/observe"
	"github.com/agkee/go-boilerplate.git/internal/web"
	"github.com/go-sql-driver/mysql"
)

var (
	addr         = flag.String("addr", "0.0.0.0:9090", "Default HTTP address")
	debug        = flag.Bool("debug", false, "Set logging level to debug")
	dsn          = flag.String("dsn", "your_DSN", "Set DSN for sentry alerting")
	dbUser       = flag.String("db_user", "root", "The user of the database")
	dbHost       = flag.String("db_host", "0.0.0.0", "The host of the database")
	dbName       = flag.String("db_name", "cafe_db", "The name of the database to connect to")
	dbPass       = flag.String("db_pass", "1234", "The password of the database to connect to")
	writeTimeout = time.Second * 30
	readTimeout  = time.Second * 30
)

func main() {
	flag.Parse()
	observe.InitLogging(*debug, *dsn)

	dbConfig := mysql.Config{
		User:   *dbUser,
		DBName: *dbName,
		Passwd: *dbPass,
		Addr:   *dbHost,
		Net:    "tcp",
	}

	dbDSN := dbConfig.FormatDSN()
	_, _ = db.GetDB(dbDSN)

	r := web.NewRouter()
	r = observe.RegisterPrometheus(r)

	// Setting timeouts and handlers for http server
	s := &http.Server{
		Addr:         *addr,
		Handler:      r,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}

	// Without graceful shutdown, the server might shutdown while handling important requests
	go func() {
		log.Printf("Running server on port %s", *addr)
		err := s.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal)
	// Notify when there is a os interrupt/kill command
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	// Block the channel here, waiting to receive that os.Interrupt or os.Kill
	sig := <-sigChan
	log.Println("Received terminate signal, gracefully shutting down", sig)

	// Wait for 30 seconds for all handlers to finish
	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)

	// Shut down after everything is over
	s.Shutdown(tc)
}
