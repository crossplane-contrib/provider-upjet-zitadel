// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

// Package v1alpha1 contains instance configuration and observation APIs.
package v1alpha1

import (
	v1 "github.com/crossplane/crossplane-runtime/v2/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/v2/pkg/resource"
)

// GetCondition returns an observation condition.
func (i *Instance) GetCondition(t v1.ConditionType) v1.Condition { return i.Status.GetCondition(t) }

// SetConditions updates observation conditions.
func (i *Instance) SetConditions(c ...v1.Condition) { i.Status.SetConditions(c...) }

// GetManagementPolicies returns the allowed read-only actions.
func (i *Instance) GetManagementPolicies() v1.ManagementPolicies { return i.Spec.ManagementPolicies }

// SetManagementPolicies updates allowed actions, validated as Observe or paused.
func (i *Instance) SetManagementPolicies(p v1.ManagementPolicies) { i.Spec.ManagementPolicies = p }

// GetProviderConfigReference returns the endpoint/credentials reference.
func (i *Instance) GetProviderConfigReference() *v1.ProviderConfigReference {
	return i.Spec.ProviderConfigReference
}

// SetProviderConfigReference updates the endpoint/credentials reference.
func (i *Instance) SetProviderConfigReference(p *v1.ProviderConfigReference) {
	i.Spec.ProviderConfigReference = p
}

// GetWriteConnectionSecretToReference implements the managed-resource contract.
func (i *Instance) GetWriteConnectionSecretToReference() *v1.LocalSecretReference {
	return i.Spec.WriteConnectionSecretToReference
}

// SetWriteConnectionSecretToReference implements the managed-resource contract.
func (i *Instance) SetWriteConnectionSecretToReference(r *v1.LocalSecretReference) {
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
