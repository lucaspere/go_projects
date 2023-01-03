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

func (c *config) Run(args []string) error {
	err := c.parseArgs(args)
	if err != nil {
		return err
	}
	name, err := getName(os.Stdout, os.Stdin)
	if err != nil {
		return err
	}

	c.GreaterName(name)

	return err
}

func (c *config) GreaterName(name string) {
	for i := 0; i < c.numTimes; i++ {
		fmt.Println("Welcome", name)
	}
}

func validateNumberArgs(args []string) error {
	if len(args) <= 1 {
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

func (c *config) parseArgs(args []string) error {
	if err := validateNumberArgs(args); err != nil {
		return err
	}
	if err := validateArgs(args); err != nil {
		return err
	}

	numTimes, err := strconv.Atoi(args[1])
	if err != nil {
		return fmt.Errorf("%s is not a int type", args[1])
	}

	c.numTimes = numTimes
	return nil
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

func main() {
	var c config
	args := os.Args[1:]
	if err := c.Run(args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}
