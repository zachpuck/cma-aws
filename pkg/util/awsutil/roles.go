package awsutil

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/spf13/viper"
)

var IAMService *iam.IAM

const (
	AmazonEKSClusterPolicy     = "arn:aws:iam::aws:policy/AmazonEKSClusterPolicy"
	AmazonEKSServicePolicy     = "arn:aws:iam::aws:policy/AmazonEKSServicePolicy"
	AmazonEKSDefaultAssumeRole = `{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "eks.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}`
)

func initializeIAM() error {
	RoleSession, err := session.NewSession(&aws.Config{
		Region:      aws.String(viper.GetString("aws-region")),
		Credentials: credentials.NewStaticCredentials(viper.GetString("aws-access-key-id"), viper.GetString("aws-secret-access-key"), ""),
	},
	)
	if err != nil {
		return err
	}

	IAMService = iam.New(RoleSession)
	return nil
}

func GenerateEKSRole(name string) (role *iam.Role, err error) {
	if IAMService == nil {
		err = initializeIAM()
		if err != nil {
			return
		}
	}

	roleResult, err := IAMService.CreateRole(&iam.CreateRoleInput{AssumeRolePolicyDocument: aws.String(AmazonEKSDefaultAssumeRole), RoleName: aws.String(name)})
	if err != nil {
		return
	}
	role = roleResult.Role
	roleName := roleResult.Role.RoleName

	_, err = IAMService.AttachRolePolicy(&iam.AttachRolePolicyInput{PolicyArn: aws.String(AmazonEKSClusterPolicy), RoleName: roleName})
	if err != nil {
		return
	}
	_, err = IAMService.AttachRolePolicy(&iam.AttachRolePolicyInput{PolicyArn: aws.String(AmazonEKSServicePolicy), RoleName: roleName})
	return

}
