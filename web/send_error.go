package web

import (
	"blockexchange/types"
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

func SendError(w http.ResponseWriter, code int, message string) {
	logrus.Trace("web.SendError: " + message)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(types.ErrorResponse{Message: message})
}
