package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	function "todoAPI/api/function"
	"github.com/gorilla/mux"
)

const port string = ":9010"

func main() {
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt)

	r := mux.NewRouter()
	r.HandleFunc("/", function.HomeHandler).Methods("GET")
	r.HandleFunc("/todo", function.FetchTodos).Methods("GET")
	r.HandleFunc("/todo", function.CreateTodo).Methods("POST")
	r.HandleFunc("/todo/{id}", function.UpdateTodo).Methods("PUT")
	r.HandleFunc("/todo/{id}", function.DeleteTodo).Methods("DELETE")

	srv := &http.Server{
		Addr:         port,
		Handler:      r,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Println("Listening on port", port)
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("Server error: %s\n", err)
		}
	}()

	<-stopChan
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if function.CheckErr(srv.Shutdown(ctx), "Server shutdown error") {
		return
	}
	log.Println("Server gracefully stopped!")
}
