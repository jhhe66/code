package main
 
import (
    "flag"
    "fmt"
)
 
var (
     openPort = flag.String("port",":8080","http listen port")
     configFile = flag.String( "config","","config file name" )
)
 
func main(){
    flag.Parse()
    fmt.Println("Configure Filename: ",*configFile)
    fmt.Println("Open Port: ",*openPort)
}