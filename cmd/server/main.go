package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ognev-dev/goplease/ws"
)

const (
	// TODO proper config
	RWTimeout = 10 * time.Second
	Host      = "127.0.0.1"
	Port      = "8090"
)

func main() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit,
		syscall.SIGTERM,
		syscall.SIGHUP,  // kill -SIGHUP
		syscall.SIGINT,  // kill -SIGINT or Ctrl+c
		syscall.SIGQUIT, // kill -SIGQUIT
	)

	hub := ws.NewHub()
	gs := ws.NewGameServer(hub)

	go hub.Run()
	go gs.Run()

	mux := http.NewServeMux()
	mux.HandleFunc("/play/", hub.ServeWS)

	addr := net.JoinHostPort(Host, Port)
	server := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  RWTimeout,
		WriteTimeout: RWTimeout,
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Printf("[goplease] server running at %s", addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrors <- err
		}
	}()

	select {
	case err := <-serverErrors:
		log.Fatalf("[goplease] server error: %v", err)

	case sig := <-quit:
		log.Printf("[goplease] %v signal received...", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			_ = server.Close()
		}

		log.Println("[goplease] bye")
	}
}
