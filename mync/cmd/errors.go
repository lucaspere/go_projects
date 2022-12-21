package cmd

import "errors"

var ErrNoServerSpecified = errors.New("You have to specify the remote server.")
var InvalidHttpMethod = errors.New("Invalid HTTP method")
