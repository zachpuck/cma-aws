package awsutil

import "github.com/aws/aws-sdk-go/service/ec2"

func GetInstanceDetails(instanceId string) (reservation *ec2.Reservation, err error) {
	instanceIds := []*string{&instanceId}
	results, err := GetInstanceListDetails(instanceIds)
	if err != nil {
		return
	}
	reservation = results[0]
	return
}

func GetInstanceListDetails(instanceIds []*string) (reservations []*ec2.Reservation, err error) {
	if EC2Service == nil {
		err = initializeAWS()
		if err != nil {
			return
		}
	}
	instanceDetails, err := EC2Service.DescribeInstances(&ec2.DescribeInstancesInput{InstanceIds: instanceIds})
	if err != nil || len(instanceDetails.Reservations) < 1 {
		return
	}

	reservations = instanceDetails.Reservations
	return
}
