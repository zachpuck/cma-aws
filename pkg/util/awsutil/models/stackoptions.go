package awsmodels

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
