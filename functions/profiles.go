package functions

import (
	"encoding/json"
	"fmt"
	"os"
)

func StartProfile() {
	args := os.Args[1:]
	if len(args) < 2 {
		fmt.Println("Incomplete command.. try 'help'")
		return
	}
	profileFunction := args[1]
	switch profileFunction {
	case "list":
		listProfiles()
	case "add":
		setupNewProfile()
	case "remove":
		removeProfile()
	case "help":
		fmt.Println("Available profile functions:")
		fmt.Println("    list: list all profiles")
		fmt.Println("    add: add a new profile")
		fmt.Println("    remove: remove a profile")
		fmt.Println("    help: show this help message")
	default:
		fmt.Println("Invalid profile function.. try 'help'")
	}
}

func removeProfile() {
	fmt.Println("please enter the name of the profile to remove:")
	var name string
	fmt.Scanln(&name)
	profiles := CheckProfilesFile()
	for i, profile := range profiles {
		if profile.(map[string]interface{})["name"] == name {
			profiles = append(profiles[:i], profiles[i+1:]...)
			profilesFile := os.Getenv("HOME") + "/.gofish/profiles.json"
			profilesData, err := json.Marshal(profiles)
			if err != nil {
				fmt.Println("Error marshalling profiles data:", err)
				os.Exit(1)
			}
			err = os.WriteFile(profilesFile, profilesData, 0644)
			if err != nil {
				fmt.Println("Error writing to profiles file:", err)
				os.Exit(1)
			}
			fmt.Println("profile removed successfully!")
			return
		}
	}
	fmt.Println("profile not found")
}

func setupNewProfile() {
	fmt.Println("welcome to gofish!")
	fmt.Println("please enter the unique profile name:")
	var name string
	fmt.Scanln(&name)
	fmt.Println("please enter your email address:")
	var email string
	fmt.Scanln(&email)
	profiles := CheckProfilesFile()
	for _, profile := range profiles {
		if profile.(map[string]interface{})["email"] == email {
			fmt.Println("Email already in use")
			return
		}
	}
	fmt.Println("please enter your display name:")
	var displayName string
	fmt.Scanln(&displayName)
	fmt.Println("please enter your password:")
	var password string
	fmt.Scanln(&password)
	fmt.Println("please enter your SMTP server:")
	var smtpServer string
	fmt.Scanln(&smtpServer)
	fmt.Println("please enter your SMTP port:")
	var smtpPort string
	fmt.Scanln(&smtpPort)
	newProfile := map[string]interface{}{
		"name":        name,
		"email":       email,
		"displayName": displayName,
		"password":    password,
		"smtpServer":  smtpServer,
		"smtpPort":    smtpPort,
	}
	profiles = append(profiles, newProfile)
	profilesFile := os.Getenv("HOME") + "/.gofish/profiles-smtp.json"
	profilesData, err := json.Marshal(profiles)
	if err != nil {
		fmt.Println("Error marshalling profiles data:", err)
		os.Exit(1)
	}
	err = os.WriteFile(profilesFile, profilesData, 0644)
	if err != nil {
		fmt.Println("Error writing to profiles file:", err)
		os.Exit(1)
	}
	fmt.Println("smtp profile added successfully!")
	fmt.Println("would you like to add imap for email fetching? (y/n)")
	var imapConfirmation string
	fmt.Scanln(&imapConfirmation)
	if imapConfirmation == "y" || imapConfirmation == "yes" {
		setupNewImapProfile(name)
	}
	fmt.Println("okay, you're all set! bye!")
}

