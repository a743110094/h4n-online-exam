package controllers

import (
	"net/http"
	"online-exam-system/database"
	"online-exam-system/middleware"
	"online-exam-system/models"
	"online-exam-system/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type SubmitAnswerRequest struct {
	QuestionID uint   `json:"question_id" binding:"required"`
	Answer     string `json:"answer" binding:"required"`
}

type SubmitExamRequest struct {
	Answers []SubmitAnswerRequest `json:"answers" binding:"required"`
}

type ExamResultResponse struct {
	Record    models.ExamRecord `json:"record"`
	Answers   []AnswerDetail    `json:"answers"`
	Score     int               `json:"score"`
	TotalScore int              `json:"total_score"`
	CorrectCount int            `json:"correct_count"`
	TotalCount   int            `json:"total_count"`
}

type AnswerDetail struct {
	Question      models.Question `json:"question"`
	StudentAnswer string          `json:"student_answer"`
	CorrectAnswer string          `json:"correct_answer"`
	IsCorrect     bool            `json:"is_correct"`
	Score         int             `json:"score"`
}

// 提交单个答案
func SubmitAnswer(c *gin.Context) {
	examID, err := strconv.ParseUint(c.Param("exam_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的考试ID"})
		return
	}

	var req SubmitAnswerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tenantID := middleware.GetTenantID(c)
	currentUserID := middleware.GetCurrentUserID(c)
	currentRole := middleware.GetCurrentUserRole(c)

	// 只有学生可以提交答案
	if currentRole != models.RoleStudent {
		c.JSON(http.StatusForbidden, gin.H{"error": "只有学生可以提交答案"})
		return
	}

	// 检查考试记录是否存在且状态正确
	var record models.ExamRecord
	if err := utils.WithTenant(database.DB, tenantID).Where("exam_id = ? AND student_id = ?", uint(examID), currentUserID).First(&record).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "考试记录不存在"})
		return
	}

	if record.Status != models.ExamInProgress {
		c.JSON(http.StatusBadRequest, gin.H{"error": "考试未进行中"})
		return
	}

	// 检查考试是否已结束
	var exam models.Exam
	if err := utils.WithTenant(database.DB, tenantID).First(&exam, uint(examID)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "考试不存在"})
		return
	}

	if time.Now().After(exam.EndTime) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "考试已结束"})
		return
	}

	// 检查题目是否存在
	var question models.Question
	if err := utils.WithTenant(database.DB, tenantID).First(&question, req.QuestionID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "题目不存在"})
		return
	}

	// 检查或创建答案记录
	var answer models.Answer
	result := database.DB.Where("exam_record_id = ? AND question_id = ?", record.ID, req.QuestionID).First(&answer)

	if result.Error != nil {
		// 创建新答案记录
		answer = models.Answer{
			ExamRecordID: record.ID,
			QuestionID:   req.QuestionID,
			Answer:       req.Answer,
		}
		if err := database.DB.Create(&answer).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "保存答案失败"})
			return
		}
	} else {
		// 更新现有答案记录
		answer.Answer = req.Answer
		if err := database.DB.Save(&answer).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "更新答案失败"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "答案提交成功",
		"answer":  answer,
	})
}

// 提交整份试卷
func SubmitExam(c *gin.Context) {
	examID, err := strconv.ParseUint(c.Param("exam_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的考试ID"})
		return
	}

	var req SubmitExamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tenantID := middleware.GetTenantID(c)
	currentUserID := middleware.GetCurrentUserID(c)
	currentRole := middleware.GetCurrentUserRole(c)

	// 只有学生可以提交试卷
	if currentRole != models.RoleStudent {
		c.JSON(http.StatusForbidden, gin.H{"error": "只有学生可以提交试卷"})
		return
	}

	// 检查考试记录是否存在且状态正确
	var record models.ExamRecord
	if err := utils.WithTenant(database.DB, tenantID).Where("exam_id = ? AND student_id = ?", uint(examID), currentUserID).First(&record).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "考试记录不存在"})
		return
	}

	if record.Status != models.ExamInProgress {
		c.JSON(http.StatusBadRequest, gin.H{"error": "考试未进行中"})
		return
	}

	// 获取考试和试卷信息
	var exam models.Exam
	if err := utils.WithTenant(database.DB, tenantID).Preload("Paper").First(&exam, uint(examID)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "考试不存在"})
		return
	}

	// 保存所有答案
	for _, answerReq := range req.Answers {
		// 检查题目是否存在
		var question models.Question
		if err := utils.WithTenant(database.DB, tenantID).First(&question, answerReq.QuestionID).Error; err != nil {
			continue // 跳过不存在的题目
		}

		// 检查或创建答案记录
		var answer models.Answer
		result := database.DB.Where("exam_record_id = ? AND question_id = ?", record.ID, answerReq.QuestionID).First(&answer)

		if result.Error != nil {
			// 创建新答案记录
			answer = models.Answer{
				ExamRecordID: record.ID,
				QuestionID:   answerReq.QuestionID,
				Answer:       answerReq.Answer,
			}
			database.DB.Create(&answer)
		} else {
			// 更新现有答案记录
			answer.Answer = answerReq.Answer
			database.DB.Save(&answer)
		}
	}

	// 计算成绩
	score, totalScore := calculateScore(record.ID)

	// 更新考试记录
	now := time.Now()
	record.EndTime = &now
	record.Score = &score
	record.Status = models.ExamCompleted

	if err := database.DB.Save(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "提交试卷失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "试卷提交成功",
		"score":       score,
		"total_score": totalScore,
		"record":      record,
	})
}

