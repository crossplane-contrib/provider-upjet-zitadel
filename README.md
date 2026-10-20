# Provider Upjet Zitadel

`provider-upjet-zitadel` is a [Crossplane](https://crossplane.io/) provider that
is built using [Upjet](https://github.com/crossplane/upjet) code
generation tools and exposes XRM-conformant managed resources for the
[Zitadel](https://zitadel.com/) API.

## Getting Started

Install the provider by using the following command after changing the image tag
to the [latest release](https://marketplace.upbound.io/providers/crossplane-contrib/provider-upjet-zitadel):
```
up ctp provider install crossplane-contrib/provider-upjet-zitadel:v0.1.0
```

Alternatively, you can use declarative installation:
```
cat <<EOF | kubectl apply -f -
---
apiVersion: pkg.crossplane.io/v1
kind: Provider
metadata:
  name: provider-upjet-zitadel
spec:
  package: xpkg.upbound.io/crossplane-contrib/provider-upjet-zitadel:v0.1.0
EOF
```

You can see the API reference [here](https://doc.crds.dev/github.com/crossplane-contrib/provider-upjet-zitadel).

## Developing

Run code-generation pipeline:
```console
go run cmd/generator/main.go "$PWD"
```

Run against a Kubernetes cluster:

```console
make run
```

Build, push, and install:

```console
make all
```

Build binary:

```console
make build
```

## Report a Bug

For filing bugs, suggesting improvements, or requesting new features, please
open an [issue](https://github.com/crossplane-contrib/provider-upjet-zitadel/issues).

## Instance observation

For Helm-bootstrapped installations, a provider-native read-only `Instance` can
discover the existing instance ID. `CustomDomain` and `TrustedDomain` can then
use native `instanceIdRef`/`instanceIdSelector` dependencies without a discovery
Job or injected IDs. See [the example and lifecycle guidance](examples/instance/README.md).

### Stable Zitadel v2 resources (pending upstream release)

The following additive resources use the embedded Terraform SDK provider and
stable Zitadel v2 services. Existing controllers continue to use their pinned
Terraform CLI backend and retain their Kubernetes API kinds.

| Namespaced resource | Zitadel service |
| --- | --- |
| `project.zitadel.m.crossplane.io/ProjectV2` | Project v2 |
| `application.zitadel.m.crossplane.io/Application` | Application v2 (OIDC/SAML/API) |
| `authorization.zitadel.m.crossplane.io/Authorization` | Authorization v2 |
| `internalpermission.zitadel.m.crossplane.io/Administrator` | InternalPermission v2 |

Each uses Kubernetes API version `v1alpha1`; that version is independent of the
underlying Zitadel service version. Project references select ProjectV2 and
organization references select Organization. User IDs remain literal because
they can identify a human or machine user; a relationship composition can resolve
that union without choosing an incorrect native reference target.

Authorization owns the complete role-key set for one user/project/organization
assignment and imports by the existing authorization ID. Administrator owns the
complete administrator role set for one explicit instance, organization, project
or project-grant target. Its identity is `user:instance`,
`user:organization:org`, `user:project:project`, or
`user:project_grant:project:grantee-org`; percent-encode colons in custom IDs.
Use only one owner per assignment. Adopt with Observe-only management policies,
verify observed identities and roles, then enable writes after the former owner
has relinquished the assignment without deletion.

This branch pins the additive Terraform implementation from
[upstream PR474](https://github.com/zitadel/terraform-provider-zitadel/pull/474)
through an exact fork commit. It is a development dependency, **not a released
upstream capability**. Replace that fork pin with the reviewed upstream release
before merging/releasing this provider. `make generate` combines the pinned
legacy CLI schema/docs with only these four resources from the exact embedded
SDK dependency. The generated SDK controllers use operation stores and support
management policies without invoking the legacy Terraform binary.

Existing legacy role, delegation, PAT and Member resources retain their existing
service implementations. Exposing these four v2 resources does not migrate
those operations or convert Actions v1 hooks to Actions v2.
