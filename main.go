package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"dahbura.me/api/config"
	"dahbura.me/api/middleware"
	"dahbura.me/api/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("> ")

	// setup env
	config.Load()

	// blank engine
	router := gin.New()

	// middleware
	router.Use(middleware.Logger())
	router.Use(middleware.Recovery())
	router.Use(middleware.Cors())

	// routes
	routes.Startup(router)

	// server config
	addr := fmt.Sprintf("%s:%s", config.Host, config.Port)

	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		IdleTimeout:  config.DefaultIdleTimeout,
		ReadTimeout:  config.DefaultReadTimeout,
		WriteTimeout: config.DefaultWriteTimeout,
	}

	go func() {
		log.Printf("Starting server on %s", config.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen: %s\n", err)
		}
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	<-done

	log.Println("Stopping server...")

	ctx, cancel := context.WithTimeout(context.Background(), config.DefaultCtxTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Print("Server forced to shutdown: ", err)
	}
}
