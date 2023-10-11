package main

import (
	"fmt"
	"os"
	"time"
)

func main() {

	// grab auth token from command line
	token := os.Args[1]
	auth := fmt.Sprintf("Bearer %s", token)

	coursesSlice := getCourses(auth)

	// iterate through each course
	for i := range coursesSlice {

		assignmentsForCourse := getAssignmentsForCourse(coursesSlice[i], auth)

		// print the course name
		fmt.Printf("%s\n\n", coursesSlice[i].Name)
		sortDates(assignmentsForCourse)
		// grab all assignments for course
		for j := range assignmentsForCourse {

			// create a time object out of the time string
			dueDate, err := time.Parse(time.RFC3339, assignmentsForCourse[j].Due_date)

			// convert to local time
			if err == nil {
				formatted := fmt.Sprintf("%s %d, %d %d:%d", dueDate.Month().String(), dueDate.Day(), dueDate.Year(), dueDate.Hour(), dueDate.Minute())

				// print assignment name and due date
				fmt.Printf("Assignment: %s, Due On: %s\n", assignmentsForCourse[j].Name, formatted)
			}
		}
		// format spacing
		fmt.Printf("\n")
	}

}
