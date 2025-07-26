package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"online-exam-system/database"
	"online-exam-system/middleware"
	"online-exam-system/models"
	"online-exam-system/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type DashboardStats struct {
	UserStats    UserStats    `json:"user_stats"`
	ExamStats    ExamStats    `json:"exam_stats"`
	QuestionStats QuestionStats `json:"question_stats"`
	RecentActivity []RecentActivity `json:"recent_activity"`
}

type UserStats struct {
	TotalUsers    int64 `json:"total_users"`
	TotalStudents int64 `json:"total_students"`
	TotalTeachers int64 `json:"total_teachers"`
	ActiveUsers   int64 `json:"active_users"` // 最近30天活跃用户
}

type ExamStats struct {
	TotalExams     int64 `json:"total_exams"`
	ActiveExams    int64 `json:"active_exams"`
	CompletedExams int64 `json:"completed_exams"`
	TotalRecords   int64 `json:"total_records"`
}

type QuestionStats struct {
	TotalQuestions    int64 `json:"total_questions"`
	TotalSubjects     int64 `json:"total_subjects"`
	TotalPapers       int64 `json:"total_papers"`
	SingleChoice      int64 `json:"single_choice"`
	MultipleChoice    int64 `json:"multiple_choice"`
	TrueFalse         int64 `json:"true_false"`
	ShortAnswer       int64 `json:"short_answer"`
}

type RecentActivity struct {
	Type        string    `json:"type"`
	Description string    `json:"description"`
	UserName    string    `json:"user_name"`
	CreatedAt   time.Time `json:"created_at"`
}

type StudentStats struct {
	TotalExams      int64   `json:"total_exams"`
	CompletedExams  int64   `json:"completed_exams"`
	AverageScore    float64 `json:"average_score"`
	BestScore       int     `json:"best_score"`
	RecentExams     []StudentExamRecord `json:"recent_exams"`
	SubjectStats    []SubjectStat `json:"subject_stats"`
}

type StudentExamRecord struct {
	Exam      models.Exam       `json:"exam"`
	Record    models.ExamRecord `json:"record"`
	Paper     models.Paper      `json:"paper"`
}

type SubjectStat struct {
	Subject       models.Subject `json:"subject"`
	ExamCount     int64          `json:"exam_count"`
	AverageScore  float64        `json:"average_score"`
	BestScore     int            `json:"best_score"`
}

type TeacherStats struct {
	TotalQuestions int64 `json:"total_questions"`
	TotalPapers    int64 `json:"total_papers"`
	TotalExams     int64 `json:"total_exams"`
	TotalStudents  int64 `json:"total_students"` // 参与过考试的学生数
	RecentExams    []TeacherExamRecord `json:"recent_exams"`
	SubjectStats   []TeacherSubjectStat `json:"subject_stats"`
}

type TeacherExamRecord struct {
	Exam           models.Exam `json:"exam"`
	Participants   int64       `json:"participants"`
	AverageScore   float64     `json:"average_score"`
	CompletionRate float64     `json:"completion_rate"`
}

type TeacherSubjectStat struct {
	Subject        models.Subject `json:"subject"`
	QuestionCount  int64          `json:"question_count"`
	PaperCount     int64          `json:"paper_count"`
	ExamCount      int64          `json:"exam_count"`
}

type ExamAnalysis struct {
	Exam            models.Exam         `json:"exam"`
	Paper           models.Paper        `json:"paper"`
	TotalStudents   int64               `json:"total_students"`
	CompletedCount  int64               `json:"completed_count"`
	AverageScore    float64             `json:"average_score"`
	HighestScore    int                 `json:"highest_score"`
	LowestScore     int                 `json:"lowest_score"`
	CompletionRate  float64             `json:"completion_rate"`
	QuestionAnalysis []QuestionAnalysis `json:"question_analysis"`
	ScoreDistribution []ScoreRange      `json:"score_distribution"`
}

type QuestionAnalysis struct {
	Question      models.Question `json:"question"`
	CorrectCount  int64           `json:"correct_count"`
	TotalCount    int64           `json:"total_count"`
	CorrectRate   float64         `json:"correct_rate"`
}

type ScoreRange struct {
	Range string `json:"range"`
	Count int64  `json:"count"`
}

