package main

import (
	"fmt"

	"github.com/go-s3sync/pkg/synchro"
)

func main() {
	fmt.Println("Synchronizing files into s3")
	c := synchro.NewTestConfig()
	fmt.Printf("Configuration : %v\n", c)

	synchro.Check(c)

}
