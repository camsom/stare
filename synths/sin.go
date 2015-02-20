package synths

import (
	"fmt"
	"log"
	"math"
)

const (
	sampleRate = 44100
	bufferSize = 2048

	hz = 220.0

	int16Max = 1<<15 - 1
)

type Sin struct {
	phase float64
	Wave  []float64
}

func NewSin() Sin {
	return Sin{
		phase: 0.0,
		Wave:  make([]float64, bufferSize),
	}
}

func sin(x float64) float64 {
	return float64(math.Sin(2 * math.Pi * float64(x/4)))
}

func (s *Sin) AddWave() {
	for i := range s.Wave {
		s.Wave[i] = s.nextFrameValue()
	}
}

func (s *Sin) nextFrameValue() float64 {
	var frame float64
	switch {
	case s.phase <= .25:
		frame = sin(s.phase * 4)
	case s.phase <= .50:
		frame = sin(1 - ((s.phase - .25) * 4))
	case s.phase <= .75:
		frame = -sin((s.phase - .50) * 4)
	case s.phase <= 1.00:
		frame = -sin((s.phase - .75) * 4)
	default:
		log.Fatal("Impossible phase")
	}

	s.phase += float64(hz) / float64(20000)
	if s.phase > 1.0 {
		fmt.Println(int16Max)
		s.phase = s.phase - 1
	}

	return frame
}
