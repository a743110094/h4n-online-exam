package main

import (
	"log"
	"online-exam-system/config"
	"online-exam-system/database"
)

func main() {
	log.Println("Starting database migration...")
	
	// 初始化配置
	config.Init()
	
	// 初始化数据库连接
	database.Init()
	
	// 运行自动迁移
	database.AutoMigrate()
	
	log.Println("Database migration completed successfully!")
}