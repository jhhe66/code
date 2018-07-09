// log.go project main.go
package main

import (
	"log4go"
)

const (
	LogConfigFile = "log.xml"
)

func Log_init() {
	log4go.LoadConfiguration(LogConfigFile)
}
