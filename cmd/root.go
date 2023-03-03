package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/jinzhu/now"
	"github.com/joho/godotenv"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

var (
	googleTokenPath  string
	googleCalendarId string
)

var rate int
var month int
var year int
var eventPrefix string

var calSvc *calendar.Service

var rootCmd = &cobra.Command{
	Use:   "walks",
	Short: "walks - a simple CLI to calculate monthly payment for dog walks based on calendar events",
	Long: `walks is a CLI built to pull dog walk events from google calendar
and calcualte the monthly due for those walks.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		b, err := ioutil.ReadFile(googleTokenPath)
		if err != nil {
			osExitErr(fmt.Sprintf("Unable to read credentials file from path: '%s'", err))
		}
		conf, err := google.JWTConfigFromJSON(b, calendar.CalendarReadonlyScope)
		if err != nil {
			osExitErr(fmt.Sprintf("Unable to load the JWT config from credentials json in environment: '%s'", err))
		}

		client := conf.Client(oauth2.NoContext)
		calSvc, err = calendar.NewService(oauth2.NoContext, option.WithHTTPClient(client))
		if err != nil {
			osExitErr(fmt.Sprintf("Unable to retrieve Calendar client: '%s'", err))
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Start Date", "End Date", "Count", "Rate", "Amount"})
		table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
		table.SetCenterSeparator("|")

		t := time.Date(year, time.Month(month), time.Now().Day(), time.Now().Hour(), time.Now().Minute(), 0, 0, time.Now().Location())

		beginningOfMonth := now.With(t).BeginningOfMonth()
		endOfMonth := now.With(t).EndOfMonth()

		events, err := calSvc.Events.
			List("primary").
			Q(eventPrefix).
			SingleEvents(true).
			ShowDeleted(false).
			TimeMin(beginningOfMonth.Format(time.RFC3339)).
			TimeMax(endOfMonth.Format(time.RFC3339)).
			MaxResults(2500).
			TimeZone(time.Now().Location().String()).
			Do()
		if err != nil {
			osExitErr(fmt.Sprintf("Unable to retrieve list of events from calendar: '%s'", err))
		}

		table.Append([]string{beginningOfMonth.Format("January 02 2006"),
			endOfMonth.Format("January 02 2006"),
			strconv.Itoa(len(events.Items)),
			strconv.Itoa(rate),
			strconv.Itoa(len(events.Items) * rate),
		})
		table.Render()
	},
}

var detailsCmd = &cobra.Command{
	Use:   "details",
	Short: "Print the day by day details of the walk count.",
	Long:  `Usually the dog walker reports the days walked, this is used for comparison against those counts.`,
	Run: func(cmd *cobra.Command, args []string) {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Month", "Date", "Count", "Amount"})
		table.SetBorder(false)
		table.SetAutoMergeCellsByColumnIndex([]int{0})
		table.SetCenterSeparator("|")

		t := time.Date(year, time.Month(month), time.Now().Day(), time.Now().Hour(), time.Now().Minute(), 0, 0, time.Now().Location())

		beginningOfMonth := now.With(t).BeginningOfMonth()
		endOfMonth := now.With(t).EndOfMonth()

		events, err := calSvc.Events.
			List("primary").
			Q(eventPrefix).
			SingleEvents(true).
			ShowDeleted(false).
			TimeMin(beginningOfMonth.Format(time.RFC3339)).
			TimeMax(endOfMonth.Format(time.RFC3339)).
			MaxResults(2500).
			TimeZone(time.Now().Location().String()).
			Do()
		if err != nil {
			osExitErr(fmt.Sprintf("Unable to retrieve list of events from calendar: '%s'", err))
		}

		eventDates := []time.Time{}
		for _, event := range events.Items {
			endWalkTime, err := time.Parse("2006-01-02T15:04:05Z", event.End.DateTime)
			if err != nil {
				log.Printf("error parsing time for end time: %v", err)
			}
			eventDates = append(eventDates, endWalkTime)
		}

		sort.Slice(eventDates, func(i, j int) bool {
			return eventDates[i].Before(eventDates[j])
		})

		for _, eventDate := range eventDates {
			table.Append([]string{eventDate.Month().String(), strconv.Itoa(eventDate.Day()), "1", strconv.Itoa(rate)})
		}
		table.SetFooter([]string{"", "", "Total", strconv.Itoa(len(events.Items) * rate)})
		table.Render()
	},
}

func init() {
	var set bool

	envFileNameArg := "setupDev.sh"
	err := godotenv.Load(envFileNameArg) //probably rip all this out and put in viper
	if err != nil {
		log.Fatalf("Error loading env file named: %v", envFileNameArg)
	}

	googleTokenPath, set = os.LookupEnv("GOOGLE_TOKEN_PATH")
	if !set {
		log.Fatal("Unable to load environment value for: GOOGLE_TOKEN_PATH")
	}

	rootCmd.AddCommand(detailsCmd)

	rootCmd.PersistentFlags().IntVarP(&rate, "rate", "r", 34, "Rate for the dog walking to apply to the event count.")
	rootCmd.PersistentFlags().IntVarP(&year, "year", "y", time.Now().Year(), "Year to pull the calendar month from.")
	rootCmd.PersistentFlags().IntVarP(&month, "month", "m", int(time.Now().Month()), "Month to pull the calendar events from.")
	rootCmd.PersistentFlags().StringVar(&eventPrefix, "prefix", `[DOG WALK]`, "Prefix of the events to pull from the calendar for counting.")
}
