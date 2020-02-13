package meetup

import (
	"net/http"

	"golang.org/x/oauth2"
)

const (
	baseURL string = "https://api.meetup.com"
)

// Client provides a meetup API client. A valid token is required to be able to
// make authenticated API calls, such as getting the list of groups belonging to
// a pro network
type Client struct {
	Token       string
	httpClient  *http.Client
	oauthConfig *oauth2.Config
}
