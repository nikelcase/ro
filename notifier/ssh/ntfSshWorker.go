package ssh

import (
	"fmt"

	"github.com/ergo-services/ergo/etf"
	"github.com/ergo-services/ergo/gen"

	"resourceObserver/names"
)

func CreateSshNotifier() gen.SupervisorBehavior {
	return &notifierFtpSup{}
}

type notifierFtpSup struct {
	gen.Supervisor
}

func (*notifierFtpSup) Init(args ...etf.Term) (gen.SupervisorSpec, error) {
	spec := gen.SupervisorSpec{
		Name: names.NotifySshSupName,
		Children: []gen.SupervisorChildSpec{
			gen.SupervisorChildSpec{
				Name:  "ntfSsh1",
				Child: &observer{},
				Args:  []etf.Term{"abc", 67890},
			},
			//			gen.SupervisorChildSpec{
			//				Name:  "ntfSsh2",
			//				Child: &observer{},
			//				Args:  []etf.Term{"abc", 67890},
			//			},
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

type observer struct {
	gen.Server
}

func (*observer) Init(process *gen.ServerProcess, args ...etf.Term) error {
	fmt.Printf("Started new process\n\tPid: %s\n\tName: %q\n\tParent: %s\n\tArgs:%#v\n",
		process.Self(),
		process.Name(),
		process.Parent().Self(),
		args)
	if err := process.MonitorEvent(names.NotifyMgrEvent); err != nil {
		return err
	}
	return nil
}

func (*observer) HandleInfo(process *gen.ServerProcess, message etf.Term) gen.ServerStatus {
	switch message.(type) {
	case names.NotifyMgrEventMsg:
		fmt.Printf("consumer %s got event: %s\r\n", process.Name(), message)
	case gen.MessageEventDown:
		fmt.Printf("%s producer has terminated\r\n", process.Name())
		return gen.ServerStatusStop
	default:
		fmt.Println("unknown message", message)
	}
	return gen.ServerStatusOK
}

func (*observer) Terminate(process *gen.ServerProcess, reason string) {
	fmt.Printf("[%s] Terminating process with reason %q\r\n", process.Name(), reason)
}
