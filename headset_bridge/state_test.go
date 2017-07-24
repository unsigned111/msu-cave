package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"gobot.io/x/gobot/platforms/neurosky"
	"testing"
)

func stubNow() {
	nowFunc = func() int { return 123456 }
}

func unstubNow() {
	resetClockImplementation()
}

func makeEEG() neurosky.EEGData {
	return neurosky.EEGData{
		Delta:    1,
		HiAlpha:  2,
		HiBeta:   3,
		LoAlpha:  4,
		LoBeta:   5,
		LoGamma:  6,
		MidGamma: 7,
		Theta:    8,
	}
}

func makeState() State {
	state := State{}
	state.UpdateEEG(makeEEG())
	state.UpdateSignal(9)
	state.UpdateAttention(10)
	state.UpdateMeditation(11)
	return state
}

func assertEqualEEG(t *testing.T, eeg neurosky.EEGData, state State) {
	assert.Equal(t, eeg.Delta, state.Delta)
	assert.Equal(t, eeg.HiAlpha, state.HiAlpha)
	assert.Equal(t, eeg.HiBeta, state.HiBeta)
	assert.Equal(t, eeg.LoAlpha, state.LoAlpha)
	assert.Equal(t, eeg.LoBeta, state.LoBeta)
	assert.Equal(t, eeg.LoGamma, state.LoGamma)
	assert.Equal(t, eeg.MidGamma, state.MidGamma)
	assert.Equal(t, eeg.Theta, state.Theta)
}

func TestUpdateEEG(t *testing.T) {
	eeg := makeEEG()

	state := State{}
	state.UpdateEEG(eeg)
	assertEqualEEG(t, eeg, state)
}

func TestUpdateSignal(t *testing.T) {
	state := State{}
	state.UpdateSignal(8)
	assert.Equal(t, 8, state.Signal)
}

func TestUpdateAttention(t *testing.T) {
	state := State{}
	state.UpdateAttention(7)
	assert.Equal(t, 7, state.Attention)
}

func TestUpdateMeditation(t *testing.T) {
	state := State{}
	state.UpdateMeditation(9)
	assert.Equal(t, 9, state.Meditation)
}

func TestJSONSerialization(t *testing.T) {
	stubNow()

	state := makeState()
	payload := state.AsPayload()

	expectedPayload := `{"timestamp":123456,"delta":1,"hiAlpha":2,`
	expectedPayload += `"hiBeta":3,"loAlpha":4,"loBeta":5,"loGamma":6,`
	expectedPayload += `"midGamma":7,"theta":8,"signal":9,`
	expectedPayload += `"attention":10,"meditation":11}`

	assert.Equal(t, expectedPayload, payload)

	unstubNow()
}

func TestLogHeader(t *testing.T) {
	buffer := bytes.NewBufferString("")
	LogHeader(buffer)

	expectedHeader := "Timestamp,Delta,HiAlpha,HiBeta,LoAlpha,"
	expectedHeader += "LoBeta,LoGamma,MidGamma,Theta,Signal,"
	expectedHeader += "Attention,Meditation\n"

	assert.Equal(t, expectedHeader, buffer.String())
}

func TestLogData(t *testing.T) {
	stubNow()

	state := makeState()
	buffer := bytes.NewBufferString("")
	state.LogData(buffer)

	expectedData := "123456,1,2,3,4,5,6,7,8,9,10,11\n"
	assert.Equal(t, expectedData, buffer.String())

	unstubNow()
}
