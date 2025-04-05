package config

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Config struct {
	Token  string  `yaml:"token"`
	Debug  bool    `yaml:"debug"`
	Admins []int64 `yaml:"admins"`
	Menu   Button  `yaml:"menu"`
}

type Button struct {
	Title   string        `yaml:"title"`
	Buttons []Button      `yaml:"buttons"`
	Command ConfigCommand `yaml:"command"`
}

type ConfigCommand struct {
	Name string   `yaml:"name"`
	Args []string `yaml:"args"`
}

func (menu *Button) Validate() {
	if menu.Title == "" {
		log.Fatal("Config menu: `title` is required field in menu.")
	}

	if menu.Buttons != nil {
		for _, item := range menu.Buttons {
			item.Validate()
		}
	}
}

func ParseFlags() (string, error) {
	var configPath string

	flag.StringVar(&configPath, "config.file", "./config.yml", "path to config file")

	flag.Parse()

	if err := validateConfigPath(configPath); err != nil {
		return "", err
	}

	return configPath, nil
}

func NewConfig(configPath string) (*Config, error) {
	config := &Config{}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)
	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	config.Menu.Validate()

	return config, nil
}

func validateConfigPath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a normal file", path)
	}
	return nil
}
