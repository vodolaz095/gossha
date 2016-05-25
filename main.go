package main

import (
	"fmt"
	"runtime"

	"github.com/vodolaz095/gossha/cli"
	"github.com/vodolaz095/gossha/config"
	"github.com/vodolaz095/gossha/models"
)

func main() {
	defer func() {
		err := recover()
		if err != nil {
			trace := make([]byte, 1024)
			count := runtime.Stack(trace, true)
			fmt.Println("====================================================")
			fmt.Println("Error! Error! Error!")
			//			fmt.Printf("Version: %v\nSubversion: %v\n\n", gossha.VERSION, gossha.SUBVERSION)
			fmt.Printf("Recover from panic: %s\n", err)
			fmt.Printf("Stack of %d bytes:\n %s\n", count, trace)
			fmt.Println("====================================================")
			fmt.Println("Please, report this error on `https://github.com/vodolaz095/gossha/issues` !\nThanks!")
		}
	}()

	cfg, err := config.InitConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error while creating config file: %s \n", err))
	}

	err = models.InitDatabase(cfg.Driver, cfg.ConnectionString, cfg.Debug)
	if err != nil {
		panic(fmt.Errorf("Fatal error while initializing database: %s \n", err))
	}

	cli.ProcessConsoleCommand()
}
