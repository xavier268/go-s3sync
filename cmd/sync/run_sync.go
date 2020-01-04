package main

import (
	"fmt"

	"github.com/go-s3sync/pkg/gosync"
)

func main() {
	fmt.Println("Synchronizing files into s3")

	c := gosync.NewConfig()
	fmt.Println(c)

	c.CheckObjects()
	c.CheckFiles()

}
