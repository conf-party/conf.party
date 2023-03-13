package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
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
	Website     string
	Location    string
	Description string
	Notes       string
}

func main() {
	srcManifests, err := ioutil.ReadDir("../conferences")
	if err != nil {
		panic(err)
	}

	indexTmpl := template.Must(template.ParseFiles("./templates/index.html"))
	eventTmpl := template.Must(template.ParseFiles("./templates/event.html"))

	confs := []Conference{}

	for _, file := range srcManifests {
		if !file.IsDir() {
			data, err := ioutil.ReadFile(path.Join("../conferences", file.Name()))
			if err != nil {
				log.Println(err)
				return
			}

			var conf Conference

			err = yaml.Unmarshal(data, &conf)
			if err != nil {
				log.Println(err)
				return
			}
			conf.Filename = strings.ReplaceAll(file.Name(), ".yaml", ".html")

			confs = append(confs, conf)
		}
	}

	for _, conf := range confs {
		outFile, err := os.Create(path.Join("../docs", conf.Filename))
		if err != nil {
			log.Println("create file: ", err)
			return
		}

		if err := eventTmpl.Execute(outFile, conf); err != nil {
			log.Println("template event: ", err)
		}
	}

	outFile, err := os.Create(path.Join("../docs", "index.html"))
	if err != nil {
		log.Println("create file: ", err)
		return
	}

	if err := indexTmpl.Execute(outFile, confs); err != nil {
		log.Println("template event: ", err)
	}
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
