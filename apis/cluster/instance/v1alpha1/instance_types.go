// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	v2 "github.com/crossplane/crossplane/apis/v2/core/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/crossplane-contrib/provider-upjet-zitadel/apis/observation"
)

// InstanceSpec configures read-only observation of the instance selected by
// the ProviderConfig's endpoint, authentication and transport headers.
type InstanceSpec struct {

	// ManagementPolicies permits observation or pausing only. This resource
	// never creates, updates or deletes a Zitadel instance.
	// +optional
	// +kubebuilder:default={"Observe"}
	// +kubebuilder:validation:XValidation:rule="self.all(p, p == 'Observe')",message="Instance supports only Observe or an empty policy list"
	ManagementPolicies v2.ManagementPolicies `json:"managementPolicies,omitempty"`
	// ProviderConfigReference selects credentials and endpoint for observation.
	// +optional
	// +kubebuilder:default={"name":"default"}
	ProviderConfigReference *v2.Reference `json:"providerConfigRef,omitempty"`
	// WriteConnectionSecretToReference is retained for managed-resource interface
	// compatibility. Instance observation does not publish credentials.
	// +optional
	WriteConnectionSecretToReference *v2.SecretReference `json:"writeConnectionSecretToRef,omitempty"`
}

// InstanceStatus contains the observed instance identity and conditions.
type InstanceStatus struct {
	v2.ManagedResourceStatus `json:",inline"`
	AtProvider               observation.Instance `json:"atProvider,omitempty"`
}

// Instance observes a Helm-bootstrapped Zitadel instance without owning its
// lifecycle. No instance ID is required. Its endpoint comes from ProviderConfig.
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:storageversion
// +kubebuilder:resource:scope=Cluster,categories={crossplane,managed,zitadel}
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="ID",type="string",JSONPath=".status.atProvider.id"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
type Instance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              InstanceSpec   `json:"spec"`
	Status            InstanceStatus `json:"status,omitempty"`
}

// InstanceList contains Instance resources.
// +kubebuilder:object:root=true
type InstanceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Instance `json:"items"`
}

// GetInstanceObservation returns the non-secret instance metadata.
func (i *Instance) GetInstanceObservation() observation.Instance { return i.Status.AtProvider }

// SetInstanceObservation updates the non-secret instance metadata.
func (i *Instance) SetInstanceObservation(o observation.Instance) { i.Status.AtProvider = o }

// InstanceGroupVersionKind identifies the observation API.
var InstanceGroupVersionKind = CRDGroupVersion.WithKind("Instance")

func init() { SchemeBuilder.Register(&Instance{}, &InstanceList{}) }
