package deezer

import "encoding/json"

type Playlist struct {
	ID            int       `json:"id,omitempty"`             // The playlist's Deezer id
	Title         string    `json:"title,omitempty"`          // The playlist's title
	Description   string    `json:"description,omitempty"`    // The playlist description
	Duration      int       `json:"duration,omitempty"`       // The playlist's duration (seconds)
	Public        bool      `json:"public,omitempty"`         // If the playlist is public or not
	IsLovedTrack  bool      `json:"is_loved_track,omitempty"` // If the playlist is the love tracks playlist
	Collaborative bool      `json:"collaborative,omitempty"`  // If the playlist is collaborative or not
	Rating        int       `json:"rating,omitempty"`         // The playlist's rate
	Fans          int       `json:"fans,omitempty"`           // The number of playlist's fans
	Link          string    `json:"link,omitempty"`           // The url of the playlist on Deezer
	Share         string    `json:"share,omitempty"`          // The share link of the playlist on Deezer
	Picture       string    `json:"picture,omitempty"`        // The url of the playlist's cover.
	Checksum      string    `json:"hecksum,omitempty"`        // The checksum for the track list
	Creator       *User     `json:"creator,omitempty"`        // User object containing : id, name
	Tracks        TrackList `json:"tracks,omitempty"`         // List of tracks
}

type extendedPlaylistList struct {
	Data     []Playlist `json:"data,omitempty"`
	Total    int        `json:"total,omitempty"`
	Checksum string     `json:"checksum,omitempty"`
	Next     string     `json:"next,omitempty"`
}

type PlaylistList []Playlist

func (p *PlaylistList) UnmarshalJSON(data []byte) error {
	ePlaylistList := extendedPlaylistList{}
	if err := json.Unmarshal(data, &ePlaylistList); err != nil {
		return err
	}

	*p = ePlaylistList.Data

	return nil
}
