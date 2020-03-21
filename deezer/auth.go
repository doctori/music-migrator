package deezer

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/go-resty/resty/v2"
	"golang.org/x/oauth2"
)

const (
	// AuthURL is the URL to Deezer Accounts Service's OAuth2 endpoint.
	AuthURL = "https://connect.deezer.com/oauth/auth.php"
	// TokenURL is the URL to the Deezer Accounts Service's OAuth2
	// token endpoint.
	TokenURL = "https://connect.deezer.com/oauth/access_token.php"
)

// Scopes let you specify exactly which types of data your application wants to access.
// The set of scopes you pass in your authentication request determines what access the
// permissions the user is asked to grant.
const (
	// ScopeBasic Access users basic information
	// Incl. name, firstname, profile picture only.
	ScopeBasic = "basic_access"
	// ScopeEmail Get the user's email
	ScopeEmail = "email"
	// ScopeOfflineAccess Access user data any time
	// Application may access user data at any time
	ScopeOfflineAccess = "offline_access"
	// ScopeManageLibrary Manage users' library
	// Add/rename a playlist. Add/order songs in the playlist.
	ScopeManageLibrary = "manage_library"
)

// Authenticator provides convenience functions for implementing the OAuth2 flow.
// You should always use `NewAuthenticator` to make them.
//
// Example:
//
//     a := deezer.NewAuthenticator(redirectURL, deezer.ScopeUserLibaryRead, deezer.ScopeUserFollowRead)
//     // direct user to Deezer to log in
//     http.Redirect(w, r, a.AuthURL("state-string"), http.StatusFound)
//
//     // then, in redirect handler:
//     token, err := a.Token(state, r)
//     client := a.NewClient(token)
//
type Authenticator struct {
	config  *oauth2.Config
	context context.Context
}

// NewAuthenticator creates an authenticator which is used to implement the
// OAuth2 authorization flow.  The redirectURL must exactly match one of the
// URLs specified in your Deezer developer account.
//
// By default, NewAuthenticator pulls your client ID and secret key from the
// SPOTIFY_ID and SPOTIFY_SECRET environment variables.  If you'd like to provide
// them from some other source, you can call `SetAuthInfo(id, key)` on the
// returned authenticator.
func NewAuthenticator(redirectURL string, scopes ...string) Authenticator {
	cfg := &oauth2.Config{
		ClientID:     os.Getenv("DEEZER_ID"),
		ClientSecret: os.Getenv("DEEZER_SECRET"),
		RedirectURL:  redirectURL,
		Scopes:       scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  AuthURL,
			TokenURL: TokenURL,
		},
	}

	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, &http.Client{})
	return Authenticator{
		config:  cfg,
		context: ctx,
	}
}

// SetAuthInfo overwrites the client ID and secret key used by the authenticator.
// You can use this if you don't want to store this information in environment variables.
func (a *Authenticator) SetAuthInfo(clientID, secretKey string) {
	a.config.ClientID = clientID
	a.config.ClientSecret = secretKey
}

// AuthURL returns a URL to the the Deezer Accounts Service's OAuth2 endpoint.
//
// State is a token to protect the user from CSRF attacks.  You should pass the
// same state to `Token`, where it will be validated.  For more info, refer to
// http://tools.ietf.org/html/rfc6749#section-10.12.
func (a Authenticator) AuthURL(state string) string {
	return a.config.AuthCodeURL(state)
}

// AuthURLWithDialog returns the same URL as AuthURL, but sets show_dialog to true
func (a Authenticator) AuthURLWithDialog(state string) string {
	return a.config.AuthCodeURL(state, oauth2.SetAuthURLParam("show_dialog", "true"))
}

// Token pulls an authorization code from an HTTP request and attempts to exchange
// it for an access token.  The standard use case is to call Token from the handler
// that handles requests to your application's redirect URL.
func (a Authenticator) Token(state string, r *http.Request) (*oauth2.Token, error) {
	values := r.URL.Query()
	if e := values.Get("error"); e != "" {
		return nil, errors.New("deezer: auth failed - " + e)
	}
	code := values.Get("code")
	if code == "" {
		return nil, errors.New("deezer: didn't get access code")
	}
	actualState := values.Get("state")
	if actualState != state {
		return nil, errors.New("deezer: redirect state parameter doesn't match")
	}
	fmt.Printf("Exchanging with this code %s\n", code)
	return a.config.Exchange(oauth2.NoContext, code)
}

// Exchange is like Token, except it allows you to manually specify the access
// code instead of pulling it out of an HTTP request.
func (a Authenticator) Exchange(code string) (*oauth2.Token, error) {
	return a.config.Exchange(a.context, code)
}

// NewClient creates a Client that will use the specified access token for its API requests.
func (a Authenticator) NewClient(token *oauth2.Token) Client {
	client := a.config.Client(a.context, token)
	restyClient := resty.NewWithClient(client)
	restyClient.SetAuthToken(token.AccessToken)

	return Client{
		http:    restyClient,
		token:   token.AccessToken,
		baseURL: baseAddress,
	}
}

// Token gets the client's current token.
func (c *Client) Token() (*oauth2.Token, error) {
	transport, ok := c.http.GetClient().Transport.(*oauth2.Transport)
	if !ok {
		return nil, errors.New("deezer: oauth2 transport type not correct")
	}
	t, err := transport.Source.Token()
	if err != nil {
		return nil, err
	}

	return t, nil
}
