package awsutil

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/spf13/viper"
)

var AutoScalingService *autoscaling.AutoScaling

func initializeAutoScaling() error {
	AutoScalingSession, err := session.NewSession(&aws.Config{
		Region:      aws.String(viper.GetString("aws-region")),
		Credentials: credentials.NewStaticCredentials(viper.GetString("aws-access-key-id"), viper.GetString("aws-secret-access-key"), ""),
	},
	)
	if err != nil {
		return err
	}

	AutoScalingService = autoscaling.New(AutoScalingSession)
	return nil
}

func GetInstancesForASGDetails(asgId string) (reservation *ec2.Reservation, err error) {
	if AutoScalingService == nil {
		err = initializeAutoScaling()
		if err != nil {
			return
		}
	}
	asgDetails, err := AutoScalingService.DescribeAutoScalingGroups(&autoscaling.DescribeAutoScalingGroupsInput{AutoScalingGroupNames: []*string{&asgId}})
	if err != nil || len(asgDetails.AutoScalingGroups) < 1 {
		return
	}
	var instanceIds []*string
	for _, j := range asgDetails.AutoScalingGroups[0].Instances {
		instanceIds = append(instanceIds, j.InstanceId)
	}

	results, err := GetInstanceListDetails(instanceIds)
	if err != nil || len(results) < 1 {
		return
	}
	reservation = results[0]
	return
}
