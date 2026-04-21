/*
Copyright 2024 The Crossplane Authors.
*/

package config

import (
	// Note(turkenh): we are importing this to embed provider schema document
	_ "embed"

	ujconfig "github.com/crossplane/upjet/v2/pkg/config"
)

const (
	resourcePrefix = "zitadel"
	modulePath     = "github.com/crossplane-contrib/provider-upjet-zitadel"
)

//go:embed schema.json
var providerSchema string

//go:embed provider-metadata.yaml
var providerMetadata string

// ExternalNameConfigs contains all external name configurations for this
// provider.
var ExternalNameConfigs = map[string]ujconfig.ExternalName{
	// Actions
	"zitadel_action":                      ujconfig.IdentifierFromProvider,
	"zitadel_action_execution_event":      ujconfig.IdentifierFromProvider,
	"zitadel_action_execution_function":   ujconfig.IdentifierFromProvider,
	"zitadel_action_execution_request":    ujconfig.IdentifierFromProvider,
	"zitadel_action_execution_response":   ujconfig.IdentifierFromProvider,
	"zitadel_action_target":               ujconfig.IdentifierFromProvider,
	"zitadel_action_target_public_key":    ujconfig.IdentifierFromProvider,
	"zitadel_trigger_actions":             ujconfig.IdentifierFromProvider,

	// Applications
	"zitadel_application_api":  ujconfig.IdentifierFromProvider,
	"zitadel_application_key":  ujconfig.IdentifierFromProvider,
	"zitadel_application_oidc": ujconfig.IdentifierFromProvider,
	"zitadel_application_saml": ujconfig.IdentifierFromProvider,
	"zitadel_active_webkey":    ujconfig.IdentifierFromProvider,
	"zitadel_webkey":           ujconfig.IdentifierFromProvider,

	// Default (instance-level) policies & settings
	"zitadel_default_domain_claimed_message_text":            ujconfig.IdentifierFromProvider,
	"zitadel_default_domain_policy":                          ujconfig.IdentifierFromProvider,
	"zitadel_default_init_message_text":                      ujconfig.IdentifierFromProvider,
	"zitadel_default_invite_user_message_text":               ujconfig.IdentifierFromProvider,
	"zitadel_default_label_policy":                           ujconfig.IdentifierFromProvider,
	"zitadel_default_lockout_policy":                         ujconfig.IdentifierFromProvider,
	"zitadel_default_login_policy":                           ujconfig.IdentifierFromProvider,
	"zitadel_default_login_texts":                            ujconfig.IdentifierFromProvider,
	"zitadel_default_notification_policy":                    ujconfig.IdentifierFromProvider,
	"zitadel_default_oidc_settings":                          ujconfig.IdentifierFromProvider,
	"zitadel_default_password_age_policy":                    ujconfig.IdentifierFromProvider,
	"zitadel_default_password_change_message_text":           ujconfig.IdentifierFromProvider,
	"zitadel_default_password_complexity_policy":             ujconfig.IdentifierFromProvider,
	"zitadel_default_password_reset_message_text":            ujconfig.IdentifierFromProvider,
	"zitadel_default_passwordless_registration_message_text": ujconfig.IdentifierFromProvider,
	"zitadel_default_privacy_policy":                         ujconfig.IdentifierFromProvider,
	"zitadel_default_security_settings":                      ujconfig.IdentifierFromProvider,
	"zitadel_default_verify_email_message_text":              ujconfig.IdentifierFromProvider,
	"zitadel_default_verify_email_otp_message_text":          ujconfig.IdentifierFromProvider,
	"zitadel_default_verify_phone_message_text":              ujconfig.IdentifierFromProvider,
	"zitadel_default_verify_sms_otp_message_text":            ujconfig.IdentifierFromProvider,

	// Domains & Organizations
	"zitadel_domain":                      ujconfig.IdentifierFromProvider,
	"zitadel_domain_claimed_message_text": ujconfig.IdentifierFromProvider,
	"zitadel_domain_policy":               ujconfig.IdentifierFromProvider,
	"zitadel_org":                         ujconfig.IdentifierFromProvider,
	"zitadel_org_member":                  ujconfig.IdentifierFromProvider,
	"zitadel_org_metadata":                ujconfig.IdentifierFromProvider,
	"zitadel_organization":                ujconfig.IdentifierFromProvider,
	"zitadel_organization_domain":         ujconfig.IdentifierFromProvider,
	"zitadel_organization_metadata":       ujconfig.IdentifierFromProvider,

	// Email & SMS providers
	"zitadel_email_provider_http":  ujconfig.IdentifierFromProvider,
	"zitadel_email_provider_smtp":  ujconfig.IdentifierFromProvider,
	"zitadel_sms_provider_http":    ujconfig.IdentifierFromProvider,
	"zitadel_sms_provider_twilio":  ujconfig.IdentifierFromProvider,
	"zitadel_smtp_config":          ujconfig.IdentifierFromProvider,

	// Identity Providers — Instance-Level
	"zitadel_idp_apple":              ujconfig.IdentifierFromProvider,
	"zitadel_idp_azure_ad":           ujconfig.IdentifierFromProvider,
	"zitadel_idp_github":             ujconfig.IdentifierFromProvider,
	"zitadel_idp_github_es":          ujconfig.IdentifierFromProvider,
	"zitadel_idp_gitlab":             ujconfig.IdentifierFromProvider,
	"zitadel_idp_gitlab_self_hosted": ujconfig.IdentifierFromProvider,
	"zitadel_idp_google":             ujconfig.IdentifierFromProvider,
	"zitadel_idp_ldap":               ujconfig.IdentifierFromProvider,
	"zitadel_idp_oauth":              ujconfig.IdentifierFromProvider,
	"zitadel_idp_oidc":               ujconfig.IdentifierFromProvider,
	"zitadel_idp_saml":               ujconfig.IdentifierFromProvider,

	// Identity Providers — Org-Level
	"zitadel_org_idp_apple":              ujconfig.IdentifierFromProvider,
	"zitadel_org_idp_azure_ad":           ujconfig.IdentifierFromProvider,
	"zitadel_org_idp_github":             ujconfig.IdentifierFromProvider,
	"zitadel_org_idp_github_es":          ujconfig.IdentifierFromProvider,
	"zitadel_org_idp_gitlab":             ujconfig.IdentifierFromProvider,
	"zitadel_org_idp_gitlab_self_hosted": ujconfig.IdentifierFromProvider,
	"zitadel_org_idp_google":             ujconfig.IdentifierFromProvider,
	"zitadel_org_idp_jwt":               ujconfig.IdentifierFromProvider,
	"zitadel_org_idp_ldap":               ujconfig.IdentifierFromProvider,
	"zitadel_org_idp_oauth":              ujconfig.IdentifierFromProvider,
	"zitadel_org_idp_oidc":               ujconfig.IdentifierFromProvider,
	"zitadel_org_idp_saml":               ujconfig.IdentifierFromProvider,

	// Instance Management
	"zitadel_instance_custom_domain":    ujconfig.IdentifierFromProvider,
	"zitadel_instance_features":         ujconfig.IdentifierFromProvider,
	"zitadel_instance_member":           ujconfig.IdentifierFromProvider,
	"zitadel_instance_restrictions":     ujconfig.IdentifierFromProvider,
	"zitadel_instance_secret_generator": ujconfig.IdentifierFromProvider,
	"zitadel_instance_trusted_domain":   ujconfig.IdentifierFromProvider,

	// Org-Level Policies & Message Texts
	"zitadel_init_message_text":                      ujconfig.IdentifierFromProvider,
	"zitadel_label_policy":                           ujconfig.IdentifierFromProvider,
	"zitadel_lockout_policy":                         ujconfig.IdentifierFromProvider,
	"zitadel_login_policy":                           ujconfig.IdentifierFromProvider,
	"zitadel_login_texts":                            ujconfig.IdentifierFromProvider,
	"zitadel_notification_policy":                    ujconfig.IdentifierFromProvider,
	"zitadel_password_age_policy":                    ujconfig.IdentifierFromProvider,
	"zitadel_password_change_message_text":           ujconfig.IdentifierFromProvider,
	"zitadel_password_complexity_policy":             ujconfig.IdentifierFromProvider,
	"zitadel_password_reset_message_text":            ujconfig.IdentifierFromProvider,
	"zitadel_passwordless_registration_message_text": ujconfig.IdentifierFromProvider,
	"zitadel_privacy_policy":                         ujconfig.IdentifierFromProvider,
	"zitadel_verify_email_message_text":              ujconfig.IdentifierFromProvider,
	"zitadel_verify_email_otp_message_text":          ujconfig.IdentifierFromProvider,
	"zitadel_verify_phone_message_text":              ujconfig.IdentifierFromProvider,
	"zitadel_verify_sms_otp_message_text":            ujconfig.IdentifierFromProvider,

	// Projects
	"zitadel_project":              ujconfig.IdentifierFromProvider,
	"zitadel_project_grant":        ujconfig.IdentifierFromProvider,
	"zitadel_project_grant_member": ujconfig.IdentifierFromProvider,
	"zitadel_project_member":       ujconfig.IdentifierFromProvider,
	"zitadel_project_role":         ujconfig.IdentifierFromProvider,

	// Users
	"zitadel_human_user":            ujconfig.IdentifierFromProvider,
	"zitadel_machine_key":           ujconfig.IdentifierFromProvider,
	"zitadel_machine_user":          ujconfig.IdentifierFromProvider,
	"zitadel_personal_access_token": ujconfig.IdentifierFromProvider,
	"zitadel_user_grant":            ujconfig.IdentifierFromProvider,
	"zitadel_user_metadata":         ujconfig.IdentifierFromProvider,

	// System
	"zitadel_system_features": ujconfig.IdentifierFromProvider,
}

