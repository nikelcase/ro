package observers

import (
	"github.com/ergo-services/ergo/etf"
	"github.com/ergo-services/ergo/gen"

	"resourceObserver/names"
	"resourceObserver/observers/ftp"
	"resourceObserver/observers/smb"
	"resourceObserver/observers/web"
)

func createObserverWorksersSup() gen.SupervisorBehavior {
	return &observerWorkersSup{}
}

type observerWorkersSup struct {
	gen.Supervisor
}

func (*observerWorkersSup) Init(args ...etf.Term) (gen.SupervisorSpec, error) {
	spec := gen.SupervisorSpec{
		Name: names.ObsWrkSupName,
		Children: []gen.SupervisorChildSpec{
			gen.SupervisorChildSpec{
				Name:  names.ObsWebSupName,
				Child: web.CreateWebObserver(),
			},
			gen.SupervisorChildSpec{
				Name:  names.ObsSmbSupName,
				Child: smb.CreateSmbObserver(),
			},
			gen.SupervisorChildSpec{
				Name:  names.ObsFtpSupName,
				Child: ftp.CreateFtpObserver(),
			},
		},
		Strategy: gen.SupervisorStrategy{
			Type:      gen.SupervisorStrategyOneForOne,
			Intensity: 5,
			Period:    5,
			Restart:   gen.SupervisorStrategyRestartTransient,
		},
	}
	return spec, nil
}
