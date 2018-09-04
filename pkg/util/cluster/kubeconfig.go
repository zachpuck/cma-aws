package cluster

import (
	"fmt"
	"github.com/spf13/viper"
	"gitlab.com/mvenezia/cma-aws/pkg/util/k8s"
	"gitlab.com/mvenezia/cma-aws/pkg/util/sshhack"
	"io/ioutil"
	"os"
)

const (
	KubeconfigK8SSecretSuffix = "-kubeconfig"
)

type SSHConnectionHost struct {
	Hostname      string
	Port          string
	KeySecretName string
	Username      string
}

type SSHConnectionOptions struct {
	BastionHost SSHConnectionHost
	TargetHost  SSHConnectionHost
}

func GetKubeConfig(clusterName string, sshConnectionOptions SSHConnectionOptions) ([]byte, error) {
	config, err := k8sutil.GetKubeconfigSecret(generateKubeconfigKey(clusterName), viper.GetString("kubernetes-namespace"))
	if err != nil {
		fmt.Printf("Secret did not exist, going to generate it\n")
		return generateKubeConfig(clusterName, sshConnectionOptions)
	} else {
		fmt.Printf("Found secret\n")
	}
	return config, nil
}

func generateKubeconfigKey(clusterName string) string {
	return clusterName + KubeconfigK8SSecretSuffix
}

func generateKubeConfig(clusterName string, sshConnectionOptions SSHConnectionOptions) ([]byte, error) {
	var targetHostKeyFilename string
	bastionHostKeyFile, err := writeSSHKeyToDiskFromK8S(generateSSHSecretKey(sshConnectionOptions.BastionHost.KeySecretName))
	if err != nil {
		return nil, err
	}
	defer os.Remove(bastionHostKeyFile) // clean up

	if sshConnectionOptions.BastionHost.KeySecretName == sshConnectionOptions.TargetHost.KeySecretName {
		targetHostKeyFilename = bastionHostKeyFile
	} else {
		targetHostKeyFilename, err = writeSSHKeyToDiskFromK8S(generateSSHSecretKey(sshConnectionOptions.TargetHost.KeySecretName))
		if err != nil {
			return nil, err
		}
		defer os.Remove(targetHostKeyFilename) // clean up
	}

	// Shell out to ssh to get the kubeconfig from the master node.  Needs to go through a bastion host
	stuff, err := sshhack.GetKubeConfig(sshhack.GetKubeConfigOptions{
		BastionHost: sshhack.SSHHostOptions{
			Hostname:    sshConnectionOptions.BastionHost.Hostname,
			Port:        sshConnectionOptions.BastionHost.Port,
			KeyFilePath: bastionHostKeyFile,
			Username:    sshConnectionOptions.BastionHost.Username,
		},
		TargetHost: sshhack.SSHHostOptions{
			Hostname:    sshConnectionOptions.TargetHost.Hostname,
			Port:        sshConnectionOptions.TargetHost.Port,
			KeyFilePath: targetHostKeyFilename,
			Username:    sshConnectionOptions.TargetHost.Username,
		},
	})
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return nil, err
	}

	err = k8sutil.CreateKubeconfigSecret(generateKubeconfigKey(clusterName), viper.GetString("kubernetes-namespace"), stuff)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return nil, err
	}
	return stuff, nil
}

func writeSSHKeyToDiskFromK8S(keyName string) (string, error) {
	keyData, err := k8sutil.GetSSHSecret(keyName, viper.GetString("kubernetes-namespace"))
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return "", err
	}
	tmpfile, err := ioutil.TempFile("", "sshfile")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return "", err
	}

	if _, err := tmpfile.Write(keyData); err != nil {
		fmt.Printf("Error: %s\n", err)
		return "", err
	}

	return tmpfile.Name(), nil
}
