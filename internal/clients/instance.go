// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package clients

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/crossplane/crossplane-runtime/v2/pkg/resource"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	tfhelper "github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"github.com/zitadel/zitadel-go/v3/pkg/client/admin"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel"
	adminpb "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/admin"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/crossplane-contrib/provider-upjet-zitadel/apis/observation"
)

// InstanceClient can only observe the instance selected by the authenticated
// endpoint. Its interface intentionally exposes no instance write operations.
type InstanceClient interface {
	Observe(context.Context) (observation.Instance, error)
	Close() error
}

type instanceClient struct{ client *admin.Client }

// NewInstanceClient uses the same credentials, authentication builder and
// transport options as the pinned Terraform provider, but no global client
// cache: every reconcile sees rotated credentials and its own ProviderConfig.
func NewInstanceClient(ctx context.Context, kube client.Client, mg resource.Managed) (InstanceClient, error) {
	spec, err := resolveProviderConfig(ctx, kube, mg)
	if err != nil {
		return nil, fmt.Errorf("cannot resolve provider config: %w", err)
	}
	data, err := resource.CommonCredentialExtractor(ctx, spec.Credentials.Source, kube, spec.Credentials.CommonCredentialSelectors)
	if err != nil {
		return nil, fmt.Errorf("cannot extract credentials: %w", err)
	}
	return instanceClientFromCredentials(ctx, data)
}

type instanceCredentials struct {
	Domain             string            `json:"domain"`
	Port               string            `json:"port"`
	Insecure           bool              `json:"insecure"`
	InsecureSkipVerify bool              `json:"insecure_skip_verify_tls"`
	AccessToken        string            `json:"access_token"`
	Token              string            `json:"token"`
	JWTFile            string            `json:"jwt_file"`
	JWTProfileFile     string            `json:"jwt_profile_file"`
	JWTProfileJSON     string            `json:"jwt_profile_json"`
	TransportHeaders   map[string]string `json:"transport_headers"`
	SystemAPI          []struct {
		KeyFile    string `json:"key_file"`
		Key        string `json:"key"`
		PrivateKey string `json:"private_key"`
		PublicKey  string `json:"public_key"`
		User       string `json:"user"`
		Audience   string `json:"audience"`
	} `json:"system_api"`
}

func instanceClientFromCredentials(ctx context.Context, data []byte) (InstanceClient, error) {
	var c instanceCredentials
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, fmt.Errorf("cannot parse instance credentials JSON")
	}
	if c.Domain == "" {
		return nil, fmt.Errorf("instance credentials require domain")
	}
	count := 0
	for _, value := range []string{c.AccessToken, c.Token, c.JWTFile, c.JWTProfileFile, c.JWTProfileJSON} {
		if value != "" {
			count++
		}
	}
	var keyFile, key, privateKey, publicKey, user, audience string
	if len(c.SystemAPI) > 1 {
		return nil, fmt.Errorf("only one system_api block may be configured")
	}
	if len(c.SystemAPI) == 1 {
		s := c.SystemAPI[0]
		keyFile, key, privateKey, publicKey, user, audience = s.KeyFile, s.Key, s.PrivateKey, s.PublicKey, s.User, s.Audience
		if keyFile != "" || key != "" || privateKey != "" || publicKey != "" {
			count++
		}
	}
	if count != 1 {
		return nil, fmt.Errorf("exactly one authentication method must be configured")
	}
	info, err := tfhelper.GetClientInfo(ctx, c.Insecure, c.Domain, c.AccessToken, c.Token, c.JWTFile, c.JWTProfileFile, c.JWTProfileJSON, keyFile, key, privateKey, publicKey, user, audience, c.Port, c.InsecureSkipVerify, c.TransportHeaders)
	// Authentication builders may report sensitive key parsing details.
	if err != nil {
		return nil, fmt.Errorf("cannot configure instance authentication")
	}
	api, err := admin.NewClient(ctx, info.Issuer, info.Domain, []string{oidc.ScopeOpenID, zitadel.ScopeZitadelAPI()}, info.Options...)
	if err != nil {
		return nil, fmt.Errorf("cannot create instance client")
	}
	return &instanceClient{client: api}, nil
}

func (c *instanceClient) Observe(ctx context.Context) (observation.Instance, error) {
	//nolint:staticcheck // Keep discovery compatible with installations using the existing Admin API; it needs no instance ID.
	response, err := c.client.GetMyInstance(ctx, &adminpb.GetMyInstanceRequest{})
	if err != nil {
		return observation.Instance{}, fmt.Errorf("cannot observe authenticated Zitadel instance: %w", err)
	}
	i := response.GetInstance()
	if i.GetId() == "" {
		return observation.Instance{}, fmt.Errorf("zitadel returned an empty instance ID")
	}
	return observation.Instance{ID: i.GetId(), Name: i.GetName(), Version: i.GetVersion(), State: i.GetState().String()}, nil
}

func (c *instanceClient) Close() error { return c.client.Connection.Close() }
