package notify

import (
	"net/http"

	"github.com/ryanmoran/stack"
)

type EveryoneHandler struct {
	errorWriter errorWriter
	notify      notifyExecutor
	strategy    Dispatcher
}

func NewEveryoneHandler(notify notifyExecutor, errWriter errorWriter, strategy Dispatcher) EveryoneHandler {
	return EveryoneHandler{
		errorWriter: errWriter,
		notify:      notify,
		strategy:    strategy,
	}
}

func (h EveryoneHandler) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	connection := context.Get("database").(DatabaseInterface).Connection()
	vcapRequestID := context.Get(VCAPRequestIDKey).(string)

	output, err := h.notify.Execute(connection, req, context, "", h.strategy, GUIDValidator{}, vcapRequestID)
	if err != nil {
		h.errorWriter.Write(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(output)
}
