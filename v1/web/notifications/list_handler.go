package notifications

import (
	"encoding/json"
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/v1/models"
	"github.com/cloudfoundry-incubator/notifications/v1/services"
	"github.com/ryanmoran/stack"
)

type NotificationsByClient map[string]Client

type Client struct {
	Name          string                  `json:"name"`
	Template      string                  `json:"template"`
	Notifications map[string]Notification `json:"notifications"`
}

type listsAllClientsAndNotifications interface {
	AllClientsAndNotifications(services.DatabaseInterface) ([]models.Client, []models.Kind, error)
}

type Notification struct {
	Description string `json:"description"`
	Template    string `json:"template"`
	Critical    bool   `json:"critical"`
}

type ListHandler struct {
	finder      listsAllClientsAndNotifications
	errorWriter errorWriter
}

func NewListHandler(notificationsFinder listsAllClientsAndNotifications, errWriter errorWriter) ListHandler {
	return ListHandler{
		finder:      notificationsFinder,
		errorWriter: errWriter,
	}
}

func (h ListHandler) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	clients, notifications, err := h.finder.AllClientsAndNotifications(context.Get("database").(DatabaseInterface))

	if err != nil {
		h.errorWriter.Write(w, err)
		return
	}

	notificationsByClient := h.constructNotifications(clients, notifications)

	writeJSON(w, http.StatusOK, notificationsByClient)
}

func (h ListHandler) constructNotifications(clients []models.Client, notifications []models.Kind) NotificationsByClient {
	notificationsByClient := NotificationsByClient{}

	for _, client := range clients {
		clientWithNotifications := Client{
			Name:     client.Description,
			Template: client.TemplateToUse(),
		}

		clientNotifications := make(map[string]Notification)
		for _, notification := range notifications {
			if notification.ClientID == client.ID {
				clientNotifications[notification.ID] = Notification{
					Description: notification.Description,
					Template:    notification.TemplateToUse(),
					Critical:    notification.Critical,
				}
			}
		}

		clientWithNotifications.Notifications = clientNotifications
		notificationsByClient[client.ID] = clientWithNotifications
	}

	return notificationsByClient
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
