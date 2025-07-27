# RabbitMQ 异步处理文档

## 概述

本项目使用 RabbitMQ 实现考试结束后的异步处理，解决大量学生同时交卷导致的系统响应缓慢问题。

## 架构设计

### 消息队列

项目中定义了以下四个消息队列：

1. **exam_results** - 考试结果队列
   - 用途：接收考试提交事件
   - 消息格式：`ExamResultMessage`

2. **score_calculation** - 成绩计算队列
   - 用途：处理成绩计算任务
   - 消息格式：`ExamResultMessage`

3. **stats_update** - 统计更新队列
   - 用途：更新用户和考试统计数据
   - 消息格式：`ExamResultMessage`

4. **report_generation** - 报告生成队列
   - 用途：生成考试报告
   - 消息格式：`ExamResultMessage`

### 消息结构

```go
type ExamResultMessage struct {
    ExamRecordID uint      `json:"exam_record_id"`
    ExamID       uint      `json:"exam_id"`
    UserID       uint      `json:"user_id"`
    TenantID     string    `json:"tenant_id"`
    SubmitTime   time.Time `json:"submit_time"`
}
```

## 异步处理流程

### 1. 考试提交流程

```
学生提交试卷 → 更新考试记录状态 → 发送异步消息 → 立即返回成功响应
```

**详细步骤：**
1. 学生调用 `SubmitExam` API 提交试卷
2. 系统保存学生答案到数据库
3. 更新考试记录状态为 `ExamCompleted`
4. 发送 `ExamResultMessage` 到 `exam_results` 队列
5. 立即返回"试卷提交成功，成绩正在计算中"响应

### 2. 异步处理流程

```
exam_results 队列 → 分发到三个处理队列 → 并行处理
                  ↓
            score_calculation (成绩计算)
            stats_update (统计更新)
            report_generation (报告生成)
```

**处理步骤：**

#### 成绩计算 (score_calculation)
1. 获取考试记录和答案
2. 计算总分和正确题目数
3. 更新考试记录的成绩信息
4. 保存计算结果到数据库

#### 统计更新 (stats_update)
1. 更新用户统计信息 (`UserStats`)
   - 参加考试次数
   - 平均分数
   - 最高分数
   - 总学习时间
2. 更新考试统计信息 (`ExamStats`)
   - 参与人数
   - 平均分数
   - 最高分数
   - 及格率

#### 报告生成 (report_generation)
1. 生成个人考试报告 (`ExamReport`)
2. 包含详细的答题分析
3. 保存报告到数据库

## 部署配置

### Docker Compose 配置

```yaml
rabbitmq:
  image: rabbitmq:3.12-management
  container_name: exam_rabbitmq
  environment:
    RABBITMQ_DEFAULT_USER: admin
    RABBITMQ_DEFAULT_PASS: admin123
  ports:
    - "5672:5672"    # AMQP 端口
    - "15672:15672"  # 管理界面端口
  volumes:
    - rabbitmq_data:/var/lib/rabbitmq
  networks:
    - exam_network
  restart: unless-stopped
  healthcheck:
    test: ["CMD", "rabbitmq-diagnostics", "ping"]
    interval: 30s
    timeout: 10s
    retries: 3
```

### 环境变量配置

后端服务需要配置以下环境变量：

```env
RABBITMQ_HOST=rabbitmq
RABBITMQ_PORT=5672
RABBITMQ_USER=admin
RABBITMQ_PASSWORD=admin123
RABBITMQ_VHOST=/
```

## 监控和管理

### RabbitMQ 管理界面

- 访问地址：http://localhost:15672
- 用户名：admin
- 密码：admin123

### 队列监控指标

1. **消息积压数量** - 监控队列中未处理的消息数
2. **消息处理速率** - 每秒处理的消息数量
3. **错误率** - 处理失败的消息比例
4. **平均处理时间** - 消息从发送到处理完成的平均时间

## 错误处理

### 消息处理失败

1. **重试机制**：失败的消息会自动重试最多 3 次
2. **死信队列**：重试失败的消息会进入死信队列
3. **错误日志**：所有处理错误都会记录到日志中

### 服务降级

当 RabbitMQ 服务不可用时：
1. 考试提交仍然正常工作
2. 成绩计算会回退到同步处理
3. 系统会记录错误日志但不影响用户体验

## 性能优化

### 消息批处理

- 消费者使用批量处理模式，一次处理多条消息
- 减少数据库连接开销
- 提高整体处理效率

### 并发处理

- 每个队列启动多个消费者协程
- 支持水平扩展
- 根据负载动态调整消费者数量

### 数据库优化

- 使用事务批量更新数据
- 优化数据库索引
- 使用连接池管理数据库连接

## 故障排查

### 常见问题

1. **消息积压**
   - 检查消费者是否正常运行
   - 增加消费者数量
   - 检查数据库性能

2. **连接失败**
   - 检查 RabbitMQ 服务状态
   - 验证网络连接
   - 检查认证信息

3. **处理超时**
   - 调整消息超时时间
   - 优化处理逻辑
   - 检查数据库查询性能

### 日志分析

关键日志位置：
- RabbitMQ 连接日志：`/var/log/rabbitmq/`
- 应用处理日志：应用程序标准输出
- 错误日志：应用程序错误输出

## 扩展计划

### 未来优化方向

1. **消息持久化**：确保消息不丢失
2. **集群部署**：RabbitMQ 集群提高可用性
3. **监控告警**：集成 Prometheus + Grafana 监控
4. **自动扩缩容**：根据队列长度自动调整消费者数量
5. **消息路由优化**：使用 Exchange 和 Routing Key 优化消息分发