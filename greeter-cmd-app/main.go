package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
)

type config struct {
	numTimes int
}

func validateNumberArgs(args []string) error {
	if len(args) != 1 {
		return errors.New("must specify a number greater than 0")
	}

	return nil
}

func validateArgs(args []string) error {
	if args[0] != "-n" {
		return fmt.Errorf("invalid %s argument. Must specify a valid argument %s", args[0], "-n")
	}

	return nil
}

func parseArgs(args []string) (config, error) {
	var c config
	if err := validateNumberArgs(args); err != nil {
		return c, err
	}
	if err := validateArgs(args); err != nil {
		return c, err
	}

	numTimes, err := strconv.Atoi(args[1])
	if err != nil {
		return c, fmt.Errorf("%s is not a int type", args[1])
	}

	c.numTimes = numTimes
	return c, nil
}
func main() {
	args := os.Args[1:]
	_, err := parseArgs(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func getName(w io.Writer, r io.Reader) (string, error) {
	msg := "Your name please? Press the Enter key when done. \n"
	fmt.Fprint(w, msg)

	scanner := bufio.NewScanner(r)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return "", err
	}

	name := scanner.Text()
	if len(name) == 0 {
		return "", errors.New("you didn't enter your name")
	}

	return name, nil
}
