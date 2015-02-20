package synths

import (
	"log"
	"math"
)

type Synth struct {
	SampleRate float64
	Frequency  float64
	phase      float64
	Wave       []float64
}

func NewSine(sr, hz float64) Synth {
	s := Synth{
		SampleRate: sr,
		Frequency:  hz,
		phase:      1.0,
		Wave:       make([]float64, bufferSize),
	}

	s.AddWave()

	return s
}

func NewSquare(sr, hz float64) {

}

func (s Synth) Add(out []float64) {
	add(s.Wave, out)
}

func (s Synth) AddWave() {
	for i := range s.Wave {
		s.Wave[i] = s.nextFrameValue()
	}
}

func (s *Synth) nextFrameValue() float64 {
	var frame float64
	switch {
	case s.phase <= .25:
		frame = sine(s.phase * 4)
	case s.phase <= .50:
		frame = sine(1 - ((s.phase - .25) * 4))
	case s.phase <= .75:
		frame = -sine((s.phase - .50) * 4)
	case s.phase <= 1.00:
		frame = -sine((s.phase - .75) * 4)
	default:
		log.Fatal("Impossible phase")
	}

	s.phase += float64(s.Frequency) / float64(s.SampleRate)
	if s.phase > 1.0 {
		s.phase = s.phase - 1
	}

	return frame
}

func sine(x float64) float64 {
	return float64(math.Sin(2 * math.Pi * float64(x/4)))
}

func add(in []float64, out []float64) {
	for i := range out {
		out[i] += in[i]
	}
}
