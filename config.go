package webgo

import (
	"encoding/json"
	htpl "html/template"
	"io/ioutil"
	"strconv"
)

// Config is used for reading app's configuration from json file
type Config struct {
	// Env is the deployment environment
	Env string `json:"environment"`
	// Host is the host on which the server is listening
	Host string `json:"host,omitempty"`
	// Port is the port number where the server has to listen for the HTTP requests
	Port string `json:"port"`

	// CertFile is the TLS/SSL certificate file path, required for HTTPS
	CertFile string `json:"certFile,omitempty"`
	// KeyFile is the filepath of private key of the certificate
	KeyFile string `json:"keyFile,omitempty"`
	// HTTPSPort is the port number where the server has to listen for the HTTP requests
	HTTPSPort string `json:"httpsPort,omitempty"`
	// HTTPSOnly if true will enable HTTPS server alone
	HTTPSOnly bool `json:"httpsOnly,omitempty"`

	// TemplatesBasePath is the base path where all the HTML templates are located
	TemplatesBasePath string `json:"templatePath,omitempty"`
}

// Load config file from the provided filepath and validate
func (cfg *Config) Load(filepath string) {
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		errLogger.Fatal(err)
	}

	err = json.Unmarshal(file, cfg)
	if err != nil {
		errLogger.Fatal(err)
	}

	cfg.Validate()
}

// Validate the config parsed into the Config struct
func (cfg *Config) Validate() {
	i, err := strconv.Atoi(cfg.Port)
	if err != nil {
		errLogger.Fatal(ErrInvalidPort)
	}
	if i <= 0 || i > 65535 {
		errLogger.Fatal(ErrInvalidPort)
	}
}

// Globals struct to hold configurations which are shared with all the request handlers via context.
type Globals struct {
	// Cfg has all the webgo configurations
	Cfg *Config

	// Templates stores all the templates pre-compiled and accessible.
	Templates map[string]*htpl.Template

	// App stores any key value configuration. This can be app specific (i.e. any app using WebGo)
	App map[string]interface{}
}

// NewGlobals returns a new Globals instance pointer with the provided configurations
func NewGlobals(cfg *Config, app map[string]interface{}, tpls map[string]*htpl.Template) (*Globals, error) {
	g := Globals{
		App:       app,
		Templates: tpls,
		Cfg:       cfg,
	}
	return &g, nil
}
