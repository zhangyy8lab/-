package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
)

func sendCode(to, subject, body string) error {
	from := "your_mail@gmail.com"
	// application authPassword see https://myaccount.google.com/apppasswords
	password := "applicatin authPassword"

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	message := []byte("To: " + to + "\r\n" + "Subject: " + subject + "\r\n" + body + "\r\n")

	// base Auth
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// client
	client, err := smtp.Dial(smtpHost + ":" + smtpPort)
	defer client.Close()
	if err != nil {
		return err
	}

	// enable STARTTLS
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpHost,
	}

	if err = client.StartTLS(tlsConfig); err != nil {
		return err
	}

	if err = client.Auth(auth); err != nil {
		return err
	}

	// set send from
	if err = client.Mail(from); err != nil {
		return err
	}

	// set send to
	if err = client.Rcpt(to); err != nil {
		return err
	}

	// read from remote repository.
	w, err := client.Data()
	if err != nil {
		return err
	}

	// set send body
	_, err = w.Write(message)
	if err != nil {
		return err
	}

	if err = w.Close(); err != nil {
		return err
	}

	return client.Quit()
}

func main() {

	to := "send_to@gmail.com"
	subject := "Auth Code"
	body := "667db1c"

	if err := sendCode(to, subject, body); err != nil {
		log.Fatal(err)
	}

	fmt.Println("send code successful")

}