// 获取仪表板统计数据（管理员）
func GetDashboardStats(c *gin.Context) {
	currentRole := middleware.GetCurrentUserRole(c)
	if currentRole != models.RoleAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "权限不足"})
		return
	}
	tenantID := middleware.GetTenantID(c)

	stats := DashboardStats{}

	// 用户统计
	utils.WithTenant(database.DB, tenantID).Model(&models.User{}).Count(&stats.UserStats.TotalUsers)
	utils.WithTenant(database.DB, tenantID).Model(&models.User{}).Where("role = ?", models.RoleStudent).Count(&stats.UserStats.TotalStudents)
	utils.WithTenant(database.DB, tenantID).Model(&models.User{}).Where("role = ?", models.RoleTeacher).Count(&stats.UserStats.TotalTeachers)
	utils.WithTenant(database.DB, tenantID).Model(&models.User{}).Where("updated_at > ?", time.Now().AddDate(0, 0, -30)).Count(&stats.UserStats.ActiveUsers)

	// 考试统计
	utils.WithTenant(database.DB, tenantID).Model(&models.Exam{}).Count(&stats.ExamStats.TotalExams)
	utils.WithTenant(database.DB, tenantID).Model(&models.Exam{}).Where("start_time <= ? AND end_time > ?", time.Now(), time.Now()).Count(&stats.ExamStats.ActiveExams)
	utils.WithTenant(database.DB, tenantID).Model(&models.Exam{}).Where("end_time <= ?", time.Now()).Count(&stats.ExamStats.CompletedExams)
	utils.WithTenant(database.DB, tenantID).Model(&models.ExamRecord{}).Count(&stats.ExamStats.TotalRecords)

	// 题目统计
	utils.WithTenant(database.DB, tenantID).Model(&models.Question{}).Count(&stats.QuestionStats.TotalQuestions)
	utils.WithTenant(database.DB, tenantID).Model(&models.Subject{}).Count(&stats.QuestionStats.TotalSubjects)
	utils.WithTenant(database.DB, tenantID).Model(&models.Paper{}).Count(&stats.QuestionStats.TotalPapers)
	utils.WithTenant(database.DB, tenantID).Model(&models.Question{}).Where("type = ?", models.SingleChoice).Count(&stats.QuestionStats.SingleChoice)
	utils.WithTenant(database.DB, tenantID).Model(&models.Question{}).Where("type = ?", models.MultipleChoice).Count(&stats.QuestionStats.MultipleChoice)
	utils.WithTenant(database.DB, tenantID).Model(&models.Question{}).Where("type = ?", models.TrueFalse).Count(&stats.QuestionStats.TrueFalse)
	utils.WithTenant(database.DB, tenantID).Model(&models.Question{}).Where("type = ?", models.ShortAnswer).Count(&stats.QuestionStats.ShortAnswer)

	// 最近活动
	stats.RecentActivity = getRecentActivity(tenantID)

	c.JSON(http.StatusOK, stats)
}

// 获取学生统计数据
func GetStudentStats(c *gin.Context) {
	currentUserID := middleware.GetCurrentUserID(c)
	currentRole := middleware.GetCurrentUserRole(c)

	if currentRole != models.RoleStudent {
		c.JSON(http.StatusForbidden, gin.H{"error": "只有学生可以查看此统计"})
		return
	}
	tenantID := middleware.GetTenantID(c)

	stats := StudentStats{}

	// 基本统计
	utils.WithTenant(database.DB, tenantID).Model(&models.ExamRecord{}).Where("student_id = ?", currentUserID).Count(&stats.TotalExams)
	utils.WithTenant(database.DB, tenantID).Model(&models.ExamRecord{}).Where("student_id = ? AND status = ?", currentUserID, models.ExamCompleted).Count(&stats.CompletedExams)

	// 平均分和最高分
	var avgScore sql.NullFloat64
	var bestScore sql.NullInt64
	utils.WithTenant(database.DB, tenantID).Model(&models.ExamRecord{}).Where("student_id = ? AND status = ?", currentUserID, models.ExamCompleted).Select("AVG(score)").Scan(&avgScore)
	utils.WithTenant(database.DB, tenantID).Model(&models.ExamRecord{}).Where("student_id = ? AND status = ?", currentUserID, models.ExamCompleted).Select("MAX(score)").Scan(&bestScore)

	if avgScore.Valid {
		stats.AverageScore = avgScore.Float64
	}
	if bestScore.Valid {
		stats.BestScore = int(bestScore.Int64)
	}

	// 最近考试记录
	var recentRecords []models.ExamRecord
	utils.WithTenant(database.DB, tenantID).Preload("Exam").Preload("Exam.Paper").Preload("Exam.Paper.Subject").Where("student_id = ?", currentUserID).Order("created_at DESC").Limit(5).Find(&recentRecords)

	for _, record := range recentRecords {
		stats.RecentExams = append(stats.RecentExams, StudentExamRecord{
			Exam:   record.Exam,
			Record: record,
			Paper:  record.Exam.Paper,
		})
	}

	// 科目统计
	stats.SubjectStats = getStudentSubjectStats(currentUserID, tenantID)

	c.JSON(http.StatusOK, stats)
}

