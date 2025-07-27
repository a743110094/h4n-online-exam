package controllers

import (
	"encoding/json"
	"net/http"
	"online-exam-system/database"
	"online-exam-system/middleware"
	"online-exam-system/models"
	"online-exam-system/services"
	"online-exam-system/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type QuestionRequest struct {
	SubjectID   uint                   `json:"subject_id" binding:"required"`
	Type        models.QuestionType    `json:"type" binding:"required"`
	Title       string                 `json:"title" binding:"required"`
	Content     string                 `json:"content"`
	Options     []string               `json:"options"`
	Answer      string                 `json:"answer" binding:"required"`
	Explanation string                 `json:"explanation"`
	Difficulty  int                    `json:"difficulty"`
	Score       int                    `json:"score"`
	Status      models.QuestionStatus  `json:"status"`
}

type QuestionListResponse struct {
	Questions []models.Question `json:"questions"`
	Total     int64             `json:"total"`
	Page      int               `json:"page"`
	Size      int               `json:"size"`
}

// 获取题目列表
func GetQuestions(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	subjectID := c.Query("subject_id")
	questionType := c.Query("type")
	search := c.Query("search")
	difficulty := c.Query("difficulty")
	tenantID := middleware.GetTenantID(c)

	offset := (page - 1) * size

	query := utils.WithTenant(database.DB, tenantID).Model(&models.Question{}).Preload("Subject").Preload("Creator")

	// 科目筛选
	if subjectID != "" {
		query = query.Where("subject_id = ?", subjectID)
	}

	// 题目类型筛选
	if questionType != "" {
		query = query.Where("type = ?", questionType)
	}

	// 难度筛选
	if difficulty != "" {
		query = query.Where("difficulty = ?", difficulty)
	}

	// 搜索筛选
	if search != "" {
		query = query.Where("title ILIKE ? OR content ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// 获取总数
	var total int64
	query.Count(&total)

	// 获取题目列表
	var questions []models.Question
	if err := query.Offset(offset).Limit(size).Order("created_at DESC").Find(&questions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取题目列表失败"})
		return
	}

	utils.SuccessResponse(c, QuestionListResponse{
		Questions: questions,
		Total:     total,
		Page:      page,
		Size:      size,
	})
}

// 获取单个题目
func GetQuestion(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的题目ID"})
		return
	}
	tenantID := middleware.GetTenantID(c)

	// 使用缓存服务获取题目信息
	cacheService := services.NewCacheService()
	question, err := cacheService.GetQuestionWithCache(tenantID, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "题目不存在"})
		return
	}

	c.JSON(http.StatusOK, question)
}

// 创建题目
func CreateQuestion(c *gin.Context) {
	var req QuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tenantID := middleware.GetTenantID(c)

	// 验证科目是否存在
	var subject models.Subject
	if err := utils.WithTenant(database.DB, tenantID).First(&subject, req.SubjectID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "科目不存在"})
		return
	}

	// 将选项转换为JSON字符串
	optionsJSON, _ := json.Marshal(req.Options)

	// 设置默认状态
	status := req.Status
	if status == "" {
		status = models.QuestionPublished
	}

	// 创建题目
	question := models.Question{
		SubjectID:   req.SubjectID,
		Type:        req.Type,
		Title:       req.Title,
		Content:     req.Content,
		Options:     string(optionsJSON),
		Answer:      req.Answer,
		Explanation: req.Explanation,
		Difficulty:  req.Difficulty,
		Score:       req.Score,
		Status:      status,
		CreatedBy:   middleware.GetCurrentUserID(c),
	}
	utils.SetTenantID(&question, tenantID)

	if err := database.DB.Create(&question).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建题目失败"})
		return
	}

	// 预加载关联数据
	utils.WithTenant(database.DB, tenantID).Preload("Subject").Preload("Creator").First(&question, question.ID)

	// 清除相关缓存
	cacheService := services.NewCacheService()
	cacheService.InvalidatePaperCache(tenantID, 0) // 清除所有试卷缓存，因为题目可能被多个试卷使用

	c.JSON(http.StatusCreated, question)
}

