package service

import (
	"encoding/json"
	"net/http"
)

func HandleSignup(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"error": "false", "message": "hello world"})
}
