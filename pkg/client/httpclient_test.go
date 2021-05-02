package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ranna-go/ranna/pkg/models"
)

func TestSpec(t *testing.T) {
	testSpec := &models.Spec{
		Image:      "golang:alpine",
		Entrypoint: "go run",
		FileName:   "main.go",
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(models.SpecMap{
			"go": testSpec,
		})
	}))
	defer ts.Close()

	client, err := New(Options{
		Endpoint:      ts.URL,
		Version:       "v1",
		Authorization: "basic test",
		UserAgent:     "test agent",
	})
	if err != nil {
		t.Fatal(err)
	}

	recSpecMap, err := client.Spec()
	if err != nil {
		t.Error(err)
	}
	recSpec, ok := recSpecMap.Get("go")
	if !ok {
		t.Error("could not recover spec map entry")
	}
	if recSpec.Image != testSpec.Image {
		t.Errorf("Image value was invalid: %s", recSpec.Image)
	}
	if recSpec.Entrypoint != testSpec.Entrypoint {
		t.Errorf("Entrypoint value was invalid: %s", recSpec.Entrypoint)
	}
	if recSpec.FileName != testSpec.FileName {
		t.Errorf("FileName value was invalid: %s", recSpec.FileName)
	}
}

func TestExec(t *testing.T) {
	testExec := &models.ExecutionResponse{
		StdOut:     "stdput",
		StdErr:     "stderr",
		ExecTimeMS: 1337,
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(testExec)
	}))
	defer ts.Close()

	client, err := New(Options{
		Endpoint:      ts.URL,
		Version:       "v1",
		Authorization: "basic test",
		UserAgent:     "test agent",
	})
	if err != nil {
		t.Fatal(err)
	}

	recExec, err := client.Exec(models.ExecutionRequest{})
	if err != nil {
		t.Error(err)
	}
	if recExec.StdOut != testExec.StdOut {
		t.Errorf("StdOut value was invalid: %s", recExec.StdOut)
	}
	if recExec.StdErr != testExec.StdErr {
		t.Errorf("StdErr value was invalid: %s", recExec.StdErr)
	}
	if recExec.ExecTimeMS != testExec.ExecTimeMS {
		t.Errorf("ExecTimeMS value was invalid: %d", recExec.ExecTimeMS)
	}
}
