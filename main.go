package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

type StreamURL struct {
	Name string
	URL  string
}

var StreamURLs = []StreamURL{
	{Name: "KEXP", URL: "https://kexp-mp3-128.streamguys1.com/kexp128.mp3"},
	{Name: "StuBru", URL: "http://icecast.vrtcdn.be/stubru-high.mp3"},
	{Name: "StuBru Bruut", URL: "http://icecast.vrtcdn.be/stubru_bruut-high.mp3"},
	{Name: "StuBru Untz", URL: "http://icecast.vrtcdn.be/stubru_untz-high.mp3"},
	{Name: "StuBru Hooray", URL: "http://icecast.vrtcdn.be/stubru_hiphophooray-high.mp3"},
}

func main() {

	// setting sample rate
	sr := beep.SampleRate(48000)
	if err := speaker.Init(sr, sr.N(time.Second/10)); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Stations:")
	for i, stream := range StreamURLs {
		fmt.Printf("%d: %s\n", i, stream.Name)
	}

	mix := &beep.Mixer{}

	speaker.Play(mix)

	var oldCtrl *beep.Ctrl

	for {

		fmt.Println("Press a number to change stations, 9 to quit")

		number, err := getNext()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if number == 9 {
			os.Exit(0)
		}

		stream := StreamURLs[number]
		fmt.Printf("Playing %s\n", stream.Name)
		res, err := http.Get(stream.URL)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		streamer, _, err := mp3.Decode(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		defer streamer.Close()
		ctrl := &beep.Ctrl{Streamer: streamer, Paused: false}
		mix.Add(ctrl)
		speaker.Lock()
		if oldCtrl != nil {
			oldCtrl.Paused = true
			oldCtrl.Streamer = nil
		}
		ctrl.Paused = true
		ctrl.Streamer = streamer
		ctrl.Paused = false
		speaker.Unlock()

		oldCtrl = ctrl
	}
}

func getNext() (int, error) {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return 9, err
	}

	number, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil {
		return 9, err
	}

	return number, nil
}
