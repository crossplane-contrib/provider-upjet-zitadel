/*
Copyright 2024 The Crossplane Authors.
*/

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	"github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/providerconfig"
)

// Setup creates all cluster-scoped controllers with the supplied logger and adds them to
// the supplied manager.
func Setup(mgr ctrl.Manager, o controller.Options) error {
	return providerconfig.Setup(mgr, o)
}

// SetupGated creates all cluster-scoped controllers with the supplied logger and adds them to
// the supplied manager, using gated setup.
func SetupGated(mgr ctrl.Manager, o controller.Options) error {
	return providerconfig.SetupGated(mgr, o)
}
