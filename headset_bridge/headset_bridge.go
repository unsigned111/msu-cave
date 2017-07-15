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

func aggregateor(
	hub Hub,
	url string,
	logFile *os.File,
) {
	state := State{}
	for {
		select {
		case eeg := <-hub.EEG:
			state.UpdateEEG(eeg)
		case headsetOn := <-hub.HeadsetOn:
			state.UpdateHeadsetOn(headsetOn)
		case attention := <-hub.Attention:
			state.UpdateAttention(attention)
		case meditation := <-hub.Meditation:
			state.UpdateMeditation(meditation)
		}
		sendData(state, url)
		logData(state, logFile)
	}
}

func makeRobot(
	device string,
	hub Hub,
) *gobot.Robot {
	adaptor := neurosky.NewNeuroskyAdaptor("neurosky", device)
	neuro := neurosky.NewNeuroskyDriver(adaptor, "neuro")

	work := func() {

		gobot.On(neuro.Event("eeg"), func(data interface{}) {
			eeg := data.(neurosky.EEG)
			hub.EEG <- eeg
		})

		onOff := MakeOnOffModel(ON_OFF_THREASHOLD, ON_OFF_WINDOW_SIZE)
		gobot.On(neuro.Event("signal"), func(data interface{}) {
			sample := data.(uint8)
			onOff.AddSample(sample)
			isOn := onOff.isOn()
			hub.HeadsetOn <- isOn
		})

		// TODO:DLM: figure out why meditation and attention
		// are not sending anything but 0
		// gobot.On(neuro.Event("attention"), func(data interface{}) {
		// 	attention := data.(uint8)
		// 	hub.Attention <- int(attention)
		// })
		//
		// gobot.On(neuro.Event("meditation"), func(data interface{}) {
		// 	meditation := data.(uint8)
		// 	hub.Meditation <- int(meditation)
		// })
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
	EEG        chan neurosky.EEG
	HeadsetOn  chan bool
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

	// setup the channels
	hub := Hub{
		EEG:        make(chan neurosky.EEG),
		HeadsetOn:  make(chan bool),
		Attention:  make(chan int),
		Meditation: make(chan int),
	}

	// init aggregator
	go aggregateor(hub, args.Url, logFile)

	// make the robot
	robot1 := makeRobot(args.Device, hub)

	// initialize gobot
	gbot := gobot.NewGobot()
	gbot.AddRobot(robot1)
	gbot.Start()

	// set the ctrl-c handler
	cleanup := cleanUpFunction(gbot)
	go cleanup()
}
