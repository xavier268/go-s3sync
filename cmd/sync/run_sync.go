package main

import (
	"fmt"

	"github.com/go-s3sync/pkg/gosync"
)

func main() {
	fmt.Println("Synchronizing files into s3")
	c := gosync.NewTestConfig()
	fmt.Printf("Configuration : %v\n", c)

	c.CheckS3()
	c.CheckFiles()

}
