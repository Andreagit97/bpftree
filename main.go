package main

import (
	"fmt"
	"os"

	"github.com/Andreagit97/bpftree/cmd"
)

// Main is a pass-through method to be compliant with testscript package.
func Main() int {
	if err := cmd.Execute(); err != nil {
		fmt.Println(("Error during command execution: " + err.Error()))
		return 1
	}
	return 0
}

func main() {
	os.Exit(Main())
}
