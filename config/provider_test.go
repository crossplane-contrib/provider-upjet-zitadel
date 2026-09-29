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
		"project organization":                  {resource: "zitadel_project", field: "org_id", want: "zitadel_organization"},
		"human_user organization":               {resource: "zitadel_human_user", field: "org_id", want: "zitadel_organization"},
		"machine_user organization":             {resource: "zitadel_machine_user", field: "org_id", want: "zitadel_organization"},
		"application_api organization":          {resource: "zitadel_application_api", field: "org_id", want: "zitadel_organization"},
		"application_oidc organization":         {resource: "zitadel_application_oidc", field: "org_id", want: "zitadel_organization"},
		"application_saml organization":         {resource: "zitadel_application_saml", field: "org_id", want: "zitadel_organization"},
		"personal_access_token machine user":    {resource: "zitadel_personal_access_token", field: "user_id", want: "zitadel_machine_user"},
		"machine_key machine user":              {resource: "zitadel_machine_key", field: "user_id", want: "zitadel_machine_user"},
		"personal_access_token org_id":          {resource: "zitadel_personal_access_token", field: "org_id", want: "zitadel_organization"},
		"machine_key org_id":                    {resource: "zitadel_machine_key", field: "org_id", want: "zitadel_organization"},
		"user_grant org_id":                     {resource: "zitadel_user_grant", field: "org_id", want: "zitadel_organization"},
		"project_grant org_id":                  {resource: "zitadel_project_grant", field: "org_id", want: "zitadel_organization"},
		"user_grant project_id":                 {resource: "zitadel_user_grant", field: "project_id", want: "zitadel_project"},
		"user_grant project_grant_id":           {resource: "zitadel_user_grant", field: "project_grant_id", want: "zitadel_project_grant"},
		"project_grant granted_org_id":          {resource: "zitadel_project_grant", field: "granted_org_id", want: "zitadel_organization"},
		"organization_domain organization_id":   {resource: "zitadel_organization_domain", field: "organization_id", want: "zitadel_organization"},
		"application_key org_id":                {resource: "zitadel_application_key", field: "org_id", want: "zitadel_organization"},
		"project_role org_id":                   {resource: "zitadel_project_role", field: "org_id", want: "zitadel_organization"},
		"project_member org_id":                 {resource: "zitadel_project_member", field: "org_id", want: "zitadel_organization"},
		"project_grant_member org_id":           {resource: "zitadel_project_grant_member", field: "org_id", want: "zitadel_organization"},
		"project_grant_member grant_id":         {resource: "zitadel_project_grant_member", field: "grant_id", want: "zitadel_project_grant"},
		"user_metadata org_id":                  {resource: "zitadel_user_metadata", field: "org_id", want: "zitadel_organization"},
		"organization_metadata organization_id": {resource: "zitadel_organization_metadata", field: "organization_id", want: "zitadel_organization"},
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
