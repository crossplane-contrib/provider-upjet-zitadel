// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"context"
	"testing"

	xpv1 "github.com/crossplane/crossplane-runtime/v2/apis/common/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	clusterorg "github.com/crossplane-contrib/provider-upjet-zitadel/apis/cluster/org/v1alpha1"
	clusteruser "github.com/crossplane-contrib/provider-upjet-zitadel/apis/cluster/user/v1alpha1"
	org "github.com/crossplane-contrib/provider-upjet-zitadel/apis/namespaced/org/v1alpha1"
	user "github.com/crossplane-contrib/provider-upjet-zitadel/apis/namespaced/user/v1alpha1"
)

var tenantLabels = map[string]string{"example.org/tenant": "application"}

func newFakeClient(t *testing.T, objs ...client.Object) client.Client {
	t.Helper()
	scheme := runtime.NewScheme()
	for _, add := range []func(*runtime.Scheme) error{org.AddToScheme, user.AddToScheme, clusterorg.AddToScheme, clusteruser.AddToScheme} {
		if err := add(scheme); err != nil {
			t.Fatal(err)
		}
	}
	return fake.NewClientBuilder().WithScheme(scheme).WithObjects(objs...).Build()
}

func namespacedOrganization(namespace, externalName string) *org.Organization {
	o := &org.Organization{ObjectMeta: metav1.ObjectMeta{
		Name: "generated-org-name", Namespace: namespace, Labels: tenantLabels,
	}}
	if externalName != "" {
		o.Annotations = map[string]string{"crossplane.io/external-name": externalName}
	}
	return o
}

func TestMachineUserOrganizationSelector(t *testing.T) {
	organization := namespacedOrganization("tenant", "organization-id")
	other := namespacedOrganization("other-tenant", "wrong-organization")
	c := newFakeClient(t, organization, other)
	machine := &user.MachineUser{ObjectMeta: metav1.ObjectMeta{Name: "bot", Namespace: "tenant"}}
	machine.Spec.ForProvider.OrgIDSelector = &xpv1.NamespacedSelector{MatchLabels: tenantLabels}
	if err := machine.ResolveReferences(context.Background(), c); err != nil {
		t.Fatal(err)
	}
	if got := machine.Spec.ForProvider.OrgID; got == nil || *got != "organization-id" {
		t.Fatalf("resolved unexpected organization: %v", got)
	}
	if got := machine.Spec.ForProvider.OrgIDRef; got == nil || got.Name != organization.Name {
		t.Fatalf("reference was not persisted: %v", got)
	}
}

// An existing literal ID must win over a selector under the default
// IfNotPresent resolve policy, so adopting references never changes the
// organization of an existing resource.
func TestMachineUserLiteralOrganizationPreserved(t *testing.T) {
	c := newFakeClient(t, namespacedOrganization("tenant", "selected-organization"))
	literal := "literal-organization"
	machine := &user.MachineUser{ObjectMeta: metav1.ObjectMeta{Name: "bot", Namespace: "tenant"}}
	machine.Spec.ForProvider.OrgID = &literal
	machine.Spec.ForProvider.OrgIDSelector = &xpv1.NamespacedSelector{MatchLabels: tenantLabels}
	if err := machine.ResolveReferences(context.Background(), c); err != nil {
		t.Fatal(err)
	}
	if got := machine.Spec.ForProvider.OrgID; got == nil || *got != literal {
		t.Fatalf("literal organization was replaced: %v", got)
	}
	if got := machine.Spec.ForProvider.OrgIDRef; got != nil {
		t.Fatalf("unexpected reference for literal organization: %v", got)
	}
}

// Zitadel falls back to the caller's own organization when org_id is empty, so
// an unresolved selector must fail rather than leave org_id unset.
func TestMachineUserUnresolvedOrganizationSelector(t *testing.T) {
	cases := map[string][]client.Object{
		"no match in namespace":        {namespacedOrganization("other-tenant", "wrong-organization")},
		"organization not yet created": {namespacedOrganization("tenant", "")},
	}
	for name, objs := range cases {
		t.Run(name, func(t *testing.T) {
			machine := &user.MachineUser{ObjectMeta: metav1.ObjectMeta{Name: "bot", Namespace: "tenant"}}
			machine.Spec.ForProvider.OrgIDSelector = &xpv1.NamespacedSelector{MatchLabels: tenantLabels}
			if err := machine.ResolveReferences(context.Background(), newFakeClient(t, objs...)); err == nil {
				t.Fatal("expected an error for an unresolved organization selector")
			}
			if got := machine.Spec.ForProvider.OrgID; got != nil {
				t.Fatalf("organization was set despite failed resolution: %v", *got)
			}
		})
	}
}

func TestClusterMachineUserOrganizationSelector(t *testing.T) {
	organization := &clusterorg.Organization{ObjectMeta: metav1.ObjectMeta{
		Name: "generated-org-name", Labels: tenantLabels,
		Annotations: map[string]string{"crossplane.io/external-name": "organization-id"},
	}}
	unlabelled := &clusterorg.Organization{ObjectMeta: metav1.ObjectMeta{
		Name:        "unlabelled-org",
		Annotations: map[string]string{"crossplane.io/external-name": "wrong-organization"},
	}}
	c := newFakeClient(t, organization, unlabelled)
	machine := &clusteruser.MachineUser{ObjectMeta: metav1.ObjectMeta{Name: "bot"}}
	machine.Spec.ForProvider.OrgIDSelector = &xpv1.Selector{MatchLabels: tenantLabels}
	if err := machine.ResolveReferences(context.Background(), c); err != nil {
		t.Fatal(err)
	}
	if got := machine.Spec.ForProvider.OrgID; got == nil || *got != "organization-id" {
		t.Fatalf("resolved unexpected organization: %v", got)
	}
	if got := machine.Spec.ForProvider.OrgIDRef; got == nil || got.Name != organization.Name {
		t.Fatalf("reference was not persisted: %v", got)
	}
}
