package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/thedevsaddam/renderer"
)

var rnd *renderer.Render
var todos = []todo{}

const port string = ":9010"

type (
	todo struct {
		ID        string    `json:"id"`
		Title     string    `json:"title"`
		Completed bool      `json:"completed"`
		CreatedAt time.Time `json:"created_at"`
	}
)

func init() {
	rnd = renderer.New()
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	err := rnd.Template(w, http.StatusOK, []string{"templete/home.tpl"}, nil)
	checkErr(err)
}

func createTodo(w http.ResponseWriter, r *http.Request) {
	var t todo

	// Decode the JSON request body
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		rnd.JSON(w, http.StatusBadRequest, renderer.M{"error": "Invalid request"})
		return
	}

	// Validate that the title is not empty
	if t.Title == "" {
		rnd.JSON(w, http.StatusBadRequest, renderer.M{"message": "The title field is required"})
		return
	}

	// Generate a unique ID and timestamp for the new Todo
	t.ID = time.Now().Format("20060102150405")
	t.CreatedAt = time.Now()
	todos = append(todos, t)

	rnd.JSON(w, http.StatusCreated, renderer.M{
		"message": "Todo created successfully",
		"todo_id": t.ID,
	})
}

func updateTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := strings.TrimSpace(vars["id"])
	var t todo

	// Decode the JSON request body
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		rnd.JSON(w, http.StatusBadRequest, renderer.M{"error": "Invalid request"})
		return
	}

	// Find and update the Todo item
	for i, item := range todos {
		if item.ID == id {
			todos[i].Title = t.Title
			todos[i].Completed = t.Completed
			rnd.JSON(w, http.StatusOK, renderer.M{"message": "Todo updated successfully"})
			return
		}
	}

	rnd.JSON(w, http.StatusNotFound, renderer.M{"message": "Todo not found"})
}

// Fetches and returns all Todo items
func fetchTodos(w http.ResponseWriter, r *http.Request) {
	rnd.JSON(w, http.StatusOK, renderer.M{"data": todos})
}

func deleteTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := strings.TrimSpace(vars["id"])

	// Find and remove the Todo item
	for i, item := range todos {
		if item.ID == id {
			todos = append(todos[:i], todos[i+1:]...)
			rnd.JSON(w, http.StatusOK, renderer.M{"message": "Todo deleted successfully"})
			return
		}
	}

	rnd.JSON(w, http.StatusNotFound, renderer.M{"message": "Todo not found"})
}

func main() {
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt)

	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler).Methods("GET")
	r.HandleFunc("/todo", fetchTodos).Methods("GET")
	r.HandleFunc("/todo", createTodo).Methods("POST")
	r.HandleFunc("/todo/{id}", updateTodo).Methods("PUT")
	r.HandleFunc("/todo/{id}", deleteTodo).Methods("DELETE")

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
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %s\n", err)
	}
	log.Println("Server gracefully stopped!")
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err) //respond with error page or message
	}
}
