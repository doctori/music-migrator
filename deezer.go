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

func printDeezerLovedTracks() {

	playlists, err := d.Client.GetCurrentUserPlaylists()
	if err != nil {
		fmt.Printf("%#v\n", err)
		return
	}
	for _, playlist := range playlists {
		if playlist.IsLovedTrack {
			pl, err := d.Client.GetPlaylist(playlist.ID)
			if err != nil {
				fmt.Println(err)
				return
			}
			for _, track := range pl.Tracks {
				fmt.Printf("%s - %s - %s\n", track.Title, track.Artist.Name, track.Album.Title)
			}
		}
	}
}
