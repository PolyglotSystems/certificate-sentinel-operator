/*
Copyright 2021 Polyglot Systems.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package helpers

import (
	"crypto/tls"
	"log"
	"net"
	"strconv"
	"time"

	defaults "github.com/kenmoini/certificate-sentinel-operator/controllers/defaults"
	mail "github.com/xhit/go-simple-mail/v2"
)

// SendSMTPMail assembles everything needed to sent an email via go-simple-mail
func SendSMTPMail(authType string, username string, password string, identity string, cramSecret string, useSSL *bool, useSTARTTLS *bool, to []string, from string, smtpServer string, subject string, textMessage string, htmlMessage string) {

	// Create a new SMTP Client
	server := mail.NewSMTPClient()

	// Set up the Server Connection
	smtpHost, smtpPort, _ := net.SplitHostPort(smtpServer)
	server.Host = smtpHost
	server.Port, _ = strconv.Atoi(smtpPort)

	// Set Authentication
	switch authType {
	case "none":
		server.Authentication = mail.AuthNone
	case "cram-md5":
		server.Authentication = mail.AuthCRAMMD5
		server.Username = username
		server.Password = cramSecret
	case "login":
		server.Authentication = mail.AuthLogin
		server.Username = username
		server.Password = password
	case "plain":
		server.Authentication = mail.AuthPlain
		server.Username = username
		server.Password = password
	default:
		server.Authentication = mail.AuthPlain
		server.Username = username
		server.Password = password
	}

	// Set other server connection variables
	// Variable to keep alive connection
	server.KeepAlive = false
	// Timeout for connect to SMTP Server
	server.ConnectTimeout = 10 * time.Second
	// Timeout for send the data and wait respond
	server.SendTimeout = 10 * time.Second

	server.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Set STARTTLS config
	if *useSTARTTLS == true {
		server.Encryption = mail.EncryptionSTARTTLS
	}
	// Set TLSConfig to provide custom TLS configuration. For example, to skip TLS verification (useful for testing):
	if *useSSL == true {
		server.Encryption = mail.EncryptionSSLTLS
		server.TLSConfig = &tls.Config{InsecureSkipVerify: false}
	}

	// Create SMTP client
	smtpClient, err := server.Connect()
	if err != nil {
		log.Printf("SMTP CLIENT CREATION ERROR! %v\n", err)
	}

	// Set up new email message
	email := mail.NewMSG()
	emailSubject := defaults.SetDefaultString(defaults.SMTPMessageSubject, subject)

	// Set message details
	email.SetFrom(from).
		SetSubject(emailSubject)

	// loop through senders, add to email message
	for _, addy := range to {
		email.AddTo(addy)
	}

	// Set additional message headers

	// Set the message body
	if htmlMessage != "" {
		email.SetBody(mail.TextHTML, htmlMessage)
	}
	/*
		if textMessage != "" {
			email.AddAlternative(mail.TextPlain, textMessage)
		}
	*/

	// Send message
	// always check error before send (y tho?)
	if email.Error != nil {
		log.Printf("ERROR SENDING SMTP MESSAGE! %v\n", email.Error)
	}

	// Call Send and pass the client
	err = email.Send(smtpClient)
	if err != nil {
		log.Println(err)
	}
	// Email sent!

}
