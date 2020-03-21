package deezer

import "fmt"

type User struct {
	ID              int    `json:"id,omitempty"`               //	The user's Deezer ID
	Name            string `json:"name,omitempty"`             //	The user's Deezer nickname
	Lastname        string `json:"lastname,omitempty"`         //	The user's last name
	Firstname       string `json:"firstname,omitempty"`        //	The user's first name
	Email           string `json:"email,omitempty"`            //	The user's email
	Status          int    `json:"status,omitempty"`           //	The user's status
	Birthday        string `json:"birthday,omitempty"`         //	The user's birthday
	InscriptionDate string `json:"inscription_date,omitempty"` //	The user's inscription date
	Gender          string `json:"gender,omitempty"`           //	The user's gender : F or M
	Link            string `json:"link,omitempty"`             //	The url of the profil for the user on Deezer
	Picture         string `json:"picture,omitempty"`          //	The url of the user's profil picture.
	Country         string `json:"country,omitempty"`          //	The user's country
	Lang            string `json:"lang,omitempty"`             //	The user's language
	Tracklist       string `json:"tracklist,omitempty"`        //	API Link to the flow of this user
}

// CurrentUser Will return the user that own the token
func (c Client) CurrentUser() (user User, err error) {
	resp, err := c.http.R().SetQueryParam("access_token", c.token).Get(fmt.Sprintf("%s/user/me", c.baseURL))
	if err != nil {
		fmt.Println(err)
		return
	}

	err = c.http.JSONUnmarshal(resp.Body(), &user)
	return
}

// GetCurrentUserPlaylists will return the playlist attached to the given Users
func (c Client) GetCurrentUserPlaylists() (list PlaylistList, err error) {
	user, err := c.CurrentUser()
	if err != nil {
		return
	}
	resp, err := c.http.R().
		SetQueryParam("acces_token", c.token).
		SetHeader("Accept", "application/json").
		Get(fmt.Sprintf("%s/user/%d/playlists", c.baseURL, user.ID))

	err = c.http.JSONUnmarshal(resp.Body(), &list)
	return
}
