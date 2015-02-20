package main

import (
	"log"
	"math"
	"time"

	"stare/synths"

	"code.google.com/p/portaudio-go/portaudio"
)

const (
	inChannels  = 0
	outChannels = 1
	sampleRate  = 44100
	bufferSize  = 2048

	hz = 220.0

	int16Max = 1<<16 - 1
)

type AudioEngine struct {
	out []float64
	syn []synths.Synth
}

func NewAudioEngine() AudioEngine {
	return AudioEngine{
		out: make([]float64, bufferSize),
	}
}

func (a *AudioEngine) AddSynth(name string, frequency float64) {
	f := frequency * math.Exp2((float64(12*0))/12)

	log.Printf("Adding %s synth with %fhz.", name, f)
	s := synths.NewSine(sampleRate, f)
	s.AddWave()

	a.syn = append(a.syn, s)
}

func (a AudioEngine) Start() error {
	stream, err := portaudio.OpenDefaultStream(0, 1, sampleRate, bufferSize, a.processAudio)
	if err != nil {
		return err
	}

	defer stream.Close()
	if err = stream.Start(); err != nil {
		return err
	}

	defer stream.Stop()
	time.Sleep(100 * time.Second)
	return nil
}

func (a *AudioEngine) processAudio(out []int16) {

	a.out = make([]float64, bufferSize)

	for _, s := range a.syn {
		s.Add(a.out)
	}

	for i := range out {
		// out[i] = int16(math.Min(1.0, math.Max(-1.0, a.out[i])) * int16Max)
		out[i] = int16(int16Max * a.out[i])
	}
}

func main() {
	err := portaudio.Initialize()
	if err != nil {
		log.Fatal(err)
	}
	defer portaudio.Terminate()

	a := NewAudioEngine()
	// a.AddSynth("sine", 110.0)
	// a.AddSynth("sine", 580.0)
	// a.AddSynth("sine", 670.0)
	a.AddSynth("sine", 300.0)

	if err = a.Start(); err != nil {
		log.Println(err)
	}

}
