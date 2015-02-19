package main

import (
	"fmt"
	"log"
	"math"
	"time"

	"code.google.com/p/portaudio-go/portaudio"
)

const (
	inChannels  = 0
	outChannels = 1
	sampleRate  = 44100
	bufferSize  = 2048

	hz = 200

	int16Max = 1<<15 - 1
)

func Start() error {
	stream, err := portaudio.OpenDefaultStream(0, 1, sampleRate, bufferSize, processAudio)
	if err != nil {
		return err
	}

	defer stream.Close()
	if err = stream.Start(); err != nil {
		return err
	}

	defer stream.Stop()
	time.Sleep(5 * time.Second)
	return nil
}

func processAudio(out []int16) {
	sine := NewSine()
	sine.addWave()
	for i := range out {
		out[i] = int16(math.Min(1.0, math.Max(-1.0, sine.wave[i])) * int16Max)
	}
}

type Sine struct {
	phase float64
	wave  []float64
}

func NewSine() Sine {
	return Sine{
		phase: 0.0,
		wave:  make([]float64, bufferSize),
	}
}

func sine(x float64) float64 {
	return float64(math.Sin(2 * math.Pi * float64(x/4)))
}

func (s *Sine) addWave() {
	for i := range s.wave {
		s.wave[i] = s.nextFrameValue()
	}
}

func (s *Sine) nextFrameValue() float64 {
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

	s.phase += float64(hz) / float64(sampleRate)
	if s.phase > 1.0 {
		fmt.Println(int16Max)
		s.phase = s.phase - 1
	}

	return frame
}

func main() {
	err := portaudio.Initialize()
	if err != nil {
		log.Fatal(err)
	}
	defer portaudio.Terminate()

	err = Start()
	log.Println(err)

}
