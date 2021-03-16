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

func SendJson(w http.ResponseWriter, o interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(o)
}

func Send(w http.ResponseWriter, o interface{}, err error) {
	if err != nil {
		SendError(w, 500, err.Error())
	} else {
		SendJson(w, o)
	}
}
