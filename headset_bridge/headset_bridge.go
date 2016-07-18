package main

import (
    "fmt"
    "os"
    "os/signal"
    "encoding/json"
    "github.com/parnurzeal/gorequest"

    "github.com/hybridgroup/gobot"
    "github.com/hybridgroup/gobot/platforms/neurosky"
)

func eegToPayload(eeg neurosky.EEG) string {
    eegRawData := map[string]int {
        "delta": eeg.Delta,
        "theta": eeg.Theta,
        "loAlpha": eeg.LoAlpha,
        "hiAlpha": eeg.HiAlpha,
        "loBeta": eeg.LoBeta,
        "hiBeta": eeg.HiBeta,
        "loGamma": eeg.LoGamma,
        "midGamma": eeg.MidGamma,
    }
    eegData, _ := json.Marshal(eegRawData)
    payload := string(eegData)
    fmt.Println(payload)
    return payload
}

func sendData(eeg neurosky.EEG) {
    payload := eegToPayload(eeg)
    request := gorequest.New()
    response, body, errors := request.Post("http://localhost:3000").
        Set("Notes","gorequst is coming!").
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

func makeRobot(label string, device string) *gobot.Robot {
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
            sendData(eeg)

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

func main() {
    gbot := gobot.NewGobot()
    robot1 := makeRobot("A", "/dev/tty.MindWaveMobile-DevA")
    gbot.AddRobot(robot1)
    gbot.Start()

    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt)
    go func(){
        for sig := range c {
            fmt.Println("Shutting down from ", sig)
            gbot.Stop()
            os.Exit(1)
        }
    }()
}
