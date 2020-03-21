package utils

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"golang.org/x/oauth2"
)

func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func GenerateRandomString(s int) (string, error) {
	b, err := generateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}

// SaveToken Will save the oauth2 token
func SaveToken(tok *oauth2.Token, tokenPath string) error {
	f, err := os.OpenFile(tokenPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	return enc.Encode(tok)
}

// ReadToken Will read the Oauth2 token
func ReadToken(tokenPath string) (*oauth2.Token, error) {
	content, err := ioutil.ReadFile(tokenPath)
	if err != nil {
		return nil, err
	}

	var tok oauth2.Token
	if err := json.Unmarshal(content, &tok); err != nil {
		return nil, err
	}
	if !tok.Valid() {
		return nil, fmt.Errorf("Token isn't valid anymore")
	}

	return &tok, nil
}
