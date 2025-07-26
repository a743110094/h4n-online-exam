#!/bin/bash

# 环境检查脚本
# 验证压力测试环境是否准备就绪

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

# 检查结果统计
PASS_COUNT=0
FAIL_COUNT=0
WARN_COUNT=0

# 记录检查结果
log_pass() {
    echo -e "${GREEN}✓ $1${NC}"
    ((PASS_COUNT++))
}

log_fail() {
    echo -e "${RED}✗ $1${NC}"
    ((FAIL_COUNT++))
}

log_warn() {
    echo -e "${YELLOW}⚠ $1${NC}"
    ((WARN_COUNT++))
}

log_info() {
    echo -e "${BLUE}ℹ $1${NC}"
}

# 检查k6是否安装
check_k6() {
    log_info "检查 k6 安装..."
    
    if command -v k6 &> /dev/null; then
        local version=$(k6 version | head -n1)
        log_pass "k6 已安装: $version"
    else
        log_fail "k6 未安装"
        echo -e "${YELLOW}安装方法:${NC}"
        echo "  macOS: brew install k6"
        echo "  Linux: sudo apt-get install k6"
        echo "  或访问: https://k6.io/docs/getting-started/installation/"
    fi
}

# 检查后端服务
check_backend() {
    log_info "检查后端服务..."
    
    if curl -s "$BASE_URL/health" > /dev/null; then
        local response=$(curl -s "$BASE_URL/health")
        log_pass "后端服务运行正常"
        echo "  响应: $response"
    else
        log_fail "后端服务未运行或无法访问"
        echo -e "${YELLOW}启动方法:${NC}"
        echo "  cd ../backend && go run main.go"
    fi
}

# 检查数据库连接
check_database() {
    log_info "检查数据库连接..."
    
    # 通过健康检查接口间接验证数据库
    local response=$(curl -s "$BASE_URL/health" 2>/dev/null || echo "")
    
    if [[ "$response" == *"ok"* ]] || [[ "$response" == *"healthy"* ]]; then
        log_pass "数据库连接正常"
    else
        log_warn "无法确认数据库状态"
        echo "  建议检查后端日志确认数据库连接"
    fi
}

# 检查测试用户
check_test_users() {
    log_info "检查测试用户..."
    
    # 尝试管理员登录
    local admin_response=$(curl -s -X POST "$BASE_URL/auth/login" \
        -H "Content-Type: application/json" \
        -H "X-Tenant-ID: $TENANT_ID" \
        -d '{
            "email": "admin@test.com",
            "password": "admin123"
        }' 2>/dev/null || echo "")
    
    if [[ "$admin_response" == *"token"* ]]; then
        log_pass "管理员用户可用"
    else
        log_fail "管理员用户不可用"
        echo "  运行: ./prepare-test-data.sh users"
    fi
    
    # 尝试学生登录
    local student_response=$(curl -s -X POST "$BASE_URL/auth/login" \
        -H "Content-Type: application/json" \
        -H "X-Tenant-ID: $TENANT_ID" \
        -d '{
            "email": "student@test.com",
            "password": "student123"
        }' 2>/dev/null || echo "")
    
    if [[ "$student_response" == *"token"* ]]; then
        log_pass "学生用户可用"
    else
        log_fail "学生用户不可用"
        echo "  运行: ./prepare-test-data.sh users"
    fi
    
    # 尝试教师登录
    local teacher_response=$(curl -s -X POST "$BASE_URL/auth/login" \
        -H "Content-Type: application/json" \
        -H "X-Tenant-ID: $TENANT_ID" \
        -d '{
            "email": "teacher@test.com",
            "password": "teacher123"
        }' 2>/dev/null || echo "")
    
    if [[ "$teacher_response" == *"token"* ]]; then
        log_pass "教师用户可用"
    else
        log_fail "教师用户不可用"
        echo "  运行: ./prepare-test-data.sh users"
    fi
}

