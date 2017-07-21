package main

import (
	"encoding/json"
	"fmt"
	"gobot.io/x/gobot/platforms/neurosky"
	"io"
	"reflect"
	"time"
)

var nowFunc func() int

func init() {
	resetClockImplementation()
}

func resetClockImplementation() {
	nowFunc = func() int {
		return int(time.Now().Unix())
	}
}

func now() int {
	return nowFunc()
}

type State struct {
	Timestamp  int  `json:"timestamp"`
	Delta      int  `json:"delta"`
	HiAlpha    int  `json:"hiAlpha"`
	HiBeta     int  `json:"hiBeta"`
	LoAlpha    int  `json:"loAlpha"`
	LoBeta     int  `json:"loBeta"`
	LoGamma    int  `json:"loGamma"`
	MidGamma   int  `json:"midGamma"`
	Theta      int  `json:"theta"`
	HeadsetOn  bool `json:"headsetOn"`
	Attention  int  `json:"attention"`
	Meditation int  `json:"meditation"`
}

func (s *State) UpdateTimestamp() {
	s.Timestamp = now()
}

func (s *State) UpdateEEG(eeg neurosky.EEGData) {
	s.UpdateTimestamp()
	s.Delta = eeg.Delta
	s.HiAlpha = eeg.HiAlpha
	s.HiBeta = eeg.HiBeta
	s.LoAlpha = eeg.LoAlpha
	s.LoBeta = eeg.LoBeta
	s.LoGamma = eeg.LoGamma
	s.MidGamma = eeg.MidGamma
	s.Theta = eeg.Theta
}

func (s *State) UpdateHeadsetOn(headsetOn bool) {
	s.UpdateTimestamp()
	s.HeadsetOn = headsetOn
}

func (s *State) UpdateAttention(attention int) {
	s.Timestamp = now()
	s.Attention = attention
}

func (s *State) UpdateMeditation(meditation int) {
	s.Timestamp = now()
	s.Meditation = meditation
}

func (s *State) AsPayload() string {
	data, _ := json.Marshal(s)
	return string(data)
}

func logDelimiter(writer io.Writer, numFields int, fieldNum int) {
	if fieldNum < numFields-1 {
		fmt.Fprintf(writer, ",")
	}
}

func LogHeader(writer io.Writer) {
	s := State{}
	v := reflect.ValueOf(s)
	numFields := v.NumField()

	for i := 0; i < numFields; i++ {
		fieldName := v.Type().Field(i).Name
		fmt.Fprint(writer, fieldName)
		logDelimiter(writer, numFields, i)
	}
	fmt.Fprintln(writer, "")
}

func (s *State) LogData(writer io.Writer) {
	v := reflect.ValueOf(*s)
	numFields := v.NumField()

	for i := 0; i < numFields; i++ {
		fieldVal := v.Field(i).Interface()
		fmt.Fprint(writer, fieldVal)
		logDelimiter(writer, numFields, i)
	}
	fmt.Fprintln(writer, "")
}
