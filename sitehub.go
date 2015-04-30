package main

import (
	"fmt"
	"log"
	"net/http"
	"path"
)

type Data struct {
	Host, Path, Method, Proto, RemoteIP string
	Header                              http.Header
	SiteConfig, SiteData                map[string]interface{}
}

type T []string

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
		w.WriteHeader(http.StatusNotFound)
		custom404 := path.Join(hp, "404.html")
		fmt.Fprint(w, "404, page not found!")
	} else {

		layout, fourOfour, templates := loadTemplates(hp)

		err := loadJSONFile(path.Join(hp, "config.json"), &data.SiteConfig)
		if err != nil {
			log.Printf("No site configuration: %s\n", err)
		}

		err := loadJSONFile(path.Join(hp, "data.json"), &data.SiteData)
		if err != nil {
			log.Printf("No site data: %s\n", err)
		}
	}
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
	hp, err := hostPath(r.Host)
	if err != nil {
		staticPath := path.Join(hp, "static")
		if validatePath(staticPath) {
			http.Handle(r.Host+"/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(staticPath))))
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
