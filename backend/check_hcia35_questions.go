package main

import (
	"encoding/json"
	"fmt"
	"log"
	"online-exam-system/config"
	"online-exam-system/database"
	"online-exam-system/models"

	"github.com/joho/godotenv"
)

func main() {
	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// 初始化配置
	config.Init()

	// 初始化数据库连接
	database.Init()

	fmt.Println("=== 检查华为AI认证(HCIA-AI)题目导入情况 ===")

	// 1. 查找华为AI认证科目
	var subject models.Subject
	err := database.DB.Where("name = ?", "华为AI认证(HCIA-AI)").First(&subject).Error
	if err != nil {
		fmt.Printf("未找到华为AI认证科目: %v\n", err)
		return
	}
	fmt.Printf("找到科目: ID=%d, Name=%s\n", subject.ID, subject.Name)

	// 2. 统计该科目下的题目数量和类型分布（只查询可用题目）
	var questions []models.Question
	err = database.DB.Where("subject_id = ? AND status = 'published'", subject.ID).Find(&questions).Error
	if err != nil {
		fmt.Printf("查询题目失败: %v\n", err)
		return
	}

	fmt.Printf("\n总题目数量: %d\n", len(questions))

	// 统计题型分布
	typeCount := make(map[string]int)
	for _, q := range questions {
		typeCount[string(q.Type)]++
	}

	fmt.Println("\n题型分布:")
	for qType, count := range typeCount {
		fmt.Printf("- %s: %d题\n", qType, count)
	}

	// 3. 检查前10道题的详细信息
	fmt.Println("\n=== 前10道题详细信息 ===")
	for i, q := range questions {
		if i >= 10 {
			break
		}

		fmt.Printf("\n题目 %d:\n", i+1)
		fmt.Printf("ID: %d\n", q.ID)
		fmt.Printf("类型: %s\n", q.Type)
		fmt.Printf("标题: %s\n", q.Title)
		fmt.Printf("内容: %s\n", q.Content)
		fmt.Printf("答案: %s\n", q.Answer)
		fmt.Printf("难度: %d\n", q.Difficulty)
		fmt.Printf("分数: %d\n", q.Score)

		// 解析选项（如果有）
		if q.Options != "" {
			var options []string
			if err := json.Unmarshal([]byte(q.Options), &options); err != nil {
				fmt.Printf("选项解析失败: %v\n", err)
				fmt.Printf("原始选项: %s\n", q.Options)
			} else {
				fmt.Println("选项:")
				for j, option := range options {
					fmt.Printf("  %c. %s\n", 'A'+j, option)
				}
			}
		}

		// 显示解释（如果有）
		if q.Explanation != "" {
			fmt.Printf("解释: %s\n", q.Explanation)
		}

		fmt.Println("---")
	}

	// 4. 检查每种题型的答案格式
	fmt.Println("\n=== 各题型答案格式检查 ===")
	for qType, count := range typeCount {
		fmt.Printf("\n%s (%d题):\n", qType, count)
		
		// 找到该类型的前3道题
		typeQuestions := make([]models.Question, 0)
		for _, q := range questions {
			if string(q.Type) == qType {
				typeQuestions = append(typeQuestions, q)
				if len(typeQuestions) >= 3 {
					break
				}
			}
		}

		for i, q := range typeQuestions {
			fmt.Printf("  示例%d: 答案='%s'\n", i+1, q.Answer)
		}
	}

	// 5. 检查是否有空答案或异常数据（只检查可用题目）
	fmt.Println("\n=== 数据完整性检查 ===")
	emptyAnswers := 0
	emptyContent := 0
	emptyTitle := 0

	for _, q := range questions {
		if q.Answer == "" {
			emptyAnswers++
		}
		if q.Content == "" {
			emptyContent++
		}
		if q.Title == "" {
			emptyTitle++
		}
	}

	fmt.Printf("空答案题目数: %d\n", emptyAnswers)
	fmt.Printf("空内容题目数: %d\n", emptyContent)
	fmt.Printf("空标题题目数: %d\n", emptyTitle)

	if emptyAnswers > 0 {
		fmt.Println("\n警告: 发现空答案题目!")
		for _, q := range questions {
			if q.Answer == "" {
				fmt.Printf("题目ID %d: %s\n", q.ID, q.Title)
			}
		}
	}

	// 6. 检查禁用的题目
	var disabledQuestions []models.Question
	err = database.DB.Where("subject_id = ? AND status = 'disabled'", subject.ID).Find(&disabledQuestions).Error
	if err == nil {
		fmt.Printf("\n已禁用题目数: %d\n", len(disabledQuestions))
		for _, q := range disabledQuestions {
			fmt.Printf("禁用题目ID %d: %s\n", q.ID, q.Title)
		}
	}
}