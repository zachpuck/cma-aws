package cluster

import (
	"github.com/spf13/viper"
	"gitlab.com/mvenezia/cma-aws/pkg/util/k8s"
)

func CleanupClusterInK8s(clusterName string) error {
	// Going to delete the SSH Key
	err := k8sutil.DeleteKubeconfigSecret(generateSSHSecretKey(clusterName), viper.GetString("kubernetes-namespace"))
	if err != nil {
		return err
	}
	// And the Kubeconfig
	err = k8sutil.DeleteSSHSecret(generateKubeconfigKey(clusterName), viper.GetString("kubernetes-namespace"))
	if err != nil {
		return err
	}
	return nil
}
