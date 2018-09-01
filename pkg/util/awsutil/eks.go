package awsutil

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/spf13/viper"
)

const (
	EKSCFSecurityGroupResourceName = "ControlPlaneSecurityGroup"
	EKSCFSubnet1ResourceName       = "Subnet01"
	EKSCFSubnet2ResourceName       = "Subnet02"
	EKSCFSubnet3ResourceName       = "Subnet03"
)

var (
	EKSService *eks.EKS
)

type EKSMasterOptions struct {
	Name          string
	RoleARN       string
	SecurityGroup string
	Subnets       []string
}

func initializeEKS() error {
	RoleSession, err := session.NewSession(&aws.Config{
		Region:      aws.String(viper.GetString("aws-region")),
		Credentials: credentials.NewStaticCredentials(viper.GetString("aws-access-key-id"), viper.GetString("aws-secret-access-key"), ""),
	},
	)
	if err != nil {
		return err
	}

	EKSService = eks.New(RoleSession)
	return nil
}

func CreateEKSControlPlaneForStack(stackName string) (cluster *eks.Cluster, err error) {
	resources, _, err := GetCFStack(stackName)
	if err != nil {
		return
	}

	// We need to find some specific things from the stack - let's go find it!
	securityGroupId, err := FindEKSSecurityGroupInStack(resources)
	if err != nil {
		return
	}
	subnets, err := FindEKSSubnetsInStack(resources)
	if err != nil {
		return
	}

	//// OK, so let's create the AWS Role for EKS...
	//role, err := GenerateEKSRole(stackName)
	//if err != nil {
	//	return
	//}

	cluster, err = CreateEKSControlPlane(EKSMasterOptions{
		Name:          stackName,
		SecurityGroup: securityGroupId,
		Subnets:       subnets,
		RoleARN:       "arn:aws:iam::297062999199:role/dumbo3",
	})
	return
}

func FindEKSSecurityGroupInStack(resources []*cloudformation.StackResource) (securityGroupId string, err error) {
	for _, j := range resources {
		if *j.LogicalResourceId == EKSCFSecurityGroupResourceName {
			securityGroupId = *j.PhysicalResourceId
			return
		}
	}
	err = fmt.Errorf("could not find security group for stack")
	return
}

func FindEKSSubnetsInStack(resources []*cloudformation.StackResource) (subnets []string, err error) {
	var subnet1, subnet2, subnet3 string
	for _, j := range resources {
		switch *j.LogicalResourceId {
		case EKSCFSubnet1ResourceName:
			subnet1 = *j.PhysicalResourceId
			break
		case EKSCFSubnet2ResourceName:
			subnet2 = *j.PhysicalResourceId
			break
		case EKSCFSubnet3ResourceName:
			subnet3 = *j.PhysicalResourceId
			break
		}
	}
	if subnet1 == "" || subnet2 == "" || subnet3 == "" {
		err = fmt.Errorf("could not find all subnets in the stack")
		return
	}
	subnets = []string{subnet1, subnet2, subnet3}
	return
}

func CreateEKSControlPlane(options EKSMasterOptions) (cluster *eks.Cluster, err error) {
	if EKSService == nil {
		err = initializeEKS()
		if err != nil {
			return
		}
	}
	eksResult, err := EKSService.CreateCluster(
		&eks.CreateClusterInput{
			Name:    aws.String(options.Name),
			RoleArn: aws.String(options.RoleARN),
			ResourcesVpcConfig: &eks.VpcConfigRequest{
				SecurityGroupIds: aws.StringSlice([]string{options.SecurityGroup}),
				SubnetIds:        aws.StringSlice(options.Subnets),
			},
		})
	if err != nil {
		return
	}
	cluster = eksResult.Cluster
	return
}
