package cluster

import (
	"github.com/spf13/viper"
	"gitlab.com/mvenezia/cma-aws/pkg/util/awsutil"
	"gitlab.com/mvenezia/cma-aws/pkg/util/awsutil/models"
	"gitlab.com/mvenezia/cma-aws/pkg/util/k8s"
)

const (
	SSHK8SSecretSuffix = "-ssh"
)

func generateSSHSecretKey(clusterName string) string {
	return clusterName + SSHK8SSecretSuffix
}

func ProvisionAdnSaveSSHKey(clusterName string, credentials awsmodels.Credentials) (string, error) {
	privateKey, err := awsutil.CreateKey(clusterName, credentials)
	if err != nil {
		return "", err
	}

	err = k8sutil.CreateSSHSecret(generateSSHSecretKey(clusterName), viper.GetString("kubernetes-namespace"), []byte(privateKey))
	if err != nil {
		// Let's try to roll back the create key on AWS...
		_ = awsutil.DeleteKey(clusterName, credentials)
		return "", err
	}

	return clusterName, nil
}
