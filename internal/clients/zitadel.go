/*
Copyright 2024 The Crossplane Authors.
*/

// Package clients provides the Zitadel API client used by managed resource controllers.
package clients

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"sync"

	"github.com/crossplane/crossplane-runtime/v2/pkg/resource"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	ujresource "github.com/crossplane/upjet/v2/pkg/resource"
	"github.com/crossplane/upjet/v2/pkg/terraform"
	tfsdk "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	tfprovider "github.com/zitadel/terraform-provider-zitadel/v2/zitadel"

	clusterv1beta1 "github.com/crossplane-contrib/provider-upjet-zitadel/apis/cluster/v1beta1"
	namespacedv1beta1 "github.com/crossplane-contrib/provider-upjet-zitadel/apis/namespaced/v1beta1"
)

const (
	errNoProviderConfig     = "no providerConfigRef provided"
	errGetProviderConfig    = "cannot get referenced ProviderConfig"
	errTrackUsage           = "cannot track ProviderConfig usage"
	errExtractCredentials   = "cannot extract credentials"
	errUnmarshalCredentials = "cannot unmarshal zitadel credentials as JSON"
)

// TerraformSetupBuilder builds a terraform.SetupFn which returns Terraform provider setup configuration.
func TerraformSetupBuilder(version, providerSource, providerVersion string) terraform.SetupFn {
	// SDK clients own gRPC connections. Reuse metadata across reconciles of a
	// ProviderConfig, replacing it when credential contents change. Constructing
	// metadata for every reconcile would leave a new connection behind each time.
	type sdkMetadata struct {
		credentials [sha256.Size]byte
		meta        any
	}
	var sdkMu sync.Mutex
	sdkCache := map[string]sdkMetadata{}
	return func(ctx context.Context, crClient client.Client, mg resource.Managed) (terraform.Setup, error) {
		ps := terraform.Setup{
			Version: version,
			Requirement: terraform.ProviderRequirement{
				Source:  providerSource,
				Version: providerVersion,
			},
		}

		pcSpec, err := resolveProviderConfig(ctx, crClient, mg)
		if err != nil {
			return terraform.Setup{}, errors.Wrap(err, "cannot resolve provider config")
		}

		data, err := resource.CommonCredentialExtractor(ctx, pcSpec.Credentials.Source, crClient, pcSpec.Credentials.CommonCredentialSelectors)
		if err != nil {
			return ps, errors.Wrap(err, errExtractCredentials)
		}

		if err := json.Unmarshal(data, &ps.Configuration); err != nil {
			return ps, errors.Wrap(err, errUnmarshalCredentials)
		}

		if r, ok := mg.(ujresource.Terraformed); ok {
			switch r.GetTerraformResourceType() {
			case "zitadel_authorization", "zitadel_administrator", "zitadel_project_v2", "zitadel_application_v2":
				// SDK connectors receive Setup.Meta directly; unlike the CLI they
				// do not configure the Terraform provider themselves.
				key := sdkProviderConfigKey(mg)
				signature := sha256.Sum256(data)
				sdkMu.Lock()
				defer sdkMu.Unlock()
				if cached, ok := sdkCache[key]; ok && cached.credentials == signature {
					ps.Meta = cached.meta
					break
				}
				p := tfprovider.Provider()
				if diagnostics := p.Configure(ctx, tfsdk.NewResourceConfigRaw(ps.Configuration)); diagnostics.HasError() {
					return ps, errors.New("cannot configure embedded Zitadel SDK provider: " + diagnostics[0].Summary)
				}
				ps.Meta = p.Meta()
				sdkCache[key] = sdkMetadata{credentials: signature, meta: ps.Meta}
			}
		}

		return ps, nil
	}
}

func sdkProviderConfigKey(mg resource.Managed) string {
	if modern, ok := mg.(resource.ModernManaged); ok {
		ref := modern.GetProviderConfigReference()
		namespace := mg.GetNamespace()
		if ref.Kind == "ClusterProviderConfig" {
			namespace = ""
		}
		return ref.Kind + "/" + namespace + "/" + ref.Name
	}
	return "legacy/" + mg.(resource.LegacyManaged).GetProviderConfigReference().Name
}

func toSharedPCSpec(pc *clusterv1beta1.ProviderConfig) (*namespacedv1beta1.ProviderConfigSpec, error) {
	if pc == nil {
		return nil, nil
	}
	data, err := json.Marshal(pc.Spec)
	if err != nil {
		return nil, err
	}
	var spec namespacedv1beta1.ProviderConfigSpec
	if err := json.Unmarshal(data, &spec); err != nil {
		return nil, err
	}
	return &spec, nil
}

func resolveProviderConfig(ctx context.Context, crClient client.Client, mg resource.Managed) (*namespacedv1beta1.ProviderConfigSpec, error) {
	switch m := mg.(type) {
	case resource.LegacyManaged:
		return resolveLegacy(ctx, crClient, m)
	case resource.ModernManaged:
		return resolveModern(ctx, crClient, m)
	default:
		return nil, errors.New("resource is not a managed resource")
	}
}

func resolveLegacy(ctx context.Context, crClient client.Client, mg resource.LegacyManaged) (*namespacedv1beta1.ProviderConfigSpec, error) {
	configRef := mg.GetProviderConfigReference()
	if configRef == nil {
		return nil, errors.New(errNoProviderConfig)
	}
	pc := &clusterv1beta1.ProviderConfig{}
	if err := crClient.Get(ctx, types.NamespacedName{Name: configRef.Name}, pc); err != nil {
		return nil, errors.Wrap(err, errGetProviderConfig)
	}

	t := resource.NewLegacyProviderConfigUsageTracker(crClient, &clusterv1beta1.ProviderConfigUsage{})
	if err := t.Track(ctx, mg); err != nil {
		return nil, errors.Wrap(err, errTrackUsage)
	}

	return toSharedPCSpec(pc)
}

func resolveModern(ctx context.Context, crClient client.Client, mg resource.ModernManaged) (*namespacedv1beta1.ProviderConfigSpec, error) {
	configRef := mg.GetProviderConfigReference()
	if configRef == nil {
		return nil, errors.New(errNoProviderConfig)
	}

	pcRuntimeObj, err := crClient.Scheme().New(namespacedv1beta1.SchemeGroupVersion.WithKind(configRef.Kind))
	if err != nil {
		return nil, errors.Wrap(err, "unknown GVK for ProviderConfig")
	}
	pcObj, ok := pcRuntimeObj.(client.Object)
	if !ok {
		return nil, errors.New("provider config type is not a client.Object")
	}

	if err := crClient.Get(ctx, types.NamespacedName{Name: configRef.Name, Namespace: mg.GetNamespace()}, pcObj); err != nil {
		return nil, errors.Wrap(err, errGetProviderConfig)
	}

	var pcSpec namespacedv1beta1.ProviderConfigSpec
	pcu := &namespacedv1beta1.ProviderConfigUsage{}
	switch pc := pcObj.(type) {
	case *namespacedv1beta1.ProviderConfig:
		pcSpec = pc.Spec
		if pcSpec.Credentials.SecretRef != nil {
			pcSpec.Credentials.SecretRef.Namespace = mg.GetNamespace()
		}
	case *namespacedv1beta1.ClusterProviderConfig:
		pcSpec = pc.Spec
	default:
		return nil, errors.New("unknown provider config type")
	}

	t := resource.NewProviderConfigUsageTracker(crClient, pcu)
	if err := t.Track(ctx, mg); err != nil {
		return nil, errors.Wrap(err, errTrackUsage)
	}
	return &pcSpec, nil
}
