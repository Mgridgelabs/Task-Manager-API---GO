package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// Task represents a task structure
type Task struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Details string `json:"details"`
}

// Seed data
var tasks = []Task{
	{ID: 1, Title: "Task 1", Details: "First task"},
	{ID: 2, Title: "Task 2", Details: "Second task"},
	{ID: 3, Title: "Task 3", Details: "Third task"},
}

// Handlers
func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, World!")
}

func getTasksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func getTaskByIDHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	for _, task := range tasks {
		if task.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(task)
			return
		}
	}
	http.Error(w, "Task not found", http.StatusNotFound)
}

func createTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
		return
	}
	var newTask Task
	if err := json.NewDecoder(r.Body).Decode(&newTask); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	newTask.ID = len(tasks) + 1 // Auto-increment ID
	tasks = append(tasks, newTask)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTask)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
		return
	}
	var updatedTask Task
	if err := json.NewDecoder(r.Body).Decode(&updatedTask); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	for i, task := range tasks {
		if task.ID == updatedTask.ID {
			tasks[i] = updatedTask
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(tasks[i])
			return
		}
	}
	http.Error(w, "Task not found", http.StatusNotFound)
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}
	for i, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...) // Remove task from slice
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.Error(w, "Task not found", http.StatusNotFound)
}

func main() {
	// Routes
	http.HandleFunc("/", helloHandler)
	http.HandleFunc("/tasks", getTasksHandler)      // Get all tasks
	http.HandleFunc("/task", getTaskByIDHandler)    // Get by ID
	http.HandleFunc("/create", createTaskHandler)   // POST to create a new task
	http.HandleFunc("/update", updateTaskHandler)   // PUT to update a task
	http.HandleFunc("/delete", deleteTaskHandler)   // DELETE a task

	// Start server
	fmt.Println("Server running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
