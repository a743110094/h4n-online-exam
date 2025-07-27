# 异步考试流程测试文档

## 概述

本文档介绍如何测试在线考试系统的异步处理流程，验证以下关键功能：

```
学生提交考试 → 立即响应成功 → 发送异步消息 → 后台处理
                                    ↓
                            [成绩计算、统计更新、报告生成]
```

## 测试文件说明

### 1. `test_async_exam_flow.go`
**完整的单元测试套件**
- 使用 testify 框架
- 创建完整的测试数据
- 测试完整的考试流程
- 验证异步处理结果
- 性能测试

### 2. `run_async_test.go`
**简化的功能测试工具**
- 独立运行的测试程序
- 快速验证异步功能
- 并发性能测试
- 实时输出测试结果

## 环境准备

### 1. 启动依赖服务

```bash
# 启动 RabbitMQ 和数据库
docker-compose up -d

# 验证服务状态
docker-compose ps
```

### 2. 检查服务连接

```bash
# 检查 RabbitMQ 管理界面
open http://localhost:15672
# 用户名: admin, 密码: admin123

# 检查数据库连接
psql -h localhost -p 5432 -U postgres -d online_exam_system
```

### 3. 安装测试依赖

```bash
cd backend
go mod tidy
go get github.com/stretchr/testify/suite
go get github.com/stretchr/testify/assert
```

## 运行测试

### 方式1: 快速功能测试（推荐）

```bash
cd backend
go run run_async_test.go test
```

**输出示例：**
```
=== 开始异步流程测试 ===
消息消费者已启动

--- 测试1: 消息发布 ---
发布考试结果消息...
✅ 考试结果消息发布成功
发布成绩计算消息...
✅ 成绩计算消息发布成功
发布统计更新消息...
✅ 统计更新消息发布成功
发布报告生成消息...
✅ 报告生成消息发布成功

--- 测试2: 考试提交流程 ---
模拟学生提交考试...
✅ 考试提交成功，响应时间: 15.2ms
📤 异步消息已发送，后台正在处理...
✅ 响应时间优秀 (< 100ms)

--- 测试3: 并发提交 ---
模拟 20 个学生同时提交考试...
✅ 并发测试完成:
   - 总提交数: 20
   - 成功数: 20
   - 成功率: 100.0%
   - 总耗时: 45.8ms
   - 平均响应时间: 2.3ms
✅ 并发性能优秀

等待异步处理完成...
=== 异步流程测试完成 ===
```

### 方式2: 完整单元测试

```bash
cd backend
go test -v -run TestAsyncExamFlowTestSuite
```

### 方式3: 运行特定测试

```bash
# 只测试消息发布
go test -v -run TestAsyncMessagePublishing

# 只测试性能
go test -v -run TestAsyncProcessingPerformance

# 测试完整流程
go test -v -run TestCompleteAsyncExamFlow
```

## 测试验证点

### 1. 响应时间验证
- ✅ 考试提交响应时间 < 100ms
- ✅ 异步消息发送成功
- ✅ 用户立即收到成功响应

### 2. 异步处理验证
- ✅ 成绩计算正确性
- ✅ 统计数据更新
- ✅ 考试报告生成
- ✅ 错误处理和重试

### 3. 并发性能验证
- ✅ 多学生同时提交不阻塞
- ✅ 消息队列处理能力
- ✅ 系统稳定性

## 故障排查

### 1. RabbitMQ 连接失败

```bash
# 检查 RabbitMQ 服务状态
docker-compose logs rabbitmq

# 重启 RabbitMQ
docker-compose restart rabbitmq

# 检查端口占用
lsof -i :5672
lsof -i :15672
```

### 2. 数据库连接失败

```bash
# 检查数据库服务
docker-compose logs postgres

# 检查环境变量
cat .env

# 测试数据库连接
psql -h localhost -p 5432 -U postgres -c "SELECT 1;"
```

### 3. 消息处理失败

```bash
# 查看应用日志
go run main.go 2>&1 | grep -i error

# 检查队列状态（RabbitMQ 管理界面）
open http://localhost:15672/#/queues

# 查看消息积压情况
curl -u admin:admin123 http://localhost:15672/api/queues
```

### 4. 测试数据问题

```bash
# 清理测试数据
psql -h localhost -p 5432 -U postgres -d online_exam_system -c "
DELETE FROM exam_records WHERE exam_id = 999;
DELETE FROM answers WHERE exam_record_id IN (SELECT id FROM exam_records WHERE exam_id = 999);
"

# 重置序列
psql -h localhost -p 5432 -U postgres -d online_exam_system -c "
SELECT setval('exam_records_id_seq', (SELECT MAX(id) FROM exam_records));
"
```

## 性能基准

### 响应时间目标
- 考试提交响应: < 100ms
- 消息发布: < 10ms
- 并发处理: 100+ 请求/秒

### 异步处理时间
- 成绩计算: < 5秒
- 统计更新: < 3秒
- 报告生成: < 10秒

### 并发能力
- 同时提交: 50+ 学生
- 消息吞吐: 1000+ 消息/分钟
- 系统稳定性: 99.9%

## 监控和调优

### 1. RabbitMQ 监控

访问管理界面: http://localhost:15672

关键指标：
- 队列长度
- 消息处理速率
- 连接数
- 内存使用

### 2. 应用监控

```bash
# 查看 Go 程序性能
go tool pprof http://localhost:8080/debug/pprof/profile

# 内存使用情况
go tool pprof http://localhost:8080/debug/pprof/heap
```

### 3. 数据库监控

```sql
-- 查看活跃连接
SELECT * FROM pg_stat_activity WHERE state = 'active';

-- 查看慢查询
SELECT query, mean_time, calls FROM pg_stat_statements ORDER BY mean_time DESC LIMIT 10;
```

## 扩展测试

### 1. 压力测试

```bash
# 使用 k6 进行压力测试
cd ../load-test
./run-load-tests.sh
```

### 2. 长时间稳定性测试

```bash
# 运行 24 小时稳定性测试
for i in {1..1440}; do
  go run run_async_test.go test
  sleep 60
done
```

### 3. 故障恢复测试

```bash
# 模拟 RabbitMQ 故障
docker-compose stop rabbitmq
sleep 30
docker-compose start rabbitmq

# 验证消息恢复处理
go run run_async_test.go test
```

## 总结

通过这些测试，我们可以验证：

1. **用户体验**: 考试提交立即响应，不会因为后台处理而阻塞
2. **系统可靠性**: 异步消息处理机制工作正常
3. **性能表现**: 系统能够处理高并发考试提交
4. **数据一致性**: 成绩计算、统计更新、报告生成都能正确完成

这确保了在大量学生同时交卷的场景下，系统能够保持良好的响应性能和用户体验。