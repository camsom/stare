package main

import (
	"bytes"
	"encoding/binary"
	"io"
	"log"
	"os/exec"

	"code.google.com/p/portaudio-go/portaudio"
)

type AudioStream struct {
	Stream *portaudio.Stream
}

type ffmpeg struct {
	in  io.ReadCloser
	cmd *exec.Cmd
}

func newFfmpeg(filename string) *ffmpeg {
	cmd := exec.Command("ffmpeg", "-i", filename, "-f", "s16le", "-")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	return &ffmpeg{stdout, cmd}
}

func (f *ffmpeg) Close() error {
	return f.in.Close()
}

func (f *ffmpeg) ProcessAudio(_, out [][]int16) {
	// int16 takes 2 bytes
	bufferSize := len(out[0]) * 4
	var pack = make([]byte, bufferSize)
	if _, err := f.in.Read(pack); err != nil {
		log.Fatal(err)
	}
	n := make([]int16, len(out[0])*2)
	for i := range n {
		var x int16
		buf := bytes.NewBuffer(pack[2*i : 2*(i+1)])
		binary.Read(buf, binary.LittleEndian, &x)
		n[i] = x
	}

	for i := range out[0] {
		out[0][i] = n[2*i]
		out[1][i] = n[2*i+1]
	}
}

func main() {
	err := portaudio.Initialize()
	if err != nil {
		log.Fatal(err)
	}

	// ff := newFfmpeg("Thots.mp3")
	// defer ff.Close()

	stream, err := portaudio.OpenDefaultStream(0, 2, 44100, 2048)
	if err != nil {
		log.Println(err)
	}
	defer stream.Close()

	// err = stream.Start()
	// if err != nil {
	// 	log.Println(err)
	// }

	// if err := ff.cmd.Wait(); err != nil {
	// 	log.Fatal(err)
	// }

	err = stream.Stop()
	if err != nil {
		log.Println(err)
	}
}
