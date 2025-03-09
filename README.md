# Todo API - Golang

## Overview
This project is a simple **Todo API** built using **Golang**, featuring CRUD operations to manage todo items. It leverages **Gorilla Mux** for routing and **TheDevSaddam Renderer** for JSON responses and templating. The API is designed to be minimal yet functional, supporting essential features like creating, updating, retrieving, and deleting todo items. Additionally, it includes graceful server shutdown handling and basic input validation.

---

## Features & Functionalities

### 1️⃣ Home Page Rendering
- The API includes a basic homepage using an HTML template (`home.tpl`).
- The `homeHandler` function renders this template when a request is made to `/`.

### 2️⃣ CRUD Operations
The API provides full **CRUD (Create, Read, Update, Delete)** functionality:

#### ✅ Create a Todo
- **Endpoint:** `POST /todo`
- **Function:** `createTodo`
- **Request Body:** JSON containing `title` (string).
- **Response:** Returns the created todo item with a unique `id` and `created_at` timestamp.
- **Validation:** The `title` field is required.

#### ✅ Fetch All Todos
- **Endpoint:** `GET /todo`
- **Function:** `fetchTodos`
- **Response:** Returns all todo items as a JSON array.

#### ✅ Update a Todo
- **Endpoint:** `PUT /todo/{id}`
- **Function:** `updateTodo`
- **Request Body:** JSON containing updated `title` and `completed` status.
- **Response:** Updates the todo item if found; returns an error if the ID is invalid.

#### ✅ Delete a Todo
- **Endpoint:** `DELETE /todo/{id}`
- **Function:** `deleteTodo`
- **Response:** Removes the specified todo item if found; returns an error if not found.

### 3️⃣ Error Handling
- Ensures that **bad requests** (e.g., missing title, incorrect JSON format) return meaningful error messages.
- Uses `http.StatusBadRequest` for invalid requests and `http.StatusNotFound` for missing todos.

### 4️⃣ Graceful Shutdown Handling
- The server listens for `os.Interrupt` signals and safely shuts down when triggered.
- Uses `context.WithTimeout()` to allow the server to finish ongoing requests before stopping.

---

## How It Works (Technical Breakdown)

### 1️⃣ Setting Up the Renderer
- The `renderer.Render` instance (`rnd`) is initialized in the `init()` function.
- It is used to render templates (`home.tpl`) and return JSON responses.

### 2️⃣ Handling API Routes
- Uses `mux.NewRouter()` to define API endpoints.
- Registers handlers for different HTTP methods (GET, POST, PUT, DELETE).

### 3️⃣ Storing and Managing Todos
- Uses an **in-memory slice (`[]todo`)** to store todo items.
- Each todo item contains:
  - `ID` (string) – Unique identifier based on timestamp.
  - `Title` (string) – Task description.
  - `Completed` (bool) – Task completion status.
  - `CreatedAt` (time.Time) – Timestamp when the task was created.

### 4️⃣ Running the HTTP Server
- The server listens on **port 9010**.
- Uses `http.Server` with timeouts for better performance and security.
- Logs server start-up and shutdown events.

---

## How to Run This Project

### **Prerequisites**
Ensure you have the following installed:
- **Golang** (1.16 or later)
- **Git** (for version control)

### **Clone the Repository**
```sh
git clone <repository-url>
cd <project-folder>
```

### **Install Dependencies**
```sh
go mod tidy
```

### **Run the Server**
```sh
go run main.go
```

The server will start on **`http://localhost:9010`**.

### **Test API Endpoints**
You can use **Postman** or `curl` to test the endpoints:

#### **Create a Todo**
```sh
curl -X POST "http://localhost:9010/todo" -H "Content-Type: application/json" -d '{"title":"Learn Golang"}'
```

#### **Fetch All Todos**
```sh
curl -X GET "http://localhost:9010/todo"
```

#### **Update a Todo**
```sh
curl -X PUT "http://localhost:9010/todo/{id}" -H "Content-Type: application/json" -d '{"title":"Updated Task", "completed":true}'
```

#### **Delete a Todo**
```sh
curl -X DELETE "http://localhost:9010/todo/{id}"
```

---


