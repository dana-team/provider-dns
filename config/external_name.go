/*
Copyright 2022 Upbound Inc.
*/

package config

import "github.com/crossplane/upjet/pkg/config"

// ExternalNameConfigs contains all external name configurations for this
// provider.
var ExternalNameConfigs = map[string]config.ExternalName{
	"dns_a_record_set":    config.TemplatedStringAsIdentifier("", "{{ .external_name }}.{{ .parameters.zone }}"),
	"dns_aaaa_record_set": config.TemplatedStringAsIdentifier("", "{{ .external_name }}.{{ .parameters.zone }}"),
	"dns_cname_record":    config.TemplatedStringAsIdentifier("", "{{ .external_name }}.{{ .parameters.zone }}"),
	"dns_mx_record_set":   config.TemplatedStringAsIdentifier("", "{{ .external_name }}.{{ .parameters.zone }}"),
	"dns_ns_record_set":   config.TemplatedStringAsIdentifier("", "{{ .external_name }}.{{ .parameters.zone }}"),
	"dns_ptr_record":      config.TemplatedStringAsIdentifier("", "{{ .external_name }}.{{ .parameters.zone }}"),
	"dns_srv_record_set":  config.TemplatedStringAsIdentifier("", "{{ .external_name }}.{{ .parameters.zone }}"),
	"dns_txt_record_set":  config.TemplatedStringAsIdentifier("", "{{ .external_name }}.{{ .parameters.zone }}"),
}

// ExternalNameConfigurations applies all external name configs listed in the
// table ExternalNameConfigs and sets the version of those resources to v1beta1
// assuming they will be tested.
func ExternalNameConfigurations() config.ResourceOption {
	return func(r *config.Resource) {
		if e, ok := ExternalNameConfigs[r.Name]; ok {
			r.ExternalName = e
		}
	}
}

// ExternalNameConfigured returns the list of all resources whose external name
// is configured manually.
func ExternalNameConfigured() []string {
	l := make([]string, len(ExternalNameConfigs))
	i := 0
	for name := range ExternalNameConfigs {
		// $ is added to match the exact string since the format is regex.
		l[i] = name + "$"
		i++
	}
	return l
}
