# SPDX-FileCopyrightText: 2026 The Crossplane Authors
# SPDX-License-Identifier: Apache-2.0

"""Native managed-resource lifecycle, verified independently through Zitadel APIs."""
import unittest
import urllib.error
from live import LiveSuite


class NativeBindings(LiveSuite):
    def test_project_and_application(self):
        identifier = self.remote_id(self.project)
        project = self.api('project.v2.ProjectService', 'GetProject', {'projectId': identifier})['project']
        self.assertEqual(project['projectId'], identifier)
        self.assertEqual(project['organizationId'], self.remote_id(self.owner))
        updated = self.managed('ProjectV2', 'product', 'project', {
            'orgIdRef': {'name': 'owner'}, 'name': 'Updated E2E product'})
        self.assertEqual(self.remote_id(updated), identifier)
        self.assertEqual(self.api('project.v2.ProjectService', 'GetProject', {'projectId': identifier})['project']['name'], 'Updated E2E product')
        fp = {'name': 'E2E public client', 'orgIdRef': {'name': 'owner'}, 'projectIdRef': {'name': 'product'},
            'oidc': [{'appType': 'OIDC_APP_TYPE_USER_AGENT', 'authMethodType': 'OIDC_AUTH_METHOD_TYPE_NONE',
                'grantTypes': ['OIDC_GRANT_TYPE_AUTHORIZATION_CODE'], 'responseTypes': ['OIDC_RESPONSE_TYPE_CODE'],
                'redirectUris': ['https://example.com/callback'], 'version': 'OIDC_VERSION_1_0'}]}
        x = self.managed('Application', 'public-client', 'application', fp)
        app_id = self.remote_id(x)
        app = self.api('application.v2.ApplicationService', 'GetApplication', {'applicationId': app_id})['application']
        self.assertEqual(app['projectId'], identifier)
        client_id = app['oidcConfiguration']['clientId']
        fp['name'] = 'Updated E2E public client'
        fp['oidc'][0]['redirectUris'] = ['https://example.com/updated-callback']
        x = self.managed('Application', 'public-client', 'application', fp)
        self.assertEqual(self.remote_id(x), app_id)
        app = self.api('application.v2.ApplicationService', 'GetApplication', {'applicationId': app_id})['application']
        self.assertEqual(app['name'], fp['name'])
        self.assertEqual(app['oidcConfiguration']['clientId'], client_id)
        self.assertEqual(app['oidcConfiguration']['redirectUris'], fp['oidc'][0]['redirectUris'])
        imported = self.managed('Application', 'client-import', 'application', fp,
                                policies=['Observe'], external_name=app_id)
        self.assertEqual(self.remote_id(imported), app_id)
        resource = 'applications.application.zitadel.m.crossplane.io'
        self.delete(resource, 'client-import')
        self.delete(resource, 'public-client')
        with self.assertRaises(urllib.error.HTTPError) as missing:
            self.api('application.v2.ApplicationService', 'GetApplication', {'applicationId': app_id})
        self.assertEqual(missing.exception.code, 404)

    def test_authorization_lifecycle(self):
        resource = 'authorizations.authorization.zitadel.m.crossplane.io'
        for name, user, org in [('human', self.human, self.owner),
                                ('machine', self.machine, self.owner),
                                ('delegated', self.machine, self.grantee)]:
            with self.subTest(subject=name):
                fp = {'userId': self.remote_id(user), 'projectId': self.remote_id(self.project),
                      'orgId': self.remote_id(org), 'roleKeys': ['reader']}
                x = self.managed('Authorization', name, 'authorization', fp)
                identifier = self.remote_id(x)
                self.assert_authorization(identifier, fp, ['reader'])
                fp['roleKeys'] = ['writer']
                x = self.managed('Authorization', name, 'authorization', fp)
                self.assertEqual(self.remote_id(x), identifier)
                self.assert_authorization(identifier, fp, ['writer'])
                # Server-side drift must be repaired, retaining the remote ID.
                self.api('authorization.v2.AuthorizationService', 'UpdateAuthorization',
                         {'id': identifier, 'roleKeys': ['reader']})
                self.kube('annotate', resource, name, '-n', self.namespace,
                          'crossplane.io/reconcile=verify-remote-drift', '--overwrite')
                self.wait(lambda: {r['key'] for a in self.authorizations(identifier) for r in a['roles']} == {'writer'}, 'remote drift correction')
                observed = self.managed('Authorization', name + '-import', 'authorization', fp,
                                        policies=['Observe'], external_name=identifier)
                self.assertEqual(self.remote_id(observed), identifier)
                self.delete(resource, name + '-import')
                self.assert_authorization(identifier, fp, ['writer'])
                self.delete(resource, name)
                self.assertEqual(self.authorizations(identifier), [])

    def test_administrator_lifecycle(self):
        resource = 'administrators.internalpermission.zitadel.m.crossplane.io'
        cases = [
            ('instance', {}, 'IAM_OWNER_VIEWER', 'IAM_USER_MANAGER'),
            ('organization', {'orgId': self.remote_id(self.owner)}, 'ORG_USER_MANAGER', 'ORG_OWNER_VIEWER'),
            ('project', {'projectId': self.remote_id(self.project)}, 'PROJECT_OWNER_VIEWER', 'PROJECT_OWNER'),
            ('project_grant', {'projectId': self.remote_id(self.project), 'orgId': self.remote_id(self.grantee)}, 'PROJECT_GRANT_OWNER_VIEWER', 'PROJECT_GRANT_OWNER')]
        for scope, target, first, second in cases:
            with self.subTest(scope=scope):
                name = scope.replace('_', '-')
                fp = dict(target, userId=self.remote_id(self.machine), resourceType=scope, roles=[first])
                x = self.managed('Administrator', name, 'internalpermission', fp)
                identifier = self.remote_id(x)
                self.assert_administrator(fp, [first])
                fp['roles'] = [second]
                x = self.managed('Administrator', name, 'internalpermission', fp)
                self.assertEqual(self.remote_id(x), identifier)
                self.assert_administrator(fp, [second])
                imported = self.managed('Administrator', name + '-import', 'internalpermission', fp,
                                       policies=['Observe'], external_name=identifier)
                self.assertEqual(self.remote_id(imported), identifier)
                self.delete(resource, name + '-import')
                self.assert_administrator(fp, [second])
                self.delete(resource, name)
                self.assertEqual(self.administrators(fp), [])


if __name__ == '__main__':
    unittest.main(verbosity=2)
