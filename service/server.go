package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
)

// Start 根据给定的服务名获取配置文件以启动一个 HTTP 服务器，并注册处理函数
func Start(ctx context.Context, serviceName, configName string, regHandlersFunc func()) (context.Context, error) {
	v := viper.New()

	v.SetConfigName(configName)
	v.AddConfigPath(".")
	v.AddConfigPath("./.env")
	v.SetConfigType("toml")

	v.SetDefault("server.host", "localhost")
	v.SetDefault("server.port", 8000)

	err := v.ReadInConfig()

	if err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			fmt.Printf("Config file for '%v' not found\n", serviceName)
		} else {
			fmt.Printf("A error occurred when reading config: %v\n", err)
		}

		os.Exit(1)
	}

	host, port := v.GetString("server.host"), v.GetString("server.port")

	if host == "" || port == "" {
		return nil, errors.New("'server.host' or 'server.port' configuration is missing")
	}

	regHandlersFunc()
	ctx = startService(ctx, serviceName, host, port)

	return ctx, nil
}

func startService(ctx context.Context, name string, host string, port string) context.Context {
	ctx, cancel := context.WithCancel(ctx)

	var srv http.Server

	srv.Addr = host + ":" + port

	// 启动 HTTP 服务器
	go func() {
		log.Println(srv.ListenAndServe())
		cancel()
	}()

	// 监听用户退出信号
	go func() {
		fmt.Printf("Serving %s on %s:%s, Press any key to stop service.\n", name, host, port)
		var input string
		fmt.Scanln(&input)
		srv.Shutdown(ctx)
		cancel()
	}()

	return ctx
}
