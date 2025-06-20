package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"MIA_P2_202202410_1VAC1S2025/fs/user"
	"MIA_P2_202202410_1VAC1S2025/models"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func getRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("Hello world!")
}

func login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"JSON inválido"}`, http.StatusBadRequest)
		return
	}

	if user.Login(req.User, req.Pass, req.Id) {
		// Código 200 OK
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(models.LoginResponse{
			Status:  "success",
			Message: "Successfully logged in.",
		})
	} else {
		// Código 401 Unauthorized
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(models.LoginResponse{
			Status:  "fail",
			Message: "Invalid credentials or partition not mounted.",
		})
	}
}
func getDisks(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Fetching disks...")
	w.Header().Set("Content-Type", "application/json")

	// Aquí podrías implementar la lógica para obtener los discos
	// Por ahora, simplemente devolvemos un mensaje de ejemplo
	disks := []string{"A", "B", "C"}

	json.NewEncoder(w).Encode(disks)
}

func main() {
	// Initialize router
	r := mux.NewRouter()

	// Route handlers / Endpoints
	r.HandleFunc("/", getRoot).Methods("GET")
	r.HandleFunc("/api/disks", getDisks).Methods("GET")
	r.HandleFunc("/api/auth/login", login).Methods("POST")

	// CORS setup
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	// Start server with CORS
	fmt.Println("Server started on port 3000")
	log.Fatal(http.ListenAndServe(":3000", handlers.CORS(originsOk, headersOk, methodsOk)(r)))
}
