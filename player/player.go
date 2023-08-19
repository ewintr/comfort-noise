package player

import (
	"net/http"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

type Station struct {
	Name string
	URL  string
}

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

func (p *Player) Select(station Station) error {
	res, err := http.Get(station.URL)
	if err != nil {
		return err
	}

	streamer, _, err := mp3.Decode(res.Body)
	if err != nil {
		return err
	}

	p.PlayStream(streamer)

	return nil
}

func (p *Player) PlayStream(streamer beep.StreamCloser) {
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
