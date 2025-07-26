package services

import (
	"fmt"
	"online-exam-system/cache"
	"online-exam-system/database"
	"online-exam-system/models"
	"online-exam-system/utils"
	"time"
)

const (
	// 缓存键前缀
	ExamCachePrefix     = "exam"
	QuestionCachePrefix = "question"
	PaperCachePrefix    = "paper"
	SubjectCachePrefix  = "subject"
	ExamListCachePrefix = "exam_list"
	UserCachePrefix     = "user"
	TokenCachePrefix    = "token"

	// 缓存过期时间
	ExamCacheTTL     = 30 * time.Minute  // 考试信息缓存30分钟
	QuestionCacheTTL = 60 * time.Minute  // 题目信息缓存1小时
	PaperCacheTTL    = 45 * time.Minute  // 试卷信息缓存45分钟
	SubjectCacheTTL  = 2 * time.Hour     // 科目信息缓存2小时
	ListCacheTTL     = 10 * time.Minute  // 列表缓存10分钟
	UserCacheTTL     = 2 * time.Hour     // 用户信息缓存2小时
	TokenCacheTTL    = 24 * time.Hour    // Token缓存24小时（与JWT过期时间一致）
)

// CacheService 缓存服务
type CacheService struct{}

// NewCacheService 创建缓存服务实例
func NewCacheService() *CacheService {
	return &CacheService{}
}

// GetExamWithCache 从缓存获取考试信息，如果不存在则从数据库获取并缓存
func (cs *CacheService) GetExamWithCache(tenantID uint, examID uint) (*models.Exam, error) {
	cacheKey := fmt.Sprintf("%s:%d", ExamCachePrefix, examID)
	
	// 尝试从缓存获取
	var exam models.Exam
	if err := cache.GetWithTenant(tenantID, cacheKey, &exam); err == nil {
		return &exam, nil
	}

	// 缓存未命中，从数据库获取
	if err := utils.WithTenant(database.DB, tenantID).Preload("Paper").Preload("Paper.Subject").Preload("Creator").First(&exam, examID).Error; err != nil {
		return nil, err
	}

	// 存入缓存
	cache.SetWithTenant(tenantID, cacheKey, exam, ExamCacheTTL)

	return &exam, nil
}

// GetQuestionWithCache 从缓存获取题目信息
func (cs *CacheService) GetQuestionWithCache(tenantID uint, questionID uint) (*models.Question, error) {
	cacheKey := fmt.Sprintf("%s:%d", QuestionCachePrefix, questionID)
	
	// 尝试从缓存获取
	var question models.Question
	if err := cache.GetWithTenant(tenantID, cacheKey, &question); err == nil {
		return &question, nil
	}

	// 缓存未命中，从数据库获取
	if err := utils.WithTenant(database.DB, tenantID).Preload("Subject").Preload("Creator").First(&question, questionID).Error; err != nil {
		return nil, err
	}

	// 存入缓存
	cache.SetWithTenant(tenantID, cacheKey, question, QuestionCacheTTL)

	return &question, nil
}

// GetPaperWithQuestionsCache 从缓存获取试卷及其题目信息
func (cs *CacheService) GetPaperWithQuestionsCache(tenantID uint, paperID uint) (*models.Paper, []models.Question, error) {
	paperCacheKey := fmt.Sprintf("%s:%d", PaperCachePrefix, paperID)
	questionsCacheKey := fmt.Sprintf("%s:%d:questions", PaperCachePrefix, paperID)
	
	// 尝试从缓存获取试卷
	var paper models.Paper
	paperCached := cache.GetWithTenant(tenantID, paperCacheKey, &paper) == nil
	
	// 尝试从缓存获取题目列表
	var questions []models.Question
	questionsCached := cache.GetWithTenant(tenantID, questionsCacheKey, &questions) == nil
	
	// 如果都缓存了，直接返回
	if paperCached && questionsCached {
		return &paper, questions, nil
	}

	// 从数据库获取试卷
	if !paperCached {
		if err := utils.WithTenant(database.DB, tenantID).Preload("Subject").First(&paper, paperID).Error; err != nil {
			return nil, nil, err
		}
		// 缓存试卷信息
		cache.SetWithTenant(tenantID, paperCacheKey, paper, PaperCacheTTL)
	}

	// 从数据库获取题目
	if !questionsCached {
		if err := utils.WithTenant(database.DB, tenantID).Preload("Subject").Model(&paper).Association("Questions").Find(&questions); err != nil {
			return nil, nil, err
		}
		// 缓存题目列表
		cache.SetWithTenant(tenantID, questionsCacheKey, questions, QuestionCacheTTL)
	}

	return &paper, questions, nil
}

