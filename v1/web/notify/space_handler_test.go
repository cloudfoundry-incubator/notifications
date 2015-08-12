package notify_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"

	"github.com/cloudfoundry-incubator/notifications/testing/fakes"
	"github.com/cloudfoundry-incubator/notifications/v1/web/notify"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SpaceHandler", func() {
	Describe("ServeHTTP", func() {
		var (
			handler     notify.SpaceHandler
			writer      *httptest.ResponseRecorder
			request     *http.Request
			notifyObj   *fakes.Notify
			context     stack.Context
			connection  *fakes.Connection
			strategy    *fakes.Strategy
			errorWriter *fakes.ErrorWriter
		)

		BeforeEach(func() {
			writer = httptest.NewRecorder()
			request = &http.Request{URL: &url.URL{Path: "/spaces/space-001"}}
			strategy = fakes.NewStrategy()
			database := fakes.NewDatabase()
			connection = fakes.NewConnection()
			database.Conn = connection
			errorWriter = fakes.NewErrorWriter()

			context = stack.NewContext()
			context.Set("database", database)
			context.Set(notify.VCAPRequestIDKey, "some-request-id")

			notifyObj = fakes.NewNotify()
			handler = notify.NewSpaceHandler(notifyObj, errorWriter, strategy)
		})

		Context("when the notifyObj.Execute returns a successful response", func() {
			It("returns the JSON representation of the response", func() {
				notifyObj.ExecuteCall.Returns.Response = []byte("whatever")
				handler.ServeHTTP(writer, request, context)

				Expect(writer.Code).To(Equal(http.StatusOK))
				Expect(writer.Body.String()).To(Equal("whatever"))
			})

			It("delegates to the notifyObj object with the correct arguments", func() {
				handler.ServeHTTP(writer, request, context)

				Expect(reflect.ValueOf(notifyObj.ExecuteCall.Receives.Connection).Pointer()).To(Equal(reflect.ValueOf(connection).Pointer()))
				Expect(notifyObj.ExecuteCall.Receives.Request).To(Equal(request))
				Expect(notifyObj.ExecuteCall.Receives.Context).To(Equal(context))
				Expect(notifyObj.ExecuteCall.Receives.GUID).To(Equal("space-001"))
				Expect(notifyObj.ExecuteCall.Receives.Strategy).To(Equal(strategy))
				Expect(notifyObj.ExecuteCall.Receives.Validator).To(BeAssignableToTypeOf(notify.GUIDValidator{}))
				Expect(notifyObj.ExecuteCall.Receives.VCAPRequestID).To(Equal("some-request-id"))
			})
		})

		Context("when the notifyObj.Execute returns an error", func() {
			It("propagates the error", func() {
				notifyObj.ExecuteCall.Returns.Error = errors.New("the error")

				handler.ServeHTTP(writer, request, context)
				Expect(errorWriter.WriteCall.Receives.Error).To(Equal(notifyObj.ExecuteCall.Returns.Error))
			})
		})
	})
})