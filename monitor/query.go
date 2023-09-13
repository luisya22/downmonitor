package monitor

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"time"
)

func QueryData(r io.Reader) error {
	entries := getData(r)

	dailyDurations := make(map[string]int)
	dailyCounts := make(map[string]int)

	today := time.Now()
	sevenDaysAgo := today.AddDate(0, 0, -7)
	thirtyDaysAgo := today.AddDate(0, -1, 0)

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

		if entryDate.After(sevenDaysAgo) {
			last7DaysDurations += entry.DurationInSeconds
			last7DaysCounts++
			uniqueDaysInLast7[entry.Date] = struct{}{}
		}

		if entryDate.After(thirtyDaysAgo) {
			last30DaysDurations += entry.DurationInSeconds
			last30DaysCounts++
			uniqueDaysInLast30[entry.Date] = struct{}{}
		}
	}

	todayEntryDurations, _ := dailyDurations[time.Now().Format("01/02/2006")]
	todayEntryCounts, _ := dailyCounts[time.Now().Format("01/02/2006")]

	avg := 0
	avg7 := 0
	avg30 := 0

	if todayEntryCounts > 0 {
		avg = todayEntryDurations / todayEntryDurations
	}

	if len(uniqueDaysInLast7) > 0 {
		avg7 = last7DaysDurations / len(uniqueDaysInLast7)
	}

	if len(uniqueDaysInLast30) > 0 {
		avg30 = last30DaysDurations / len(uniqueDaysInLast30)
	}

	//TODO: return data instead of printing

	fmt.Printf("Today: Average downtime: %d seconds, Downtimes: %d\n", avg, todayEntryCounts)
	fmt.Printf("Last 7 days: Average daily downtime: %d seconds, Downtimes: %d\n", avg7, last7DaysCounts)
	fmt.Printf("Last 30 days: Average daily downtime: %d seconds, Downtimes: %d\n", avg30, last30DaysCounts)
	fmt.Printf("---------\n\n----------\n")

	return nil
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
