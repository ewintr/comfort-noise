package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

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

	radio := player.NewPlayer()
	if err := radio.Init(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Stations:")
	for i, stream := range stations {
		fmt.Printf("%d: %s\n", i, stream.Name)
	}

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

		station := stations[number]
		fmt.Printf("Playing %s\n", station.Name)

		if err := radio.Select(station); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

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
