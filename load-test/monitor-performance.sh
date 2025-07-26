#!/bin/bash

# 性能监控脚本 - 配合压力测试使用
# 监控CPU、内存、网络等系统资源

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 配置
MONITOR_INTERVAL=5  # 监控间隔（秒）
BACKEND_PORT=8080
FRONTEND_PORT=5173
LOG_FILE="performance-monitor.log"

# 检查系统类型
OS_TYPE=$(uname -s)

# 获取CPU使用率
get_cpu_usage() {
    if [[ "$OS_TYPE" == "Darwin" ]]; then
        # macOS
        top -l 1 -n 0 | grep "CPU usage" | awk '{print $3}' | sed 's/%//'
    else
        # Linux
        top -bn1 | grep "Cpu(s)" | awk '{print $2}' | sed 's/%us,//'
    fi
}

# 获取内存使用率
get_memory_usage() {
    if [[ "$OS_TYPE" == "Darwin" ]]; then
        # macOS
        vm_stat | awk '
        /Pages free/ { free = $3 }
        /Pages active/ { active = $3 }
        /Pages inactive/ { inactive = $3 }
        /Pages speculative/ { speculative = $3 }
        /Pages wired/ { wired = $3 }
        END {
            total = free + active + inactive + speculative + wired
            used = active + inactive + wired
            printf "%.1f", (used / total) * 100
        }'
    else
        # Linux
        free | awk '/Mem:/ { printf "%.1f", ($3/$2) * 100.0 }'
    fi
}

# 获取进程CPU和内存使用情况
get_process_stats() {
    local process_name=$1
    local port=$2
    
    if [[ "$OS_TYPE" == "Darwin" ]]; then
        # macOS - 通过端口查找进程
        local pid=$(lsof -ti:$port 2>/dev/null | head -1)
        if [[ -n "$pid" ]]; then
            ps -p $pid -o %cpu,%mem,comm | tail -1
        else
            echo "0.0 0.0 $process_name(未运行)"
        fi
    else
        # Linux
        local pid=$(ss -tlnp | grep ":$port" | awk '{print $6}' | cut -d',' -f2 | cut -d'=' -f2 | head -1)
        if [[ -n "$pid" ]]; then
            ps -p $pid -o %cpu,%mem,comm --no-headers
        else
            echo "0.0 0.0 $process_name(未运行)"
        fi
    fi
}

# 获取网络连接数
get_connection_count() {
    local port=$1
    if [[ "$OS_TYPE" == "Darwin" ]]; then
        netstat -an | grep ":$port" | grep ESTABLISHED | wc -l | tr -d ' '
    else
        ss -tn | grep ":$port" | grep ESTAB | wc -l
    fi
}

# 检查后端服务响应时间
check_backend_response() {
    local start_time=$(date +%s%3N)
    if curl -s -o /dev/null -w "%{http_code}" http://localhost:$BACKEND_PORT/api/v1/health > /dev/null 2>&1; then
        local end_time=$(date +%s%3N)
        local response_time=$((end_time - start_time))
        echo "$response_time"
    else
        echo "timeout"
    fi
}

# 显示实时监控信息
show_monitor_header() {
    clear
    echo -e "${BLUE}=== 在线考试系统性能监控 ===${NC}"
    echo -e "${BLUE}监控间隔: ${MONITOR_INTERVAL}秒 | 日志文件: $LOG_FILE${NC}"
    echo -e "${YELLOW}按 Ctrl+C 停止监控${NC}"
    echo ""
    printf "%-20s %-10s %-10s %-15s %-15s %-10s %-15s\n" \
           "时间" "CPU%" "内存%" "后端CPU%" "后端内存%" "连接数" "响应时间(ms)"
    echo "--------------------------------------------------------------------------------"
}

# 记录监控数据
log_performance_data() {
    local timestamp=$(date '+%Y-%m-%d %H:%M:%S')
    local cpu_usage=$(get_cpu_usage)
    local memory_usage=$(get_memory_usage)
    local backend_stats=$(get_process_stats "backend" $BACKEND_PORT)
    local backend_cpu=$(echo $backend_stats | awk '{print $1}')
    local backend_mem=$(echo $backend_stats | awk '{print $2}')
    local connections=$(get_connection_count $BACKEND_PORT)
    local response_time=$(check_backend_response)
    
    # 显示到终端
    printf "%-20s %-10s %-10s %-15s %-15s %-10s %-15s\n" \
           "$timestamp" "$cpu_usage" "$memory_usage" "$backend_cpu" "$backend_mem" "$connections" "$response_time"
    
    # 记录到日志文件
    echo "$timestamp,$cpu_usage,$memory_usage,$backend_cpu,$backend_mem,$connections,$response_time" >> "$LOG_FILE"
}

