package functions

import (
	"crypto/tls"
	"fmt"
	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
	"os"
	"strings"
)

func StartSend() {
	args := os.Args[1:]
	if len(args) < 2 {
		fmt.Println("Incomplete command.. try 'help'")
		return
	}
	sendFunction := args[1]
	switch sendFunction {
	case "help":
		fmt.Println("Available send functions:")
		fmt.Println("    help: show this help message")
		fmt.Println("Usage of send:")
		fmt.Println("    <profile_name> <receiving_address> \"<subject>\" \"<message>\"")
		os.Exit(0)
	}
	if len(args) < 5 {
		fmt.Println("Incomplete command.. try 'help'")
		return
	}
	profileName := args[1]
	receivingAddress := args[2]
	subject := args[3]
	message := args[4]
	getProfiles := CheckProfilesFile()
	ourProfile := make(map[string]interface{})
	for _, profile := range getProfiles {
		if profile.(map[string]interface{})["name"] == profileName {
			ourProfile = profile.(map[string]interface{})
			break
		}
	}
	if len(ourProfile) == 0 {
		fmt.Println("profile not found.. try 'profile list'")
	}
	fmt.Println("you want to send an email with the following details:")
	fmt.Println("from:", ourProfile["email"])
	fmt.Println("to:", receivingAddress)
	fmt.Println("subject:", subject)
	fmt.Println("message:", message)
	fmt.Println("is this correct? (y/n)")
	var confirmation string
	fmt.Scanln(&confirmation)
	if confirmation == "n" || confirmation == "no" {
		fmt.Println("email sending cancelled")
		return
	}
	fmt.Println("sending..")
	senderEmail := ourProfile["email"].(string)
	senderPassword := ourProfile["password"].(string)
	senderDisplayName := ourProfile["displayName"].(string)
	smtpServer := ourProfile["smtpServer"].(string)
	smtpPort := ourProfile["smtpPort"].(string)
	tlsConfig := &tls.Config{
		ServerName: smtpServer,
	}

	client, err := smtp.DialTLS(smtpServer+":"+smtpPort, tlsConfig)
	if err != nil {
		fmt.Printf("Failed to connect: %v", err)
	}
	defer client.Close()

	auth := sasl.NewPlainClient("", senderEmail, senderPassword)
	toUsers := []string{receivingAddress}
	msgToSend := strings.NewReader("From: " + senderDisplayName + " <" + senderEmail + ">\r\n" +
		"To: " + receivingAddress + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		message + "\r\n")
	err2 := smtp.SendMailTLS(smtpServer+":"+smtpPort, auth, senderEmail, toUsers, msgToSend)
	if err2 != nil {
		fmt.Println("Error sending email:", err2)
		os.Exit(1)
	}
	fmt.Println("email sent successfully!")
}
