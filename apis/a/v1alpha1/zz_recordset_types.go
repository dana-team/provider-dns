// SPDX-FileCopyrightText: 2023 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

/*
Copyright 2022 Upbound Inc.
*/

// Code generated by upjet. DO NOT EDIT.

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	v1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
)

type RecordSetInitParameters struct {

	// (Set of String) The IPv4 addresses this record set will point to.
	// The IPv4 addresses this record set will point to.
	Addresses []*string `json:"addresses,omitempty" tf:"addresses,omitempty"`

	// (Number) The TTL of the record set. Defaults to 3600.
	// The TTL of the record set. Defaults to `3600`.
	TTL *float64 `json:"ttl,omitempty" tf:"ttl,omitempty"`
}

type RecordSetObservation struct {

	// (Set of String) The IPv4 addresses this record set will point to.
	// The IPv4 addresses this record set will point to.
	Addresses []*string `json:"addresses,omitempty" tf:"addresses,omitempty"`

	// (String) The ID of this resource.
	ID *string `json:"id,omitempty" tf:"id,omitempty"`

	// (Number) The TTL of the record set. Defaults to 3600.
	// The TTL of the record set. Defaults to `3600`.
	TTL *float64 `json:"ttl,omitempty" tf:"ttl,omitempty"`

	// (String) DNS zone the record set belongs to. It must be an FQDN, that is, include the trailing dot.
	// DNS zone the record set belongs to. It must be an FQDN, that is, include the trailing dot.
	Zone *string `json:"zone,omitempty" tf:"zone,omitempty"`
}

type RecordSetParameters struct {

	// (Set of String) The IPv4 addresses this record set will point to.
	// The IPv4 addresses this record set will point to.
	// +kubebuilder:validation:Optional
	Addresses []*string `json:"addresses,omitempty" tf:"addresses,omitempty"`

	// (Number) The TTL of the record set. Defaults to 3600.
	// The TTL of the record set. Defaults to `3600`.
	// +kubebuilder:validation:Optional
	TTL *float64 `json:"ttl,omitempty" tf:"ttl,omitempty"`

	// (String) DNS zone the record set belongs to. It must be an FQDN, that is, include the trailing dot.
	// DNS zone the record set belongs to. It must be an FQDN, that is, include the trailing dot.
	// +kubebuilder:validation:Required
	Zone *string `json:"zone" tf:"zone,omitempty"`
}

// RecordSetSpec defines the desired state of RecordSet
type RecordSetSpec struct {
	v1.ResourceSpec `json:",inline"`
	ForProvider     RecordSetParameters `json:"forProvider"`
	// THIS IS A BETA FIELD. It will be honored
	// unless the Management Policies feature flag is disabled.
	// InitProvider holds the same fields as ForProvider, with the exception
	// of Identifier and other resource reference fields. The fields that are
	// in InitProvider are merged into ForProvider when the resource is created.
	// The same fields are also added to the terraform ignore_changes hook, to
	// avoid updating them after creation. This is useful for fields that are
	// required on creation, but we do not desire to update them after creation,
	// for example because of an external controller is managing them, like an
	// autoscaler.
	InitProvider RecordSetInitParameters `json:"initProvider,omitempty"`
}

// RecordSetStatus defines the observed state of RecordSet.
type RecordSetStatus struct {
	v1.ResourceStatus `json:",inline"`
	AtProvider        RecordSetObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// RecordSet is the Schema for the RecordSets API. Creates an A type DNS record set.
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="EXTERNAL-NAME",type="string",JSONPath=".metadata.annotations.crossplane\\.io/external-name"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,managed,dns}
type RecordSet struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	// +kubebuilder:validation:XValidation:rule="!('*' in self.managementPolicies || 'Create' in self.managementPolicies || 'Update' in self.managementPolicies) || has(self.forProvider.addresses) || (has(self.initProvider) && has(self.initProvider.addresses))",message="spec.forProvider.addresses is a required parameter"
	Spec   RecordSetSpec   `json:"spec"`
	Status RecordSetStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// RecordSetList contains a list of RecordSets
type RecordSetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RecordSet `json:"items"`
}

// Repository type metadata.
var (
	RecordSet_Kind             = "RecordSet"
	RecordSet_GroupKind        = schema.GroupKind{Group: CRDGroup, Kind: RecordSet_Kind}.String()
	RecordSet_KindAPIVersion   = RecordSet_Kind + "." + CRDGroupVersion.String()
	RecordSet_GroupVersionKind = CRDGroupVersion.WithKind(RecordSet_Kind)
)

func init() {
	SchemeBuilder.Register(&RecordSet{}, &RecordSetList{})
}