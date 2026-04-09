package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	
	"neighbourly/server/internal/database"
	"neighbourly/server/internal/models"
)

type ServiceData struct {
	ID    string `json:"id"`
	Icon  string `json:"icon"`
	Title string `json:"title"`
	Desc  string `json:"desc"`
	Price string `json:"price"`
}

type ReviewData struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	Rating int    `json:"rating"`
	Text   string `json:"text"`
}

type ProviderProfileData struct {
	ID           string        `json:"id"`
	Name         string        `json:"name"`
	Title        string        `json:"title"`
	Status       string        `json:"status"`
	Rating       string        `json:"rating"`
	ReviewsCount int           `json:"reviewsCount"`
	ResponseTime string        `json:"responseTime"`
	Experience   string        `json:"experience"`
	Avatar       string        `json:"avatar"`
	Services     []ServiceData `json:"services"`
	Reviews      []ReviewData  `json:"reviews"`
}

type SeekerStats struct {
	Bookings    int `json:"bookings"`
	ReviewsLeft int `json:"reviewsLeft"`
}

type SeekerProfileData struct {
	Name        string      `json:"name"`
	Email       string      `json:"email"`
	MemberSince string      `json:"memberSince"`
	Avatar      string      `json:"avatar"`
	Stats       SeekerStats `json:"stats"`
}

type ProfileResponse struct {
	Success bool        `json:"success"`
	Profile interface{} `json:"profile,omitempty"`
	Message string      `json:"message,omitempty"`
}

// ProfileHandler handles fetching user profiles
func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// 1. Extract userid from the URL parameters
	userID := r.PathValue("userid")
	
	// Fallback if r.PathValue is empty (just in case)
	if userID == "" {
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) > 0 {
			userID = parts[len(parts)-1]
		}
	}

	// 2. If userid === "null", decode the token and extract the ID
	if userID == "null" {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			sendProfileError(w, http.StatusUnauthorized, "Unauthorized access: Token missing")
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		
		jwtSecret := os.Getenv("JWT_SECRET")
		if jwtSecret == "" {
			jwtSecret = "my_super_secret_jwt_key_123"
		}

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			// Ensure token method is HMAC
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, http.ErrAbortHandler
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			log.Printf("Token validation error: %v", err)
			sendProfileError(w, http.StatusUnauthorized, "Unauthorized access: Invalid token")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			sendProfileError(w, http.StatusUnauthorized, "Unauthorized access: Invalid token claims")
			return
		}

		// Extract ID from JWT payload
		if idClaim, exists := claims["id"].(string); exists && idClaim != "" {
			userID = idClaim
		} else {
			sendProfileError(w, http.StatusNotFound, "Failed to extract user ID from token")
			return
		}
	}

	if userID == "" {
		sendProfileError(w, http.StatusBadRequest, "User ID is required")
		return
	}

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		sendProfileError(w, http.StatusBadRequest, "Invalid User ID format")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	usersCollection := database.Client.Database("neighbourly").Collection("users")

	var user models.User
	err = usersCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		sendProfileError(w, http.StatusNotFound, "User not found")
		return
	}

	if !user.IsProfileComplete {
		sendProfileError(w, http.StatusInternalServerError, "incomplete profile")
		return
	}

	if strings.ToLower(user.Role) == "provider" || strings.ToLower(user.Role) == "service provider" {
		// For providers, services must not be empty to be considered complete
		if len(user.Services) == 0 {
			sendProfileError(w, http.StatusInternalServerError, "incomplete profile")
			return
		}

		title := user.Title
		if title == "" {
			title = "Neighbourly Provider"
		}
		status := user.Status
		if status == "" {
			status = "Available"
		}
		experience := user.Experience
		if experience == "" {
			experience = "-"
		}

		// Map DB services to response format
		services := make([]ServiceData, len(user.Services))
		for i, s := range user.Services {
			services[i] = ServiceData{
				ID:    s.ID,
				Icon:  s.Icon,
				Title: s.Title,
				Desc:  s.Desc,
				Price: s.Price,
			}
		}

		// Map DB reviews to response format
		reviews := make([]ReviewData, len(user.Reviews))
		for i, r := range user.Reviews {
			reviews[i] = ReviewData{
				ID:     r.ID,
				Name:   r.Name,
				Avatar: r.Avatar,
				Rating: r.Rating,
				Text:   r.Text,
			}
		}

		profile := ProviderProfileData{
			ID:           user.ID.Hex(),
			Name:         user.FullName,
			Title:        title,
			Status:       status,
			Rating:       "5.0",
			ReviewsCount: len(reviews),
			ResponseTime: "-",
			Experience:   experience,
			Avatar:       user.Avatar,
			Services:     services,
			Reviews:      reviews,
		}
		json.NewEncoder(w).Encode(ProfileResponse{
			Success: true,
			Profile: profile,
		})
	} else {
		profile := SeekerProfileData{
			Name:        user.FullName,
			Email:       user.Email,
			MemberSince: "Oct 2023", // Basic default since we don't track creation date yet
			Avatar:      user.Avatar,
			Stats: SeekerStats{
				Bookings:    0,
				ReviewsLeft: 0,
			},
		}
		json.NewEncoder(w).Encode(ProfileResponse{
			Success: true,
			Profile: profile,
		})
	}
}

