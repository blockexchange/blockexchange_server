package web

import (
	"blockexchange/types"
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

func SendError(w http.ResponseWriter, message string) {
	logrus.Trace("web.SendError: " + message)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(types.ErrorResponse{Message: message})
}
