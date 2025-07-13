package main

import (
	"os"
	"time"

	"github.com/gopxl/beep/speaker"
	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/wav"
)

func PlayWav(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	streamer, format, err := wav.Decode(f)
	if err != nil {
		return err
	}
	defer streamer.Close()
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))

	<-done
	return nil
}
