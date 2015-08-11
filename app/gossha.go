package main

import (
	"fmt"
	"github.com/vodolaz095/gossha"
)

func main() {
	//	fmt.Println(gossha.Greet())

	cfg, _, err := gossha.InitConfig()
	gossha.RuntimeConfig = &cfg
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
	gossha.ProcessConsoleCommand(cfg)
}
