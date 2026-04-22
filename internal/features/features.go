/*
Copyright 2024 The Crossplane Authors.
*/

// Package features defines feature flags for this provider.
package features

import (
	xpfeature "github.com/crossplane/crossplane-runtime/v2/pkg/feature"
)

const (
	// EnableBetaManagementPolicies enables beta management policy support.
	EnableBetaManagementPolicies xpfeature.Flag = xpfeature.EnableBetaManagementPolicies
)
