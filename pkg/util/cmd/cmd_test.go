package cmd

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestWhoami(t *testing.T) {
	cmd := New("whoami",
		nil,
		time.Duration(10)*time.Second,
		nil,
	)

	output, err := cmd.Run()
	fmt.Printf("I am: %s", output.String())
	if err != nil {
		t.Errorf("expected nil, got: %s", err)
	}

	t.Logf("output is %s", output.String())
}

func TestCmdEksctl(t *testing.T) {
	cmd := New("eksctl",
		[]string{"--help"},
		time.Duration(10)*time.Second,
		nil,
	)

	output, err := cmd.Run()
	if err != nil {
		t.Errorf("expected nil, got: %s", err)
	}

	t.Logf("output is %s", output.String())
	if strings.Contains(output.String(), "eksctl create") == false {
		t.Errorf("Expected eksctl create string in output, got: %s", output.String())
	}
}

func TestCmdBadRegion(t *testing.T) {
	cmd := New("eksctl",
		[]string{"create", "cluster", "--name", "buzz", "--region", "us-east-7000"},
		time.Duration(10)*time.Second,
		[]string{"AWS_ACCESS_KEY_ID=FOO", "AWS_SECRET_ACCESS_KEY=BAR", "AWS_DEFAULT_REGION=us-east-7000"},
	)

	output, err := cmd.Run()
	t.Logf("output is %s", output.String())
	if err == nil {
		t.Error("error from Run = nil; want error")
	}
}
