package main

import (
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

	fmt.Println("=== 修复华为AI认证题目中的空答案 ===")

	// 1. 查找华为AI认证科目
	var subject models.Subject
	err := database.DB.Where("name = ?", "华为AI认证(HCIA-AI)").First(&subject).Error
	if err != nil {
		fmt.Printf("未找到华为AI认证科目: %v\n", err)
		return
	}
	fmt.Printf("找到科目: ID=%d, Name=%s\n", subject.ID, subject.Name)

	// 2. 查找空答案的题目
	var emptyAnswerQuestions []models.Question
	err = database.DB.Where("subject_id = ? AND (answer = '' OR answer IS NULL)", subject.ID).Find(&emptyAnswerQuestions).Error
	if err != nil {
		fmt.Printf("查询空答案题目失败: %v\n", err)
		return
	}

	fmt.Printf("\n找到 %d 道空答案题目:\n", len(emptyAnswerQuestions))

	for _, q := range emptyAnswerQuestions {
		fmt.Printf("\n题目ID: %d\n", q.ID)
		fmt.Printf("题型: %s\n", q.Type)
		fmt.Printf("标题: %s\n", q.Title)
		fmt.Printf("当前答案: '%s'\n", q.Answer)

		// 根据题目内容判断如何处理
		if q.ID == 626 { // 这是那道匹配题
			fmt.Println("这是一道匹配题，建议删除或设置为不可用状态")
			
			// 选项1：删除这道题
			// err = database.DB.Delete(&q).Error
			// if err != nil {
			//     fmt.Printf("删除题目失败: %v\n", err)
			// } else {
			//     fmt.Println("已删除匹配题")
			// }

			// 选项2：设置为不可用状态
			err = database.DB.Model(&q).Update("status", "disabled").Error
			if err != nil {
				fmt.Printf("更新题目状态失败: %v\n", err)
			} else {
				fmt.Println("已将匹配题设置为不可用状态")
			}
		} else {
			// 其他空答案题目的处理逻辑
			fmt.Println("需要手动检查和修复")
		}
	}

	// 3. 验证修复结果
	var remainingEmptyQuestions []models.Question
	err = database.DB.Where("subject_id = ? AND (answer = '' OR answer IS NULL) AND status != 'disabled'", subject.ID).Find(&remainingEmptyQuestions).Error
	if err != nil {
		fmt.Printf("验证失败: %v\n", err)
		return
	}

	fmt.Printf("\n修复完成！剩余空答案题目数: %d\n", len(remainingEmptyQuestions))

	// 4. 重新统计题目情况
	var totalQuestions int64
	var activeQuestions int64
	database.DB.Model(&models.Question{}).Where("subject_id = ?", subject.ID).Count(&totalQuestions)
	database.DB.Model(&models.Question{}).Where("subject_id = ? AND status = 'published'", subject.ID).Count(&activeQuestions)

	fmt.Printf("总题目数: %d\n", totalQuestions)
	fmt.Printf("可用题目数: %d\n", activeQuestions)
}