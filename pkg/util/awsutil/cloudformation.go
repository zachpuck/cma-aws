package awsutil

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/spf13/viper"
	awsutilmodels "gitlab.com/mvenezia/cma-aws/pkg/util/awsutil/models"
	"time"
)

var CFService *cloudformation.CloudFormation

const (
	CFStackGiveUpThreshold = 10 * time.Minute
	CFStackRetestInterval  = 5 * time.Second
)

func initializeCF() error {
	CFSession, err := session.NewSession(&aws.Config{
		Region:      aws.String(viper.GetString("aws-region")),
		Credentials: credentials.NewStaticCredentials(viper.GetString("aws-access-key-id"), viper.GetString("aws-secret-access-key"), ""),
	},
	)
	if err != nil {
		return err
	}

	CFService = cloudformation.New(CFSession)
	return nil
}

func createCFServiceFromCredentials(input awsutilmodels.Credentials) (*cloudformation.CloudFormation, error) {
	tempSession, err := session.NewSession(&aws.Config{
		Region:      aws.String(input.Region),
		Credentials: credentials.NewStaticCredentials(input.AccessKeyId, input.SecertAccessKey, ""),
	},
	)
	if err != nil {
		return nil, err
	}

	return cloudformation.New(tempSession), nil
}

func GetCFList() (output []*cloudformation.StackSummary, err error) {
	if CFService == nil {
		err = initializeCF()
		if err != nil {
			return
		}
	}
	stacks, err := CFService.ListStacks(&cloudformation.ListStacksInput{})
	output = stacks.StackSummaries
	return
}

func GetCFStack(name string) (resources []*cloudformation.StackResource, parameters []*cloudformation.Parameter, err error) {
	if CFService == nil {
		err = initializeCF()
		if err != nil {
			return
		}
	}

	output, err := CFService.DescribeStackResources(&cloudformation.DescribeStackResourcesInput{StackName: &name})
	if err != nil {
		return
	}
	resources = output.StackResources
	instanceDetails, err := CFService.DescribeStacks(&cloudformation.DescribeStacksInput{StackName: &name})
	if err != nil {
		return
	}
	parameters = instanceDetails.Stacks[0].Parameters
	return
}

func DeployHeptioCFStack(cfTemplateOptions awsutilmodels.K8SCFTemplateOptions) (awsStackId string, err error) {
	if CFService == nil {
		err = initializeCF()
		if err != nil {
			return
		}
	}

	cfStackInput := &cloudformation.CreateStackInput{
		Capabilities: []*string{aws.String(awsutilmodels.K8SCFTemplateKeyIAMCapability)},
		TemplateBody: aws.String(awsutilmodels.K8SCFTemplate),
		StackName:    aws.String(cfTemplateOptions.Name),
		Parameters: []*cloudformation.Parameter{
			{ParameterKey: aws.String(awsutilmodels.K8SCFTemplateNameParameter), ParameterValue: aws.String("")},
			{ParameterKey: aws.String(awsutilmodels.K8SCFTemplateKeyNameParameter), ParameterValue: &cfTemplateOptions.KeyName},
			{ParameterKey: aws.String(awsutilmodels.K8SCFTemplateVPCIdParameter), ParameterValue: &cfTemplateOptions.VPCId},
			{ParameterKey: aws.String(awsutilmodels.K8SCFTemplateAvailabilityZoneParameter), ParameterValue: &cfTemplateOptions.AvailabilityZone},
			{ParameterKey: aws.String(awsutilmodels.K8SCFTemplateClusterSubnetIdParameter), ParameterValue: &cfTemplateOptions.ClusterSubnetId},
			{ParameterKey: aws.String(awsutilmodels.K8SCFTemplateLoadBalancerSubnetIdParameter), ParameterValue: &cfTemplateOptions.LoadBalancerSubnetId},
			{ParameterKey: aws.String(awsutilmodels.K8SCFTemplateSSHLocationParameter), ParameterValue: &cfTemplateOptions.SSHLocation},
			{ParameterKey: aws.String(awsutilmodels.K8SCFTemplateAPILBLocationParameter), ParameterValue: &cfTemplateOptions.APILBLocation},
			{ParameterKey: aws.String(awsutilmodels.K8SCFTemplateInstanceTypeParameter), ParameterValue: &cfTemplateOptions.InstanceType},
			{ParameterKey: aws.String(awsutilmodels.K8SCFTemplateDiskSizeParameter), ParameterValue: &cfTemplateOptions.DiskSize},
			{ParameterKey: aws.String(awsutilmodels.K8SCFTemplateKeyK8sNodeCapacityParameter), ParameterValue: &cfTemplateOptions.K8sNodeCapacity},
		},
	}

	output, err := CFService.CreateStack(cfStackInput)
	if err == nil {
		awsStackId = *output.StackId
	}
	return
}

