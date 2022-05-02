package main

import "os"

const HELP_TEXT = `
_______________________________________________________________

  MyInsights CLI Client
  CPAD G-53

  Version 0.0.1
	
_______________________________________________________________

Usage:
  cli-client <command> [args]

Commands:
  end              Send an event to the server
  last             Get the last event from the server
  conf             Show current configuration

Configuration:
  INSIGHT_SERVER   Env variable which holds the server URL
                   Example: http://localhost:5000
                   (no trailing slash)

_______________________________________________________________

`


func show_help () {
	println(HELP_TEXT)
	os.Exit(1)
}


func main () {
	if len(os.Args) == 1 {
		show_help()
	} else {
		address := os.Getenv("INSIGHT_SERVER")
		if address == "" {
			println("ERROR: INSIGHT_SERVER environment variable not set!")
			println("       This must be configured as a fully qualified URL")
			os.Exit(1)
		}
		switch os.Args[1] {
		case "send":
			println("Sending ...")
		case "last":
			println("Getting ...")
		}
	}
}