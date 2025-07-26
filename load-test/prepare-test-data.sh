#!/bin/bash

# 测试数据准备脚本
# 为压力测试创建必要的用户、科目、题目等数据

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 配置
BASE_URL="http://localhost:8080/api/v1"
TENANT_ID="100"
ADMIN_EMAIL="admin@test.com"
ADMIN_PASSWORD="admin123"

# 检查后端服务
check_backend() {
    echo -e "${BLUE}检查后端服务...${NC}"
    if ! curl -s "$BASE_URL/health" > /dev/null; then
        echo -e "${RED}错误: 后端服务未运行${NC}"
        echo -e "${YELLOW}请先启动后端服务: cd ../backend && go run main.go${NC}"
        exit 1
    fi
    echo -e "${GREEN}后端服务运行正常${NC}"
}

# 管理员登录获取token
get_admin_token() {
    echo -e "${BLUE}管理员登录...${NC}"
    
    local response=$(curl -s -X POST "$BASE_URL/auth/login" \
        -H "Content-Type: application/json" \
        -H "X-Tenant-ID: $TENANT_ID" \
        -d "{
            \"email\": \"$ADMIN_EMAIL\",
            \"password\": \"$ADMIN_PASSWORD\"
        }")
    
    local token=$(echo "$response" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
    
    if [[ -z "$token" ]]; then
        echo -e "${RED}管理员登录失败${NC}"
        echo "响应: $response"
        exit 1
    fi
    
    echo -e "${GREEN}管理员登录成功${NC}"
    echo "$token"
}

# 创建测试用户
create_test_users() {
    local token=$1
    echo -e "${BLUE}创建测试用户...${NC}"
    
    # 创建教师用户
    local teacher_response=$(curl -s -X POST "$BASE_URL/admin/users" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $token" \
        -H "X-Tenant-ID: $TENANT_ID" \
        -d "{
            \"name\": \"测试教师\",
            \"email\": \"teacher@test.com\",
            \"password\": \"teacher123\",
            \"role\": \"teacher\"
        }")
    
    # 创建学生用户
    for i in {1..5}; do
        curl -s -X POST "$BASE_URL/admin/users" \
            -H "Content-Type: application/json" \
            -H "Authorization: Bearer $token" \
            -H "X-Tenant-ID: $TENANT_ID" \
            -d "{
                \"name\": \"测试学生$i\",
                \"email\": \"student$i@test.com\",
                \"password\": \"student123\",
                \"role\": \"student\"
            }" > /dev/null
    done
    
    # 创建默认学生用户（用于压测脚本）
    curl -s -X POST "$BASE_URL/admin/users" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $token" \
        -H "X-Tenant-ID: $TENANT_ID" \
        -d "{
            \"name\": \"默认学生\",
            \"email\": \"student@test.com\",
            \"password\": \"student123\",
            \"role\": \"student\"
        }" > /dev/null
    
    echo -e "${GREEN}测试用户创建完成${NC}"
}

# 创建科目
create_subjects() {
    local token=$1
    echo -e "${BLUE}创建测试科目...${NC}"
    
    local subjects=("数学" "语文" "英语" "物理" "化学")
    
    for subject in "${subjects[@]}"; do
        curl -s -X POST "$BASE_URL/admin/subjects" \
            -H "Content-Type: application/json" \
            -H "Authorization: Bearer $token" \
            -H "X-Tenant-ID: $TENANT_ID" \
            -d "{
                \"name\": \"$subject\",
                \"description\": \"$subject 科目\"
            }" > /dev/null
    done
    
    echo -e "${GREEN}测试科目创建完成${NC}"
}

# 获取教师token
get_teacher_token() {
    echo -e "${BLUE}教师登录...${NC}"
    
    local response=$(curl -s -X POST "$BASE_URL/auth/login" \
        -H "Content-Type: application/json" \
        -H "X-Tenant-ID: $TENANT_ID" \
        -d "{
            \"email\": \"teacher@test.com\",
            \"password\": \"teacher123\"
        }")
    
    local token=$(echo "$response" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
    
    if [[ -z "$token" ]]; then
        echo -e "${RED}教师登录失败${NC}"
        return 1
    fi
    
    echo -e "${GREEN}教师登录成功${NC}"
    echo "$token"
}

# 创建题目
create_questions() {
    local token=$1
    echo -e "${BLUE}创建测试题目...${NC}"
    
    local difficulties=("easy" "medium" "hard")
    local types=("single_choice" "multiple_choice" "true_false")
    
    # 为每个科目创建题目
    for subject_id in {1..5}; do
        for i in {1..20}; do
            local difficulty=${difficulties[$((RANDOM % 3))]}
            local type=${types[$((RANDOM % 3))]}
            
            local options='["选项A", "选项B", "选项C", "选项D"]'
            local correct_answer='["A"]'
            
            if [[ "$type" == "multiple_choice" ]]; then
                correct_answer='["A", "B"]'
            elif [[ "$type" == "true_false" ]]; then
                options='["正确", "错误"]'
                correct_answer='["正确"]'
            fi
            
            curl -s -X POST "$BASE_URL/teacher/questions" \
                -H "Content-Type: application/json" \
                -H "Authorization: Bearer $token" \
                -H "X-Tenant-ID: $TENANT_ID" \
                -d "{
                    \"subject_id\": $subject_id,
                    \"type\": \"$type\",
                    \"difficulty\": \"$difficulty\",
                    \"content\": \"测试题目 $i - 科目 $subject_id ($difficulty)\",
                    \"options\": $options,
                    \"correct_answer\": $correct_answer,
                    \"explanation\": \"这是题目 $i 的解析\",
                    \"score\": 5
                }" > /dev/null
        done
    done
    
    echo -e "${GREEN}测试题目创建完成 (每个科目20道题)${NC}"
}

