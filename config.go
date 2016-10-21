package webgo

import (
	"encoding/json"
	htpl "html/template"
	"io/ioutil"
	"strconv"
)

// Config struct for reading app's configuration from json file
type Config struct {
	Env               string `json:"environment"`
	Host              string `json:"host,omitempty"`
	Port              string `json:"port"`
	TemplatesBasePath string `json:"templatePath"`

	DBC DBConfig `json:"dbConfig"`

	// Data holds the full json config file data as bytes
	Data []byte
}

// Load config file from the provided filepath
func (cfg *Config) Load(filepath string) {
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		Log.Fatal(err)
	}

	err = json.Unmarshal(file, cfg)
	if err != nil {
		Log.Fatal(err)
	}

	cfg.Data = file

	cfg.Validate()
}

// Validate the config parsed into the Config struct
func (cfg *Config) Validate() {
	if cfg.Env != "production" && cfg.Env != "development" {
		Log.Fatal(C003)
	}

	i, err := strconv.Atoi(cfg.Port)
	if err != nil {
		Log.Fatal(C004)
	}
	if i <= 0 || i > 65535 {
		Log.Fatal(C004)
	}
}

//Globals struct to hold configurations which are shared with all the request handlers via context.
type Globals struct {
	// Multiplexer params
	Params map[string]string

	// All the app configurations
	Cfg *Config

	// All templates, which can be accessed anywhere from the app
	Templates map[string]*htpl.Template

	// Data store handler from the Database handling library
	Db *DataStore

	// This can be used to add any app specifc data, which needs to be shared
	// E.g. This can be used to plug in a new DB driver, if someone does not want to use MongoDb
	App map[string]interface{}
}

// Add a custom global config
func (g *Globals) Add(key string, data interface{}) {
	g.App[key] = data
}

//Init initializes the Context and set appropriate values
func (g *Globals) Init(cfg *Config, tpls map[string]*htpl.Template, ds *DataStore) {
	g.App = make(map[string]interface{})
	g.Templates = make(map[string]*htpl.Template)
	g.Cfg = cfg
	g.Templates = tpls
	g.Db = ds
	g.Params = make(map[string]string)
}
