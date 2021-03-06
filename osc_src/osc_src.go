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
	"math/rand"

	"github.com/hypebeast/go-osc/osc"
)

type AppArgs struct {
	ReplayFileName  	string
	Port			int
	Delay		   	int
	Toggle			bool
	Similarity		bool
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
		false,
		"Toggle headset on/off simulation",
	)
	similarity := flag.Bool(
		"s",
		false,
		"Toggle similarity simulation",
	)
	flag.Parse()

	args := AppArgs{
		ReplayFileName:	*fileName,
		Port:		*port,
		Delay:		*delay,
		Toggle:		*toggle,
		Similarity:	*similarity,
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
	if args.Toggle {
		go toggleSender()
	}
	
	if args.Similarity {
		go similaritySender()
	}

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

        msg := osc.NewMessage("/onoff")
        msg.Append(int32(1))
        messageChannel <- msg
	
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
		fmt.Println(endpoint, ":\t\t", strings.Join(record, " "))
		time.Sleep(time.Duration(args.Delay) * time.Second)
	}
}

func toggleSender() {
	endpoint := "/onoff"
	onoff := true
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
		time.Sleep(time.Duration(args.Delay * 20) * time.Second)
	}
}

func similaritySender() {
	endpoint := "/similarity"
	for {
		msg := osc.NewMessage(endpoint)
		value := rand.Float64()
		msg.Append(value)
		messageChannel <- msg
		fmt.Println(endpoint, ":\t", value)
		time.Sleep(time.Duration(args.Delay * 5) * time.Second)
	}
}
