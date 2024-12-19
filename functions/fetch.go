package functions

import (
	"fmt"
	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapclient"
	"github.com/emersion/go-message/charset"
	"log"
	"mime"
	"os"
	"strings"
)

type Email struct {
	From    string
	Subject string
	Body    string
	For     string
}

func StartFetch() {
	args := os.Args[1:]
	fetchFunction := "help"
	if len(args) > 1 {
		fetchFunction = args[1]
	} else {
		fetchFunction = "nothing"
	}
	switch fetchFunction {
	case "help":
		fmt.Println("Available fetch functions:")
		fmt.Println("    help: show this help message")
		fmt.Println("Usage of fetch:")
		fmt.Println("    <profile_name>")
		fmt.Println("    leave empty for all profiles")
		os.Exit(0)
	}
	userWantsProfile := false
	profileToFetch := ""
	if len(args) > 1 {
		userWantsProfile = true
		profileToFetch = args[1]
	}
	imapProfiles := CheckImapProfilesFile()
	collectionOfMessages := []Email{}
	for _, profile := range imapProfiles {
		if userWantsProfile {
			if profile.(map[string]interface{})["name"] == profileToFetch {
				collectionOfMessages = append(fetchEmails(profile.(map[string]interface{}), 20), collectionOfMessages...)
				break
			}
		} else {
			collectionOfMessages = append(fetchEmails(profile.(map[string]interface{}), 20), collectionOfMessages...)
		}
	}
	fmt.Println("| For | From | Subject | Body |")
	for _, message := range collectionOfMessages {
		fmt.Println("|", message.For, "|", message.From, "|", message.Subject, "|", message.Body, "|")
	}
}

func fetchEmails(profile map[string]interface{}, numberOfEmailsToFetch int) []Email {
	options := &imapclient.Options{
		WordDecoder: &mime.WordDecoder{CharsetReader: charset.Reader},
	}
	client, err := imapclient.DialTLS(profile["server"].(string)+":"+profile["port"].(string), options)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		os.Exit(1)
	}
	defer client.Close()
	if err := client.Login(profile["username"].(string), profile["password"].(string)).Wait(); err != nil {
		log.Fatalf("Failed to log in: %v", err)
		os.Exit(1)
	}
	defer client.Logout()

	selectedMbox, err := client.Select("INBOX", nil).Wait()
	if err != nil {
		fmt.Printf("Failed to select INBOX: %v\n", err)
		os.Exit(1)
	}

	if selectedMbox.NumMessages == 0 {
		fmt.Println("No messages to fetch")
		return []Email{}
	}

	seqSet := imap.SeqSet{}
	convertThatToUint := uint32(numberOfEmailsToFetch) + 1
	seqSet.AddRange(selectedMbox.NumMessages-convertThatToUint, selectedMbox.NumMessages)
	fetchOptions := &imap.FetchOptions{
		Envelope:      true,
		Flags:         true,
		BodyStructure: &imap.FetchItemBodyStructure{Extended: true},
		BodySection:   []*imap.FetchItemBodySection{{Peek: true, Part: []int{1}}},
	}

	messagesToReturn := []Email{}

	messages, err := client.Fetch(seqSet, fetchOptions).Collect()
	if err != nil {
		fmt.Printf("Failed to fetch messages in INBOX: %v\n", err)
		os.Exit(1)
	}
	if len(messages) > 0 {
		for _, msg := range messages {
			for _, section := range msg.BodySection {
				if section != nil {
					bodyContent := string(section)
					bodyToReturn := bodyContent
					if msg.BodyStructure.MediaType() == "text/html" {
						bodyToReturn = "HTML Content"
					}
					if len(bodyContent) > 0 {
						if strings.Contains(bodyContent, "<head>") {
							bodyToReturn = "HTML Content"
						}
						if strings.Contains(bodyContent, "<!DOCTYPE html>") {
							bodyToReturn = "HTML Content"
						}
						if strings.Contains(bodyContent, "<html>") {
							bodyToReturn = "HTML Content"
						}
						bodyToReturn = strings.ReplaceAll(bodyToReturn, "\n", "")
						bodyToReturn = strings.ReplaceAll(bodyToReturn, "\r", "")
						bodyToReturn = strings.ReplaceAll(bodyToReturn, "|", "")
					}
					messagesToReturn = append(messagesToReturn, Email{
						From:    msg.Envelope.From[0].Name,
						Subject: msg.Envelope.Subject,
						Body:    bodyToReturn,
						For:     profile["name"].(string),
					})
				}
			}
		}
	}
	return messagesToReturn
}
