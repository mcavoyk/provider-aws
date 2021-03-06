package ec2

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/aws/aws-sdk-go-v2/service/ec2"

	"github.com/crossplane/provider-aws/apis/ec2/v1beta1"
	awsclients "github.com/crossplane/provider-aws/pkg/clients"
)

const (
	// VPCIDNotFound is the code that is returned by ec2 when the given VPCID is not valid
	VPCIDNotFound = "InvalidVpcID.NotFound"
)

// VPCClient is the external client used for VPC Custom Resource
type VPCClient interface {
	CreateVpcRequest(*ec2.CreateVpcInput) ec2.CreateVpcRequest
	DeleteVpcRequest(*ec2.DeleteVpcInput) ec2.DeleteVpcRequest
	DescribeVpcsRequest(*ec2.DescribeVpcsInput) ec2.DescribeVpcsRequest
	DescribeVpcAttributeRequest(*ec2.DescribeVpcAttributeInput) ec2.DescribeVpcAttributeRequest
	ModifyVpcAttributeRequest(*ec2.ModifyVpcAttributeInput) ec2.ModifyVpcAttributeRequest
	CreateTagsRequest(*ec2.CreateTagsInput) ec2.CreateTagsRequest
	ModifyVpcTenancyRequest(*ec2.ModifyVpcTenancyInput) ec2.ModifyVpcTenancyRequest
}

// NewVPCClient returns a new client using AWS credentials as JSON encoded data.
func NewVPCClient(cfg aws.Config) VPCClient {
	return ec2.New(cfg)
}

// IsVPCNotFoundErr returns true if the error is because the item doesn't exist
func IsVPCNotFoundErr(err error) bool {
	if awsErr, ok := err.(awserr.Error); ok {
		if awsErr.Code() == VPCIDNotFound {
			return true
		}
	}

	return false
}

// IsVpcUpToDate returns true if there is no update-able difference between desired
// and observed state of the resource.
func IsVpcUpToDate(spec v1beta1.VPCParameters, vpc ec2.Vpc, attributes ec2.DescribeVpcAttributeOutput) bool {
	if aws.StringValue(spec.InstanceTenancy) != string(vpc.InstanceTenancy) {
		return false
	}

	if aws.BoolValue(spec.EnableDNSHostNames) != aws.BoolValue(attributes.EnableDnsHostnames.Value) ||
		aws.BoolValue(spec.EnableDNSSupport) != aws.BoolValue(attributes.EnableDnsSupport.Value) {
		return false
	}

	return v1beta1.CompareTags(spec.Tags, vpc.Tags)
}

// GenerateVpcObservation is used to produce v1beta1.VPCObservation from
// ec2.Vpc.
func GenerateVpcObservation(vpc ec2.Vpc) v1beta1.VPCObservation {
	o := v1beta1.VPCObservation{
		IsDefault:     aws.BoolValue(vpc.IsDefault),
		DHCPOptionsID: aws.StringValue(vpc.DhcpOptionsId),
		OwnerID:       aws.StringValue(vpc.OwnerId),
		VPCState:      string(vpc.State),
	}

	if len(vpc.CidrBlockAssociationSet) > 0 {
		o.CIDRBlockAssociationSet = make([]v1beta1.VPCCIDRBlockAssociation, len(vpc.CidrBlockAssociationSet))
		for i, v := range vpc.CidrBlockAssociationSet {
			o.CIDRBlockAssociationSet[i] = v1beta1.VPCCIDRBlockAssociation{
				AssociationID: aws.StringValue(v.AssociationId),
				CIDRBlock:     aws.StringValue(v.CidrBlock),
			}
			o.CIDRBlockAssociationSet[i].CIDRBlockState = v1beta1.VPCCIDRBlockState{
				State:         string(v.CidrBlockState.State),
				StatusMessage: aws.StringValue(v.CidrBlockState.StatusMessage),
			}
		}
	}

	if len(vpc.Ipv6CidrBlockAssociationSet) > 0 {
		o.IPv6CIDRBlockAssociationSet = make([]v1beta1.VPCIPv6CidrBlockAssociation, len(vpc.Ipv6CidrBlockAssociationSet))
		for i, v := range vpc.Ipv6CidrBlockAssociationSet {
			o.IPv6CIDRBlockAssociationSet[i] = v1beta1.VPCIPv6CidrBlockAssociation{
				AssociationID:      aws.StringValue(v.AssociationId),
				IPv6CIDRBlock:      aws.StringValue(v.Ipv6CidrBlock),
				IPv6Pool:           aws.StringValue(v.Ipv6Pool),
				NetworkBorderGroup: aws.StringValue(v.NetworkBorderGroup),
			}
			o.IPv6CIDRBlockAssociationSet[i].IPv6CIDRBlockState = v1beta1.VPCCIDRBlockState{
				State:         string(v.Ipv6CidrBlockState.State),
				StatusMessage: aws.StringValue(v.Ipv6CidrBlockState.StatusMessage),
			}
		}
	}

	return o
}

// LateInitializeVPC fills the empty fields in *v1beta1.VPCParameters with
// the values seen in ec2.Vpc and ec2.DescribeVpcAttributeOutput.
func LateInitializeVPC(in *v1beta1.VPCParameters, v *ec2.Vpc, attributes *ec2.DescribeVpcAttributeOutput) { // nolint:gocyclo
	if v == nil {
		return
	}

	in.CIDRBlock = awsclients.LateInitializeString(in.CIDRBlock, v.CidrBlock)
	in.InstanceTenancy = awsclients.LateInitializeStringPtr(in.InstanceTenancy, aws.String(string(v.InstanceTenancy)))
	if attributes.EnableDnsHostnames != nil {
		in.EnableDNSHostNames = awsclients.LateInitializeBoolPtr(in.EnableDNSHostNames, attributes.EnableDnsHostnames.Value)
	}
	if attributes.EnableDnsHostnames != nil {
		in.EnableDNSSupport = awsclients.LateInitializeBoolPtr(in.EnableDNSSupport, attributes.EnableDnsSupport.Value)
	}
}
