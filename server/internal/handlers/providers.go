package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
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

// Haversine distance calculation
func calculateDistance(lat1, lng1, lat2, lng2 float64) float64 {
	const R = 3958.8 // Radius of the Earth in miles
	dLat := (lat2 - lat1) * math.Pi / 180
	dLng := (lng2 - lng1) * math.Pi / 180
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLng/2)*math.Sin(dLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}

type ProviderCard struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	ServiceText string `json:"serviceText"`
	Rating      string `json:"rating"`
	Image       string `json:"image"`
}

type ProvidersListResponse struct {
	Success   bool           `json:"success"`
	Providers []ProviderCard `json:"providers,omitempty"`
	Message   string         `json:"message,omitempty"`
}

// ListProvidersHandler handles GET /api/providers
func ListProvidersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// Validate JWT token
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ProvidersListResponse{
			Success: false,
			Message: "Unauthorized access",
		})
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
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ProvidersListResponse{
			Success: false,
			Message: "Unauthorized access",
		})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ProvidersListResponse{
			Success: false,
			Message: "Unauthorized access",
		})
		return
	}

	viewerID, _ := claims["id"].(string)

	// Query DB for all completed provider profiles
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	usersCollection := database.Client.Database("neighbourly").Collection("users")

	// Get viewer info for location
	var viewer models.User
	if viewerID != "" {
		vObjID, _ := primitive.ObjectIDFromHex(viewerID)
		usersCollection.FindOne(ctx, bson.M{"_id": vObjID}).Decode(&viewer)
	}

	// Match providers with completed profiles
	filter := bson.M{
		"role": bson.M{
			"$in": []string{"provider", "Provider", "service provider", "Service Provider"},
		},
		"isProfileComplete": true,
	}

	cursor, err := usersCollection.Find(ctx, filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ProvidersListResponse{
			Success: false,
			Message: "Error fetching providers",
		})
		return
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err := cursor.All(ctx, &users); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ProvidersListResponse{
			Success: false,
			Message: "Error reading providers",
		})
		return
	}

	// Build provider cards
	providers := make([]ProviderCard, 0, len(users))
	for _, u := range users {
		// Only include providers who have at least one service
		if len(u.Services) == 0 {
			continue
		}

		title := u.Title
		if title == "" {
			title = "Service Provider"
		}
		serviceText := strings.ToUpper(title)

		if viewer.Location != nil && u.Location != nil {
			dist := calculateDistance(viewer.Location.Lat, viewer.Location.Lng, u.Location.Lat, u.Location.Lng)
			serviceText += " • " + strings.ToUpper(strings.ReplaceAll(strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.1f", dist), "0"), "."), ".0", "")) + " MI"
		}

		providers = append(providers, ProviderCard{
			ID:          u.ID.Hex(),
			Name:        u.FullName,
			ServiceText: serviceText,
			Rating:      "5.0",
			Image:       u.Avatar,
		})
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ProvidersListResponse{
		Success:   true,
		Providers: providers,
	})
}
