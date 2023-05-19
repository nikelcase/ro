package downloader

import (
	"github.com/ergo-services/ergo/etf"
	"github.com/ergo-services/ergo/gen"

	"resourceObserver/downloader/ftp"
	"resourceObserver/downloader/smb"
	"resourceObserver/downloader/web"
	"resourceObserver/names"
)

func createDownloaderWorksersSup() gen.SupervisorBehavior {
	return &downloaderWorkersSup{}
}

type downloaderWorkersSup struct {
	gen.Supervisor
}

func (*downloaderWorkersSup) Init(args ...etf.Term) (gen.SupervisorSpec, error) {
	spec := gen.SupervisorSpec{
		Name: names.DownloadWrkSupName,
		Children: []gen.SupervisorChildSpec{
			gen.SupervisorChildSpec{
				Name:  names.DownloadWebSupName,
				Child: web.CreateWebDownloader(),
			},
			gen.SupervisorChildSpec{
				Name:  names.DownloadSmbSupName,
				Child: smb.CreateSmbDownloader(),
			},
			gen.SupervisorChildSpec{
				Name:  names.DownloadFtpSupName,
				Child: ftp.CreateFtpDownloader(),
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
