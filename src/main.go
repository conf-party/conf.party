package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"gopkg.in/yaml.v3"
)

var (
	rootDir string
	outDir  string
	action  string = "build"
)

func init() {
	defaultDir, _ := os.Getwd()

	flag.StringVar(&rootDir, "rootdir", defaultDir, "")
	flag.StringVar(&outDir, "out", "/out", "")
	flag.Parse()

	args := os.Args[1:]
	if len(args) > 0 {
		action = args[0]
	}
}

func main() {
	switch action {
	case "build":
		build()
	case "verify":
		verify()
	default:
		build()
	}
}

func build() {
	srcManifests, err := ioutil.ReadDir(path.Join(rootDir, "conferences"))
	if err != nil {
		panic(err)
	}

	indexTmpl := template.Must(template.ParseFiles(path.Join(rootDir, "./src/templates/index.html")))
	eventTmpl := template.Must(template.ParseFiles(path.Join(rootDir, "./src/templates/event.html")))

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

			confs = append(confs, conf)
		}
	}

	for _, conf := range confs {
		outFile, err := os.Create(path.Join(rootDir, "out", conf.Filename))
		if err != nil {
			log.Println("create file: ", err)
			return
		}

		if err := eventTmpl.Execute(outFile, conf); err != nil {
			log.Println("template event: ", err)
		}
	}

	outFile, err := os.Create(path.Join(rootDir, "out", "index.html"))
	if err != nil {
		log.Println("create file: ", err)
		return
	}

	if err := indexTmpl.Execute(outFile, confs); err != nil {
		log.Println("template event: ", err)
	}

	// Copy files
	path.Join(rootDir, "out")

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
