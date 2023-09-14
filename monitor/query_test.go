package monitor

import (
	"fmt"
	"io"
	"strings"
	"testing"
	"time"
)

func TestQueryData(t *testing.T) {

	entryJsonFormat := "{\"date\":\"%v\",\"downTime\":\"%v\",\"restoreTime\":\"%v\",\"durationInSeconds\":%v}\n"

	io.MultiReader()

	testMap := []struct {
		name          string
		reader        io.Reader
		actualTime    time.Time
		wantsResponse QueryResponse
	}{
		{
			name: "Returns correct query response",
			reader: io.MultiReader(
				strings.NewReader(
					fmt.Sprintf(
						entryJsonFormat,
						"09/12/2023",
						time.Date(2023, time.September, 12, 15, 34, 0, 0, time.UTC).Format(time.RFC3339),
						time.Date(2023, time.September, 12, 16, 34, 0, 0, time.UTC).Format(time.RFC3339),
						3600,
					),
				),
				strings.NewReader(
					fmt.Sprintf(
						entryJsonFormat,
						"09/06/2023",
						time.Date(2023, time.September, 6, 15, 34, 0, 0, time.UTC).Format(time.RFC3339),
						time.Date(2023, time.September, 6, 16, 34, 0, 0, time.UTC).Format(time.RFC3339),
						3600,
					),
				),
				strings.NewReader(
					fmt.Sprintf(
						entryJsonFormat,
						"08/13/2023",
						time.Date(2023, time.August, 13, 15, 34, 0, 0, time.UTC).Format(time.RFC3339),
						time.Date(2023, time.August, 13, 16, 34, 0, 0, time.UTC).Format(time.RFC3339),
						3600,
					),
				),
			),
			actualTime: time.Date(2023, time.September, 12, 20, 34, 0, 0, time.UTC),
			wantsResponse: QueryResponse{
				Today: QueryRow{
					Avg:    3600,
					Amount: 1,
				},
				Days7: QueryRow{
					Avg:    3600,
					Amount: 2,
				},
				Days30: QueryRow{
					Avg:    3600,
					Amount: 3,
				},
			},
		},
	}

	for _, tt := range testMap {
		t.Run(tt.name, func(t *testing.T) {

			clock := &MockTime{
				time: tt.actualTime,
			}

			res, err := QueryData(tt.reader, clock)
			if err != nil {
				t.Fatal(err)
			}

			if !res.Equal(tt.wantsResponse) {
				t.Errorf("got: %v; want: %v", res, tt.wantsResponse)
			}
		})
	}
}
