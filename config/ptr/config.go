package ptr

import "github.com/crossplane/upjet/pkg/config"

// Configure configures individual resources by adding custom ResourceConfigurators.
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("dns_ptr_record", func(r *config.Resource) {
		r.ShortGroup = "ptr"
	})
}
