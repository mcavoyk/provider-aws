/*
Copyright 2019 The Crossplane Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	runtimev1alpha1 "github.com/crossplane/crossplane-runtime/apis/core/v1alpha1"

	ec2v1beta1 "github.com/crossplane/provider-aws/apis/ec2/v1beta1"
)

// NatGatewayState is the state of a Nat Gateway.
type NatGatewayState string

// Nat Gateway states.
const (
	NatGatewayStatePending   NatGatewayState = "pending"
	NatGatewayStateFailed    NatGatewayState = "failed"
	NatGatewayStateAvailable NatGatewayState = "available"
	NatGatewayStateDeleting  NatGatewayState = "deleting"
	NatGatewayStateDeleted   NatGatewayState = "deleted"
)

// NatGatewayParameters define the desired state of an AWS VPC Nat
// Gateway.
type NatGatewayParameters struct {
	// The allocation ID of an Elastic IP address to associate with the NAT gateway.
	// If the Elastic IP address is associated with another resource, you must first
	// disassociate it.
	// +immutable
	AllocationID *string `json:"allocationId,omitempty"`

	// TODO(mcavoyk): Add AllocationIdRef/AllocationIdSelector when EIP's are a managed resource

	// The subnet in which to create the NAT gateway.
	// +optional
	// +immutable
	SubnetID *string `json:"subnetId,omitempty"`

	// SubnetIDRef references a subnet and retrieves its ID.
	// +optional
	// +immutable
	SubnetIDRef *runtimev1alpha1.Reference `json:"subnetIdRef,omitempty"`

	// SubnetIDSelector selects a reference to a subnet and retrieves its ID.
	// +optional
	SubnetIDSelector *runtimev1alpha1.Selector `json:"subnetIdSelector,omitempty"`

	// Tags represents to current ec2 tags.
	// +optional
	Tags []ec2v1beta1.Tag `json:"tags,omitempty"`
}

// A NatGatewaySpec defines the desired state of a NatGateway.
type NatGatewaySpec struct {
	runtimev1alpha1.ResourceSpec `json:",inline"`
	ForProvider                  NatGatewayParameters `json:"forProvider"`
}

// NatGatewayAddress describes the IP addresses and network interface associated with a NAT gateway.
type NatGatewayAddress struct {
	// The allocation ID of the Elastic IP address that's associated with the NAT
	// gateway.
	AllocationID string `json:"allocationId,omitempty"`

	// The ID of the network interface associated with the NAT gateway.
	NetworkInterfaceID string `json:"networkInterfaceId,omitempty"`

	// The private IP address associated with the Elastic IP address.
	PrivateIP string `json:"privateIp,omitempty"`

	// The Elastic IP address associated with the NAT gateway.
	PublicIP string `json:"publicIp,omitempty"`
}

// NatGatewayObservation keeps the state for the external resource
type NatGatewayObservation struct {
	// The date and time the NAT gateway was created.
	CreateTime *metav1.Time `json:"createTime,omitempty"`

	// The date and time the NAT gateway was deleted, if applicable.
	DeleteTime *metav1.Time `json:"deleteTime,omitempty"`

	// If the NAT gateway could not be created, specifies the error code for the
	// failure. (InsufficientFreeAddressesInSubnet | Gateway.NotAttached | InvalidAllocationID.NotFound
	// | Resource.AlreadyAssociated | InternalError | InvalidSubnetID.NotFound)
	FailureCode string `json:"failureCode,omitempty"`

	// If the NAT gateway could not be created, specifies the error message for
	// the failure, that corresponds to the error code.
	//
	//    * For InsufficientFreeAddressesInSubnet: "Subnet has insufficient free
	//    addresses to create this NAT gateway"
	//
	//    * For Gateway.NotAttached: "Network vpc-xxxxxxxx has no Internet gateway
	//    attached"
	//
	//    * For InvalidAllocationID.NotFound: "Elastic IP address eipalloc-xxxxxxxx
	//    could not be associated with this NAT gateway"
	//
	//    * For Resource.AlreadyAssociated: "Elastic IP address eipalloc-xxxxxxxx
	//    is already associated"
	//
	//    * For InternalError: "Network interface eni-xxxxxxxx, created and used
	//    internally by this NAT gateway is in an invalid state. Please try again."
	//
	//    * For InvalidSubnetID.NotFound: "The specified subnet subnet-xxxxxxxx
	//    does not exist or could not be found."
	FailureMessage string `json:"failureMessage,omitempty"`

	// Information about the IP addresses and network interface associated with
	// the NAT gateway.
	NatGatewayAddresses []NatGatewayAddress `json:"natGatewayAddresses,omitempty"`

	// The ID of the NAT gateway.
	NatGatewayID string `json:"natGatewayId,omitempty"`

	// The state of the NAT gateway.
	//
	//    * pending: The NAT gateway is being created and is not ready to process
	//    traffic.
	//
	//    * failed: The NAT gateway could not be created. Check the failureCode
	//    and failureMessage fields for the reason.
	//
	//    * available: The NAT gateway is able to process traffic. This status remains
	//    until you delete the NAT gateway, and does not indicate the health of
	//    the NAT gateway.
	//
	//    * deleting: The NAT gateway is in the process of being terminated and
	//    may still be processing traffic.
	//
	//    * deleted: The NAT gateway has been terminated and is no longer processing
	//    traffic.
	State NatGatewayState `json:"state,omitempty"`

	// The ID of the subnet in which the NAT gateway is located.
	SubnetID string `json:"subnetId,omitempty"`

	// The ID of the VPC in which the NAT gateway is located.
	VpcID string `json:"vpcId,omitempty"`
}

// A NatGatewayStatus represents the observed state of a NatGateway.
type NatGatewayStatus struct {
	runtimev1alpha1.ResourceStatus `json:",inline"`
	AtProvider                     NatGatewayObservation `json:"atProvider"`
}

// +kubebuilder:object:root=true

// A NatGateway is a managed resource that represents an AWS VPC Nat
// Gateway.
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="ID",type="string",JSONPath=".metadata.annotations.crossplane\\.io/external-name"
// +kubebuilder:printcolumn:name="SUBNET",type="string",JSONPath=".spec.forProvider.subnetId"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,managed,aws}
type NatGateway struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NatGatewaySpec   `json:"spec"`
	Status NatGatewayStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// NatGatewayList contains a list of NatGateway
type NatGatewayList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NatGateway `json:"items"`
}
