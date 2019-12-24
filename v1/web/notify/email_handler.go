package notify

import (
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/v1/services"
	"github.com/ryanmoran/stack"
)

const (
	VCAPRequestIDKey    = "vcap_request_id"
	RequestReceivedTime = "request_received_time"
)

type notifyExecutor interface {
	Execute(conn ConnectionInterface, req *http.Request, context stack.Context, guid string, strategy Dispatcher, validator ValidatorInterface, vcapRequestID string) (response []byte, err error)
}

type errorWriter interface {
	Write(writer http.ResponseWriter, err error)
}

type Dispatcher interface {
	Dispatch(dispatch services.Dispatch) ([]services.Response, error)
}

type EmailHandler struct {
	errorWriter errorWriter
	notify      notifyExecutor
	strategy    Dispatcher
}

func NewEmailHandler(notify notifyExecutor, errWriter errorWriter, strategy Dispatcher) EmailHandler {
	return EmailHandler{
		errorWriter: errWriter,
		notify:      notify,
		strategy:    strategy,
	}
}

func (h EmailHandler) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	vcapRequestID := context.Get(VCAPRequestIDKey).(string)
	database := context.Get("database").(DatabaseInterface)
	conn := database.Connection()

	output, err := h.notify.Execute(conn, req, context, "", h.strategy, EmailValidator{}, vcapRequestID)
	if err != nil {
		h.errorWriter.Write(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(output)
}
