package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

		// grab all assignments for course
		for j := range assignmentsForCourse {

			// create a time object out of the time string
			dueDate, err := time.Parse(time.RFC3339, assignmentsForCourse[j].Due_date)

			// convert to local time
			// localizing is being handled in getAssignmentsForCourse currently
			// dueDate = dueDate.Local()
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

func instructureRequest(requestUrl string, authorizationToken string) *http.Response {

	// craft reqeust
	req, err := http.NewRequest("GET", requestUrl, nil)
	req.Header.Add("Authorization", authorizationToken)
	if err != nil {
		fmt.Printf(err.Error())
	}

	// send request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf(err.Error())
	}

	// I may want to handle possible errors here
	return resp
}

// eureka!
// should return: array of all course IDs
func getCourses(auth string) []Course {
	exampleCoursesUrl := "https://capital.instructure.com/api/v1/courses/"

	resp := instructureRequest(exampleCoursesUrl, auth)
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf(err.Error())
		}

		// snippet only
		// unmarshal into slice of Course struct
		var result []Course
		if err := json.Unmarshal(bodyBytes, &result); err != nil { // Parse []byte to go struct pointer
			fmt.Println("Can not unmarshal JSON")
		}

		return result
	}
	return nil
}

func getAssignmentsForCourse(course Course, auth string) []Assignment {

	assignmentsUrlTemplate := fmt.Sprintf("https://capital.instructure.com/api/v1/courses/%d/assignments", course.ID)

	resp := instructureRequest(assignmentsUrlTemplate, auth)
	defer resp.Body.Close()

	// handle the response
	if resp.StatusCode == http.StatusOK {
		// pass response body to io reader
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf(err.Error())
		}

		// unmarshal JSON into slice of custom struct
		var result []Assignment
		if err := json.Unmarshal(bodyBytes, &result); err != nil {
			fmt.Println("Can not unmarshal JSON")
		}

		for i := range result {
			dueDate, err := time.Parse(time.RFC3339, result[i].Due_date)
			dueDate = dueDate.Local()
			if err == nil {
				current := time.Now().Local()
				if dueDate.After(current) {
					result[i].course = course.Name
					result[i].courseKey = course.ID
					result[i].Due_date = dueDate.Format(time.RFC3339)
				} else {
					result[i].course = ""
					result[i].Due_date = ""
					result[i].Name = ""
					result[i].courseKey = 0
				}
			}
		}
		return result
	} else {
		fmt.Printf(resp.Status)
	}
	return nil
}
