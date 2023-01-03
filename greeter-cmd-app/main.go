package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
)

type config struct {
	numTimes int
}

func validateArgs(args []string) error {
	if len(args) != 1 {
		return errors.New("Must specify a number greater than 0")
	}

	return nil
}

func main() {
	commands := os.Args[1:]
	fmt.Println(commands)
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
