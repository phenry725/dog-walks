package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/jinzhu/now"
	"github.com/joho/godotenv"
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

var dogWalkingRate int

var rootCmd = &cobra.Command{
	Use:   "walks",
	Short: "walks - a simple CLI to calculate monthly payment for dog walks based on calendar events",
	Long: `walks is a CLI built to pull dog walk events from google calendar
and calcualte the monthly due for those walks.`,
	Run: func(cmd *cobra.Command, args []string) {
		b, err := ioutil.ReadFile(googleTokenPath)
		if err != nil {
			osExitErr(fmt.Sprintf("Unable to read credentials file from path: '%s'", err))
		}
		conf, err := google.JWTConfigFromJSON(b, calendar.CalendarReadonlyScope)
		if err != nil {
			osExitErr(fmt.Sprintf("Unable to load the JWT config from credentials json in environment: '%s'", err))
		}

		client := conf.Client(oauth2.NoContext)
		srv, err := calendar.NewService(oauth2.NoContext, option.WithHTTPClient(client))
		if err != nil {
			osExitErr(fmt.Sprintf("Unable to retrieve Calendar client: '%s'", err))
		}

		fmt.Fprintf(os.Stdout, "Looking up dog walks from %v to %v\n", now.BeginningOfMonth().Format(time.RFC3339), now.EndOfMonth().Format(time.RFC3339))

		events, err := srv.Events.List("primary").Q(`[DOG WALK]`).SingleEvents(true).ShowDeleted(false).TimeMin(now.BeginningOfMonth().Format(time.RFC3339)).TimeMax(now.EndOfMonth().Format(time.RFC3339)).Do()
		if err != nil {
			osExitErr(fmt.Sprintf("Unable to retrieve list of events from calendar: '%s'", err))
		}
		fmt.Fprintf(os.Stdout, "%v walks will cost: %v\n", time.Now().Month().String(), len(events.Items)*dogWalkingRate)

	},
}

func init() {
	var set bool

	envFileNameArg := "setupDev.sh"
	err := godotenv.Load(envFileNameArg)
	if err != nil {
		log.Fatalf("Error loading env file named: %v", envFileNameArg)
	}

	googleTokenPath, set = os.LookupEnv("GOOGLE_TOKEN_PATH")
	if !set {
		log.Fatal("Unable to load environment value for: GOOGLE_TOKEN_PATH")
	}
	rootCmd.LocalFlags().IntVarP(&dogWalkingRate, "rate", "r", 34, "Rate for the dog walking to apply to the event count.")
}
