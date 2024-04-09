package recordset

import "github.com/crossplane/upjet/pkg/config"

const (
	apiVersion = "v1alpha1"
	shortGroup = "recordset"
)

// Configure configures individual resources by adding custom ResourceConfigurators.
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("dns_a_record_set", func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "ARecordSet"
		r.Version = apiVersion
	})

	p.AddResourceConfigurator("dns_aaaa_record_set", func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "AAAARecordSet"
		r.Version = apiVersion
	})

	p.AddResourceConfigurator("dns_cname_record_set", func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "CNAMERecordSet"
		r.Version = apiVersion
	})

	p.AddResourceConfigurator("dns_mx_record_set", func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "MXRecordSet"
		r.Version = apiVersion
	})

	p.AddResourceConfigurator("dns_ns_record_set", func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "NSRecordSet"
		r.Version = apiVersion
	})

	p.AddResourceConfigurator("dns_srv_record_set", func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "SRVRecordSet"
		r.Version = apiVersion
	})

	p.AddResourceConfigurator("dns_txt_record_set", func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "TXTRecordSet"
		r.Version = apiVersion
	})

}
