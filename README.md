# Airbrake hook for [Logrus](https://github.com/Sirupsen/logrus)

## Setup

```sh
go get github.com/Appsdeck/logrus-airbrake
```

## Example

```go
package main

import (
	"fmt"
	"net/http"
	"math/rand"
	
	"github.com/Sirupsen/logrus"
	airbrake "github.com/AlekSi/airbrake-go"
	logrus_airbrake "github.com/Appsdeck/logrus-airbrake"
)


func main() {
	airbrake.ApiKey = "123456ABCD"
	airbrake.Environment = "testing"

	logger := logrus.New()
	logger.Hooks.Add(logrus_airbrake.Hook{})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := fmt.Errorf("something wrong happened in the database")

		logger.WithFields(
			logrus.Fields{"req": r, "error": err, "extra-data", rand.Int()},
		).Error("Something is really wrong")
	})

	http.ListenAndServe(":31313", nil)
}
```
