package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/thedevsaddam/renderer"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var rnd *renderer.Render
var client *mongo.Client
var todoCollection *mongo.Collection

type Todo struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title     string             `json:"title" bson:"title"`
	Completed bool               `json:"completed" bson:"completed"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}

func init() {
	rnd = renderer.New()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	CheckErr(err, "Error connecting to MongoDB")
	todoCollection = client.Database("tododb").Collection("todos")
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	err := rnd.Template(w, http.StatusOK, []string{"templete/home.tpl"}, nil)
	CheckErr(err, "Failed")
}

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	var t Todo
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		rnd.JSON(w, http.StatusBadRequest, renderer.M{"error": "Invalid request"})
		return
	}

	if t.Title == "" {
		rnd.JSON(w, http.StatusBadRequest, renderer.M{"message": "The title field is required"})
		return
	}

	t.CreatedAt = time.Now()
	res, err := todoCollection.InsertOne(context.TODO(), t)
	if CheckErr(err, "Failed to create todo") {
		return
	}

	rnd.JSON(w, http.StatusCreated, renderer.M{"message": "Todo created successfully", "todo_id": res.InsertedID})
}

func FetchTodos(w http.ResponseWriter, r *http.Request) {
	cursor, err := todoCollection.Find(context.TODO(), bson.M{})
	if CheckErr(err, "Failed to fetch todos") {
		return
	}
	defer cursor.Close(context.TODO())

	var todos []Todo
	if CheckErr(cursor.All(context.TODO(), &todos), "Failed to parse todos") {
		return
	}

	rnd.JSON(w, http.StatusOK, renderer.M{"data": todos})
}

func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := strings.TrimSpace(vars["id"])
	objID, err := primitive.ObjectIDFromHex(id)
	if CheckErr(err, "Invalid ID") {
		rnd.JSON(w, http.StatusBadRequest, renderer.M{"error": "Invalid ID"})
		return
	}

	var t Todo
	if CheckErr(json.NewDecoder(r.Body).Decode(&t), "Invalid request") {
		rnd.JSON(w, http.StatusBadRequest, renderer.M{"error": "Invalid request"})
		return
	}

	update := bson.M{"$set": bson.M{"title": t.Title, "completed": t.Completed}}
	res, err := todoCollection.UpdateOne(context.TODO(), bson.M{"_id": objID}, update)
	if CheckErr(err, "Failed to update todo") || res.MatchedCount == 0 {
		rnd.JSON(w, http.StatusNotFound, renderer.M{"message": "Todo not found"})
		return
	}

	rnd.JSON(w, http.StatusOK, renderer.M{"message": "Todo updated successfully"})
}

func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := strings.TrimSpace(vars["id"])
	objID, err := primitive.ObjectIDFromHex(id)
	if CheckErr(err, "Invalid ID") {
		rnd.JSON(w, http.StatusBadRequest, renderer.M{"error": "Invalid ID"})
		return
	}

	res, err := todoCollection.DeleteOne(context.TODO(), bson.M{"_id": objID})
	if CheckErr(err, "Failed to delete todo") || res.DeletedCount == 0 {
		rnd.JSON(w, http.StatusNotFound, renderer.M{"message": "Todo not found"})
		return
	}

	rnd.JSON(w, http.StatusOK, renderer.M{"message": "Todo deleted successfully"})
}

func CheckErr(err error, msg string) bool {
	if err != nil {
		log.Println(msg, ":", err)
		return true
	}
	return false
}
