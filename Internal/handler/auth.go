package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"strings"
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

// ===== REGISTER =====

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var user model.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid request")
		return
	}

	errors := make(map[string]string)

	if user.Name == "" {
		errors["name"] = "is required"
	}
	if user.Email == "" {
		errors["email"] = "is required"
	}
	if user.Email != "" && !strings.Contains(user.Email, "@") { 
		errors["email"] = "invalid format"
	}
	if user.Password == "" {
		errors["password"] = "is required"
	}

	if len(errors) > 0 {
		writeValidationError(w, errors)
		return
	}

	user.ID = uuid.New().String()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

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
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			writeValidationError(w, map[string]string{
				"email": "already exists",
			})
			return
		}

		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
	})
}

// ===== LOGIN =====

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input model.User

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid request")
		return
	}

	errors := make(map[string]string)

	if input.Email == "" {
		errors["email"] = "is required"
	}
	if input.Password == "" {
		errors["password"] = "is required"
	}

	if len(errors) > 0 {
		writeValidationError(w, errors)
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
		writeError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(input.Password))
	if err != nil {
		writeError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		writeError(w, http.StatusInternalServerError, "missing jwt secret")
		return
	}

	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": storedUser.ID,
		"email":   storedUser.Email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

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