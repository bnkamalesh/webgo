package webgo

import (
	"html/template"
	"io/ioutil"
)

// Templates stores all the templates the app will use, and is given to all the request handlers via context
type Templates struct {
	Tpls map[string]*template.Template
}

// Load templates, all the HTML files are parsed and made available for instant use
func (t *Templates) Load(files map[string]string) {
	t.Tpls = make(map[string]*template.Template)

	// Looping through list of template files
	for key, filePath := range files {
		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			Log.Fatal(err)
			return
		}

		// Parsing the file into html template.
		c, err := template.New(key).Parse(string(content))
		if err != nil {
			Log.Fatal(err)
		}
		t.Tpls[key] = c
	}
}