// 获取考试结果
func GetExamResult(c *gin.Context) {
	examID, err := strconv.ParseUint(c.Param("exam_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的考试ID"})
		return
	}
	tenantID := middleware.GetTenantID(c)
	currentUserID := middleware.GetCurrentUserID(c)
	currentRole := middleware.GetCurrentUserRole(c)

	// 学生只能查看自己的成绩
	var record models.ExamRecord
	if currentRole == models.RoleStudent {
		if err := utils.WithTenant(database.DB, tenantID).Where("exam_id = ? AND student_id = ?", uint(examID), currentUserID).First(&record).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "考试记录不存在"})
			return
		}
	} else {
		// 教师和管理员需要指定学生ID
		studentID := c.Query("student_id")
		if studentID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "请指定学生ID"})
			return
		}

		studentIDUint, err := strconv.ParseUint(studentID, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的学生ID"})
			return
		}

		if err := utils.WithTenant(database.DB, tenantID).Where("exam_id = ? AND student_id = ?", uint(examID), uint(studentIDUint)).First(&record).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "考试记录不存在"})
			return
		}
	}

	if record.Status != models.ExamCompleted {
		c.JSON(http.StatusBadRequest, gin.H{"error": "考试尚未完成"})
		return
	}

	// 获取考试和试卷信息
	var exam models.Exam
	if err := utils.WithTenant(database.DB, tenantID).Preload("Paper").First(&exam, uint(examID)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "考试不存在"})
		return
	}

	// 获取试卷题目
	var questions []models.Question
	if err := database.DB.Preload("Subject").Model(&exam.Paper).Association("Questions").Find(&questions); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取题目失败"})
		return
	}

	// 获取学生答案
	var answers []models.Answer
	if err := database.DB.Where("exam_record_id = ?", record.ID).Find(&answers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取答案失败"})
		return
	}

	// 构建答案映射
	answerMap := make(map[uint]string)
	for _, answer := range answers {
		answerMap[answer.QuestionID] = answer.Answer
	}

	// 构建详细答案列表
	var answerDetails []AnswerDetail
	var correctCount int
	var totalScore int

	for _, question := range questions {
		studentAnswer := answerMap[question.ID]
		isCorrect := checkAnswer(question, studentAnswer)
		score := 0
		if isCorrect {
			score = question.Score
			correctCount++
		}
		totalScore += question.Score

		answerDetails = append(answerDetails, AnswerDetail{
			Question:      question,
			StudentAnswer: studentAnswer,
			CorrectAnswer: question.Answer,
			IsCorrect:     isCorrect,
			Score:         score,
		})
	}

	c.JSON(http.StatusOK, ExamResultResponse{
		Record:       record,
		Answers:      answerDetails,
		Score:        getIntValue(record.Score),
		TotalScore:   totalScore,
		CorrectCount: correctCount,
		TotalCount:   len(questions),
	})
}

// 获取学生答案（考试进行中）
func GetStudentAnswers(c *gin.Context) {
	examID, err := strconv.ParseUint(c.Param("exam_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的考试ID"})
		return
	}
	tenantID := middleware.GetTenantID(c)
	currentUserID := middleware.GetCurrentUserID(c)
	currentRole := middleware.GetCurrentUserRole(c)

	// 只有学生可以获取自己的答案
	if currentRole != models.RoleStudent {
		c.JSON(http.StatusForbidden, gin.H{"error": "只有学生可以获取答案"})
		return
	}

	// 检查考试记录是否存在
	var record models.ExamRecord
	if err := utils.WithTenant(database.DB, tenantID).Where("exam_id = ? AND student_id = ?", uint(examID), currentUserID).First(&record).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "考试记录不存在"})
		return
	}

	// 获取学生答案
	var answers []models.Answer
	if err := database.DB.Preload("Question").Where("exam_record_id = ?", record.ID).Find(&answers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取答案失败"})
		return
	}

	c.JSON(http.StatusOK, answers)
}

// 计算成绩
func calculateScore(examRecordID uint) (int, int) {
	var answers []models.Answer
	database.DB.Preload("Question").Where("exam_record_id = ?", examRecordID).Find(&answers)

	var score, totalScore int
	for _, answer := range answers {
		totalScore += answer.Question.Score
		if checkAnswer(answer.Question, answer.Answer) {
			score += answer.Question.Score
		}
	}

	return score, totalScore
}

// 检查答案是否正确
func checkAnswer(question models.Question, studentAnswer string) bool {
	switch question.Type {
	case models.SingleChoice, models.TrueFalse:
		return question.Answer == studentAnswer
	case models.MultipleChoice:
		// 多选题需要完全匹配
		return question.Answer == studentAnswer
	case models.ShortAnswer:
		// 简答题可以进行简单的字符串比较，实际项目中可能需要更复杂的判断逻辑
		return question.Answer == studentAnswer
	default:
		return false
	}
}

// getIntValue 获取指针类型int的值
func getIntValue(ptr *int) int {
	if ptr != nil {
		return *ptr
	}
	return 0
}