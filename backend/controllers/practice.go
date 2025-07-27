package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"online-exam-system/database"
	"online-exam-system/middleware"
	"online-exam-system/models"
	"online-exam-system/utils"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// 获取推荐练习列表
func GetPracticeRecommendations(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未登录"})
		return
	}

	// 获取租户ID
	tenantID := middleware.GetTenantID(c)

	// 获取用户信息
	var user models.User
	if err := utils.WithTenant(database.DB, tenantID).First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 获取推荐练习列表
	var recommendations []models.PracticeRecommendation
	if err := utils.WithTenant(database.DB, tenantID).Preload("Subject").Where("is_active = ?", true).Find(&recommendations).Error; err != nil {
		utils.InternalServerErrorResponse(c, "获取推荐练习失败")
		return
	}

	// 计算每个推荐练习的完成进度
	type RecommendationWithProgress struct {
		models.PracticeRecommendation
		Progress float64 `json:"progress"`
	}

	var result []RecommendationWithProgress
	// 确保返回空数组而不是null
	if recommendations == nil {
		result = []RecommendationWithProgress{}
	} else {
		for _, rec := range recommendations {
				// 计算用户在该科目的练习进度
				var totalPractices int64
				var completedPractices int64
				utils.WithTenant(database.DB, tenantID).Model(&models.PracticeRecord{}).Where("user_id = ? AND subject_id = ?", userID, rec.SubjectID).Count(&totalPractices)
				utils.WithTenant(database.DB, tenantID).Model(&models.PracticeRecord{}).Where("user_id = ? AND subject_id = ? AND is_completed = ?", userID, rec.SubjectID, true).Count(&completedPractices)

			progress := 0.0
			if totalPractices > 0 {
				progress = math.Round(float64(completedPractices)/float64(totalPractices)*100*100) / 100
			}

			result = append(result, RecommendationWithProgress{
				PracticeRecommendation: rec,
				Progress:               progress,
			})
		}
	}

	utils.SuccessResponse(c, result)
}

// 开始练习
func StartPractice(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		utils.UnauthorizedResponse(c, "用户未登录")
		return
	}

	// 获取租户ID
	tenantID := middleware.GetTenantID(c)

	type StartPracticeRequest struct {
		SubjectID    uint   `json:"subject_id" binding:"required"`
		Difficulty   int    `json:"difficulty"`
		QuestionType string `json:"question_type"`
		PracticeType string `json:"practice_type"`
		QuestionCount int   `json:"question_count"`
	}

	var req StartPracticeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	// 设置默认值
	if req.PracticeType == "" {
		req.PracticeType = "sequence"
	}
	if req.QuestionCount == 0 {
		req.QuestionCount = 10
	}

	// 构建查询条件
	query := utils.WithTenant(database.DB, tenantID).Where("subject_id = ? AND status = ?", req.SubjectID, models.QuestionPublished)
	if req.Difficulty > 0 {
		query = query.Where("difficulty = ?", req.Difficulty)
	}
	if req.QuestionType != "" {
		query = query.Where("type = ?", req.QuestionType)
	}

	// 如果是错题重练，只获取用户答错的题目
	if req.PracticeType == "wrong" {
		// 获取用户答错的题目ID
		var wrongQuestionIDs []uint
		utils.WithTenant(database.DB, tenantID).Table("practice_answers").
			Joins("JOIN practice_records ON practice_answers.practice_record_id = practice_records.id").
			Where("practice_records.user_id = ? AND practice_answers.is_correct = ? AND practice_records.tenant_id = ? AND practice_answers.tenant_id = ?", userID, false, tenantID, tenantID).
			Pluck("practice_answers.question_id", &wrongQuestionIDs)

		if len(wrongQuestionIDs) == 0 {
			utils.BadRequestResponse(c, "暂无错题可供练习")
			return
		}
		query = query.Where("id IN ?", wrongQuestionIDs)
	}

	// 获取题目列表
	var questions []models.Question
	if req.PracticeType == "random" {
		query = query.Order("RANDOM()")
	} else {
		query = query.Order("id ASC")
	}
	query = query.Limit(req.QuestionCount)

	if err := query.Find(&questions).Error; err != nil {
		utils.InternalServerErrorResponse(c, "获取题目失败")
		return
	}

	if len(questions) == 0 {
		utils.BadRequestResponse(c, "没有找到符合条件的题目")
		return
	}

	// 创建练习记录
	questionIDs := make([]uint, len(questions))
	for i, q := range questions {
		questionIDs[i] = q.ID
	}
	questionIDsJSON, _ := json.Marshal(questionIDs)

	practiceRecord := models.PracticeRecord{
		UserID:       userID,
		SubjectID:    req.SubjectID,
		Title:        "练习记录",
		QuestionIDs:  string(questionIDsJSON),
		TotalCount:   len(questions),
		Difficulty:   req.Difficulty,
		PracticeType: req.PracticeType,
	}

	// 设置租户ID
	utils.SetTenantID(&practiceRecord, tenantID)

	if err := database.DB.Create(&practiceRecord).Error; err != nil {
		utils.InternalServerErrorResponse(c, "创建练习记录失败")
		return
	}

	utils.SuccessResponse(c, gin.H{
		"practice_id": practiceRecord.ID,
		"questions":   questions,
	})
}

