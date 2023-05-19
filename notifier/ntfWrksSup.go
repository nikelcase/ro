package notifier

import (
	"github.com/ergo-services/ergo/etf"
	"github.com/ergo-services/ergo/gen"

	"resourceObserver/names"
	"resourceObserver/notifier/ssh"
	"resourceObserver/notifier/stdout"
	"resourceObserver/notifier/telnet"
	"resourceObserver/notifier/web"
	"resourceObserver/notifier/winsent"
)

func createNotifyWorksersSup() gen.SupervisorBehavior {
	return &observerWorkersSup{}
}

type observerWorkersSup struct {
	gen.Supervisor
}

func (*observerWorkersSup) Init(args ...etf.Term) (gen.SupervisorSpec, error) {
	spec := gen.SupervisorSpec{
		Name: names.NotifyWrkSupName,
		Children: []gen.SupervisorChildSpec{
			gen.SupervisorChildSpec{
				Name:  names.NotifyWebSupName,
				Child: web.CreateWebNotifier(),
			},
			gen.SupervisorChildSpec{
				Name:  names.NotifySshSupName,
				Child: ssh.CreateSshNotifier(),
			},
			gen.SupervisorChildSpec{
				Name:  names.NotifyWinsentSupName,
				Child: winsent.CreateWinsentNotifier(),
			},
			gen.SupervisorChildSpec{
				Name:  names.NotifyTelnetSupName,
				Child: telnet.CreateTelnetNotifier(),
			},
			gen.SupervisorChildSpec{
				Name:  names.NotifyStdoutSupName,
				Child: stdout.CreateStdoutNotifier(),
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
