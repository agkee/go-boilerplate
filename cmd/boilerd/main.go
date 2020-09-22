package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/agkee/go-boilerplate.git/internal/web"
)

var (
	addr  = flag.String("addr", "0.0.0.0:9090", "Default HTTP address")
	debug = flag.Bool("debug", false, "Set logging level to debug")
	dsn   = flag.String("dsn", "", "Set DSN for sentry alerting")

	writeTimeout = time.Second * 30
	readTimeout  = time.Second * 30
)

func main() {
	r := web.NewRouter()
	// r = observe.RegisterPrometheus(r)

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
