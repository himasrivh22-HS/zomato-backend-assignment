package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"zomato-backend-assignment/internal/model"

	"golang.org/x/crypto/bcrypt"
	"github.com/google/uuid"

	"github.com/golang-jwt/jwt/v5"
)

type AuthHandler struct {
	DB *sql.DB
}

// REGISTER
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var user model.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	user.ID = uuid.New().String()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	query := "INSERT INTO users (id, name, email, password) VALUES ($1, $2, $3, $4)"

	_, err = h.DB.Exec(query, user.ID, user.Name, user.Email, string(hashedPassword))
	if err != nil {
		fmt.Println("DB ERROR:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(user)
}

// LOGIN
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input model.User

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var storedUser model.User

	query := "SELECT id, name, email, password FROM users WHERE email=$1"

	err = h.DB.QueryRow(query, input.Email).Scan(
		&storedUser.ID,
		&storedUser.Name,
		&storedUser.Email,
		&storedUser.Password,
	)

	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(input.Password))
	if err != nil {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	"user_id": storedUser.ID,
	"email":   storedUser.Email,
	"exp":     time.Now().Add(time.Hour * 24).Unix(),
})

// Secret key (you can change later)
tokenString, err := token.SignedString([]byte("mysecretkey"))
if err != nil {
	http.Error(w, "Error generating token", http.StatusInternalServerError)
	return
}

// Send token as response
json.NewEncoder(w).Encode(map[string]string{
	"token": tokenString,
})
}