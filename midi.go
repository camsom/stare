package main

import "github.com/rakyll/portmidi"

func init() {
	portmidi.Initialize()
}

func DeviceCount() int {
	return portmidi.CountDevices()
}

func NewMidiStream(id portmidi.DeviceId, bufferSize int64) (*portmidi.Stream, error) {
	return portmidi.NewInputStream(id, bufferSize)
}
