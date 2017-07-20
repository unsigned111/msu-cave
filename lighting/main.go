package main

import (
	"flag"
	"github.com/dordille/artnet"
	"github.com/hypebeast/go-osc/osc"
	"log"
	"strconv"
	"sync"
	"time"
)

var lastColor LightingColor
var settings Settings
var headsetOn bool
var headsetLock sync.Mutex
var similarity float64
var similarityLock sync.Mutex
var colorChannel = make(chan LightingColor, 512)
var debug bool

// Runs the OSC to ArtNet packet service.
func main() {
	handleRuntimeFlags()
	initLighting()
	go artnetSend()
	go idleLighting()
	oscListen()
}

// Handles setting up runtime flags and processing their values.
func handleRuntimeFlags() {
	var configFlag = flag.String("f", "./settings.json", "The configuration file to use")
	var debugFlag = flag.Bool("d", false, "Enable debug output")
	flag.Parse()
	settings = ParseConfig(*configFlag)
	debug = *debugFlag
	headsetLock.Lock()
	headsetOn = true
	headsetLock.Unlock()
}

// Starts an OSC listener (server) based on the settings struct.
func oscListen() {
	addr := settings.ListenIP + ":" + strconv.Itoa(settings.ListenPort)
	server := &osc.Server{Addr: addr}

	log.Printf("OSC server starting at : " + addr)
	log.Printf("OSC listening on : %s (%s)", "/eeg", "EEG channel values")
	log.Printf("OSC listening on : %s (%s)", "/onoff", "Headset on/off")
	log.Printf("OSC listening on : %s (%s)", "/similarity", "Similarity")
	server.Handle("/eeg", func(msg *osc.Message) {
		value := msg.Arguments[settings.HeadsetChannel]
		activeLighting(int(value.(int32)))
	})
	server.Handle("/onoff", func(msg *osc.Message) {
		value := msg.Arguments[0]
		handleOnOff(int(value.(int32)))

	})
	server.Handle("/similarity", func(msg *osc.Message) {
		value := msg.Arguments[0]
		handleSimilarity(float64(value.(int32)))
	})
	server.ListenAndServe()
}

// Starts the ArtNet sender based on the settings struct.
// Binds to the UDP port, then loops infinitely, pulling
// lighting states out of the lighting queue and transmits
// those instructions as ArtNet frames.
func artnetSend() {
	addr := settings.ControlIP + ":" + strconv.Itoa(settings.ControlPort)
	log.Print("Artnet sender connecting to : " + addr)
	err, artnetInterface := artnet.NewNode(addr)
	check(err)

	// calculate how long it takes to send a frame
	frameTimespan := time.Duration(1000/settings.FPS) * time.Millisecond
	var data [512]uint8
	artnetInterface.Dmx(1, data)
	for {
		// start the execution duration timer
		start := time.Now()
		var data [512]uint8
		currentColor := <-colorChannel
		// stuff color values into the dmx frame
		data[0] = uint8(currentColor.Red)
		data[1] = uint8(currentColor.Green)
		data[2] = uint8(currentColor.Blue)
		data[3] = uint8(currentColor.White)
		data[5] = uint8(currentColor.Intensity)
		artnetInterface.Dmx(0, data)
		if debug {
			log.Print("SENDING ARTNET FRAME : " + currentColor.String())
		}

		// calculate execution duration
		elapsed := time.Since(start)

		// if the duration is less than the time allotted, sleep for the rest of the allowed time
		if elapsed < frameTimespan {
			duration := frameTimespan - elapsed
			time.Sleep(duration)
		}
	}
}

func handleOnOff(value int) {
	headsetLock.Lock()
	if value == 0 {
		headsetOn = false
	} else {
		headsetOn = true
	}
	if debug {
		log.Printf("Turning headset listening %t\n", headsetOn)
	}
	headsetLock.Unlock()
}

func handleSimilarity(value float64) {
	similarityLock.Lock()
	similarity = value
	if debug {
		log.Printf("RECEIVED SIMILARITY : %f\n", similarity)
	}
	similarityLock.Unlock()
}

// Given a specific channel value, generates a collection of lighting
// states linearly interpolated between the last broadcast lighting state to
// the next desired lighting state.
func activeLighting(value int) {
	if readHeadsetState() == true {
		if debug {
			log.Printf("RECIEVED VALUE (channel %d): %d\n", settings.HeadsetChannel, value)
		}

		flushQueue()

		var nextColor LightingColor
		var target float32
		var intensity = int(float32(float64(settings.MaxIntensity - settings.MinIntensity) * similarity)) + settings.MinIntensity

		// if we have received a 0 value, use the last frame instead
		if value == 0 {
			nextColor = lastColor
		} else {
			var colors []LightingColor
			var redSpec = float32(settings.EndColor.Red - settings.StartColor.Red)
			var greenSpec = float32(settings.EndColor.Green - settings.StartColor.Green)
			var blueSpec = float32(settings.EndColor.Blue - settings.StartColor.Blue)
			var whiteSpec = float32(settings.EndColor.White - settings.StartColor.White)

			target = float32(value) / float32(settings.Scaler)
			nextColor = LightingColor{
				int(target*redSpec) + settings.StartColor.Red,
				int(target*greenSpec) + settings.StartColor.Green,
				int(target*blueSpec) + settings.StartColor.Blue,
				int(target*whiteSpec) + settings.StartColor.White,
				intensity,
			}
			colors = lastColor.Interpolate(nextColor, settings.FPS-1)
			queueColors(colors)
			queueColor(nextColor)
		}
		if debug {
			log.Printf("Target: %f", target)
			log.Printf("%s", nextColor)
		}
		lastColor = nextColor
	}
}

// Generates pulses in intensity of the "headset off" lighting color
// including a delay in between pulses. Continually generates these pulses
// until the headset flag is flipped back on.
func idleLighting() {
	halfPulseLength := settings.FPS * (settings.PulseLength >> 1)
	pauseDuration := time.Duration(settings.PulsePause) * time.Second
	for {
		if readHeadsetState() == false {
			// calculate the number of messages to send using the FPS and the Pulse Length from settings
			rampUp := settings.OffStartColor.Interpolate(settings.OffEndColor, halfPulseLength)
			queueColors(rampUp)
			rampDown := settings.OffEndColor.Interpolate(settings.OffStartColor, halfPulseLength)
			queueColors(rampDown)
			lastColor = settings.OffStartColor
			time.Sleep(pauseDuration)
		}
	}
}

// Initializes the lighting system.
func initLighting() {
	blackout := Blackout()
	queueColor(blackout)
	lastColor = blackout
}

// Enqueues a collection of LightingColors into the lighting queue.
func queueColors(colors []LightingColor) {
	for _, value := range colors {
		colorChannel <- value
	}
}

// Enqueues a single LightingColor into the lighting queue.
func queueColor(color LightingColor) {
	colorChannel <- color
}

// Flushes the entirety of the lighting queue.
func flushQueue() {
	for {
		select {
		case _ = <-colorChannel:
		default:
			return
		}
	}
}

// readHeadsetState locks a mutex, retrieves the value, and returns it
func readHeadsetState() bool {
	headsetLock.Lock()
	val := headsetOn
	headsetLock.Unlock()
	return val
}

// Checks for an error and bails if an error is found.
// In general we don't want to fail gracefully in this app.
func check(err error) {
	if err != nil {
		log.Fatal()
	}
}
