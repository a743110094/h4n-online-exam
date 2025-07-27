package main

import (
	"fmt"
	"log"
	"online-exam-system/config"
	"online-exam-system/database"
	"online-exam-system/models"
	"online-exam-system/utils"

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

	practiceRecordID := uint(257)

	fmt.Printf("=== 调试练习记录ID %d 的错题查询 ===\n", practiceRecordID)

	// 1. 先查询所有练习记录，找到ID为257的记录
	var allRecords []models.PracticeRecord
	database.DB.Find(&allRecords)
	fmt.Printf("数据库中总练习记录数: %d\n", len(allRecords))

	var practiceRecord models.PracticeRecord
	var found bool
	for _, record := range allRecords {
		if record.ID == practiceRecordID {
			practiceRecord = record
			found = true
			break
		}
		if record.ID >= 250 { // 显示ID接近257的记录
			fmt.Printf("练习记录: ID=%d, UserID=%d, TenantID=%d, IsCompleted=%t\n", record.ID, record.UserID, record.TenantID, record.IsCompleted)
		}
	}

	if !found {
		fmt.Printf("练习记录ID %d 不存在\n", practiceRecordID)
		return
	}

	fmt.Printf("找到练习记录: ID=%d, UserID=%d, TenantID=%d, IsCompleted=%t\n", practiceRecord.ID, practiceRecord.UserID, practiceRecord.TenantID, practiceRecord.IsCompleted)
	tenantID := practiceRecord.TenantID

	// 2. 检查该练习记录的所有答题记录
	var allAnswers []models.PracticeAnswer
	err := utils.WithTenant(database.DB, tenantID).Where("practice_record_id = ?", practiceRecordID).Find(&allAnswers).Error
	if err != nil {
		fmt.Printf("查询答题记录失败: %v\n", err)
		return
	}
	fmt.Printf("总答题记录数: %d\n", len(allAnswers))

	// 3. 统计正确和错误答案
	correctCount := 0
	wrongCount := 0
	for _, answer := range allAnswers {
		if answer.IsCorrect != nil {
			if *answer.IsCorrect {
				correctCount++
			} else {
				wrongCount++
			}
			fmt.Printf("答题记录: QuestionID=%d, Answer='%s', IsCorrect=%t\n", answer.QuestionID, answer.Answer, *answer.IsCorrect)
		} else {
			fmt.Printf("答题记录: QuestionID=%d, Answer='%s', IsCorrect=nil\n", answer.QuestionID, answer.Answer)
		}
	}
	fmt.Printf("正确答案数: %d, 错误答案数: %d\n", correctCount, wrongCount)

	// 4. 使用与API相同的查询逻辑查询错题
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
	}

	var wrongQuestions []WrongQuestionDetail
	query := utils.WithTenant(database.DB, tenantID).Table("practice_answers").
		Select("practice_answers.*, questions.title, questions.content, questions.options, questions.answer as correct_answer, questions.explanation, questions.type, questions.difficulty, subjects.name as subject_name").
		Joins("JOIN practice_records ON practice_answers.practice_record_id = practice_records.id").
		Joins("JOIN questions ON practice_answers.question_id = questions.id").
		Joins("JOIN subjects ON questions.subject_id = subjects.id").
		Where("practice_records.user_id = ? AND practice_answers.is_correct = ? AND practice_records.tenant_id = ? AND practice_answers.tenant_id = ? AND questions.tenant_id = ? AND subjects.tenant_id = ?", practiceRecord.UserID, false, tenantID, tenantID, tenantID, tenantID).
		Where("practice_answers.practice_record_id = ?", practiceRecordID)

	// 获取总数
	var total int64
	query.Count(&total)
	fmt.Printf("错题查询总数: %d\n", total)

	// 获取详细数据
	query.Order("practice_answers.created_at DESC").Limit(10).Scan(&wrongQuestions)
	fmt.Printf("错题详细数据数量: %d\n", len(wrongQuestions))

	for i, wq := range wrongQuestions {
		fmt.Printf("错题 %d: QuestionID=%d, Title='%s', UserAnswer='%s', CorrectAnswer='%s'\n", i+1, wq.QuestionID, wq.Title, wq.UserAnswer, wq.CorrectAnswer)
	}

	// 5. 检查租户ID是否正确
	fmt.Printf("\n=== 检查租户ID ===\n")
	var practiceAnswerWithTenant models.PracticeAnswer
	err2 := database.DB.Where("practice_record_id = ?", practiceRecordID).First(&practiceAnswerWithTenant).Error
	if err2 == nil {
		fmt.Printf("答题记录的租户ID: %d\n", practiceAnswerWithTenant.TenantID)
	}

	var practiceRecordWithTenant models.PracticeRecord
	err3 := database.DB.Where("id = ?", practiceRecordID).First(&practiceRecordWithTenant).Error
	if err3 == nil {
		fmt.Printf("练习记录的租户ID: %d\n", practiceRecordWithTenant.TenantID)
	}
}