package main

import (
	"fmt"
	"github.com/vodolaz095/gossha"
	"runtime"
)

func main() {
	defer func() {
		err := recover()
		if err != nil {
			trace := make([]byte, 1024)
			count := runtime.Stack(trace, true)
			fmt.Println("====================================================")
			fmt.Println("Error! Error! Error!")
			fmt.Printf("Version: %v\nSubversion: %v\n\n", gossha.VERSION, gossha.SUBVERSION)
			fmt.Printf("Recover from panic: %s\n", err)
			fmt.Printf("Stack of %d bytes:\n %s\n", count, trace)
			fmt.Println("====================================================")
			fmt.Println("Please, report this error on `https://github.com/vodolaz095/gossha/issues` !\nThanks!")
		}
	}()

	cfg, err := gossha.InitConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	gossha.RuntimeConfig = &cfg
	err = gossha.InitDatabase(cfg.Driver, cfg.ConnectionString)
	if err != nil {
		panic(fmt.Errorf("Fatal error initializing database: %s \n", err))
	}
	gossha.ProcessConsoleCommand(cfg)
}
