package main

func loadTemplates(p) (string, string, []string) {
	glb := path.Join(p, "*.html")
	files, err := filepath.Glob(glb)
	if err != nil {
		return nil, nil, nil
	}
}
