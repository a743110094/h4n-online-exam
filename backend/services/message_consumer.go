package services

import (
	"encoding/json"
	"fmt"
	"log"
	"online-exam-system/database"
	"online-exam-system/models"
	"strconv"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
)

// MessageConsumer 消息消费者
type MessageConsumer struct {
	db      *gorm.DB
	channel *amqp.Channel
}

// NewMessageConsumer 创建新的消息消费者
func NewMessageConsumer() *MessageConsumer {
	db := database.GetDB()
	rabbitmq := GetRabbitMQService()
	if rabbitmq == nil {
		log.Fatal("RabbitMQ service not initialized")
	}

	return &MessageConsumer{
		db:      db,
		channel: rabbitmq.channel,
	}
}

// StartConsumers 启动所有消费者
func (mc *MessageConsumer) StartConsumers() {
	go mc.consumeScoreCalculation()
	go mc.consumeStatsUpdate()
	go mc.consumeReportGeneration()
	log.Println("All message consumers started")
}

// consumeScoreCalculation 消费成绩计算消息
func (mc *MessageConsumer) consumeScoreCalculation() {
	msgs, err := mc.channel.Consume(
		"score_calculation", // queue
		"",                  // consumer
		false,               // auto-ack
		false,               // exclusive
		false,               // no-local
		false,               // no-wait
		nil,                 // args
	)
	if err != nil {
		log.Printf("Failed to register score calculation consumer: %v", err)
		return
	}

	log.Println("Score calculation consumer started")

	for msg := range msgs {
		var examResult ExamResultMessage
		err := json.Unmarshal(msg.Body, &examResult)
		if err != nil {
			log.Printf("Failed to unmarshal score calculation message: %v", err)
			msg.Nack(false, false)
			continue
		}

		// 处理成绩计算
		err = mc.processScoreCalculation(examResult)
		if err != nil {
			log.Printf("Failed to process score calculation: %v", err)
			msg.Nack(false, true) // 重新入队
			continue
		}

		msg.Ack(false)
		log.Printf("Score calculation completed for user %d, exam %d", examResult.UserID, examResult.ExamID)
	}
}

// consumeStatsUpdate 消费统计更新消息
func (mc *MessageConsumer) consumeStatsUpdate() {
	msgs, err := mc.channel.Consume(
		"stats_update",
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Printf("Failed to register stats update consumer: %v", err)
		return
	}

	log.Println("Stats update consumer started")

	for msg := range msgs {
		var result ScoreCalculationResult
		err := json.Unmarshal(msg.Body, &result)
		if err != nil {
			log.Printf("Failed to unmarshal stats update message: %v", err)
			msg.Nack(false, false)
			continue
		}

		// 处理统计更新
		err = mc.processStatsUpdate(result)
		if err != nil {
			log.Printf("Failed to process stats update: %v", err)
			msg.Nack(false, true)
			continue
		}

		msg.Ack(false)
		log.Printf("Stats update completed for user %d, exam %d", result.UserID, result.ExamID)
	}
}

// consumeReportGeneration 消费报告生成消息
func (mc *MessageConsumer) consumeReportGeneration() {
	msgs, err := mc.channel.Consume(
		"report_generation",
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Printf("Failed to register report generation consumer: %v", err)
		return
	}

	log.Println("Report generation consumer started")

	for msg := range msgs {
		var result ScoreCalculationResult
		err := json.Unmarshal(msg.Body, &result)
		if err != nil {
			log.Printf("Failed to unmarshal report generation message: %v", err)
			msg.Nack(false, false)
			continue
		}

		// 处理报告生成
		err = mc.processReportGeneration(result)
		if err != nil {
			log.Printf("Failed to process report generation: %v", err)
			msg.Nack(false, true)
			continue
		}

		msg.Ack(false)
		log.Printf("Report generation completed for user %d, exam %d", result.UserID, result.ExamID)
	}
}

