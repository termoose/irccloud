package logging

import (
	"fmt"
	"log"
	"os"
)

type logger struct {
	handle *os.File
}

func CreateLogger(name string) *logger {
	file, err := os.OpenFile(fmt.Sprintf("%s.log", name), os.O_CREATE | os.O_APPEND, 0644)

	if err == nil {
		return &logger{
			handle: file,
		}
	}

	return nil
}

func (l *logger) Info(line string) {
	log.SetOutput(l.handle)
	log.Print(line)
}