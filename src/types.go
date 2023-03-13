package main

import "time"

type Conference struct {
	Filename string
	Name     string
	Date     string
	Website  string
	Location string
	Parties  []Party
}

type Party struct {
	Name        string
	Date        string
	Website     string
	Location    string
	Description string
	Notes       string
}

const dateLayout = "2006-01-02"

func (c Conference) PrettyDate() string {
	t, _ := time.Parse(dateLayout, c.Date)
	return t.Format("Monday, 2 January")
}

func (p Party) PrettyDate() string {
	t, _ := time.Parse(dateLayout, p.Date)
	return t.Format("Monday, 2 January")
}
