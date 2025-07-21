package main

import (
	"context"
	"fmt"
	stlog "log"
	"simple-distributed/log"
	"simple-distributed/service"
)

func main() {
	log.Run("./distributed.log")

	ctx, err := service.Start(
		context.Background(),
		"Log Service",
		"logservice",
		log.RegisterHandlers,
	)

	if err != nil {
		stlog.Fatalln(err)
	}

	<-ctx.Done() // 等待服务停止

	fmt.Println("Shutting down service...")
}
