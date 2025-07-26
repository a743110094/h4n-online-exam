package main

import (
	"log"
	"online-exam-system/config"
	"online-exam-system/database"
	"online-exam-system/models"
	"time"
)

func main() {
	// 初始化配置
	config.Init()
	
	// 初始化数据库
	database.Init()
	
	// 插入推荐练习数据
	insertPracticeRecommendations()
	
	// 插入示例练习记录
	insertSamplePracticeRecords()
	
	log.Println("练习数据初始化完成")
}

func insertPracticeRecommendations() {
	recommendations := []models.PracticeRecommendation{
		{
			Title:          "二叉树遍历专项练习",
			Description:    "深入理解二叉树的前序、中序、后序遍历",
			SubjectID:      1,
			Difficulty:     3,
			QuestionCount:  15,
			EstimatedTime:  20,
			Rating:         4.8,
			KnowledgePoint: "二叉树遍历",
			QuestionTypes:  `["single_choice", "multiple_choice"]`,
			IsActive:       true,
		},
		{
			Title:          "排序算法基础",
			Description:    "掌握冒泡、选择、插入等基础排序算法",
			SubjectID:      1,
			Difficulty:     2,
			QuestionCount:  20,
			EstimatedTime:  25,
			Rating:         4.6,
			KnowledgePoint: "排序算法",
			QuestionTypes:  `["single_choice", "true_false"]`,
			IsActive:       true,
		},
		{
			Title:          "进程调度算法",
			Description:    "理解FCFS、SJF、RR等调度算法",
			SubjectID:      2,
			Difficulty:     4,
			QuestionCount:  12,
			EstimatedTime:  30,
			Rating:         4.9,
			KnowledgePoint: "进程调度",
			QuestionTypes:  `["single_choice", "short_answer"]`,
			IsActive:       true,
		},
		{
			Title:          "数据库索引优化",
			Description:    "掌握数据库索引的创建和优化技巧",
			SubjectID:      3,
			Difficulty:     3,
			QuestionCount:  18,
			EstimatedTime:  35,
			Rating:         4.7,
			KnowledgePoint: "数据库索引",
			QuestionTypes:  `["single_choice", "multiple_choice"]`,
			IsActive:       true,
		},
		{
			Title:          "网络协议基础",
			Description:    "理解TCP/IP、HTTP等网络协议",
			SubjectID:      4,
			Difficulty:     2,
			QuestionCount:  25,
			EstimatedTime:  40,
			Rating:         4.5,
			KnowledgePoint: "网络协议",
			QuestionTypes:  `["single_choice", "true_false"]`,
			IsActive:       true,
		},
	}
	
	for _, rec := range recommendations {
		// 检查是否已存在
		var count int64
		database.DB.Model(&models.PracticeRecommendation{}).Where("title = ?", rec.Title).Count(&count)
		if count == 0 {
			if err := database.DB.Create(&rec).Error; err != nil {
				log.Printf("创建推荐练习失败: %v", err)
			} else {
				log.Printf("创建推荐练习: %s", rec.Title)
			}
		}
	}
}

func insertSamplePracticeRecords() {
	// 为学生用户(ID=3)创建一些示例练习记录
	records := []models.PracticeRecord{
		{
			UserID:       3,
			SubjectID:    1,
			Title:        "栈和队列基础",
			Description:  "练习栈和队列的基本操作",
			QuestionIDs:  `[1,2,3,4,5,6,7,8,9,10]`,
			TotalCount:   10,
			CorrectCount: 8,
			WrongCount:   2,
			Score:        80,
			Duration:     900,
			Difficulty:   2,
			PracticeType: "sequence",
			IsCompleted:  true,
			CreatedAt:    time.Now().Add(-2 * time.Hour),
			UpdatedAt:    time.Now().Add(-2 * time.Hour),
		},
		{
			UserID:       3,
			SubjectID:    1,
			Title:        "递归算法练习",
			Description:  "练习递归算法的实现",
			QuestionIDs:  `[11,12,13,14,15,16,17,18]`,
			TotalCount:   8,
			CorrectCount: 7,
			WrongCount:   1,
			Score:        87,
			Duration:     720,
			Difficulty:   3,
			PracticeType: "random",
			IsCompleted:  true,
			CreatedAt:    time.Now().Add(-24 * time.Hour),
			UpdatedAt:    time.Now().Add(-24 * time.Hour),
		},
		{
			UserID:       3,
			SubjectID:    2,
			Title:        "内存管理",
			Description:  "操作系统内存管理相关题目",
			QuestionIDs:  `[19,20,21,22,23,24,25,26,27,28,29,30,31,32,33]`,
			TotalCount:   15,
			CorrectCount: 13,
			WrongCount:   2,
			Score:        86,
			Duration:     1080,
			Difficulty:   3,
			PracticeType: "sequence",
			IsCompleted:  true,
			CreatedAt:    time.Now().Add(-48 * time.Hour),
			UpdatedAt:    time.Now().Add(-48 * time.Hour),
		},
	}
	
	for _, record := range records {
		// 检查是否已存在
		var count int64
		database.DB.Model(&models.PracticeRecord{}).Where("user_id = ? AND title = ?", record.UserID, record.Title).Count(&count)
		if count == 0 {
			if err := database.DB.Create(&record).Error; err != nil {
				log.Printf("创建练习记录失败: %v", err)
			} else {
				log.Printf("创建练习记录: %s", record.Title)
			}
		}
	}
}