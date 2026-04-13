package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"zomato-backend-assignment/internal/model"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	DB *sql.DB
}

//
// ===== COMMON RESPONSE HELPERS =====
//

func respondValidationError(w http.ResponseWriter, fields map[string]string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"error":  "validation failed",
		"fields": fields,
	})
}

func respondError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}

//
// ===== REGISTER =====
//

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var user model.User

	// Decode request
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid request")
		return
	}

	// Validation
	errors := make(map[string]string)

	if user.Name == "" {
		errors["name"] = "is required"
	}
	if user.Email == "" {
		errors["email"] = "is required"
	}
	if user.Password == "" {
		errors["password"] = "is required"
	}

	if len(errors) > 0 {
		respondValidationError(w, errors)
		return
	}

	// Generate UUID
	user.ID = uuid.New().String()

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	// Insert into DB (with created_at)
	query := `
	INSERT INTO users (id, name, email, password, created_at)
	VALUES ($1, $2, $3, $4, $5)
	`

	_, err = h.DB.Exec(
		query,
		user.ID,
		user.Name,
		user.Email,
		string(hashedPassword),
		time.Now(),
	)

	if err != nil {
		// Handle duplicate email
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			respondValidationError(w, map[string]string{
				"email": "already exists",
			})
			return
		}

		respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	// Response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
	})
}

//
// ===== LOGIN =====
//

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input model.User

	// Decode request
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid request")
		return
	}

	// Validation
	errors := make(map[string]string)

	if input.Email == "" {
		errors["email"] = "is required"
	}
	if input.Password == "" {
		errors["password"] = "is required"
	}

	if len(errors) > 0 {
		respondValidationError(w, errors)
		return
	}

	var storedUser model.User

	// Fetch user
	query := "SELECT id, name, email, password FROM users WHERE email=$1"

	err = h.DB.QueryRow(query, input.Email).Scan(
		&storedUser.ID,
		&storedUser.Name,
		&storedUser.Email,
		&storedUser.Password,
	)

	if err != nil {
		respondError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}

	// Compare password
	err = bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(input.Password))
	if err != nil {
		respondError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}

	// JWT Secret from ENV
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		respondError(w, http.StatusInternalServerError, "missing JWT secret")
		return
	}

	// Generate token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": storedUser.ID,
		"email":   storedUser.Email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	// Response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"token": tokenString,
		"user": map[string]string{
			"id":    storedUser.ID,
			"name":  storedUser.Name,
			"email": storedUser.Email,
		},
	})
}