# 检查测试数据
check_test_data() {
    log_info "检查测试数据..."
    
    # 获取管理员token
    local admin_response=$(curl -s -X POST "$BASE_URL/auth/login" \
        -H "Content-Type: application/json" \
        -H "X-Tenant-ID: $TENANT_ID" \
        -d '{
            "email": "admin@test.com",
            "password": "admin123"
        }' 2>/dev/null || echo "")
    
    local token=$(echo "$admin_response" | grep -o '"token":"[^"]*"' | cut -d'"' -f4 2>/dev/null || echo "")
    
    if [[ -z "$token" ]]; then
        log_warn "无法获取管理员token，跳过数据检查"
        return
    fi
    
    # 检查科目
    local subjects_response=$(curl -s "$BASE_URL/subjects" \
        -H "Authorization: Bearer $token" \
        -H "X-Tenant-ID: $TENANT_ID" 2>/dev/null || echo "")
    
    if [[ "$subjects_response" == *"数学"* ]] || [[ "$subjects_response" == *"name"* ]]; then
        log_pass "测试科目数据可用"
    else
        log_fail "测试科目数据不足"
        echo "  运行: ./prepare-test-data.sh subjects"
    fi
    
    # 检查题目（通过教师token）
    local teacher_response=$(curl -s -X POST "$BASE_URL/auth/login" \
        -H "Content-Type: application/json" \
        -H "X-Tenant-ID: $TENANT_ID" \
        -d '{
            "email": "teacher@test.com",
            "password": "teacher123"
        }' 2>/dev/null || echo "")
    
    local teacher_token=$(echo "$teacher_response" | grep -o '"token":"[^"]*"' | cut -d'"' -f4 2>/dev/null || echo "")
    
    if [[ -n "$teacher_token" ]]; then
        local questions_response=$(curl -s "$BASE_URL/questions" \
            -H "Authorization: Bearer $teacher_token" \
            -H "X-Tenant-ID: $TENANT_ID" 2>/dev/null || echo "")
        
        if [[ "$questions_response" == *"content"* ]] || [[ "$questions_response" == *"测试题目"* ]]; then
            log_pass "测试题目数据可用"
        else
            log_fail "测试题目数据不足"
            echo "  运行: ./prepare-test-data.sh questions"
        fi
    else
        log_warn "无法验证题目数据"
    fi
}

# 检查系统资源
check_system_resources() {
    log_info "检查系统资源..."
    
    # 检查CPU
    local cpu_count=$(sysctl -n hw.ncpu 2>/dev/null || nproc 2>/dev/null || echo "unknown")
    if [[ "$cpu_count" != "unknown" ]] && [[ $cpu_count -ge 2 ]]; then
        log_pass "CPU核心数充足: $cpu_count 核"
    else
        log_warn "CPU核心数可能不足: $cpu_count"
    fi
    
    # 检查内存
    if command -v free &> /dev/null; then
        local mem_gb=$(free -g | awk '/^Mem:/{print $2}')
        if [[ $mem_gb -ge 4 ]]; then
            log_pass "内存充足: ${mem_gb}GB"
        else
            log_warn "内存可能不足: ${mem_gb}GB"
        fi
    elif command -v vm_stat &> /dev/null; then
        # macOS
        local mem_mb=$(vm_stat | grep "Pages free" | awk '{print $3}' | sed 's/\.//' | awk '{print $1 * 4096 / 1024 / 1024}')
        log_pass "系统内存检查完成 (macOS)"
    else
        log_warn "无法检查内存状态"
    fi
    
    # 检查磁盘空间
    local disk_usage=$(df -h . | tail -1 | awk '{print $5}' | sed 's/%//')
    if [[ $disk_usage -lt 90 ]]; then
        log_pass "磁盘空间充足 (已使用 ${disk_usage}%)"
    else
        log_warn "磁盘空间不足 (已使用 ${disk_usage}%)"
    fi
}

# 检查网络连接
check_network() {
    log_info "检查网络连接..."
    
    # 检查本地端口
    if netstat -an 2>/dev/null | grep -q ":8080.*LISTEN" || lsof -i :8080 2>/dev/null | grep -q LISTEN; then
        log_pass "端口 8080 正在监听"
    else
        log_fail "端口 8080 未监听"
        echo "  确保后端服务在 8080 端口运行"
    fi
    
    # 检查网络延迟
    local ping_time=$(ping -c 1 localhost 2>/dev/null | grep "time=" | awk -F"time=" '{print $2}' | awk '{print $1}' || echo "unknown")
    if [[ "$ping_time" != "unknown" ]]; then
        log_pass "本地网络延迟: ${ping_time}ms"
    else
        log_warn "无法测试网络延迟"
    fi
}

