package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	fmt.Printf("dir: %s\n", filepath.Dir(os.Args[1]))
	path, _ := filepath.Abs(os.Args[1])
	fmt.Printf("path: %s \n", path)
}
