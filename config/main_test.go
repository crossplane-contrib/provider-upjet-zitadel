// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"os"
	"testing"

	apisCluster "github.com/crossplane-contrib/provider-upjet-zitadel/apis/cluster"
	apisNamespaced "github.com/crossplane-contrib/provider-upjet-zitadel/apis/namespaced"
	resolverapis "github.com/crossplane-contrib/provider-upjet-zitadel/internal/apis"
)

// TestMain registers the API types with the resolver scheme, as the provider
// does at startup, so tests can call the generated ResolveReferences methods.
func TestMain(m *testing.M) {
	if err := resolverapis.BuildScheme(apisCluster.AddToSchemes); err != nil {
		panic(err)
	}
	if err := resolverapis.BuildScheme(apisNamespaced.AddToSchemes); err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}
