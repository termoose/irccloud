package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
)

type Confdata struct {
	Username string   `yaml:"username"`
	Password string   `yaml:"password"`
	Triggers []string `yaml:"triggers"`
}

func Parse() Confdata {
	var result Confdata
	currUser, _ := user.Current()
	confDir := filepath.Join(currUser.HomeDir, "/.config/irccloud/")
	filename := filepath.Join(confDir, "config.yaml")

	// Don't care if this fails
	_ = os.MkdirAll(confDir, 0700)

	f, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Could not find config, creating dummy in %s\n", filename)
		return writeDummyConfig(filename)
	}

	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(&result); err != nil {
		panic("Invalid config format!")
	}

	return result
}

func writeDummyConfig(filename string) Confdata {
	dummy := Confdata{
		Username: "your_username_here",
		Password: "secret_password_here",
		Triggers: []string{},
	}
	content, _ := yaml.Marshal(&dummy)
	if err := ioutil.WriteFile(filename, content, 0600); err != nil {
		fmt.Printf("Could not write config to to file %s\n", filename)
	}

	return dummy
}
