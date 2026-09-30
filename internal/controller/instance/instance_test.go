// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package instance

import (
	"context"
	"errors"
	"testing"

	xpv1 "github.com/crossplane/crossplane-runtime/v2/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/v2/pkg/meta"
	"github.com/crossplane/crossplane-runtime/v2/pkg/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	cluster "github.com/crossplane-contrib/provider-upjet-zitadel/apis/cluster/instance/v1alpha1"
	namespaced "github.com/crossplane-contrib/provider-upjet-zitadel/apis/namespaced/instance/v1alpha1"
	"github.com/crossplane-contrib/provider-upjet-zitadel/apis/observation"
	"github.com/crossplane-contrib/provider-upjet-zitadel/internal/clients"
)

type fakeAPI struct {
	value         observation.Instance
	err           error
	reads, closes int
}

func (a *fakeAPI) Observe(context.Context) (observation.Instance, error) {
	a.reads++
	return a.value, a.err
}
func (a *fakeAPI) Close() error { a.closes++; return nil }

func testClient(t *testing.T, objects ...client.Object) client.Client {
	t.Helper()
	scheme := runtime.NewScheme()
	if err := namespaced.AddToScheme(scheme); err != nil {
		t.Fatal(err)
	}
	if err := cluster.AddToScheme(scheme); err != nil {
		t.Fatal(err)
	}
	return fake.NewClientBuilder().WithScheme(scheme).WithObjects(objects...).Build()
}

func TestObservationLifecycle(t *testing.T) {
	for _, scope := range []string{"namespaced", "cluster"} {
		t.Run(scope, func(t *testing.T) {
			var inst observation.ObservedInstance = &namespaced.Instance{ObjectMeta: metav1.ObjectMeta{Name: "existing", Namespace: "tenant"}}
			if scope == "cluster" {
				inst = &cluster.Instance{ObjectMeta: metav1.ObjectMeta{Name: "existing"}}
			}
			c := testClient(t, inst)
			api := &fakeAPI{value: observation.Instance{ID: "external-instance", Name: "Helm instance", Version: "4", State: "STATE_RUNNING"}}
			e := &external{kube: c, api: api}
			for n := 0; n < 2; n++ {
				o, err := e.Observe(context.Background(), inst)
				if err != nil || !o.ResourceExists || !o.ResourceUpToDate {
					t.Fatalf("observation failed: %+v, %v", o, err)
				}
			}
			if inst.GetInstanceObservation() != api.value || meta.GetExternalName(inst) != "external-instance" {
				t.Fatal("observation was not recorded/bound")
			}
			inst.SetConditions(xpv1.ReconcileSuccess())
			if observation.InstanceID()(inst) != "external-instance" {
				t.Fatal("healthy observation is not referenceable")
			}
			if _, err := e.Create(context.Background(), inst); err == nil {
				t.Fatal("creation must be rejected")
			}
			if _, err := e.Update(context.Background(), inst); err == nil {
				t.Fatal("update must be rejected")
			}
			before := api.reads
			if _, err := e.Delete(context.Background(), inst); err != nil {
				t.Fatal(err)
			}
			if api.reads != before {
				t.Fatal("deletion called the remote API")
			}
			if err := e.Disconnect(context.Background()); err != nil || api.closes != 1 {
				t.Fatal("connection not closed")
			}
		})
	}
}

func TestObservationFailureAndIdentityChange(t *testing.T) {
	for _, failure := range []string{"empty ID", "unavailable", "different instance"} {
		t.Run(failure, func(t *testing.T) {
			inst := &namespaced.Instance{ObjectMeta: metav1.ObjectMeta{Name: "existing", Namespace: "tenant"}}
			meta.SetExternalName(inst, "original")
			inst.Status.AtProvider = observation.Instance{ID: "original"}
			inst.SetConditions(xpv1.Available(), xpv1.ReconcileSuccess())
			api := &fakeAPI{}
			if failure == "unavailable" {
				api.err = errors.New("API unavailable")
			}
			if failure == "different instance" {
				api.value.ID = "replacement"
			}
			e := &external{kube: testClient(t, inst), api: api}
			if _, err := e.Observe(context.Background(), inst); err == nil {
				t.Fatal("expected failure")
			}
			if meta.GetExternalName(inst) != "original" || inst.Status.AtProvider.ID != "original" {
				t.Fatal("bound identity changed")
			}
			if observation.InstanceID()(inst) != "" {
				t.Fatal("failed observation must not resolve references")
			}
			// API availability recovers without accepting another identity.
			api.err = nil
			api.value.ID = "original"
			if _, err := e.Observe(context.Background(), inst); err != nil {
				t.Fatal(err)
			}
			if observation.InstanceID()(inst) != "original" {
				t.Fatal("observation did not recover")
			}
		})
	}
}

func TestConnectRejectsMutatingPoliciesBeforeCredentials(t *testing.T) {
	for _, policy := range []xpv1.ManagementAction{xpv1.ManagementActionAll, xpv1.ManagementActionCreate, xpv1.ManagementActionDelete, xpv1.ManagementActionUpdate, xpv1.ManagementActionLateInitialize} {
		t.Run(string(policy), func(t *testing.T) {
			called := false
			c := &connector{newClient: func(context.Context, client.Client, resource.Managed) (clients.InstanceClient, error) {
				called = true
				return nil, nil
			}}
			inst := &namespaced.Instance{}
			inst.Spec.ManagementPolicies = xpv1.ManagementPolicies{policy}
			if _, err := c.Connect(context.Background(), inst); err == nil || called {
				t.Fatal("mutating policy reached the credential/client path")
			}
		})
	}
}
