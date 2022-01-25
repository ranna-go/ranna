package main

import (
	"fmt"

	ranna "github.com/ranna-go/ranna/pkg/client"
	"github.com/ranna-go/ranna/pkg/models"
)

const (
	code = `
import sys
import os

print('Hello World!')
print('args:', sys.argv)
print('env:', os.environ.get("TEST"))
`
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	c, err := ranna.New(ranna.Options{
		Endpoint: "https://api.ranna.dev",
	})
	must(err)

	fmt.Println(c.Spec())

	// This is a normal code execution
	fmt.Println(c.Exec(models.ExecutionRequest{
		Language:         "python3",
		Code:             code,
		InlineExpression: false,
		Arguments:        []string{"these", "are", "some", "args", "man"},
		Environment: map[string]string{
			"TEST": "some stuff over here",
		},
	}))

	// ranna also supports inline expressions to run code without having to provide the boilerplate code around it
	// (package declaration, main function)
	//
	// It is important to know though that this feature is not supported by all languages and the packages which are
	// imported by default are limited. Please have a look at the language specification of your ranna instance.
	fmt.Println(c.Exec(models.ExecutionRequest{
		Language:         "go",
		Code:             `fmt.Println("Hello, world!")`, // Note how this is only one small expression
		InlineExpression: true,
		Arguments:        []string{},
		Environment:      map[string]string{},
	}))
}
