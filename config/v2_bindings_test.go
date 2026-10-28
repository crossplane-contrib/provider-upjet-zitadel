// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"context"
	"testing"

	xpv1 "github.com/crossplane/crossplane-runtime/v2/apis/common/v1"
	ujconfig "github.com/crossplane/upjet/v2/pkg/config"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	authorization "github.com/crossplane-contrib/provider-upjet-zitadel/apis/namespaced/authorization/v1alpha1"
	permission "github.com/crossplane-contrib/provider-upjet-zitadel/apis/namespaced/internalpermission/v1alpha1"
	org "github.com/crossplane-contrib/provider-upjet-zitadel/apis/namespaced/org/v1alpha1"
	project "github.com/crossplane-contrib/provider-upjet-zitadel/apis/namespaced/project/v1alpha1"
)

func TestStableV2BindingsUseEmbeddedSDK(t *testing.T) {
	for _, p := range []*ujconfig.Provider{GetProvider(), GetProviderNamespaced()} {
		if err := p.TerraformProvider.InternalValidate(); err != nil {
			t.Fatal(err)
		}
		for _, name := range []string{"zitadel_authorization", "zitadel_administrator", "zitadel_project_v2", "zitadel_application_v2"} {
			r := p.Resources[name]
			if r == nil || !r.ShouldUseTerraformPluginSDKClient() || r.UseAsync {
				t.Fatalf("%s must use the synchronous embedded SDK", name)
			}
			if r.TerraformResource != p.TerraformProvider.ResourcesMap[name] {
				t.Fatalf("%s is not backed by the embedded resource", name)
			}
			if _, exists := r.References["user_id"]; exists {
				t.Fatalf("%s must not restrict the polymorphic subject to one user kind", name)
			}
		}
		for _, name := range []string{"zitadel_user_grant", "zitadel_org_member", "zitadel_project", "zitadel_application_oidc"} {
			if p.Resources[name].ShouldUseTerraformPluginSDKClient() {
				t.Fatalf("legacy backend changed for %s", name)
			}
		}
		for _, name := range []string{"zitadel_authorization", "zitadel_administrator", "zitadel_application_v2"} {
			if got := p.Resources[name].References["project_id"].TerraformName; got != "zitadel_project_v2" {
				t.Fatalf("%s uses legacy project target %s", name, got)
			}
			if got := p.Resources[name].References["org_id"].TerraformName; got != "zitadel_organization" {
				t.Fatalf("%s uses legacy organization target %s", name, got)
			}
		}
	}
}

func TestStableV2AssignmentReferences(t *testing.T) {
	scheme := runtime.NewScheme()
	for _, add := range []func(*runtime.Scheme) error{org.AddToScheme, project.AddToScheme, authorization.AddToScheme, permission.AddToScheme} {
		if err := add(scheme); err != nil {
			t.Fatal(err)
		}
	}
	organization := &org.Organization{ObjectMeta: metav1.ObjectMeta{Name: "tenant", Namespace: "app", Annotations: map[string]string{"crossplane.io/external-name": "organization-id"}}}
	product := &project.ProjectV2{ObjectMeta: metav1.ObjectMeta{Name: "product", Namespace: "app", Annotations: map[string]string{"crossplane.io/external-name": "project-id"}}}
	otherProduct := &project.ProjectV2{ObjectMeta: metav1.ObjectMeta{Name: "product", Namespace: "other", Annotations: map[string]string{"crossplane.io/external-name": "wrong-project-id"}}}
	c := fake.NewClientBuilder().WithScheme(scheme).WithObjects(organization, product, otherProduct).Build()
	a := &authorization.Authorization{ObjectMeta: metav1.ObjectMeta{Name: "roles", Namespace: "app"}}
	a.Spec.ForProvider.ProjectIDRef = &xpv1.NamespacedReference{Name: "product"}
	a.Spec.ForProvider.OrgIDRef = &xpv1.NamespacedReference{Name: "tenant"}
	if err := a.ResolveReferences(context.Background(), c); err != nil {
		t.Fatal(err)
	}
	if a.Spec.ForProvider.ProjectID == nil || *a.Spec.ForProvider.ProjectID != "project-id" || a.Spec.ForProvider.OrgID == nil || *a.Spec.ForProvider.OrgID != "organization-id" {
		t.Fatalf("wrong assignment target: %#v", a.Spec.ForProvider)
	}
	m := &permission.Administrator{ObjectMeta: metav1.ObjectMeta{Name: "admin", Namespace: "app"}}
	m.Spec.ForProvider.OrgIDRef = &xpv1.NamespacedReference{Name: "tenant"}
	if err := m.ResolveReferences(context.Background(), c); err != nil {
		t.Fatal(err)
	}
	if m.Spec.ForProvider.OrgID == nil || *m.Spec.ForProvider.OrgID != "organization-id" {
		t.Fatal("administrator did not resolve canonical Organization")
	}
	a.Spec.ForProvider.ProjectID = nil
	a.Spec.ForProvider.ProjectIDRef.Name = "missing"
	if err := a.ResolveReferences(context.Background(), c); err == nil {
		t.Fatal("missing reference silently accepted")
	}
}
