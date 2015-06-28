package main

import (
	"bitbucket.org/vodolaz095/gossha"
	"fmt"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	cfg, args, err := gossha.InitConfig()
	gossha.RuntimeConfig = &cfg
	gossha.PrintHelpOnCli()
	gossha.InitDatabase(cfg.Driver, cfg.ConnectionString)
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	if !cfg.Debug {
		defer func() {
			err := recover()
			fmt.Println("Oops!!!", err)
		}()
	}

	if len(args) == 0 {
		if cfg.Debug {
			go func() {
				fmt.Println(http.ListenAndServe("localhost:3000", nil))
			}()
		}
		if len(gossha.RuntimeConfig.BindTo) > 0 {
			for _, v := range gossha.RuntimeConfig.BindTo[:(len(gossha.RuntimeConfig.BindTo) - 1)] {
				go func() {
					gossha.StartSSHD(v)
				}()
			}
			gossha.StartSSHD(gossha.RuntimeConfig.BindTo[len(gossha.RuntimeConfig.BindTo)-1])
		} else {
			gossha.StartSSHD(fmt.Sprintf("0.0.0.0:%v", gossha.RuntimeConfig.Port))
		}
	} else {
		gossha.ProcessConsoleCommand(args)
	}
}
