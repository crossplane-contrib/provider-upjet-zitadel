// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

// Package v1alpha1 contains instance configuration and observation APIs.
package v1alpha1

import (
	"github.com/crossplane/crossplane-runtime/v2/pkg/resource"
	v2 "github.com/crossplane/crossplane/apis/v2/core/v2"
)

// GetCondition returns an observation condition.
func (i *Instance) GetCondition(t v2.ConditionType) v2.Condition { return i.Status.GetCondition(t) }

// SetConditions updates observation conditions.
func (i *Instance) SetConditions(c ...v2.Condition) { i.Status.SetConditions(c...) }

// GetManagementPolicies returns the allowed read-only actions.
func (i *Instance) GetManagementPolicies() v2.ManagementPolicies { return i.Spec.ManagementPolicies }

// SetManagementPolicies updates allowed actions, validated as Observe or paused.
func (i *Instance) SetManagementPolicies(p v2.ManagementPolicies) { i.Spec.ManagementPolicies = p }

// GetProviderConfigReference returns the endpoint/credentials reference.
func (i *Instance) GetProviderConfigReference() *v2.Reference { return i.Spec.ProviderConfigReference }

// SetProviderConfigReference updates the endpoint/credentials reference.
func (i *Instance) SetProviderConfigReference(p *v2.Reference) { i.Spec.ProviderConfigReference = p }

// GetWriteConnectionSecretToReference implements the managed-resource contract.
func (i *Instance) GetWriteConnectionSecretToReference() *v2.SecretReference {
	return i.Spec.WriteConnectionSecretToReference
}

// SetWriteConnectionSecretToReference implements the managed-resource contract.
func (i *Instance) SetWriteConnectionSecretToReference(r *v2.SecretReference) {
	i.Spec.WriteConnectionSecretToReference = r
}

// GetItems returns managed observations.
func (l *InstanceList) GetItems() []resource.Managed {
	items := make([]resource.Managed, len(l.Items))
	for n := range l.Items {
		items[n] = &l.Items[n]
	}
	return items
}

// GetDeletionPolicy keeps legacy interface compatibility without exposing a
// deletionPolicy API or granting ownership of the external instance.
func (i *Instance) GetDeletionPolicy() v2.DeletionPolicy { return v2.DeletionOrphan }

// SetDeletionPolicy cannot grant ownership of an observed instance.
func (i *Instance) SetDeletionPolicy(v2.DeletionPolicy) {}
