package models

import (
	"time"
)

// 租户模型
type Tenant struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null"`
	Code        string    `json:"code" gorm:"uniqueIndex;not null"`
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// 用户角色枚举
type UserRole string

const (
	RoleAdmin   UserRole = "admin"
	RoleTeacher UserRole = "teacher"
	RoleStudent UserRole = "student"
)

// 题目类型枚举
type QuestionType string

const (
	SingleChoice   QuestionType = "single_choice"
	MultipleChoice QuestionType = "multiple_choice"
	TrueFalse      QuestionType = "true_false"
	ShortAnswer    QuestionType = "short_answer"
)

// 考试状态枚举
type ExamStatus string

const (
	ExamDraft     ExamStatus = "draft"
	ExamPublished ExamStatus = "published"
	ExamStarted   ExamStatus = "started"
	ExamEnded     ExamStatus = "ended"
)

// 考试记录状态枚举
type ExamRecordStatus string

const (
	ExamNotStarted ExamRecordStatus = "not_started"
	ExamInProgress ExamRecordStatus = "in_progress"
	ExamCompleted  ExamRecordStatus = "completed"
	ExamTimeout    ExamRecordStatus = "timeout"
)

// 用户模型
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	TenantID  uint      `json:"tenant_id" gorm:"not null;index;default:100"`
	Username  string    `json:"username" gorm:"uniqueIndex;not null"`
	Email     string    `json:"email" gorm:"uniqueIndex;not null"`
	Password  string    `json:"-" gorm:"not null"`
	Role      UserRole  `json:"role" gorm:"not null;default:'student'"`
	Name      string    `json:"name" gorm:"not null"`
	Avatar    string    `json:"avatar"`
	IsActive  bool      `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// 科目模型
type Subject struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	TenantID    uint      `json:"tenant_id" gorm:"not null;index;default:100"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	ParentID    *uint     `json:"parent_id"`
	Parent      *Subject  `json:"parent" gorm:"foreignKey:ParentID"`
	Children    []Subject `json:"children" gorm:"foreignKey:ParentID"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// 题目状态枚举
type QuestionStatus string

const (
	QuestionDraft     QuestionStatus = "draft"
	QuestionPublished QuestionStatus = "published"
	QuestionArchived  QuestionStatus = "archived"
)

// 题目模型
type Question struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	TenantID       uint           `json:"tenant_id" gorm:"not null;index;default:100"`
	SubjectID      uint           `json:"subject_id" gorm:"not null"`
	Subject        Subject        `json:"subject" gorm:"foreignKey:SubjectID"`
	Type           QuestionType   `json:"type" gorm:"not null"`
	Title          string         `json:"title" gorm:"not null"`
	Content        string         `json:"content" gorm:"type:text"`
	Options        string         `json:"options" gorm:"type:text"` // JSON格式存储选项
	Answer         string         `json:"answer" gorm:"not null"`
	Explanation    string         `json:"explanation" gorm:"type:text"`
	Difficulty     int            `json:"difficulty" gorm:"default:1"` // 1-5难度等级
	Score          int            `json:"score" gorm:"default:1"`      // 题目分值
	Status         QuestionStatus `json:"status" gorm:"default:'published'"`
	KnowledgePoint string         `json:"knowledge_point" gorm:"default:''"`  // 知识点
	UsageCount     int            `json:"usage_count" gorm:"default:0"`       // 使用次数
	CorrectRate    float64        `json:"correct_rate" gorm:"default:0"`      // 正确率(0-1)
	CreatedBy      uint           `json:"created_by"`
	Creator        User           `json:"creator" gorm:"foreignKey:CreatedBy"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

