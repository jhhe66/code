// QssServer project main.go
package main

import (
//"os"
//"runtime/pprof"
)

func main() {

	//f, _ := os.Create("profile")

	//pprof.StartCPUProfile(f)

	InitConf()

	Log_init()

	ConfTest()

	Run()

	//pprof.StopCPUProfile()
}
