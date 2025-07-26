<template>
  <div class="monitoring-view">
    <!-- 页面标题 -->
    <div class="page-header">
      <h1 class="page-title">监控中心</h1>
      <p class="page-subtitle">实时监控考试进行状态和学生答题情况</p>
    </div>

    <!-- 监控统计卡片 -->
    <div class="monitoring-stats">
      <div class="stats-grid">
        <div class="stat-card">
          <div class="stat-icon online">
            <el-icon><User /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ onlineStudents }}</div>
            <div class="stat-label">在线学生</div>
          </div>
        </div>
        
        <div class="stat-card">
          <div class="stat-icon active">
            <el-icon><Document /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ activeExams }}</div>
            <div class="stat-label">进行中考试</div>
          </div>
        </div>
        
        <div class="stat-card">
          <div class="stat-icon warning">
            <el-icon><Warning /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ anomalies }}</div>
            <div class="stat-label">异常行为</div>
          </div>
        </div>
        
        <div class="stat-card">
          <div class="stat-icon success">
            <el-icon><CircleCheck /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ completedToday }}</div>
            <div class="stat-label">今日完成</div>
          </div>
        </div>
      </div>
    </div>

    <!-- 实时监控表格 -->
    <div class="monitoring-table">
      <div class="table-header">
        <h2>实时监控</h2>
        <div class="table-actions">
          <el-button type="primary" @click="refreshData">
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
        </div>
      </div>
      
      <el-table :data="monitoringData" style="width: 100%" stripe>
        <el-table-column prop="studentName" label="学生姓名" width="120" />
        <el-table-column prop="examTitle" label="考试名称" width="200" />
        <el-table-column prop="progress" label="答题进度" width="150">
          <template #default="scope">
            <el-progress :percentage="scope.row.progress" :color="getProgressColor(scope.row.progress)" />
          </template>
        </el-table-column>
        <el-table-column prop="timeRemaining" label="剩余时间" width="120" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="scope">
            <el-tag :type="getStatusType(scope.row.status)">{{ scope.row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="lastActivity" label="最后活动" width="150" />
        <el-table-column label="操作" width="150">
          <template #default="scope">
            <el-button size="small" @click="viewDetails(scope.row)">查看详情</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <!-- 异常行为监控 -->
    <div class="anomaly-section" v-if="anomalyList.length > 0">
      <h2>异常行为监控</h2>
      <div class="anomaly-list">
        <div v-for="anomaly in anomalyList" :key="anomaly.id" class="anomaly-item">
          <div class="anomaly-icon">
            <el-icon><Warning /></el-icon>
          </div>
          <div class="anomaly-content">
            <div class="anomaly-title">{{ anomaly.type }}</div>
            <div class="anomaly-desc">{{ anomaly.description }}</div>
            <div class="anomaly-time">{{ anomaly.time }}</div>
          </div>
          <div class="anomaly-actions">
            <el-button size="small" type="warning" @click="handleAnomaly(anomaly)">处理</el-button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { ElMessage } from 'element-plus'
import { User, Document, Warning, CircleCheck, Refresh } from '@element-plus/icons-vue'

// 监控统计数据
const onlineStudents = ref(45)
const activeExams = ref(3)
const anomalies = ref(2)
const completedToday = ref(28)

// 监控表格数据
const monitoringData = ref([
  {
    id: 1,
    studentName: '张三',
    examTitle: '高等数学期末考试',
    progress: 75,
    timeRemaining: '25分钟',
    status: '正常',
    lastActivity: '2分钟前'
  },
  {
    id: 2,
    studentName: '李四',
    examTitle: '英语四级模拟考试',
    progress: 45,
    timeRemaining: '45分钟',
    status: '正常',
    lastActivity: '1分钟前'
  },
  {
    id: 3,
    studentName: '王五',
    examTitle: '计算机基础测试',
    progress: 90,
    timeRemaining: '10分钟',
    status: '异常',
    lastActivity: '10分钟前'
  }
])

// 异常行为列表
const anomalyList = ref([
  {
    id: 1,
    type: '长时间无操作',
    description: '学生王五已10分钟未进行任何操作',
    time: '10分钟前'
  },
  {
    id: 2,
    type: '频繁切换窗口',
    description: '学生赵六在5分钟内切换窗口15次',
    time: '5分钟前'
  }
])

// 定时器
let refreshTimer: NodeJS.Timeout | null = null

// 获取进度条颜色
const getProgressColor = (progress: number) => {
  if (progress < 30) return '#f56c6c'
  if (progress < 70) return '#e6a23c'
  return '#67c23a'
}

// 获取状态标签类型
const getStatusType = (status: string) => {
  switch (status) {
    case '正常': return 'success'
    case '异常': return 'danger'
    case '离线': return 'info'
    default: return 'warning'
  }
}

// 刷新数据
const refreshData = () => {
  ElMessage.success('数据已刷新')
  // 这里可以调用API刷新数据
}

// 查看详情
const viewDetails = (row: any) => {
  ElMessage.info(`查看学生 ${row.studentName} 的详细信息`)
  // 这里可以打开详情弹窗或跳转到详情页面
}

// 处理异常
const handleAnomaly = (anomaly: any) => {
  ElMessage.warning(`处理异常: ${anomaly.type}`)
  // 这里可以处理异常行为
}

// 自动刷新数据
const startAutoRefresh = () => {
  refreshTimer = setInterval(() => {
    // 模拟数据更新
    onlineStudents.value = Math.floor(Math.random() * 10) + 40
  }, 30000) // 每30秒刷新一次
}

onMounted(() => {
  startAutoRefresh()
})

onUnmounted(() => {
  if (refreshTimer) {
    clearInterval(refreshTimer)
  }
})
</script>

<style scoped>
.monitoring-view {
  padding: var(--spacing-lg);
  background: var(--color-bg-soft);
  min-height: 100vh;
}

.page-header {
  margin-bottom: var(--spacing-xl);
}

.page-title {
  font-size: 28px;
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0 0 var(--spacing-sm) 0;
}

.page-subtitle {
  font-size: 16px;
  color: var(--color-text-secondary);
  margin: 0;
}

.monitoring-stats {
  margin-bottom: var(--spacing-xl);
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: var(--spacing-lg);
}

.stat-card {
  background: white;
  border-radius: var(--radius-lg);
  padding: var(--spacing-lg);
  box-shadow: var(--shadow-sm);
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
}

.stat-icon {
  width: 60px;
  height: 60px;
  border-radius: var(--radius-lg);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  color: white;
}

.stat-icon.online {
  background: linear-gradient(135deg, #67c23a, #85ce61);
}

.stat-icon.active {
  background: linear-gradient(135deg, #409eff, #66b1ff);
}

.stat-icon.warning {
  background: linear-gradient(135deg, #e6a23c, #ebb563);
}

.stat-icon.success {
  background: linear-gradient(135deg, #67c23a, #85ce61);
}

.stat-content {
  flex: 1;
}

.stat-value {
  font-size: 32px;
  font-weight: 600;
  color: var(--color-text-primary);
  line-height: 1;
}

.stat-label {
  font-size: 14px;
  color: var(--color-text-secondary);
  margin-top: var(--spacing-xs);
}

.monitoring-table {
  background: white;
  border-radius: var(--radius-lg);
  padding: var(--spacing-lg);
  box-shadow: var(--shadow-sm);
  margin-bottom: var(--spacing-xl);
}

.table-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--spacing-lg);
}

.table-header h2 {
  font-size: 20px;
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0;
}

.anomaly-section {
  background: white;
  border-radius: var(--radius-lg);
  padding: var(--spacing-lg);
  box-shadow: var(--shadow-sm);
}

.anomaly-section h2 {
  font-size: 20px;
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0 0 var(--spacing-lg) 0;
}

.anomaly-list {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
}

.anomaly-item {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
  padding: var(--spacing-md);
  background: var(--color-bg-soft);
  border-radius: var(--radius-md);
  border-left: 4px solid var(--color-warning);
}

.anomaly-icon {
  width: 40px;
  height: 40px;
  border-radius: var(--radius-md);
  background: var(--color-warning);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 18px;
}

.anomaly-content {
  flex: 1;
}

.anomaly-title {
  font-weight: 600;
  color: var(--color-text-primary);
  margin-bottom: var(--spacing-xs);
}

.anomaly-desc {
  color: var(--color-text-secondary);
  font-size: 14px;
  margin-bottom: var(--spacing-xs);
}

.anomaly-time {
  color: var(--color-text-tertiary);
  font-size: 12px;
}

@media (max-width: 768px) {
  .monitoring-view {
    padding: var(--spacing-md);
  }
  
  .stats-grid {
    grid-template-columns: 1fr;
  }
  
  .table-header {
    flex-direction: column;
    gap: var(--spacing-md);
    align-items: stretch;
  }
  
  .anomaly-item {
    flex-direction: column;
    align-items: flex-start;
  }
}
</style>