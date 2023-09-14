package monitor

import (
	"bytes"
	"fmt"
	"testing"
	"time"
)

func TestTrackInternetOutage(t *testing.T) {

	entryJsonFormat := "{\"date\":\"%v\",\"downTime\":\"%v\",\"restoreTime\":\"%v\",\"durationInSeconds\":%v}\n"

	testMap := []struct {
		name             string
		lastStateWasDown bool
		entry            *DownEntry
		validationFunc   func() bool
		eventTime        time.Time
		wants            string
		wantEntry        DownEntry
	}{
		{
			name:             "Test Internet Down - Last Was Up",
			lastStateWasDown: false,
			entry:            &DownEntry{},
			validationFunc:   func() bool { return false },
			eventTime:        time.Date(2023, time.September, 12, 15, 34, 0, 0, time.UTC),
			wants:            "",
			wantEntry: DownEntry{
				Date:              "09/12/2023",
				DownTime:          time.Date(2023, time.September, 12, 15, 34, 0, 0, time.UTC),
				RestoreTime:       time.Time{},
				DurationInSeconds: 0,
			},
		},
		{
			name:             "Test Internet Down - Last Was Down",
			lastStateWasDown: true,
			entry: &DownEntry{
				Date:              "09/12/2023",
				DownTime:          time.Date(2023, time.September, 12, 15, 34, 0, 0, time.UTC),
				RestoreTime:       time.Time{},
				DurationInSeconds: 0,
			},
			validationFunc: func() bool { return false },
			eventTime:      time.Date(2023, time.September, 13, 16, 34, 0, 0, time.UTC),
			wants:          "",
			wantEntry: DownEntry{
				Date:              "09/12/2023",
				DownTime:          time.Date(2023, time.September, 12, 15, 34, 0, 0, time.UTC),
				RestoreTime:       time.Time{},
				DurationInSeconds: 0,
			},
		},
		{
			name:             "Test Internet Up - Last Was Down",
			lastStateWasDown: true,
			entry: &DownEntry{
				Date:              "09/12/2023",
				DownTime:          time.Date(2023, time.September, 12, 15, 34, 0, 0, time.UTC),
				RestoreTime:       time.Time{},
				DurationInSeconds: 0,
			},
			validationFunc: func() bool { return true },
			eventTime:      time.Date(2023, time.September, 12, 16, 34, 0, 0, time.UTC),
			wants: fmt.Sprintf(
				entryJsonFormat,
				"09/12/2023",
				time.Date(2023, time.September, 12, 15, 34, 0, 0, time.UTC).Format(time.RFC3339),
				time.Date(2023, time.September, 12, 16, 34, 0, 0, time.UTC).Format(time.RFC3339),
				3600,
			),
			wantEntry: DownEntry{},
		},
		{
			name:             "Test Internet Up - Last Was Up",
			lastStateWasDown: false,
			entry:            &DownEntry{},
			validationFunc:   func() bool { return true },
			eventTime:        time.Now(),
			wants:            "",
			wantEntry:        DownEntry{},
		},
	}

	for _, tt := range testMap {
		t.Run(tt.name, func(t *testing.T) {

			var buffer bytes.Buffer

			clock := MockTime{
				time: tt.eventTime,
			}

			trackInternetOutage(&tt.lastStateWasDown, tt.entry, &buffer, tt.validationFunc, &clock, 0)

			// Validate buffer
			if buffer.String() != tt.wants {
				t.Errorf("got: %q; want: %q", buffer.String(), tt.wants)
			}

			// Validate Entry
			if !tt.entry.Equal(tt.wantEntry) {
				t.Errorf("got: %q; want: %q", tt.entry, tt.wantEntry)
			}
		})
	}
}