// 提交练习答案
func SubmitPracticeAnswer(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		utils.UnauthorizedResponse(c, "用户未登录")
		return
	}

	// 获取租户ID
	tenantID := middleware.GetTenantID(c)

	practiceID, err := strconv.ParseUint(c.Param("practice_id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "无效的练习ID")
		return
	}

	type SubmitAnswerRequest struct {
		QuestionID uint   `json:"question_id" binding:"required"`
		Answer     string `json:"answer" binding:"required"`
		TimeSpent  int    `json:"time_spent"`
	}

	var req SubmitAnswerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	// 验证练习记录是否存在且属于当前用户
	var practiceRecord models.PracticeRecord
	if err := utils.WithTenant(database.DB, tenantID).Where("id = ? AND user_id = ?", practiceID, userID).First(&practiceRecord).Error; err != nil {
		utils.NotFoundResponse(c, "练习记录不存在")
		return
	}

	// 获取题目信息
	var question models.Question
	if err := utils.WithTenant(database.DB, tenantID).First(&question, req.QuestionID).Error; err != nil {
		utils.NotFoundResponse(c, "题目不存在")
		return
	}

	// 判断答案是否正确
	isCorrect := checkPracticeAnswer(question, req.Answer)
	score := 0
	if isCorrect {
		score = question.Score
	}

	// 创建或更新答题记录
	var practiceAnswer models.PracticeAnswer
	err = utils.WithTenant(database.DB, tenantID).Where("practice_record_id = ? AND question_id = ?", practiceID, req.QuestionID).First(&practiceAnswer).Error
	if err != nil {
		// 创建新的答题记录
		practiceAnswer = models.PracticeAnswer{
			PracticeRecordID: uint(practiceID),
			QuestionID:       req.QuestionID,
			Answer:           req.Answer,
			IsCorrect:        &isCorrect,
			Score:            &score,
			TimeSpent:        req.TimeSpent,
		}
		// 设置租户ID
		utils.SetTenantID(&practiceAnswer, tenantID)
		if err := database.DB.Create(&practiceAnswer).Error; err != nil {
			utils.InternalServerErrorResponse(c, "保存答案失败")
			return
		}
	} else {
		// 更新已有的答题记录
		practiceAnswer.Answer = req.Answer
		practiceAnswer.IsCorrect = &isCorrect
		practiceAnswer.Score = &score
		practiceAnswer.TimeSpent = req.TimeSpent
		if err := database.DB.Save(&practiceAnswer).Error; err != nil {
			utils.InternalServerErrorResponse(c, "更新答案失败")
			return
		}
	}

	utils.SuccessResponse(c, gin.H{
		"is_correct": isCorrect,
		"score":      score,
		"explanation": question.Explanation,
	})
}

