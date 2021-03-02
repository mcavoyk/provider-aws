/*
Copyright 2021 The Crossplane Authors.

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

// Code generated by ack-generate. DO NOT EDIT.

package v1alpha1

import (
	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// AuthorizerParameters defines the desired state of Authorizer
type AuthorizerParameters struct {
	// Region is which region the Authorizer will be created.
	// +kubebuilder:validation:Required
	Region string `json:"region"`

	AuthorizerCredentialsARN *string `json:"authorizerCredentialsARN,omitempty"`

	AuthorizerPayloadFormatVersion *string `json:"authorizerPayloadFormatVersion,omitempty"`

	AuthorizerResultTtlInSeconds *int64 `json:"authorizerResultTtlInSeconds,omitempty"`

	// +kubebuilder:validation:Required
	AuthorizerType *string `json:"authorizerType"`

	AuthorizerURI *string `json:"authorizerURI,omitempty"`

	EnableSimpleResponses *bool `json:"enableSimpleResponses,omitempty"`

	// +kubebuilder:validation:Required
	IDentitySource []*string `json:"identitySource"`

	IDentityValidationExpression *string `json:"identityValidationExpression,omitempty"`

	JWTConfiguration *JWTConfiguration `json:"jwtConfiguration,omitempty"`

	// +kubebuilder:validation:Required
	Name                       *string `json:"name"`
	CustomAuthorizerParameters `json:",inline"`
}

// AuthorizerSpec defines the desired state of Authorizer
type AuthorizerSpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       AuthorizerParameters `json:"forProvider"`
}

// AuthorizerObservation defines the observed state of Authorizer
type AuthorizerObservation struct {
	AuthorizerID *string `json:"authorizerID,omitempty"`
}

// AuthorizerStatus defines the observed state of Authorizer.
type AuthorizerStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	AtProvider          AuthorizerObservation `json:"atProvider"`
}

// +kubebuilder:object:root=true

// Authorizer is the Schema for the Authorizers API
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="EXTERNAL-NAME",type="string",JSONPath=".metadata.annotations.crossplane\\.io/external-name"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,managed,aws}
type Authorizer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              AuthorizerSpec   `json:"spec,omitempty"`
	Status            AuthorizerStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// AuthorizerList contains a list of Authorizers
type AuthorizerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Authorizer `json:"items"`
}

// Repository type metadata.
var (
	AuthorizerKind             = "Authorizer"
	AuthorizerGroupKind        = schema.GroupKind{Group: Group, Kind: AuthorizerKind}.String()
	AuthorizerKindAPIVersion   = AuthorizerKind + "." + GroupVersion.String()
	AuthorizerGroupVersionKind = GroupVersion.WithKind(AuthorizerKind)
)

func init() {
	SchemeBuilder.Register(&Authorizer{}, &AuthorizerList{})
}