package awsutil

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/spf13/viper"
)

var EC2Session *session.Session
var EC2Service *ec2.EC2

func initializeAWS() error {
	EC2Session, err := session.NewSession(&aws.Config{
		Region:      aws.String(viper.GetString("aws-region")),
		Credentials: credentials.NewStaticCredentials(viper.GetString("aws-access-key-id"), viper.GetString("aws-secret-access-key"), ""),
	},
	)
	if err != nil {
		return err
	}

	EC2Service = ec2.New(EC2Session)
	return nil
}