// 试卷模型
type Paper struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	TenantID    uint      `json:"tenant_id" gorm:"not null;index;default:100"`
	Title       string    `json:"title" gorm:"not null"`
	Description string    `json:"description"`
	SubjectID   uint      `json:"subject_id"`
	Subject     Subject   `json:"subject" gorm:"foreignKey:SubjectID"`
	TotalScore  int       `json:"total_score" gorm:"default:0"`
	Duration    int       `json:"duration" gorm:"default:60"` // 考试时长(分钟)
	CreatedBy   uint      `json:"created_by"`
	Creator     User      `json:"creator" gorm:"foreignKey:CreatedBy"`
	Questions   []Question `json:"questions" gorm:"many2many:paper_questions;"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// 试卷题目关联模型
type PaperQuestion struct {
	ID         uint     `json:"id" gorm:"primaryKey"`
	PaperID    uint     `json:"paper_id" gorm:"not null"`
	QuestionID uint     `json:"question_id" gorm:"not null"`
	Score      float64  `json:"score" gorm:"default:1"` // 题目在该试卷中的分值
	Order      int      `json:"order" gorm:"default:0"` // 题目顺序
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// 考试模型
type Exam struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	TenantID    uint       `json:"tenant_id" gorm:"not null;index;default:100"`
	Title       string     `json:"title" gorm:"not null"`
	Description string     `json:"description"`
	PaperID     uint       `json:"paper_id" gorm:"not null"`
	Paper       Paper      `json:"paper" gorm:"foreignKey:PaperID"`
	StartTime   time.Time  `json:"start_time"`
	EndTime     time.Time  `json:"end_time"`
	Duration    int        `json:"duration"` // 考试时长(分钟)
	Status      ExamStatus `json:"status" gorm:"default:'draft'"`
	StudentIDs  string     `json:"student_ids" gorm:"type:text"` // JSON格式存储学生ID列表
	CreatedBy   uint       `json:"created_by"`
	Creator     User       `json:"creator" gorm:"foreignKey:CreatedBy"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// 考试参与记录
type ExamRecord struct {
	ID           uint              `json:"id" gorm:"primaryKey"`
	TenantID     uint              `json:"tenant_id" gorm:"not null;index;default:100"`
	ExamID       uint              `json:"exam_id" gorm:"not null"`
	Exam         Exam              `json:"exam" gorm:"foreignKey:ExamID"`
	StudentID    uint              `json:"student_id" gorm:"not null"`
	Student      User              `json:"student" gorm:"foreignKey:StudentID"`
	UserID       uint              `json:"user_id" gorm:"not null"`
	User         User              `json:"user" gorm:"foreignKey:UserID"` 
	Answers      string            `json:"answers" gorm:"type:text"` // JSON格式存储答案
	Score        float64           `json:"score" gorm:"default:0"`
	CorrectCount int               `json:"correct_count" gorm:"default:0"`
	TotalCount   int               `json:"total_count" gorm:"default:0"`
	SubmitTime   time.Time         `json:"submit_time"`
	StartTime    time.Time         `json:"start_time"`
	EndTime      *time.Time        `json:"end_time"`
	TotalScore   int               `json:"total_score"`
	Status       ExamRecordStatus  `json:"status" gorm:"default:'not_started'"`
	IsFinished   bool              `json:"is_finished" gorm:"default:false"`
	ExtraTime    int               `json:"extra_time" gorm:"default:0"` // 额外时间(分钟)
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
}

// 答题记录
type Answer struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	TenantID     uint      `json:"tenant_id" gorm:"not null;index;default:100"`
	ExamRecordID uint      `json:"exam_record_id" gorm:"not null"`
	ExamRecord   ExamRecord `json:"exam_record" gorm:"foreignKey:ExamRecordID"`
	QuestionID   uint      `json:"question_id" gorm:"not null"`
	Question     Question  `json:"question" gorm:"foreignKey:QuestionID"`
	Answer       string    `json:"answer"`
	IsCorrect    *bool     `json:"is_correct"`
	Score        *int      `json:"score"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// 练习记录
type PracticeRecord struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	TenantID     uint      `json:"tenant_id" gorm:"not null;index;default:100"`
	UserID       uint      `json:"user_id" gorm:"not null"`
	User         User      `json:"user" gorm:"foreignKey:UserID"`
	SubjectID    uint      `json:"subject_id" gorm:"not null"`
	Subject      Subject   `json:"subject" gorm:"foreignKey:SubjectID"`
	Title        string    `json:"title" gorm:"not null"`
	Description  string    `json:"description"`
	QuestionIDs  string    `json:"question_ids" gorm:"type:text"` // JSON格式存储题目ID列表
	TotalCount   int       `json:"total_count" gorm:"default:0"`
	CorrectCount int       `json:"correct_count" gorm:"default:0"`
	WrongCount   int       `json:"wrong_count" gorm:"default:0"`
	Score        int       `json:"score" gorm:"default:0"`
	Duration     int       `json:"duration" gorm:"default:0"` // 练习时长(秒)
	Difficulty   int       `json:"difficulty" gorm:"default:1"` // 1-5难度等级
	PracticeType string    `json:"practice_type" gorm:"default:'sequence'"` // sequence, random, wrong
	IsCompleted  bool      `json:"is_completed" gorm:"default:false"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// 练习答题记录
