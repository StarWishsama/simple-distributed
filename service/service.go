package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"simple-distributed/registry"
)

// Start 根据给定的服务名获取配置文件以启动一个 HTTP 服务器，并注册处理函数
func Start(ctx context.Context, reg registry.Registration, v *viper.Viper, regHandlersFunc func()) (context.Context, error) {
	host, port := v.GetString("server.host"), v.GetString("server.port")

	if host == "" || port == "" {
		return nil, errors.New("'server.host' or 'server.port' configuration is missing")
	}

	regHandlersFunc()
	ctx = startService(ctx, reg.ServiceName, host, port)
	if err := registry.RegisterService(reg); err != nil {
		return ctx, err
	}

	return ctx, nil
}

func startService(ctx context.Context, serviceName registry.ServiceName, host string, port string) context.Context {
	ctx, cancel := context.WithCancel(ctx)

	var srv http.Server

	srv.Addr = host + ":" + port

	// 启动 HTTP 服务器
	go func() {
		log.Println(srv.ListenAndServe())
		if err := unregisterService(host, port); err != nil {
			log.Printf("Failed to unregister service: %v\n", err)
		}
		cancel()
	}()

	// 监听用户退出信号
	go func() {
		fmt.Printf("Serving %s on %s:%s, Press any key to stop service.\n", serviceName, host, port)
		var input string
		fmt.Scanln(&input)
		if err := unregisterService(host, port); err != nil {
			log.Printf("Failed to unregister service: %v\n", err)
		}
		srv.Shutdown(ctx)
		cancel()
	}()

	return ctx
}

// unregisterService 从注册中心注销服务
func unregisterService(host, port string) error {
	url := fmt.Sprintf("http://%s:%s", host, port)
	if err := registry.UnregisterService(url); err != nil {
		return fmt.Errorf("failed to unregister service at %s: %w", url, err)
	}
	return nil
}