func DeployNewVPCHeptioCFStack(cfTemplateOptions awsutilmodels.NewVPCK8SCFTemplateOptions, credentials awsutilmodels.Credentials) (awsStackId string, err error) {
	service, err := createCFServiceFromCredentials(credentials)
	if err != nil {
		return
	}
	cfStackInput := &cloudformation.CreateStackInput{
		Capabilities: []*string{aws.String(awsutilmodels.NewVPCK8SCFTemplateKeyIAMCapability)},
		TemplateBody: aws.String(awsutilmodels.NewVPCHeptioCFTemplate),
		StackName:    aws.String(cfTemplateOptions.Name),
		Parameters: []*cloudformation.Parameter{
			{ParameterKey: aws.String(awsutilmodels.NewVPCK8SCFTemplateKeyNameParameter), ParameterValue: &cfTemplateOptions.KeyName},
			{ParameterKey: aws.String(awsutilmodels.NewVPCK8SCFTemplateAvailabilityZoneParameter), ParameterValue: &cfTemplateOptions.AvailabilityZone},
			{ParameterKey: aws.String(awsutilmodels.NewVPCK8SCFTemplateAdminIngressLocationParameter), ParameterValue: &cfTemplateOptions.AdminIngressLocation},
			{ParameterKey: aws.String(awsutilmodels.NewVPCK8SCFTemplateInstanceTypeParameter), ParameterValue: &cfTemplateOptions.InstanceType},
			{ParameterKey: aws.String(awsutilmodels.NewVPCK8SCFTemplateDiskSizeParameter), ParameterValue: &cfTemplateOptions.DiskSize},
			{ParameterKey: aws.String(awsutilmodels.NewVPCK8SCFTemplateKeyK8sNodeCapacityParameter), ParameterValue: &cfTemplateOptions.K8sNodeCapacity},
		},
	}

	output, err := service.CreateStack(cfStackInput)
	if err == nil {
		awsStackId = *output.StackId
	}
	return
}

func DeployJasonCFStack(cfTemplateOptions awsutilmodels.TestK8SCFTemplateOptions) (awsStackId string, err error) {
	if CFService == nil {
		err = initializeCF()
		if err != nil {
			return
		}
	}

	cfStackInput := &cloudformation.CreateStackInput{
		Capabilities: []*string{aws.String(awsutilmodels.TestK8SCFTemplateKeyIAMCapability)},
		TemplateBody: aws.String(awsutilmodels.TestK8SCFTemplate),
		StackName:    aws.String(cfTemplateOptions.Name),
		Parameters: []*cloudformation.Parameter{
			{ParameterKey: aws.String(awsutilmodels.TestK8SCFTemplateNameParameter), ParameterValue: &cfTemplateOptions.Name},
			{ParameterKey: aws.String(awsutilmodels.TestK8SCFTemplateKeyNameParameter), ParameterValue: &cfTemplateOptions.KeyName},
			{ParameterKey: aws.String(awsutilmodels.TestK8SCFTemplateUsernameParameter), ParameterValue: &cfTemplateOptions.Username},
			{ParameterKey: aws.String(awsutilmodels.TestK8SCFTemplateAvailabilityZoneParameter), ParameterValue: &cfTemplateOptions.AvailabilityZone},
			{ParameterKey: aws.String(awsutilmodels.TestK8SCFTemplateSSHLocationParameter), ParameterValue: &cfTemplateOptions.SSHLocation},
			{ParameterKey: aws.String(awsutilmodels.TestK8SCFTemplateInstanceTypeParameter), ParameterValue: &cfTemplateOptions.InstanceType},
			{ParameterKey: aws.String(awsutilmodels.TestK8SCFTemplateDiskSizeParameter), ParameterValue: &cfTemplateOptions.DiskSize},
			{ParameterKey: aws.String(awsutilmodels.TestK8SCFTemplateKeyK8sNodeCapacityParameter), ParameterValue: &cfTemplateOptions.K8sNodeCapacity},
		},
	}

	output, err := CFService.CreateStack(cfStackInput)
	if err == nil {
		awsStackId = *output.StackId
	}
	return
}

