package main

import (
	"flag"
	"fmt"

	"github.com/ergo-services/ergo"
	"github.com/ergo-services/ergo/etf"
	"github.com/ergo-services/ergo/gen"
	"github.com/ergo-services/ergo/node"

	"resourceObserver/downloader"
	"resourceObserver/names"
	"resourceObserver/notifier"
	"resourceObserver/observers"
)

const appname = "ResourceObserver"

func main() {
	//RunApp()
	RunSup()
}

func RunApp() {
	flag.Parse()

	apps := []gen.ApplicationBehavior{
		CreateRoApp(),
	}
	opts := node.Options{
		StaticRoutesOnly: true,
		Applications:     apps,
	}
	fmt.Println("Start node " + appname + "@127.0.0.1")
	roNode, err := ergo.StartNode(appname+"@127.0.0.1", appname, opts)
	if err != nil {
		panic(err)
	}
	fmt.Println("Started application process")
	roNode.Wait()
}

func RunSup() {
	flag.Parse()

	fmt.Println("Start node " + appname + "@127.0.0.1")
	roNode, err := ergo.StartNode(appname+"@127.0.0.1", appname, node.Options{
		StaticRoutesOnly: true,
	})
	if err != nil {
		panic(err)
	}
	appSup := CreateAppSup()

	sup, err := roNode.Spawn(appname+"AppSup", gen.ProcessOptions{}, appSup)
	if err != nil {
		panic(err)
	}
	fmt.Println("Started supervisor process", sup.Self())
	sup.Wait()
	//roNode.Wait()
}

func CreateRoApp() gen.ApplicationBehavior {
	return &roApp{}
}

type roApp struct {
	gen.Application
}

func (*roApp) Load(args ...etf.Term) (gen.ApplicationSpec, error) {
	return gen.ApplicationSpec{
		Name:        appname + "App",
		Description: "Resource Observer App",
		Version:     "v0.0.1",
		Children: []gen.ApplicationChildSpec{
			gen.ApplicationChildSpec{
				Name:  names.AppSupName,
				Child: CreateAppSup(),
			},
		},
	}, nil
}

func (*roApp) Start(process gen.Process, args ...etf.Term) {
	fmt.Printf("Application started with Pid %s!\n", process.Self())
}

func CreateAppSup() gen.SupervisorBehavior {
	return &appSup{}
}

type appSup struct {
	gen.Supervisor
}

func (*appSup) Init(args ...etf.Term) (gen.SupervisorSpec, error) {
	spec := gen.SupervisorSpec{
		Name: names.AppSupName,
		Children: []gen.SupervisorChildSpec{
			gen.SupervisorChildSpec{
				Name:  names.ObsSupName,
				Child: observers.CreateObserverSup(),
			},
			gen.SupervisorChildSpec{
				Name:  names.DownloadSupName,
				Child: downloader.CreateDownloaderSup(),
			},
			gen.SupervisorChildSpec{
				Name:  names.NotifySupName,
				Child: notifier.CreateNotifierSup(),
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
