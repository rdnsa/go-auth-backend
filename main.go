package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("secret-key")

type User struct {
	Name     string
	Email    string
	Password string
}

var users = map[string]User{} // penyimpanan sementara

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var u User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil || u.Email == "" || u.Password == "" {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if _, exists := users[u.Email]; exists {
		http.Error(w, "Email already registered", http.StatusConflict)
		return
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	u.Password = string(hashed)
	users[u.Email] = u

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Register successful"})
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil || creds.Email == "" || creds.Password == "" {
		http.Error(w, "Invalid credentials", http.StatusBadRequest)
		return
	}

	user, exists := users[creds.Email]
	if !exists || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)) != nil {
		http.Error(w, "Email or password is wrong", http.StatusUnauthorized)
		return
	}

	expiration := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Email: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiration),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(jwtKey)

	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/api/register", registerHandler).Methods("POST")
	r.HandleFunc("/api/login", loginHandler).Methods("POST")

	fmt.Println("🚀 Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
