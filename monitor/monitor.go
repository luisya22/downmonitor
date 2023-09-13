package monitor

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"time"
)

type DownEntry struct {
	Date              string    `json:"date"`
	DownTime          time.Time `json:"downTime"`
	RestoreTime       time.Time `json:"restoreTime"`
	DurationInSeconds int       `json:"durationInSeconds"`
}

func (de *DownEntry) resetValues() {
	de.Date = ""
	de.DownTime = time.Time{}
	de.RestoreTime = time.Time{}
	de.DurationInSeconds = 0
}

func (de *DownEntry) Equal(c DownEntry) bool {
	if de.Date != c.Date {
		return false
	}

	if !de.DownTime.Equal(c.DownTime) {
		return false
	}

	if !de.RestoreTime.Equal(c.RestoreTime) {
		return false
	}

	if de.DurationInSeconds != c.DurationInSeconds {
		return false
	}

	return true
}

func Start() {
	lastStateWasDown := false
	entry := &DownEntry{}
	validationFunc := googleRequest

	file, err := os.OpenFile("./downmonitor.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("error opening file: %s", err.Error())
	}
	defer file.Close()

	c := RealTime{}

	for {
		trackInternetOutage(&lastStateWasDown, entry, file, validationFunc, &c, 10)
	}
}

func trackInternetOutage(lastStateWasDown *bool, entry *DownEntry, w io.Writer, isInternetAvailable func() bool, clock Clock, waitSeconds int) {
	isDown := !isInternetAvailable()

	if isDown && !*lastStateWasDown {
		entry.Date = clock.Now().Format("01/02/2006")
		entry.DownTime = clock.Now()
		*lastStateWasDown = true
	} else if !isDown && *lastStateWasDown {
		entry.RestoreTime = clock.Now()
		duration := entry.RestoreTime.Sub(entry.DownTime)
		entry.DurationInSeconds = int(math.Floor(duration.Seconds()))

		logDownTime(entry, w)

		entry.resetValues()

		*lastStateWasDown = false
	}

	time.Sleep(time.Duration(waitSeconds) * time.Second)

}

func googleRequest() bool {

	client := &http.Client{
		Timeout: time.Second * 5,
	}

	res, err := client.Get("https://www.google.com/")
	if err != nil {
		return false
	}
	res.Body.Close()

	return res.StatusCode == 200
}

func logDownTime(entry *DownEntry, w io.Writer) error {

	entryJson, err := json.Marshal(*entry)
	if err != nil {
		return err
	}

	out := fmt.Sprintln(string(entryJson))
	fmt.Fprintf(w, "%s", out)

	return nil
}
