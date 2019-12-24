package preferences

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/cloudfoundry-incubator/notifications/v1/models"
	"github.com/cloudfoundry-incubator/notifications/v1/services"
	"github.com/cloudfoundry-incubator/notifications/v1/web/webutil"
	"github.com/cloudfoundry-incubator/notifications/valiant"
	"github.com/ryanmoran/stack"
)

type UpdateUserPreferencesHandler struct {
	preferences preferenceUpdater
	errorWriter errorWriter
}

func NewUpdateUserPreferencesHandler(preferences preferenceUpdater, errWriter errorWriter) UpdateUserPreferencesHandler {
	return UpdateUserPreferencesHandler{
		preferences: preferences,
		errorWriter: errWriter,
	}
}

func (h UpdateUserPreferencesHandler) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	database := context.Get("database").(DatabaseInterface)
	connection := database.Connection()

	userGUID := regexp.MustCompile(".*/user_preferences/(.*)").FindStringSubmatch(req.URL.Path)[1]

	builder := services.NewPreferencesBuilder()
	validator := valiant.NewValidator(req.Body)
	err := validator.Validate(&builder)
	if err != nil {
		h.errorWriter.Write(w, webutil.ValidationError{Err: err})
		return
	}

	preferences, err := builder.ToPreferences()
	if err != nil {
		h.errorWriter.Write(w, err)
		return
	}

	transaction := connection.Transaction()
	transaction.Begin()
	err = h.preferences.Update(transaction, preferences, builder.GlobalUnsubscribe, userGUID)
	if err != nil {
		transaction.Rollback()

		switch err.(type) {
		case services.MissingKindOrClientError, services.CriticalKindError:
			h.errorWriter.Write(w, webutil.ValidationError{Err: err})
		default:
			h.errorWriter.Write(w, err)
		}
		return
	}

	err = transaction.Commit()
	if err != nil {
		h.errorWriter.Write(w, models.TransactionCommitError{Err: err})
		return
	}

	w.WriteHeader(http.StatusNoContent)
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
