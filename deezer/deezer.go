package deezer

import (
	"fmt"
	"log"
	"net/http"

	"github.com/doctori/music-migrator/utils"
	"github.com/go-resty/resty/v2"
)

const baseAddress = "https://api.deezer.com"
const tokenPath = ".deez"

// Client is a client for working with the Deezer Web API.
// To create an authenticated client, use the `Authenticator.NewClient` method.
type Client struct {
	http      *resty.Client
	baseURL   string
	token     string
	AutoRetry bool
}

// deezerRedirectURI is the OAuth redirect URI for the application.
// You must register an application at Spotify's developer portal
// and enter this value.
const deezerRedirectURI = "http://localhost:8080/deezer-callback"

// Deezer holds deezer related infos
type Deezer struct {
	// Connected is true if the login phase has appened
	auth      Authenticator
	Client    *Client
	authURL   string
	Connected bool
	ch        chan *Client
	state     string
}

// Init will initial the deezer client
func (d *Deezer) Init() (err error) {
	d.auth = NewAuthenticator(deezerRedirectURI, ScopeBasic, ScopeOfflineAccess, ScopeManageLibrary)
	d.Connected = false
	d.ch = make(chan *Client)
	d.state, err = utils.GenerateRandomString(32)
	if err != nil {
		return err
	}

	d.authURL = d.auth.AuthURL(d.state)
	return
}

func (d *Deezer) Connect() {
	tok, err := utils.ReadToken(tokenPath)
	if err != nil {
		fmt.Println("Please log in to Deezer by visiting the following page in your browser:", d.authURL)
		// wait for auth to complete
		d.Client = <-d.ch

	} else {
		client := d.auth.NewClient(tok)
		d.Client = &client
	}

	// use the client to make calls that require authorization
	user, err := d.Client.CurrentUser()
	d.Connected = true
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("You are logged in as:", user.Name)
}

func (d Deezer) CompleteAuth(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Got request for: %s\n", r.URL.String())
	tok, err := d.auth.Token(d.state, r)
	if err != nil {
		fmt.Printf("Couldn't get the token because %s", err)
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	/*
		if st := r.FormValue("state"); st != d.state {
			http.NotFound(w, r)
			log.Fatalf("State mismatch: %s != %s\n", st, d.state)
		}
	*/
	// use the token to get an authenticated client
	client := d.auth.NewClient(tok)
	utils.SaveToken(tok, tokenPath)
	fmt.Fprintf(w, "Login Completed!")
	d.ch <- &client
}