func setupNewImapProfile(name string) {
	fmt.Println("please enter your IMAP server:")
	var imapServer string
	fmt.Scanln(&imapServer)
	fmt.Println("please enter your IMAP port:")
	var imapPort string
	fmt.Scanln(&imapPort)
	fmt.Println("please enter your IMAP username:")
	var imapUsername string
	fmt.Scanln(&imapUsername)
	fmt.Println("please enter your IMAP password:")
	var imapPassword string
	fmt.Scanln(&imapPassword)
	newImapProfile := map[string]interface{}{
		"name":     name,
		"server":   imapServer,
		"port":     imapPort,
		"username": imapUsername,
		"password": imapPassword,
	}
	profiles := CheckImapProfilesFile()
	profiles = append(profiles, newImapProfile)
	profilesFile := os.Getenv("HOME") + "/.gofish/profiles-imap.json"
	profilesData, err := json.Marshal(profiles)
	if err != nil {
		fmt.Println("Error marshalling profiles data:", err)
		os.Exit(1)
	}
	err = os.WriteFile(profilesFile, profilesData, 0644)
	if err != nil {
		fmt.Println("Error writing to profiles file:", err)
		os.Exit(1)
	}
	fmt.Println("imap profile added successfully!")
}

func CheckImapProfilesFile() []interface{} {
	userHome := os.Getenv("HOME")
	profilesDir := userHome + "/.gofish"
	profilesFile := profilesDir + "/profiles-imap.json"
	if _, err := os.Stat(profilesDir); os.IsNotExist(err) {
		err := os.Mkdir(profilesDir, 0755)
		if err != nil {
			fmt.Println("Error creating profiles directory:", err)
			os.Exit(1)
		}
	}
	if _, err := os.Stat(profilesFile); os.IsNotExist(err) {
		file, err := os.Create(profilesFile)
		if err != nil {
			fmt.Println("Error creating profiles file:", err)
			os.Exit(1)
		}
		defer file.Close()
		_, err = file.WriteString("[]")
		if err != nil {
			fmt.Println("Error writing to profiles file:", err)
			os.Exit(1)
		}
	}
	data, err := os.ReadFile(profilesFile)
	if err != nil {
		fmt.Println("Error reading profiles file:", err)
		os.Exit(1)
	}
	profiles := []interface{}{}
	json.Unmarshal(data, &profiles)
	return profiles
}

func CheckProfilesFile() []interface{} {
	userHome := os.Getenv("HOME")
	profilesDir := userHome + "/.gofish"
	profilesFile := profilesDir + "/profiles-smtp.json"
	if _, err := os.Stat(profilesDir); os.IsNotExist(err) {
		err := os.Mkdir(profilesDir, 0755)
		if err != nil {
			fmt.Println("Error creating profiles directory:", err)
			os.Exit(1)
		}
	}
	if _, err := os.Stat(profilesFile); os.IsNotExist(err) {
		file, err := os.Create(profilesFile)
		if err != nil {
			fmt.Println("Error creating profiles file:", err)
			os.Exit(1)
		}
		defer file.Close()
		_, err = file.WriteString("[]")
		if err != nil {
			fmt.Println("Error writing to profiles file:", err)
			os.Exit(1)
		}
	}
	data, err := os.ReadFile(profilesFile)
	if err != nil {
		fmt.Println("Error reading profiles file:", err)
		os.Exit(1)
	}
	profiles := []interface{}{}
	json.Unmarshal(data, &profiles)
	return profiles
}

func listProfiles() {
	profiles := CheckProfilesFile()
	fmt.Println("smtp:")
	for _, profile := range profiles {
		fmt.Println(profile.(map[string]interface{})["name"], ":", profile.(map[string]interface{})["email"], ":", profile.(map[string]interface{})["smtpServer"], ":", profile.(map[string]interface{})["smtpPort"])
	}
	if len(profiles) == 0 {
		fmt.Println("no smtp profiles found. try adding a new profile with 'add'")
	}
	imapProfiles := CheckImapProfilesFile()
	fmt.Println("imap:")
	for _, profile := range imapProfiles {
		fmt.Println(profile.(map[string]interface{})["name"], ":", profile.(map[string]interface{})["server"], ":", profile.(map[string]interface{})["port"], ":", profile.(map[string]interface{})["username"])
	}
	if len(imapProfiles) == 0 {
		fmt.Println("no imap profiles found. try adding a new profile with 'add'")
	}
}
