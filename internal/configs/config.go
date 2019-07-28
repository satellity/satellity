package configs

import (
	"io/ioutil"
	"path"

	"gopkg.in/yaml.v3"
)

const (
	BuildVersion = "BUILD_VERSION"
)

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
	Operators []string `yaml:"operators"`

	Environment string
	OperatorSet map[string]bool
}

var AppConfig *Option

func Init(dir, env string) error {
	data, err := ioutil.ReadFile(path.Join(dir, "./config.yaml"))
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