// 完成练习
func CompletePractice(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		utils.UnauthorizedResponse(c, "用户未登录")
		return
	}

	// 获取租户ID
	tenantID := middleware.GetTenantID(c)

	practiceID, err := strconv.ParseUint(c.Param("practice_id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "无效的练习ID")
		return
	}

	// 验证练习记录是否存在且属于当前用户
	var practiceRecord models.PracticeRecord
	if err := utils.WithTenant(database.DB, tenantID).Where("id = ? AND user_id = ?", practiceID, userID).First(&practiceRecord).Error; err != nil {
		utils.NotFoundResponse(c, "练习记录不存在")
		return
	}

	// 统计答题情况
	var totalAnswers int64
	var correctAnswers int64
	var totalScore int

	utils.WithTenant(database.DB, tenantID).Model(&models.PracticeAnswer{}).Where("practice_record_id = ?", practiceID).Count(&totalAnswers)
	utils.WithTenant(database.DB, tenantID).Model(&models.PracticeAnswer{}).Where("practice_record_id = ? AND is_correct = ?", practiceID, true).Count(&correctAnswers)
	utils.WithTenant(database.DB, tenantID).Model(&models.PracticeAnswer{}).Where("practice_record_id = ?", practiceID).Select("COALESCE(SUM(score), 0)").Scan(&totalScore)

	// 更新练习记录
	practiceRecord.CorrectCount = int(correctAnswers)
	practiceRecord.WrongCount = int(totalAnswers) - int(correctAnswers)
	practiceRecord.Score = totalScore
	practiceRecord.IsCompleted = true

	if err := database.DB.Save(&practiceRecord).Error; err != nil {
		utils.InternalServerErrorResponse(c, "更新练习记录失败")
		return
	}

	accuracy := 0.0
	if totalAnswers > 0 {
		accuracy = math.Round(float64(correctAnswers)/float64(totalAnswers)*100*100) / 100
	}

	utils.SuccessResponse(c, gin.H{
		"total_questions": totalAnswers,
		"correct_count":   correctAnswers,
		"wrong_count":     int(totalAnswers) - int(correctAnswers),
		"score":           totalScore,
		"accuracy":        accuracy,
	})
}

// 获取练习历史
func GetPracticeHistory(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		utils.UnauthorizedResponse(c, "用户未登录")
		return
	}

	// 获取租户ID
	tenantID := middleware.GetTenantID(c)

	// 分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	// 获取练习历史
	var practiceRecords []models.PracticeRecord
	var total int64

	// 先获取总数
	utils.WithTenant(database.DB, tenantID).Model(&models.PracticeRecord{}).Where("user_id = ? AND is_completed = ?", userID, true).Count(&total)
	
	// 再获取分页数据
	utils.WithTenant(database.DB, tenantID).Preload("Subject").Where("user_id = ? AND is_completed = ?", userID, true).Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&practiceRecords)

	// 使用分页响应结构
	utils.SuccessPaginationResponse(c, practiceRecords, total, page, pageSize)
}

// 获取错题列表
func GetWrongQuestions(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		utils.UnauthorizedResponse(c, "用户未登录")
		return
	}

	// 获取租户ID
	tenantID := middleware.GetTenantID(c)

	// 分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	subjectID := c.Query("subject_id")
	practiceRecordID := c.Query("practice_record_id")
	offset := (page - 1) * pageSize

	// 构建查询条件
	query := database.DB.Table("practice_answers").
		Select("practice_answers.*, questions.title, questions.content, questions.options, questions.answer as correct_answer, questions.explanation, questions.type, questions.difficulty, subjects.name as subject_name").
		Joins("JOIN practice_records ON practice_answers.practice_record_id = practice_records.id").
		Joins("JOIN questions ON practice_answers.question_id = questions.id").
		Joins("JOIN subjects ON questions.subject_id = subjects.id").
		Where("practice_records.user_id = ? AND practice_answers.is_correct = ? AND practice_records.tenant_id = ? AND practice_answers.tenant_id = ? AND questions.tenant_id = ? AND subjects.tenant_id = ?", userID, false, tenantID, tenantID, tenantID, tenantID)

	if subjectID != "" {
		query = query.Where("questions.subject_id = ?", subjectID)
	}

	// 如果指定了练习记录ID，只获取该次练习的错题
	if practiceRecordID != "" {
		query = query.Where("practice_answers.practice_record_id = ?", practiceRecordID)
	}

	// 获取总数
	var total int64
	query.Count(&total)

	// 获取分页数据
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

	var wrongQuestions []WrongQuestionDetail
	query.Select("practice_answers.id, practice_answers.question_id, questions.title, questions.content, questions.options, questions.type, questions.difficulty, subjects.name as subject_name, practice_answers.answer as user_answer, questions.answer as correct_answer, questions.explanation, practice_answers.time_spent, practice_answers.created_at").
		Order("practice_answers.created_at DESC").
		Offset(offset).Limit(pageSize).
		Scan(&wrongQuestions)

	utils.SuccessPaginationResponse(c, wrongQuestions, total, page, pageSize)
}