// GetSubjectWithCache 从缓存获取科目信息
func (cs *CacheService) GetSubjectWithCache(tenantID uint, subjectID uint) (*models.Subject, error) {
	cacheKey := fmt.Sprintf("%s:%d", SubjectCachePrefix, subjectID)
	
	// 尝试从缓存获取
	var subject models.Subject
	if err := cache.GetWithTenant(tenantID, cacheKey, &subject); err == nil {
		return &subject, nil
	}

	// 缓存未命中，从数据库获取
	if err := utils.WithTenant(database.DB, tenantID).First(&subject, subjectID).Error; err != nil {
		return nil, err
	}

	// 存入缓存
	cache.SetWithTenant(tenantID, cacheKey, subject, SubjectCacheTTL)

	return &subject, nil
}

// GetStudentExamListCache 获取学生考试列表缓存
func (cs *CacheService) GetStudentExamListCache(tenantID uint, studentID uint, status string, page int, size int) (interface{}, bool) {
	cacheKey := fmt.Sprintf("%s:student:%d:status:%s:page:%d:size:%d", ExamListCachePrefix, studentID, status, page, size)
	
	var result interface{}
	if err := cache.GetWithTenant(tenantID, cacheKey, &result); err == nil {
		return result, true
	}
	return nil, false
}

// SetStudentExamListCache 设置学生考试列表缓存
func (cs *CacheService) SetStudentExamListCache(tenantID uint, studentID uint, status string, page int, size int, data interface{}) {
	cacheKey := fmt.Sprintf("%s:student:%d:status:%s:page:%d:size:%d", ExamListCachePrefix, studentID, status, page, size)
	cache.SetWithTenant(tenantID, cacheKey, data, ListCacheTTL)
}

// InvalidateExamCache 使考试相关缓存失效
func (cs *CacheService) InvalidateExamCache(tenantID uint, examID uint) {
	cacheKey := fmt.Sprintf("%s:%d", ExamCachePrefix, examID)
	cache.DeleteWithTenant(tenantID, cacheKey)
	
	// 清除相关的列表缓存
	cs.InvalidateExamListCache(tenantID)
}

// InvalidateQuestionCache 使题目相关缓存失效
func (cs *CacheService) InvalidateQuestionCache(tenantID uint, questionID uint) {
	cacheKey := fmt.Sprintf("%s:%d", QuestionCachePrefix, questionID)
	cache.DeleteWithTenant(tenantID, cacheKey)
}

// InvalidatePaperCache 使试卷相关缓存失效
func (cs *CacheService) InvalidatePaperCache(tenantID uint, paperID uint) {
	paperCacheKey := fmt.Sprintf("%s:%d", PaperCachePrefix, paperID)
	questionsCacheKey := fmt.Sprintf("%s:%d:questions", PaperCachePrefix, paperID)
	
	cache.DeleteWithTenant(tenantID, paperCacheKey)
	cache.DeleteWithTenant(tenantID, questionsCacheKey)
}

// InvalidateExamListCache 使考试列表缓存失效
func (cs *CacheService) InvalidateExamListCache(tenantID uint) {
	// 这里可以使用模式匹配删除所有相关的列表缓存
	// 由于Redis的KEYS命令在生产环境中性能较差，这里简化处理
	// 在实际生产中，可以考虑使用Redis的发布订阅或者标签系统
}

// GetUserWithCache 从缓存获取用户信息
func (cs *CacheService) GetUserWithCache(tenantID uint, userID uint) (*models.User, error) {
	cacheKey := fmt.Sprintf("%s:%d", UserCachePrefix, userID)
	
	// 尝试从缓存获取
	var user models.User
	if err := cache.GetWithTenant(tenantID, cacheKey, &user); err == nil {
		return &user, nil
	}

	// 缓存未命中，从数据库获取
	if err := utils.WithTenant(database.DB, tenantID).First(&user, userID).Error; err != nil {
		return nil, err
	}

	// 清除密码字段后存入缓存
	userForCache := user
	userForCache.Password = ""
	cache.SetWithTenant(tenantID, cacheKey, userForCache, UserCacheTTL)

	// 返回时也清除密码
	user.Password = ""
	return &user, nil
}

