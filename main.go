package main

import (
	"github.com/GGBusuioc/go-cache/cache"
	"github.com/GGBusuioc/go-cache/config"
	"github.com/GGBusuioc/go-cache/handler"
	"github.com/GGBusuioc/go-cache/logger"
	"github.com/GGBusuioc/go-cache/router"
	"github.com/GGBusuioc/go-cache/server"
)

func main() {
	c := config.NewConfig()
	l := logger.NewLogger(c)
	ch := cache.NewCache()
	h := handler.NewHandler(ch, *l)
	r := router.NewRouter(h, *l)

	s := server.NewServer(r, h, *l)
	s.Start()
}
