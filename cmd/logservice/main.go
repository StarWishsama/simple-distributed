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
)

func main() {
	log.Run("./distributed.log")

	r := registry.Registration{
		ServiceName: "Log Service",
		ServiceURL:  "",
	}

	v, err := util.InitViper("logservice")

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

	ctx, err := service.Start(
		context.Background(),
		r,
		v,
		log.RegisterHandlers,
	)

	if err != nil {
		stlog.Fatalln(err)
	}

	<-ctx.Done() // 等待服务停止

	fmt.Println("Shutting down service...")
}
