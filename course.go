package main

type Course struct {
	ID   int    `json:"id"`
	Name string `json:"name,omitempty"`
}

type Assignment struct {
	Due_date  string `json:"due_at"`
	Name      string `json:"name"`
	course    string
	courseKey int
}