// 获取教师统计数据
func GetTeacherStats(c *gin.Context) {
	currentUserID := middleware.GetCurrentUserID(c)
	currentRole := middleware.GetCurrentUserRole(c)

	if currentRole != models.RoleTeacher {
		c.JSON(http.StatusForbidden, gin.H{"error": "只有教师可以查看此统计"})
		return
	}
	tenantID := middleware.GetTenantID(c)

	stats := TeacherStats{}

	// 基本统计
	utils.WithTenant(database.DB, tenantID).Model(&models.Question{}).Where("created_by = ?", currentUserID).Count(&stats.TotalQuestions)
	utils.WithTenant(database.DB, tenantID).Model(&models.Paper{}).Where("created_by = ?", currentUserID).Count(&stats.TotalPapers)
	utils.WithTenant(database.DB, tenantID).Model(&models.Exam{}).Where("created_by = ?", currentUserID).Count(&stats.TotalExams)

	// 参与过考试的学生数
	utils.WithTenant(database.DB, tenantID).Model(&models.ExamRecord{}).Joins("JOIN exams ON exam_records.exam_id = exams.id").Where("exams.created_by = ?", currentUserID).Distinct("student_id").Count(&stats.TotalStudents)

	// 最近考试记录
	stats.RecentExams = getTeacherRecentExams(currentUserID, tenantID)

	// 科目统计
	stats.SubjectStats = getTeacherSubjectStats(currentUserID, tenantID)

	c.JSON(http.StatusOK, stats)
}

// 获取考试分析
func GetExamAnalysis(c *gin.Context) {
	examID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的考试ID"})
		return
	}
	tenantID := middleware.GetTenantID(c)
	currentUserID := middleware.GetCurrentUserID(c)
	currentRole := middleware.GetCurrentUserRole(c)

	// 获取考试信息
	var exam models.Exam
	if err := utils.WithTenant(database.DB, tenantID).Preload("Paper").Preload("Paper.Subject").First(&exam, uint(examID)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "考试不存在"})
		return
	}

	// 权限检查
	if currentRole == models.RoleTeacher && exam.CreatedBy != currentUserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "没有权限查看此考试分析"})
		return
	}

	analysis := ExamAnalysis{
		Exam:  exam,
		Paper: exam.Paper,
	}

	// 统计参与学生数
	if exam.StudentIDs != "" {
		// 指定学生考试
		var studentIDs []uint
		json.Unmarshal([]byte(exam.StudentIDs), &studentIDs)
		analysis.TotalStudents = int64(len(studentIDs))
	} else {
		// 所有学生都可参加
		utils.WithTenant(database.DB, tenantID).Model(&models.User{}).Where("role = ?", models.RoleStudent).Count(&analysis.TotalStudents)
	}

	// 完成情况统计
	utils.WithTenant(database.DB, tenantID).Model(&models.ExamRecord{}).Where("exam_id = ? AND status = ?", uint(examID), models.ExamCompleted).Count(&analysis.CompletedCount)

	if analysis.TotalStudents > 0 {
		analysis.CompletionRate = float64(analysis.CompletedCount) / float64(analysis.TotalStudents) * 100
	}

	// 成绩统计
	if analysis.CompletedCount > 0 {
		var avgScore sql.NullFloat64
		var maxScore, minScore sql.NullInt64

		utils.WithTenant(database.DB, tenantID).Model(&models.ExamRecord{}).Where("exam_id = ? AND status = ?", uint(examID), models.ExamCompleted).Select("AVG(score)").Scan(&avgScore)
		utils.WithTenant(database.DB, tenantID).Model(&models.ExamRecord{}).Where("exam_id = ? AND status = ?", uint(examID), models.ExamCompleted).Select("MAX(score)").Scan(&maxScore)
		utils.WithTenant(database.DB, tenantID).Model(&models.ExamRecord{}).Where("exam_id = ? AND status = ?", uint(examID), models.ExamCompleted).Select("MIN(score)").Scan(&minScore)

		if avgScore.Valid {
			analysis.AverageScore = avgScore.Float64
		}
		if maxScore.Valid {
			analysis.HighestScore = int(maxScore.Int64)
		}
		if minScore.Valid {
			analysis.LowestScore = int(minScore.Int64)
		}
	}

	// 题目分析
	analysis.QuestionAnalysis = getQuestionAnalysis(uint(examID), tenantID)

	// 分数分布
	analysis.ScoreDistribution = getScoreDistribution(uint(examID), tenantID)

	c.JSON(http.StatusOK, analysis)
}

