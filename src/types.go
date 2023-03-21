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
	Name     string  `yaml:"name"`
	Date     string  `yaml:"date"`
	EndDate  string  `yaml:"endDate"`
	Website  string  `yaml:"website"`
	Location string  `yaml:"location"`
	Parties  []Party `yaml:"parties"`
}

type Party struct {
	Name        string `yaml:"name"`
	Date        string `yaml:"date"`
	Time        string `yaml:"time"`
	Website     string `yaml:"website"`
	Location    string `yaml:"location"`
	Description string `yaml:"description"`
	Notes       string `yaml:"notes"`
}

const dateLayout = "2006-01-02"

func (c Conference) PrettyDate() string {
	t, _ := time.Parse(dateLayout, c.Date)
	d := t.Format("Monday, 2 January")

	if c.EndDate != "" {
		t, _ := time.Parse(dateLayout, c.EndDate)
		d = fmt.Sprintf("%s → %s", d, t.Format("Monday, 2 January"))
	}

	return d
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
