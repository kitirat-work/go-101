package main

import (
	"context"
	exportlargeexcel "export-large-excel"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	ctx := context.Background()
	r := mux.NewRouter()

	store := exportlargeexcel.NewStore()
	service := exportlargeexcel.NewService(store)
	handler := exportlargeexcel.NewHandler(service)

	r.HandleFunc("/export", handler.Export).Methods(http.MethodGet)

	srv := &http.Server{
		Handler:      r,
		Addr:         "0.0.0.0:8000",
		WriteTimeout: 5 * time.Minute,
		ReadTimeout:  5 * time.Minute,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("fail %v!!", err)
		}
	}()
	log.Println("service is ready...")

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	<-c
	srv.Shutdown(ctx)
	log.Println("shutting down...")
	os.Exit(0)
}
