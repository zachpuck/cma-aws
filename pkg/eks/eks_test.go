package eks

import "testing"

func TestCreate(t *testing.T) {
	ekscluster := New(AwsCredentials{
		AccessKeyId:     "secret",
		SecretAccessKey: "secret",
		Region:          "us-east-1",
	}, "test1", "us-east-1", "")

	nodepools := make([]NodePool, 1)
	nodepools[0].Name = "nodepool1"
	nodepools[0].Nodes = 2
	nodepools[0].Type = "m5.large"

	err := ekscluster.CreateCluster(CreateClusterInput{
		Name:              ekscluster.clusterName,
		Version:           "1.11",
		Region:            ekscluster.region,
		AvailabilityZones: nil,
		NodePools:         nodepools,
	})
	if err != nil {
		t.Errorf("ekscluster.Create() = %s; want nil", err)
	}
}
