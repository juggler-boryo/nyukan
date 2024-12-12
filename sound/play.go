package sound

import (
	"log"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

var (
	bb      = "sound/bb.mp3"
	in      = "sound/in.mp3"
	connect = "sound/connect.mp3"
	error   = "sound/error.mp3"
	try     = "sound/try.mp3"
)

func init() {
	// Initialize speaker with standard sample rate
	sr := beep.SampleRate(44100)
	err := speaker.Init(sr, sr.N(time.Millisecond*100))
	if err != nil {
		log.Printf("Error initializing speaker: %v", err)
	}
}

// Play plays the specified sound file and waits for it to complete
func Play(soundFile string) {
	f, err := os.Open(soundFile)
	if err != nil {
		log.Printf("Error opening sound file: %v", err)
		return
	}
	defer f.Close()

	streamer, _, err := mp3.Decode(f)
	if err != nil {
		log.Printf("Error decoding MP3: %v", err)
		return
	}
	defer streamer.Close()

	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))

	// Wait for playback to complete
	<-done
}

// Convenience functions for playing specific sounds
func PlayBB()      { Play(bb) }
func PlayIn()      { Play(in) }
func PlayConnect() { Play(connect) }
func PlayError()   { Play(error) }
func PlayTry()     { Play(try) }
