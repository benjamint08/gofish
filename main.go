package main

import (
	"fmt"
	"gofish/functions"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println("Incomplete command.. try 'help'")
		return
	}
	function := args[0]
	switch function {
	case "profile":
		functions.StartProfile()
	case "send":
		functions.StartSend()
	case "fetch":
		functions.StartFetch()
	case "help":
		fmt.Println("available functions (run help on any function for help):")
		fmt.Println("    profile: manage email profiles")
		fmt.Println("    fetch: fetch new emails for all profiles")
		fmt.Println("    send: send a new email")
		fmt.Println("    help: show this help message")
	}
}
