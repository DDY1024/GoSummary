package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// 注册初始路由
	router.GET("/hello", func(c *gin.Context) {
		c.String(200, "Hello, World!")
	})

	// 创建 HTTP 服务器
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	quit := make(chan struct{}, 1)

	// 启动服务器
	go func() {
		// http.ErrServerClosed: http server 正常关闭
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("Failed to start server:", err)
		}
		quit <- struct{}{}
	}()

	go func() {
		time.Sleep(5 * time.Second)

		fmt.Println("Start Shutdown!")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		// 关闭服务器
		// gracefully shutdown
		// 1. close listener
		// 2. close idle connection
		// 3. close active connection
		if err := server.Shutdown(ctx); err != nil {
			fmt.Println("Failed to shutdown server:", err)
		}
		fmt.Println("Server gracefully stopped")
	}()

	<-quit
	fmt.Println("End!")

	if err := server.ListenAndServe(); err != nil {
		fmt.Println("Failed to start server:", err)
	}

	// 监听操作系统的信号
	// quit := make(chan os.Signal, 1)
	// signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// // 等待信号
	// <-quit
}
