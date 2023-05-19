package observers

import (
	"github.com/ergo-services/ergo/etf"
	"github.com/ergo-services/ergo/gen"

	"resourceObserver/names"
	"resourceObserver/observers/scheduler"
)

func CreateObserverSup() gen.SupervisorBehavior {
	return &observerSup{}
}

type observerSup struct {
	gen.Supervisor
}

func (*observerSup) Init(args ...etf.Term) (gen.SupervisorSpec, error) {
	spec := gen.SupervisorSpec{
		Name: names.ObsSupName,
		Children: []gen.SupervisorChildSpec{
			gen.SupervisorChildSpec{
				Name:  names.ObsSchedName,
				Child: scheduler.CreateScheduler(),
			},
			gen.SupervisorChildSpec{
				Name:  names.ObsWrkSupName,
				Child: createObserverWorksersSup(),
			},
		},
		Strategy: gen.SupervisorStrategy{
			Type:      gen.SupervisorStrategyOneForAll,
			Intensity: 5,
			Period:    5,
			Restart:   gen.SupervisorStrategyRestartTransient,
		},
	}
	return spec, nil
}