type PracticeAnswer struct {
	ID               uint           `json:"id" gorm:"primaryKey"`
	TenantID         uint           `json:"tenant_id" gorm:"not null;index;default:100"`
	PracticeRecordID uint           `json:"practice_record_id" gorm:"not null"`
	PracticeRecord   PracticeRecord `json:"practice_record" gorm:"foreignKey:PracticeRecordID"`
	QuestionID       uint           `json:"question_id" gorm:"not null"`
	Question         Question       `json:"question" gorm:"foreignKey:QuestionID"`
	Answer           string         `json:"answer"`
	IsCorrect        *bool          `json:"is_correct"`
	Score            *int           `json:"score"`
	TimeSpent        int            `json:"time_spent" gorm:"default:0"` // 答题用时(秒)
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
}

// 推荐练习
type PracticeRecommendation struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	TenantID       uint      `json:"tenant_id" gorm:"not null;index;default:100"`
	Title          string    `json:"title" gorm:"not null"`
	Description    string    `json:"description"`
	SubjectID      uint      `json:"subject_id" gorm:"not null"`
	Subject        Subject   `json:"subject" gorm:"foreignKey:SubjectID"`
	Difficulty     int       `json:"difficulty" gorm:"default:1"` // 1-5难度等级
	QuestionCount  int       `json:"question_count" gorm:"default:10"`
	EstimatedTime  int       `json:"estimated_time" gorm:"default:15"` // 预估时间(分钟)
	Rating         float64   `json:"rating" gorm:"default:4.5"`
	KnowledgePoint string    `json:"knowledge_point"`
	QuestionTypes  string    `json:"question_types" gorm:"type:text"` // JSON格式存储题目类型
	IsActive       bool      `json:"is_active" gorm:"default:true"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// AI问答记录
type AIChat struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	TenantID  uint      `json:"tenant_id" gorm:"not null;index;default:100"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	User      User      `json:"user" gorm:"foreignKey:UserID"`
	Message   string    `json:"message" gorm:"type:text;not null"`
	Response  string    `json:"response" gorm:"type:text"`
	Context   string    `json:"context" gorm:"type:text"`
	CreatedAt time.Time `json:"created_at"`
}

// 用户统计模型
type UserStats struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	TenantID     uint      `json:"tenant_id" gorm:"not null;index;default:100"`
	UserID       uint      `json:"user_id" gorm:"not null;uniqueIndex:idx_user_tenant"`
	User         User      `json:"user" gorm:"foreignKey:UserID"`
	ExamCount    int       `json:"exam_count" gorm:"default:0"`
	TotalScore   float64   `json:"total_score" gorm:"default:0"`
	AverageScore float64   `json:"average_score" gorm:"default:0"`
	HighestScore float64   `json:"highest_score" gorm:"default:0"`
	LowestScore  float64   `json:"lowest_score" gorm:"default:0"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// 考试统计模型
type ExamStats struct {
	ID               uint      `json:"id" gorm:"primaryKey"`
	TenantID         uint      `json:"tenant_id" gorm:"not null;index;default:100"`
	ExamID           uint      `json:"exam_id" gorm:"not null;uniqueIndex:idx_exam_tenant"`
	Exam             Exam      `json:"exam" gorm:"foreignKey:ExamID"`
	ParticipantCount int       `json:"participant_count" gorm:"default:0"`
	AverageScore     float64   `json:"average_score" gorm:"default:0"`
	HighestScore     float64   `json:"highest_score" gorm:"default:0"`
	LowestScore      float64   `json:"lowest_score" gorm:"default:0"`
	PassRate         float64   `json:"pass_rate" gorm:"default:0"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// 考试报告模型
type ExamReport struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	TenantID     uint      `json:"tenant_id" gorm:"not null;index;default:100"`
	ExamID       uint      `json:"exam_id" gorm:"not null"`
	Exam         Exam      `json:"exam" gorm:"foreignKey:ExamID"`
	UserID       uint      `json:"user_id" gorm:"not null"`
	User         User      `json:"user" gorm:"foreignKey:UserID"`
	Score        float64   `json:"score"`
	CorrectCount int       `json:"correct_count"`
	TotalCount   int       `json:"total_count"`
	PassRate     float64   `json:"pass_rate"`
	Rank         int       `json:"rank"`
	ReportData   string    `json:"report_data" gorm:"type:text"` // JSON格式存储详细报告数据
	GeneratedAt  time.Time `json:"generated_at"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}