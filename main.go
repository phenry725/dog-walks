package main

import (
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"

	"github.com/jinzhu/now"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

var (
	googleTokenPath  string
	googleCalendarId string

	initOnce sync.Once
)

const (
	dogWalkingRate = 34
)

func init() {
	var set bool
	var err error
	initOnce.Do(func() {
		envFileNameArg := "setupDev.sh"
		err = godotenv.Load(envFileNameArg)
		if err != nil {
			log.Fatalf("Error loading env file named: %v", envFileNameArg)
		}

		googleTokenPath, set = os.LookupEnv("GOOGLE_TOKEN_PATH")
		if !set {
			log.Fatal("Unable to load environment value for: GOOGLE_TOKEN_PATH")
		}
	})
}

func main() {
	// Error checks are discarded for brevity.
	b, err := ioutil.ReadFile(googleTokenPath)
	if err != nil {
		panic(err)
	}
	conf, err := google.JWTConfigFromJSON(b, calendar.CalendarReadonlyScope)
	if err != nil {
		panic(err)
	}

	// Create the client.
	client := conf.Client(oauth2.NoContext)
	srv, err := calendar.NewService(oauth2.NoContext, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}

	log.Printf("Looking up dog walks from %v to %v", now.BeginningOfMonth().Format(time.RFC3339), now.EndOfMonth().Format(time.RFC3339))

	events, err := srv.Events.List("primary").Q(`[DOG WALK]`).SingleEvents(true).ShowDeleted(false).TimeMin(now.BeginningOfMonth().Format(time.RFC3339)).TimeMax(now.EndOfMonth().Format(time.RFC3339)).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve next list of the user's events: %v", err)
	}
	log.Printf("%v walks will cost: %v", time.Now().Month().String(), len(events.Items)*dogWalkingRate)
}
