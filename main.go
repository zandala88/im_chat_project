package main

import (
	"go.uber.org/zap"
	"im/config"
	_ "im/config"
	"im/internal/kafka"
	"im/internal/router"
	server "im/internal/websocket"
	"im/internal/websocket/constant"
	_ "im/public"
	"net/http"
	"time"
)

func main() {
	if config.Configs.ChannelType.ChannelType == constant.KAFKA {
		kafka.InitProducer(config.Configs.ChannelType.KafkaTopic, config.Configs.ChannelType.KafkaHosts)
		kafka.InitConsumer(config.Configs.ChannelType.KafkaHosts)
		go kafka.ConsumerMsg(server.ConsumerKafkaMsg)
	}

	zap.S().Info("start server", zap.String("start", "start web sever..."))

	newRouter := router.NewRouter()

	go server.MyServer.Start()

	s := &http.Server{
		Addr:           ":8888",
		Handler:        newRouter,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	err := s.ListenAndServe()
	if nil != err {
		zap.S().Error("server error", zap.Any("serverError", err))
	}
}
