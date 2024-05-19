package controller

import (
	"encoding/json"
	"net/http"
)

func GetRate(w http.ResponseWriter, r *http.Request) {
	rate := 27.5

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rate)
}
