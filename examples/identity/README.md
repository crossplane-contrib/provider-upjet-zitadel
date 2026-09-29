# Identity references and generated credentials

`labels-and-vault.yaml` uses native Crossplane selectors to connect an
Organization, MachineUser and AccessToken. No Zitadel ID needs to be copied into
another manifest. The provider generates the PAT and publishes it to a Kubernetes
connection Secret; ESO PushSecret persists it to Vault. The provider does not need
Vault credentials. Change the SecretStore and destination to use another ESO
backend.

Use labels that uniquely identify the intended dependency. Standard Crossplane
selector semantics apply, including `matchControllerRef` and resolution policies.
A namespaced selector defaults to the managed resource's namespace; use an
explicit namespace when intentionally referencing another namespace. Existing
literal IDs remain supported. Organization references target the current
`org.zitadel.m.crossplane.io/Organization` resource, not the legacy `Org` resource.
AccessToken and Key user references target MachineUser.

ESO requires the PushSecret CRD and write permission for the remote destination.
`deletionPolicy: None` retains the Vault value on deletion; choose `Delete` only
when removal of the remote credential is intended. This example does not manage
the SecretStore or namespace. In a composition, derive the connection Secret name
and PushSecret source from the same value rather than asking consumers to copy it.
