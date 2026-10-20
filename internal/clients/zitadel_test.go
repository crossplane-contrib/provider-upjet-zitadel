// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package clients

import (
	"context"
	"encoding/json"
	"testing"

	xpv1 "github.com/crossplane/crossplane-runtime/v2/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/v2/pkg/resource"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	application "github.com/crossplane-contrib/provider-upjet-zitadel/apis/namespaced/application/v1alpha1"
	authorization "github.com/crossplane-contrib/provider-upjet-zitadel/apis/namespaced/authorization/v1alpha1"
	permission "github.com/crossplane-contrib/provider-upjet-zitadel/apis/namespaced/internalpermission/v1alpha1"
	project "github.com/crossplane-contrib/provider-upjet-zitadel/apis/namespaced/project/v1alpha1"
	pc "github.com/crossplane-contrib/provider-upjet-zitadel/apis/namespaced/v1beta1"
)

func TestSDKSetupUsesCurrentProviderCredentials(t *testing.T) {
	ctx := context.Background()
	scheme := runtime.NewScheme()
	for _, add := range []func(*runtime.Scheme) error{corev1.AddToScheme, pc.SchemeBuilder.AddToScheme, application.AddToScheme, authorization.AddToScheme, permission.AddToScheme, project.AddToScheme} {
		if err := add(scheme); err != nil {
			t.Fatal(err)
		}
	}
	credentials := &corev1.Secret{ObjectMeta: metav1.ObjectMeta{Name: "credentials", Namespace: "app"}}
	config := &pc.ProviderConfig{ObjectMeta: metav1.ObjectMeta{Name: "default", Namespace: "app"}, Spec: pc.ProviderConfigSpec{Credentials: pc.ProviderCredentials{Source: xpv1.CredentialsSourceSecret, CommonCredentialSelectors: xpv1.CommonCredentialSelectors{SecretRef: &xpv1.SecretKeySelector{SecretReference: xpv1.SecretReference{Name: credentials.Name, Namespace: credentials.Namespace}, Key: "credentials"}}}}}
	kube := fake.NewClientBuilder().WithScheme(scheme).WithObjects(credentials, config).Build()
	setup := TerraformSetupBuilder("1.5.7", "zitadel/zitadel", "2.12.8")
	var previous *helper.ClientInfo
	for _, token := range []string{"first", "rotated"} {
		var current *helper.ClientInfo
		data, err := json.Marshal(map[string]any{"domain": token + ".example.test", "access_token": token, "transport_headers": map[string]string{"x-test": token}})
		if err != nil {
			t.Fatal(err)
		}
		credentials.Data = map[string][]byte{"credentials": data}
		if err := kube.Update(ctx, credentials); err != nil {
			t.Fatal(err)
		}
		for _, mg := range []resource.ModernManaged{&project.ProjectV2{}, &application.Application{}, &authorization.Authorization{}, &permission.Administrator{}} {
			mg.SetName("assignment")
			mg.SetUID("assignment")
			mg.SetNamespace("app")
			mg.SetProviderConfigReference(&xpv1.ProviderConfigReference{Name: "default", Kind: "ProviderConfig"})
			got, err := setup(ctx, kube, mg)
			if err != nil {
				t.Fatal(err)
			}
			info, ok := got.Meta.(*helper.ClientInfo)
			if !ok || info == nil {
				t.Fatalf("%T SDK setup omitted configured ClientInfo", mg)
			}
			bearer, err := info.TokenSource.Token()
			if err != nil || bearer.AccessToken != token || info.Domain != token+".example.test:443" || info.TransportHeaders["x-test"] != token {
				t.Fatal("SDK metadata reused credentials or omitted transport configuration")
			}
			if current == nil {
				current = info
			} else if info != current {
				t.Fatal("reconciles with unchanged credentials did not reuse metadata")
			}
			if info == previous {
				t.Fatal("credential rotation reused old provider metadata")
			}
		}
		previous = current
	}
	alias := config.DeepCopy()
	alias.Name = "other"
	alias.ResourceVersion = ""
	if err := kube.Create(ctx, alias); err != nil {
		t.Fatal(err)
	}
	aliased := &project.ProjectV2{}
	aliased.SetName("aliased")
	aliased.SetUID("aliased")
	aliased.SetNamespace("app")
	aliased.SetProviderConfigReference(&xpv1.ProviderConfigReference{Name: "other", Kind: "ProviderConfig"})
	if got, err := setup(ctx, kube, aliased); err != nil || got.Meta == previous {
		t.Fatal("different ProviderConfigs shared SDK metadata")
	}
	legacy := &project.Project{}
	legacy.SetName("legacy")
	legacy.SetUID("legacy")
	legacy.SetNamespace("app")
	legacy.SetProviderConfigReference(&xpv1.ProviderConfigReference{Name: "default", Kind: "ProviderConfig"})
	got, err := setup(ctx, kube, legacy)
	if err != nil || got.Meta != nil || got.Requirement.Version != "2.12.8" {
		t.Fatal("legacy CLI setup changed")
	}
	credentials.Data["credentials"] = []byte(`{"domain":"example.test"}`)
	if err := kube.Update(ctx, credentials); err != nil {
		t.Fatal(err)
	}
	invalid := &project.ProjectV2{}
	invalid.SetName("invalid")
	invalid.SetUID("invalid")
	invalid.SetNamespace("app")
	invalid.SetProviderConfigReference(&xpv1.ProviderConfigReference{Name: "default", Kind: "ProviderConfig"})
	if _, err := setup(ctx, kube, invalid); err == nil {
		t.Fatal("invalid SDK credentials accepted")
	}
}
