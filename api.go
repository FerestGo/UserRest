package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

func StartServer(router *mux.Router, config Config) {
	// start http server
	log.Printf("starting HTTP server")
	srv := listenAndServe(fmt.Sprintf(":%s", config["APP_PORT"]), router)

	// Setting up signal capturing
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// Waiting for SIGINT (pkill -2)
	<-stop

	log.Printf("stopping server")
	// graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		panic(err)
	}
	log.Printf("exit")
}

func listenAndServe(port string, router *mux.Router) *http.Server {
	srv := &http.Server{Addr: port, Handler: router}
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %s", err)
		}
	}()
	return srv
}
