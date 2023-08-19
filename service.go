package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"player/player"
)

func main() {

	stations := []player.Station{
		{Name: "KEXP", URL: "https://kexp-mp3-128.streamguys1.com/kexp128.mp3"},
		{Name: "StuBru", URL: "http://icecast.vrtcdn.be/stubru-high.mp3"},
		{Name: "StuBru Bruut", URL: "http://icecast.vrtcdn.be/stubru_bruut-high.mp3"},
		{Name: "StuBru Untz", URL: "http://icecast.vrtcdn.be/stubru_untz-high.mp3"},
		{Name: "StuBru Hooray", URL: "http://icecast.vrtcdn.be/stubru_hiphophooray-high.mp3"},
	}

	radio := player.NewPlayer(stations)
	if err := radio.Init(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	go http.ListenAndServe(":8080", player.NewAPI(radio))

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
}
