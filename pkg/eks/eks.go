package eks

import "fmt"

type EKS struct {
	awsCreds       AwsCredentials
	clusterName    string
	region         string
	kubeConfigPath string
}

func New(creds AwsCredentials, name string, region string, kubeConfigPath string) *EKS {
	return &EKS{
		awsCreds:       creds,
		clusterName:    name,
		region:         region,
		kubeConfigPath: kubeConfigPath,
	}
}

func (e *EKS) CreateCluster(in CreateClusterInput) error {
	fmt.Printf("Create input %+v\n", in)

	return nil
}
