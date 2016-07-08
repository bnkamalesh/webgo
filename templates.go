package webgo

import (
	htpl "html/template"
	"io/ioutil"
)

// Templates stores all the templates the app will use, and is given to all the request handlers via context
type Templates struct {
	Tpls map[string]*htpl.Template
}

// Load templates, all the HTML files are parsed and made available for instant use
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
		c, err := htpl.New(key).Parse(string(content))
		if err != nil {
			Err.Log.Fatal("templates.go", "Load()", err)
		}
		t.Tpls[key] = c
	}
}
