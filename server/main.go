package main

import (
	"encoding/json"
	"log"
	"net/http"
	"path/filepath"
	
	"wa-1/game"
)

func main() {
	// 1. Load Dictionary
	dict, err := game.NewDictionary(filepath.Join("data", "places.json"))
	if err != nil {
		log.Fatalf("Failed to load dictionary: %v", err)
	}
	log.Println("Dictionary loaded.")

	// 2. Setup User Manager
	um, err := game.NewUserManager(filepath.Join("data", "users.json"))
	if err != nil {
		log.Fatalf("Failed to load user manager: %v", err)
	}
	log.Println("User Manager loaded.")

	// 3. Setup Game Manager
	manager := game.NewManager(dict, um)

	// 4. Setup Routes
	http.HandleFunc("/ws", manager.HandleWS)
	http.HandleFunc("/api/register", handleRegister(um))
	http.HandleFunc("/api/login", handleLogin(um))
	
	// Serve Frontend (Vue build)
	fs := http.FileServer(http.Dir("../client/dist"))
	http.Handle("/", fs)

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func handleRegister(um *game.UserManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		if r.Method == "OPTIONS" { return }
		
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req AuthRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		user, err := um.Register(req.Username, req.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}
}

func handleLogin(um *game.UserManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		if r.Method == "OPTIONS" { return }

		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req AuthRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		user, err := um.Login(req.Username, req.Password)
		if err != nil {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}