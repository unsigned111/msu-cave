package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
	"strings"

	"github.com/hypebeast/go-osc/osc"
)

type AppArgs struct {
	ReplayFileName  	string
	Port			int
	Delay		   	int
	Toggle			bool
}

func parseArgs() AppArgs {
	fileName := flag.String(
		"f",
		"default.csv",
		"Brainwave log file from which to play back",
	)
	port := flag.Int(
		"p",
		7770,
		"Port to send data",
	)
	delay := flag.Int(
		"d",
		1,
		"Delay between messages",
	)
	toggle := flag.Bool(
		"t",
		true,
		"Toggle headset on/off enabled",
	)
	flag.Parse()

	args := AppArgs{
		ReplayFileName:	*fileName,
		Port:		*port,
		Delay:		*delay,
		Toggle:		*toggle,
	}

	fmt.Println("Reading from:", args.ReplayFileName)
	fmt.Println("Destination Port:", args.Port)
	fmt.Println("Delay:", args.Delay)
	fmt.Println("Headset Toggle On/Off:", args.Toggle)
	return args
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

var args AppArgs
var messageChannel = make(chan osc.Packet)

func main() {
	// parse the command line arguments
	args = parseArgs()

	client := osc.NewClient("localhost", args.Port)
	go eegSender()
	go toggleSender()

	for {
		client.Send(<-messageChannel)
	}
}

func eegSender() {
	endpoint := "/eeg"

	// open the file
	file, err := os.Open(args.ReplayFileName)
	check(err)
	defer file.Close()
	
	reader := csv.NewReader(file)
	// read header
	reader.Read()
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		check(err)
		msg := osc.NewMessage(endpoint)
		for value := range record {
			i, err := strconv.Atoi(record[value])
			check(err)
			msg.Append(int32(i))
		}
		messageChannel <- msg
		fmt.Println(endpoint, ":\t", strings.Join(record, " "))
		time.Sleep(time.Duration(args.Delay) * time.Second)
	}
}

func toggleSender() {
	endpoint := "/onoff"
	onoff := false
	for {
		var value int
		if onoff {
			value = 1
		} else {
			value = 0
		}
		msg := osc.NewMessage(endpoint)
		msg.Append(int32(value))
		messageChannel <- msg
		onoff = !onoff
		fmt.Println(endpoint, ":\t", value)
		time.Sleep(time.Duration(args.Delay * 10) * time.Second)
	}
}
