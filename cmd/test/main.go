package main

// JUST FOR TESTING PURPOSE

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/zekroTJA/ranna/internal/models"
	"github.com/zekroTJA/ranna/internal/sandbox/docker"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	godotenv.Load()

	dsp, err := docker.NewDockerSandboxProvider()
	if err != nil {
		fmt.Println("client creation failed:", err)
		return
	}

	sbx, err := dsp.CreateSandbox(models.Spec{
		Image:      "golang:alpine",
		Entrypoint: "go run",
		Cmd:        "main.go",
		Subdir:     "random123",
		HostDir:    "/home/mgr/exec",
	})
	if err != nil {
		fmt.Println("sandbox creation failed:", err)
		return
	}
	defer sbx.Delete()

	stdout, stderr, err := sbx.Run()
	if err != nil {
		fmt.Println("sandbox run failed:", err)
		return
	}
	fmt.Println(stdout, stderr)
}
