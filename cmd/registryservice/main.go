package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"simple-distributed/registry"
	"simple-distributed/util"
)

func main() {
	service := registry.RegistryService{}

	http.Handle("/services", service)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	v, err := util.InitViper("regservice")

	if err != nil {
		fmt.Println("Error initializing configuration:", err)
		return
	}

	var srv http.Server
	host, port := v.GetString("server.host"), v.GetString("server.port")
	srv.Addr = host + ":" + port

	// 服务启动异常监听
	go func() {
		log.Println(srv.ListenAndServe())
		cancel()
	}()

	// 监听用户退出信号
	go func() {
		fmt.Printf("Serving %s on %s:%s, Press any key to stop service.\n", "Registry Service", host, port)
		var input string
		fmt.Scanln(&input)
		srv.Shutdown(ctx)
		cancel()
	}()

	<-ctx.Done() // 等待服务停止

	fmt.Println("Shutting down service...")
}
