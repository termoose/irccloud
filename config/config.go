package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
)

type confdata struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func Parse() confdata {
	var result confdata
	curr_user, _ := user.Current()
	conf_dir := filepath.Join(curr_user.HomeDir, "/.irccloud/")
	filename := filepath.Join(conf_dir, "config.yaml")

	// Don't care if this fails
	os.Mkdir(conf_dir, 0700)

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

func writeDummyConfig(filename string) confdata {
	dummy := confdata{
		Username: "your_username_here",
		Password: "secret_password_here",
	}

	content, _ := yaml.Marshal(&dummy)
	if err := ioutil.WriteFile(filename, content, 0600); err != nil {
		fmt.Printf("Could not write config to to file %s\n", filename)
	}

	return dummy
}
