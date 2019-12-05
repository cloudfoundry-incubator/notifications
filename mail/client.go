package mail

import (
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"net"
	"net/smtp"
	"strings"
	"time"

	"github.com/pivotal-golang/lager"
)

const (
	SMTPAuthNone    = "none"
	SMTPAuthPlain   = "plain"
	SMTPAuthCRAMMD5 = "cram-md5"
	SMTPAuthLogin   = "login"
)

var SMTPAuthMechanisms = []string{SMTPAuthNone, SMTPAuthPlain, SMTPAuthCRAMMD5, SMTPAuthLogin}

type AuthMechanism int

type Client struct {
	config Config
	client *smtp.Client
}

type Config struct {
	Host              string
	Port              string
	User              string
	Pass              string
	Secret            string
	SMTPAuthMechanism string
	TestMode          bool
	SkipVerifySSL     bool
	DisableTLS        bool
	ConnectTimeout    time.Duration
	LoggingEnabled    bool
}

type connection struct {
	client *smtp.Client
	err    error
}

// login auth
type loginAuth struct {
	username, password string
}

func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte{}, nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, errors.New("Unkown fromServer")
		}
	}
	return nil, nil
}

// login auth

func NewClient(config Config) *Client {
	client := &Client{config: config}

	if client.config.ConnectTimeout == 0 {
		client.config.ConnectTimeout = 15 * time.Second
	}

	return client
}

func (c Client) createLoggerSession(logger lager.Logger) lager.Logger {
	if strings.HasSuffix(logger.SessionName(), ".smtp") {
		return logger
	}

	return logger.Session("smtp")
}

func (c *Client) Connect(logger lager.Logger) error {
	logger = c.createLoggerSession(logger)

	c.PrintLog(logger, "connecting")
	if c.config.TestMode {
		c.PrintLog(logger, "test-mode-not-connected")
		return nil
	}

	if c.client != nil {
		c.PrintLog(logger, "existing-connection")
		return nil
	}

	select {
	case connection := <-c.connect():
		c.PrintLog(logger, "connected")
		if connection.err != nil {
			return connection.err
		}

		c.client = connection.client
	case <-time.After(c.config.ConnectTimeout):
		c.PrintLog(logger, "connection-timeout", lager.Data{"timeout-duration": c.config.ConnectTimeout})
		return errors.New("server timeout")
	}

	return nil
}

func (c *Client) connect() chan connection {
	channel := make(chan connection)

	go func() {
		client, err := smtp.Dial(net.JoinHostPort(c.config.Host, c.config.Port))
		channel <- connection{
			client: client,
			err:    err,
		}
	}()

	return channel
}

func (c *Client) Send(msg Message, logger lager.Logger) error {
	logger = c.createLoggerSession(logger)

	if c.config.TestMode {
		logger.Info("test-mode")
		return nil
	}

	err := c.Connect(logger)
	if err != nil {
		return c.Error(logger, err)
	}

	c.PrintLog(logger, "hello-initiating")
	err = c.Hello()
	if err != nil {
		return c.Error(logger, err)
	}
	c.PrintLog(logger, "hello-complete")

	if !c.config.DisableTLS {
		c.PrintLog(logger, "tls-starting")
		err = c.StartTLS()
		if err != nil {
			return c.Error(logger, err)
		}
		c.PrintLog(logger, "tls-connected")

		c.PrintLog(logger, "authentication-starting")
		err = c.Auth(logger)
		if err != nil {
			return c.Error(logger, err)
		}
		c.PrintLog(logger, "authenticated")
	}

	c.PrintLog(logger, "setting-msg-from", lager.Data{"from": msg.From})
	err = c.client.Mail(msg.From)
	if err != nil {
		return c.Error(logger, err)
	}

	c.PrintLog(logger, "setting-msg-to", lager.Data{"to": msg.To})
	err = c.client.Rcpt(msg.To)
	if err != nil {
		return c.Error(logger, err)
	}

	c.PrintLog(logger, "setting-msg-data", lager.Data{"message-data": base64.StdEncoding.EncodeToString([]byte(msg.Data()))})
	err = c.Data(msg)
	if err != nil {
		return c.Error(logger, err)
	}
	c.PrintLog(logger, "msg-data-sent")

	c.PrintLog(logger, "quiting")
	err = c.Quit()
	if err != nil {
		return c.Error(logger, err)
	}
	c.PrintLog(logger, "disconnected")

	return nil
}

func (c *Client) Hello() error {
	err := c.client.Hello("localhost")
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Extension(name string) (bool, string) {
	return c.client.Extension(name)
}

func (c *Client) StartTLS() error {
	if ok, _ := c.Extension("STARTTLS"); ok {
		err := c.client.StartTLS(&tls.Config{
			ServerName:         c.config.Host,
			InsecureSkipVerify: c.config.SkipVerifySSL,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) Auth(logger lager.Logger) error {
	if ok, _ := c.Extension("AUTH"); ok {
		if mechanism := c.AuthMechanism(logger); mechanism != nil {
			err := c.client.Auth(mechanism)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *Client) AuthMechanism(logger lager.Logger) smtp.Auth {
	switch c.config.SMTPAuthMechanism {
	case SMTPAuthCRAMMD5:
		c.PrintLog(logger, "crammd5-authentication")
		return smtp.CRAMMD5Auth(c.config.User, c.config.Secret)
	case SMTPAuthPlain:
		c.PrintLog(logger, "plain-authentication")
		return smtp.PlainAuth("", c.config.User, c.config.Pass, c.config.Host)
	case SMTPAuthLogin:
		c.PrintLog(logger, "login-authentication")
		return LoginAuth(c.config.User, c.config.Pass)
	default:
		c.PrintLog(logger, "no-authentication")
		return nil
	}
}

func (c *Client) Data(msg Message) error {
	wc, err := c.client.Data()
	if err != nil {
		return err
	}

	data := strings.Replace(string(msg.Data()), "%", "%%", -1)
	_, err = fmt.Fprintf(wc, data)
	if err != nil {
		return err
	}

	err = wc.Close()
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Quit() error {
	err := c.client.Quit()
	c.client = nil
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Error(logger lager.Logger, err error) error {
	if c.client != nil {
		failure := c.Quit()
		if failure != nil {
			return failure
		}
	}

	logger.Error("failed", err)

	return err
}

func (c *Client) PrintLog(logger lager.Logger, action string, data ...lager.Data) {
	if c.config.LoggingEnabled {
		logger.Info(action, data...)
	}
}
