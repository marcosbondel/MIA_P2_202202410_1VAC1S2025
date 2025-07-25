package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"MIA_P2_202202410_1VAC1S2025/fs/analyzer"
	"MIA_P2_202202410_1VAC1S2025/fs/file_system"
	"MIA_P2_202202410_1VAC1S2025/fs/user"
	"MIA_P2_202202410_1VAC1S2025/models"

	"io"
	"os"
	"strings"

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

	var buffer_string string = ""

	if user.Login(req.User, req.Pass, req.Id, &buffer_string) {
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

func logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var buffer_string string = ""

	if !user.Logout(&buffer_string) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.LoginResponse{
			Status:  "fail",
			Message: "Failed to log out. User not logged in or partition not mounted.",
		})
	} else {
		// Código 200 OK
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(models.LoginResponse{
			Status:  "success",
			Message: "Successfully logged out.",
		})
	}

	// json.NewEncoder(w).Encode(models.LoginResponse{
	// 	Status:  "success",
	// 	Message: "Successfully logged out.",
	// })
}

func doExecute(w http.ResponseWriter, r *http.Request) {
	// Aquí podrías implementar la lógica para ejecutar comandos
	// Por ahora, simplemente devolvemos un mensaje de ejemplo
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("Executing command...")

	var req models.ExecuteStringRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"JSON inválido"}`, http.StatusBadRequest)
		return
	}

	fmt.Println("Command String:", req.CommandString)
	response_string := analyzer.AnalyzeHTTPInput(req.CommandString)

	// json.NewEncoder(w).Encode("Command executed successfully.")
	json.NewEncoder(w).Encode(response_string)
}

func getDisks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	disks, err := file_system.ListDisks()
	if err != nil {
		http.Error(w, `{"error":"No se pudieron listar los discos"}`, http.StatusInternalServerError)
		return
	}

	response := models.DiskResponse{Disks: disks}
	json.NewEncoder(w).Encode(response)
}

// main.go o handlers/partitions.go
func getDiskPartitions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Obtener letra del disco desde query param: /api/disks/A/partitions
	vars := mux.Vars(r)
	letter := vars["letter"]

	partitions, err := file_system.GetDiskPartitions(letter)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(partitions)
}

func uploadSDAAFile(w http.ResponseWriter, r *http.Request) {
	// Limita tamaño del archivo a 10MB
	r.ParseMultipartForm(10 << 20)

	// Recupera el archivo
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error leyendo el archivo", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Verifica extensión
	if !strings.HasSuffix(handler.Filename, ".sdaa") {
		http.Error(w, "Solo se permiten archivos .sdaa", http.StatusBadRequest)
		return
	}

	// Ruta donde se almacenará
	savePath := "./uploads/" + handler.Filename

	// Crea el archivo en el servidor
	dst, err := os.Create(savePath)
	if err != nil {
		http.Error(w, "No se pudo guardar el archivo", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copia el contenido del archivo recibido al nuevo archivo
	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Error al guardar el archivo", http.StatusInternalServerError)
		return
	}

	fmt.Println("✅ Archivo guardado en:", savePath)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Archivo subido exitosamente",
		"path":    savePath,
	})
}

func main() {
	// Initialize router
	r := mux.NewRouter()

	// Route handlers / Endpoints
	r.HandleFunc("/", getRoot).Methods("GET")
	r.HandleFunc("/api/disks", getDisks).Methods("GET")
	r.HandleFunc("/api/auth/login", login).Methods("POST")
	r.HandleFunc("/api/auth/logout", logout).Methods("POST")
	r.HandleFunc("/api/run_command", doExecute).Methods("POST")
	r.HandleFunc("/api/disks/{letter}/partitions", getDiskPartitions).Methods("GET")
	r.HandleFunc("/api/fs", file_system.GetFileSystem).Methods("GET")
	r.HandleFunc("/api/upload", uploadSDAAFile).Methods("POST")

	// CORS setup
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	// Start server with CORS
	fmt.Println("Server started on port 3000")
	log.Fatal(http.ListenAndServe(":3000", handlers.CORS(originsOk, headersOk, methodsOk)(r)))
}