func sendProfileError(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ProfileResponse{
		Success: false,
		Message: message,
	})
}

// --- Update Profile Endpoints ---

type UpdateServiceData struct {
	ID    string `json:"id"`
	Icon  string `json:"icon"`
	Title string `json:"title"`
	Desc  string `json:"desc"`
	Price string `json:"price"`
}

type UpdateLocationData struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type UpdateProfileRequest struct {
	Avatar     string              `json:"avatar,omitempty"`
	Title      string              `json:"title,omitempty"`
	Experience string              `json:"experience,omitempty"`
	BasePrice  string              `json:"basePrice,omitempty"`
	Status     string              `json:"status,omitempty"`
	Location   *UpdateLocationData `json:"location,omitempty"`
	Services   []UpdateServiceData `json:"services,omitempty"`
}

type UpdateProfileResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func sendUpdateProfileError(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(UpdateProfileResponse{
		Success: false,
		Message: message,
	})
}

// UpdateProfileHandler handles POST /api/profile/update
func UpdateProfileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		sendUpdateProfileError(w, http.StatusUnauthorized, "Unauthorized access")
		return
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "my_super_secret_jwt_key_123"
	}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, http.ErrAbortHandler
		}
		return []byte(jwtSecret), nil
	})

	if err != nil || !token.Valid {
		sendUpdateProfileError(w, http.StatusUnauthorized, "Unauthorized access")
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		sendUpdateProfileError(w, http.StatusUnauthorized, "Unauthorized access")
		return
	}

	userID, ok := claims["id"].(string)
	if !ok || userID == "" {
		sendUpdateProfileError(w, http.StatusUnauthorized, "Unauthorized access")
		return
	}

	var req UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendUpdateProfileError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		sendUpdateProfileError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	usersCollection := database.Client.Database("neighbourly").Collection("users")

	var user models.User
	err = usersCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		sendUpdateProfileError(w, http.StatusUnauthorized, "Unauthorized access")
		return
	}

	updateFields := bson.M{
		"isProfileComplete": true,
	}

	if req.Avatar != "" {
		updateFields["avatar"] = req.Avatar
	}

	if strings.ToLower(user.Role) == "provider" || strings.ToLower(user.Role) == "service provider" {
		if req.Title != "" {
			updateFields["title"] = req.Title
		}
		if req.Experience != "" {
			updateFields["experience"] = req.Experience
		}
		if req.BasePrice != "" {
			updateFields["basePrice"] = req.BasePrice
		}
		if req.Status != "" {
			updateFields["status"] = req.Status
		}
		if len(req.Services) > 0 {
			// Convert to models.Service for DB storage
			dbServices := make([]models.Service, len(req.Services))
			for i, s := range req.Services {
				dbServices[i] = models.Service{
					ID:    s.ID,
					Icon:  s.Icon,
					Title: s.Title,
					Desc:  s.Desc,
					Price: s.Price,
				}
			}
			updateFields["services"] = dbServices
		}
		if req.Location != nil {
			updateFields["location"] = models.Location{
				Lat: req.Location.Lat,
				Lng: req.Location.Lng,
			}
		}
	}

	update := bson.M{"$set": updateFields}

	_, err = usersCollection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		sendUpdateProfileError(w, http.StatusInternalServerError, "Internal server error saving profile")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(UpdateProfileResponse{
		Success: true,
		Message: "Profile updated successfully",
	})
}
