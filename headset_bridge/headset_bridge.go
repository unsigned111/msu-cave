package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"os"
	"os/signal"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/neurosky"
)

func eegToPayload(eeg neurosky.EEG) string {
	eegRawData := map[string]int{
		"delta":    eeg.Delta,
		"theta":    eeg.Theta,
		"loAlpha":  eeg.LoAlpha,
		"hiAlpha":  eeg.HiAlpha,
		"loBeta":   eeg.LoBeta,
		"hiBeta":   eeg.HiBeta,
		"loGamma":  eeg.LoGamma,
		"midGamma": eeg.MidGamma,
	}
	eegData, _ := json.Marshal(eegRawData)
	payload := string(eegData)
	fmt.Println(payload)
	return payload
}

func sendData(eeg neurosky.EEG, url string) {
	payload := eegToPayload(eeg)
	request := gorequest.New()
	response, body, errors := request.Post(url).
		Set("Notes", "gorequst is coming!").
		Send(payload).
		End()
	if len(errors) > 0 {
		fmt.Println("Error sending data: ", errors)
	} else if response.StatusCode != 200 {
		fmt.Println(response.Status, ": ", body)
	} else {
		fmt.Println("Request sent successfully")
	}
}

func makeRobot(device string, url string) *gobot.Robot {
	adaptor := neurosky.NewNeuroskyAdaptor("neurosky", device)
	neuro := neurosky.NewNeuroskyDriver(adaptor, "neuro")
	work := func() {
		gobot.On(neuro.Event("eeg"), func(data interface{}) {
			eeg := data.(neurosky.EEG)
			fmt.Println("Delta", eeg.Delta)
			fmt.Println("Theta", eeg.Theta)
			fmt.Println("LoAlpha", eeg.LoAlpha)
			fmt.Println("HiAlpha", eeg.HiAlpha)
			fmt.Println("LoBeta", eeg.LoBeta)
			fmt.Println("HiBeta", eeg.HiBeta)
			fmt.Println("LoGamma", eeg.LoGamma)
			fmt.Println("MidGamma", eeg.MidGamma)
			fmt.Println("\n")
			sendData(eeg, url)
		})
	}
	robot := gobot.NewRobot(
		"brainBot",
		[]gobot.Connection{adaptor},
		[]gobot.Device{neuro},
		work,
	)
	return robot
}

type AppArgs struct {
	Device string
	Url    string
}

func parseArgs() AppArgs {
	device := flag.String(
		"d",
		"/dev/tty.MindWaveMobile-DevA",
		"Address of the device from where the headset conncets",
	)
	url := flag.String(
		"p",
		"http://localhost:3000",
		"Url to send headset data",
	)
	flag.Parse()
	args := AppArgs{Device: *device, Url: *url}

	fmt.Println("Connecting to:", args.Device)
	fmt.Println("Sending data to:", args.Url)
	return args
}

func cleanUpFunction(gbot *gobot.Gobot) func() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	return func() {
		for sig := range c {
			fmt.Println("Shutting down from ", sig)
			gbot.Stop()
			os.Exit(1)
		}
	}
}

func main() {
	// parse the command line arguments
	args := parseArgs()

	// make the robot
	robot1 := makeRobot(args.Device, args.Url)

	// initialize gobot
	gbot := gobot.NewGobot()
	gbot.AddRobot(robot1)
	gbot.Start()

	// set the ctrl-c handler
	cleanup := cleanUpFunction(gbot)
	go cleanup()
}
