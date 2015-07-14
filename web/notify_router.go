package web

import (
	"github.com/cloudfoundry-incubator/notifications/metrics"
	"github.com/cloudfoundry-incubator/notifications/services"
	"github.com/cloudfoundry-incubator/notifications/web/handlers"
	"github.com/gorilla/mux"
	"github.com/ryanmoran/stack"
)

func NewNotifyRouter(notify handlers.NotifyInterface,
	errorWriter handlers.ErrorWriterInterface,
	userStrategy services.StrategyInterface,
	logging RequestLogging,
	notificationsWriteAuthenticator Authenticator,
	databaseAllocator DatabaseAllocator,
	spaceStrategy services.StrategyInterface,
	organizationStrategy services.StrategyInterface,
	everyoneStrategy services.StrategyInterface,
	uaaScopeStrategy services.StrategyInterface,
	emailStrategy services.StrategyInterface,
	emailsWriteAuthenticator Authenticator) *mux.Router {

	router := mux.NewRouter()
	requestCounter := NewRequestCounter(router, metrics.DefaultLogger)

	router.Handle("/users/{user_id}", stack.NewStack(handlers.NewNotifyUser(notify, errorWriter, userStrategy)).Use(logging, requestCounter, notificationsWriteAuthenticator, databaseAllocator)).Methods("POST").Name("POST /users/{user_id}")
	router.Handle("/spaces/{space_id}", stack.NewStack(handlers.NewNotifySpace(notify, errorWriter, spaceStrategy)).Use(logging, requestCounter, notificationsWriteAuthenticator, databaseAllocator)).Methods("POST").Name("POST /spaces/{space_id}")
	router.Handle("/organizations/{org_id}", stack.NewStack(handlers.NewNotifyOrganization(notify, errorWriter, organizationStrategy)).Use(logging, requestCounter, notificationsWriteAuthenticator, databaseAllocator)).Methods("POST").Name("POST /organizations/{org_id}")
	router.Handle("/everyone", stack.NewStack(handlers.NewNotifyEveryone(notify, errorWriter, everyoneStrategy)).Use(logging, requestCounter, notificationsWriteAuthenticator, databaseAllocator)).Methods("POST").Name("POST /everyone")
	router.Handle("/uaa_scopes/{scope}", stack.NewStack(handlers.NewNotifyUAAScope(notify, errorWriter, uaaScopeStrategy)).Use(logging, requestCounter, notificationsWriteAuthenticator, databaseAllocator)).Methods("POST").Name("POST /uaa_scopes/{scope}")
	router.Handle("/emails", stack.NewStack(handlers.NewNotifyEmail(notify, errorWriter, emailStrategy)).Use(logging, requestCounter, emailsWriteAuthenticator, databaseAllocator)).Methods("POST").Name("POST /emails")

	return router
}