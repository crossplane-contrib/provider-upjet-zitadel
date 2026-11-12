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

### Organization login-policy bootstrap

Login policies use their organization ID as their singleton resource ID. Before
initial observation, the provider now uses an explicit `org_id` (including a
resolved organization reference), or discovers the authenticated organization
through a read-only Management API request for PAT-backed provider configs.
Existing imported external names remain authoritative. The discovered ID also
fills the Terraform `org_id` parameter so a successful creation returns an ID.
This allows normal observation to adopt a pre-existing policy after a partial
creation, rather than repeatedly issuing AddCustomLoginPolicy.

Existing policies and MFA settings are not replaced by this discovery. Provider
transport headers preserve virtual-instance routing; redirects are rejected and
error responses/credentials are not included in errors. Other authentication
methods retain their previous discovery behavior and can supply an explicit
organization ID/reference or import an external name.

`go test ./config` covers retained IDs, explicit organization IDs, PAT discovery,
repeat lookup, scoped transport headers and sanitized failures. The optional
`TestLiveLoginPolicyOrganization` reads test configuration from stdin only and
checks read-only discovery against an expected retained policy ID. It does not
claim a complete live create/update/delete lifecycle test.
