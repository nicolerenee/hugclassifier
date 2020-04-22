package cmd

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/nicolerenee/hugclassifier/pkg/classifier"
	"github.com/nicolerenee/hugclassifier/pkg/meetup"
	"github.com/nicolerenee/hugclassifier/pkg/onepass"
	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		generateCSV()
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().IntP("days", "d", 30, "Number of days back to search for events")
	viper.BindPFlag("duration.days", generateCmd.Flags().Lookup("days"))
	generateCmd.Flags().String("start", "", "start date for events, if set --days will be ignored. Example (2020-04-01)")
	viper.BindPFlag("duration.start", generateCmd.Flags().Lookup("start"))
	generateCmd.Flags().String("end", "", "end date for events, defaults to today")
	viper.BindPFlag("duration.end", generateCmd.Flags().Lookup("end"))
	generateCmd.Flags().String("output", "./results.csv", "file to store CSV output")
	viper.BindPFlag("output", generateCmd.Flags().Lookup("output"))
}

func generateCSV() {
	var dateFrom string
	now := time.Now().UTC()
	dateTo := now.Format(meetup.DateFormat)

	if viper.GetString("duration.start") != "" {
		// start date specified, use it.
		s, err := time.Parse("2006-01-02", viper.GetString("duration.start"))
		if err != nil {
			log.Fatal("error processing --start paramater: ", err)
		}
		dateFrom = s.Format(meetup.DateFormat)

		if viper.GetString("duration.end") != "" {
			e, err := time.Parse("2006-01-02", viper.GetString("duration.end"))
			if err != nil {
				log.Fatal("error processing --end paramater: ", err)
			}
			dateTo = e.Format(meetup.DateFormat)
		}
	} else {
		// no start date specified, use duration.days
		days := viper.GetInt("duration.days")
		dateFrom = now.AddDate(0, 0, 0-days).Format(meetup.DateFormat)
	}

	op, err := onepass.NewClient(viper.GetString("onepassword.account"))
	if err != nil {
		log.Fatal("error: ", err)
	}
	apiCreds, err := op.GetLogin(viper.GetString("onepassword.login"), viper.GetString("onepassword.vault"))
	if err != nil {
		log.Printf("error: 1Password failed, ensure you are logged in by running: eval $(op signin %s)\n", viper.GetString("onepassword.account"))
		log.Fatal("error: ", err)
	}

	m, err := meetup.NewClient(apiCreds.Username, apiCreds.Password)
	if err != nil {
		log.Fatal("error: ", err)
	}

	if err := m.Authenticate(); err != nil {
		log.Fatal("error: failed to auth to meetup: ", err)
	}

	log.Printf("info: successfully authenticated, retreiving data\n")
	log.Printf("info: collecting events between %s and %s\n", dateFrom, dateTo)

	var events []*meetup.Event
	groups, err := m.HUGs()
	if err != nil {
		log.Fatal("error: failed to retreive HUG groups from meetup: ", err)
	}

	for _, group := range groups {
		e, err := m.EventsForGroup(group.MeetupURLName, dateFrom, dateTo)
		if err != nil {
			log.Fatal("error: failed to retreive events for a group from meetup: ", err)
		}
		events = append(events, e...)
		// Keep from pounding the API to hard
		time.Sleep(250 * time.Millisecond)
	}

	file, err := os.Create(viper.GetString("output"))
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

	for _, event := range events {
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
	log.Printf("info: found a total of %d events in %d groups\n", len(events), len(groups))
}
