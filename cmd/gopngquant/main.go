package main

import (
	"flag"
	"fmt"
	"os"

	"code.ivysaur.me/imagequant"
)

func main() {
	ShouldDisplayVersion := flag.Bool("Version", false, "")

	flag.Parse()

	if *ShouldDisplayVersion {
		fmt.Printf("libimagequant %d\n", imagequant.GetLibraryVersion())
		os.Exit(1)
	}
}
