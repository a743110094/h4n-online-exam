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

	practiceRecordID := "257"
	userID := uint(3)
	tenantID := uint(100)

	fmt.Printf("=== 测试错题查询API逻辑 ===\n")

	// 使用与API相同的查询逻辑
	type WrongQuestionDetail struct {
		ID             uint   `json:"id"`
		QuestionID     uint   `json:"question_id"`
		Title          string `json:"title"`
		Content        string `json:"content"`
		Options        string `json:"options"`
		Type           string `json:"type"`
		Difficulty     int    `json:"difficulty"`
		SubjectName    string `json:"subject_name"`
		UserAnswer     string `json:"user_answer"`
		CorrectAnswer  string `json:"correct_answer"`
		Explanation    string `json:"explanation"`
		TimeSpent      int    `json:"time_spent"`
		CreatedAt      string `json:"created_at"`
	}

	// 构建查询条件（与修复后的API相同）
	query := database.DB.Table("practice_answers").
		Select("practice_answers.*, questions.title, questions.content, questions.options, questions.answer as correct_answer, questions.explanation, questions.type, questions.difficulty, subjects.name as subject_name").
		Joins("JOIN practice_records ON practice_answers.practice_record_id = practice_records.id").
		Joins("JOIN questions ON practice_answers.question_id = questions.id").
		Joins("JOIN subjects ON questions.subject_id = subjects.id").
		Where("practice_records.user_id = ? AND practice_answers.is_correct = ? AND practice_records.tenant_id = ? AND practice_answers.tenant_id = ? AND questions.tenant_id = ? AND subjects.tenant_id = ?", userID, false, tenantID, tenantID, tenantID, tenantID)

	// 如果指定了练习记录ID，只获取该次练习的错题
	if practiceRecordID != "" {
		query = query.Where("practice_answers.practice_record_id = ?", practiceRecordID)
	}

	// 获取总数
	var total int64
	query.Count(&total)
	fmt.Printf("错题查询总数: %d\n", total)

	// 获取详细数据
	var wrongQuestions []WrongQuestionDetail
	query.Select("practice_answers.id, practice_answers.question_id, questions.title, questions.content, questions.options, questions.type, questions.difficulty, subjects.name as subject_name, practice_answers.answer as user_answer, questions.answer as correct_answer, questions.explanation, practice_answers.time_spent, practice_answers.created_at").
		Order("practice_answers.created_at DESC").
		Limit(10).
		Scan(&wrongQuestions)

	fmt.Printf("错题详细数据数量: %d\n", len(wrongQuestions))

	for i, wq := range wrongQuestions {
		fmt.Printf("错题 %d: QuestionID=%d, Title='%s', UserAnswer='%s', CorrectAnswer='%s'\n", i+1, wq.QuestionID, wq.Title, wq.UserAnswer, wq.CorrectAnswer)
	}

	if len(wrongQuestions) > 0 {
		fmt.Printf("\n✅ 修复成功！现在可以正确查询到错题了。\n")
	} else {
		fmt.Printf("\n❌ 仍然没有查询到错题，可能还有其他问题。\n")
	}
}