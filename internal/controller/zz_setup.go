/*
Copyright 2022 Upbound Inc.
*/

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/pkg/controller"

	recordset "github.com/dana-team/provider-dns/internal/controller/a/recordset"
	recordsetaaaa "github.com/dana-team/provider-dns/internal/controller/aaaa/recordset"
	record "github.com/dana-team/provider-dns/internal/controller/cname/record"
	recordsetmx "github.com/dana-team/provider-dns/internal/controller/mx/recordset"
	recordsetns "github.com/dana-team/provider-dns/internal/controller/ns/recordset"
	providerconfig "github.com/dana-team/provider-dns/internal/controller/providerconfig"
	recordptr "github.com/dana-team/provider-dns/internal/controller/ptr/record"
	recordsetsrv "github.com/dana-team/provider-dns/internal/controller/srv/recordset"
	recordsettxt "github.com/dana-team/provider-dns/internal/controller/txt/recordset"
)

// Setup creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		recordset.Setup,
		recordsetaaaa.Setup,
		record.Setup,
		recordsetmx.Setup,
		recordsetns.Setup,
		providerconfig.Setup,
		recordptr.Setup,
		recordsetsrv.Setup,
		recordsettxt.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
