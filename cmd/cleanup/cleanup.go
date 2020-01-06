package main

import (
	"fmt"

"github.com/xavier268/go-s3sync/pkg/gosync"
)

func main() {
	fmt.Println("Cleaning empty directories after sync")

	c := gosync.NewConfig().SetMode(gosync.ModeCleanEmptyDirs)
	fmt.Println(c)

	fmt.Printf("If that configuration is correct, type 'yes' to continue:")
	yes := ""
	fmt.Scanln(&yes)
	if yes == "yes" {
		c.RemoveAllEmptyDirs()
	} else {
		fmt.Println("Aborting ...")
	}

}
