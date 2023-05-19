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
	if err := process.RegisterEvent(names.DownloadMgrEvent, names.DownloadMgrEventMsg{}); err != nil {
		lib.Warning("can't register event %q: %s", names.DownloadMgrEvent, err)
	}
	fmt.Printf("process %s registered event %s\n", process.Self(), names.DownloadMgrEvent)
	if err := process.RegisterEvent(names.DownloadEvent, names.DownloadEventMsg{}); err != nil {
		lib.Warning("can't register event %q: %s", names.DownloadEvent, err)
	}
	fmt.Printf("process %s registered event %s\n", process.Self(), names.DownloadEvent)
	if err := process.MonitorEvent(names.UpdateEvent); err != nil {
		return err
	}
	return nil
}

func (*managerSrv) HandleInfo(process *gen.ServerProcess, message etf.Term) gen.ServerStatus {
	switch message.(type) {
	case names.UpdateEventMsg:
		fmt.Printf("consumer %s got event: %s\r\n", process.Name(), message)
		msg := message.(names.UpdateEventMsg)
		event := names.DownloadMgrEventMsg{Pid: process.Self(), Msg: msg}

		fmt.Printf("... producing event: %#v\n", event)
		if err := process.SendEventMessage(names.DownloadMgrEvent, event); err != nil {
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
	msg, ok := message.(names.DownloadEventMsg)
	if ok == true {
		fmt.Printf("... producing event: %#v\n", msg)
		if err := process.SendEventMessage(names.DownloadEvent, msg); err != nil {
			fmt.Println("can't send event:", err)
		}
		//if msg.state == error -> try again

	} else {
		fmt.Println("unknown message", message)
	}
	return gen.ServerStatusOK
}

func (*managerSrv) Terminate(process *gen.ServerProcess, reason string) {
	process.UnregisterEvent(names.DownloadMgrEvent)
	process.UnregisterEvent(names.DownloadEvent)
	process.DemonitorEvent(names.UpdateEvent)
	fmt.Printf("[%s] Terminating process with reason %q\r\n", process.Name(), reason)
}
