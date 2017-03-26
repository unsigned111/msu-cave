package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/hypebeast/go-osc/osc"
)

type AppArgs struct {
	ReplayFileName  string
	Port            int
        Delay           int
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
	flag.Parse()

	args := AppArgs{
		ReplayFileName: *fileName,
		Port:           *port,
                Delay:          *delay,
	}

	fmt.Println("Reading from:", args.ReplayFileName)
	fmt.Println("Destination Port:", args.Port)
        fmt.Println("Delay:", args.Delay)
	return args
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	// parse the command line arguments
	args := parseArgs()

	// open the file
	file, err := os.Open(args.ReplayFileName)
	check(err)
	defer file.Close()

	client := osc.NewClient("localhost", args.Port)
	reader := csv.NewReader(file)
	// read header
	reader.Read()
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		check(err)
		fmt.Println(record)

		msg := osc.NewMessage("/eeg")
		for value := range record {
			i, err := strconv.Atoi(record[value])
			check(err)
			msg.Append(int32(i))
		}
		client.Send(msg)
		time.Sleep(time.Duration(args.Delay) * time.Second)
	}
}
