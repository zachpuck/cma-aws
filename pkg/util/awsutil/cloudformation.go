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

func DeployEKSVPCCFStack(cfTemplateOptions awsutilmodels.EKSVPCOptions) (awsStackId string, err error) {
	if CFService == nil {
		err = initializeCF()
		if err != nil {
			return
		}
	}

	cfStackInput := &cloudformation.CreateStackInput{
		TemplateURL: aws.String(awsutilmodels.EKSVPCCFTemplateURL),
		StackName:   aws.String(cfTemplateOptions.Name),
		Parameters:  []*cloudformation.Parameter{},
	}

	output, err := CFService.CreateStack(cfStackInput)
	if err == nil {
		awsStackId = *output.StackId
	}
	return
}

func DeleteCFStack(name string) (err error) {
	if CFService == nil {
		err = initializeCF()
		if err != nil {
			return
		}
	}

	deleteCFStackInput := &cloudformation.DeleteStackInput{StackName: &name}
	_, err = CFService.DeleteStack(deleteCFStackInput)
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
