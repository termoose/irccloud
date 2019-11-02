package ui

import (
	"io/ioutil"
	"log"
)

func readFile(filename string) string {
	dat, err := ioutil.ReadFile(filename)

	if err != nil {
		log.Fatal("Could not open file!")
	}

	return string(dat)
}