// 获取最近活动
func getRecentActivity(tenantID uint) []RecentActivity {
	var activities []RecentActivity

	// 最近创建的考试
	var recentExams []models.Exam
	utils.WithTenant(database.DB, tenantID).Preload("Creator").Order("created_at DESC").Limit(5).Find(&recentExams)
	for _, exam := range recentExams {
		activities = append(activities, RecentActivity{
			Type:        "exam_created",
			Description: "创建了考试: " + exam.Title,
			UserName:    exam.Creator.Username,
			CreatedAt:   exam.CreatedAt,
		})
	}

	// 最近完成的考试
	var recentRecords []models.ExamRecord
	utils.WithTenant(database.DB, tenantID).Preload("Student").Preload("Exam").Where("status = ?", models.ExamCompleted).Order("end_time DESC").Limit(5).Find(&recentRecords)
	for _, record := range recentRecords {
		activities = append(activities, RecentActivity{
			Type:        "exam_completed",
			Description: "完成了考试: " + record.Exam.Title,
			UserName:    record.Student.Username,
			CreatedAt:   *record.EndTime,
		})
	}

	return activities
}

// 获取学生科目统计
func getStudentSubjectStats(studentID uint, tenantID uint) []SubjectStat {
	var stats []SubjectStat

	var subjects []models.Subject
	utils.WithTenant(database.DB, tenantID).Find(&subjects)

	for _, subject := range subjects {
		var examCount int64
		var avgScore sql.NullFloat64
		var bestScore sql.NullInt64

		// 统计该科目的考试数量和成绩
		utils.WithTenant(database.DB, tenantID).Model(&models.ExamRecord{}).
			Joins("JOIN exams ON exam_records.exam_id = exams.id").
			Joins("JOIN papers ON exams.paper_id = papers.id").
			Where("exam_records.student_id = ? AND papers.subject_id = ? AND exam_records.status = ?", studentID, subject.ID, models.ExamCompleted).
			Count(&examCount)

		if examCount > 0 {
			utils.WithTenant(database.DB, tenantID).Model(&models.ExamRecord{}).
				Joins("JOIN exams ON exam_records.exam_id = exams.id").
				Joins("JOIN papers ON exams.paper_id = papers.id").
				Where("exam_records.student_id = ? AND papers.subject_id = ? AND exam_records.status = ?", studentID, subject.ID, models.ExamCompleted).
				Select("AVG(exam_records.score)").Scan(&avgScore)

			utils.WithTenant(database.DB, tenantID).Model(&models.ExamRecord{}).
				Joins("JOIN exams ON exam_records.exam_id = exams.id").
				Joins("JOIN papers ON exams.paper_id = papers.id").
				Where("exam_records.student_id = ? AND papers.subject_id = ? AND exam_records.status = ?", studentID, subject.ID, models.ExamCompleted).
				Select("MAX(exam_records.score)").Scan(&bestScore)

			stat := SubjectStat{
				Subject:   subject,
				ExamCount: examCount,
			}

			if avgScore.Valid {
				stat.AverageScore = avgScore.Float64
			}
			if bestScore.Valid {
				stat.BestScore = int(bestScore.Int64)
			}

			stats = append(stats, stat)
		}
	}

	return stats
}

