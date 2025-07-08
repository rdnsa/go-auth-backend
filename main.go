package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

var users = []User{}

func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	if r.Method == "OPTIONS" {
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Hanya POST yang diperbolehkan", http.StatusMethodNotAllowed)
		return
	}

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Format JSON salah", http.StatusBadRequest)
		return
	}

	// Simpan user ke array
	users = append(users, user)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Register berhasil"})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	if r.Method == "OPTIONS" {
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Hanya POST yang diperbolehkan", http.StatusMethodNotAllowed)
		return
	}

	var input User
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Format JSON salah", http.StatusBadRequest)
		return
	}

	for _, u := range users {
		if u.Email == input.Email && u.Password == input.Password {
			json.NewEncoder(w).Encode(map[string]string{"token": "token123"})
			return
		}
	}

	http.Error(w, "Login gagal, email atau password salah", http.StatusUnauthorized)
}

func main() {
	http.HandleFunc("/api/register", RegisterHandler)
	http.HandleFunc("/api/login", LoginHandler)

	fmt.Println("🚀 Server running at http://0.0.0.0:8080")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}
