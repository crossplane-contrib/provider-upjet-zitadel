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
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	org "github.com/crossplane-contrib/provider-upjet-zitadel/apis/namespaced/org/v1alpha1"
	user "github.com/crossplane-contrib/provider-upjet-zitadel/apis/namespaced/user/v1alpha1"
)

func TestMachineUserOrganizationSelector(t *testing.T) {
	scheme := runtime.NewScheme()
	if err := org.AddToScheme(scheme); err != nil {
		t.Fatal(err)
	}
	if err := user.AddToScheme(scheme); err != nil {
		t.Fatal(err)
	}
	organization := &org.Organization{ObjectMeta: metav1.ObjectMeta{
		Name: "generated-org-name", Namespace: "tenant", Labels: map[string]string{"app": "harmony"},
		Annotations: map[string]string{"crossplane.io/external-name": "organization-id"},
	}}
	other := organization.DeepCopy()
	other.Namespace = "other-tenant"
	other.Annotations["crossplane.io/external-name"] = "wrong-organization"
	client := fake.NewClientBuilder().WithScheme(scheme).WithObjects(organization, other).Build()
	machine := &user.MachineUser{ObjectMeta: metav1.ObjectMeta{Name: "bot", Namespace: "tenant"}}
	machine.Spec.ForProvider.OrgIDSelector = &xpv1.NamespacedSelector{MatchLabels: map[string]string{"app": "harmony"}}
	if err := machine.ResolveReferences(context.Background(), client); err != nil {
		t.Fatal(err)
	}
	if got := machine.Spec.ForProvider.OrgID; got == nil || *got != "organization-id" {
		t.Fatalf("resolved unexpected organization: %v", got)
	}
	if got := machine.Spec.ForProvider.OrgIDRef; got == nil || got.Name != organization.Name {
		t.Fatalf("reference was not persisted: %v", got)
	}
}
