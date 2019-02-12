package awsutil

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	awsutilmodels "github.com/samsung-cnct/cma-aws/pkg/util/awsutil/models"
	"time"
)

const (
	CFStackGiveUpThreshold = 10 * time.Minute
	CFStackRetestInterval  = 5 * time.Second
)

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

func GetCFList(credentials awsutilmodels.Credentials) (output []*cloudformation.StackSummary, err error) {
	service, err := createCFServiceFromCredentials(credentials)
	if err != nil {
		return
	}
	stacks, err := service.ListStacks(&cloudformation.ListStacksInput{})
	output = stacks.StackSummaries
	return
}

func GetCFStack(name string, credentials awsutilmodels.Credentials) (resources []*cloudformation.StackResource, parameters []*cloudformation.Parameter, err error) {
	service, err := createCFServiceFromCredentials(credentials)
	if err != nil {
		return
	}

	output, err := service.DescribeStackResources(&cloudformation.DescribeStackResourcesInput{StackName: &name})
	if err != nil {
		return
	}
	resources = output.StackResources
	instanceDetails, err := service.DescribeStacks(&cloudformation.DescribeStacksInput{StackName: &name})
	if err != nil {
		return
	}
	parameters = instanceDetails.Stacks[0].Parameters
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

func WaitUntilStackIsReady(stackName string, printProgress bool, credentials awsutilmodels.Credentials) error {
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
		resources, _, err := GetCFStack(stackName, credentials)
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

func GetHeptioCFStackOutput(name string, credentials awsutilmodels.Credentials) (output awsutilmodels.ClusterOutput, err error) {
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
		output = awsutilmodels.ClusterOutput{
			Status: *stack.StackStatus,
		}
		return
	}

	// OK, we should feel safe to grab outputs
	output = createHeptioCFStackOutput(stack.Outputs)
	output.Status = cloudformation.StackStatusCreateComplete
	return
}

func createHeptioCFStackOutput(outputs []*cloudformation.Output) awsutilmodels.ClusterOutput {
	output := awsutilmodels.ClusterOutput{}
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
