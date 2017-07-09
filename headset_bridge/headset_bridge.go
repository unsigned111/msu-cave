package main

import (
	"flag"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"os"
	"os/signal"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/neurosky"
)

var DEFAULT_LOG_FILE_NAME string = ""

const ON_OFF_THREASHOLD = .5
const ON_OFF_WINDOW_SIZE = 3

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
		fmt.Println("Request sent successfully")
	}
}

func logHeader(logFile *os.File) {
	LogHeader(logFile)
}

func logData(state State, logFile *os.File) {
	state.LogData(logFile)
}

func makeRobot(device string, url string, logFile *os.File) *gobot.Robot {
	adaptor := neurosky.NewNeuroskyAdaptor("neurosky", device)
	neuro := neurosky.NewNeuroskyDriver(adaptor, "neuro")

	eegChan := make(chan neurosky.EEG)
	headsetOnChan := make(chan bool)

	state := State{}
	aggregateor := func() {
		for {
			select {
			case eeg := <-eegChan:
				state.UpdateEEG(eeg)
			case headsetOn := <-headsetOnChan:
				state.TestUpdateHeadsetOn(headsetOn)
			}
			sendData(state, url)
			logData(state, logFile)
		}
	}
	go aggregateor()

	work := func() {
		gobot.On(neuro.Event("eeg"), func(data interface{}) {
			eeg := data.(neurosky.EEG)
			eegChan <- eeg
		})

		onOff := MakeOnOffModel(ON_OFF_THREASHOLD, ON_OFF_WINDOW_SIZE)
		gobot.On(neuro.Event("signal"), func(data interface{}) {
			sample := data.(uint8)
			onOff.AddSample(sample)
			isOn := onOff.isOn()
			headsetOnChan <- isOn
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

	// make the robot
	robot1 := makeRobot(args.Device, args.Url, logFile)

	// initialize gobot
	gbot := gobot.NewGobot()
	gbot.AddRobot(robot1)
	gbot.Start()

	// set the ctrl-c handler
	cleanup := cleanUpFunction(gbot)
	go cleanup()
}