// 开始错题复习
func StartWrongQuestionReview(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		utils.UnauthorizedResponse(c, "用户未登录")
		return
	}

	// 获取租户ID
	tenantID := middleware.GetTenantID(c)

	type StartReviewRequest struct {
		SubjectID        uint   `json:"subject_id"`
		PracticeRecordID uint   `json:"practice_record_id"` // 指定练习记录ID
		QuestionIDs      []uint `json:"question_ids"`      // 可选，指定要复习的错题ID
		MaxQuestions     int    `json:"max_questions"`     // 最大题目数量，默认20
	}

	var req StartReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	// 设置默认值
	if req.MaxQuestions == 0 {
		req.MaxQuestions = 20
	}

	var questionIDs []uint

	if len(req.QuestionIDs) > 0 {
		// 使用指定的错题ID
		questionIDs = req.QuestionIDs
	} else {
		// 获取用户的错题
		query := database.DB.Table("practice_answers").
			Select("DISTINCT practice_answers.question_id").
			Joins("JOIN practice_records ON practice_answers.practice_record_id = practice_records.id").
			Joins("JOIN questions ON practice_answers.question_id = questions.id").
			Where("practice_records.user_id = ? AND practice_answers.is_correct = ? AND practice_records.tenant_id = ? AND practice_answers.tenant_id = ? AND questions.tenant_id = ?", userID, false, tenantID, tenantID, tenantID)

		if req.SubjectID > 0 {
			query = query.Where("questions.subject_id = ?", req.SubjectID)
		}

		// 如果指定了练习记录ID，只获取该次练习的错题
		if req.PracticeRecordID > 0 {
			query = query.Where("practice_answers.practice_record_id = ?", req.PracticeRecordID)
		}

		query.Order("practice_answers.created_at DESC").Limit(req.MaxQuestions).Pluck("question_id", &questionIDs)
	}

	if len(questionIDs) == 0 {
		utils.BadRequestResponse(c, "没有找到错题")
		return
	}

	// 获取题目详情
	var questions []models.Question
	if err := utils.WithTenant(database.DB, tenantID).Preload("Subject").Where("id IN ?", questionIDs).Find(&questions).Error; err != nil {
		utils.InternalServerErrorResponse(c, "获取题目失败")
		return
	}

	// 创建复习练习记录
	questionIDsJSON, _ := json.Marshal(questionIDs)
	practiceRecord := models.PracticeRecord{
		UserID:       userID,
		SubjectID:    req.SubjectID,
		Title:        "错题复习",
		QuestionIDs:  string(questionIDsJSON),
		TotalCount:   len(questions),
		Difficulty:   2, // 默认中等难度
		PracticeType: "review", // 标记为复习类型
	}

	// 设置租户ID
	utils.SetTenantID(&practiceRecord, tenantID)

	if err := database.DB.Create(&practiceRecord).Error; err != nil {
		utils.InternalServerErrorResponse(c, "创建复习记录失败")
		return
	}

	utils.SuccessResponse(c, gin.H{
		"practice_id": practiceRecord.ID,
		"questions":   questions,
		"type":        "review",
	})
}

