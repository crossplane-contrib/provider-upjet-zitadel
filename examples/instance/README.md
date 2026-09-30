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

## Reproduce with Hops from source

Start the local control plane from your Hops cluster configuration first:

```sh
hops local up /path/to/local/cluster.yaml --context kind-hops
```

That configuration must install Crossplane, deploy a working Zitadel instance
(for example through its Helm chart), and configure provider credentials. Wait
for setup and credentials to be ready. Stop the `hops local up` watcher while
leaving the cluster running so its release pin does not overwrite the test build.
If it runs as a service, stop that service instead.

From this provider checkout:

```sh
go test ./...
go vet ./...
hops provider install --path "$PWD" --context kind-hops \
  --cluster-provider kind --docker-provider dory --version-prefix v0.999.4
```

Use your Docker backend instead of `dory` if different. When Dory's socket is not
your default Docker endpoint, set `DOCKER_HOST=unix://$HOME/.dory/dory.sock`.
Check Provider/ProviderRevision health and the running pod's image ID against the
fresh local build; the presence of new CRDs alone does not establish that the new
controller is running.

Edit `observation.yaml` to use a test namespace, unique domain names, an IAM-admin
ProviderConfig for discovery/trusted domains, and a System API ProviderConfig for
custom domains. Apply it and wait for all three resources to become Ready/Synced:

```sh
kubectl --context kind-hops apply -f examples/instance/observation.yaml
kubectl --context kind-hops get instances.instance.zitadel.m.crossplane.io \
  -n crossplane-system -o yaml
kubectl --context kind-hops get customdomains.instance.zitadel.m.crossplane.io,trusteddomains.instance.zitadel.m.crossplane.io \
  -n crossplane-system -o yaml
```

Both domain specs should resolve `instanceId` to the observer's typed ID even
though the input manifest has only a reference or selector. Crossplane persists
the resolved ID/reference in the live spec; this is expected. For an ordering
test, apply the domain documents first, check their unresolved conditions and
absence of external IDs, then apply the observer. They should recover without
manual ID patches. Reapply unchanged manifests and check IDs/generations stay
stable. The legacy API supports the same test using cluster-scoped resources and
legacy `ProviderConfig` references.

Clean up the domains first, wait for their finalizers, then delete the observer.
Recreating only the observer must report the same remote instance ID. Remove the
test observer again, restore your prior released Provider package/runtime config,
remove the source install's ImageConfig/runtime config/RBAC overrides, and restart
the local watcher. Existing identity IDs and credential Secret data should remain
unchanged throughout.
