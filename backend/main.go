package main

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"campus-message-board/internal/config"
	"campus-message-board/internal/handler"
	"campus-message-board/internal/model"
	"campus-message-board/internal/repository"
	"campus-message-board/internal/service"
	"github.com/gin-gonic/gin"
)

//go:embed schema.sql
var schemaSQL string

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	settings, err := config.Load()
	if err != nil {
		return err
	}
	database, err := sql.Open("mysql", settings.DSN())
	if err != nil {
		return fmt.Errorf("创建数据库连接失败：%w", err)
	}
	defer database.Close()
	database.SetMaxOpenConns(5)
	database.SetMaxIdleConns(5)
	database.SetConnMaxLifetime(3 * time.Minute)

	// 先检查指定数据库连接并创建缺失的表，成功后才开放 HTTP 端口。
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := database.PingContext(ctx); err != nil {
		return fmt.Errorf("连接项目数据库失败：%w", err)
	}
	if _, err := database.ExecContext(ctx, schemaSQL); err != nil {
		return fmt.Errorf("初始化 messages 表失败：%w", err)
	}

	// 显式连接三个具体模块，顺着这些构造函数即可找到下一层。
	messageRepository := repository.New(database)
	messageService := service.New(messageRepository)
	messageHandler := handler.New(messageService)

	router := gin.New()
	router.Use(gin.Logger(), gin.CustomRecovery(func(c *gin.Context, recovered any) {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.ErrorResponse{Message: "服务暂时出错，请稍后重试"})
	}))
	if err := router.SetTrustedProxies(nil); err != nil {
		return err
	}
	router.HandleMethodNotAllowed = true
	router.NoRoute(func(c *gin.Context) { c.JSON(http.StatusNotFound, model.ErrorResponse{Message: "接口不存在"}) })
	router.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, model.ErrorResponse{Message: "请求方法不支持"})
	})
	router.GET("/api/messages", messageHandler.List)
	router.POST("/api/messages", messageHandler.Create)

	address := net.JoinHostPort("127.0.0.1", settings.Port)
	log.Printf("Go/Gin 后端启动：http://%s；数据库：%s:%s/%s", address, settings.DBHost, settings.DBPort, settings.DBName)
	server := &http.Server{Addr: address, Handler: router, ReadHeaderTimeout: 5 * time.Second}
	return server.ListenAndServe()
}
