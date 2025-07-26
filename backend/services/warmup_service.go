package services

import (
	"log"
	"online-exam-system/database"
	"online-exam-system/models"
	"online-exam-system/utils"
	"time"
)

// WarmupService 缓存预热服务
type WarmupService struct {
	cacheService *CacheService
}

// NewWarmupService 创建缓存预热服务实例
func NewWarmupService() *WarmupService {
	return &WarmupService{
		cacheService: NewCacheService(),
	}
}

// StartWarmupScheduler 启动缓存预热调度器
func (ws *WarmupService) StartWarmupScheduler() {
	// 每10分钟执行一次预热任务
	ticker := time.NewTicker(10 * time.Minute)
	go func() {
		for {
			select {
			case <-ticker.C:
				ws.WarmupUpcomingExams()
			}
		}
	}()

	log.Println("缓存预热调度器已启动")
}

// WarmupUpcomingExams 预热即将开始的考试数据
func (ws *WarmupService) WarmupUpcomingExams() {
	log.Println("开始预热即将开始的考试数据...")

	// 获取所有租户
	var tenants []models.Tenant
	if err := database.DB.Find(&tenants).Error; err != nil {
		log.Printf("获取租户列表失败: %v", err)
		return
	}

	// 为每个租户预热数据
	for _, tenant := range tenants {
		ws.warmupTenantExams(tenant.ID)
	}

	log.Println("考试数据预热完成")
}

// warmupTenantExams 预热指定租户的考试数据
func (ws *WarmupService) warmupTenantExams(tenantID uint) {
	now := time.Now()
	// 预热未来30分钟内开始的考试
	startTime := now
	endTime := now.Add(30 * time.Minute)

	// 获取即将开始的考试
	var upcomingExams []models.Exam
	if err := utils.WithTenant(database.DB, tenantID).Where("start_time BETWEEN ? AND ?", startTime, endTime).Find(&upcomingExams).Error; err != nil {
		log.Printf("租户 %d 获取即将开始的考试失败: %v", tenantID, err)
		return
	}

	if len(upcomingExams) == 0 {
		return
	}

	log.Printf("租户 %d 发现 %d 个即将开始的考试，开始预热缓存", tenantID, len(upcomingExams))

	// 预热每个考试的数据
	for _, exam := range upcomingExams {
		ws.warmupExamData(tenantID, exam)
	}
}

// warmupExamData 预热单个考试的数据
func (ws *WarmupService) warmupExamData(tenantID uint, exam models.Exam) {
	log.Printf("预热考试 %d (%s) 的数据", exam.ID, exam.Title)

	// 预热考试基本信息
	if _, err := ws.cacheService.GetExamWithCache(tenantID, exam.ID); err != nil {
		log.Printf("预热考试 %d 基本信息失败: %v", exam.ID, err)
	}

	// 预热试卷和题目信息
	if _, _, err := ws.cacheService.GetPaperWithQuestionsCache(tenantID, exam.PaperID); err != nil {
		log.Printf("预热考试 %d 试卷信息失败: %v", exam.ID, err)
	}

	// 预热科目信息（如果试卷有科目）
	var paper models.Paper
	if err := utils.WithTenant(database.DB, tenantID).First(&paper, exam.PaperID).Error; err == nil {
		if paper.SubjectID != 0 {
			if _, err := ws.cacheService.GetSubjectWithCache(tenantID, paper.SubjectID); err != nil {
				log.Printf("预热考试 %d 科目信息失败: %v", exam.ID, err)
			}
		}
	}

	log.Printf("考试 %d 数据预热完成", exam.ID)
}

// WarmupExamOnDemand 按需预热指定考试的数据（考试开始时调用）
func (ws *WarmupService) WarmupExamOnDemand(tenantID uint, examID uint) {
	log.Printf("按需预热考试 %d 的数据", examID)

	// 获取考试信息
	var exam models.Exam
	if err := utils.WithTenant(database.DB, tenantID).First(&exam, examID).Error; err != nil {
		log.Printf("获取考试 %d 信息失败: %v", examID, err)
		return
	}

	// 预热考试数据
	ws.warmupExamData(tenantID, exam)
}

// WarmupPopularQuestions 预热热门题目数据
func (ws *WarmupService) WarmupPopularQuestions(tenantID uint) {
	log.Printf("预热租户 %d 的热门题目数据", tenantID)

	// 获取最近创建的题目（假设这些是热门题目）
	var questions []models.Question
	if err := utils.WithTenant(database.DB, tenantID).Order("created_at DESC").Limit(50).Find(&questions).Error; err != nil {
		log.Printf("获取热门题目失败: %v", err)
		return
	}

	// 预热题目缓存
	for _, question := range questions {
		if _, err := ws.cacheService.GetQuestionWithCache(tenantID, question.ID); err != nil {
			log.Printf("预热题目 %d 失败: %v", question.ID, err)
		}
	}

	log.Printf("租户 %d 热门题目数据预热完成", tenantID)
}

// WarmupAllSubjects 预热所有科目数据
func (ws *WarmupService) WarmupAllSubjects(tenantID uint) {
	log.Printf("预热租户 %d 的所有科目数据", tenantID)

	var subjects []models.Subject
	if err := utils.WithTenant(database.DB, tenantID).Find(&subjects).Error; err != nil {
		log.Printf("获取科目列表失败: %v", err)
		return
	}

	// 预热科目缓存
	for _, subject := range subjects {
		if _, err := ws.cacheService.GetSubjectWithCache(tenantID, subject.ID); err != nil {
			log.Printf("预热科目 %d 失败: %v", subject.ID, err)
		}
	}

	log.Printf("租户 %d 科目数据预热完成", tenantID)
}

// PerformFullWarmup 执行全量数据预热（系统启动时调用）
func (ws *WarmupService) PerformFullWarmup() {
	log.Println("开始执行全量数据预热...")

	// 获取所有租户
	var tenants []models.Tenant
	if err := database.DB.Find(&tenants).Error; err != nil {
		log.Printf("获取租户列表失败: %v", err)
		return
	}

	// 为每个租户预热数据
	for _, tenant := range tenants {
		log.Printf("开始预热租户 %d 的数据", tenant.ID)
		
		// 预热即将开始的考试
		ws.warmupTenantExams(tenant.ID)
		
		// 预热热门题目
		ws.WarmupPopularQuestions(tenant.ID)
		
		// 预热所有科目
		ws.WarmupAllSubjects(tenant.ID)
		
		log.Printf("租户 %d 数据预热完成", tenant.ID)
	}

	log.Println("全量数据预热完成")
}