// checkPracticeAnswer 检查练习答案是否正确
func checkPracticeAnswer(question models.Question, userAnswer string) bool {
	correctAnswer := strings.TrimSpace(question.Answer)
	userAnswer = strings.TrimSpace(userAnswer)

	// 添加调试日志
	log.Printf("[DEBUG] Question ID: %d, Type: %s", question.ID, question.Type)
	log.Printf("[DEBUG] Correct Answer: '%s'", correctAnswer)
	log.Printf("[DEBUG] User Answer: '%s'", userAnswer)

	// 如果是多选题，需要特殊处理
	if question.Type == "multiple" {
		var correctOptions []string
		var userOptions []string

		// 解析正确答案（可能是JSON数组格式如[0,1]或逗号分隔格式如A,B）
		if strings.HasPrefix(correctAnswer, "[") && strings.HasSuffix(correctAnswer, "]") {
			// JSON数组格式，先解析为interface{}数组
			var rawOptions []interface{}
			if err := json.Unmarshal([]byte(correctAnswer), &rawOptions); err != nil {
				log.Printf("[DEBUG] Failed to parse correct answer as JSON: %v", err)
				return false
			}
			// 转换为字符串数组
			for _, option := range rawOptions {
				correctOptions = append(correctOptions, fmt.Sprintf("%v", option))
			}
		} else {
			// 逗号分隔格式
			correctOptions = strings.Split(correctAnswer, ",")
			for i := range correctOptions {
				correctOptions[i] = strings.TrimSpace(correctOptions[i])
			}
		}

		// 解析用户答案（通常是逗号分隔格式如A,B）
		userOptions = strings.Split(userAnswer, ",")
		for i := range userOptions {
			userOptions[i] = strings.TrimSpace(userOptions[i])
		}

		// 将数字索引转换为字母选项（如果正确答案是数字格式）
		if len(correctOptions) > 0 {
			// 检查第一个元素是否为数字
			if _, err := strconv.Atoi(correctOptions[0]); err == nil {
				// 正确答案是数字索引，转换为字母
				for i, option := range correctOptions {
					if idx, err := strconv.Atoi(option); err == nil && idx >= 0 && idx < 26 {
						correctOptions[i] = string(rune('A' + idx))
					}
				}
			}
		}

		sort.Strings(correctOptions)
		sort.Strings(userOptions)

		// 添加调试日志
		log.Printf("[DEBUG] Correct Options after processing: %v", correctOptions)
		log.Printf("[DEBUG] User Options after processing: %v", userOptions)

		// 比较排序后的选项
		if len(correctOptions) != len(userOptions) {
			log.Printf("[DEBUG] Length mismatch: correct=%d, user=%d", len(correctOptions), len(userOptions))
			return false
		}

		for i := range correctOptions {
			if correctOptions[i] != userOptions[i] {
				log.Printf("[DEBUG] Option mismatch at index %d: correct='%s', user='%s'", i, correctOptions[i], userOptions[i])
				return false
			}
		}
		log.Printf("[DEBUG] Multiple choice answer is correct")
		return true
	}

	// 其他题型直接比较
	result := strings.EqualFold(correctAnswer, userAnswer)
	log.Printf("[DEBUG] Other question type result: %t", result)
	return result
}

// 获取练习统计
func GetPracticeStats(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		utils.UnauthorizedResponse(c, "用户未登录")
		return
	}

	// 获取租户ID
	tenantID := middleware.GetTenantID(c)

	// 总练习次数
	var totalPracticed int64
	utils.WithTenant(database.DB, tenantID).Model(&models.PracticeRecord{}).Where("user_id = ? AND is_completed = ?", userID, true).Count(&totalPracticed)

	// 今日练习次数
	today := time.Now().Format("2006-01-02")
	var todayPracticed int64
	utils.WithTenant(database.DB, tenantID).Model(&models.PracticeRecord{}).Where("user_id = ? AND is_completed = ? AND DATE(created_at) = ?", userID, true, today).Count(&todayPracticed)

	// 总正确率
	var totalCorrect, totalQuestions int64
	utils.WithTenant(database.DB, tenantID).Model(&models.PracticeRecord{}).Where("user_id = ? AND is_completed = ?", userID, true).Select("SUM(correct_count)").Scan(&totalCorrect)
	utils.WithTenant(database.DB, tenantID).Model(&models.PracticeRecord{}).Where("user_id = ? AND is_completed = ?", userID, true).Select("SUM(total_count)").Scan(&totalQuestions)

	correctRate := 0.0
	if totalQuestions > 0 {
		correctRate = math.Round(float64(totalCorrect)/float64(totalQuestions)*100*100) / 100
	}

	// 错题数量
	var wrongQuestions int64
	utils.WithTenant(database.DB, tenantID).Model(&models.PracticeAnswer{}).
		Joins("JOIN practice_records ON practice_answers.practice_record_id = practice_records.id").
		Where("practice_records.user_id = ? AND practice_answers.is_correct = ? AND practice_records.tenant_id = ? AND practice_answers.tenant_id = ?", userID, false, tenantID, tenantID).
		Count(&wrongQuestions)

	utils.SuccessResponse(c, gin.H{
		"total_practiced": totalPracticed,
		"today_practiced": todayPracticed,
		"correct_rate":    correctRate,
		"wrong_questions": wrongQuestions,
	})
}