// 获取教师最近考试记录
func getTeacherRecentExams(teacherID uint, tenantID uint) []TeacherExamRecord {
	var records []TeacherExamRecord

	var exams []models.Exam
	utils.WithTenant(database.DB, tenantID).Where("created_by = ?", teacherID).Order("created_at DESC").Limit(5).Find(&exams)

	for _, exam := range exams {
		var participants int64
		var avgScore sql.NullFloat64
		var totalStudents int64

		utils.WithTenant(database.DB, tenantID).Model(&models.ExamRecord{}).Where("exam_id = ?", exam.ID).Count(&participants)
		utils.WithTenant(database.DB, tenantID).Model(&models.ExamRecord{}).Where("exam_id = ? AND status = ?", exam.ID, models.ExamCompleted).Select("AVG(score)").Scan(&avgScore)

		// 计算完成率
		if exam.StudentIDs != "" {
			var studentIDs []uint
			json.Unmarshal([]byte(exam.StudentIDs), &studentIDs)
			totalStudents = int64(len(studentIDs))
		} else {
			utils.WithTenant(database.DB, tenantID).Model(&models.User{}).Where("role = ?", models.RoleStudent).Count(&totalStudents)
		}

		record := TeacherExamRecord{
			Exam:         exam,
			Participants: participants,
		}

		if avgScore.Valid {
			record.AverageScore = avgScore.Float64
		}

		if totalStudents > 0 {
			record.CompletionRate = float64(participants) / float64(totalStudents) * 100
		}

		records = append(records, record)
	}

	return records
}

// 获取教师科目统计
func getTeacherSubjectStats(teacherID uint, tenantID uint) []TeacherSubjectStat {
	var stats []TeacherSubjectStat

	var subjects []models.Subject
	utils.WithTenant(database.DB, tenantID).Find(&subjects)

	for _, subject := range subjects {
		var questionCount, paperCount, examCount int64

		utils.WithTenant(database.DB, tenantID).Model(&models.Question{}).Where("subject_id = ? AND created_by = ?", subject.ID, teacherID).Count(&questionCount)
		utils.WithTenant(database.DB, tenantID).Model(&models.Paper{}).Where("subject_id = ? AND created_by = ?", subject.ID, teacherID).Count(&paperCount)
		utils.WithTenant(database.DB, tenantID).Model(&models.Exam{}).Joins("JOIN papers ON exams.paper_id = papers.id").Where("papers.subject_id = ? AND exams.created_by = ?", subject.ID, teacherID).Count(&examCount)

		if questionCount > 0 || paperCount > 0 || examCount > 0 {
			stats = append(stats, TeacherSubjectStat{
				Subject:       subject,
				QuestionCount: questionCount,
				PaperCount:    paperCount,
				ExamCount:     examCount,
			})
		}
	}

	return stats
}

// 获取题目分析
func getQuestionAnalysis(examID uint, tenantID uint) []QuestionAnalysis {
	var analysis []QuestionAnalysis

	// 获取考试的试卷题目
	var exam models.Exam
	utils.WithTenant(database.DB, tenantID).Preload("Paper").Preload("Paper.Questions").First(&exam, examID)

	questions := exam.Paper.Questions

	for _, question := range questions {
		var totalCount, correctCount int64

		// 统计答题总数
		utils.WithTenant(database.DB, tenantID).Model(&models.Answer{}).
			Joins("JOIN exam_records ON answers.exam_record_id = exam_records.id").
			Where("exam_records.exam_id = ? AND answers.question_id = ?", examID, question.ID).
			Count(&totalCount)

		// 统计正确答案数（这里简化处理，实际可能需要更复杂的逻辑）
		utils.WithTenant(database.DB, tenantID).Model(&models.Answer{}).
			Joins("JOIN exam_records ON answers.exam_record_id = exam_records.id").
			Where("exam_records.exam_id = ? AND answers.question_id = ? AND answers.answer = ?", examID, question.ID, question.Answer).
			Count(&correctCount)

		correctRate := 0.0
		if totalCount > 0 {
			correctRate = float64(correctCount) / float64(totalCount) * 100
		}

		analysis = append(analysis, QuestionAnalysis{
			Question:     question,
			CorrectCount: correctCount,
			TotalCount:   totalCount,
			CorrectRate:  correctRate,
		})
	}

	return analysis
}

// 获取分数分布
func getScoreDistribution(examID uint, tenantID uint) []ScoreRange {
	var distribution []ScoreRange

	// 定义分数区间
	ranges := []struct {
		name string
		min  int
		max  int
	}{
		{"0-59", 0, 59},
		{"60-69", 60, 69},
		{"70-79", 70, 79},
		{"80-89", 80, 89},
		{"90-100", 90, 100},
	}

	for _, r := range ranges {
		var count int64
		utils.WithTenant(database.DB, tenantID).Model(&models.ExamRecord{}).
			Where("exam_id = ? AND status = ? AND score >= ? AND score <= ?", examID, models.ExamCompleted, r.min, r.max).
			Count(&count)

		distribution = append(distribution, ScoreRange{
			Range: r.name,
			Count: count,
		})
	}

	return distribution
}