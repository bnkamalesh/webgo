package webgo

import (
	htpl "html/template"
	"io/ioutil"
	"strconv"

	simplejson "github.com/bitly/go-simplejson"
)

// Struct for reading app's configuration from json file
type Config struct {
	Env               string `json:"environment"`
	Host              string `json:"host,omitempty"`
	Port              string `json:"port"`
	TemplatesBasePath string `json:"templatePath"`

	DBConfig struct {
		Name     string `json:"name"`
		Host     string `json:"host"`
		Port     string `json:"port"`
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"dbConfig"`

	// The whole json file read from the config file
	Data *simplejson.Json

	// If the app needs to add some extra info to the config, simple string key, value pairs
	Misc map[string]string
}

// ===

// Load config file from the provided filepath
func (cfg *Config) Load(filepath string) {
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		Err.Fatal("config.go", "Load() [1] - could not read file", err)
	}

	decode, decodeErr := simplejson.NewJson(file)
	if decodeErr != nil {
		Err.Fatal("config.go", "Load() [2] - could not decode json file", err)
	}
	// Assigning the whole json data fetched from the config file
	cfg.Data = decode

	// Setting the config values gotten from the file
	cfg.Env = cfg.Data.Get("environment").MustString()
	cfg.Host = cfg.Data.Get("host").MustString()
	cfg.Port = cfg.Data.Get("port").MustString()
	cfg.TemplatesBasePath = cfg.Data.Get("templatePath").MustString()

	// Reading dbConfig from the json file
	dbConf, err := cfg.Data.Get("dbConfig").Map()
	if err != nil {
		Err.Fatal("config.go", "Load() [3]- could not get dbConfig from file", err)
	}

	// Setting DB config values
	cfg.DBConfig.Name = dbConf["name"].(string)
	cfg.DBConfig.Host = dbConf["host"].(string)
	cfg.DBConfig.Port = dbConf["port"].(string)
	cfg.DBConfig.Username = dbConf["username"].(string)
	cfg.DBConfig.Password = dbConf["password"].(string)

	cfg.Misc = make(map[string]string)

	cfg.Validate()
}

// ===

func (cfg *Config) Validate() {
	if cfg.Env != "production" && cfg.Env != "development" {
		Err.Fatal("webgo - config.go", "Validate() - [1]", Err.C003)
	}

	i, err := strconv.Atoi(cfg.Port)
	if err != nil {
		Err.Fatal("webgo - config.go", "Validate() - [2]", Err.C004)
	}
	if i <= 0 || i > 65535 {
		Err.Fatal("webgo - config.go", "Validate() - [3]", Err.C004)
	}
}

//	Add any global app configurations here.
//	They're available to every single request handler, via context.
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

// ===

// Add a custom global config
func (g *Globals) Add(key string, data interface{}) {
	g.App[key] = data
}

// ===

// Initialize the Context and set appropriate values
func (g *Globals) Init(cfg *Config, tpls map[string]*htpl.Template, ds *DataStore) {
	g.App = make(map[string]interface{})
	g.Templates = make(map[string]*htpl.Template)
	g.Cfg = cfg
	g.Templates = tpls
	g.Db = ds
	g.Params = make(map[string]string)
}

// ===
