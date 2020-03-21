package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/doctori/music-migrator/utils"
	"github.com/zmb3/spotify"
)

// spotifyRedirectURI is the OAuth redirect URI for the application.
// You must register an application at Spotify's developer portal
// and enter this value.
const spotifyRedirectURI = "http://localhost:8080/spotify-callback"
const spotifyTokenPath = ".spot"

// Spotify holds spotify related infos
type Spotify struct {
	auth   spotify.Authenticator
	client *spotify.Client
	// Connected is true if the login phase has appened
	Connected bool
	authURL   string
	ch        chan *spotify.Client
	state     string
}

// Init Initialize the spotify client
func (s *Spotify) Init() (err error) {
	s.auth = spotify.NewAuthenticator(
		spotifyRedirectURI,
		spotify.ScopeUserReadPrivate,
		spotify.ScopePlaylistReadPrivate,
		spotify.ScopePlaylistReadCollaborative,
		spotify.ScopeUserLibraryRead,
	)
	s.Connected = false
	s.ch = make(chan *spotify.Client)
	s.state, err = utils.GenerateRandomString(32)
	if err != nil {
		return err
	}
	s.authURL = s.auth.AuthURL(s.state)
	return
}

// Connect will proceed to the authentication of the user
func (s *Spotify) Connect() {
	tok, err := utils.ReadToken(spotifyTokenPath)
	if err != nil {
		fmt.Println("Please log in to Spotify by visiting the following page in your browser:", s.authURL)
		// wait for auth to complete
		s.client = <-s.ch
	} else {
		client := s.auth.NewClient(tok)
		s.client = &client
	}

	// use the client to make calls that require authorization
	user, err := s.client.CurrentUser()
	s.Connected = true
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("You are logged in as:", user.ID)
}

// ListPlaylists will simply list all playlist of the current user
func (s *Spotify) ListPlaylists() {

	// fetching all playlists with pagination
	playlistCounter := 0

	offset := 0
	limit := 50
	opt := spotify.Options{
		Offset: &offset,
		Limit:  &limit,
	}

	for {
		playlists, err := s.client.CurrentUsersPlaylistsOpt(&opt)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		for _, playlist := range playlists.Playlists {
			playlistCounter++
			fmt.Println(playlist.Name)
		}
		if playlistCounter == playlists.Total {
			break
		}
		offset += limit
	}

}
func (s *Spotify) PrintSpotifyLovedTracks() {
	trackCounter := 0

	offset := 0
	limit := 50
	opt := spotify.Options{
		Offset: &offset,
		Limit:  &limit,
	}
	s.client.CurrentUsersTracks()
	for {
		tracks, err := s.client.CurrentUsersTracksOpt(&opt)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		for _, track := range tracks.Tracks {
			trackCounter++
			fmt.Printf("%s - %s - %s\n", track.Name, track.Artists[0].Name, track.Album.Name)
		}
		if trackCounter == tracks.Total {
			break
		}
		offset += limit
	}
}

func (s Spotify) completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := s.auth.Token(s.state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}

	if st := r.FormValue("state"); st != s.state {
		http.NotFound(w, r)
		fmt.Printf("State mismatch: %s != %s\n", st, s.state)
		log.Fatalf("State mismatch: %s != %s\n", st, s.state)
	}

	// use the token to get an authenticated client
	client := s.auth.NewClient(tok)
	fmt.Fprintf(w, "Login Completed!")
	fmt.Printf("Login Completed!")
	utils.SaveToken(tok, spotifyTokenPath)
	s.ch <- &client
}
