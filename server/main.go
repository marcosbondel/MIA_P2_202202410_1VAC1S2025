package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// Book struct (Model)
type Book struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Author *Author `json:"author"`
}

// Author struct
type Author struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

// Init books var as a slice Book struct
var books []Book

// Get all books
func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

// Get single book
func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // Get params
	// Loop through books and find with id
	for _, item := range books {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Book{})
}

// Create a new book
func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	book.ID = strconv.Itoa(len(books) + 1) // Mock ID - not safe for production
	books = append(books, book)
	json.NewEncoder(w).Encode(book)
}

// Update book
func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range books {
		if item.ID == params["id"] {
			books = append(books[:index], books[index+1:]...)
			var book Book
			_ = json.NewDecoder(r.Body).Decode(&book)
			book.ID = params["id"]
			books = append(books, book)
			json.NewEncoder(w).Encode(book)
			return
		}
	}
	json.NewEncoder(w).Encode(books)
}

func getRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("Hello world!")
}

func login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// For simplicity, we are not implementing actual login logic here.
	// In a real application, you would validate user credentials and return a token.
	fmt.Println("Login endpoint hit")
	json.NewEncoder(w).Encode("Login endpoint - implement authentication logic here")
}

func main() {
	// Initialize router
	r := mux.NewRouter()

	// Mock data
	books = append(books, Book{ID: "1", Title: "Book One", Author: &Author{Firstname: "John", Lastname: "Doe"}})
	books = append(books, Book{ID: "2", Title: "Book Two", Author: &Author{Firstname: "Steve", Lastname: "Smith"}})

	// Route handlers / Endpoints
	r.HandleFunc("/", getRoot).Methods("GET")
	r.HandleFunc("/api/auth/login", login).Methods("POST")
	r.HandleFunc("/api/books", getBooks).Methods("GET")
	r.HandleFunc("/api/books/{id}", getBook).Methods("GET")
	r.HandleFunc("/api/books", createBook).Methods("POST")
	r.HandleFunc("/api/books/{id}", updateBook).Methods("PUT")

	// CORS setup
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	// Start server with CORS
	fmt.Println("Server started on port 3000")
	log.Fatal(http.ListenAndServe(":3000", handlers.CORS(originsOk, headersOk, methodsOk)(r)))
}
