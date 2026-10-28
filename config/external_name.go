/*
Copyright 2024 The Crossplane Authors.
*/

// Package config contains the provider configuration and resource external-name mappings.
package config

import (
	"slices"

	"github.com/crossplane/upjet/v2/pkg/config"
)

// New stable v2 resources use the embedded SDK provider. Existing controllers
// continue to use the Terraform CLI, preserving their reconciliation backend.
var sdkV2Resources = []string{
	"zitadel_authorization$", "zitadel_administrator$",
	"zitadel_project_v2$", "zitadel_application_v2$",
}

func terraformCLIConfigured() []string {
	var result []string
	for _, name := range ExternalNameConfigured() {
		if !slices.Contains(sdkV2Resources, name) {
			result = append(result, name)
		}
	}
	return result
}

// ExternalNameConfigurations applies all external name configs listed in the
// table ExternalNameConfigs and sets the version of those resources to v1beta1
// assuming they will be tested.
func ExternalNameConfigurations() config.ResourceOption {
	return func(r *config.Resource) {
		if e, ok := ExternalNameConfigs[r.Name]; ok {
			r.ExternalName = e
		}
	}
}

// ExternalNameConfigured returns the list of all resources whose external name
// is configured manually.
func ExternalNameConfigured() []string {
	l := make([]string, len(ExternalNameConfigs))
	i := 0
	for name := range ExternalNameConfigs {
		// $ is added to match the exact string since the format is regex.
		l[i] = name + "$"
		i++
	}
	return l
}
