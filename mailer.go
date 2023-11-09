package main

import (
	"fmt"
	"os"
	"profgo/internal/emailservice"
)

var (
	host       = os.Getenv("EMAIL_HOST")
	username   = os.Getenv("EMAIL_USERNAME")
	password   = os.Getenv("EMAIL_PASSWORD")
	portNumber = os.Getenv("EMAIL_PORT")
)

func main() {
	client, _ := emailservice.NewClient(host, username, password, portNumber)
	sender := emailservice.NewEmailService(*client)
	m := emailservice.NewMessage("Test program", "Golang mailsender")
	m.From = "Yuor email address"
	m.To = []string{"receiver1", "reciever2"}
	m.CC = []string{""}
	m.BCC = []string{""}
	m.AttachFile("path to file")
	fmt.Println(sender.SendMessage(m))
}
