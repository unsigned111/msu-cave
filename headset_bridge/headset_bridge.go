package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"os"
	"os/signal"
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/neurosky"
)

var DEFAULT_LOG_FILE_NAME string = ""

const (
	TIMESTAMP = iota
	DELTA
	HI_ALPHA
	HI_BETA
	LO_ALPHA
	LO_BETA
	LO_GAMMA
	MID_GAMMA
	THETA
	N_PARAMS
)

var nameMap = map[int]string{
	TIMESTAMP: "timestamp",
	DELTA:     "delta",
	HI_ALPHA:  "hiAlpha",
	HI_BETA:   "hiBeta",
	LO_ALPHA:  "loAlpha",
	LO_BETA:   "loBeta",
	LO_GAMMA:  "loGamma",
	MID_GAMMA: "midGamma",
	THETA:     "theta",
}

func eegRawData(eeg neurosky.EEG) map[string]int {
	return map[string]int{
		nameMap[TIMESTAMP]: int(time.Now().Unix()),
		nameMap[DELTA]:     eeg.Delta,
		nameMap[HI_ALPHA]:  eeg.HiAlpha,
		nameMap[HI_BETA]:   eeg.HiBeta,
		nameMap[LO_ALPHA]:  eeg.LoAlpha,
		nameMap[LO_BETA]:   eeg.LoBeta,
		nameMap[LO_GAMMA]:  eeg.LoGamma,
		nameMap[MID_GAMMA]: eeg.MidGamma,
		nameMap[THETA]:     eeg.Theta,
	}
}

func eegToPayload(eeg neurosky.EEG) string {
	eegRawData := eegRawData(eeg)
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

func format(i int, base string) string {
	var format_str = base
	if i < N_PARAMS-1 {
		format_str += ","
	}
	return format_str
}

func logHeader(logFile *os.File) {
	for i := 0; i < N_PARAMS; i++ {
		fmt.Fprintf(logFile, format(i, "%s"), nameMap[i])
	}
	fmt.Fprintf(logFile, "\n")
}

func logData(eeg neurosky.EEG, logFile *os.File) {
	payload := eegRawData(eeg)
	for i := 0; i < N_PARAMS; i++ {
		key := nameMap[i]
		fmt.Fprintf(logFile, format(i, "%d"), payload[key])
	}
	fmt.Fprintf(logFile, "\n")
}

func makeRobot(device string, url string, logFile *os.File) *gobot.Robot {
	adaptor := neurosky.NewNeuroskyAdaptor("neurosky", device)
	neuro := neurosky.NewNeuroskyDriver(adaptor, "neuro")
	work := func() {
		gobot.On(neuro.Event("eeg"), func(data interface{}) {
			eeg := data.(neurosky.EEG)
			fmt.Println("Delta", eeg.Delta)
			fmt.Println("HiAlpha", eeg.HiAlpha)
			fmt.Println("HiBeta", eeg.HiBeta)
			fmt.Println("LoAlpha", eeg.LoAlpha)
			fmt.Println("LoBeta", eeg.LoBeta)
			fmt.Println("LoGamma", eeg.LoGamma)
			fmt.Println("MidGamma", eeg.MidGamma)
			fmt.Println("Theta", eeg.Theta)
			fmt.Println("\n")
			sendData(eeg, url)
			logData(eeg, logFile)
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
