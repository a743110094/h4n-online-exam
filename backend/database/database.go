package database

import (
	"log"
	"online-exam-system/config"
	"online-exam-system/models"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Init() {
	var err error
	dsn := config.GetConfig().DatabaseURL
	
	// 判断是否使用SQLite
	if strings.HasSuffix(dsn, ".db") {
		DB, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
	} else {
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
	}
	
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	
	log.Println("Database connected successfully")
}

func AutoMigrate() {
	err := DB.AutoMigrate(
		&models.User{},
		&models.Subject{},
		&models.Question{},
		&models.Paper{},
		&models.Exam{},
		&models.ExamRecord{},
		&models.Answer{},
		&models.AIChat{},
		&models.PracticeRecord{},
		&models.PracticeAnswer{},
		&models.PracticeRecommendation{},
	)
	
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	
	log.Println("Database migration completed")
	
	// 创建默认管理员账户
	createDefaultAdmin()
}

func createDefaultAdmin() {
	var count int64
	DB.Model(&models.User{}).Where("role = ?", models.RoleAdmin).Count(&count)
	
	if count == 0 {
		admin := models.User{
			Username: "admin",
			Email:    "admin@exam.com",
			Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password
			Role:     models.RoleAdmin,
			Name:     "系统管理员",
			IsActive: true,
		}
		
		if err := DB.Create(&admin).Error; err != nil {
			log.Printf("Failed to create default admin: %v", err)
		} else {
			log.Println("Default admin created: username=admin, password=password")
		}
	}
	
	// 创建默认教师用户
	createDefaultTeacher()
	// 创建默认学生用户
	createDefaultStudent()
}

func createDefaultTeacher() {
	var count int64
	DB.Model(&models.User{}).Where("username = ?", "teacher1").Count(&count)
	
	if count == 0 {
		teacher := models.User{
			Username: "teacher1",
			Email:    "teacher1@example.com",
			Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password
			Role:     models.RoleTeacher,
			Name:     "张老师",
			IsActive: true,
		}
		
		if err := DB.Create(&teacher).Error; err != nil {
			log.Printf("Failed to create default teacher: %v", err)
		} else {
			log.Println("Default teacher created: username=teacher1, password=password")
		}
	}
}

func createDefaultStudent() {
	var count int64
	DB.Model(&models.User{}).Where("username = ?", "student1").Count(&count)
	
	if count == 0 {
		student := models.User{
			Username: "student1",
			Email:    "student1@example.com",
			Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password
			Role:     models.RoleStudent,
			Name:     "王同学",
			IsActive: true,
		}
		
		if err := DB.Create(&student).Error; err != nil {
			log.Printf("Failed to create default student: %v", err)
		} else {
			log.Println("Default student created: username=student1, password=password")
		}
	}
}

func GetDB() *gorm.DB {
	return DB
}