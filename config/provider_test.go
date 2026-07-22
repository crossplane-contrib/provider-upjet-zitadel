// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"reflect"
	"testing"

	ujresource "github.com/crossplane/upjet/v2/pkg/resource"

	applicationv1alpha1 "github.com/crossplane-contrib/provider-upjet-zitadel/apis/namespaced/application/v1alpha1"
)

func TestRequiredCrossResourceReferences(t *testing.T) {
	p := GetProviderNamespaced()
	tests := map[string]struct {
		resource     string
		field        string
		want         string
		wantRef      string
		wantSelector string
	}{
		"login policy identity providers": {
			resource:     "zitadel_login_policy",
			field:        "idps",
			want:         "zitadel_org_idp_google",
			wantRef:      "IdpGoogleRefs",
			wantSelector: "IdpGoogleSelector",
		},
		"trigger actions": {
			resource: "zitadel_trigger_actions",
			field:    "action_ids",
			want:     "zitadel_action",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := p.Resources[tt.resource].References[tt.field]
			if got.TerraformName != tt.want {
				t.Errorf("reference target = %q, want %q", got.TerraformName, tt.want)
			}
			if got.RefFieldName != tt.wantRef {
				t.Errorf("reference field = %q, want %q", got.RefFieldName, tt.wantRef)
			}
			if got.SelectorFieldName != tt.wantSelector {
				t.Errorf("selector field = %q, want %q", got.SelectorFieldName, tt.wantSelector)
			}
		})
	}
}

func TestOIDCClientCredentialsConnectionDetails(t *testing.T) {
	mapping := (&applicationv1alpha1.Oidc{}).GetConnectionDetailsMapping()

	got, err := ujresource.GetSensitiveAttributes(map[string]any{
		"client_id":     "generated-client-id",
		"client_secret": "generated-client-secret",
	}, mapping)
	if err != nil {
		t.Fatalf("GetSensitiveAttributes() error = %v", err)
	}

	want := map[string][]byte{
		"attribute.client_id":     []byte("generated-client-id"),
		"attribute.client_secret": []byte("generated-client-secret"),
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("connection details = %#v, want %#v", got, want)
	}
}
