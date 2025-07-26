#!/bin/bash

# 在线考试系统压力测试脚本
# 要求：支持50VU并发，95%响应时间不超过100ms

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 检查k6是否安装
check_k6() {
    if ! command -v k6 &> /dev/null; then
        echo -e "${RED}错误: k6 未安装${NC}"
        echo -e "${YELLOW}请访问 https://k6.io/docs/getting-started/installation/ 安装k6${NC}"
        echo -e "${YELLOW}macOS: brew install k6${NC}"
        echo -e "${YELLOW}Ubuntu: sudo apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys C5AD17C747E3415A3642D57D77C6C491D6AC1D69${NC}"
        echo -e "${YELLOW}        echo 'deb https://dl.k6.io/deb stable main' | sudo tee /etc/apt/sources.list.d/k6.list${NC}"
        echo -e "${YELLOW}        sudo apt-get update && sudo apt-get install k6${NC}"
        exit 1
    fi
}

# 检查后端服务是否运行
check_backend() {
    echo -e "${BLUE}检查后端服务状态...${NC}"
    if ! curl -s http://localhost:8080/api/v1/health > /dev/null; then
        echo -e "${RED}错误: 后端服务未运行或无法访问${NC}"
        echo -e "${YELLOW}请确保后端服务在 http://localhost:8080 运行${NC}"
        echo -e "${YELLOW}在backend目录执行: go run main.go${NC}"
        exit 1
    fi
    echo -e "${GREEN}后端服务运行正常${NC}"
}

# 创建测试报告目录
setup_reports() {
    REPORT_DIR="reports/$(date +%Y%m%d_%H%M%S)"
    mkdir -p "$REPORT_DIR"
    echo -e "${BLUE}测试报告将保存到: $REPORT_DIR${NC}"
}

# 运行综合压力测试
run_comprehensive_test() {
    echo -e "${GREEN}=== 运行综合API压力测试 ===${NC}"
    echo -e "${YELLOW}目标: 50VU并发，95%响应时间<100ms${NC}"
    
    k6 run --out json="$REPORT_DIR/comprehensive-test.json" \
           --out csv="$REPORT_DIR/comprehensive-test.csv" \
           k6-load-test.js
    
    echo -e "${GREEN}综合测试完成${NC}"
}

# 运行核心业务接口测试
run_core_apis_test() {
    echo -e "${GREEN}=== 运行核心业务接口压力测试 ===${NC}"
    echo -e "${YELLOW}重点测试: 练习、考试等高频接口${NC}"
    
    k6 run --out json="$REPORT_DIR/core-apis-test.json" \
           --out csv="$REPORT_DIR/core-apis-test.csv" \
           k6-core-apis-test.js
    
    echo -e "${GREEN}核心接口测试完成${NC}"
}

# 运行管理接口测试
run_admin_teacher_test() {
    echo -e "${GREEN}=== 运行管理接口压力测试 ===${NC}"
    echo -e "${YELLOW}测试: 教师和管理员专用接口${NC}"
    
    k6 run --out json="$REPORT_DIR/admin-teacher-test.json" \
           --out csv="$REPORT_DIR/admin-teacher-test.csv" \
           k6-admin-teacher-test.js
    
    echo -e "${GREEN}管理接口测试完成${NC}"
}

# 运行快速验证测试
run_quick_test() {
    echo -e "${GREEN}=== 运行快速验证测试 ===${NC}"
    echo -e "${YELLOW}快速验证系统性能${NC}"
    
    # 创建快速测试配置
    cat > quick-test.js << 'EOF'
import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  stages: [
    { duration: '30s', target: 20 },
    { duration: '1m', target: 50 },
    { duration: '30s', target: 0 },
  ],
  thresholds: {
    http_req_duration: ['p(95)<100'],
    http_req_failed: ['rate<0.1'],
  },
};

const BASE_URL = 'http://localhost:8080/api/v1';

export default function () {
  // 健康检查
  const healthResponse = http.get(`${BASE_URL}/health`);
  check(healthResponse, {
    'health check status 200': (r) => r.status === 200,
    'health check RT < 50ms': (r) => r.timings.duration < 50,
  });
  
  sleep(1);
}
EOF
    
    k6 run --out json="$REPORT_DIR/quick-test.json" quick-test.js
    rm quick-test.js
    
    echo -e "${GREEN}快速测试完成${NC}"
}

# 生成测试报告摘要
generate_summary() {
    echo -e "${BLUE}=== 生成测试报告摘要 ===${NC}"
    
    SUMMARY_FILE="$REPORT_DIR/test-summary.txt"
    
    cat > "$SUMMARY_FILE" << EOF
在线考试系统压力测试报告
测试时间: $(date)
测试目标: 50VU并发，95%响应时间<100ms

测试文件:
EOF
    
    ls -la "$REPORT_DIR"/*.json >> "$SUMMARY_FILE" 2>/dev/null || true
    
    echo -e "${GREEN}测试报告摘要已生成: $SUMMARY_FILE${NC}"
}

# 显示帮助信息
show_help() {
    echo -e "${BLUE}在线考试系统压力测试工具${NC}"
    echo -e "${YELLOW}用法: $0 [选项]${NC}"
    echo ""
    echo "选项:"
    echo "  all              运行所有测试"
    echo "  comprehensive    运行综合API测试"
    echo "  core             运行核心业务接口测试"
    echo "  admin            运行管理接口测试"
    echo "  quick            运行快速验证测试"
    echo "  help             显示此帮助信息"
    echo ""
    echo "示例:"
    echo "  $0 all           # 运行所有测试"
    echo "  $0 core          # 只运行核心接口测试"
    echo "  $0 quick         # 快速验证"
}

# 主函数
main() {
    echo -e "${BLUE}在线考试系统压力测试工具${NC}"
    echo -e "${BLUE}================================${NC}"
    
    # 检查依赖
    check_k6
    check_backend
    
    # 设置报告目录
    setup_reports
    
    case "${1:-all}" in
        "all")
            run_quick_test
            run_comprehensive_test
            run_core_apis_test
            run_admin_teacher_test
            ;;
        "comprehensive")
            run_comprehensive_test
            ;;
        "core")
            run_core_apis_test
            ;;
        "admin")
            run_admin_teacher_test
            ;;
        "quick")
            run_quick_test
            ;;
        "help")
            show_help
            exit 0
            ;;
        *)
            echo -e "${RED}未知选项: $1${NC}"
            show_help
            exit 1
            ;;
    esac
    
    # 生成报告摘要
    generate_summary
    
    echo -e "${GREEN}=== 所有测试完成 ===${NC}"
    echo -e "${BLUE}测试报告位置: $REPORT_DIR${NC}"
    echo -e "${YELLOW}建议检查以下指标:${NC}"
    echo -e "${YELLOW}  - http_req_duration p(95) < 100ms${NC}"
    echo -e "${YELLOW}  - http_req_failed rate < 10%${NC}"
    echo -e "${YELLOW}  - 各接口响应时间分布${NC}"
}

# 执行主函数
main "$@"