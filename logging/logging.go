package logging

import (
	"log"
	"os"
	"sync"
)

var once sync.Once
var logger *log.Logger

func GetLogger() *log.Logger {

	once.Do(func() {
		logger = log.New(os.Stdout, "logger: ", log.Lmsgprefix|log.Lshortfile)
	})

	return logger
}
