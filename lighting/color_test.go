package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBlackOut(t *testing.T) {
	x := Blackout()
	assert.ObjectsAreEqual(x, LightingColor{0, 0, 0, 0, 0})
}

func TestMin(t *testing.T) {
	// test left-handed value larger
	x := min(1, 0)
	assert.True(t, x == 0)

	// test right-handed value larger
	x = min(0, 1)
	assert.True(t, x == 0)

	// test both values equal
	x = min(1, 1)
	assert.True(t, x == 1)
}

func TestRound(t *testing.T) {
	// round a value already equal to an integer value
	x := round(1.0)
	assert.Equal(t, 1, x)

	// round down (positive)
	x = round(1.1)
	assert.Equal(t, 1, x)

	// round up (positive)
	x = round(0.9)
	assert.Equal(t, 1, x)

	// round up (positive) at value ending with .5
	x = round(1.5)
	assert.Equal(t, 2, x)

	// round down (positive) at value ending with .499999
	x = round(1.49999)
	assert.Equal(t, 1, x)

	// round a value already equal to an integer value
	x = round(-1.0)
	assert.Equal(t, -1, x)

	// round down (negative)
	x = round(-1.6)
	assert.Equal(t, -2, x)

	// round up (positive)
	x = round(-1.4)
	assert.Equal(t, -1, x)

	// round up (positive) at value ending with .5
	x = round(-1.5)
	assert.Equal(t, -2, x)

	// round down (positive) at value ending with .499999
	x = round(-1.49999)
	assert.Equal(t, -1, x)
}

func TestMakeRange(t *testing.T) {
	// check for when the desired count is smaller than 2
	x, err := makeRange(0, 1, 1)
	assert.Nil(t, x)
	assert.NotNil(t, err)

	// check for when the desired count is exactly 2
	x, err = makeRange(0, 1, 2)
	assert.Equal(t, []int{0, 1}, x)
	assert.Nil(t, err)

	// check for when the range width is smaller than the number of steps (each step is < 1)
	x, err = makeRange(0, 1, 4)
	assert.Equal(t, []int{0, 0, 1, 1}, x)
	assert.Nil(t, err)

	// check for when the range width is the same as the number of steps (each step is 1)
	x, err = makeRange(0, 3, 4)
	assert.Equal(t, []int{0, 1, 2, 3}, x)
	assert.Nil(t, err)

	// check for when the range width is larger than the number of steps (each step is > 1)
	x, err = makeRange(0, 7, 4)
	assert.Equal(t, []int{0, 2, 4, 7}, x)
	assert.Nil(t, err)

	// check for when range is negative and each step is < 1
	x, err = makeRange(1, 0, 4)
	assert.Equal(t, []int{1, 1, 0, 0}, x)
	assert.Nil(t, err)

	// check for when range is negative and each step is == 1
	x, err = makeRange(3, 0, 4)
	assert.Equal(t, []int{3, 2, 1, 0}, x)
	assert.Nil(t, err)

	// check for when range is negative and each step is > 1
	x, err = makeRange(7, 0, 4)
	assert.Equal(t, []int{7, 5, 3, 0}, x)
	assert.Nil(t, err)
}

func TestLightingColor_Interpolate(t *testing.T) {
	// check for when the starting color is smaller than the ending color
	x := LightingColor{0, 0, 0, 0, 0}
	y := LightingColor{255, 255, 255, 255, 255}
	z := x.Interpolate(y, 256)
	for i := 1; i < 256; i++ {
		assert.True(t, z[i].Red >= z[i-1].Red)
		assert.True(t, z[i].Green >= z[i-1].Green)
		assert.True(t, z[i].Blue >= z[i-1].Blue)
		assert.True(t, z[i].White >= z[i-1].White)
		assert.True(t, z[i].Intensity >= z[i-1].Intensity)
	}
	assert.Equal(t, 256, len(z))

	// check for when the starting color is greater than the ending color
	x = LightingColor{255, 255, 255, 255, 255}
	y = LightingColor{0, 0, 0, 0, 0}
	z = x.Interpolate(y, 256)
	for i := 1; i < 256; i++ {
		assert.True(t, z[i].Red <= z[i-1].Red)
		assert.True(t, z[i].Green <= z[i-1].Green)
		assert.True(t, z[i].Blue <= z[i-1].Blue)
		assert.True(t, z[i].White <= z[i-1].White)
		assert.True(t, z[i].Intensity <= z[i-1].Intensity)
	}
	assert.Equal(t, 256, len(z))

	// check for when the starting color and the ending color are the same
	x = LightingColor{255, 255, 255, 255, 255}
	y = LightingColor{255, 255, 255, 255, 255}
	z = x.Interpolate(y, 16)
	for i := 0; i < 16; i++ {
		assert.Equal(t, x, z[i])
		assert.Equal(t, y, z[i])
	}
}

func TestLightingColor_String(t *testing.T) {
	assert.Equal(t, "[0, 0, 0, 0] @ 0", LightingColor{0, 0, 0, 0, 0}.String())
	assert.Equal(t, "[255, 255, 255, 255] @ 255", LightingColor{255, 255, 255, 255, 255}.String())
}
