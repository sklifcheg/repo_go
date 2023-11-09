package emailservice

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

type Message struct {
	From        string
	To          []string
	CC          []string
	BCC         []string
	Subject     string
	Body        string
	Attachments map[string][]byte
}

type EmailService struct {
	client Client
}

func NewEmailService(client Client) *EmailService {
	return &EmailService{
		client: client,
	}
}

func NewMessage(s, b string) *Message {
	return &Message{Subject: s, Body: b, Attachments: make(map[string][]byte)}
}

func (srv *EmailService) SendMessage(m *Message) error {
	err := srv.client.Client.Mail(m.From)
	if err != nil {
		return err
	}

	for _, to := range m.To {
		err = srv.client.Client.Rcpt(to)
		if err != nil {
			return err
		}
	}

	w, err := srv.client.Client.Data()
	if err != nil {
		return err
	}

	_, err = w.Write(m.ToBytes())
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	err = srv.client.Client.Quit()
	if err != nil {
		return err
	}

	return nil
}

func (m *Message) AttachFile(src string) error {
	b, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	_, fileName := filepath.Split(src)
	m.Attachments[fileName] = b
	return nil
}

func (m *Message) ToBytes() []byte {
	buf := bytes.NewBuffer(nil)
	withAttachments := len(m.Attachments) > 0
	buf.WriteString(fmt.Sprintf("From: %s\r\n", m.From))
	buf.WriteString(fmt.Sprintf("Subject: %s\r\n", m.Subject))
	buf.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(m.To, ";")))
	if len(m.CC) > 0 {
		buf.WriteString(fmt.Sprintf("Cc: %s\r\n", strings.Join(m.CC, ";")))
	}

	if len(m.BCC) > 0 {
		buf.WriteString(fmt.Sprintf("Bcc: %s\r\n", strings.Join(m.BCC, ";")))
	}

	buf.WriteString("MIME-Version: 1.0\r\n")
	writer := multipart.NewWriter(buf)
	boundary := writer.Boundary()
	if withAttachments {
		buf.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=\"%s\"\r\n", boundary))
		buf.WriteString(fmt.Sprintf("\r\n--%s\n", boundary))
	} else {
		buf.WriteString("Content-Type: text/plain; charset=utf-8\n")
	}

	buf.WriteString("Content-Type: text/html; charset=\"utf-8\"\r\n")
	buf.WriteString("Content-Transfer-Encoding: 7bit\r\n")

	buf.WriteString(m.Body)
	if withAttachments {
		for k, v := range m.Attachments {
			buf.WriteString(fmt.Sprintf("\r\n--%s\n", boundary))
			buf.WriteString("Content-Type: text/plain; charset=\"utf-8\"\r\n")
			buf.WriteString("Content-Transfer-Encoding: base64\n")
			buf.WriteString("Content-Disposition: attachment;filename=\"" + k + "\"\r\n")
			buf.WriteString("\r\n")
			b := make([]byte, base64.StdEncoding.EncodedLen(len(v)))
			base64.StdEncoding.Encode(b, v)
			buf.Write(b)
			buf.WriteString(fmt.Sprintf("\n--%s", boundary))
		}

		buf.WriteString("--")
	}

	return buf.Bytes()
}
