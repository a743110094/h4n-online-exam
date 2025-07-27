package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// 数据库模型
type Subject struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	TenantID    uint      `json:"tenant_id" gorm:"not null;index;default:100"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	ParentID    *uint     `json:"parent_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Question struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	TenantID       uint      `json:"tenant_id" gorm:"not null;index;default:100"`
	SubjectID      uint      `json:"subject_id" gorm:"not null"`
	Type           string    `json:"type" gorm:"not null"`
	Title          string    `json:"title" gorm:"not null"`
	Content        string    `json:"content" gorm:"type:text"`
	Options        string    `json:"options" gorm:"type:text"`
	Answer         string    `json:"answer" gorm:"not null"`
	Explanation    string    `json:"explanation" gorm:"type:text"`
	Difficulty     int       `json:"difficulty" gorm:"default:1"`
	Score          int       `json:"score" gorm:"default:1"`
	Status         string    `json:"status" gorm:"default:'published'"`
	KnowledgePoint string    `json:"knowledge_point" gorm:"default:''"`
	UsageCount     int       `json:"usage_count" gorm:"default:0"`
	CorrectRate    float64   `json:"correct_rate" gorm:"default:0"`
	CreatedBy      uint      `json:"created_by"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// 题目数据结构
type QuestionData struct {
	Number      int
	Type        string
	Title       string
	Options     []string
	Answer      string
	Explanation string
}

func main() {
	// 连接数据库
	dsn := "host=localhost user=postgres password=password dbname=online_exam_system port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// 创建华为AI认证科目
	subject := Subject{
		TenantID:    100,
		Name:        "华为AI认证(HCIA-AI)",
		Description: "华为认证人工智能工程师(HCIA-AI)认证考试题库",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// 检查科目是否已存在
	var existingSubject Subject
	result := db.Where("name = ? AND tenant_id = ?", subject.Name, subject.TenantID).First(&existingSubject)
	if result.Error == nil {
		fmt.Printf("科目 '%s' 已存在，ID: %d\n", existingSubject.Name, existingSubject.ID)
		subject = existingSubject
	} else {
		// 创建新科目
		if err := db.Create(&subject).Error; err != nil {
			log.Fatal("Failed to create subject:", err)
		}
		fmt.Printf("成功创建科目 '%s'，ID: %d\n", subject.Name, subject.ID)
	}

	// 读取并解析题目文件
	filePath := "/Users/lrx/solo-test/data_import/HCIA35/doc.md"
	questions, err := parseQuestions(filePath)
	if err != nil {
		log.Fatal("Failed to parse questions:", err)
	}

	fmt.Printf("解析到 %d 道题目\n", len(questions))

	// 导入题目到数据库
	successCount := 0
	for i, q := range questions {
		question := Question{
			TenantID:       100,
			SubjectID:      subject.ID,
			Type:           q.Type,
			Title:          q.Title,
			Content:        q.Title,
			Options:        formatOptions(q.Options),
			Answer:         q.Answer,
			Explanation:    q.Explanation,
			Difficulty:     2, // 默认难度为2
			Score:          getScoreByType(q.Type),
			Status:         "published",
			KnowledgePoint: "华为AI认证",
			UsageCount:     0,
			CorrectRate:    0.0,
			CreatedBy:      2, // 默认创建者ID为2
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		if err := db.Create(&question).Error; err != nil {
			fmt.Printf("Failed to create question %d: %v\n", i+1, err)
		} else {
			successCount++
			if successCount%50 == 0 {
				fmt.Printf("已导入 %d 道题目...\n", successCount)
			}
		}
	}

	fmt.Printf("\n导入完成！成功导入 %d 道题目到科目 '%s'\n", successCount, subject.Name)
}

// 解析题目文件
func parseQuestions(filePath string) ([]QuestionData, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var questions []QuestionData
	scanner := bufio.NewScanner(file)

	// 正则表达式
	questionRe := regexp.MustCompile(`^(\d+)\.(.+)$`)
	optionRe := regexp.MustCompile(`^([A-Z])、(.+)$`)
	answerRe := regexp.MustCompile(`^答案：(.+)$`)

	currentSection := ""
	currentQuestion := QuestionData{}
	currentOptions := []string{}
	questionStarted := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		// 检测题目类型
		if strings.Contains(line, "一、单选题") {
			currentSection = "single_choice"
			continue
		} else if strings.Contains(line, "二、多选题") {
			currentSection = "multiple_choice"
			continue
		} else if strings.Contains(line, "三、判断题") {
			currentSection = "true_false"
			continue
		} else if strings.Contains(line, "填空题") {
			currentSection = "short_answer"
			continue
		} else if strings.Contains(line, "简答题") {
			currentSection = "short_answer"
			continue
		}

		// 解析题目
		if currentSection != "" {
			// 检查是否是新题目开始
			if matches := questionRe.FindStringSubmatch(line); matches != nil && len(matches) == 3 {
				// 保存上一题
				if questionStarted {
					currentQuestion.Options = currentOptions
					questions = append(questions, currentQuestion)
				}

				// 开始新题目
				number, _ := strconv.Atoi(matches[1])
				currentQuestion = QuestionData{
					Number: number,
					Type:   currentSection,
					Title:  strings.TrimSpace(matches[2]),
				}
				currentOptions = []string{}
				questionStarted = true
				continue
			}

			// 解析选项
			if matches := optionRe.FindStringSubmatch(line); matches != nil && len(matches) == 3 {
				option := fmt.Sprintf("%s. %s", matches[1], strings.TrimSpace(matches[2]))
				currentOptions = append(currentOptions, option)
				continue
			}

			// 解析答案
			if matches := answerRe.FindStringSubmatch(line); matches != nil && len(matches) == 2 {
				currentQuestion.Answer = strings.TrimSpace(matches[1])
				continue
			}

			// 对于判断题，如果没有选项，添加默认选项
			if currentSection == "true_false" && questionStarted && len(currentOptions) == 0 {
				currentOptions = []string{"A. 正确", "B. 错误"}
			}
		}
	}

	// 保存最后一题
	if questionStarted {
		currentQuestion.Options = currentOptions
		questions = append(questions, currentQuestion)
	}

	return questions, nil
}

// 格式化选项为JSON字符串
func formatOptions(options []string) string {
	if len(options) == 0 {
		return "[]"
	}

	optionsJSON, _ := json.Marshal(options)
	return string(optionsJSON)
}

// 根据题目类型获取分值
func getScoreByType(questionType string) int {
	switch questionType {
	case "single_choice":
		return 2
	case "multiple_choice":
		return 3
	case "true_false":
		return 1
	case "short_answer":
		return 5
	default:
		return 2
	}
}