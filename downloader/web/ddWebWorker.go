package web

import (
	"fmt"

	"github.com/ergo-services/ergo/etf"
	"github.com/ergo-services/ergo/gen"

	"resourceObserver/names"
)

func CreateWebDownloader() gen.SupervisorBehavior {
	return &downloaderWebSup{}
}

type downloaderWebSup struct {
	gen.Supervisor
}

func (*downloaderWebSup) Init(args ...etf.Term) (gen.SupervisorSpec, error) {
	spec := gen.SupervisorSpec{
		Name: names.DownloadWebSupName,
		Children: []gen.SupervisorChildSpec{
			gen.SupervisorChildSpec{
				Name:  "ddPortal",
				Child: &observer{},
				Args:  []etf.Term{"abc", 67890},
			},
			//			gen.SupervisorChildSpec{
			//				Name:  "ddSupport",
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
	if err := process.MonitorEvent(names.DownloadMgrEvent); err != nil {
		return err
	}
	return nil
}

func (*observer) HandleInfo(process *gen.ServerProcess, message etf.Term) gen.ServerStatus {
	pname := process.Name()
	switch message.(type) {
	case names.DownloadMgrEventMsg:
		fmt.Printf("consumer %s got event: %s\r\n", pname, message)
		msg := message.(names.DownloadMgrEventMsg)
		err := process.Cast(msg.Pid, names.DownloadEventMsg{
			State:  names.DdSuccess,
			Unpack: true,
			Title:  pname,
			Path:   pname,
		})
		fmt.Printf("consumer %s error: %#v\r\n", pname, err)
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
