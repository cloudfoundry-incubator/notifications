package notify

import (
	"net/http"
	"strings"

	"github.com/ryanmoran/stack"
)

type OrganizationHandler struct {
	errorWriter errorWriter
	notify      notifyExecutor
	strategy    Dispatcher
}

func NewOrganizationHandler(notify notifyExecutor, errWriter errorWriter, strategy Dispatcher) OrganizationHandler {
	return OrganizationHandler{
		errorWriter: errWriter,
		notify:      notify,
		strategy:    strategy,
	}
}

func (h OrganizationHandler) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	conn := context.Get("database").(DatabaseInterface).Connection()
	orgGUID := strings.TrimPrefix(req.URL.Path, "/organizations/")
	vcapRequestID := context.Get(VCAPRequestIDKey).(string)

	output, err := h.notify.Execute(conn, req, context, orgGUID, h.strategy, GUIDValidator{}, vcapRequestID)
	if err != nil {
		h.errorWriter.Write(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(output)
}
