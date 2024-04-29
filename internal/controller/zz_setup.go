/*
Copyright 2022 Upbound Inc.
*/

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/pkg/controller"

	providerconfig "github.com/dana-team/provider-dns/internal/controller/providerconfig"
	cnamerecord "github.com/dana-team/provider-dns/internal/controller/record/cnamerecord"
	ptrrecord "github.com/dana-team/provider-dns/internal/controller/record/ptrrecord"
	aaaarecordset "github.com/dana-team/provider-dns/internal/controller/recordset/aaaarecordset"
	arecordset "github.com/dana-team/provider-dns/internal/controller/recordset/arecordset"
	mxrecordset "github.com/dana-team/provider-dns/internal/controller/recordset/mxrecordset"
	nsrecordset "github.com/dana-team/provider-dns/internal/controller/recordset/nsrecordset"
	srvrecordset "github.com/dana-team/provider-dns/internal/controller/recordset/srvrecordset"
	txtrecordset "github.com/dana-team/provider-dns/internal/controller/recordset/txtrecordset"
)

// Setup creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		providerconfig.Setup,
		cnamerecord.Setup,
		ptrrecord.Setup,
		aaaarecordset.Setup,
		arecordset.Setup,
		mxrecordset.Setup,
		nsrecordset.Setup,
		srvrecordset.Setup,
		txtrecordset.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
