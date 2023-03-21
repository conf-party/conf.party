package main

import (
	"fmt"
	"html/template"
	"net/url"
	"strings"
	"time"
)

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
	Time        string
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

func (c Conference) Validate() error {
	var err error

	if c.Name == "" {
		return fmt.Errorf("name must be provided")
	}
	if c.Website == "" {
		return fmt.Errorf("website must be provided")
	}
	if c.Date == "" {
		return fmt.Errorf("date must be provided")
	}

	_, err = time.Parse(dateLayout, c.Date)
	if err != nil {
		return fmt.Errorf("date is not in the expected format: %v", err)
	}

	_, err = url.ParseRequestURI(c.Website)
	if err != nil {
		return fmt.Errorf("failed to parse website URL: %v", err)
	}

	for _, p := range c.Parties {
		if err := p.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func (p Party) PrettyDate() string {
	t, _ := time.Parse(dateLayout, p.Date)
	d := t.Format("Monday, 2 January")
	if p.Time != "" {
		d = fmt.Sprintf("%s - %s", d, p.Time)
	}
	return d
}

func (p Party) ParsedDescription() template.HTML {
	return template.HTML(strings.ReplaceAll(p.Description, "\n\n", "<br/><br/>"))
}

func (p Party) ParsedNotes() template.HTML {
	return template.HTML(strings.ReplaceAll(p.Notes, "\n\n", "<br/><br/>"))
}

func (p Party) Validate() error {
	var err error

	if p.Name == "" {
		return fmt.Errorf("name must be provided")
	}
	if p.Website == "" {
		return fmt.Errorf("website must be provided")
	}
	if p.Date == "" {
		return fmt.Errorf("date must be provided")
	}

	_, err = time.Parse(dateLayout, p.Date)
	if err != nil {
		return fmt.Errorf("date is not in the expected format: %v", err)
	}

	_, err = url.ParseRequestURI(p.Website)
	if err != nil {
		return fmt.Errorf("failed to parse website URL: %v", err)
	}

	return nil
}
