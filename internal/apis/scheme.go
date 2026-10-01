// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

// Package apis looks up reference targets for the generated resolvers.
//
// The upjet resolver transformer rewrites the generated ResolveReferences
// methods to call GetManagedResource instead of importing the target API
// package. This breaks the import cycles that cross-group references would
// otherwise create, e.g. org.Member -> user.HumanUser -> org.Organization.
package apis

import (
	xpresource "github.com/crossplane/crossplane-runtime/v2/pkg/resource"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var s = runtime.NewScheme()

// GetManagedResource returns new, empty instances of the managed resource
// kind and its list kind for the given group and version. The kinds must have
// been registered with BuildScheme.
func GetManagedResource(group, version, kind, listKind string) (xpresource.Managed, xpresource.ManagedList, error) {
	gv := schema.GroupVersion{Group: group, Version: version}
	kindGVK := gv.WithKind(kind)
	o, err := s.New(kindGVK)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to get a new API object of GVK %q from the runtime scheme", kindGVK)
	}
	listGVK := gv.WithKind(listKind)
	lo, err := s.New(listGVK)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to get a new API object list of GVK %q from the runtime scheme", listGVK)
	}
	m, ok := o.(xpresource.Managed)
	if !ok {
		return nil, nil, errors.Errorf("GVK %q is not a managed resource", kindGVK)
	}
	l, ok := lo.(xpresource.ManagedList)
	if !ok {
		return nil, nil, errors.Errorf("GVK %q is not a managed resource list", listGVK)
	}
	return m, l, nil
}

// BuildScheme registers API types with the scheme used by GetManagedResource.
func BuildScheme(sb runtime.SchemeBuilder) error {
	return errors.Wrap(sb.AddToScheme(s), "failed to register the GVKs with the runtime scheme")
}
