package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMakeOnOffModel(t *testing.T) {
	m := MakeOnOffModel(.5, 3)
	assert.Equal(t, .5, m.Threashold, "sets threashold")
	assert.Equal(t, []uint8{255, 255, 255}, m.Samples, "sets sample buffer")
	assert.Equal(t, 0, m.NextSlot, "initializes the tail of the buffer")
}

func TestAddSample(t *testing.T) {
	// given
	m := MakeOnOffModel(.5, 3)

	// add first
	m.AddSample(8)
	assert.Equal(t, 1, m.NextSlot, "incraments next slot")
	assert.Equal(t, []uint8{8, 255, 255}, m.Samples, "updates 0th elemnt")

	// add second
	m.AddSample(7)
	assert.Equal(t, 2, m.NextSlot, "incraments next slot")
	assert.Equal(t, []uint8{8, 7, 255}, m.Samples, "updates 1st element")

	// add third
	m.AddSample(6)
	assert.Equal(t, 0, m.NextSlot, "wraps next slot")
	assert.Equal(t, []uint8{8, 7, 6}, m.Samples, "updates 2nd elt")

	// add fouth and wrap
	m.AddSample(5)
	assert.Equal(t, 1, m.NextSlot, "incraments the next slot")
	assert.Equal(t, []uint8{5, 7, 6}, m.Samples, "updates 0th elt")
}

func TestNumSamples(t *testing.T) {
	m := MakeOnOffModel(.5, 4)
	assert.Equal(t, 4, m.NumSamples())
}

func TestIsOn(t *testing.T) {
	m := MakeOnOffModel(.5, 2)

	m.Samples = []uint8{4, 10}
	assert.Equal(t, true, m.isOn(), "it is on for low values")

	m.Samples = []uint8{225, 250}
	assert.Equal(t, false, m.isOn(), "it is off for high values")
}
