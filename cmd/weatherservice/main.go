package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	stlog "log"
	"os"
	"simple-distributed/log"
	"simple-distributed/registry"
	"simple-distributed/service"
	"simple-distributed/util"
	"simple-distributed/weather"
)

func main() {
	r := registry.Registration{
		ServiceName: registry.Weather,
		RequiredServices: []registry.ServiceName{
			registry.Log,
		},
	}

	v, err := util.InitViper("weatherservice")

	if err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			fmt.Printf("Config file for '%v' not found\n", r.ServiceName)
		} else {
			fmt.Printf("A error occurred when reading config: %v\n", err)
		}

		os.Exit(1)
	}

	r.ServiceURL = fmt.Sprintf("http://%v:%v", v.GetString("server.host"), v.GetString("server.port"))
	r.ServiceUpdateURL = fmt.Sprintf("%v/services", r.ServiceURL)

	ctx, err := service.Start(
		context.Background(),
		r,
		v,
		weather.RegisterHandlers,
	)

	if err != nil {
		stlog.Fatalln(err)
	}

	if logProv, err := registry.GetProvider(registry.Log); err == nil {
		fmt.Printf("Log service found: %s\n", logProv)
		log.SetClientLogger(logProv, r.ServiceName)
	}

	<-ctx.Done() // 等待服务停止

	fmt.Println("Shutting down service...")
}
