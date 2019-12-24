package notify

import (
	"net/http"
	"strings"

	"github.com/ryanmoran/stack"
)

type UserHandler struct {
	errorWriter errorWriter
	notify      notifyExecutor
	strategy    Dispatcher
}

func NewUserHandler(notify notifyExecutor, errWriter errorWriter, strategy Dispatcher) UserHandler {
	return UserHandler{
		errorWriter: errWriter,
		notify:      notify,
		strategy:    strategy,
	}
}

func (h UserHandler) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	conn := context.Get("database").(DatabaseInterface).Connection()
	userGUID := strings.TrimPrefix(req.URL.Path, "/users/")
	vcapRequestID := context.Get(VCAPRequestIDKey).(string)

	output, err := h.notify.Execute(conn, req, context, userGUID, h.strategy, GUIDValidator{}, vcapRequestID)
	if err != nil {
		h.errorWriter.Write(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(output)
}
