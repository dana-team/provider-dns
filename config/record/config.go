package record

import "github.com/crossplane/upjet/pkg/config"

const (
	apiVersion = "v1alpha1"
	shortGroup = "record"
)

// Configure configures individual resources by adding custom ResourceConfigurators.
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("dns_cname_record", func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "CNAMERecord"
		r.Version = apiVersion
	})

	p.AddResourceConfigurator("dns_ptr_record", func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "PTRRecord"
		r.Version = apiVersion
	})
}
