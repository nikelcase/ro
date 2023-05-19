package notifier

import (
	"github.com/ergo-services/ergo/etf"
	"github.com/ergo-services/ergo/gen"

	"resourceObserver/names"
	"resourceObserver/notifier/manager"
)

func CreateNotifierSup() gen.SupervisorBehavior {
	return &notifierSup{}
}

type notifierSup struct {
	gen.Supervisor
}

func (*notifierSup) Init(args ...etf.Term) (gen.SupervisorSpec, error) {
	spec := gen.SupervisorSpec{
		Name: names.NotifySupName,
		Children: []gen.SupervisorChildSpec{
			gen.SupervisorChildSpec{
				Name:  names.NotifyMgrName,
				Child: manager.CreateManager(),
			},
			gen.SupervisorChildSpec{
				Name:  names.NotifyWrkSupName,
				Child: createNotifyWorksersSup(),
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