# 检查压测脚本
check_test_scripts() {
    log_info "检查压测脚本..."
    
    local scripts=("k6-load-test.js" "k6-core-apis-test.js" "k6-admin-teacher-test.js")
    
    for script in "${scripts[@]}"; do
        if [[ -f "$script" ]]; then
            log_pass "脚本存在: $script"
        else
            log_fail "脚本缺失: $script"
        fi
    done
    
    # 检查运行脚本
    if [[ -f "run-load-tests.sh" ]] && [[ -x "run-load-tests.sh" ]]; then
        log_pass "运行脚本可执行: run-load-tests.sh"
    else
        log_fail "运行脚本问题: run-load-tests.sh"
    fi
    
    # 检查数据准备脚本
    if [[ -f "prepare-test-data.sh" ]] && [[ -x "prepare-test-data.sh" ]]; then
        log_pass "数据准备脚本可执行: prepare-test-data.sh"
    else
        log_fail "数据准备脚本问题: prepare-test-data.sh"
    fi
}

# 性能预检查
check_performance_baseline() {
    log_info "性能基线检查..."
    
    if command -v k6 &> /dev/null && curl -s "$BASE_URL/health" > /dev/null; then
        echo "  运行快速性能测试..."
        
        # 创建临时测试脚本
        cat > /tmp/quick-test.js << 'EOF'
import http from 'k6/http';
import { check } from 'k6';

export let options = {
  vus: 5,
  duration: '10s',
};

export default function() {
  let response = http.get('http://localhost:8080/api/v1/health', {
    headers: { 'X-Tenant-ID': '100' }
  });
  
  check(response, {
    'status is 200': (r) => r.status === 200,
    'response time < 100ms': (r) => r.timings.duration < 100,
  });
}
EOF
        
        # 运行快速测试
        local test_result=$(k6 run /tmp/quick-test.js 2>&1 | tail -5)
        
        if echo "$test_result" | grep -q "http_req_duration.*avg="; then
            local avg_time=$(echo "$test_result" | grep "http_req_duration" | grep -o "avg=[0-9.]*" | cut -d'=' -f2)
            if (( $(echo "$avg_time < 100" | bc -l 2>/dev/null || echo "1") )); then
                log_pass "基线性能良好 (平均响应时间: ${avg_time}ms)"
            else
                log_warn "基线性能较慢 (平均响应时间: ${avg_time}ms)"
            fi
        else
            log_warn "无法获取性能基线数据"
        fi
        
        # 清理临时文件
        rm -f /tmp/quick-test.js
    else
        log_warn "跳过性能基线检查 (k6或后端服务不可用)"
    fi
}

# 生成报告
generate_report() {
    echo ""
    echo -e "${BLUE}=== 环境检查报告 ===${NC}"
    echo -e "${GREEN}通过: $PASS_COUNT${NC}"
    echo -e "${RED}失败: $FAIL_COUNT${NC}"
    echo -e "${YELLOW}警告: $WARN_COUNT${NC}"
    echo ""
    
    if [[ $FAIL_COUNT -eq 0 ]]; then
        echo -e "${GREEN}✓ 环境检查通过！可以开始压力测试${NC}"
        echo ""
        echo -e "${YELLOW}建议的测试步骤:${NC}"
        echo "1. 准备测试数据: ./prepare-test-data.sh"
        echo "2. 运行压力测试: ./run-load-tests.sh"
        echo "3. 监控性能: ./monitor-performance.sh (另开终端)"
        return 0
    else
        echo -e "${RED}✗ 环境检查失败！请解决上述问题后重试${NC}"
        echo ""
        echo -e "${YELLOW}常见解决方案:${NC}"
        echo "1. 安装 k6: brew install k6"
        echo "2. 启动后端: cd ../backend && go run main.go"
        echo "3. 准备数据: ./prepare-test-data.sh"
        return 1
    fi
}

# 显示帮助
show_help() {
    echo -e "${BLUE}环境检查脚本${NC}"
    echo -e "${YELLOW}用法: $0 [选项]${NC}"
    echo ""
    echo "选项:"
    echo "  check            完整环境检查 (默认)"
    echo "  quick            快速检查"
    echo "  performance      性能基线检查"
    echo "  help             显示帮助"
    echo ""
    echo "示例:"
    echo "  $0               # 完整检查"
    echo "  $0 quick         # 快速检查"
    echo "  $0 performance   # 性能检查"
}

# 主函数
main() {
    echo -e "${BLUE}=== 压力测试环境检查 ===${NC}"
    echo ""
    
    case "${1:-check}" in
        "check")
            check_k6
            check_backend
            check_database
            check_test_users
            check_test_data
            check_system_resources
            check_network
            check_test_scripts
            check_performance_baseline
            generate_report
            ;;
        "quick")
            check_k6
            check_backend
            check_test_scripts
            generate_report
            ;;
        "performance")
            check_performance_baseline
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