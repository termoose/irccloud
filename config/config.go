package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
)

type ConfigData struct {
	Username string   `yaml:"username"`
	Password string   `yaml:"password"`
	Triggers []string `yaml:"triggers"`
	LastChan string   `yaml:"last_chan"`
}

func Parse() ConfigData {
	var result ConfigData
	//currUser, _ := user.Current()
	//confDir := filepath.Join(currUser.HomeDir, "/.config/irccloud/")
	filename, configDir := getPaths()

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

func WriteLatestChannel(data ConfigData, latest string) {
	data.LastChan = latest
	filename, _ := getPaths()

	writeConfig(filename, data)
}

func writeConfig(filename string, data ConfigData) {
	content, _ := yaml.Marshal(&data)
	if err := ioutil.WriteFile(filename, content, 0600); err != nil {
		fmt.Printf("Could not write config to to file %s\n", filename)
	}
}

func writeDummyConfig(filename string) ConfigData {
	dummy := ConfigData{
		Username: "your_username_here",
		Password: "secret_password_here",
		Triggers: []string{},
		LastChan: "",
	}

	writeConfig(filename, dummy)
	//content, _ := yaml.Marshal(&dummy)
	//if err := ioutil.WriteFile(filename, content, 0600); err != nil {
	//	fmt.Printf("Could not write config to to file %s\n", filename)
	//}

	return dummy
}
