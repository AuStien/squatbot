package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/austien/squatbot/squat"
	log "github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

func main() {
	apiKey, ok := os.LookupEnv("SQUATBOT_API_KEY")
	if !ok {
		log.Fatal("missing API key")
	}

	publishID, ok := os.LookupEnv("SQUATBOT_PUBLISH_ID")
	if !ok {
		log.Fatal("missing ID to post message")
	}

	var maxWaitDuration time.Duration
	maxWaitDurationStr, ok := os.LookupEnv("SQUATBOT_MAX_WAIT_DURATION")
	if !ok {
		maxWaitDuration = time.Hour * 8 * 1000
	} else {
		duration, err := time.ParseDuration(maxWaitDurationStr)
		if err != nil {
			log.WithError(err).Fatalf("failed to parse %q as time.Duration. Make sure duration has specified unit. Valid time units are \"ns\", \"us\" (or \"µs\"), \"ms\", \"s\", \"m\", \"h\".", maxWaitDurationStr)
		}
		maxWaitDuration = duration
	}

	reportID := os.Getenv("SQUATBOT_REPORT_ID")

	go func() {
		r := http.NewServeMux()
		r.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) { return })
		r.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) { return })

		log.Info("Server running on port 8080")
		log.Fatal(http.ListenAndServe(":8080", r))
	}()

	api := slack.New(apiKey)

	squatMessages := make([]string, len(squat.SquatMessages))
	copy(squatMessages, squat.SquatMessages)
	restMessages := make([]string, len(squat.RestMessages))
	copy(restMessages, squat.RestMessages)

	rand.Seed(time.Now().Unix())
	for {
		wait := rand.Int63n(int64(maxWaitDuration))
		restWait := (time.Hour * 24) - time.Duration(wait)
		log.Infof("Sleeping until ", time.Unix(0, time.Now().UnixNano()+wait))
		time.Sleep(time.Duration(wait))

		squats, err := squat.Squats(time.Now().Day())
		if err != nil {
			if reportID != "" {
				log.WithError(err).Error("failed to get amount of squats")
				channel, _, postErr := api.PostMessage(reportID, slack.MsgOptionText(err.Error(), false))
				if postErr != nil {
					log.WithError(postErr).Errorf("failed to send error message to %q", channel)
				}
			} else {
				log.WithError(err).Fatal("failed to get amount of squats")
			}
			time.Sleep(restWait)
			continue
		}

		var msg string
		if squats == 0 {
			i := rand.Intn(len(restMessages))
			msg = restMessages[i]
			restMessages = append(restMessages[:i], restMessages[i+1:]...)
		} else {
			i := rand.Intn(len(squatMessages))
			msg = fmt.Sprintf(squatMessages[i], squats)
			squatMessages = append(squatMessages[:i], squatMessages[i+1:]...)
		}

		channel, _, err := api.PostMessage(publishID, slack.MsgOptionText(msg, false))
		if err != nil {
			if reportID != "" {
				log.WithError(err).Errorf("failed to send message: %q", msg)
				channel, _, postErr := api.PostMessage(reportID, slack.MsgOptionText(err.Error(), false))
				if postErr != nil {
					log.WithError(postErr).Errorf("failed to send error message to %q", channel)
				}
			} else {
				log.WithError(err).Fatalf("failed to send message: %q", msg)
			}
			time.Sleep(restWait)
			continue
		}

		log.WithField("squatMessages", squatMessages).Infof("(%s): %s", channel, msg)

		if len(squatMessages) == 0 {
			tmp := make([]string, len(squat.SquatMessages))
			copy(tmp, squat.SquatMessages)
			squatMessages = tmp
			log.Info("Squat messages have been reset")
		}

		if len(restMessages) == 0 {
			tmp := make([]string, len(squat.RestMessages))
			copy(tmp, squat.RestMessages)
			restMessages = tmp
			log.Info("Rest messages have been reset")
		}

		// This makes sure that the random duration always starts at the same time of day
		log.Infof("Starting next iteration at %s", time.Unix(0, time.Now().UnixNano()+restWait.Nanoseconds()))
		time.Sleep(restWait)
	}
}
