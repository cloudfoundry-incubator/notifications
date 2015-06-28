package strategies_test

import (
	"errors"
	"time"

	"github.com/cloudfoundry-incubator/notifications/cf"
	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/postal"
	"github.com/cloudfoundry-incubator/notifications/postal/strategies"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Everyone Strategy", func() {
	var (
		strategy            strategies.EveryoneStrategy
		tokenLoader         *fakes.TokenLoader
		allUsers            *fakes.AllUsers
		enqueuer            *fakes.Enqueuer
		conn                *fakes.Connection
		requestReceivedTime time.Time
	)

	BeforeEach(func() {
		requestReceivedTime, _ = time.Parse(time.RFC3339Nano, "2015-06-08T14:38:03.180764129-07:00")
		conn = fakes.NewConnection()
		tokenHeader := map[string]interface{}{
			"alg": "FAST",
		}
		tokenClaims := map[string]interface{}{
			"client_id": "mister-client",
			"exp":       int64(3404281214),
			"scope":     []string{"notifications.write"},
		}
		tokenLoader = fakes.NewTokenLoader()
		tokenLoader.Token = fakes.BuildToken(tokenHeader, tokenClaims)
		enqueuer = fakes.NewEnqueuer()
		allUsers = fakes.NewAllUsers()
		allUsers.AllUserGUIDsCall.Returns = []string{"user-380", "user-319"}
		strategy = strategies.NewEveryoneStrategy(tokenLoader, allUsers, enqueuer)
	})

	Describe("Dispatch", func() {
		It("call enqueuer.Enqueue with the correct arguments for an organization", func() {
			_, err := strategy.Dispatch(strategies.Dispatch{
				Connection: conn,
				Kind: strategies.Kind{
					ID:          "welcome_user",
					Description: "Your Official Welcome",
				},
				Client: strategies.Client{
					ID:          "my-client",
					Description: "Welcome system",
				},
				Message: strategies.Message{
					ReplyTo: "reply-to@example.com",
					Subject: "this is the subject",
					To:      "dr@strangelove.com",
					Text:    "Welcome to the system, now get off my lawn.",
					HTML: strategies.HTML{
						BodyContent:    "<p>Welcome to the system, now get off my lawn.</p>",
						BodyAttributes: "some-html-body-attributes",
						Head:           "<head></head>",
						Doctype:        "<html>",
					},
				},
				VCAPRequest: strategies.VCAPRequest{
					ID:          "some-vcap-request-id",
					ReceiptTime: requestReceivedTime,
				},
			})
			Expect(err).NotTo(HaveOccurred())

			var users []strategies.User
			for _, guid := range allUsers.AllUserGUIDsCall.Returns {
				users = append(users, strategies.User{GUID: guid})
			}

			Expect(enqueuer.EnqueueCall.Args.Connection).To(Equal(conn))
			Expect(enqueuer.EnqueueCall.Args.Users).To(Equal(users))
			Expect(enqueuer.EnqueueCall.Args.Options).To(Equal(postal.Options{
				ReplyTo:           "reply-to@example.com",
				Subject:           "this is the subject",
				To:                "dr@strangelove.com",
				KindID:            "welcome_user",
				KindDescription:   "Your Official Welcome",
				SourceDescription: "Welcome system",
				Text:              "Welcome to the system, now get off my lawn.",
				HTML: postal.HTML{
					BodyContent:    "<p>Welcome to the system, now get off my lawn.</p>",
					BodyAttributes: "some-html-body-attributes",
					Head:           "<head></head>",
					Doctype:        "<html>",
				},
				Endorsement: strategies.EveryoneEndorsement,
			}))
			Expect(enqueuer.EnqueueCall.Args.Space).To(Equal(cf.CloudControllerSpace{}))
			Expect(enqueuer.EnqueueCall.Args.Org).To(Equal(cf.CloudControllerOrganization{}))
			Expect(enqueuer.EnqueueCall.Args.Client).To(Equal("my-client"))
			Expect(enqueuer.EnqueueCall.Args.Scope).To(Equal(""))
			Expect(enqueuer.EnqueueCall.Args.VCAPRequestID).To(Equal("some-vcap-request-id"))
			Expect(enqueuer.EnqueueCall.Args.RequestReceived).To(Equal(requestReceivedTime))
		})
	})

	Context("failure cases", func() {
		Context("when token loader fails to return a token", func() {
			It("returns an error", func() {
				tokenLoader.LoadError = errors.New("BOOM!")
				_, err := strategy.Dispatch(strategies.Dispatch{})

				Expect(err).To(Equal(errors.New("BOOM!")))
			})
		})

		Context("when allUsers fails to load users", func() {
			It("returns the error", func() {
				allUsers.AllUserGUIDsCall.Error = errors.New("BOOM!")
				_, err := strategy.Dispatch(strategies.Dispatch{})

				Expect(err).To(Equal(errors.New("BOOM!")))
			})
		})
	})
})