// 更新题目
func UpdateQuestion(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的题目ID"})
		return
	}
	tenantID := middleware.GetTenantID(c)

	var req QuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var question models.Question
	if err := utils.WithTenant(database.DB, tenantID).First(&question, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "题目不存在"})
		return
	}

	// 检查权限（只有创建者或管理员可以修改）
	currentUserID := middleware.GetCurrentUserID(c)
	currentRole := middleware.GetCurrentUserRole(c)
	if question.CreatedBy != currentUserID && currentRole != models.RoleAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "没有权限修改此题目"})
		return
	}

	// 验证科目是否存在
	var subject models.Subject
	if err := utils.WithTenant(database.DB, tenantID).First(&subject, req.SubjectID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "科目不存在"})
		return
	}

	// 将选项转换为JSON字符串
	optionsJSON, _ := json.Marshal(req.Options)

	// 更新题目
	question.SubjectID = req.SubjectID
	question.Type = req.Type
	question.Title = req.Title
	question.Content = req.Content
	question.Options = string(optionsJSON)
	question.Answer = req.Answer
	question.Explanation = req.Explanation
	question.Difficulty = req.Difficulty
	question.Score = req.Score
	if req.Status != "" {
		question.Status = req.Status
	}

	if err := database.DB.Save(&question).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新题目失败"})
		return
	}

	// 预加载关联数据
	utils.WithTenant(database.DB, tenantID).Preload("Subject").Preload("Creator").First(&question, question.ID)

	// 清除相关缓存
	cacheService := services.NewCacheService()
	cacheService.InvalidateQuestionCache(tenantID, uint(id))
	cacheService.InvalidatePaperCache(tenantID, 0) // 清除所有试卷缓存

	c.JSON(http.StatusOK, question)
}

// 删除题目
func DeleteQuestion(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的题目ID"})
		return
	}
	tenantID := middleware.GetTenantID(c)

	var question models.Question
	if err := utils.WithTenant(database.DB, tenantID).First(&question, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "题目不存在"})
		return
	}

	// 检查权限（只有创建者或管理员可以删除）
	currentUserID := middleware.GetCurrentUserID(c)
	currentRole := middleware.GetCurrentUserRole(c)
	if question.CreatedBy != currentUserID && currentRole != models.RoleAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "没有权限删除此题目"})
		return
	}

	if err := database.DB.Delete(&question).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除题目失败"})
		return
	}

	// 清除相关缓存
	cacheService := services.NewCacheService()
	cacheService.InvalidateQuestionCache(tenantID, uint(id))
	cacheService.InvalidatePaperCache(tenantID, 0) // 清除所有试卷缓存

	c.JSON(http.StatusOK, gin.H{"message": "题目删除成功"})
}

// 批量导入题目
func BatchImportQuestions(c *gin.Context) {
	var req struct {
		Questions []QuestionRequest `json:"questions" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tenantID := middleware.GetTenantID(c)

	var successCount int
	var errors []string
	currentUserID := middleware.GetCurrentUserID(c)

	for i, questionReq := range req.Questions {
		// 验证科目是否存在
		var subject models.Subject
		if err := utils.WithTenant(database.DB, tenantID).First(&subject, questionReq.SubjectID).Error; err != nil {
			errors = append(errors, "第"+strconv.Itoa(i+1)+"题：科目不存在")
			continue
		}

		// 将选项转换为JSON字符串
		optionsJSON, _ := json.Marshal(questionReq.Options)

		// 创建题目
		question := models.Question{
			SubjectID:   questionReq.SubjectID,
			Type:        questionReq.Type,
			Title:       questionReq.Title,
			Content:     questionReq.Content,
			Options:     string(optionsJSON),
			Answer:      questionReq.Answer,
			Explanation: questionReq.Explanation,
			Difficulty:  questionReq.Difficulty,
			Score:       questionReq.Score,
			Status:      questionReq.Status,
			CreatedBy:   currentUserID,
		}
		utils.SetTenantID(&question, tenantID)

		if err := database.DB.Create(&question).Error; err != nil {
			errors = append(errors, "第"+strconv.Itoa(i+1)+"题：创建失败")
			continue
		}

		successCount++
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "批量导入完成",
		"success_count": successCount,
		"errors":        errors,
	})
}

// 获取题目统计信息
func GetQuestionStats(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)

	var stats struct {
		Total          int64 `json:"total"`
		SingleChoice   int64 `json:"single_choice"`
		MultipleChoice int64 `json:"multiple_choice"`
		TrueFalse      int64 `json:"true_false"`
		ShortAnswer    int64 `json:"short_answer"`
	}

	// 总题目数
	utils.WithTenant(database.DB, tenantID).Model(&models.Question{}).Count(&stats.Total)

	// 各类型题目数
	utils.WithTenant(database.DB, tenantID).Model(&models.Question{}).Where("type = ?", models.SingleChoice).Count(&stats.SingleChoice)
	utils.WithTenant(database.DB, tenantID).Model(&models.Question{}).Where("type = ?", models.MultipleChoice).Count(&stats.MultipleChoice)
	utils.WithTenant(database.DB, tenantID).Model(&models.Question{}).Where("type = ?", models.TrueFalse).Count(&stats.TrueFalse)
	utils.WithTenant(database.DB, tenantID).Model(&models.Question{}).Where("type = ?", models.ShortAnswer).Count(&stats.ShortAnswer)

	c.JSON(http.StatusOK, stats)
}