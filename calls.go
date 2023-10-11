package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

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
