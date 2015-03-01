package main

import (
	"log"
	"math"
	"time"

	"github.com/rakyll/portmidi"

	"stare/synths"

	"code.google.com/p/portaudio-go/portaudio"
)

const (
	inChannels  = 0
	outChannels = 1
	sampleRate  = 44100
	bufferSize  = 16384

	hz = 300.0

	int16Max = 4000
)

var (
	Midi         *portmidi.Stream
	MidiSep      portmidi.Timestamp
	MidiVelocity int64

	Sin    []float64
	Square []float64
)

type AudioEngine struct {
	stream *portaudio.Stream

	out []float64
	syn []synths.Synth
}

func NewAudioEngine() AudioEngine {
	return AudioEngine{
		out: make([]float64, sampleRate),
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

	a.stream = stream
	a.add(Sin)
	a.add(Square)

	if err := a.stream.Start(); err != nil {
		return err
	}

	dur := math.Min(120, math.Max(float64(MidiVelocity), 150))
	time.Sleep(time.Duration(dur) * time.Millisecond)

	a.stream.Close()

	return nil
}

func (a AudioEngine) Stop() error {
	err := a.stream.Stop()
	if err != nil {
		return err
	}

	return nil
}

func (a *AudioEngine) processAudio(out []int16) {
	sample := make([]float64, sampleRate)

	for i := range a.out {
		sample[i] = a.out[i]
	}

	for i := range out {
		out[i] = int16(math.Min(1.0, sample[i]) * int16Max)
	}

}

func staticSine(hz float64) []float64 {
	phaseInc := hz * (2 * math.Pi) / sampleRate
	phase := 0.0

	sampleLength := sampleRate
	sample := make([]float64, sampleLength)
	for i := 0; i < sampleLength; i++ {
		sample[i] = math.Sin(phase)
		phase += phaseInc
		if phase >= (2 * math.Pi) {
			phase -= (2 * math.Pi)
		}
	}

	return sample
}

func staticSquare(hz float64) []float64 {
	phaseInc := hz * (2 * math.Pi) / sampleRate
	twoDPi := 2.0 / math.Pi
	phase := 0.0

	sample := make([]float64, sampleRate)
	for i := 0; i < sampleRate; i++ {
		triValue := phase * twoDPi
		if triValue < 0 {
			triValue += 1.0
		} else {
			triValue = 1.0 - triValue
		}
		sample[i] = triValue
		phase += phaseInc
		if phase >= math.Pi {
			phase -= (2 * math.Pi)
		}
	}

	return sample
}

func (a *AudioEngine) add(out []float64) {
	for i := range out {
		a.out[i] = out[i]
	}
}

func main() {
	err := portaudio.Initialize()
	if err != nil {
		log.Fatal(err)
	}
	defer portaudio.Terminate()

	a := NewAudioEngine()

	err = portmidi.Initialize()
	if err != nil {
		log.Println(err)
	}

	Midi, err = NewMidiStream(1, 1024)
	if err == nil {
		lasttime := portmidi.Timestamp(0)

		for event := range Midi.Listen() {
			Sin = staticSine(25 * float64(event.Data1))
			Square = staticSquare(25 * float64(event.Data1))

			MidiSep = event.Timestamp - lasttime
			MidiVelocity = event.Data2

			if MidiSep > 100 && MidiVelocity >= 0 {
				if err = a.Start(); err != nil {
					log.Println(err)
				}
			}
			lasttime = event.Timestamp
		}
	} else {
		log.Println(err)
		Sin = staticSquare(25 * float64(10))

		if err = a.Start(); err != nil {
			log.Println(err)
		}
	}
}
