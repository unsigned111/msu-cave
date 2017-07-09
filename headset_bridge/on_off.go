package main

type OnOffModel struct {
	Threashold float64
	Samples    []uint8
	NextSlot   int
}

func MakeOnOffModel(threashold float64, numSamples int) OnOffModel {
	samples := make([]uint8, numSamples)
	for i, _ := range samples {
		samples[i] = 255
	}

	return OnOffModel{
		Threashold: threashold,
		Samples:    samples,
		NextSlot:   0,
	}
}

func (m OnOffModel) NumSamples() int {
	return len(m.Samples)
}

func (m *OnOffModel) AddSample(sample uint8) {
	nextSlot := m.NextSlot
	m.Samples[nextSlot] = sample
	m.NextSlot = (nextSlot + 1) % m.NumSamples()
}

func (m OnOffModel) isOn() bool {
	agg := 0.0
	for _, val := range m.Samples {
		agg += float64(val) / 255.0
	}
	avgSignal := agg / float64(m.NumSamples())
	return avgSignal < m.Threashold
}