// processScoreCalculation 处理成绩计算
func (mc *MessageConsumer) processScoreCalculation(examResult ExamResultMessage) error {
	// 转换TenantID
	tenantID, err := strconv.ParseUint(examResult.TenantID, 10, 32)
	if err != nil {
		return fmt.Errorf("invalid tenant_id: %v", err)
	}

	// 获取考试信息
	var exam models.Exam
	err = mc.db.Where("id = ? AND tenant_id = ?", examResult.ExamID, uint(tenantID)).First(&exam).Error
	if err != nil {
		return fmt.Errorf("failed to get exam: %v", err)
	}

	// 获取试卷信息
	var paper models.Paper
	err = mc.db.Where("id = ? AND tenant_id = ?", exam.PaperID, uint(tenantID)).First(&paper).Error
	if err != nil {
		return fmt.Errorf("failed to get paper: %v", err)
	}

	// 获取试卷题目
	var paperQuestions []models.PaperQuestion
	err = mc.db.Where("paper_id = ?", paper.ID).Find(&paperQuestions).Error
	if err != nil {
		return fmt.Errorf("failed to get paper questions: %v", err)
	}

	// 解析用户答案
	var userAnswers map[string]string
	err = json.Unmarshal([]byte(examResult.Answers), &userAnswers)
	if err != nil {
		return fmt.Errorf("failed to parse user answers: %v", err)
	}

	log.Printf("Processing score calculation for exam %d, user %d", examResult.ExamID, examResult.UserID)
	log.Printf("User answers: %+v", userAnswers)

	// 计算成绩
	correctCount := 0
	totalScore := 0.0
	totalCount := len(paperQuestions)

	for _, pq := range paperQuestions {
		// 获取题目信息
		var question models.Question
		err = mc.db.Where("id = ? AND tenant_id = ?", pq.QuestionID, uint(tenantID)).First(&question).Error
		if err != nil {
			log.Printf("Failed to get question %d: %v", pq.QuestionID, err)
			continue
		}

		// 检查答案
		questionKey := fmt.Sprintf("%d", question.ID)
		userAnswer, exists := userAnswers[questionKey]
		log.Printf("Question %d: key=%s, userAnswer=%s, correctAnswer=%s, exists=%v, match=%v", 
			question.ID, questionKey, userAnswer, question.Answer, exists, userAnswer == question.Answer)
		
		if exists && userAnswer == question.Answer {
			correctCount++
			totalScore += pq.Score
			log.Printf("Question %d correct! Added %.1f points", question.ID, pq.Score)
		}
	}

	log.Printf("Final calculation: correctCount=%d, totalScore=%.1f, totalCount=%d", correctCount, totalScore, totalCount)

	// 更新考试记录
	var examRecord models.ExamRecord
	err = mc.db.Where("exam_id = ? AND user_id = ? AND tenant_id = ?", examResult.ExamID, examResult.UserID, uint(tenantID)).First(&examRecord).Error
	if err != nil {
		return fmt.Errorf("failed to find exam record: %v", err)
	}

	// 只更新分数相关字段，避免覆盖其他已更新的字段
	err = mc.db.Model(&examRecord).Updates(map[string]interface{}{
		"answers":       examResult.Answers,
		"score":         totalScore,
		"correct_count": correctCount,
		"total_count":   totalCount,
		"submit_time":   examResult.SubmitTime,
	}).Error
	if err != nil {
		return fmt.Errorf("failed to update exam record: %v", err)
	}

	// 发布统计更新消息
	result := ScoreCalculationResult{
		ExamID:       examResult.ExamID,
		UserID:       examResult.UserID,
		Score:        totalScore,
		CorrectCount: correctCount,
		TotalCount:   totalCount,
		TenantID:     examResult.TenantID,
	}

	rabbitmq := GetRabbitMQService()
	if rabbitmq != nil {
		rabbitmq.PublishStatsUpdate(result)
		rabbitmq.PublishReportGeneration(result)
	}

	return nil
}

