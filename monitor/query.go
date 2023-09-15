package monitor

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"time"
)

type QueryRow struct {
	Avg      int
	Amount   int
	TotalSum int
}

type QueryResponse struct {
	Today  QueryRow
	Days7  QueryRow
	Days30 QueryRow
}

func (qr *QueryRow) Equal(c QueryRow) bool {

	if qr.Avg != c.Avg {
		return false
	}

	if qr.Amount != c.Amount {
		return false
	}

	if qr.TotalSum != c.TotalSum {
		return false
	}

	return true
}

func (qr *QueryResponse) Equal(c QueryResponse) bool {

	if !qr.Today.Equal(c.Today) {
		return false
	}

	if !qr.Days7.Equal(c.Days7) {
		return false
	}

	if !qr.Days30.Equal(c.Days30) {
		return false
	}

	return true
}

func QueryData(r io.Reader, clock Clock) (QueryResponse, error) {
	entries := getData(r)

	dailyDurations := make(map[string]int)
	dailyCounts := make(map[string]int)

	now := clock.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	sevenDaysAgo := today.AddDate(0, 0, -7)
	thirtyDaysAgo := today.AddDate(0, 0, -30)

	var todayTotalSum int
	var last7DaysDurations int
	var last7DaysCounts int
	var last30DaysDurations int
	var last30DaysCounts int
	uniqueDaysInLast7 := make(map[string]struct{})
	uniqueDaysInLast30 := make(map[string]struct{})

	for _, entry := range entries {
		entryDate, _ := time.Parse("01/02/2006", entry.Date)

		dailyDurations[entry.Date] += entry.DurationInSeconds
		dailyCounts[entry.Date]++

		if entry.Date == clock.Now().Format("01/02/2006") {
			todayTotalSum += entry.DurationInSeconds
		}

		if entryDate.After(sevenDaysAgo) || entryDate.Equal(sevenDaysAgo) {
			last7DaysDurations += entry.DurationInSeconds
			last7DaysCounts++
			uniqueDaysInLast7[entry.Date] = struct{}{}
		}

		if entryDate.After(thirtyDaysAgo) || entryDate.Equal(thirtyDaysAgo) {
			last30DaysDurations += entry.DurationInSeconds
			last30DaysCounts++
			uniqueDaysInLast30[entry.Date] = struct{}{}
		}
	}

	todayEntryDurations, _ := dailyDurations[clock.Now().Format("01/02/2006")]
	todayEntryCounts, _ := dailyCounts[clock.Now().Format("01/02/2006")]

	avg := 0
	avg7 := 0
	avg30 := 0

	if todayEntryCounts > 0 {
		avg = todayEntryDurations / todayEntryCounts
	}

	if len(uniqueDaysInLast7) > 0 {
		avg7 = last7DaysDurations / len(uniqueDaysInLast7)
	}

	if len(uniqueDaysInLast30) > 0 {
		avg30 = last30DaysDurations / len(uniqueDaysInLast30)
	}

	t := QueryRow{
		Avg:      avg,
		Amount:   todayEntryCounts,
		TotalSum: todayTotalSum,
	}

	days7 := QueryRow{
		Avg:      avg7,
		Amount:   last7DaysCounts,
		TotalSum: last7DaysDurations,
	}

	days30 := QueryRow{
		Avg:      avg30,
		Amount:   last30DaysCounts,
		TotalSum: last30DaysDurations,
	}

	res := QueryResponse{
		Today:  t,
		Days7:  days7,
		Days30: days30,
	}

	return res, nil
}

func getData(r io.Reader) []DownEntry {
	scanner := bufio.NewScanner(r)
	entries := []DownEntry{}

	for scanner.Scan() {
		line := scanner.Bytes()

		if len(line) == 0 {
			continue
		}

		entry := DownEntry{}
		err := json.Unmarshal(line, &entry)
		if err != nil {
			fmt.Printf("error processing line: %s. Error: %v\n", string(line), err.Error())
		}

		entries = append(entries, entry)
	}

	return entries
}
