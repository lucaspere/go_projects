package cmd

import "errors"

var ErrNoServerSpecified = errors.New("you have to specify the remote server.")
var InvalidHttpMethod = errors.New("invalid HTTP method")
var InvalidJsonBody = errors.New("invalid json body")
