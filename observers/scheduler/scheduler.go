package scheduler

import (
	"fmt"
	"time"

	"github.com/ergo-services/ergo/etf"
	"github.com/ergo-services/ergo/gen"
	"github.com/ergo-services/ergo/lib"

	"resourceObserver/names"
)

func CreateScheduler() gen.ServerBehavior {
	return &schedulerSrv{}
}

type schedulerSrv struct {
	gen.Server
}

func (*schedulerSrv) Init(process *gen.ServerProcess, args ...etf.Term) error {
	if err := process.RegisterEvent(names.SchedEvent, names.SchedEventMsg{}); err != nil {
		lib.Warning("can't register event %q: %s", names.SchedEvent, err)
	}
	fmt.Printf("process %s registered event %s\n", process.Self(), names.SchedEvent)
	if err := process.RegisterEvent(names.UpdateEvent, names.UpdateEventMsg{}); err != nil {
		lib.Warning("can't register event %q: %s", names.UpdateEvent, err)
	}
	fmt.Printf("process %s registered event %s\n", process.Self(), names.UpdateEvent)
	process.SendAfter(process.Self(), 1, time.Second)
	return nil
}

func (*schedulerSrv) HandleInfo(process *gen.ServerProcess, message etf.Term) gen.ServerStatus {
	n := message.(int)
	if n > 2 {
		return gen.ServerStatusStop
	}
	// sending message with delay 1 second
	process.SendAfter(process.Self(), n+1, time.Second)
	event := names.SchedEventMsg{Pid: process.Self()}

	fmt.Printf("%s... producing event: %s\n", process.Name(), event)
	if err := process.SendEventMessage(names.SchedEvent, event); err != nil {
		fmt.Println("can't send event:", err)
	}
	return gen.ServerStatusOK
}

func (*schedulerSrv) HandleCast(process *gen.ServerProcess, message etf.Term) gen.ServerStatus {
	fmt.Printf("[%s] HandleCast: %#v\n", process.Name(), message)
	msg, ok := message.(names.UpdateEventMsg)
	if ok == true {
		fmt.Printf("%s... producing event: %#v\n", process.Name(), msg)
		if err := process.SendEventMessage(names.UpdateEvent, msg); err != nil {
			fmt.Println("can't send event:", err)
		}
	} else {
		fmt.Println("unknown message", message)
	}
	return gen.ServerStatusOK
}

func (*schedulerSrv) Terminate(process *gen.ServerProcess, reason string) {
	process.UnregisterEvent(names.SchedEvent)
	process.UnregisterEvent(names.UpdateEvent)
	fmt.Printf("[%s] Terminating process with reason %q\r\n", process.Name(), reason)
}
