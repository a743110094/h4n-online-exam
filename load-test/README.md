# 在线考试系统压力测试

本目录包含了针对在线考试系统的完整压力测试套件，确保系统能够承受50个虚拟用户(VU)的并发访问，且95%的响应时间不超过100ms。

## 📋 测试目标

- **并发用户数**: 50 VU
- **95%响应时间**: < 100ms
- **错误率**: < 5%
- **测试工具**: k6

## 🚀 快速开始

### 1. 安装依赖

#### macOS
```bash
brew install k6
```

#### Ubuntu/Debian
```bash
sudo apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys C5AD17C747E3415A3642D57D77C6C491D6AC1D69
echo "deb https://dl.k6.io/deb stable main" | sudo tee /etc/apt/sources.list.d/k6.list
sudo apt-get update && sudo apt-get install k6
```

#### Windows
```powershell
choco install k6
```

### 2. 启动后端服务

确保后端服务在 `http://localhost:8080` 运行：

```bash
cd ../backend
go run main.go
```

### 3. 运行压力测试

```bash
# 给脚本执行权限
chmod +x run-load-tests.sh

# 运行所有测试
./run-load-tests.sh all

# 或运行特定测试
./run-load-tests.sh quick      # 快速验证
./run-load-tests.sh core       # 核心接口测试
./run-load-tests.sh admin      # 管理接口测试
```

## 📁 文件说明

### 测试脚本

| 文件 | 描述 | 目标用户 | 并发数 |
|------|------|----------|--------|
| `k6-load-test.js` | 综合API压力测试 | 所有角色 | 50 VU |
| `k6-core-apis-test.js` | 核心业务接口测试 | 主要是学生 | 50 VU |
| `k6-admin-teacher-test.js` | 管理接口测试 | 教师/管理员 | 15 VU |
| `run-load-tests.sh` | 测试执行脚本 | - | - |

### 测试覆盖的接口

#### 🎓 学生接口 (高频)
- `POST /api/v1/auth/login` - 用户登录
- `GET /api/v1/practice/recommendations` - 获取练习推荐
- `POST /api/v1/practice/start` - 开始练习
- `POST /api/v1/practice/{id}/answer` - 提交练习答案
- `GET /api/v1/practice/history` - 练习历史
- `GET /api/v1/practice/stats` - 练习统计
- `GET /api/v1/exams/student` - 学生考试列表
- `POST /api/v1/exams/{id}/start` - 开始考试
- `GET /api/v1/user/profile` - 用户资料

#### 👨‍🏫 教师接口 (中频)
- `GET /api/v1/questions` - 题目列表
- `POST /api/v1/teacher/questions` - 创建题目
- `GET /api/v1/papers` - 试卷列表
- `POST /api/v1/teacher/papers/auto` - 自动组卷
- `GET /api/v1/exams` - 考试列表
- `POST /api/v1/teacher/exams` - 创建考试
- `GET /api/v1/stats/teacher` - 教师统计

#### 👑 管理员接口 (低频)
- `GET /api/v1/admin/dashboard` - 管理员仪表板
- `GET /api/v1/admin/users` - 用户管理
- `GET /api/v1/subjects/all` - 科目管理

#### 🔧 公共接口
- `GET /api/v1/health` - 健康检查
- `GET /api/v1/subjects` - 科目列表

## 📊 性能指标

### 关键指标

| 指标 | 目标值 | 说明 |
|------|--------|------|
| `http_req_duration p(95)` | < 100ms | 95%请求响应时间 |
| `http_req_failed rate` | < 5% | 请求失败率 |
| `http_reqs` | - | 总请求数 |
| `vus` | 50 | 虚拟用户数 |
| `vus_max` | 50 | 最大虚拟用户数 |

### 自定义指标

- `errors` - 自定义错误率
- `practice_start_errors` - 练习开始错误率
- `exam_start_errors` - 考试开始错误率
- `create_question_errors` - 创建题目错误率
- `create_exam_errors` - 创建考试错误率

## 🧪 测试场景

### 1. 综合压力测试 (`k6-load-test.js`)

**测试阶段**:
- 预热: 2分钟内增加到10用户
- 压力: 5分钟内增加到50用户
- 稳定: 保持50用户10分钟
- 冷却: 2分钟内减少到0用户

