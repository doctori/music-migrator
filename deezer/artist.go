package deezer

type Artist struct {
	ID        int       `json:"id,omitempty"`        //	The artist's Deezer id
	Name      string    `json:"name,omitempty"`      //	name	The artist's name
	Link      string    `json:"link,omitempty"`      //	The url of the artist on Deezer
	Share     string    `json:"share,omitempty"`     //	The share link of the artist on Deezer
	Picture   string    `json:"picture,omitempty"`   //	The url of the artist picture.
	NbAlbum   int       `json:"nb_album,omitempty"`  //	The number of artist's albums
	NbFan     int       `json:"nb_fan,omitempty"`    //	The number of artist's fans
	Radio     bool      `json:"radio,omitempty"`     //	true if the artist has a smartradio
	Tracklist string    `json:"tracklist,omitempty"` //	API Link to the top of this artist
	Role      string    `json:"role,omitempty"`      //	The artist's role in a track or album
	Albums    AlbumList `json:"albums,omitempty"`
	Tracks    TrackList `json:"tracks,omitempty"`
}
