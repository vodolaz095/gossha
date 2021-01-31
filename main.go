package main

import (
	"fmt"
	"runtime"

	"github.com/vodolaz095/gossha/cli"
	"github.com/vodolaz095/gossha/config"
	"github.com/vodolaz095/gossha/models"
	"github.com/vodolaz095/gossha/version"
)

func main() {
	defer func() {
		err := recover()
		if err != nil {
			trace := make([]byte, 1024)
			count := runtime.Stack(trace, true)
			fmt.Println("====================================================")
			fmt.Println("Error! Error! Error!")
			fmt.Printf("Version: %v\n", version.Version)
			fmt.Printf("Recover from panic: %s\n", err)
			fmt.Printf("Stack of %d bytes:\n %s\n", count, trace)
			fmt.Println("====================================================")
			fmt.Println("Please, report this error on `https://github.com/vodolaz095/gossha/issues` !\nThanks!")
		}
	}()

	cfg, err := config.InitConfig()
	if err != nil {
		panic(fmt.Errorf("%s - creating config file", err))
	}

	err = models.InitDatabase(cfg.Driver, cfg.ConnectionString, cfg.Debug)
	if err != nil {
		panic(fmt.Errorf("%s - creating config file", err))
	}
	cli.ProcessConsoleCommand()
}