# 生成性能报告
generate_report() {
    if [[ ! -f "$LOG_FILE" ]]; then
        echo -e "${RED}没有找到监控日志文件${NC}"
        return
    fi
    
    echo -e "\n${GREEN}=== 性能监控报告 ===${NC}"
    
    # 统计信息
    local total_records=$(wc -l < "$LOG_FILE")
    echo -e "${BLUE}总监控记录数: $total_records${NC}"
    
    if [[ $total_records -gt 0 ]]; then
        # CPU使用率统计
        local avg_cpu=$(awk -F',' '{sum+=$2; count++} END {if(count>0) printf "%.1f", sum/count}' "$LOG_FILE")
        local max_cpu=$(awk -F',' 'BEGIN{max=0} {if($2>max) max=$2} END {printf "%.1f", max}' "$LOG_FILE")
        
        # 内存使用率统计
        local avg_mem=$(awk -F',' '{sum+=$3; count++} END {if(count>0) printf "%.1f", sum/count}' "$LOG_FILE")
        local max_mem=$(awk -F',' 'BEGIN{max=0} {if($3>max) max=$3} END {printf "%.1f", max}' "$LOG_FILE")
        
        # 后端进程统计
        local avg_backend_cpu=$(awk -F',' '{sum+=$4; count++} END {if(count>0) printf "%.1f", sum/count}' "$LOG_FILE")
        local max_backend_cpu=$(awk -F',' 'BEGIN{max=0} {if($4>max) max=$4} END {printf "%.1f", max}' "$LOG_FILE")
        
        # 连接数统计
        local avg_connections=$(awk -F',' '{sum+=$6; count++} END {if(count>0) printf "%.0f", sum/count}' "$LOG_FILE")
        local max_connections=$(awk -F',' 'BEGIN{max=0} {if($6>max) max=$6} END {printf "%.0f", max}' "$LOG_FILE")
        
        # 响应时间统计（排除timeout）
        local avg_response=$(awk -F',' '$7 != "timeout" {sum+=$7; count++} END {if(count>0) printf "%.0f", sum/count}' "$LOG_FILE")
        local max_response=$(awk -F',' '$7 != "timeout" && $7 > max {max=$7} END {printf "%.0f", max}' "$LOG_FILE")
        
        echo ""
        echo -e "${YELLOW}系统资源使用情况:${NC}"
        echo "  CPU使用率    - 平均: ${avg_cpu}%, 最高: ${max_cpu}%"
        echo "  内存使用率   - 平均: ${avg_mem}%, 最高: ${max_mem}%"
        echo ""
        echo -e "${YELLOW}后端服务性能:${NC}"
        echo "  CPU使用率    - 平均: ${avg_backend_cpu}%, 最高: ${max_backend_cpu}%"
        echo "  并发连接数   - 平均: ${avg_connections}, 最高: ${max_connections}"
        echo "  响应时间     - 平均: ${avg_response}ms, 最高: ${max_response}ms"
        
        # 性能评估
        echo ""
        echo -e "${YELLOW}性能评估:${NC}"
        
        if (( $(echo "$avg_response < 100" | bc -l) )); then
            echo -e "  响应时间: ${GREEN}优秀${NC} (平均 ${avg_response}ms < 100ms)"
        elif (( $(echo "$avg_response < 200" | bc -l) )); then
            echo -e "  响应时间: ${YELLOW}良好${NC} (平均 ${avg_response}ms)"
        else
            echo -e "  响应时间: ${RED}需要优化${NC} (平均 ${avg_response}ms > 200ms)"
        fi
        
        if (( $(echo "$avg_cpu < 70" | bc -l) )); then
            echo -e "  CPU使用率: ${GREEN}正常${NC} (平均 ${avg_cpu}%)"
        elif (( $(echo "$avg_cpu < 90" | bc -l) )); then
            echo -e "  CPU使用率: ${YELLOW}较高${NC} (平均 ${avg_cpu}%)"
        else
            echo -e "  CPU使用率: ${RED}过高${NC} (平均 ${avg_cpu}%)"
        fi
        
        if (( $(echo "$avg_mem < 80" | bc -l) )); then
            echo -e "  内存使用率: ${GREEN}正常${NC} (平均 ${avg_mem}%)"
        else
            echo -e "  内存使用率: ${YELLOW}较高${NC} (平均 ${avg_mem}%)"
        fi
    fi
    
    echo ""
    echo -e "${BLUE}详细日志文件: $LOG_FILE${NC}"
}

# 清理函数
cleanup() {
    echo -e "\n${YELLOW}停止监控...${NC}"
    generate_report
    echo -e "${GREEN}监控已停止${NC}"
    exit 0
}

# 显示帮助信息
show_help() {
    echo -e "${BLUE}性能监控脚本${NC}"
    echo -e "${YELLOW}用法: $0 [选项]${NC}"
    echo ""
    echo "选项:"
    echo "  start            开始监控 (默认)"
    echo "  report           生成性能报告"
    echo "  clean            清理日志文件"
    echo "  help             显示帮助信息"
    echo ""
    echo "示例:"
    echo "  $0 start         # 开始实时监控"
    echo "  $0 report        # 查看性能报告"
    echo "  $0 clean         # 清理日志"
}

# 清理日志文件
clean_logs() {
    if [[ -f "$LOG_FILE" ]]; then
        rm "$LOG_FILE"
        echo -e "${GREEN}日志文件已清理${NC}"
    else
        echo -e "${YELLOW}没有找到日志文件${NC}"
    fi
}

# 主函数
main() {
    case "${1:-start}" in
        "start")
            # 设置信号处理
            trap cleanup SIGINT SIGTERM
            
            # 初始化日志文件
            echo "时间,CPU%,内存%,后端CPU%,后端内存%,连接数,响应时间(ms)" > "$LOG_FILE"
            
            # 显示监控界面
            show_monitor_header
            
            # 开始监控循环
            while true; do
                log_performance_data
                sleep $MONITOR_INTERVAL
            done
            ;;
        "report")
            generate_report
            ;;
        "clean")
            clean_logs
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

# 检查依赖
if ! command -v bc &> /dev/null; then
    echo -e "${YELLOW}警告: bc 命令未找到，某些计算功能可能不可用${NC}"
fi

# 执行主函数
main "$@"