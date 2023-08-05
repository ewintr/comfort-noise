package server

import (
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
)

type Player struct {
	oldStreamer beep.StreamCloser
	oldCtrl     *beep.Ctrl
	mixer       *beep.Mixer
}

func NewPlayer() *Player {
	return &Player{
		mixer: &beep.Mixer{},
	}
}

func (p *Player) Init() error {
	sr := beep.SampleRate(48000)
	if err := speaker.Init(sr, sr.N(time.Second/10)); err != nil {
		return err
	}

	speaker.Play(p.mixer)

	return nil
}

func (p *Player) Play(streamer beep.StreamCloser) {
	ctrl := &beep.Ctrl{
		Streamer: streamer,
	}
	p.mixer.Add(ctrl)

	speaker.Lock()

	if p.oldCtrl != nil {
		p.oldCtrl.Paused = true
		p.oldCtrl.Streamer = nil
		p.oldStreamer.Close()
	}

	ctrl.Paused = true
	ctrl.Streamer = streamer
	ctrl.Paused = false

	speaker.Unlock()

	p.oldCtrl = ctrl
	p.oldStreamer = streamer
}
