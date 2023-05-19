package downloader

import (
	"github.com/ergo-services/ergo/etf"
	"github.com/ergo-services/ergo/gen"

	"resourceObserver/downloader/manager"
	"resourceObserver/names"
)

func CreateDownloaderSup() gen.SupervisorBehavior {
	return &downloaderSup{}
}

type downloaderSup struct {
	gen.Supervisor
}

func (*downloaderSup) Init(args ...etf.Term) (gen.SupervisorSpec, error) {
	spec := gen.SupervisorSpec{
		Name: names.DownloadSupName,
		Children: []gen.SupervisorChildSpec{
			gen.SupervisorChildSpec{
				Name:  names.DownloadMgrName,
				Child: manager.CreateManager(),
			},
			gen.SupervisorChildSpec{
				Name:  names.DownloadWrkSupName,
				Child: createDownloaderWorksersSup(),
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
