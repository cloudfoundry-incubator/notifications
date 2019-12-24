package webutil

import (
	"encoding/json"
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/cf"
	"github.com/cloudfoundry-incubator/notifications/v1/collections"
	"github.com/cloudfoundry-incubator/notifications/v1/models"
	"github.com/cloudfoundry-incubator/notifications/v1/services"
)

type ErrorWriter struct{}

func NewErrorWriter() ErrorWriter {
	return ErrorWriter{}
}

func (writer ErrorWriter) Write(w http.ResponseWriter, err error) {

	w.Header().Set("Content-Type", "application/json")

	switch err.(type) {
	case UAAScopesError, CriticalNotificationError, collections.TemplateAssignmentError, MissingUserTokenError, ValidationError:
		w.WriteHeader(422)
	case services.CCDownError:
		w.WriteHeader(http.StatusBadGateway)
	case services.CCNotFoundError, models.NotFoundError, cf.NotFoundError:
		w.WriteHeader(http.StatusNotFound)
	case ParseError, SchemaError:
		w.WriteHeader(http.StatusBadRequest)
	case models.DuplicateError:
		w.WriteHeader(http.StatusConflict)
	case services.DefaultScopeError:
		w.WriteHeader(http.StatusNotAcceptable)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(map[string][]string{
		"errors": []string{err.Error()},
	})
}
