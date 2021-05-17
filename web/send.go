package web

import (
	"blockexchange/types"
	"bytes"
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

func SendJson(w http.ResponseWriter, o interface{}) []byte {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	buf := bytes.NewBuffer([]byte{})
	json.NewEncoder(buf).Encode(o)
	w.Write(buf.Bytes())
	return buf.Bytes()
}

func SendRawJson(w http.ResponseWriter, data []byte) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func Send(w http.ResponseWriter, o interface{}, err error) {
	if err != nil {
		SendError(w, 500, err.Error())
	} else {
		SendJson(w, o)
	}
}
