package awsmodels

const (
	EKSVPCCFTemplateURL = "https://amazon-eks.s3-us-west-2.amazonaws.com/cloudformation/2018-08-21/amazon-eks-vpc-sample.yaml"

	EKSVPCCFTemplateNameParameter = "Name"
)

type EKSVPCOptions struct {
	Name string
}
