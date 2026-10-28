// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package apis_test

import (
	"testing"

	apisNamespaced "github.com/crossplane-contrib/provider-upjet-zitadel/apis/namespaced"
	user "github.com/crossplane-contrib/provider-upjet-zitadel/apis/namespaced/user/v1alpha1"
	"github.com/crossplane-contrib/provider-upjet-zitadel/internal/apis"
)

func TestGetManagedResource(t *testing.T) {
	if err := apis.BuildScheme(apisNamespaced.AddToSchemes); err != nil {
		t.Fatal(err)
	}

	m, l, err := apis.GetManagedResource("user.zitadel.m.crossplane.io", "v1alpha1", "HumanUser", "HumanUserList")
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := m.(*user.HumanUser); !ok {
		t.Errorf("managed resource type = %T, want *HumanUser", m)
	}
	if _, ok := l.(*user.HumanUserList); !ok {
		t.Errorf("managed resource list type = %T, want *HumanUserList", l)
	}

	for name, gvk := range map[string][4]string{
		"unregistered group": {"user.zitadel.crossplane.io", "v1alpha1", "HumanUser", "HumanUserList"},
		"unknown kind":       {"user.zitadel.m.crossplane.io", "v1alpha1", "User", "UserList"},
		"not managed":        {"zitadel.m.crossplane.io", "v1beta1", "ProviderConfig", "ProviderConfigList"},
	} {
		t.Run(name, func(t *testing.T) {
			if _, _, err := apis.GetManagedResource(gvk[0], gvk[1], gvk[2], gvk[3]); err == nil {
				t.Fatal("expected an error")
			}
		})
	}
}
