package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/jarxorg/tree"
)

const (
	cmdName = "tq"
)

var (
	isHelp = flag.Bool("help", false, "help for "+cmdName)
)

func main() {
	flag.Parse()

	if *isHelp {
		flag.Usage()
		return
	}

	if err := exec(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}

func exec() error {
	in, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return err
	}
	node, err := tree.UnmarshalJSON(in)
	if err != nil {
		return err
	}
	query, err := tree.ParseQuery(flag.Arg(0))
	if err != nil {
		return err
	}
	result, err := query.Exec(node)
	if err != nil {
		return err
	}
	out, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(out))
	return nil
}
