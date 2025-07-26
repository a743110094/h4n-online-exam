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
	"time"

	"github.com/gin-gonic/gin"
)

type ExamRequest struct {
	PaperID     uint      `json:"paper_id" binding:"required"`
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"start_time" binding:"required"`
	EndTime     time.Time `json:"end_time" binding:"required"`
	StudentIDs  []uint    `json:"student_ids"` // 指定学生ID列表，为空则所有学生可参加
}

type ExamListResponse struct {
	Exams []models.Exam `json:"exams"`
	Total int64         `json:"total"`
	Page  int           `json:"page"`
	Size  int           `json:"size"`
}

type ExamDetailResponse struct {
	Exam      models.Exam       `json:"exam"`
	Paper     models.Paper      `json:"paper"`
	Questions []models.Question `json:"questions"`
	Record    *models.ExamRecord `json:"record,omitempty"` // 学生的考试记录
}

type StudentExamListResponse struct {
	Exams []StudentExamInfo `json:"exams"`
	Total int64             `json:"total"`
	Page  int               `json:"page"`
	Size  int               `json:"size"`
}

type StudentExamInfo struct {
	Exam   models.Exam       `json:"exam"`
	Paper  models.Paper      `json:"paper"`
	Record *models.ExamRecord `json:"record,omitempty"`
	Status string            `json:"status"` // not_started, in_progress, completed, expired
}

// 获取考试列表（教师/管理员）
func GetExams(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	status := c.Query("status")
	search := c.Query("search")

	// 获取租户ID
	tenantID := middleware.GetTenantID(c)

	offset := (page - 1) * size

	query := utils.WithTenant(database.DB, tenantID).Model(&models.Exam{}).Preload("Paper").Preload("Paper.Subject").Preload("Creator")

	// 如果是教师，只能看到自己创建的考试
	currentRole := middleware.GetCurrentUserRole(c)
	if currentRole == models.RoleTeacher {
		query = query.Where("created_by = ?", middleware.GetCurrentUserID(c))
	}

	// 状态筛选
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 搜索筛选
	if search != "" {
		query = query.Where("title ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// 获取总数
	var total int64
	query.Count(&total)

	// 获取考试列表
	var exams []models.Exam
	if err := query.Offset(offset).Limit(size).Order("created_at DESC").Find(&exams).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取考试列表失败"})
		return
	}

	c.JSON(http.StatusOK, ExamListResponse{
		Exams: exams,
		Total: total,
		Page:  page,
		Size:  size,
	})
}

// 获取学生考试列表
func GetStudentExams(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	status := c.Query("status")

	// 获取租户ID
	tenantID := middleware.GetTenantID(c)
	currentUserID := middleware.GetCurrentUserID(c)

	// 尝试从缓存获取
	cacheService := services.NewCacheService()
	if cachedResult, found := cacheService.GetStudentExamListCache(tenantID, currentUserID, status, page, size); found {
		c.JSON(http.StatusOK, cachedResult)
		return
	}

	offset := (page - 1) * size

	// 获取学生可参加的考试
	query := utils.WithTenant(database.DB, tenantID).Model(&models.Exam{}).Preload("Paper").Preload("Paper.Subject")

	// 状态筛选
	if status != "" {
		switch status {
		case "not_started":
			query = query.Where("start_time > ?", time.Now())
		case "in_progress":
			query = query.Where("start_time <= ? AND end_time > ?", time.Now(), time.Now())
		case "expired":
			query = query.Where("end_time <= ?", time.Now())
		}
	}

	// 获取总数
	var total int64
	query.Count(&total)

	// 获取考试列表
	var exams []models.Exam
	if err := query.Offset(offset).Limit(size).Order("start_time DESC").Find(&exams).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取考试列表失败"})
		return
	}

	// 获取学生的考试记录
	var examIDs []uint
	for _, exam := range exams {
		examIDs = append(examIDs, exam.ID)
	}

	var records []models.ExamRecord
	if len(examIDs) > 0 {
		utils.WithTenant(database.DB, tenantID).Where("exam_id IN ? AND student_id = ?", examIDs, currentUserID).Find(&records)
	}

	// 构建记录映射
	recordMap := make(map[uint]*models.ExamRecord)
	for i := range records {
		recordMap[records[i].ExamID] = &records[i]
	}

	// 构建响应数据
	var studentExams []StudentExamInfo
	for _, exam := range exams {
		record := recordMap[exam.ID]
		status := getExamStatus(exam, record)

		// 检查学生是否有权限参加此考试
		if !canStudentTakeExam(exam, currentUserID) {
			continue
		}

		studentExams = append(studentExams, StudentExamInfo{
			Exam:   exam,
			Paper:  exam.Paper,
			Record: record,
			Status: status,
		})
	}

	response := StudentExamListResponse{
		Exams: studentExams,
		Total: int64(len(studentExams)),
		Page:  page,
		Size:  size,
	}
	
	// 缓存结果
	cacheService.SetStudentExamListCache(tenantID, currentUserID, status, page, size, response)
	
	c.JSON(http.StatusOK, response)
}

// 获取单个考试详情
func GetExam(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的考试ID"})
		return
	}

	// 获取租户ID
	tenantID := middleware.GetTenantID(c)

	// 使用缓存服务获取考试信息
	cacheService := services.NewCacheService()
	exam, err := cacheService.GetExamWithCache(tenantID, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "考试不存在"})
		return
	}

	// 使用缓存服务获取试卷和题目信息
	paper, questions, err := cacheService.GetPaperWithQuestionsCache(tenantID, exam.PaperID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取试卷题目失败"})
		return
	}
	
	// 更新exam中的Paper信息
	exam.Paper = *paper

	// 如果是学生查看，隐藏答案和解析
	if middleware.GetCurrentUserRole(c) == models.RoleStudent {
		for i := range questions {
			questions[i].Answer = ""
			questions[i].Explanation = ""
		}
	}

	// 如果是学生，获取考试记录
	var record *models.ExamRecord
	if middleware.GetCurrentUserRole(c) == models.RoleStudent {
		currentUserID := middleware.GetCurrentUserID(c)
		
		// 检查学生是否有权限参加此考试
		if !canStudentTakeExam(*exam, currentUserID) {
			c.JSON(http.StatusForbidden, gin.H{"error": "没有权限参加此考试"})
			return
		}

		var examRecord models.ExamRecord
		if err := utils.WithTenant(database.DB, tenantID).Where("exam_id = ? AND student_id = ?", uint(id), currentUserID).First(&examRecord).Error; err == nil {
			record = &examRecord
		}
	}

	c.JSON(http.StatusOK, ExamDetailResponse{
		Exam:      *exam,
		Paper:     exam.Paper,
		Questions: questions,
		Record:    record,
	})
}

// 创建考试
func CreateExam(c *gin.Context) {
	var req ExamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取租户ID
	tenantID := middleware.GetTenantID(c)

	// 验证时间
	if req.StartTime.After(req.EndTime) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "开始时间不能晚于结束时间"})
		return
	}

	// 验证试卷是否存在
	var paper models.Paper
	if err := utils.WithTenant(database.DB, tenantID).First(&paper, req.PaperID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "试卷不存在"})
		return
	}

	// 验证学生ID（如果指定了）
	if len(req.StudentIDs) > 0 {
		var studentCount int64
		utils.WithTenant(database.DB, tenantID).Model(&models.User{}).Where("id IN ? AND role = ?", req.StudentIDs, models.RoleStudent).Count(&studentCount)
		if int(studentCount) != len(req.StudentIDs) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "部分学生ID无效"})
			return
		}
	}

	// 将学生ID列表转换为JSON字符串
	studentIDsJSON := ""
	if len(req.StudentIDs) > 0 {
		studentIDsBytes, _ := json.Marshal(req.StudentIDs)
		studentIDsJSON = string(studentIDsBytes)
	}

	// 创建考试
	exam := models.Exam{
		PaperID:     req.PaperID,
		Title:       req.Title,
		Description: req.Description,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		Status:      models.ExamDraft,
		StudentIDs:  studentIDsJSON,
		CreatedBy:   middleware.GetCurrentUserID(c),
	}

	// 设置租户ID
	utils.SetTenantID(&exam, tenantID)

	if err := database.DB.Create(&exam).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建考试失败"})
		return
	}

	// 预加载关联数据
	utils.WithTenant(database.DB, tenantID).Preload("Paper").Preload("Paper.Subject").Preload("Creator").First(&exam, exam.ID)

	// 清除相关缓存
	cacheService := services.NewCacheService()
	cacheService.InvalidateExamListCache(tenantID)

	c.JSON(http.StatusCreated, exam)
}

// 更新考试
func UpdateExam(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的考试ID"})
		return
	}

	// 获取租户ID
	tenantID := middleware.GetTenantID(c)

	var req ExamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var exam models.Exam
	if err := utils.WithTenant(database.DB, tenantID).First(&exam, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "考试不存在"})
		return
	}

	// 检查权限（只有创建者或管理员可以修改）
	currentUserID := middleware.GetCurrentUserID(c)
	currentRole := middleware.GetCurrentUserRole(c)
	if exam.CreatedBy != currentUserID && currentRole != models.RoleAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "没有权限修改此考试"})
		return
	}

	// 检查考试是否已开始
	if exam.Status != models.ExamDraft {
		c.JSON(http.StatusBadRequest, gin.H{"error": "考试已开始，无法修改"})
		return
	}

	// 验证时间
	if req.StartTime.After(req.EndTime) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "开始时间不能晚于结束时间"})
		return
	}

	// 验证试卷是否存在
	var paper models.Paper
	if err := utils.WithTenant(database.DB, tenantID).First(&paper, req.PaperID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "试卷不存在"})
		return
	}

	// 验证学生ID（如果指定了）
	if len(req.StudentIDs) > 0 {
		var studentCount int64
		utils.WithTenant(database.DB, tenantID).Model(&models.User{}).Where("id IN ? AND role = ?", req.StudentIDs, models.RoleStudent).Count(&studentCount)
		if int(studentCount) != len(req.StudentIDs) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "部分学生ID无效"})
			return
		}
	}

	// 将学生ID列表转换为JSON字符串
	studentIDsJSON := ""
	if len(req.StudentIDs) > 0 {
		studentIDsBytes, _ := json.Marshal(req.StudentIDs)
		studentIDsJSON = string(studentIDsBytes)
	}

	// 更新考试
	exam.PaperID = req.PaperID
	exam.Title = req.Title
	exam.Description = req.Description
	exam.StartTime = req.StartTime
	exam.EndTime = req.EndTime
	exam.StudentIDs = studentIDsJSON

	if err := database.DB.Save(&exam).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新考试失败"})
		return
	}

	// 预加载关联数据
	utils.WithTenant(database.DB, tenantID).Preload("Paper").Preload("Paper.Subject").Preload("Creator").First(&exam, exam.ID)

	// 清除相关缓存
	cacheService := services.NewCacheService()
	cacheService.InvalidateExamCache(tenantID, uint(id))
	cacheService.InvalidateExamListCache(tenantID)

	c.JSON(http.StatusOK, exam)
}

// 删除考试
func DeleteExam(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的考试ID"})
		return
	}

	// 获取租户ID
	tenantID := middleware.GetTenantID(c)

	var exam models.Exam
	if err := utils.WithTenant(database.DB, tenantID).First(&exam, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "考试不存在"})
		return
	}

	// 检查权限（只有创建者或管理员可以删除）
	currentUserID := middleware.GetCurrentUserID(c)
	currentRole := middleware.GetCurrentUserRole(c)
	if exam.CreatedBy != currentUserID && currentRole != models.RoleAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "没有权限删除此考试"})
		return
	}

	// 检查考试是否已开始
	if exam.Status != models.ExamDraft {
		c.JSON(http.StatusBadRequest, gin.H{"error": "考试已开始，无法删除"})
		return
	}

	if err := database.DB.Delete(&exam).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除考试失败"})
		return
	}

	// 清除相关缓存
	cacheService := services.NewCacheService()
	cacheService.InvalidateExamCache(tenantID, uint(id))
	cacheService.InvalidateExamListCache(tenantID)

	c.JSON(http.StatusOK, gin.H{"message": "考试删除成功"})
}

// 开始考试
func StartExam(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的考试ID"})
		return
	}

	// 获取租户ID
	tenantID := middleware.GetTenantID(c)

	currentUserID := middleware.GetCurrentUserID(c)
	currentRole := middleware.GetCurrentUserRole(c)

	// 只有学生可以开始考试
	if currentRole != models.RoleStudent {
		c.JSON(http.StatusForbidden, gin.H{"error": "只有学生可以开始考试"})
		return
	}

	var exam models.Exam
	if err := utils.WithTenant(database.DB, tenantID).Preload("Paper").First(&exam, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "考试不存在"})
		return
	}

	// 检查学生是否有权限参加此考试
	if !canStudentTakeExam(exam, currentUserID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "没有权限参加此考试"})
		return
	}

	// 检查考试时间
	now := time.Now()
	if now.Before(exam.StartTime) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "考试尚未开始"})
		return
	}
	if now.After(exam.EndTime) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "考试已结束"})
		return
	}

	// 检查是否已有考试记录
	var existingRecord models.ExamRecord
	if err := database.DB.Where("exam_id = ? AND student_id = ?", uint(id), currentUserID).First(&existingRecord).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "已经参加过此考试"})
		return
	}

	// 创建考试记录
	record := models.ExamRecord{
		ExamID:    uint(id),
		StudentID: currentUserID,
		StartTime: now,
		Status:    models.ExamInProgress,
	}

	if err := database.DB.Create(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "开始考试失败"})
		return
	}

	// 触发考试数据的按需预热
	warmupService := services.NewWarmupService()
	go warmupService.WarmupExamOnDemand(tenantID, uint(id))

	c.JSON(http.StatusOK, gin.H{
		"message": "考试开始成功",
		"record":  record,
	})
}

// 辅助函数：获取考试状态
func getExamStatus(exam models.Exam, record *models.ExamRecord) string {
	now := time.Now()

	if record != nil {
		if record.Status == models.ExamCompleted {
			return "completed"
		}
		if record.Status == models.ExamInProgress {
			if now.After(exam.EndTime) {
				return "expired"
			}
			return "in_progress"
		}
	}

	if now.Before(exam.StartTime) {
		return "not_started"
	}
	if now.After(exam.EndTime) {
		return "expired"
	}
	return "in_progress"
}

// 辅助函数：检查学生是否可以参加考试
func canStudentTakeExam(exam models.Exam, studentID uint) bool {
	// 如果没有指定学生列表，所有学生都可以参加
	if exam.StudentIDs == "" {
		return true
	}

	var studentIDs []uint
	if err := json.Unmarshal([]byte(exam.StudentIDs), &studentIDs); err != nil {
		return false
	}

	for _, id := range studentIDs {
		if id == studentID {
			return true
		}
	}

	return false
}