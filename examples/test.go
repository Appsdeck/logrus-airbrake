package main

import (
	"fmt"
	"math/rand"
	"net/http"

	airbrake "github.com/AlekSi/airbrake-go"
	logrus_airbrake "github.com/Appsdeck/logrus-airbrake"
	"github.com/Sirupsen/logrus"
)

func main() {
	airbrake.ApiKey = "123456ABCD"
	airbrake.Environment = "testing"

	logger := logrus.New()
	logger.Hooks.Add(logrus_airbrake.Hook{})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := fmt.Errorf("something wrong happened in the database")

		logger.WithFields(
			logrus.Fields{"req": r, "error": err, "extra-data": rand.Int()},
		).Error("Something is really wrong")
	})

	http.ListenAndServe(":31313", nil)
}
