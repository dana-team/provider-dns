package txt

import "github.com/crossplane/upjet/pkg/config"

// Configure configures individual resources by adding custom ResourceConfigurators.
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("dns_txt_record_set", func(r *config.Resource) {
		r.ShortGroup = "txt"
	})
}
