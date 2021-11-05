package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
)

type Data struct {
	Username     string   `yaml:"username"`
	Password     string   `yaml:"password"`
	Triggers     []string `yaml:"triggers"`
	LastChan     string   `yaml:"last_chan"`
	OnlyMessages bool     `yaml:"only_messages"`
}

func ParseCustom(filename string) Data {
	return parseData(filename)
}

func Parse() Data {
	filename, _ := getPaths()
	return parseData(filename)
}

func parseData(filename string) Data {
	var result Data
	configDir := filepath.Base(filename)

	// Don't care if this fails
	_ = os.MkdirAll(configDir, 0700)

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

func getPaths() (string, string) {
	currUser, _ := user.Current()
	confDir := filepath.Join(currUser.HomeDir, "/.config/irccloud/")
	return filepath.Join(confDir, "config.yaml"), confDir
}

func WriteLatestChannel(data Data, latest string) {
	data.LastChan = latest
	filename, _ := getPaths()

	writeConfig(filename, data)
}

func writeConfig(filename string, data Data) {
	content, _ := yaml.Marshal(&data)
	if err := ioutil.WriteFile(filename, content, 0600); err != nil {
		fmt.Printf("Could not write config to to file %s\n", filename)
	}
}

func writeDummyConfig(filename string) Data {
	dummy := Data{
		Username: "your_username_here",
		Password: "secret_password_here",
		Triggers: []string{},
		LastChan: "",
		OnlyMessages: false,
	}

	writeConfig(filename, dummy)
	return dummy
}
