/*
Copyright 2021 Upbound Inc.
*/

package config

import (
	// Note(turkenh): we are importing this to embed provider schema document
	_ "embed"
	"github.com/dana-team/provider-dns/config/record"
	"github.com/dana-team/provider-dns/config/recordset"

	ujconfig "github.com/crossplane/upjet/pkg/config"
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
		record.Configure,
		recordset.Configure,
	} {
		configure(pc)
	}

	pc.ConfigureResources()
	return pc
}
