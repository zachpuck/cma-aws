package awsmodels

type K8SCFTemplateOptions struct {
	Name                 string
	KeyName              string
	VPCId                string
	AvailabilityZone     string
	ClusterSubnetId      string
	LoadBalancerSubnetId string
	SSHLocation          string
	APILBLocation        string
	InstanceType         string
	DiskSize             string
	K8sNodeCapacity      string
}

const (
	K8SCFTemplateNameParameter                 = "ClusterAssociation"
	K8SCFTemplateKeyNameParameter              = "KeyName"
	K8SCFTemplateVPCIdParameter                = "VPCID"
	K8SCFTemplateAvailabilityZoneParameter     = "AvailabilityZone"
	K8SCFTemplateClusterSubnetIdParameter      = "ClusterSubnetId"
	K8SCFTemplateLoadBalancerSubnetIdParameter = "LoadBalancerSubnetId"
	K8SCFTemplateSSHLocationParameter          = "SSHLocation"
	K8SCFTemplateAPILBLocationParameter        = "ApiLbLocation"
	K8SCFTemplateInstanceTypeParameter         = "InstanceType"
	K8SCFTemplateDiskSizeParameter             = "DiskSizeGb"
	K8SCFTemplateKeyK8sNodeCapacityParameter   = "K8sNodeCapacity"

	K8SCFTemplateKeyIAMCapability = "CAPABILITY_IAM"
)

type NewVPCK8SCFTemplateOptions struct {
	Name                 string
	KeyName              string
	AvailabilityZone     string
	AdminIngressLocation string
	InstanceType         string
	DiskSize             string
	K8sNodeCapacity      string
}

const (
	NewVPCK8SCFTemplateKeyNameParameter              = "KeyName"
	NewVPCK8SCFTemplateAvailabilityZoneParameter     = "AvailabilityZone"
	NewVPCK8SCFTemplateAdminIngressLocationParameter = "AdminIngressLocation"
	NewVPCK8SCFTemplateInstanceTypeParameter         = "InstanceType"
	NewVPCK8SCFTemplateDiskSizeParameter             = "DiskSizeGb"
	NewVPCK8SCFTemplateKeyK8sNodeCapacityParameter   = "K8sNodeCapacity"

	NewVPCK8SCFTemplateKeyIAMCapability = "CAPABILITY_IAM"
)

type TestK8SCFTemplateOptions struct {
	Name             string
	KeyName          string
	Username         string
	AvailabilityZone string
	SSHLocation      string
	InstanceType     string
	DiskSize         string
	K8sNodeCapacity  string
}

const (
	TestK8SCFTemplateNameParameter               = "CmsId"
	TestK8SCFTemplateKeyNameParameter            = "KeyName"
	TestK8SCFTemplateUsernameParameter           = "username"
	TestK8SCFTemplateAvailabilityZoneParameter   = "AvailabilityZone"
	TestK8SCFTemplateSSHLocationParameter        = "SSHLocation"
	TestK8SCFTemplateInstanceTypeParameter       = "InstanceType"
	TestK8SCFTemplateDiskSizeParameter           = "DiskSizeGb"
	TestK8SCFTemplateKeyK8sNodeCapacityParameter = "K8sNodeCapacity"

	TestK8SCFTemplateKeyIAMCapability = "CAPABILITY_IAM"
)
