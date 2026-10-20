# Live native binding tests

These opt-in tests exercise the installed controllers and verify their results
independently through Zitadel's Authorization/InternalPermission v2 HTTP APIs.
They cover human/machine and cross-organization assignments, all four
administrator scopes, full role replacement with unchanged identity,
Observe-only imports, server-side authorization drift correction, and remote
deletion. Fixtures include one shared delegation, owned separately from bindings.
The ProjectV2/Application case additionally checks a project rename and public
OIDC application create/update/import/delete, preserving application/client IDs.
It does not verify confidential-client secret restoration, SAML or API apps.

Requirements: Python 3 (standard library only), kubectl, a running Zitadel 4.x
instance, and a healthy provider built with these new native resources. Existing
cluster-scoped ProviderConfig `default` must use the same administrative
credentials/instance as the API endpoint. The referenced Secret's `credentials`
key must hold provider JSON containing `access_token` and optional
`transport_headers`; the test reads it only in memory.

```sh
export ZITADEL_E2E_CONTEXT=<explicit-test-context>
export ZITADEL_E2E_API_URL=https://id.example.com
export ZITADEL_E2E_CREDENTIALS_SECRET=<namespace>/<provider-secret>
# Optional: ZITADEL_E2E_PROVIDER_CONFIG, ZITADEL_E2E_CA_FILE
make test-e2e-bindings
```

The Make target fails if the explicit opt-in settings are absent. Direct
unittest discovery reports a skip when they are absent; a skip is not live
validation. All fixture namespaces and names are generated per run. Cleanup
deletes bindings before delegation, roles, project, users and organizations;
cleanup errors fail the suite and retain the namespace for diagnosis. No API
credentials are printed or written to files. HTTP can be used explicitly for
a loopback port-forward; HTTPS verifies certificates by default.

This existing-cluster suite does not provision a cluster/package or replace the
separately permission-gated Uptest `make e2e` workflow. Uptest uses its own
credentials and example selection. Passing package installation, schema checks,
unit/fake-client tests or this suite without opt-in is not a passing live e2e run.
