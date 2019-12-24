package notify

import (
	"net/http"
	"strings"

	"github.com/ryanmoran/stack"
)

type UAAScopeHandler struct {
	errorWriter errorWriter
	notify      notifyExecutor
	strategy    Dispatcher
}

func NewUAAScopeHandler(notify notifyExecutor, errWriter errorWriter, strategy Dispatcher) UAAScopeHandler {
	return UAAScopeHandler{
		errorWriter: errWriter,
		notify:      notify,
		strategy:    strategy,
	}
}

func (h UAAScopeHandler) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	conn := context.Get("database").(DatabaseInterface).Connection()
	scope := strings.TrimPrefix(req.URL.Path, "/uaa_scopes/")
	vcapRequestID := context.Get(VCAPRequestIDKey).(string)

	output, err := h.notify.Execute(conn, req, context, scope, h.strategy, GUIDValidator{}, vcapRequestID)
	if err != nil {
		h.errorWriter.Write(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(output)
}
