package player

import (
	"fmt"
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

type Status struct {
	Playing  bool
	Station  string
	Channels []string
}

type Player struct {
	stations    []Station
	current     int
	oldStreamer beep.StreamCloser
	oldCtrl     *beep.Ctrl
	mixer       *beep.Mixer
	sr          beep.SampleRate
}

func NewPlayer(stations []Station) *Player {
	return &Player{
		stations: stations,
		current:  -1,
		mixer:    &beep.Mixer{},
		sr:       beep.SampleRate(48000),
	}
}

func (p *Player) Init() error {

	if err := speaker.Init(p.sr, p.sr.N(time.Second/10)); err != nil {
		return err
	}

	speaker.Play(p.mixer)

	return nil
}

func (p *Player) Select(channel int) error {
	if channel < 0 || channel >= len(p.stations) {
		return fmt.Errorf("unknown channel: %d", channel)
	}

	res, err := http.Get(p.stations[channel].URL)
	if err != nil {
		return err
	}

	streamer, format, err := mp3.Decode(res.Body)
	if err != nil {
		return err
	}

	resampled := Resample(4, format.SampleRate, p.sr, streamer)

	p.PlayStream(resampled)
	p.current = channel

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

func (p *Player) Stop() {
	speaker.Lock()
	if p.oldCtrl != nil {
		p.oldCtrl.Paused = true
		p.oldCtrl.Streamer = nil
		p.oldStreamer.Close()
	}
	speaker.Unlock()

	p.current = -1
}

func (p *Player) Status() Status {
	channels := make([]string, len(p.stations))
	for i, stat := range p.stations {
		channels[i] = stat.Name
	}

	if p.current < 0 || p.current >= len(p.stations) {
		return Status{Channels: channels}
	}

	return Status{
		Playing:  true,
		Station:  p.stations[p.current].Name,
		Channels: channels,
	}
}
