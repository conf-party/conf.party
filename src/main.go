package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"sort"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

var (
	rootDir string
	outDir  string
	action  string = "build"
)

const dateLayout = "2006-01-02"

func init() {
	defaultDir, _ := os.Getwd()

	flag.StringVar(&action, "action", "build", "")
	flag.StringVar(&rootDir, "rootdir", defaultDir, "")
	flag.StringVar(&outDir, "out", "/out", "")
	flag.Parse()
}

func main() {
	switch action {
	case "build":
		build()
	case "verify":
		verify()
	case "serve":
		build()
		serve()
	default:
		build()
	}
}

func build() {
	srcManifests, err := ioutil.ReadDir(path.Join(rootDir, "conferences"))
	if err != nil {
		panic(err)
	}

	templateFuncs := template.FuncMap{
		"parse": func(str string) template.HTML {
			return template.HTML(strings.ReplaceAll(str, "\n\n", "<br/><br/>"))
		},
		"minus": func(a, b int) int {
			return a - b
		},
		"day": func(date string) string {
			t, _ := time.Parse(dateLayout, date)
			return t.Format("Monday")
		},
		"prettyDate": func(date, endDate, dateTime string) string {
			t, _ := time.Parse(dateLayout, date)
			d := t.Format("Monday, 2 January")

			if endDate != "" {
				t, _ := time.Parse(dateLayout, endDate)
				d = fmt.Sprintf("%s → %s", d, t.Format("Monday, 2 January"))
			}

			if dateTime != "" {
				d = fmt.Sprintf("%s - %s", d, dateTime)
			}

			return d
		},
	}

	indexFile, _ := os.ReadFile(path.Join(rootDir, "./src/templates/index.html"))
	eventFile, _ := os.ReadFile(path.Join(rootDir, "./src/templates/event.html"))

	indexTmpl, _ := template.New("index").Funcs(templateFuncs).Parse(string(indexFile))
	eventTmpl, _ := template.New("event").Funcs(templateFuncs).Parse(string(eventFile))

	confs := []Conference{}

	for _, file := range srcManifests {
		if !file.IsDir() {
			data, err := ioutil.ReadFile(path.Join(rootDir, "conferences", file.Name()))
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

			if conf.EndDate > time.Now().Add(25*time.Hour).Format("2006-01-02") {
				confs = append(confs, conf)
			}
		}
	}

	sort.Slice(confs, func(i, j int) bool {
		return confs[i].Date > confs[j].Date
	})

	for _, conf := range confs {
		outFile, err := os.Create(path.Join(outDir, conf.Filename))
		if err != nil {
			log.Println("create file: ", err)
			return
		}

		if err := eventTmpl.Execute(outFile, conf); err != nil {
			log.Println("template event: ", err)
		}
	}

	outFile, err := os.Create(path.Join(outDir, "index.html"))
	if err != nil {
		log.Println("create file: ", err)
		return
	}

	if err := indexTmpl.Execute(outFile, confs); err != nil {
		log.Println("template event: ", err)
	}

	// // Copy static files
	cp := exec.Command("sh", "-c", fmt.Sprintf("cp -r %s/* %s", path.Join(rootDir, "out"), outDir))
	err = cp.Run()
	if err != nil {
		panic(err)
	}
}

func verify() {
	srcManifests, err := ioutil.ReadDir(path.Join(rootDir, "conferences"))
	if err != nil {
		panic(err)
	}
	for _, file := range srcManifests {
		if !file.IsDir() {
			data, err := ioutil.ReadFile(path.Join(rootDir, "conferences", file.Name()))
			if err != nil {
				panic(err)
			}

			fmt.Printf("Validating %s\n", file.Name())

			var conf Conference
			err = yaml.Unmarshal(data, &conf)
			if err != nil {
				panic(err)
			}

			if err := conf.Validate(); err != nil {
				panic(err)
			}
		}
	}
}

func serve() {
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir(outDir))))
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
