package main

import (
	"fmt"
	"github.com/satori/go.uuid"
)

func main() {
	fmt.Printf("ID1: %v\n", uuid.NewV1())
	//fmt.Printf("ID2: %v\n", uuid.NewV2())
	//fmt.Printf("ID3: %v\n", uuid.NewV3())
	fmt.Printf("ID4: %v\n", uuid.NewV4())
	//fmt.Printf("ID5: %v\n", uuid.NewV5())
}
