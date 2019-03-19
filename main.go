package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/mozqnet/go-exploitdb/commands"
)

// Name :
const Name string = "go-exploitdb"

var version = "1.0.0"

func main() {
	var v = flag.Bool("v", false, "Show version")

	if envArgs := os.Getenv("GOVAL_DICTIONARY_ARGS"); 0 < len(envArgs) {
		flag.CommandLine.Parse(strings.Fields(envArgs))
	} else {
		flag.Parse()
	}

	if *v {
		fmt.Printf("go-exploitdb %s \n", version)
		os.Exit(0)
	}

	if err := commands.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
