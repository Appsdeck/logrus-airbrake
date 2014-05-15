package logrus_airbrake

import (
	"errors"
	"fmt"
	"net/http"

	airbrake "github.com/AlekSi/airbrake-go"
	"github.com/Sirupsen/logrus"
)

type Hook struct{}

func (hook Hook) Fire(entry *logrus.Entry) error {
	var err error
	var req *http.Request

	if r, ok := entry.Data["req"]; ok {
		req, ok = r.(*http.Request)
		if ok {
			// We don't want to log credentials
			req.Header.Del("Authorization")

			entry.Data["req"] = fmt.Sprintf(
				"%s %s %s %s",
				req.Method, req.URL, req.UserAgent(), req.RemoteAddr,
			)
		}
	} else {
		// If there is no request, we build one in order to send
		// all the variables to airbrake
		req = new(http.Request)
		req.Header = make(http.Header)
	}

	// All the fields which aren't level|msg|error|time|req are added
	// to the headers of the request which will be sent to Airbrake
	// The main goal is to be able to see all the values on Airbrake dashboard
	for val, key := range entry.Data {
		if val != "level" && val != "msg" && val != "error" && val != "time" && val != "req" {
			req.Header.Add("log-"+val, fmt.Sprintf("%v", key))
		}
	}

	// If there is an error field, we want it to be part of Airbrake ticket name
	var errorMsg error
	if entry.Data["error"] != nil {
		errorMsg = fmt.Errorf("%v - %v",
			entry.Data["error"].(error),
			entry.Data["msg"].(string),
		)
	} else {
		errorMsg = errors.New(entry.Data["msg"].(string))
	}
	err = airbrake.Error(errorMsg, req)

	if err != nil {
		log := logrus.New()
		log.WithFields(logrus.Fields{
			"source":   "airbrake",
			"endpoint": airbrake.Endpoint,
			"error":    err,
		}).Warn("Failed to send error to Airbrake")
	}

	return nil
}

func (hook Hook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.Error,
		logrus.Fatal,
		logrus.Panic,
	}
}
