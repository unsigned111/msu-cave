package main

import (
	"bytes"
	"github.com/hybridgroup/gobot/platforms/neurosky"
	"github.com/stretchr/testify/assert"
	"testing"
)

func stubNow() {
	nowFunc = func() int { return 123456 }
}

func unstubNow() {
	resetClockImplementation()
}

func makeEEG() neurosky.EEG {
	return neurosky.EEG{
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
	state.UpdateHeadsetOn(true)
	state.UpdateEEG(makeEEG())
	return state
}

func assertEqualEEG(t *testing.T, eeg neurosky.EEG, state State) {
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

func TestUpdateHeadsetOn(t *testing.T) {
	state := State{}
	state.UpdateHeadsetOn(true)
	assert.Equal(t, true, state.HeadsetOn)
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
	expectedPayload += `"midGamma":7,"theta":8,"headsetOn":true}`

	assert.Equal(t, expectedPayload, payload)

	unstubNow()
}

func TestLogHeader(t *testing.T) {
	buffer := bytes.NewBufferString("")
	LogHeader(buffer)

	expectedHeader := "Timestamp,Delta,HiAlpha,HiBeta,LoAlpha,"
	expectedHeader += "LoBeta,LoGamma,MidGamma,Theta,HeadsetOn\n"

	assert.Equal(t, expectedHeader, buffer.String())
}

func TestLogData(t *testing.T) {
	stubNow()

	state := makeState()
	buffer := bytes.NewBufferString("")
	state.LogData(buffer)

	expectedData := "123456,1,2,3,4,5,6,7,8,true\n"
	assert.Equal(t, expectedData, buffer.String())

	unstubNow()
}
