package ec2

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/aws/aws-sdk-go-v2/service/ec2"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/crossplane/provider-aws/apis/ec2/v1alpha1"
	"github.com/crossplane/provider-aws/apis/ec2/v1beta1"
	awsclients "github.com/crossplane/provider-aws/pkg/clients"
)

// NatGatewayClient is the external client used for NatGateway Custom Resource
type NatGatewayClient interface {
	CreateNatGatewayRequest(input *ec2.CreateNatGatewayInput) ec2.CreateNatGatewayRequest
	DeleteNatGatewayRequest(input *ec2.DeleteNatGatewayInput) ec2.DeleteNatGatewayRequest
	DescribeNatGatewaysRequest(input *ec2.DescribeNatGatewaysInput) ec2.DescribeNatGatewaysRequest
	CreateTagsRequest(input *ec2.CreateTagsInput) ec2.CreateTagsRequest
}

// NewNatGatewayClient returns a new client using AWS credentials as JSON encoded data.
func NewNatGatewayClient(ctx context.Context, credentials []byte, region string, auth awsclients.AuthMethod) (NatGatewayClient, error) {
	cfg, err := auth(ctx, credentials, awsclients.DefaultSection, region)
	if cfg == nil {
		return nil, err
	}
	return ec2.New(*cfg), nil
}

// IsNatGatewayNotFoundErr returns true if the error is because the item doesn't exist
func IsNatGatewayNotFoundErr(err error) bool {
	if awsErr, ok := err.(awserr.Error); ok {
		if awsErr.Code() == VPCIDNotFound {
			return true
		}
	}
	return false
}

// IsNatGatewayUpToDate returns true if there is no update-able difference between desired
// and observed state of the resource.
func IsNatGatewayUpToDate(spec v1alpha1.NatGatewayParameters, ng ec2.NatGateway) bool {
	return v1beta1.CompareTags(spec.Tags, ng.Tags)
}

// GenerateNatGatewayObservation is used to produce v1alpha1.NatGatewayObservation from
// ec2.NatGateway.
func GenerateNatGatewayObservation(ng ec2.NatGateway) v1alpha1.NatGatewayObservation {
	o := v1alpha1.NatGatewayObservation{
		FailureCode:         aws.StringValue(ng.FailureCode),
		FailureMessage:      aws.StringValue(ng.FailureMessage),
		NatGatewayAddresses: make([]v1alpha1.NatGatewayAddress, len(ng.NatGatewayAddresses)),
		NatGatewayID:        aws.StringValue(ng.NatGatewayId),
		State:               v1alpha1.NatGatewayState(ng.State),
		SubnetID:            aws.StringValue(ng.SubnetId),
		VpcID:               aws.StringValue(ng.VpcId),
	}

	if ng.CreateTime != nil {
		o.CreateTime = &metav1.Time{Time: *ng.CreateTime}
	}

	if ng.DeleteTime != nil {
		o.DeleteTime = &metav1.Time{Time: *ng.DeleteTime}
	}

	for i, addr := range ng.NatGatewayAddresses {
		o.NatGatewayAddresses[i] = v1alpha1.NatGatewayAddress{
			AllocationID:       aws.StringValue(addr.AllocationId),
			NetworkInterfaceID: aws.StringValue(addr.NetworkInterfaceId),
			PrivateIP:          aws.StringValue(addr.PrivateIp),
			PublicIP:           aws.StringValue(addr.PublicIp),
		}
	}
	return o
}

// LateInitializeNatGateway fills the empty fields in *v1alpha1.NatGatewayParameters with
// the values seen in ec2.NatGateway.
func LateInitializeNatGateway(in *v1alpha1.NatGatewayParameters, ng *ec2.NatGateway) {
	if ng == nil {
		return
	}
	in.SubnetID = awsclients.LateInitializeStringPtr(in.SubnetID, ng.SubnetId)
}
