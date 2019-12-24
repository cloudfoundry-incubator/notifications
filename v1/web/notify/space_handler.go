package notify

import (
	"net/http"
	"strings"

	"github.com/ryanmoran/stack"
)

type SpaceHandler struct {
	errorWriter errorWriter
	notify      notifyExecutor
	strategy    Dispatcher
}

func NewSpaceHandler(notify notifyExecutor, errWriter errorWriter, strategy Dispatcher) SpaceHandler {
	return SpaceHandler{
		errorWriter: errWriter,
		notify:      notify,
		strategy:    strategy,
	}
}

func (h SpaceHandler) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	conn := context.Get("database").(DatabaseInterface).Connection()
	spaceGUID := strings.TrimPrefix(req.URL.Path, "/spaces/")
	vcapRequestID := context.Get(VCAPRequestIDKey).(string)

	output, err := h.notify.Execute(conn, req, context, spaceGUID, h.strategy, GUIDValidator{}, vcapRequestID)
	if err != nil {
		h.errorWriter.Write(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(output)
}
