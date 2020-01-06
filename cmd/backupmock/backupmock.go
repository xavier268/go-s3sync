package main

import (
	"fmt"

	"github.com/xavier268/go-s3sync/pkg/gosync"
)

func main() {
	fmt.Println("Simulating backing up files into s3")

	c := gosync.NewConfig().SetMode(gosync.ModeBackupMock)
	fmt.Println(c)

	c.ProcessObjects()
	c.ProcessFiles()

}
