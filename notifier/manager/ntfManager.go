package manager

import (
	"fmt"

	"github.com/ergo-services/ergo/etf"
	"github.com/ergo-services/ergo/gen"
	"github.com/ergo-services/ergo/lib"

	"resourceObserver/names"
)

func CreateManager() gen.ServerBehavior {
	return &managerSrv{}
}

type managerSrv struct {
	gen.Server
}

func (*managerSrv) Init(process *gen.ServerProcess, args ...etf.Term) error {
	if err := process.RegisterEvent(names.NotifyMgrEvent, names.NotifyMgrEventMsg{}); err != nil {
		lib.Warning("can't register event %q: %s", names.NotifyMgrEvent, err)
	}
	fmt.Printf("process %s registered event %s\n", process.Self(), names.NotifyMgrEvent)
	if err := process.MonitorEvent(names.UpdateEvent); err != nil {
		lib.Warning("can't monitor event %q: %s", names.UpdateEvent, err)
		return nil
	}
	if err := process.MonitorEvent(names.DownloadEvent); err != nil {
		lib.Warning("can't monitor event %q: %s", names.DownloadEvent, err)
		return nil
	}
	return nil
}

func (*managerSrv) HandleInfo(process *gen.ServerProcess, message etf.Term) gen.ServerStatus {
	switch message.(type) {
	case names.UpdateEventMsg, names.DownloadEventMsg:
		fmt.Printf("consumer %s got event: %s\r\n", process.Name(), message)
		event := names.NotifyMgrEventMsg{Pid: process.Self(), Msg: message}
		fmt.Printf("... producing event: %#v\n", event)
		if err := process.SendEventMessage(names.NotifyMgrEvent, event); err != nil {
			fmt.Println("can't send event:", err)
		}
	case gen.MessageEventDown:
		fmt.Printf("%s producer has terminated\r\n", process.Name())
		return gen.ServerStatusStop
	default:
		fmt.Println("unknown message", message)
	}
	return gen.ServerStatusOK
}

func (*managerSrv) HandleCast(process *gen.ServerProcess, message etf.Term) gen.ServerStatus {
	fmt.Printf("[%s] HandleCast: %#v\n", process.Name(), message)
	msg, ok := message.(names.NotifyMsg)
	if ok == true {
		fmt.Printf("... producing event: %#v\n", msg)
		//if msg.state == error -> try again
	} else {
		fmt.Println("unknown message", message)
	}
	return gen.ServerStatusOK
}

func (*managerSrv) Terminate(process *gen.ServerProcess, reason string) {
	process.UnregisterEvent(names.NotifyMgrEvent)
	process.DemonitorEvent(names.UpdateEvent)
	process.DemonitorEvent(names.DownloadEvent)
	fmt.Printf("[%s] Terminating process with reason %q\r\n", process.Name(), reason)
}
