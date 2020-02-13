package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/nicolerenee/hugclassifier/pkg/classifier"
	"github.com/nicolerenee/hugclassifier/pkg/meetup"
	"github.com/nicolerenee/hugclassifier/pkg/onepass"
)

const (
	onepassAccount        = "hashicorp"
	onepassVault          = "vfa4wwzimconvcuafzvecxzera" // UUID for Developer Advocate Group vault
	onepassMeetupAPILogin = "f7ih4yyotjal7kcobvqa7er3lu" // UUID for Meetup Oauth Creds login item
	// onepassMeetupLogin    = "mn2hjkb5jvf4lcrfzzb7ndd5jq" // UUID for shared Meetup login
)

func main() {
	op, err := onepass.NewClient(onepassAccount)
	if err != nil {
		log.Fatal("error: ", err)
	}
	apiCreds, err := op.GetLogin(onepassMeetupAPILogin, onepassVault)
	if err != nil {
		log.Fatal("error: ", err)
	}

	m, err := meetup.NewClient(apiCreds.Username, apiCreds.Password)
	if err != nil {
		log.Fatal("error: ", err)
	}

	if err := m.Authenticate(); err != nil {
		log.Fatal("error: failed to auth to meetup: ", err)
	}

	fmt.Fprintf(os.Stdout, "info: successfully authenticated, retreiving data\n")

	var events []*meetup.Event
	groups, err := m.HUGs()
	if err != nil {
		log.Fatal("error: failed to retreive HUG groups from meetup: ", err)
	}

	for _, group := range groups {
		e, err := m.EventsForGroup(group.MeetupURLName)
		if err != nil {
			log.Fatal("error: failed to retreive events for a group from meetup: ", err)
		}
		events = append(events, e...)
	}

	file, err := os.Create("result.csv")
	if err != nil {
		log.Fatal("cannot create file", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"Group", "Event Name", "Event Hosts", "Event UTC Time", "Local Date", "Local Time", "Status", "Yes RSVPs", "Link", "Topics", "Description"}
	err = writer.Write(header)
	if err != nil {
		log.Fatal("failed to write to csv file", err)
	}

	// fmt.Println("=====================================================")
	for _, event := range events {
		// fmt.Printf("       Group: %s\n", event.Group.Name)
		// fmt.Printf("        Name: %s\n", event.Name)
		// fmt.Printf("        Date: %s\n", event.Time())
		// fmt.Printf("  Local Time: %s\n", event.LocalTime)
		// fmt.Printf("      Status: %s\n", event.Status)
		// fmt.Printf("  RSVP Count: %d\n", event.YesRSVPCount)
		// fmt.Printf("        Link: %s\n", event.Link)
		// fmt.Println("=====================================================")

		values := []string{
			event.Group.Name,
			event.Name,
			event.Hosts(),
			event.Time().UTC().String(),
			event.LocalDate,
			event.LocalTime,
			event.Status,
			fmt.Sprintf("%d", event.YesRSVPCount),
			event.Link,
			strings.Join(classifier.Classify(event.Description), ", "),
			event.Description,
		}
		err := writer.Write(values)
		if err != nil {
			log.Fatal("failed to write to file", err)
		}
	}
	fmt.Printf("Found a total of %d events in %d groups\n", len(events), len(groups))
}
