package main

import (
	"errors"
	. "flag-parse/flag-parse"
	"fmt"
	"os"
)

func main() {
	c, err := ParseArgs(os.Stderr, os.Args[1:])
	if err != nil {
		if errors.Is(err, ErrInvalidPosArgSpecified) {
			fmt.Fprintln(os.Stdout, err)
		}
		os.Exit(1)
	}
	err = ValidateArgs(c)
	if err != nil {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}
	err = RunCmd(os.Stdin, os.Stdout, c)
	if err != nil {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}
}
