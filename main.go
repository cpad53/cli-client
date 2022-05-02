package main

import "fmt"
import "bytes"
import "os"
import "strings"
import "errors"
import "time"
import "encoding/json"
import "net/http"
import "io/ioutil"
import "strconv"


const HELP_TEXT = `
_______________________________________________________________

  MyInsights CLI Client
  CPAD G-53

  Version 0.0.1
	
_______________________________________________________________

Usage:
  cli-client <command> [args]

Commands:
  send             Send an event to the server
  last             Get the last event from the server
  conf             Show current configuration

Configuration:
  INSIGHT_SERVER   Env variable which holds the server URL
                   Example: http://localhost:5000
                   (no trailing slash)

_______________________________________________________________

`


type EventPayload struct {
	Ts string `json:"ts"`
	Msg string `json:"msg"`
}



func show_help () {
	println(HELP_TEXT)
	os.Exit(1)
}


func sendEvent(address string, args... string) (bool, error) {
	if len(args) == 0 {
		return false, errors.New("ERROR: No input received\n")
	}
	println("Sending ...")
	event := &EventPayload{
		Ts: strconv.FormatInt(time.Now().UnixMilli(), 10),
		Msg: strings.Join(args, " "),
	}
	jMsg, err := json.Marshal(event)
	if err == nil {
		target := address + "/event"
		jBytes := bytes.NewBuffer(jMsg)
		resp, err := http.Post(target, "application/json", jBytes)
		if err != nil {
			return false, err
		} else {
			println("StatusCode:", resp.StatusCode)
		}
	} else {
		return false, err
	}
	return true, nil
}


func getLastEvent(address string) (bool, error) {
	target := address + "/lastEvent"
	resp, err := http.Get(target)
	if err != nil {
		return false, err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	// println(string(data))
	var event EventPayload
	err = json.Unmarshal(data, &event)
	if err != nil {
		return false, err
	}
	hr_suffix := ""
	ts_int64, err := strconv.ParseInt(event.Ts, 10, 64)
	if err == nil {
		ts_time := time.UnixMilli(ts_int64)
		hr_suffix = "( " + ts_time.Format(time.RFC850) + " )"
	} else {
		fmt.Println(err.Error())
	}
	fmt.Println("Timestamp:", event.Ts, hr_suffix)
	fmt.Println("Event Msg:", event.Msg)
	return true, nil
}


func showConfig() {
	println("INSIGHT_SERVER Address Is:", os.Getenv("INSIGHT_SERVER"))
	os.Exit(0)
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
			ok, err := sendEvent(address, os.Args[2:]...)
			if !ok {
				println(err.Error())
				os.Exit(1)
			}
		case "last":
			ok, err := getLastEvent(address)
			if !ok {
				println(err.Error())
				os.Exit(1)
			}
		}
	}
}