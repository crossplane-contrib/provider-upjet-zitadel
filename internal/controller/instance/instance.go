// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

// Package instance reconciles read-only, provider-native Instance resources.
package instance

import (
	"context"
	"fmt"
	"time"

	xpv1 "github.com/crossplane/crossplane-runtime/v2/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/v2/pkg/event"
	"github.com/crossplane/crossplane-runtime/v2/pkg/meta"
	"github.com/crossplane/crossplane-runtime/v2/pkg/ratelimiter"
	"github.com/crossplane/crossplane-runtime/v2/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/v2/pkg/resource"
	tjcontroller "github.com/crossplane/upjet/v2/pkg/controller"
	"k8s.io/apimachinery/pkg/runtime/schema"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	cluster "github.com/crossplane-contrib/provider-upjet-zitadel/apis/cluster/instance/v1alpha1"
	namespaced "github.com/crossplane-contrib/provider-upjet-zitadel/apis/namespaced/instance/v1alpha1"
	"github.com/crossplane-contrib/provider-upjet-zitadel/apis/observation"
	"github.com/crossplane-contrib/provider-upjet-zitadel/internal/clients"
)

// SetupGated starts each scope after its CRD is available.
func SetupGated(mgr ctrl.Manager, o tjcontroller.Options) error {
	for _, gvk := range []schema.GroupVersionKind{cluster.InstanceGroupVersionKind, namespaced.InstanceGroupVersionKind} {
		o.Gate.Register(func() {
			if err := setup(mgr, o, gvk); err != nil {
				mgr.GetLogger().Error(err, "cannot start Instance observer", "gvk", gvk)
			}
		}, gvk)
	}
	return nil
}

// Setup starts both scopes without SafeStart gating.
func Setup(mgr ctrl.Manager, o tjcontroller.Options) error {
	for _, gvk := range []schema.GroupVersionKind{cluster.InstanceGroupVersionKind, namespaced.InstanceGroupVersionKind} {
		if err := setup(mgr, o, gvk); err != nil {
			return err
		}
	}
	return nil
}

func setup(mgr ctrl.Manager, o tjcontroller.Options, gvk schema.GroupVersionKind) error {
	name := managed.ControllerName(gvk.String())
	opts := []managed.ReconcilerOption{
		managed.WithExternalConnecter(&connector{kube: mgr.GetClient(), newClient: clients.NewInstanceClient}),
		managed.WithLogger(o.Logger.WithValues("controller", name)),
		managed.WithRecorder(event.NewAPIRecorder(mgr.GetEventRecorderFor(name))),
		managed.WithInitializers(), // Kubernetes name is not a Zitadel ID.
		managed.WithManagementPolicies(),
		managed.WithPollInterval(o.PollInterval),
		managed.WithTimeout(30 * time.Second),
	}
	if o.MetricOptions != nil {
		opts = append(opts, managed.WithMetricRecorder(o.MetricOptions.MRMetrics))
	}
	if o.PollJitter != 0 {
		opts = append(opts, managed.WithPollJitterHook(o.PollJitter))
	}
	reconciler := managed.NewReconciler(mgr, resource.ManagedKind(gvk), opts...)
	var object client.Object = &namespaced.Instance{}
	if gvk == cluster.InstanceGroupVersionKind {
		object = &cluster.Instance{}
	}
	return ctrl.NewControllerManagedBy(mgr).Named(name).WithOptions(o.ForControllerRuntime()).For(object).
		WithEventFilter(resource.DesiredStateChanged()).Complete(ratelimiter.NewReconciler(name, reconciler, o.GlobalRateLimiter))
}

type connector struct {
	kube      client.Client
	newClient func(context.Context, client.Client, resource.Managed) (clients.InstanceClient, error)
}

func (c *connector) Connect(ctx context.Context, mg resource.Managed) (managed.ExternalClient, error) {
	if _, ok := mg.(observation.ObservedInstance); !ok {
		return nil, fmt.Errorf("managed resource is not an Instance")
	}
	for _, policy := range mg.GetManagementPolicies() {
		if policy != xpv1.ManagementActionObserve {
			return nil, fmt.Errorf("instance permits only Observe management policy")
		}
	}
	api, err := c.newClient(ctx, c.kube, mg)
	if err != nil {
		return nil, err
	}
	return &external{kube: c.kube, api: api}, nil
}

type external struct {
	kube client.Client
	api  clients.InstanceClient
}

func (e *external) Observe(ctx context.Context, mg resource.Managed) (managed.ExternalObservation, error) {
	inst, ok := mg.(observation.ObservedInstance)
	if !ok {
		return managed.ExternalObservation{}, fmt.Errorf("managed resource is not an Instance")
	}
	observed, err := e.api.Observe(ctx)
	if err != nil {
		inst.SetConditions(xpv1.Unavailable())
		return managed.ExternalObservation{}, err
	}
	if observed.ID == "" {
		inst.SetConditions(xpv1.Unavailable())
		return managed.ExternalObservation{}, fmt.Errorf("empty instance ID")
	}
	// An endpoint/credential change or restored database must not silently
	// retarget dependent domains to a different instance. Recreate the observer
	// explicitly when accepting a different external identity.
	if bound := meta.GetExternalName(inst); bound != "" && bound != observed.ID {
		inst.SetConditions(xpv1.Unavailable())
		return managed.ExternalObservation{}, fmt.Errorf("observed instance ID differs from the bound external identity; explicit adoption is required")
	}
	if meta.GetExternalName(inst) == "" {
		meta.SetExternalName(inst, observed.ID)
		if err := e.kube.Update(ctx, inst); err != nil {
			return managed.ExternalObservation{}, fmt.Errorf("cannot persist observed instance identity: %w", err)
		}
	}
	inst.SetInstanceObservation(observed)
	inst.SetConditions(xpv1.Available())
	return managed.ExternalObservation{ResourceExists: true, ResourceUpToDate: true}, nil
}

func (e *external) Create(context.Context, resource.Managed) (managed.ExternalCreation, error) {
	return managed.ExternalCreation{}, fmt.Errorf("instance observation cannot create a Zitadel instance")
}
func (e *external) Update(context.Context, resource.Managed) (managed.ExternalUpdate, error) {
	return managed.ExternalUpdate{}, fmt.Errorf("instance observation cannot update a Zitadel instance")
}

// Delete only forgets the observation. There is deliberately no remote API call.
func (e *external) Delete(context.Context, resource.Managed) (managed.ExternalDelete, error) {
	return managed.ExternalDelete{}, nil
}
func (e *external) Disconnect(context.Context) error { return e.api.Close() }
