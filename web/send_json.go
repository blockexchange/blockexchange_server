package web

import (
	"encoding/json"
	"net/http"
)

func SendJson(w http.ResponseWriter, o interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(o)
}
