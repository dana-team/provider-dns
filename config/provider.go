/*
Copyright 2021 Upbound Inc.
*/

package config

import (
	// Note(turkenh): we are importing this to embed provider schema document
	_ "embed"

	ujconfig "github.com/crossplane/upjet/pkg/config"
	"github.com/dana-team/provider-dns/config/a"
	"github.com/dana-team/provider-dns/config/aaaa"
	"github.com/dana-team/provider-dns/config/cname"
	"github.com/dana-team/provider-dns/config/mx"
	"github.com/dana-team/provider-dns/config/ns"
	"github.com/dana-team/provider-dns/config/ptr"
	"github.com/dana-team/provider-dns/config/srv"
	"github.com/dana-team/provider-dns/config/txt"
)

const (
	resourcePrefix = "dns"
	modulePath     = "github.com/dana-team/provider-dns"
)

//go:embed schema.json
var providerSchema string

//go:embed provider-metadata.yaml
var providerMetadata string

// GetProvider returns provider configuration
func GetProvider() *ujconfig.Provider {
	pc := ujconfig.NewProvider([]byte(providerSchema), resourcePrefix, modulePath, []byte(providerMetadata),
		ujconfig.WithRootGroup("dns.crossplane.io"),
		ujconfig.WithIncludeList(ExternalNameConfigured()),
		ujconfig.WithFeaturesPackage("internal/features"),
		ujconfig.WithDefaultResourceOptions(
			ExternalNameConfigurations(),
		))

	for _, configure := range []func(provider *ujconfig.Provider){
		// add custom config functions
		a.Configure,
		aaaa.Configure,
		cname.Configure,
		mx.Configure,
		ns.Configure,
		ptr.Configure,
		srv.Configure,
		txt.Configure,
	} {
		configure(pc)
	}

	pc.ConfigureResources()
	return pc
}
