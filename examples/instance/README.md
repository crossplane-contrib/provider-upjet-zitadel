# Observe an installed instance

Zitadel's Helm chart creates the first instance during setup. The provider-native
`Instance` resource observes the current instance through the Admin API using the
referenced ProviderConfig. It needs no externally assigned ID and has no create,
update or delete operations. The observation is periodically refreshed by the
provider's normal poll interval. The namespace-scoped and legacy cluster-scoped
APIs expose the same behavior.

`spec.managementPolicies` defaults to `[Observe]`; `[]` pauses reconciliation.
Admission rejects other policies, and the connector also checks them before
loading credentials. Deleting the resource only removes the Kubernetes
observation, even when the external service is unavailable. It cannot delete the
Helm release, database or remote instance.

The endpoint, authentication and transport settings use the existing
ProviderConfig credential JSON. Its credential must be authorized to read the
current instance (`GetMyInstance`); a normal application user is insufficient.
The observer reuses the pinned Terraform provider's authentication builder for
PAT, JWT profile and System API credentials. It does not use Terraform's
`zitadel_instance` data source, because that data source requires the ID before
reading. This controller adds the discovery step above the existing Terraform
resource controllers.

JWT-profile discovery and token exchange use the reconcile context, including
cancellation, and a dedicated HTTP client with a 30-second fallback timeout.
Other authentication methods retain the Terraform helper's behavior.

`GetMyInstance` is a deprecated Admin API in the current Zitadel documentation;
it remains available in the pinned SDK and was verified against Zitadel v4.13.1.
Using it keeps discovery compatible with the provider's existing authentication
stack. This change does not migrate the Terraform-backed domain APIs to v2.

Only non-secret ID/name/version/state metadata is published in
`status.atProvider`. The discovered ID is also bound in `crossplane.io/external-name`.
An already bound resource refuses to switch to another ID if credentials, routing
or the database change. Adopt a replacement intentionally by declaring a new
observer and then updating consumers after checking their replacement behavior.
An explicitly supplied external-name annotation is treated as an expected ID to
verify, not as proof of an observation.

`CustomDomain` and `TrustedDomain` accept `instanceIdRef` and
`instanceIdSelector`. They extract typed status only from a Ready, Synced, current
generation observation that is not deleting. Missing/unready targets cause normal
reference-resolution retries before Terraform calls. Existing explicit
`instanceId` values retain the standard default `IfNotPresent` behavior. Use
`policy.resolve: Always` when intentionally following an observation; changing an
instance ID can replace a domain resource.

Crossplane persists the resolved ID/reference in the live domain spec. Seeing
`instanceId` alongside a reference or selector is expected; only the reference
or selector needs to be authored in the input manifest.

Namespaced references/selectors default to the consumer's namespace and support
explicit namespace targeting. Cluster-scoped consumers reference cluster-scoped
observers. Labels must uniquely identify the intended observation; standard
Crossplane selectors do not enforce uniqueness. The credentials used by the
observer and a domain consumer may differ (IAM admin versus System API), but must
target the same instance. References do not validate endpoint equality.

`observation.yaml` shows a namespaced observer, a named custom-domain reference and
a labelled trusted-domain selector. Set actual domain names and provide the
matching administrator/System API ProviderConfigs before applying it. Generated
credentials remain in Kubernetes connection Secrets and can be exported through
ESO PushSecrets as in `../identity/labels-and-vault.yaml`.

For deletion, remove dependent domain resources and wait for their finalizers
before removing the observer or its namespace. This preserves provider access
while remote domain cleanup runs.
