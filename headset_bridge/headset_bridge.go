package main

import (
	"flag"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"os"
	"os/signal"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/neurosky"
)

var DEFAULT_LOG_FILE_NAME string = ""

func sendData(state State, url string) {
	payload := state.AsPayload()
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
		fmt.Println("Request sent successfully", state.Timestamp)
	}
}

func logHeader(logFile *os.File) {
	LogHeader(logFile)
}

func logData(state State, logFile *os.File) {
	state.LogData(logFile)
}

func aggregateor(
	hub Hub,
	url string,
	logFile *os.File,
) {
	state := State{}
	lastSend := state.Timestamp
	for {
		select {
		case eeg := <-hub.EEG:
			state.UpdateEEG(eeg)
		case signal := <-hub.Signal:
			state.UpdateSignal(signal)
		case attention := <-hub.Attention:
			state.UpdateAttention(attention)
		case meditation := <-hub.Meditation:
			state.UpdateMeditation(meditation)
		}

		if lastSend < state.Timestamp {
			sendData(state, url)
			logData(state, logFile)
			lastSend = state.Timestamp
		}
	}
}

func makeRobot(
	device string,
	hub Hub,
) *gobot.Robot {
	adaptor := neurosky.NewAdaptor(device)
	neuro := neurosky.NewDriver(adaptor)

	work := func() {

		neuro.On(neuro.Event("eeg"), func(data interface{}) {
			eeg := data.(neurosky.EEGData)
			hub.EEG <- eeg
		})

		neuro.On(neuro.Event("signal"), func(data interface{}) {
			signal := data.(uint8)
			hub.Signal <- int(signal)
		})

		neuro.On(neuro.Event("attention"), func(data interface{}) {
			attention := data.(uint8)
			hub.Attention <- int(attention)
		})

		neuro.On(neuro.Event("meditation"), func(data interface{}) {
			meditation := data.(uint8)
			hub.Meditation <- int(meditation)
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
	Device      string
	Url         string
	LogFileName string
}

type Hub struct {
	EEG        chan neurosky.EEGData
	Signal     chan int
	Attention  chan int
	Meditation chan int
}

func parseArgs() AppArgs {
	device := flag.String(
		"d",
		"/dev/tty.MindWaveMobile-DevA",
		"Address of the device from where the headset connects",
	)
	url := flag.String(
		"p",
		"http://localhost:3000",
		"Url to send headset data",
	)
	logFileName := flag.String(
		"l",
		DEFAULT_LOG_FILE_NAME,
		"Log file name, if not provided, no logging will occur",
	)
	flag.Parse()

	args := AppArgs{
		Device:      *device,
		Url:         *url,
		LogFileName: *logFileName,
	}

	fmt.Println("Connecting to:", args.Device)
	fmt.Println("Sending data to:", args.Url)
	fmt.Println("Logging to:", args.LogFileName)
	return args
}

func cleanUpFunction(gbot *gobot.Robot) func() {
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

func openLogFile(logFileName string) (*os.File, bool) {
	opened := false
	var logFile *os.File
	var err error
	if logFileName != DEFAULT_LOG_FILE_NAME {
		logFile, err = os.Create(logFileName)
		opened = err == nil
		if !opened {
			fmt.Println("There was an error opening the log file, logging disabled")
			fmt.Println(err)
		}
	}
	return logFile, opened
}

func main() {
	// parse the command line arguments
	args := parseArgs()

	// init logfile (if requested)
	logFile, logging := openLogFile(args.LogFileName)
	if logging {
		logHeader(logFile)
		defer logFile.Close()
	}

	// setup the channels
	hub := Hub{
		EEG:        make(chan neurosky.EEGData),
		Signal:     make(chan int),
		Attention:  make(chan int),
		Meditation: make(chan int),
	}

	// init aggregator
	go aggregateor(hub, args.Url, logFile)

	// make the robot
	robot := makeRobot(args.Device, hub)
	robot.Start()

	// set the ctrl-c handler
	cleanup := cleanUpFunction(robot)
	go cleanup()
}
