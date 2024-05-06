/*
Copyright 2021 Upbound Inc.
*/

package config

import (
	"context"

	_ "embed" // needed for upjet

	ujconfig "github.com/crossplane/upjet/pkg/config"
	"github.com/crossplane/upjet/pkg/registry/reference"
	"github.com/dana-team/provider-dns/config/record"
	"github.com/dana-team/provider-dns/config/recordset"
	"github.com/hashicorp/terraform-provider-dns/xpprovider"
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
func GetProvider(ctx context.Context) (*ujconfig.Provider, error) {
	fwProvider, p := xpprovider.GetProvider(ctx)

	pc := ujconfig.NewProvider([]byte(providerSchema), resourcePrefix, modulePath, []byte(providerMetadata),
		ujconfig.WithRootGroup("dns.crossplane.io"),
		ujconfig.WithShortName("dns"),
		ujconfig.WithIncludeList(resourceList(cliReconciledExternalNameConfigs)),
		ujconfig.WithTerraformPluginSDKIncludeList(resourceList(terraformPluginSDKExternalNameConfigs)),
		ujconfig.WithTerraformPluginFrameworkIncludeList(resourceList(terraformPluginFrameworkExternalNameConfigs)),
		ujconfig.WithDefaultResourceOptions(
			resourceConfigurator(),
		),
		ujconfig.WithReferenceInjectors([]ujconfig.ReferenceInjector{reference.NewInjector(modulePath)}),
		ujconfig.WithFeaturesPackage("internal/features"),
		ujconfig.WithTerraformProvider(p),
		ujconfig.WithTerraformPluginFrameworkProvider(fwProvider),
	)

	for _, configure := range []func(provider *ujconfig.Provider){
		// add custom config functions
		record.Configure,
		recordset.Configure,
	} {
		configure(pc)
	}

	pc.ConfigureResources()
	return pc, nil
}

// resourceList returns the list of resources that have external
// name configured in the specified table.
func resourceList(t map[string]ujconfig.ExternalName) []string {
	l := make([]string, len(t))
	i := 0
	for n := range t {
		// Expected format is regex and we'd like to have exact matches.
		l[i] = n + "$"
		i++
	}
	return l
}
