package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path"
	"strings"
)

type Data struct {
	Host, Path, Method, Proto, RemoteIP string
	Header                              http.Header
	SiteConfig, SiteData                map[string]interface{}
}

func handle404(w http.ResponseWriter, r *http.Request, data *Data) {
	w.WriteHeader(http.StatusNotFound)
	custom404 := path.Join(data.Host, "404.html")
	if validateFile(custom404) {
		// server custom404
		log.Printf("404 page not found!")
		t, err := template.ParseFiles(custom404)
		if err != nil {
			log.Printf("Could not parse template file: %s\n", err)
			fmt.Fprint(w, "404 page not found!")
		} else {
			t.Execute(w, data)
		}
	} else {
		log.Printf("404 page not found!")
		fmt.Fprint(w, "404 page not found!")
	}
}

func mainHandler(w http.ResponseWriter, r *http.Request) {

	data := &Data{
		Host:     host(r.Host),
		Path:     r.URL.Path,
		Method:   r.Method,
		Proto:    r.Proto,
		RemoteIP: strings.Split(r.RemoteAddr, ":")[0],
		Header:   r.Header,
	}

	hp, err := hostPath(r.Host)
	if err != nil {
		//404
		handle404(w, r, data)
	} else {
		log.Println("Loading templates...")
		templates, err := loadTemplates(hp)
		if err != nil {
			log.Println("Could not load templates.")
		} else {

			err := loadJSONFile(path.Join(hp, "config.json"), &data.SiteConfig)
			if err != nil {
				log.Printf("No site configuration: %s\n", err)
			}

			err = loadJSONFile(path.Join(hp, "data.json"), &data.SiteData)
			if err != nil {
				log.Printf("No site data: %s\n", err)
			}

			if len(templates) > 0 {
				log.Printf("Parsing %d templates.", len(templates))
				t, err := template.ParseFiles(templates...)
				if err != nil {
					log.Printf("Could not parse required templates: %s\n", err)
					fmt.Fprint(w, "Could not parse required templates!")
				} else {
					t.ExecuteTemplate(w, "layout", data)
				}
			} else {
				log.Printf("The site at %s has no templates. Please create some templates.", data.Host)
				handle404(w, r, data)
			}
		}
	}
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
	hp, err := hostPath(r.Host)
	if err == nil {
		staticPath := path.Join(hp, "static")
		if validatePath(staticPath) {
			fs := http.StripPrefix("/static/", http.FileServer(http.Dir(staticPath)))
			fs.ServeHTTP(w, r)
		}
	}
}

func singleHandler(w http.ResponseWriter, r *http.Request) {
	hp, err := hostPath(r.Host)
	if err == nil {
		filePath := path.Join(hp, r.URL.Path)
		if validateFile(filePath) {
			http.ServeFile(w, r, filePath)
		}
	}
}

func main() {
	err := checkWorkDir()
	if err != nil {
		log.Fatalf("Could not find, or create the working directory at %s\n", config.workDir())
	}

	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/static/", staticHandler)
	http.HandleFunc("/sitemap.xml", singleHandler)
	http.HandleFunc("/favicon.ico", singleHandler)
	http.HandleFunc("/robots.txt", singleHandler)

	err = http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil)
	if err != nil {
		log.Fatalf("Server failed to start. Returned error %s\n", err)
	}
}
