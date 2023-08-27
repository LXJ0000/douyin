package main

import (
	"douyin/config"
	"douyin/router"
	"fmt"
)

func main() {
	r := router.InitDouyinRouter()
	if err := r.Run(fmt.Sprintf(":%d", config.Info.Server.Port)); err != nil {
		return
	}
}
