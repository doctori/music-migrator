package deezer

import "encoding/json"

type Genre struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type ExtendedGenreList struct {
	Data []Genre `json:"data,omitempty"`
}

type GenreList []Genre

func (g *GenreList) UnmarshalJSON(data []byte) error {
	extendedGenreList := ExtendedGenreList{}
	if err := json.Unmarshal(data, &extendedGenreList); err != nil {
		return err
	}

	*g = extendedGenreList.Data

	return nil
}
