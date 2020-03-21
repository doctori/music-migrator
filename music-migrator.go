package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/doctori/music-migrator/deezer"
)

var s Spotify
var d deezer.Deezer

func mainCompleter(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "spotify", Description: "spotify related stuff"},
		{Text: "deezer", Description: "deezer related stuff"},
		{Text: "quit", Description: "leave me alone !"},
	}
	operands := strings.Split(d.TextBeforeCursor(), " ")
	if len(operands) > 1 {
		switch operands[0] {
		case "spotify":
			s = []prompt.Suggest{
				{Text: "login", Description: "login to spotify"},
				{Text: "playlist", Description: "list playlists"},
				{Text: "loved-tracks", Description: "list all the loved tracks"},
			}
		case "deezer":
			s = []prompt.Suggest{
				{Text: "login", Description: "login do deezer"},
				{Text: "playlist", Description: "list playlists"},
				{Text: "loved-tracks", Description: "List all the loved tracks"},
			}
		}

	}

	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}
func mainSwitcher(input string) {

	switch strings.Trim(input, " ") {
	case "quit":
		fmt.Println("bye bye ")
		os.Exit(0)
	case "spotify":
		fmt.Println("Spotify mode")
	case "spotify login":
		s.Connect()

	case "spotify playlist":
		if !s.Connected {
			s.Connect()
		}
		s.ListPlaylists()
	case "spotify loved-tracks":
		if !s.Connected {
			s.Connect()
		}
		s.PrintSpotifyLovedTracks()
	case "deezer login":
		d.Connect()
	case "deezer playlist":
		if !d.Connected {
			d.Connect()
		}
		printDeezerPlaylists()
	case "deezer loved-tracks":
		if !d.Connected {
			d.Connect()
		}
		printDeezerLovedTracks()
	default:
		fmt.Printf("What do you want me to do with %s\n", input)
	}

}
func main() {

	fmt.Println("Hello world")

	p := prompt.New(
		mainSwitcher,
		mainCompleter,
	)
	p.Run()
	os.Exit(0)
}

func init() {
	d.Init()
	s.Init()
	// first start an HTTP server
	http.HandleFunc("/spotify-callback", s.completeAuth)
	http.HandleFunc("/deezer-callback", d.CompleteAuth)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})
	go http.ListenAndServe("127.0.0.1:8080", nil)
}
