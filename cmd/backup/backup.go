package main

import (
	"fmt"

	"github.com/xavier268/go-s3sync/pkg/gosync"
)

func main() {
	fmt.Println("Backing up files into s3")

	c := gosync.NewConfig().SetMode(gosync.ModeBackup)
	fmt.Println(c)

	fmt.Printf("If that configuration is correct, type 'yes' to continue:")
	yes := ""
	fmt.Scanln(&yes)
	if yes == "yes" {
		c.ProcessObjects()
		c.ProcessFiles()
	} else {
		fmt.Println("Aborting ...")
	}

}
