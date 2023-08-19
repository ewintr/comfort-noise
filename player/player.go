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
	sr          beep.SampleRate
}

func NewPlayer() *Player {
	return &Player{
		mixer: &beep.Mixer{},
		sr:    beep.SampleRate(48000),
	}
}

func (p *Player) Init() error {

	if err := speaker.Init(p.sr, p.sr.N(time.Second/10)); err != nil {
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

	streamer, format, err := mp3.Decode(res.Body)
	if err != nil {
		return err
	}

	resampled := Resample(4, format.SampleRate, p.sr, streamer)

	p.PlayStream(resampled)

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
