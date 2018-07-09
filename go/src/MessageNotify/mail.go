package main

import (
	"fmt"
	"net/smtp"
	"strings"
)

// "Content-Type: text/"+ mailtype + "; charset=UTF-8"
//

type Mail struct {
	Host     string "exmail.boyaa.com"
	Port     int
	From     string "noreply@boyaamail.com"
	To       string "AustinChen@boyaa.com"
	Msg      string "Test Mail"
	Subject  string "Test"
	Encoding string "Content-Type: text/html; charset=UTF-8"
	User     string "queren"
	Password string "jakquan"
}

func (m *Mail) Send() error {
	a := smtp.PlainAuth("", m.User, m.Password, "exmail.boyaa.com")
	msg := []byte("To: " + m.To +
		"\r\nFrom: " +
		m.User +
		"<" + m.User +
		">\r\nSubject: " +
		m.Subject + "\r\n" +
		m.Encoding +
		"\r\n\r\n" +
		m.Msg)

	err := smtp.SendMail("exmail.boyaa.com:25", a, m.User, strings.Split(m.To, ";"), msg)

	if err == nil {
		fmt.Printf("send success.\n")
	} else {
		fmt.Printf("send failed: %v\n", err)
	}

	return nil
}
