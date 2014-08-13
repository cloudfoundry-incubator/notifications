package handlers

import (
    "encoding/json"
    "net/http"

    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/cloudfoundry-incubator/notifications/postal"
    "github.com/cloudfoundry-incubator/notifications/web/handlers/params"
)

type ErrorWriterInterface interface {
    Write(http.ResponseWriter, error)
}

type ErrorWriter struct{}

func NewErrorWriter() ErrorWriter {
    return ErrorWriter{}
}

func (writer ErrorWriter) Write(w http.ResponseWriter, err error) {
    switch err.(type) {
    case postal.CCDownError:
        writer.write(w, http.StatusBadGateway, []string{"Cloud Controller is unavailable"})
    case postal.CCNotFoundError:
        writer.write(w, http.StatusNotFound, []string{err.Error()})
    case postal.UAADownError:
        writer.write(w, http.StatusBadGateway, []string{"UAA is unavailable"})
    case postal.UAAGenericError:
        writer.write(w, http.StatusBadGateway, []string{err.Error()})
    case postal.TemplateLoadError:
        writer.write(w, http.StatusInternalServerError, []string{"An email template could not be loaded"})
    case params.ParseError:
        writer.write(w, 422, []string{err.Error()})
    case params.ValidationError:
        writer.write(w, 422, err.(params.ValidationError).Errors())
    case models.ErrDuplicateRecord:
        writer.write(w, 409, []string{err.Error()})
    case models.ErrRecordNotFound:
        writer.write(w, 404, []string{err.Error()})
    default:
        panic(err)
    }
}

func (writer ErrorWriter) write(w http.ResponseWriter, code int, errors []string) {
    response, err := json.Marshal(map[string][]string{
        "errors": errors,
    })
    if err != nil {
        panic(err)
    }

    w.WriteHeader(code)
    w.Write(response)
}
