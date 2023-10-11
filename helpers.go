package main

import "time"

func sortDates(assignments []Assignment) {

	for i := range assignments {
		for j := 0; j < len(assignments)-i-1; j++ {
			date, _ := time.Parse(time.RFC3339, assignments[j].Due_date)
			nextDate, _ := time.Parse(time.RFC3339, assignments[j+1].Due_date)
			if date.After(nextDate) {
				assignments[j], assignments[j+1] = assignments[j+1], assignments[j]
			}
		}
	}
}
