package uaa

import (
	"fmt"
	"net/url"

	"github.com/pivotal-cf-experimental/warrant"
	uaaSSOGolang "github.com/pivotal-cf/uaa-sso-golang/uaa"
)

type ZonedUAAClient struct {
	clientID       string
	clientSecret   string
	verifySSL      bool
	tokenValidator *TokenValidator
}

func NewZonedUAAClient(clientID, clientSecret string, verifySSL bool, validator *TokenValidator) (client ZonedUAAClient) {
	return ZonedUAAClient{
		clientID:       clientID,
		clientSecret:   clientSecret,
		verifySSL:      verifySSL,
		tokenValidator: validator,
	}
}

func (z ZonedUAAClient) GetTokenKey(uaaHost string) (string, error) {
	uaaClient := warrant.New(warrant.Config{
		Host:          uaaHost,
		SkipVerifySSL: !z.verifySSL,
	})

	signingKey, err := uaaClient.Tokens.GetSigningKey()
	if err != nil {
		return "", err
	}

	return signingKey.Value, nil
}

func (z ZonedUAAClient) GetClientToken(host string) (string, error) {
	uaaClient := warrant.New(warrant.Config{
		Host:          host,
		SkipVerifySSL: !z.verifySSL,
	})

	return uaaClient.Clients.GetToken(z.clientID, z.clientSecret)
}

func (z ZonedUAAClient) UsersEmailsByIDs(token string, ids ...string) ([]User, error) {
	uaaHost, err := z.tokenHost(token)
	if err != nil {
		return nil, err
	}

	uaaClient := uaaSSOGolang.NewUAA("", uaaHost, z.clientID, z.clientSecret, "")
	uaaClient.VerifySSL = z.verifySSL
	uaaClient.SetToken(token)

	var myUsers []User
	users, err := uaaClient.UsersEmailsByIDs(ids...)
	if err != nil {
		return myUsers, err
	}

	for _, user := range users {
		myUsers = append(myUsers, newUserFromSSOGolangUser(user))
	}

	return myUsers, nil
}

func (z ZonedUAAClient) tokenHost(token string) (string, error) {
	parsedToken, err := z.tokenValidator.Parse(token)

	if err != nil {
		return "", err
	}

	tokenIssuerURL, err := url.Parse(parsedToken.Claims["iss"].(string))
	if err != nil {
		return "", err
	}

	return tokenIssuerURL.Scheme + "://" + tokenIssuerURL.Host, nil
}

func (z ZonedUAAClient) AllUsers(token string) ([]User, error) {
	uaaHost, err := z.tokenHost(token)
	if err != nil {
		return nil, err
	}

	uaaSSOGolangClient := uaaSSOGolang.NewUAA("", uaaHost, z.clientID, z.clientSecret, "")
	uaaSSOGolangClient.VerifySSL = z.verifySSL
	uaaSSOGolangClient.SetToken(token)
	users, err := uaaSSOGolangClient.AllUsers()

	var myUsers []User
	for _, user := range users {
		myUsers = append(myUsers, newUserFromSSOGolangUser(user))
	}

	return myUsers, err
}

func (z ZonedUAAClient) UsersGUIDsByScope(token string, scope string) ([]string, error) {
	uaaHost, err := z.tokenHost(token)
	if err != nil {
		return nil, err
	}

	uaaSSOGolangClient := uaaSSOGolang.NewUAA("", uaaHost, z.clientID, z.clientSecret, "")
	uaaSSOGolangClient.VerifySSL = z.verifySSL

	return uaaSSOGolangClient.UsersGUIDsByScope(scope)
}

func newUserFromWarrantUser(warrantUser warrant.User) User {
	user := User{}
	user.ID = warrantUser.ID
	user.Emails = warrantUser.Emails

	return user
}

func newUserFromSSOGolangUser(uaaUser uaaSSOGolang.User) User {
	user := User{}
	user.ID = uaaUser.ID
	user.Emails = uaaUser.Emails

	return user
}

type User struct {
	ID     string
	Emails []string
}

type Failure struct {
	code    int
	message string
}

func NewFailure(code int, message []byte) Failure {
	return Failure{
		code:    code,
		message: string(message),
	}
}

func (failure Failure) Code() int {
	return failure.code
}

func (failure Failure) Message() string {
	return failure.message
}

func (failure Failure) Error() string {
	return fmt.Sprintf("UAA Wrapper Failure: %d %s", failure.code, failure.message)
}