// processStatsUpdate 处理统计更新
func (mc *MessageConsumer) processStatsUpdate(result ScoreCalculationResult) error {
	// 转换TenantID
	tenantID, err := strconv.ParseUint(result.TenantID, 10, 32)
	if err != nil {
		return fmt.Errorf("invalid tenant_id: %v", err)
	}

	// 更新用户统计
	var userStats models.UserStats
	err = mc.db.Where("user_id = ? AND tenant_id = ?", result.UserID, uint(tenantID)).First(&userStats).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 创建新的用户统计
			userStats = models.UserStats{
				UserID:       result.UserID,
				ExamCount:    1,
				TotalScore:   result.Score,
				AverageScore: result.Score,
				TenantID:     uint(tenantID),
			}
			err = mc.db.Create(&userStats).Error
		} else {
			return fmt.Errorf("failed to get user stats: %v", err)
		}
	} else {
		// 更新现有统计
		userStats.ExamCount++
		userStats.TotalScore += result.Score
		userStats.AverageScore = userStats.TotalScore / float64(userStats.ExamCount)
		err = mc.db.Save(&userStats).Error
	}

	if err != nil {
		return fmt.Errorf("failed to update user stats: %v", err)
	}

	// 更新考试统计
	var examStats models.ExamStats
	err = mc.db.Where("exam_id = ? AND tenant_id = ?", result.ExamID, uint(tenantID)).First(&examStats).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 创建新的考试统计
			examStats = models.ExamStats{
				ExamID:           result.ExamID,
				ParticipantCount: 1,
				AverageScore:     result.Score,
				HighestScore:     result.Score,
				LowestScore:      result.Score,
				TenantID:         uint(tenantID),
			}
			err = mc.db.Create(&examStats).Error
		} else {
			return fmt.Errorf("failed to get exam stats: %v", err)
		}
	} else {
		// 更新现有统计
		examStats.ParticipantCount++

		// 重新计算平均分
		var totalScore float64
		mc.db.Model(&models.ExamRecord{}).Where("exam_id = ? AND tenant_id = ?", result.ExamID, uint(tenantID)).Select("COALESCE(SUM(score), 0)").Scan(&totalScore)
		examStats.AverageScore = totalScore / float64(examStats.ParticipantCount)

		// 更新最高分和最低分
		if result.Score > examStats.HighestScore {
			examStats.HighestScore = result.Score
		}
		if result.Score < examStats.LowestScore {
			examStats.LowestScore = result.Score
		}

		err = mc.db.Save(&examStats).Error
	}

	return err
}

// processReportGeneration 处理报告生成
func (mc *MessageConsumer) processReportGeneration(result ScoreCalculationResult) error {
	// 转换TenantID
	tenantID, err := strconv.ParseUint(result.TenantID, 10, 32)
	if err != nil {
		return fmt.Errorf("invalid tenant_id: %v", err)
	}

	// 生成考试报告
	report := models.ExamReport{
		ExamID:       result.ExamID,
		UserID:       result.UserID,
		Score:        result.Score,
		CorrectCount: result.CorrectCount,
		TotalCount:   result.TotalCount,
		PassRate:     float64(result.CorrectCount) / float64(result.TotalCount) * 100,
		GeneratedAt:  time.Now(),
		TenantID:     uint(tenantID),
	}

	// 计算排名
	var rank int64
	mc.db.Model(&models.ExamRecord{}).Where("exam_id = ? AND score > ? AND tenant_id = ?",
		result.ExamID, result.Score, uint(tenantID)).Count(&rank)
	report.Rank = int(rank) + 1

	// 保存报告
	err = mc.db.Create(&report).Error
	if err != nil {
		return fmt.Errorf("failed to save exam report: %v", err)
	}

	log.Printf("Generated exam report for user %d, exam %d, rank: %d",
		result.UserID, result.ExamID, report.Rank)

	return nil
}
