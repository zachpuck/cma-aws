package sshhack

import (
	"os/exec"
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
	stuff, err := exec.Command(
		"/bin/bash",
		[]string{
			"-c",
			`ssh -i ` + options.BastionHost.KeyFilePath + ` ` + StrictHostKeyCheckingOption + ` -o ProxyCommand="` + options.GenerateProxyCommand() + `" ` + options.TargetHost.GenerateUserAtHost() + ` cat ./kubeconfig`}...).Output()
	if err != nil {
		return nil, err
	}
	return stuff, nil
}
