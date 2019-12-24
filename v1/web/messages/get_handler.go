package messages

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/cloudfoundry-incubator/notifications/v1/services"
	"github.com/ryanmoran/stack"
)

type GetHandler struct {
	finder      messageFinder
	errorWriter errorWriter
}

type errorWriter interface {
	Write(writer http.ResponseWriter, err error)
}

type messageFinder interface {
	Find(services.DatabaseInterface, string) (services.Message, error)
}

func NewGetHandler(finder messageFinder, errWriter errorWriter) GetHandler {
	return GetHandler{
		finder:      finder,
		errorWriter: errWriter,
	}
}

func (h GetHandler) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	messageID := strings.Split(req.URL.Path, "/messages/")[1]

	message, err := h.finder.Find(context.Get("database").(DatabaseInterface), messageID)
	if err != nil {
		h.errorWriter.Write(w, err)
		return
	}

	var document struct {
		Status string `json:"status"`
	}
	document.Status = message.Status

	writeJSON(w, http.StatusOK, document)
}

func writeJSON(w http.ResponseWriter, status int, object interface{}) {
	output, err := json.Marshal(object)
	if err != nil {
		panic(err) // No JSON we write into a response should ever panic
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(output)
}
