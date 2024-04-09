package recordset

import "github.com/crossplane/upjet/pkg/config"

// Configure configures individual resources by adding custom ResourceConfigurators.
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("dns_a_record_set", func(r *config.Resource) {
		r.ShortGroup = "recordset"
		r.Kind = "ARecordSet"
		r.Version = "v1alpha1"
	})

	p.AddResourceConfigurator("dns_aaaa_record_set", func(r *config.Resource) {
		r.ShortGroup = "recordset"
		r.Kind = "AAAARecordSet"
		r.Version = "v1alpha1"
	})

	p.AddResourceConfigurator("dns_cname_record_set", func(r *config.Resource) {
		r.ShortGroup = "recordset"
		r.Kind = "CNAMERecordSet"
		r.Version = "v1alpha1"
	})

	p.AddResourceConfigurator("dns_mx_record_set", func(r *config.Resource) {
		r.ShortGroup = "recordset"
		r.Kind = "MXRecordSet"
		r.Version = "v1alpha1"
	})

	p.AddResourceConfigurator("dns_ns_record_set", func(r *config.Resource) {
		r.ShortGroup = "recordset"
		r.Kind = "NSRecordSet"
		r.Version = "v1alpha1"
	})

	p.AddResourceConfigurator("dns_srv_record_set", func(r *config.Resource) {
		r.ShortGroup = "recordset"
		r.Kind = "SRVRecordSet"
		r.Version = "v1alpha1"
	})

	p.AddResourceConfigurator("dns_txt_record_set", func(r *config.Resource) {
		r.ShortGroup = "recordset"
		r.Kind = "TXTRecordSet"
		r.Version = "v1alpha1"
	})

}
