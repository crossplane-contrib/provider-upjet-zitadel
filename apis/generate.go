//go:build generate
// +build generate

/*
Copyright 2024 The Crossplane Authors.
*/

// NOTE: See the below link for details on what is happening here.
// https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module

// Remove existing CRDs
//go:generate rm -rf ../package/crds

// Remove generated files
//go:generate bash -c "find . -iname 'zz_*' ! -iname 'zz_generated.managed*.go' -delete"
//go:generate bash -c "find . -type d -empty -delete"
//go:generate bash -c "find ../internal/controller/cluster ../internal/controller/namespaced -iname 'zz_*' -delete"
//go:generate bash -c "find ../internal/controller/cluster ../internal/controller/namespaced -type d -empty -delete"
//go:generate rm -rf ../examples-generated

// Generate documentation from Terraform docs.
//go:generate go run github.com/crossplane/upjet/v2/cmd/scraper -n ${TERRAFORM_PROVIDER_SOURCE} -r ../.work/${TERRAFORM_PROVIDER_SOURCE}/${TERRAFORM_DOCS_PATH} -o ../config/provider-metadata.yaml --prelude-xpath "//text()[contains(., \"page_title\")]"

// Run Upjet generator
//go:generate go run ../cmd/generator/main.go ..

// Generate deepcopy methodsets and CRD manifests
//go:generate go run -tags generate sigs.k8s.io/controller-tools/cmd/controller-gen object:headerFile=../hack/boilerplate.go.txt paths=./... crd:allowDangerousTypes=true,crdVersions=v1 output:artifacts:config=../package/crds

// Generate crossplane-runtime methodsets (resource.Claim, etc)
//go:generate go run -tags generate github.com/crossplane/crossplane-tools/cmd/angryjet generate-methodsets --header-file=../hack/boilerplate.go.txt ./...

// Workaround: fix angryjet bug where *float64 single-resolution fields are emitted
// using ptr.FromFloatPtrValue/ptr.ToFloatPtrValue (k8s.io/utils/ptr) instead of
// reference.FromFloatPtrValue/reference.ToFloatPtrValue (crossplane-runtime/pkg/reference).
//go:generate bash ../hack/fix-angryjet-float.sh

// Rewrite the generated resolvers to look up reference targets through
// internal/apis instead of importing their API packages. Cross-group
// references (e.g. org.Member -> user.HumanUser -> org.Organization) would
// otherwise create import cycles. See https://github.com/crossplane/upjet/issues/96
//go:generate go run github.com/crossplane/upjet/v2/cmd/resolver -g zitadel.crossplane.io -a github.com/crossplane-contrib/provider-upjet-zitadel/internal/apis -s -p ./cluster/...
//go:generate go run github.com/crossplane/upjet/v2/cmd/resolver -g zitadel.m.crossplane.io -a github.com/crossplane-contrib/provider-upjet-zitadel/internal/apis -s -p ./namespaced/...

package apis

import (
	_ "sigs.k8s.io/controller-tools/cmd/controller-gen" //nolint:typecheck

	_ "github.com/crossplane/crossplane-tools/cmd/angryjet" //nolint:typecheck
)
