package main

import (
	"fmt"
	"runtime"

	"github.com/kardianos/service"
	"github.com/thanhfphan/ws-hello/config"
	"github.com/thanhfphan/ws-hello/pkg/log"
)

func main() {
	appCfg, err := config.Build()
	if err != nil {
		panic(err)
	}

	logger, err := log.New(appCfg.LogConfig)
	if err != nil {
		panic(err)
	}

	cfg := &service.Config{
		Name:        "ws-hello",
		DisplayName: "Windows Service Hello Application",
		Description: "An example of a Windows Service Application",
	}

	p := &program{
		log: logger,
	}
	s, err := service.New(p, cfg)
	if err != nil {
		logger.Fatal(fmt.Sprintf("service.New got err=%v", err))
		return
	}

	err = s.Run()
	if err != nil {
		logger.Fatal(fmt.Sprintf("s.Run got err=%v", err))
		return
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("caught panic", r)
			}
			logger.Stop()
		}()
		defer func() {
			// log panic and then re-raise the panic
			logger.StopOnPanic()
		}()
	}()
}

type program struct {
	log log.Logger
}

func (p program) Start(s service.Service) error {
	p.log.Info("Starting the program ...")
	p.log.Info(runtime.GOOS)

	return nil
}

func (p program) Stop(s service.Service) error {
	p.log.Info("Stopping the program ...")
	return nil
}
