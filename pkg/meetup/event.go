package meetup

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// Event represents an event as returned from the Meetup API. Not all fields are
// available, but the ones we needed have been added, along with a few others
// that might be useful in the future.
type Event struct {
	MilliTime    int64       `json:"time"`
	Duration     int         `json:"duration"` // in milliseconds
	Name         string      `json:"name"`
	ID           string      `json:"id"`
	Status       string      `json:"status"`
	LocalDate    string      `json:"local_date"`
	LocalTime    string      `json:"local_time"`
	UTCOffset    int         `json:"utc_offset"` // in milliseconds
	YesRSVPCount int         `json:"yes_rsvp_count"`
	Link         string      `json:"link"`
	Description  string      `json:"description"`
	Group        Group       `json:"Group"`
	EventHosts   []EventHost `json:"event_hosts"`
}

// EventHost contains the data for the host of an event as returned from the
// Meetup API
type EventHost struct {
	Name string `json:"name"`
}

// EventsForGroup will return the events for a group given the group name, which
// is the url portion of the meetup group.
func (c *Client) EventsForGroup(groupName string, dateFrom, dateTo string) ([]*Event, error) {
	uri := fmt.Sprintf("%s/%s/events?no_earlier_than=%s&no_later_than=%s&fields=event_hosts&status=upcoming,past",
		baseURL, groupName, dateFrom, dateTo)

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}
	// req.Header.Add("Authorization", c.authValue())
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	var events []*Event

	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&events)
	switch {
	case err == io.EOF:
		// empty body no events found
		return events, nil
	case err != nil:
		fmt.Println("ERROR")
		fmt.Println(uri)
		fmt.Println(resp.StatusCode)
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(body))

		return nil, err
	}

	return events, nil
}

// Time returns the time of an event. Converting the millisecond unix time that
// is returned from the API to a valid go time object
func (e *Event) Time() time.Time {
	secs := e.MilliTime / 1000 // convert milliseconds to seconds
	return time.Unix(secs, 0)
}

// Hosts returns the a comma seperated string with all the hosts for the event
func (e *Event) Hosts() string {
	hosts := []string{}
	for _, h := range e.EventHosts {
		hosts = append(hosts, h.Name)
	}
	return strings.Join(hosts, ", ")
}
