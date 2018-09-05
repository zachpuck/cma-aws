package awsutil

import (
	"github.com/aws/aws-sdk-go/service/ec2"
	"gitlab.com/mvenezia/cma-aws/pkg/util/awsutil/models"
)

func GetInstanceDetails(instanceId string, credentials awsmodels.Credentials) (reservation *ec2.Reservation, err error) {
	instanceIds := []*string{&instanceId}
	results, err := GetInstanceListDetails(instanceIds, credentials)
	if err != nil {
		return
	}
	reservation = results[0]
	return
}

func GetInstanceListDetails(instanceIds []*string, credentials awsmodels.Credentials) (reservations []*ec2.Reservation, err error) {
	service, err := createEC2ServiceFromCredentials(credentials)
	if err != nil {
		return
	}
	instanceDetails, err := service.DescribeInstances(&ec2.DescribeInstancesInput{InstanceIds: instanceIds})
	if err != nil || len(instanceDetails.Reservations) < 1 {
		return
	}

	reservations = instanceDetails.Reservations
	return
}
