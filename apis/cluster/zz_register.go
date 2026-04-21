/*
Copyright 2024 The Crossplane Authors.
*/

package cluster

import (
	"k8s.io/apimachinery/pkg/runtime"

	v1beta1 "github.com/crossplane-contrib/provider-upjet-zitadel/apis/cluster/v1beta1"
)

func init() {
	// Register the types with the Scheme so the packages are imported.
	AddToScheme = func(s *runtime.Scheme) error {
		return v1beta1.SchemeBuilder.AddToScheme(s)
	}
}

// AddToScheme adds all registered types to the scheme.
var AddToScheme func(s *runtime.Scheme) error
