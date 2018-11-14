package sshhack

import (
	"fmt"
	"os/exec"
	"gopkg.in/alessio/shellescape.v1"
)

const (
	StrictHostKeyCheckingOption = "-o StrictHostKeyChecking=no"
)

type SSHHostOptions struct {
	Hostname    string
	Port        string
	KeyFilePath string
	Username    string
}

func (t *SSHHostOptions) GenerateUserAtHost() string {
	return t.Username + "@" + t.Hostname
}

func (t *SSHHostOptions) GenerateNetCatCommand() string {
	return "nc " + t.Hostname + " " + t.Port
}

type GetKubeConfigOptions struct {
	TargetHost  SSHHostOptions
	BastionHost SSHHostOptions
}

func (t *GetKubeConfigOptions) GenerateProxyCommand() string {
	// Should look like ssh -o StrictHostKeyChecking=no -i \"${SSH_KEY}\" -p 22 ubuntu@18.207.4.203 nc 10.0.5.199 22"
	return "ssh " + StrictHostKeyCheckingOption + " -i " + t.BastionHost.KeyFilePath + " -p " + t.BastionHost.Port + " " + t.BastionHost.GenerateUserAtHost() + " " + t.TargetHost.GenerateNetCatCommand()
}

func GetKubeConfig(options GetKubeConfigOptions) ([]byte, error) {

	// Notes about commandString:
	// the following commandString is used to retrieve the master elb URL from ./.kube/config, modify it to use port 6443, and set it as the server url in ./kubeconfig.
	// it will then return the ./kubeconfig file with the replaced "server: " url instead of an ip address, this will allow you to always connect to the master no mater if the ip changes
	// additionally 'shellescape' is used to escape all the of the single quotes in commandString, for exec.Command to work correctly
	commandString := fmt.Sprintf("sed 's#server: .*#server: '$(sed -re 's#(https?://.*):([0-9]+)#\\1:6\\2#' %s | grep -Po 'https?://.*')'#' %s", "./.kube/config", "./kubeconfig")

	stuff, err := exec.Command(
		"/bin/bash",
		[]string{
			"-c",
			`ssh -i ` + options.BastionHost.KeyFilePath + ` ` + StrictHostKeyCheckingOption + ` -o ProxyCommand="` + options.GenerateProxyCommand() + `" ` + options.TargetHost.GenerateUserAtHost() + ` ` + shellescape.Quote(commandString)}...).Output()
	if err != nil {
		return nil, err
	}
	return stuff, nil
}
