package apiserver

import (
	"fmt"
	pb "github.com/samsung-cnct/cma-aws/pkg/generated/api"
	"github.com/samsung-cnct/cma-aws/pkg/util/awsutil"
	"github.com/samsung-cnct/cma-aws/pkg/util/awsutil/models"
	"github.com/samsung-cnct/cma-aws/pkg/util/cluster"
	"golang.org/x/net/context"
	"strconv"
)

func (s *Server) CreateCluster(ctx context.Context, in *pb.CreateClusterMsg) (*pb.CreateClusterReply, error) {
	// Quick validation
	if in.Provider.GetAws() == nil {
		return nil, fmt.Errorf("AWS Configuration not provided, bailing")
	}
	if len(in.Provider.GetAws().DataCenter.AvailabilityZones) < 1 {
		return nil, fmt.Errorf("Need an availability zone")
	}
	if len(in.Provider.GetAws().InstanceGroups) < 1 {
		return nil, fmt.Errorf("Need an instance group")
	}

	clusterName := in.Name
	credentials := generateCredentials(in.Provider.GetAws().Credentials)

	// Going to create the SSH Key and store it...
	keyName, err := cluster.ProvisionAndSaveSSHKey(clusterName, credentials)
	if err != nil {
		return nil, err
	}

	stackId, err := awsutil.DeployNewVPCHeptioCFStack(awsmodels.NewVPCK8SCFTemplateOptions{
		Name:                 clusterName,
		KeyName:              keyName,
		AvailabilityZone:     in.Provider.GetAws().DataCenter.AvailabilityZones[0],
		AdminIngressLocation: "0.0.0.0/0",
		InstanceType:         in.Provider.GetAws().InstanceGroups[0].Type,
		DiskSize:             "60",
		K8sNodeCapacity:      strconv.Itoa(int(in.Provider.GetAws().InstanceGroups[0].MinQuantity)),
	}, credentials)
	if err != nil {
		// Going to try to roll back the key creation...

		return nil, err
	}

	fmt.Printf("Creating stack %s\n", stackId)
	return &pb.CreateClusterReply{
		Ok: true,
		Cluster: &pb.ClusterItem{
			Id:     clusterName,
			Name:   clusterName,
			Status: pb.ClusterStatus_PROVISIONING,
		},
	}, nil
}

func (s *Server) GetCluster(ctx context.Context, in *pb.GetClusterMsg) (*pb.GetClusterReply, error) {
	stackId := in.Name
	credentials := generateCredentials(in.Credentials)

	outputs, err := awsutil.GetHeptioCFStackOutput(stackId, credentials)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return nil, err
	}
	kubeconfig, err := cluster.GetKubeConfig(stackId, cluster.SSHConnectionOptions{
		BastionHost: cluster.SSHConnectionHost{
			Hostname:      outputs.BastionHostPublicIp,
			Port:          "22",
			Username:      "ubuntu",
			KeySecretName: stackId,
		},
		TargetHost: cluster.SSHConnectionHost{
			Hostname:      outputs.MasterPrivateIp,
			Port:          "22",
			Username:      "ubuntu",
			KeySecretName: stackId,
		},
	})
	return &pb.GetClusterReply{
		Ok: true,
		Cluster: &pb.ClusterDetailItem{
			Id:         stackId,
			Name:       stackId,
			Status:     pb.ClusterStatus_RUNNING,
			Kubeconfig: string(kubeconfig),
		},
	}, nil
}

func (s *Server) DeleteCluster(ctx context.Context, in *pb.DeleteClusterMsg) (*pb.DeleteClusterReply, error) {
	stackId := in.Name
	credentials := generateCredentials(in.Credentials)

	err := awsutil.DeleteCFStack(stackId, credentials)
	if err != nil {
		return nil, err
	}

	// TODO: Should we continue to clean up if the initial thing fails?
	err = cluster.CleanupClusterInK8s(stackId)
	if err != nil {
		return nil, err
	}
	err = awsutil.DeleteKey(stackId, credentials)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteClusterReply{Ok: true, Status: "Deleting"}, nil
}

func (s *Server) GetClusterList(ctx context.Context, in *pb.GetClusterListMsg) (reply *pb.GetClusterListReply, err error) {
	reply = &pb.GetClusterListReply{}
	return
}

func generateCredentials(credentials *pb.AWSCredentials) awsmodels.Credentials {
	return awsmodels.Credentials{
		Region:          credentials.Region,
		AccessKeyId:     credentials.SecretKeyId,
		SecertAccessKey: credentials.SecretAccessKey,
	}
}
