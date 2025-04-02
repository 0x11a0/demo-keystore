// @title Crypto Vault API
// @version 1.0
// @description This is a secure microservice to store and retrieve encrypted private keys.
// @host localhost:8080
// @BasePath /

package main // Main package for the executable program

import (
	"encoding/json"                 // Standard lib for JSON encoding/decoding
	"fmt"                           // For printing to stdout
	"net/http"                      // HTTP server and routing
	"sync"                          // Provides mutex for concurrency safety

	"crypto-vault/crypto"           // Local module for encryption/decryption logic
	_ "crypto-vault/docs"           // Import generated Swagger docs (used for side effects only)
	httpSwagger "github.com/swaggo/http-swagger" // Handler to serve Swagger UI
)

// vault is an in-memory "database" with a thread-safe lock to prevent race conditions
var vault = struct {
	sync.RWMutex           // Read/Write lock for safe concurrent access
	data map[string]string // Map of username to encrypted private key
}{data: make(map[string]string)} // Initialize the map when the app starts

// StoreRequest defines the expected JSON structure when a user wants to store a key
type StoreRequest struct {
	Username   string `json:"username"`     // Username of the user
	PrivateKey string `json:"private_key"`  // The private key to encrypt and store
	Password   string `json:"password"`     // The password used for encryption
}

// RetrieveRequest defines the expected JSON structure when a user wants to retrieve their key
type RetrieveRequest struct {
	Username string `json:"username"` // Username of the user
	Password string `json:"password"` // Password used to decrypt the stored private key
}

// storeKey handles HTTP POST requests to /store
// @Summary Store encrypted private key
// @Description Encrypt and store a private key for a user
// @Accept json
// @Produce plain
// @Param request body StoreRequest true "Store Request"
// @Success 200 {string} string "Stored successfully"
// @Failure 500 {string} string "Encryption failed"
// @Router /store [post]
func storeKey(w http.ResponseWriter, r *http.Request) {
	var req StoreRequest // Create an empty StoreRequest to fill from the body

	json.NewDecoder(r.Body).Decode(&req) // Decode JSON body into req struct

	// Encrypt the private key using the user's password
	encrypted, err := crypto.Encrypt(req.PrivateKey, req.Password)
	if err != nil {
		// If encryption fails, return 500 error to the client
		http.Error(w, "Encryption failed", http.StatusInternalServerError)
		return
	}

	vault.Lock()                         // Lock for safe write
	vault.data[req.Username] = encrypted // Store encrypted key in memory
	vault.Unlock()                       // Unlock after writing

	w.Write([]byte("Stored successfully")) // Send confirmation response
}

// retrieveKey handles HTTP POST requests to /retrieve
// @Summary Retrieve decrypted private key
// @Description Decrypt and return a stored private key
// @Accept json
// @Produce plain
// @Param request body RetrieveRequest true "Retrieve Request (only username and password are required)"
// @Success 200 {string} string "Decrypted key"
// @Failure 404 {string} string "No such user"
// @Failure 401 {string} string "Decryption failed"
// @Router /retrieve [post]
func retrieveKey(w http.ResponseWriter, r *http.Request) {
	var req RetrieveRequest // Create an empty RetrieveRequest to fill from body

	json.NewDecoder(r.Body).Decode(&req) // Decode incoming JSON into struct

	vault.RLock()                        // Lock for safe read
	encrypted, ok := vault.data[req.Username] // Try to get encrypted key from memory
	vault.RUnlock()                      // Unlock after reading

	if !ok {
		// If user is not found, return 404
		http.Error(w, "No such user", http.StatusNotFound)
		return
	}

	// Attempt to decrypt using the password
	decrypted, err := crypto.Decrypt(encrypted, req.Password)
	if err != nil {
		// If decryption fails (wrong password or tampered data), return 401
		http.Error(w, "Decryption failed", http.StatusUnauthorized)
		return
	}

	w.Write([]byte(decrypted)) // Return the decrypted private key
}

// main is the entry point of the application
func main() {
	http.HandleFunc("/store", storeKey)     // Route for storing keys
	http.HandleFunc("/retrieve", retrieveKey) // Route for retrieving keys
	http.Handle("/swagger/", httpSwagger.WrapHandler) // Route for serving Swagger UI

	fmt.Println("Server running on :8080")  // Print startup message
	http.ListenAndServe(":8080", nil)       // Start HTTP server on port 8080
}
