package web

import (
	"blockexchange/types"
	"encoding/json"
	"net/http"
)

func SendError(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(types.ErrorResponse{Message: message})
}
