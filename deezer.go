package main

import (
	"fmt"
)

func printDeezerPlaylists() {

	playlists, err := d.Client.GetCurrentUserPlaylists()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, playlist := range playlists {
		fmt.Println(playlist.Title)
	}
}