# 创建试卷
create_papers() {
    local token=$1
    echo -e "${BLUE}创建测试试卷...${NC}"
    
    # 为每个科目创建试卷
    for subject_id in {1..5}; do
        curl -s -X POST "$BASE_URL/teacher/papers/auto" \
            -H "Content-Type: application/json" \
            -H "Authorization: Bearer $token" \
            -H "X-Tenant-ID: $TENANT_ID" \
            -d "{
                \"name\": \"科目 $subject_id 测试试卷\",
                \"subject_id\": $subject_id,
                \"total_score\": 100,
                \"question_config\": {
                    \"single_choice\": { \"count\": 10, \"score\": 5 },
                    \"multiple_choice\": { \"count\": 5, \"score\": 10 }
                },
                \"difficulty_distribution\": {
                    \"easy\": 0.3,
                    \"medium\": 0.5,
                    \"hard\": 0.2
                }
            }" > /dev/null
    done
    
    echo -e "${GREEN}测试试卷创建完成${NC}"
}

# 创建考试
create_exams() {
    local token=$1
    echo -e "${BLUE}创建测试考试...${NC}"
    
    # 创建几个考试
    for i in {1..3}; do
        local start_time=$(date -d "+1 hour" -Iseconds 2>/dev/null || date -v+1H -Iseconds)
        local end_time=$(date -d "+3 hours" -Iseconds 2>/dev/null || date -v+3H -Iseconds)
        
        curl -s -X POST "$BASE_URL/teacher/exams" \
            -H "Content-Type: application/json" \
            -H "Authorization: Bearer $token" \
            -H "X-Tenant-ID: $TENANT_ID" \
            -d "{
                \"title\": \"压力测试考试 $i\",
                \"paper_id\": $i,
                \"start_time\": \"$start_time\",
                \"end_time\": \"$end_time\",
                \"duration\": 120,
                \"description\": \"这是压力测试用的考试 $i\"
            }" > /dev/null
    done
    
    echo -e "${GREEN}测试考试创建完成${NC}"
}

# 验证数据
verify_data() {
    local token=$1
    echo -e "${BLUE}验证测试数据...${NC}"
    
    # 检查用户数量
    local users_response=$(curl -s "$BASE_URL/admin/users" \
        -H "Authorization: Bearer $token" \
        -H "X-Tenant-ID: $TENANT_ID")
    
    # 检查科目数量
    local subjects_response=$(curl -s "$BASE_URL/subjects" \
        -H "Authorization: Bearer $token" \
        -H "X-Tenant-ID: $TENANT_ID")
    
    # 检查题目数量
    local questions_response=$(curl -s "$BASE_URL/questions" \
        -H "Authorization: Bearer $token" \
        -H "X-Tenant-ID: $TENANT_ID")
    
    echo -e "${GREEN}数据验证完成${NC}"
    echo -e "${YELLOW}建议检查后端日志确认数据创建情况${NC}"
}

# 显示帮助
show_help() {
    echo -e "${BLUE}测试数据准备脚本${NC}"
    echo -e "${YELLOW}用法: $0 [选项]${NC}"
    echo ""
    echo "选项:"
    echo "  prepare          准备所有测试数据 (默认)"
    echo "  users            只创建用户"
    echo "  subjects         只创建科目"
    echo "  questions        只创建题目"
    echo "  papers           只创建试卷"
    echo "  exams            只创建考试"
    echo "  verify           验证数据"
    echo "  help             显示帮助"
    echo ""
    echo "示例:"
    echo "  $0 prepare       # 准备所有测试数据"
    echo "  $0 users         # 只创建用户"
    echo "  $0 verify        # 验证数据"
}

# 主函数
main() {
    echo -e "${BLUE}=== 测试数据准备工具 ===${NC}"
    
    # 检查后端服务
    check_backend
    
    case "${1:-prepare}" in
        "prepare")
            echo -e "${YELLOW}开始准备所有测试数据...${NC}"
            
            # 获取管理员token
            admin_token=$(get_admin_token)
            
            # 创建基础数据
            create_test_users "$admin_token"
            create_subjects "$admin_token"
            
            # 获取教师token
            teacher_token=$(get_teacher_token)
            
            # 创建教学数据
            create_questions "$teacher_token"
            create_papers "$teacher_token"
            create_exams "$teacher_token"
            
            # 验证数据
            verify_data "$admin_token"
            
            echo -e "${GREEN}=== 所有测试数据准备完成 ===${NC}"
            echo -e "${YELLOW}现在可以运行压力测试了！${NC}"
            ;;
        "users")
            admin_token=$(get_admin_token)
            create_test_users "$admin_token"
            ;;
        "subjects")
            admin_token=$(get_admin_token)
            create_subjects "$admin_token"
            ;;
        "questions")
            teacher_token=$(get_teacher_token)
            create_questions "$teacher_token"
            ;;
        "papers")
            teacher_token=$(get_teacher_token)
            create_papers "$teacher_token"
            ;;
        "exams")
            teacher_token=$(get_teacher_token)
            create_exams "$teacher_token"
            ;;
        "verify")
            admin_token=$(get_admin_token)
            verify_data "$admin_token"
            ;;
        "help")
            show_help
            ;;
        *)
            echo -e "${RED}未知选项: $1${NC}"
            show_help
            exit 1
            ;;
    esac
}

# 执行主函数
main "$@"