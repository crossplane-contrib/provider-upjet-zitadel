// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

// Package observation contains shared read-only instance observations.
package observation

import (
	xpv1 "github.com/crossplane/crossplane-runtime/v2/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/v2/pkg/reference"
	"github.com/crossplane/crossplane-runtime/v2/pkg/resource"
)

// Instance contains non-secret metadata returned by Zitadel.
type Instance struct {
	ID      string `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
	State   string `json:"state,omitempty"`
}

// ObservedInstance is implemented by both scopes of the Instance resource.
type ObservedInstance interface {
	resource.Managed
	GetInstanceObservation() Instance
	SetInstanceObservation(Instance)
}

// InstanceID extracts an ID only after a successful observation. An annotation
// or stale status alone is not evidence that the configured instance exists.
func InstanceID() reference.ExtractValueFn {
	return func(mg resource.Managed) string {
		inst, ok := mg.(ObservedInstance)
		if !ok || inst.GetCondition(xpv1.TypeReady).Status != "True" || inst.GetCondition(xpv1.TypeSynced).Status != "True" || inst.GetCondition(xpv1.TypeSynced).ObservedGeneration < inst.GetGeneration() || inst.GetDeletionTimestamp() != nil {
			return ""
		}
		return inst.GetInstanceObservation().ID
	}
}
