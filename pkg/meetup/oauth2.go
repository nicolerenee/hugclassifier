package meetup

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/oauth2"
)

const (
	authURL     string = "https://secure.meetup.com/oauth2/authorize"
	tokenURL    string = "https://secure.meetup.com/oauth2/access"
	oauthState  string = "vftxcvo546yhg35bvg42h2489evijodsm"
	redirectURL string = "http://127.0.0.1:14565/oauth/callback"
)

// NewClient will configure a client with the ID and Secret you pass in
func NewClient(clientID, clientSecret string) (*Client, error) {
	c := &Client{}
	c.oauthConfig = &oauth2.Config{
		RedirectURL:  redirectURL,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{"basic"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  authURL,
			TokenURL: tokenURL,
		},
	}
	c.httpClient = &http.Client{
		Timeout: time.Second * 10,
	}
	return c, nil
}

func (c *Client) handleOAuthCallback(w http.ResponseWriter, r *http.Request, tc chan string) {
	state := r.FormValue("state")
	if state != oauthState {
		fmt.Println("ERROR: oauth state doesn't match")
		w.WriteHeader(http.StatusBadRequest)
		tc <- ""
	}
	code := r.FormValue("code")

	token, err := c.oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		fmt.Printf("ERROR: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		tc <- ""
	}

	w.Write([]byte("Success. You can now close this window."))
	tc <- token.AccessToken
}

// Authenticate starts a server and listens for an oauth2 callback and will
// return the API token to the caller
func (c *Client) Authenticate() error {
	tc := make(chan string)

	server := &http.Server{Addr: ":14565"}
	http.HandleFunc("/oauth/callback", func(w http.ResponseWriter, r *http.Request) {
		c.handleOAuthCallback(w, r, tc)
	})

	go func() {
		if err := server.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				return
			}
			fmt.Printf("ERROR: %s\n", err.Error())
			tc <- ""
		}
	}()

	fmt.Printf("To authenticate visit: %s\n", c.oauthConfig.AuthCodeURL(oauthState))

	token := <-tc

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		return err
	}

	if token == "" {
		return errors.New("failed to get a token")

	}

	c.Token = token

	return nil
}

func (c *Client) authValue() string {
	return fmt.Sprintf("Bearer %s", c.Token)
}