**测试内容**: 覆盖所有角色的常用接口

### 2. 核心业务测试 (`k6-core-apis-test.js`)

**测试阶段**:
- 快速预热: 1分钟到20用户
- 压力增加: 3分钟到50用户
- 高强度: 保持50用户15分钟
- 逐步降低: 3分钟到0用户

**测试内容**: 重点测试学生练习和考试流程

### 3. 管理接口测试 (`k6-admin-teacher-test.js`)

**测试阶段**:
- 预热: 1分钟到5用户
- 增加: 3分钟到15用户
- 稳定: 保持15用户10分钟
- 结束: 2分钟到0用户

**测试内容**: 教师和管理员的管理操作

## 📈 测试报告

测试完成后，报告将保存在 `reports/YYYYMMDD_HHMMSS/` 目录下：

```
reports/
└── 20241201_143022/
    ├── comprehensive-test.json    # 综合测试详细数据
    ├── comprehensive-test.csv     # 综合测试CSV格式
    ├── core-apis-test.json        # 核心接口测试数据
    ├── core-apis-test.csv         # 核心接口测试CSV
    ├── admin-teacher-test.json    # 管理接口测试数据
    ├── admin-teacher-test.csv     # 管理接口测试CSV
    └── test-summary.txt           # 测试摘要
```

## 🔍 结果分析

### 查看实时结果

```bash
# 运行测试时会显示实时统计
k6 run k6-load-test.js
```

### 分析JSON报告

```bash
# 使用jq分析JSON报告
cat reports/*/comprehensive-test.json | jq '.metrics.http_req_duration'
```

### 关键检查点

1. **响应时间**: 确保95%请求 < 100ms
2. **错误率**: 确保 < 5%
3. **吞吐量**: 检查每秒请求数(RPS)
4. **资源使用**: 监控CPU、内存使用情况

## 🛠️ 故障排除

### 常见问题

1. **连接被拒绝**
   ```
   错误: dial tcp 127.0.0.1:8080: connect: connection refused
   ```
   **解决**: 确保后端服务正在运行

2. **认证失败**
   ```
   错误: login status is not 200
   ```
   **解决**: 检查测试用户是否存在，密码是否正确

3. **响应时间过长**
   ```
   阈值失败: http_req_duration p(95) > 100ms
   ```
   **解决**: 检查数据库索引、服务器资源

### 调试技巧

1. **增加日志输出**:
   ```javascript
   console.log('Response:', response.body);
   ```

2. **降低并发数**:
   ```javascript
   stages: [{ duration: '1m', target: 10 }]
   ```

3. **单独测试接口**:
   ```bash
   k6 run --vus 1 --duration 30s k6-load-test.js
   ```

## 📝 自定义测试

### 修改测试参数

```javascript
// 在测试文件中修改
export const options = {
  stages: [
    { duration: '2m', target: 20 }, // 修改目标用户数
  ],
  thresholds: {
    http_req_duration: ['p(95)<150'], // 修改响应时间阈值
  },
};
```

### 添加新的测试接口

```javascript
function testNewAPI(token) {
  const params = { headers: getAuthHeaders(token) };
  const response = http.get(`${BASE_URL}/new-endpoint`, params);
  
  check(response, {
    'new API status 200': (r) => r.status === 200,
    'new API RT < 100ms': (r) => r.timings.duration < 100,
  });
}
```

## 🎯 性能优化建议

基于测试结果，可能的优化方向：

1. **数据库优化**
   - 添加适当的索引
   - 优化查询语句
   - 使用连接池

2. **缓存策略**
   - Redis缓存热点数据
   - 静态资源CDN
   - 接口响应缓存

3. **代码优化**
   - 减少数据库查询次数
   - 异步处理非关键操作
   - 优化算法复杂度

4. **基础设施**
   - 增加服务器资源
   - 负载均衡
   - 数据库读写分离

## 📞 支持

如有问题，请检查：
1. 后端服务是否正常运行
2. 测试用户数据是否正确
3. 网络连接是否稳定
4. k6版本是否最新

---

**注意**: 请在测试环境中运行压力测试，避免对生产环境造成影响。