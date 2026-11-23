package api

import (
	"blockexchange/core"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

// the default origin check only allows same-origin connections,
// required because the session cookie authenticates the upgrade request
var wasmUpgrader = websocket.Upgrader{}

// websocket proxy for the luanti wasm client, only to this server (see core.WasmProxy)
func (api *Api) WasmWebsocket(w http.ResponseWriter, r *http.Request, ctx *SecureContext) {
	schema_uid := mux.Vars(r)["schema_uid"]

	if api.loadOwnedSchema(w, ctx, schema_uid, true) == nil {
		return
	}

	conn, err := wasmUpgrader.Upgrade(w, r, nil)
	if err != nil {
		// the upgrader already sent an error response
		return
	}

	log := logrus.WithFields(logrus.Fields{"schema_uid": schema_uid, "user_uid": ctx.Claims.UserUID})
	err = api.core.NewWasmProxy().Handle(conn)
	if errors.Is(err, core.ErrWasmProxyNotAllowed) {
		log.WithError(err).Warn("wasm proxy: refused")
	} else if err != nil {
		log.WithError(err).Debug("wasm proxy: connection closed")
	}
}
