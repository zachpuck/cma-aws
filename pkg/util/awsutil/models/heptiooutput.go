package awsmodels

const (
	HeptioStackOutputBastionHostPublicDNS = "BastionHostPublicDNS"
	HeptioStackOutputVPCID                = "VPCID"
	HeptioStackOutputMasterInstanceId     = "MasterInstanceId"
	HeptioStackOutputBastionHostPublicIp  = "BastionHostPublicIp"
	HeptioStackOutputNodeGroupInstanceId  = "NodeGroupInstanceId"
	HeptioStackOutputMasterPrivateIp      = "MasterPrivateIp"
)

type HeptioStackOutput struct {
	BastionHostPublicDNS string
	VPCID                string
	MasterInstanceId     string
	BastionHostPublicIp  string
	NodeGroupInstanceId  string
	MasterPrivateIp      string
}
