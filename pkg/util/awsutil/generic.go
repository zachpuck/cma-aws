package awsutil

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"gitlab.com/mvenezia/cma-aws/pkg/util/awsutil/models"
)

func createEC2ServiceFromCredentials(input awsmodels.Credentials) (*ec2.EC2, error) {
	tempSession, err := session.NewSession(&aws.Config{
		Region:      aws.String(input.Region),
		Credentials: credentials.NewStaticCredentials(input.AccessKeyId, input.SecertAccessKey, ""),
	},
	)
	if err != nil {
		return nil, err
	}

	return ec2.New(tempSession), nil
}
