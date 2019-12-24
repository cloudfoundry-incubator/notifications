package templates

import (
	"encoding/json"
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/v1/collections"
	"github.com/cloudfoundry-incubator/notifications/v1/web/webutil"
	"github.com/ryanmoran/stack"
)

type errorWriter interface {
	Write(writer http.ResponseWriter, err error)
}

type templateCreator interface {
	Create(connection collections.ConnectionInterface, template collections.Template) (collections.Template, error)
}

type CreateHandler struct {
	creator     templateCreator
	errorWriter errorWriter
}

func NewCreateHandler(creator templateCreator, errWriter errorWriter) CreateHandler {
	return CreateHandler{
		creator:     creator,
		errorWriter: errWriter,
	}
}

func (h CreateHandler) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	templateParams, err := NewTemplateParams(req.Body)
	if err != nil {
		h.errorWriter.Write(w, err)
		return
	}

	connection := context.Get("database").(DatabaseInterface).Connection()

	template, err := h.creator.Create(connection, collections.Template{
		Name:     templateParams.Name,
		Text:     templateParams.Text,
		HTML:     templateParams.HTML,
		Subject:  templateParams.Subject,
		Metadata: string(templateParams.Metadata),
	})
	if err != nil {
		h.errorWriter.Write(w, webutil.TemplateCreateError{})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"template_id":"` + template.ID + `"}`))
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