func DeleteCFStack(name string, credentials awsutilmodels.Credentials) (err error) {
	service, err := createCFServiceFromCredentials(credentials)
	if err != nil {
		return
	}

	deleteCFStackInput := &cloudformation.DeleteStackInput{StackName: &name}
	_, err = service.DeleteStack(deleteCFStackInput)
	return
}

func IsStackReady(resources []*cloudformation.StackResource) (ready bool, masterReady bool, workerReady bool) {
	for _, item := range resources {
		if *item.LogicalResourceId == "K8sMaster" && *item.ResourceStatus == cloudformation.ChangeSetStatusCreateComplete {
			masterReady = true
			continue
		}
		if *item.LogicalResourceId == "K8sNodeGroup" && *item.ResourceStatus == cloudformation.ChangeSetStatusCreateComplete {
			workerReady = true
		}
	}

	if workerReady && masterReady {
		ready = true
	}
	return
}

func WaitUntilStackIsReady(stackName string, printProgress bool) error {
	timePassed := 0 * time.Second
	complete := false
	for complete == false {
		if timePassed > CFStackGiveUpThreshold {
			if printProgress {
				fmt.Printf("Gave up waiting, sorry\n")
			}
			return fmt.Errorf("gave up waiting, for stack %s, sorry", stackName)
			break
		}
		resources, _, err := GetCFStack(stackName)
		if err == nil {
			complete, masterOk, workerOK := IsStackReady(resources)
			if complete == true {
				if printProgress {
					fmt.Printf("Cluster is ready!\n")
				}
				return nil
			}
			if printProgress {
				fmt.Printf("Cluster is not ready; ")
				if masterOk {
					fmt.Printf("Master is ready; ")
				} else {
					fmt.Printf("Master is NOT ready; ")
				}
				if workerOK {
					fmt.Printf("Worker Nodes are ready ")
				} else {
					fmt.Printf("Worker Nodes are NOT ready ")
				}
				fmt.Printf("\n")
			}
		} else {
			if printProgress {
				fmt.Printf("Error with stack: %s\n", err)
			}
		}
		time.Sleep(CFStackRetestInterval)
		timePassed += CFStackRetestInterval
	}

	return nil
}

func GetHeptioCFStackOutput(name string, credentials awsutilmodels.Credentials) (output awsutilmodels.HeptioStackOutput, err error) {
	service, err := createCFServiceFromCredentials(credentials)
	if err != nil {
		return
	}

	query, err := service.DescribeStacks(&cloudformation.DescribeStacksInput{StackName: &name})
	if err != nil {
		return
	}
	if len(query.Stacks) < 1 {
		err = fmt.Errorf("no stack found with name -->%s<-- ", name)
		return
	}
	stack := query.Stacks[0]
	if *stack.StackStatus != cloudformation.StackStatusCreateComplete {
		err = fmt.Errorf("stack -->%s<-- is not complete, rather it has status %s", name, *stack.StackStatus)
		return
	}

	// OK, we should feel safe to grab outputs
	output = createHeptioCFStackOutput(stack.Outputs)
	return
}

func createHeptioCFStackOutput(outputs []*cloudformation.Output) awsutilmodels.HeptioStackOutput {
	output := awsutilmodels.HeptioStackOutput{}
	for _, j := range outputs {
		switch *j.OutputKey {
		case awsutilmodels.HeptioStackOutputBastionHostPublicDNS:
			output.BastionHostPublicDNS = *j.OutputValue
			break
		case awsutilmodels.HeptioStackOutputVPCID:
			output.VPCID = *j.OutputValue
			break
		case awsutilmodels.HeptioStackOutputMasterInstanceId:
			output.MasterInstanceId = *j.OutputValue
			break
		case awsutilmodels.HeptioStackOutputBastionHostPublicIp:
			output.BastionHostPublicIp = *j.OutputValue
			break
		case awsutilmodels.HeptioStackOutputNodeGroupInstanceId:
			output.NodeGroupInstanceId = *j.OutputValue
			break
		case awsutilmodels.HeptioStackOutputMasterPrivateIp:
			output.MasterPrivateIp = *j.OutputValue
			break
		}
	}
	return output
}
