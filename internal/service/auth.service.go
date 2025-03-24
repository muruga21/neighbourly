package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/redis/go-redis/v9"
)

func HandleSignup(w http.ResponseWriter, r *http.Request) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // No password set
		DB:       0,  // Use default DB
		Protocol: 2,  // Connection protocol
	})
	fmt.Println(client)
	json.NewEncoder(w).Encode(map[string]string{"error": "false", "message": "hello world"})
}
