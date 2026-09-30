// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"context"
	"errors"
	"testing"

	xpv1 "github.com/crossplane/crossplane-runtime/v2/apis/common/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	cluster "github.com/crossplane-contrib/provider-upjet-zitadel/apis/cluster/instance/v1alpha1"
	namespaced "github.com/crossplane-contrib/provider-upjet-zitadel/apis/namespaced/instance/v1alpha1"
	"github.com/crossplane-contrib/provider-upjet-zitadel/apis/observation"
)

func instanceRefClient(t *testing.T, objects ...client.Object) client.Client {
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
func instanceTarget(ns, id string) *namespaced.Instance {
	i := &namespaced.Instance{ObjectMeta: metav1.ObjectMeta{Name: "installed", Namespace: ns, Labels: map[string]string{"example.org/instance": "installed"}}}
	i.Status.AtProvider = observation.Instance{ID: id}
	i.SetConditions(xpv1.Available(), xpv1.ReconcileSuccess())
	return i
}

func TestNamespacedInstanceReferences(t *testing.T) {
	for _, method := range []string{"reference", "selector", "cross namespace"} {
		t.Run(method, func(t *testing.T) {
			target := instanceTarget("tenant", "instance-id")
			c := instanceRefClient(t, target, instanceTarget("other", "wrong-instance"))
			d := &namespaced.CustomDomain{ObjectMeta: metav1.ObjectMeta{Name: "domain", Namespace: "tenant"}}
			switch method {
			case "reference":
				d.Spec.ForProvider.InstanceIDRef = &xpv1.NamespacedReference{Name: "installed"}
			case "selector":
				d.Spec.ForProvider.InstanceIDSelector = &xpv1.NamespacedSelector{MatchLabels: target.Labels}
			case "cross namespace":
				ns := "tenant"
				d.Namespace = "consumer"
				d.Spec.ForProvider.InstanceIDRef = &xpv1.NamespacedReference{Name: "installed", Namespace: ns}
			}
			if err := d.ResolveReferences(context.Background(), c); err != nil {
				t.Fatal(err)
			}
			if d.Spec.ForProvider.InstanceID == nil || *d.Spec.ForProvider.InstanceID != "instance-id" {
				t.Fatal("wrong resolved ID")
			}
		})
	}
}

func TestUnobservedInstanceCannotResolve(t *testing.T) {
	for _, failure := range []string{"missing", "wrong namespace", "empty ID", "not Ready", "not Synced", "stale generation", "deleting"} {
		t.Run(failure, func(t *testing.T) {
			target := instanceTarget("tenant", "instance-id")
			switch failure {
			case "wrong namespace":
				target.Namespace = "other"
			case "empty ID":
				target.Status.AtProvider.ID = ""
			case "not Ready":
				target.SetConditions(xpv1.Unavailable())
			case "not Synced":
				target.SetConditions(xpv1.ReconcileError(errors.New("not synced")))
			case "stale generation":
				target.Generation = 2
			case "deleting":
				now := metav1.Now()
				target.DeletionTimestamp = &now
				target.Finalizers = []string{"test"}
			}
			objs := []client.Object{target}
			if failure == "missing" {
				objs = nil
			}
			d := &namespaced.TrustedDomain{ObjectMeta: metav1.ObjectMeta{Name: "browser", Namespace: "tenant"}}
			d.Spec.ForProvider.InstanceIDSelector = &xpv1.NamespacedSelector{MatchLabels: target.Labels}
			if err := d.ResolveReferences(context.Background(), instanceRefClient(t, objs...)); err == nil {
				t.Fatal("expected unresolved instance")
			}
			if d.Spec.ForProvider.InstanceID != nil {
				t.Fatal("an unresolved observation supplied an ID")
			}
		})
	}
}

func TestLiteralInstanceIDRemainsSupported(t *testing.T) {
	literal := "existing-instance"
	d := &namespaced.TrustedDomain{ObjectMeta: metav1.ObjectMeta{Name: "browser", Namespace: "tenant"}}
	d.Spec.ForProvider.InstanceID = &literal
	d.Spec.ForProvider.InstanceIDSelector = &xpv1.NamespacedSelector{MatchLabels: instanceTarget("tenant", "").Labels}
	if err := d.ResolveReferences(context.Background(), instanceRefClient(t)); err != nil {
		t.Fatal(err)
	}
	if *d.Spec.ForProvider.InstanceID != literal || d.Spec.ForProvider.InstanceIDRef != nil {
		t.Fatal("literal instance ID changed")
	}
}

func TestClusterInstanceReferences(t *testing.T) {
	target := &cluster.Instance{ObjectMeta: metav1.ObjectMeta{Name: "installed", Labels: map[string]string{"example.org/instance": "installed"}}}
	target.Status.AtProvider = observation.Instance{ID: "cluster-instance"}
	target.SetConditions(xpv1.Available(), xpv1.ReconcileSuccess())
	for _, method := range []string{"reference", "selector"} {
		t.Run(method, func(t *testing.T) {
			d := &cluster.CustomDomain{}
			if method == "reference" {
				d.Spec.ForProvider.InstanceIDRef = &xpv1.Reference{Name: target.Name}
			} else {
				d.Spec.ForProvider.InstanceIDSelector = &xpv1.Selector{MatchLabels: target.Labels}
			}
			if err := d.ResolveReferences(context.Background(), instanceRefClient(t, target)); err != nil {
				t.Fatal(err)
			}
			if d.Spec.ForProvider.InstanceID == nil || *d.Spec.ForProvider.InstanceID != "cluster-instance" {
				t.Fatal("incorrect cluster-scoped ID")
			}
		})
	}
}
