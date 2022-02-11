package main

import (
	"os"

	"github.com/ahmadrosid/heline/cmd"
)

func main() {
	os.Exit(cmd.Run(os.Args[1:]))
}
