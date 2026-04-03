package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"

	"neighbourly/server/internal/database"
	"neighbourly/server/internal/models"
)

type SignupRequest struct {
	FullName    string `json:"fullName"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	Password    string `json:"password"`
	Role        string `json:"role"`
}

type AuthResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Token   string `json:"token,omitempty"`
}

type LoginRequest struct {
	PhoneNumber string `json:"phoneNumber"`
	Password    string `json:"password"`
	Role        string `json:"role,omitempty"`
}

// SignupHandler registers a new user
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var req SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendAuthResponse(w, http.StatusBadRequest, false, "Invalid request body", "")
		return
	}

	// Validate required fields
	if req.Email == "" || req.Password == "" || req.FullName == "" {
		sendAuthResponse(w, http.StatusBadRequest, false, "Missing required fields", "")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Hardcoding database name "neighbourly" for now
	usersCollection := database.Client.Database("neighbourly").Collection("users")

	// Check if user already exists
	count, err := usersCollection.CountDocuments(ctx, bson.M{"email": req.Email})
	if err != nil {
		sendAuthResponse(w, http.StatusInternalServerError, false, "Database error checking user", "")
		return
	}
	if count > 0 {
		sendAuthResponse(w, http.StatusConflict, false, "Email already exists", "")
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		sendAuthResponse(w, http.StatusInternalServerError, false, "Error hashing password", "")
		return
	}

	// Create user model
	user := models.User{
		FullName:    req.FullName,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		Password:    string(hashedPassword),
		Role:        req.Role,
	}

	// Insert into DB
	_, err = usersCollection.InsertOne(ctx, user)
	if err != nil {
		sendAuthResponse(w, http.StatusInternalServerError, false, "Error creating user", "")
		return
	}

	// Generate JWT token
	token, err := generateJWT(user.Email, user.Role)
	if err != nil {
		sendAuthResponse(w, http.StatusInternalServerError, false, "User created but error generating token", "")
		return
	}

	sendAuthResponse(w, http.StatusCreated, true, "Account created successfully.", token)
}

// LoginHandler authenticates an existing user
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendAuthResponse(w, http.StatusBadRequest, false, "Invalid request body", "")
		return
	}

	// Validate required fields
	if req.PhoneNumber == "" || req.Password == "" {
		sendAuthResponse(w, http.StatusBadRequest, false, "Missing required fields", "")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	usersCollection := database.Client.Database("neighbourly").Collection("users")

	var user models.User
	err := usersCollection.FindOne(ctx, bson.M{"phoneNumber": req.PhoneNumber}).Decode(&user)
	if err != nil {
		sendAuthResponse(w, http.StatusUnauthorized, false, "Invalid phone number or password", "")
		return
	}

	// Verify the password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		sendAuthResponse(w, http.StatusUnauthorized, false, "Invalid phone number or password", "")
		return
	}

	// Generate JWT token
	token, err := generateJWT(user.Email, user.Role)
	if err != nil {
		sendAuthResponse(w, http.StatusInternalServerError, false, "Error generating token", "")
		return
	}

	sendAuthResponse(w, http.StatusOK, true, "Login successful", token)
}

func sendAuthResponse(w http.ResponseWriter, status int, success bool, message, token string) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(AuthResponse{
		Success: success,
		Message: message,
		Token:   token,
	})
}

func generateJWT(email string, role string) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "my_super_secret_jwt_key_123" // Fallback secret
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"role":  role,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	})

	return token.SignedString([]byte(jwtSecret))
}
