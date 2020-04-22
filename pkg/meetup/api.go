package meetup

import (
	"net/http"

	"golang.org/x/oauth2"
)

const (
	baseURL string = "https://api.meetup.com"
	// DateFormat provides a formatting string for dates passed into the Meetup API
	DateFormat string = "2006-01-02T15:04:05.999"
)

// Client provides a meetup API client. A valid token is required to be able to
// make authenticated API calls, such as getting the list of groups belonging to
// a pro network
type Client struct {
	Token       string
	httpClient  *http.Client
	oauthConfig *oauth2.Config
}
