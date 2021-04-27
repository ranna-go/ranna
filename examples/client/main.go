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
		Endpoint: "https://api.ranna.zekro.de",
	})
	must(err)

	fmt.Println(c.Spec())
	fmt.Println(c.Exec(models.ExecutionRequest{
		Language:  "python3",
		Code:      code,
		Arguments: []string{"these", "are", "some", "args", "man"},
		Environment: map[string]string{
			"TEST": "some stuff over here",
		},
	}))
}
