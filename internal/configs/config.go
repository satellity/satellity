package configs

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// Constants of configs
const (
	BuildVersion = "BUILD_VERSION"
)

// Option for configurations
type Option struct {
	Name string `yaml:"name"`
	HTTP struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"http"`
	Database struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Name     string `yaml:"name"`
	} `yaml:"database"`
	Github struct {
		ClientID     string `yaml:"client_id"`
		ClientSecret string `yaml:"client_secret"`
	} `yaml:"github"`
	System struct {
		Attachments struct {
			Storage string `yaml:"storage"`
			Path    string `yaml:"path"`
		} `yaml:"attachments"`
	} `yaml:"system"`
	Recaptcha struct {
		URL     string `yaml:"url"`
		SiteKey string `yaml:"site_key"`
		Secret  string `yaml:"secret"`
	} `yaml:"recaptcha"`
	Operators []string `yaml:"operators"`
	Email     struct {
		Verification struct {
			Title string `yaml:"title"`
			Reset string `yaml:"reset"`
			Body  string `yaml:"body"`
		} `yaml:"verification"`
	} `yaml:"email"`
	Mailgun struct {
		Domain string `yaml:"domain"`
		Key    string `yaml:"key"`
		Sender string `yaml:"sender"`
	}

	Environment string
	OperatorSet map[string]bool
}

// AppConfig is the configs for the whole application
var AppConfig *Option

// Init is using to initialize the configs
func Init(file, env string) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	var options map[string]Option
	err = yaml.Unmarshal(data, &options)
	if err != nil {
		return err
	}
	opt := options[env]
	opt.Environment = env
	opt.OperatorSet = make(map[string]bool)
	for _, operator := range opt.Operators {
		opt.OperatorSet[operator] = true
	}
	AppConfig = &opt
	return nil
}
