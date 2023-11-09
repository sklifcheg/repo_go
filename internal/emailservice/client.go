package emailservice

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
)

type Client struct {
	Client *smtp.Client
}

func NewClient(host, username, password, portNumber string) (*Client, error) {
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}
	auth := smtp.PlainAuth("", username, password, host)
	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%s", host, portNumber), tlsconfig)
	if err != nil {
		return nil, err
	}
	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return nil, err
	}
	// Auth
	if err = client.Auth(auth); err != nil {
		return nil, err
	}
	return &Client{client}, nil
}
