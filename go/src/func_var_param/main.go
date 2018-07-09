package main

import (
	"fmt"
)

func Debug(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}

func Info(format string, v ...interface{}) {
	str := fmt.Sprintf(format, v...)
	fmt.Printf(str)
}

func main() {
	Debug("Debug: %s %f \n", "hello", 1.2322)
	Info("Info: %s %f\n", "world", 2.00)
}
