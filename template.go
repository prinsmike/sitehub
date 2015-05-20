package main

import (
	"errors"
	"log"
	"path"
	"path/filepath"
)

func loadTemplates(p string) (templates []string, err error) {

	layout := path.Join(p, "layout.html")
	templates = []string{layout}

	if !validateFile(layout) {
		return nil, errors.New("No layout file. You should create a layout.html file for this website.")
	}

	// Get the file names.
	glb := path.Join(p, "*.html")
	files, err := filepath.Glob(glb)
	if err != nil {
		log.Printf("Could not get template files for path %s: %s\n", p, err)
		return nil, err
	}

	// Filter.
	for _, f := range files {
		if filepath.Base(f) != "layout.html" && filepath.Base(f) != "404.html" {
			templates = append(templates, f)
		}
	}

	return templates, nil
}
