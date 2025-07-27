package main

import (
	"log"
	"online-exam-system/cache"
	"online-exam-system/config"
	"online-exam-system/database"
	"online-exam-system/routes"
	"online-exam-system/services"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// 初始化配置
	config.Init()

	// 初始化数据库
	database.Init()

	// 初始化Redis缓存
	cache.InitRedis()

	// 初始化RabbitMQ
	if err := services.InitRabbitMQ(); err != nil {
		log.Printf("Failed to initialize RabbitMQ: %v", err)
		log.Println("Continuing without RabbitMQ...")
	} else {
		log.Println("RabbitMQ initialized successfully")
		
		// 启动消息消费者
		consumer := services.NewMessageConsumer()
		consumer.StartConsumers()
	}

	// 自动迁移数据库表
	database.AutoMigrate()

	// 启动缓存预热服务
	warmupService := services.NewWarmupService()
	warmupService.StartWarmupScheduler()
	
	// 执行初始数据预热
	go warmupService.PerformFullWarmup()

	// 设置Gin模式
	gin.SetMode(gin.DebugMode)

	// 创建路由
	r := gin.Default()

	// 设置CORS中间件
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Tenant-ID")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// 注册路由
	routes.SetupRoutes(r)

	// 设置优雅关闭
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Shutting down gracefully...")
		
		// 关闭RabbitMQ连接
		if err := services.CloseRabbitMQ(); err != nil {
			log.Printf("Error closing RabbitMQ: %v", err)
		}
		
		os.Exit(0)
	}()

	// 启动服务器
	port := config.GetConfig().Port
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on port %s", port)
	log.Fatal(r.Run(":" + port))
}