func newProvider(rootGroup string) *ujconfig.Provider {
	pc := ujconfig.NewProvider([]byte(providerSchema), resourcePrefix, modulePath, []byte(providerMetadata),
		ujconfig.WithRootGroup(rootGroup),
		ujconfig.WithShortName("zitadel"),
		ujconfig.WithIncludeList(ExternalNameConfigured()),
		ujconfig.WithFeaturesPackage("internal/features"),
		ujconfig.WithExampleManifestConfiguration(ujconfig.ExampleManifestConfiguration{
			ManagedResourceNamespace: "crossplane-system",
		}),
		ujconfig.WithDefaultResourceOptions(
			ExternalNameConfigurations(),
		))

	// Applications reference a project
	pc.AddResourceConfigurator("zitadel_application_api", func(r *ujconfig.Resource) {
		r.ShortGroup = "application"
		r.References["project_id"] = ujconfig.Reference{
			TerraformName: "zitadel_project",
		}
	})
	pc.AddResourceConfigurator("zitadel_application_oidc", func(r *ujconfig.Resource) {
		r.ShortGroup = "application"
		r.References["project_id"] = ujconfig.Reference{
			TerraformName: "zitadel_project",
		}
	})
	pc.AddResourceConfigurator("zitadel_application_saml", func(r *ujconfig.Resource) {
		r.ShortGroup = "application"
		r.References["project_id"] = ujconfig.Reference{
			TerraformName: "zitadel_project",
		}
	})
	pc.AddResourceConfigurator("zitadel_application_key", func(r *ujconfig.Resource) {
		r.ShortGroup = "application"
		r.References["project_id"] = ujconfig.Reference{
			TerraformName: "zitadel_project",
		}
	})

	// Projects
	pc.AddResourceConfigurator("zitadel_project", func(r *ujconfig.Resource) {
		r.ShortGroup = "project"
	})
	pc.AddResourceConfigurator("zitadel_project_role", func(r *ujconfig.Resource) {
		r.ShortGroup = "project"
		r.References["project_id"] = ujconfig.Reference{
			TerraformName: "zitadel_project",
		}
	})
	pc.AddResourceConfigurator("zitadel_project_member", func(r *ujconfig.Resource) {
		r.ShortGroup = "project"
		r.References["project_id"] = ujconfig.Reference{
			TerraformName: "zitadel_project",
		}
	})
	pc.AddResourceConfigurator("zitadel_project_grant", func(r *ujconfig.Resource) {
		r.ShortGroup = "project"
		r.References["project_id"] = ujconfig.Reference{
			TerraformName: "zitadel_project",
		}
	})
	pc.AddResourceConfigurator("zitadel_project_grant_member", func(r *ujconfig.Resource) {
		r.ShortGroup = "project"
		r.References["project_id"] = ujconfig.Reference{
			TerraformName: "zitadel_project",
		}
	})

	// Organizations
	pc.AddResourceConfigurator("zitadel_org", func(r *ujconfig.Resource) {
		r.ShortGroup = "org"
	})
	pc.AddResourceConfigurator("zitadel_org_member", func(r *ujconfig.Resource) {
		r.ShortGroup = "org"
		r.References["org_id"] = ujconfig.Reference{
			TerraformName: "zitadel_org",
		}
	})
	pc.AddResourceConfigurator("zitadel_org_metadata", func(r *ujconfig.Resource) {
		r.ShortGroup = "org"
		r.References["org_id"] = ujconfig.Reference{
			TerraformName: "zitadel_org",
		}
	})
	pc.AddResourceConfigurator("zitadel_organization", func(r *ujconfig.Resource) {
		r.ShortGroup = "org"
	})

	// Users
	pc.AddResourceConfigurator("zitadel_human_user", func(r *ujconfig.Resource) {
		r.ShortGroup = "user"
	})
	pc.AddResourceConfigurator("zitadel_machine_user", func(r *ujconfig.Resource) {
		r.ShortGroup = "user"
	})
	pc.AddResourceConfigurator("zitadel_machine_key", func(r *ujconfig.Resource) {
		r.ShortGroup = "user"
	})
	pc.AddResourceConfigurator("zitadel_personal_access_token", func(r *ujconfig.Resource) {
		r.ShortGroup = "user"
	})
	pc.AddResourceConfigurator("zitadel_user_grant", func(r *ujconfig.Resource) {
		r.ShortGroup = "user"
	})
	pc.AddResourceConfigurator("zitadel_user_metadata", func(r *ujconfig.Resource) {
		r.ShortGroup = "user"
	})

	// Instance
	pc.AddResourceConfigurator("zitadel_instance_member", func(r *ujconfig.Resource) {
		r.ShortGroup = "instance"
	})
	pc.AddResourceConfigurator("zitadel_instance_features", func(r *ujconfig.Resource) {
		r.ShortGroup = "instance"
	})
	pc.AddResourceConfigurator("zitadel_instance_custom_domain", func(r *ujconfig.Resource) {
		r.ShortGroup = "instance"
	})
	pc.AddResourceConfigurator("zitadel_instance_restrictions", func(r *ujconfig.Resource) {
		r.ShortGroup = "instance"
	})
	pc.AddResourceConfigurator("zitadel_instance_secret_generator", func(r *ujconfig.Resource) {
		r.ShortGroup = "instance"
	})
	pc.AddResourceConfigurator("zitadel_instance_trusted_domain", func(r *ujconfig.Resource) {
		r.ShortGroup = "instance"
	})

	pc.ConfigureResources()
	return pc
}

// GetProvider returns cluster-scoped provider configuration.
func GetProvider() *ujconfig.Provider {
	return newProvider("zitadel.crossplane.io")
}

// GetProviderNamespaced returns namespaced MR provider configuration.
func GetProviderNamespaced() *ujconfig.Provider {
	return newProvider("zitadel.m.crossplane.io")
}
