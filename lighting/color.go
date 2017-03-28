package main

import (
	"errors"
	"fmt"
)

// ColorMaxValue is the largest value a lighting channel can have
const ColorMaxValue int = 255

// LightingColor represents a given lighting state.
// They contain values for 4 channels:
//	red, green, blue, and an intensity channel.
type LightingColor struct {
	Red, Green, Blue, White, Intensity int
}

// Blackout creates a blacked out lighting state.
func Blackout() LightingColor {
	return LightingColor{0, 0, 0, 0, 0}
}

// Returns a string notation for a lighting state.
func (lc LightingColor) String() string {
	return fmt.Sprintf("[%d, %d, %d, %d] @ %d", lc.Red, lc.Green, lc.Blue, lc.White, lc.Intensity)
}

// Interpolate linearly interpolates between two lighting states.
func (lc LightingColor) Interpolate(end LightingColor, count int) (colors []LightingColor) {
	colors = make([]LightingColor, count)
	reds, _ := makeRange(lc.Red, end.Red, count)
	greens, _ := makeRange(lc.Green, end.Green, count)
	blues, _ := makeRange(lc.Blue, end.Blue, count)
	whites, _ := makeRange(lc.White, end.White, count)
	intensities, _ := makeRange(lc.Intensity, end.Intensity, count)
	for i := 0; i < count; i++ {
		colors[i].Red = reds[i]
		colors[i].Green = greens[i]
		colors[i].Blue = blues[i]
		colors[i].White = whites[i]
		colors[i].Intensity = intensities[i]
	}
	return
}

// Creates a range, given minimum and maximum values and the number of elements in the range.
func makeRange(min, max, count int) ([]int, error) {
	if count < 2 {
		return nil, errors.New("Cannot generate a range with fewer than 2 members")
	} else if count == 2 {
		return []int{min, max}, nil
	}
	a := make([]int, count)
	a[0] = min
	a[len(a)-1] = max
	var step float32
	if max < min {
		step = float32(max-min-1) / float32(count)
	} else {
		step = float32(max-min) / float32(count)
	}

	current := float32(min) + step
	for i := 1; i < len(a)-1; i++ {
		a[i] = round(current)
		current = current + step
	}
	return a, nil
}

//  Returns the minimum of the input values.
func min(x int, y int) int {
	if x > y {
		return y
	}
	return x
}

// round rounds float32s to nearest integer
func round(val float32) int {
	if val < 0 {
		return int(val - 0.5)
	}
	return int(val + 0.5)
}
