package webgo

import (
	htpl "html/template"
	"io/ioutil"
)

// This stores all the templates the app will use, and is given to all the request handlers
// It's accessed using the `context`.
type Templates struct {
	Tpls map[string]*htpl.Template
}

// ===

// Function to load templates, parse them and keep them ready in the memory
func (t *Templates) Load(files map[string]string) {
	t.Tpls = make(map[string]*htpl.Template)

	// Looping through list of template files
	for key, filePath := range files {
		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			Err.Log.Fatal("templates.go", "Load()", err)
			return
		}

		// Parsing the file into html template.
		c, err := htpl.New("Error-Template").Parse(string(content))
		if err != nil {
			Err.Log.Fatal("templates.go", "Load()", err)
		}
		t.Tpls[key] = c
		// ===
	}
	// ===
}

// ===
