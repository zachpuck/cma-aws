package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/juju/loggo"
	log "github.com/samsung-cnct/cma-aws/pkg/util"
)

var (
	logger loggo.Logger
)

type Cmd struct {
	// Name is the command name
	Name string
	// Args are the optional list of command args
	Args []string
	// Timeout is used to kill the command process if it runs too long
	Timeout time.Duration
	// Environ is optional environmentals needed by the command
	// Environ format is FOO=value for each string
	Environ []string
}

func New(name string, args []string, timeout time.Duration, environ []string) *Cmd {
	return &Cmd{
		Name:    name,
		Args:    args,
		Timeout: timeout,
		Environ: environ,
	}
}

func (c *Cmd) Run() (bytes.Buffer, error) {
	var streamOut, streamErr bytes.Buffer
	logger = log.GetModuleLogger("pkg.util.cmd", loggo.INFO)
	logger.Infof("Running command \"%v %v\"\n", c.Name, strings.Join(c.Args, " "))
	cmd := exec.Command(c.Name, c.Args...)

	// not sure we should include all existing os envs in case AWS creds exist
	//cmd.Env = os.Environ()
	for _, env := range c.Environ {
		cmd.Env = append(cmd.Env, env)
	}
	// including the os PATH
	cmd.Env = append(cmd.Env, "PATH="+os.Getenv("PATH"))

	cmd.Stdout = &streamOut
	cmd.Stderr = &streamErr

	err := cmd.Start()
	if err != nil {
		return streamOut, err
	}

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-time.After(c.Timeout):
		// We do not print stdout because it may contain secrets.
		fmt.Fprintf(os.Stderr, "Command %v stderr: %v\n", c.Name, string(streamErr.Bytes()))

		if err := cmd.Process.Kill(); err != nil {
			panic(fmt.Sprintf("Failed to kill command %v, err %v", c.Name, err))
		}
		return streamOut, fmt.Errorf("Command %v timed out", c.Name)
	case err := <-done:
		// We do not print stdout because it may contain secrets.
		fmt.Fprintf(os.Stderr, "Command %v stderr: %v\n", c.Name, string(streamErr.Bytes()))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Command %v returned err %v\n", c.Name, err)
			return streamOut, err
		}
	}
	logger.Infof("Command %v completed successfully\n", c.Name)

	return streamOut, nil
}
