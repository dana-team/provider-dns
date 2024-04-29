/*
Copyright 2022 Upbound Inc.
*/

package config

import "github.com/crossplane/upjet/pkg/config"

// terraformPluginSDKExternalNameConfigs contains all external name configurations for this
// provider.
var terraformPluginSDKExternalNameConfigs = map[string]config.ExternalName{
	"dns_a_record_set":    config.TemplatedStringAsIdentifier("", "{{ .external_name }}.{{ .parameters.zone }}"),
	"dns_aaaa_record_set": config.TemplatedStringAsIdentifier("", "{{ .external_name }}.{{ .parameters.zone }}"),
}

var terraformPluginFrameworkExternalNameConfigs = map[string]config.ExternalName{
	"dns_mx_record_set":  config.TemplatedStringAsIdentifier("", "{{ .external_name }}.{{ .parameters.zone }}"),
	"dns_ns_record_set":  config.TemplatedStringAsIdentifier("", "{{ .external_name }}.{{ .parameters.zone }}"),
	"dns_srv_record_set": config.TemplatedStringAsIdentifier("", "{{ .external_name }}.{{ .parameters.zone }}"),
	"dns_txt_record_set": config.TemplatedStringAsIdentifier("", "{{ .external_name }}.{{ .parameters.zone }}"),
	"dns_cname_record":   config.TemplatedStringAsIdentifier("", "{{ .external_name }}.{{ .parameters.zone }}"),
	"dns_ptr_record":     config.TemplatedStringAsIdentifier("", "{{ .external_name }}.{{ .parameters.zone }}"),
}

// cliReconciledExternalNameConfigs contains all external name configurations
// belonging to Terraform resources to be reconciled under the CLI-based
// architecture for this provider.
var cliReconciledExternalNameConfigs = map[string]config.ExternalName{}

// resourceConfigurator applies all external name configs listed in
// the table terraformPluginSDKExternalNameConfigs,
// cliReconciledExternalNameConfigs, and
// terraformPluginFrameworkExternalNameConfigs and sets the version of
// those resources to v1beta1.
func resourceConfigurator() config.ResourceOption {
	return func(r *config.Resource) {
		// If an external name is configured for multiple architectures,
		// Terraform Plugin Framework takes precedence over Terraform
		// Plugin SDKv2, which takes precedence over CLI architecture.
		e, configured := terraformPluginFrameworkExternalNameConfigs[r.Name]
		if !configured {
			e, configured = terraformPluginSDKExternalNameConfigs[r.Name]
			if !configured {
				e, configured = cliReconciledExternalNameConfigs[r.Name]
			}
		}
		if !configured {
			return
		}
		r.Version = "v1beta1"
		r.ExternalName = e
	}
}
