package meetup

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Group represents a meetup group
type Group struct {
	MeetupURLName string `json:"urlname"`
	Name          string `json:"name"`
	Location      string `json:"localized_location"`
	Region        string `json:"region"`
	Timezone      string `json:"timezone"`
}

// HUGs will return a list of HUGs from the Meetup API. This is accomplished by
// using the API to get a list of all the groups that belong to the HUG pro
// network.
func (c *Client) HUGs() ([]*Group, error) {
	uri := fmt.Sprintf("%s/pro/%s/groups", baseURL, "hugs")
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", c.authValue())
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("received a %d status code", resp.StatusCode)
	}

	var groups []*Group

	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&groups)
	if err != nil {
		return nil, err
	}

	return groups, nil
}
