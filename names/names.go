package names

import (
	"time"

	"github.com/ergo-services/ergo/etf"
	"github.com/ergo-services/ergo/gen"
)

const AppSupName = "appSup"

const (
	ObsSupName              = "obsSup"
	ObsSchedName            = "obsSched"
	ObsWrkSupName           = "obsWrkSup"
	ObsWebSupName           = "obsWebSup"
	ObsSmbSupName           = "obsSmbSup"
	ObsFtpSupName           = "obsFtpSup"
	SchedEvent    gen.Event = "observe"
)

type SchedEventMsg struct {
	Pid etf.Pid
}

const UpdateEvent gen.Event = "update"

type ObsState int

const (
	ObsSuccess ObsState = iota
	ObsError
)

type UpdateEventMsg struct {
	Name      string
	State     ObsState
	Title     string
	Text      string
	BaseName  string
	FullAddr  string
	ShortAddr string
	Time      time.Time
}

const (
	NotifySupName                  = "notifySup"
	NotifyMgrName                  = "notifyMgr"
	NotifyWrkSupName               = "notifyWrkSup"
	NotifyWebSupName               = "notifyWebSup"
	NotifySshSupName               = "notifySshSup"
	NotifyStdoutSupName            = "notifyStdoutSup"
	NotifyTelnetSupName            = "notifyTelnetSup"
	NotifyWinsentSupName           = "obsWinsentSup"
	NotifyMgrEvent       gen.Event = "notify"
)

type NotifyMgrEventMsg struct {
	Pid etf.Pid
	Msg interface{}
}

type NtfState int

const (
	NtfSuccess NtfState = iota
	NtfError
)

type NotifyMsg struct {
	Pid   etf.Pid
	State NtfState
	Name  string
	Text  string
}

const (
	DownloadSupName              = "downloadSup"
	DownloadMgrName              = "downloadMgr"
	DownloadWrkSupName           = "downloadWrkSup"
	DownloadWebSupName           = "downloadWebSup"
	DownloadSmbSupName           = "downloadSmbSup"
	DownloadFtpSupName           = "downloadFtpSup"
	DownloadMgrEvent   gen.Event = "download"
)

type DownloadMgrEventMsg struct {
	Pid etf.Pid
	Msg UpdateEventMsg
}

const DownloadEvent gen.Event = "success"

type DdState int

const (
	DdSuccess DdState = iota
	DdError
)

type DownloadEventMsg struct {
	State  DdState
	Unpack bool
	Title  string
	Path   string
}
