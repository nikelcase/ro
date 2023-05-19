package web

import (
	"fmt"
	"time"

	"github.com/ergo-services/ergo/etf"
	"github.com/ergo-services/ergo/gen"

	"resourceObserver/names"
)

func CreateWebObserver() gen.SupervisorBehavior {
	return &observerWebSup{}
}

type observerWebSup struct {
	gen.Supervisor
}

func (*observerWebSup) Init(args ...etf.Term) (gen.SupervisorSpec, error) {
	spec := gen.SupervisorSpec{
		Name: names.ObsWebSupName,
		Children: []gen.SupervisorChildSpec{
			gen.SupervisorChildSpec{
				Name:  "Portal",
				Child: &observer{},
				Args:  []etf.Term{"abc", 67890},
			},
			//			gen.SupervisorChildSpec{
			//				Name:  "Support",
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
	if err := process.MonitorEvent(names.SchedEvent); err != nil {
		return err
	}
	return nil
}

func (*observer) HandleInfo(process *gen.ServerProcess, message etf.Term) gen.ServerStatus {
	pname := process.Name()
	switch message.(type) {
	case names.SchedEventMsg:
		fmt.Printf("consumer %s got event: %s\r\n", pname, message)
		msg := message.(names.SchedEventMsg)
		err := process.Cast(msg.Pid, names.UpdateEventMsg{
			Name:      pname,
			State:     names.ObsSuccess,
			Time:      time.Now(),
			Text:      pname,
			Title:     pname,
			BaseName:  pname,
			FullAddr:  pname,
			ShortAddr: pname,
		})
		fmt.Printf("consumer %s error: %#v\r\n", pname, err)
	case gen.MessageEventDown:
		fmt.Printf("%s producer has terminated\r\n", pname)
		return gen.ServerStatusStop
	default:
		fmt.Println("unknown message", message)
	}
	return gen.ServerStatusOK
}

func (*observer) Terminate(process *gen.ServerProcess, reason string) {
	fmt.Printf("[%s] Terminating process with reason %q\r\n", process.Name(), reason)
}