// GetUserByUsernameWithCache 通过用户名从缓存获取用户信息
func (cs *CacheService) GetUserByUsernameWithCache(tenantID uint, username string) (*models.User, error) {
	cacheKey := fmt.Sprintf("%s:username:%s", UserCachePrefix, username)
	
	// 尝试从缓存获取
	var user models.User
	if err := cache.GetWithTenant(tenantID, cacheKey, &user); err == nil {
		return &user, nil
	}

	// 缓存未命中，从数据库获取
	if err := utils.WithTenant(database.DB, tenantID).Where("username = ? AND is_active = ?", username, true).First(&user).Error; err != nil {
		return nil, err
	}

	// 清除密码字段后存入缓存
	userForCache := user
	userForCache.Password = ""
	cache.SetWithTenant(tenantID, cacheKey, userForCache, UserCacheTTL)

	// 返回完整用户信息（包含密码用于验证）
	return &user, nil
}

// SetUserCache 设置用户缓存
func (cs *CacheService) SetUserCache(tenantID uint, user *models.User) {
	// 按ID缓存
	cacheKey := fmt.Sprintf("%s:%d", UserCachePrefix, user.ID)
	userForCache := *user
	userForCache.Password = "" // 清除密码
	cache.SetWithTenant(tenantID, cacheKey, userForCache, UserCacheTTL)
	
	// 按用户名缓存
	usernameCacheKey := fmt.Sprintf("%s:username:%s", UserCachePrefix, user.Username)
	cache.SetWithTenant(tenantID, usernameCacheKey, userForCache, UserCacheTTL)
}

// InvalidateUserCache 使用户缓存失效
func (cs *CacheService) InvalidateUserCache(tenantID uint, userID uint, username string) {
	// 删除按ID的缓存
	cacheKey := fmt.Sprintf("%s:%d", UserCachePrefix, userID)
	cache.DeleteWithTenant(tenantID, cacheKey)
	
	// 删除按用户名的缓存
	if username != "" {
		usernameCacheKey := fmt.Sprintf("%s:username:%s", UserCachePrefix, username)
		cache.DeleteWithTenant(tenantID, usernameCacheKey)
	}
}

// SetTokenCache 缓存token信息（用于快速验证）
func (cs *CacheService) SetTokenCache(tenantID uint, tokenHash string, userID uint, username string, role models.UserRole) {
	cacheKey := fmt.Sprintf("%s:%s", TokenCachePrefix, tokenHash)
	tokenInfo := map[string]interface{}{
		"user_id":  userID,
		"username": username,
		"role":     role,
		"tenant_id": tenantID,
	}
	cache.SetWithTenant(tenantID, cacheKey, tokenInfo, TokenCacheTTL)
}

// GetTokenCache 从缓存获取token信息
func (cs *CacheService) GetTokenCache(tenantID uint, tokenHash string) (map[string]interface{}, bool) {
	cacheKey := fmt.Sprintf("%s:%s", TokenCachePrefix, tokenHash)
	
	var tokenInfo map[string]interface{}
	if err := cache.GetWithTenant(tenantID, cacheKey, &tokenInfo); err == nil {
		return tokenInfo, true
	}
	return nil, false
}

// InvalidateTokenCache 使token缓存失效
func (cs *CacheService) InvalidateTokenCache(tenantID uint, tokenHash string) {
	cacheKey := fmt.Sprintf("%s:%s", TokenCachePrefix, tokenHash)
	cache.DeleteWithTenant(tenantID, cacheKey)
}

// WarmupExamCache 预热考试缓存（在考试开始前预加载热点数据）
func (cs *CacheService) WarmupExamCache(tenantID uint, examID uint) error {
	// 预加载考试信息
	exam, err := cs.GetExamWithCache(tenantID, examID)
	if err != nil {
		return err
	}

	// 预加载试卷和题目信息
	_, _, err = cs.GetPaperWithQuestionsCache(tenantID, exam.PaperID)
	if err != nil {
		return err
	}

	return nil
}

// BatchWarmupExamCache 批量预热即将开始的考试缓存
func (cs *CacheService) BatchWarmupExamCache(tenantID uint) error {
	// 获取即将开始的考试（未来1小时内）
	var exams []models.Exam
	now := time.Now()
	oneHourLater := now.Add(time.Hour)
	
	if err := utils.WithTenant(database.DB, tenantID).Where("start_time BETWEEN ? AND ?", now, oneHourLater).Find(&exams).Error; err != nil {
		return err
	}

	// 预热每个考试的缓存
	for _, exam := range exams {
		cs.WarmupExamCache(tenantID, exam.ID)
	}

	return nil
}