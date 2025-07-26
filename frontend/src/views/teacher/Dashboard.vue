<template>
  <div class="teacher-dashboard">
    <!-- 欢迎区域和快速操作 -->
    <div class="welcome-and-actions">
      <!-- 左侧欢迎卡片 -->
      <div class="welcome-section dopamine-card">
        <div class="welcome-content">
          <div class="welcome-text">
            <h1 class="welcome-title">
              欢迎回来，{{ authStore.user?.realName || '老师' }}！
            </h1>
            <p class="welcome-subtitle">
              今天是 {{ formatDate(new Date()) }}，开始您的教学工作吧
            </p>
          </div>
          <div class="welcome-stats">
            <div class="stat-item">
              <div class="stat-value">{{ teacherStats.totalStudents }}</div>
              <div class="stat-label">管理学生</div>
            </div>
            <div class="stat-item">
              <div class="stat-value">{{ teacherStats.totalExams }}</div>
              <div class="stat-label">创建考试</div>
            </div>
            <div class="stat-item">
              <div class="stat-value">{{ teacherStats.totalQuestions }}</div>
              <div class="stat-label">题库题目</div>
            </div>
          </div>
        </div>
        <div class="welcome-illustration">
          <div class="illustration-bg">
            <el-icon :size="80" class="illustration-icon">
              <School />
            </el-icon>
          </div>
        </div>
      </div>
      
      <!-- 右侧快速操作 -->
      <div class="quick-actions">
        <div class="section-header">
          <h2>快速操作</h2>
        </div>
        <div class="action-grid-2x2">
          <div
            v-for="action in quickActions"
            :key="action.key"
            class="action-card-compact dopamine-card"
            @click="handleQuickAction(action.key)"
          >
            <div class="action-icon-compact" :style="{ background: action.color }">
              <el-icon :size="16">
                <component :is="action.icon" />
              </el-icon>
            </div>
            <div class="action-content-compact">
              <div class="action-title-compact">{{ action.title }}</div>
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 进行中的考试 -->
    <div class="ongoing-exams">
      <div class="section-header">
        <h2>进行中的考试</h2>
        <el-button type="primary" @click="$router.push('/teacher/monitoring')">
          监控中心
        </el-button>
      </div>
      <div class="exam-list">
        <div
          v-for="exam in ongoingExams"
          :key="exam.id"
          class="exam-card dopamine-card"
        >
          <div class="exam-header">
            <div class="exam-subject">{{ exam.subject }}</div>
            <div class="exam-status status-progress">
              进行中
            </div>
          </div>
          <div class="exam-title">{{ exam.title }}</div>
          <div class="exam-stats">
            <div class="exam-stat">
              <el-icon><User /></el-icon>
              <span>{{ exam.participantCount }}/{{ exam.totalStudents }} 人参与</span>
            </div>
            <div class="exam-stat">
              <el-icon><Clock /></el-icon>
              <span>剩余 {{ exam.remainingTime }} 分钟</span>
            </div>
            <div class="exam-stat">
              <el-icon><Document /></el-icon>
              <span>{{ exam.submittedCount }} 人已提交</span>
            </div>
          </div>
          <div class="exam-progress">
            <div class="progress-label">
              <span>参与进度</span>
              <span>{{ Math.round((exam.participantCount / exam.totalStudents) * 100) }}%</span>
            </div>
            <el-progress
              :percentage="Math.round((exam.participantCount / exam.totalStudents) * 100)"
              :color="getProgressColor(Math.round((exam.participantCount / exam.totalStudents) * 100))"
            />
          </div>
          <div class="exam-actions">
            <el-button
              type="primary"
              size="small"
              @click="monitorExam(exam.id)"
            >
              实时监控
            </el-button>
            <el-button
              type="warning"
              size="small"
              @click="endExam(exam.id)"
            >
              结束考试
            </el-button>
          </div>
        </div>
        
        <div v-if="ongoingExams.length === 0" class="empty-state">
          <el-icon :size="48" class="empty-icon"><Clock /></el-icon>
          <p>暂无进行中的考试</p>
        </div>
      </div>
    </div>
    
    <!-- 数据统计 -->
    <div class="statistics-section">
      <div class="section-header">
        <h2>数据统计</h2>
      </div>
      <div class="stats-grid">
        <!-- 考试统计图表 -->
        <div class="stats-card dopamine-card">
          <h3>本月考试统计</h3>
          <div class="chart-container">
            <div ref="examStatsChart" class="chart"></div>
          </div>
        </div>
        
        <!-- 成绩分布 -->
        <div class="stats-card dopamine-card">
          <h3>最近考试成绩分布</h3>
          <div class="chart-container">
            <div ref="gradeDistributionChart" class="chart"></div>
          </div>
        </div>
        
        <!-- 题目统计 -->
        <div class="stats-card dopamine-card">
          <h3>题库统计</h3>
          <div class="question-stats">
            <div class="question-type" v-for="type in questionTypes" :key="type.type">
              <div class="type-header">
                <span class="type-name">{{ type.name }}</span>
                <span class="type-count">{{ type.count }}</span>
              </div>
              <el-progress
                :percentage="Math.round((type.count / totalQuestions) * 100)"
                :color="type.color"
                :stroke-width="8"
              />
            </div>
          </div>
        </div>
        
        <!-- 学生活跃度 -->
        <div class="stats-card dopamine-card">
          <h3>学生活跃度</h3>
          <div class="activity-stats">
            <div class="activity-item">
              <div class="activity-label">今日活跃</div>
              <div class="activity-value">{{ studentActivity.today }}</div>
            </div>
            <div class="activity-item">
              <div class="activity-label">本周活跃</div>
              <div class="activity-value">{{ studentActivity.thisWeek }}</div>
            </div>
            <div class="activity-item">
              <div class="activity-label">本月活跃</div>
              <div class="activity-value">{{ studentActivity.thisMonth }}</div>
            </div>
            <div class="activity-chart">
              <div ref="activityChart" class="chart"></div>
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 待办事项 -->
    <div class="todo-section">
      <div class="section-header">
        <h2>待办事项</h2>
        <el-button type="primary" size="small" @click="addTodo">
          <el-icon><Plus /></el-icon>
          添加
        </el-button>
      </div>
      <div class="todo-list">
        <div
          v-for="todo in todoList"
          :key="todo.id"
          class="todo-item dopamine-card"
          :class="{ completed: todo.completed }"
        >
          <el-checkbox
            v-model="todo.completed"
            @change="updateTodo(todo.id)"
          />
          <div class="todo-content">
            <div class="todo-title" :class="{ completed: todo.completed }">
              {{ todo.title }}
            </div>
            <div class="todo-desc">{{ todo.description }}</div>
            <div class="todo-meta">
              <span class="todo-priority" :class="`priority-${todo.priority}`">
                {{ getPriorityText(todo.priority) }}
              </span>
              <span class="todo-date">{{ formatDate(todo.dueDate) }}</span>
            </div>
          </div>
          <el-button
            type="danger"
            size="small"
            circle
            @click="deleteTodo(todo.id)"
          >
            <el-icon><Delete /></el-icon>
          </el-button>
        </div>
        
        <div v-if="todoList.length === 0" class="empty-state">
          <el-icon :size="48" class="empty-icon"><List /></el-icon>
          <p>暂无待办事项</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  School, ArrowRight, User, Clock, Document, Plus, Delete, List,
  EditPen, Files, Notebook, DataAnalysis
} from '@element-plus/icons-vue'
import { useAuthStore } from '@/stores/auth'
import * as echarts from 'echarts'

const router = useRouter()
const authStore = useAuthStore()

// 教师统计数据
const teacherStats = ref({
  totalStudents: 156,
  totalExams: 23,
  totalQuestions: 342
})

// 快速操作
const quickActions = ref([
  {
    key: 'create-question',
    title: '创建题目',
    description: '快速添加新的考试题目',
    icon: 'EditPen',
    color: 'var(--gradient-primary)'
  },
  {
    key: 'create-paper',
    title: '组卷',
    description: '创建新的试卷',
    icon: 'Files',
    color: 'var(--gradient-success)'
  },
  {
    key: 'create-exam',
    title: '安排考试',
    description: '创建新的考试安排',
    icon: 'Notebook',
    color: 'var(--gradient-warning)'
  },
  {
    key: 'view-stats',
    title: '查看统计',
    description: '查看详细的教学数据',
    icon: 'DataAnalysis',
    color: 'var(--gradient-danger)'
  }
])

// 进行中的考试
const ongoingExams = ref([
  {
    id: 1,
    title: '数据结构期中考试',
    subject: '数据结构',
    participantCount: 45,
    totalStudents: 50,
    submittedCount: 12,
    remainingTime: 45
  },
  {
    id: 2,
    title: '算法设计随堂测验',
    subject: '算法设计',
    participantCount: 28,
    totalStudents: 30,
    submittedCount: 8,
    remainingTime: 15
  }
])

// 题目类型统计
const questionTypes = ref([
  { type: 'single', name: '单选题', count: 156, color: '#4ECDC4' },
  { type: 'multiple', name: '多选题', count: 89, color: '#FF6B8A' },
  { type: 'judge', name: '判断题', count: 67, color: '#FFD93D' },
  { type: 'fill', name: '填空题', count: 30, color: '#6BCF7F' }
])

const totalQuestions = ref(342)

// 学生活跃度
const studentActivity = ref({
  today: 89,
  thisWeek: 134,
  thisMonth: 156
})

// 待办事项
const todoList = ref([
  {
    id: 1,
    title: '批改期中考试试卷',
    description: '数据结构期中考试，共50份试卷需要批改',
    priority: 'high',
    dueDate: new Date(Date.now() + 2 * 24 * 60 * 60 * 1000),
    completed: false
  },
  {
    id: 2,
    title: '准备下周课程内容',
    description: '算法设计第8章：动态规划',
    priority: 'medium',
    dueDate: new Date(Date.now() + 5 * 24 * 60 * 60 * 1000),
    completed: false
  },
  {
    id: 3,
    title: '更新题库',
    description: '添加新的算法题目到题库中',
    priority: 'low',
    dueDate: new Date(Date.now() + 7 * 24 * 60 * 60 * 1000),
    completed: true
  }
])

// 图表引用
const examStatsChart = ref<HTMLElement>()
const gradeDistributionChart = ref<HTMLElement>()
const activityChart = ref<HTMLElement>()

// 格式化日期
const formatDate = (date: Date): string => {
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
}

// 获取进度颜色
const getProgressColor = (progress: number): string => {
  if (progress >= 80) return '#67C23A'
  if (progress >= 60) return '#E6A23C'
  return '#F56C6C'
}

// 获取优先级文本
const getPriorityText = (priority: string): string => {
  const textMap: Record<string, string> = {
    high: '高优先级',
    medium: '中优先级',
    low: '低优先级'
  }
  return textMap[priority] || ''
}

// 处理快速操作
const handleQuickAction = (key: string) => {
  switch (key) {
    case 'create-question':
      router.push('/teacher/questions')
      break
    case 'create-paper':
      router.push('/teacher/papers')
      break
    case 'create-exam':
      router.push('/teacher/exams')
      break
    case 'view-stats':
      ElMessage.info('查看统计功能开发中')
      break
  }
}

// 监控考试
const monitorExam = (examId: number) => {
  router.push(`/teacher/monitoring?examId=${examId}`)
}

// 结束考试
const endExam = async (examId: number) => {
  try {
    await ElMessageBox.confirm(
      '确定要结束这场考试吗？结束后学生将无法继续答题。',
      '结束考试确认',
      {
        confirmButtonText: '确定结束',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    ElMessage.success('考试已结束')
    // 这里调用API结束考试
  } catch {
    // 用户取消
  }
}

// 添加待办事项
const addTodo = () => {
  ElMessage.info('添加待办事项功能开发中')
}

// 更新待办事项
const updateTodo = (todoId: number) => {
  const todo = todoList.value.find(t => t.id === todoId)
  if (todo) {
    ElMessage.success(todo.completed ? '任务已完成' : '任务已标记为未完成')
  }
}

// 删除待办事项
const deleteTodo = async (todoId: number) => {
  try {
    await ElMessageBox.confirm(
      '确定要删除这个待办事项吗？',
      '删除确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    const index = todoList.value.findIndex(t => t.id === todoId)
    if (index > -1) {
      todoList.value.splice(index, 1)
      ElMessage.success('待办事项已删除')
    }
  } catch {
    // 用户取消
  }
}

// 初始化考试统计图表
const initExamStatsChart = () => {
  if (!examStatsChart.value) return
  
  const chart = echarts.init(examStatsChart.value)
  const option = {
    tooltip: {
      trigger: 'axis'
    },
    legend: {
      data: ['创建考试', '参与人数']
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      data: ['第1周', '第2周', '第3周', '第4周']
    },
    yAxis: {
      type: 'value'
    },
    series: [
      {
        name: '创建考试',
        type: 'bar',
        data: [3, 5, 2, 4],
        itemStyle: {
          color: '#4ECDC4'
        }
      },
      {
        name: '参与人数',
        type: 'line',
        data: [120, 180, 90, 150],
        itemStyle: {
          color: '#FF6B8A'
        }
      }
    ]
  }
  
  chart.setOption(option)
  
  window.addEventListener('resize', () => {
    chart.resize()
  })
}

// 初始化成绩分布图表
const initGradeDistributionChart = () => {
  if (!gradeDistributionChart.value) return
  
  const chart = echarts.init(gradeDistributionChart.value)
  const option = {
    tooltip: {
      trigger: 'item'
    },
    series: [
      {
        name: '成绩分布',
        type: 'pie',
        radius: ['40%', '70%'],
        data: [
          { value: 15, name: '优秀(90-100)' },
          { value: 25, name: '良好(80-89)' },
          { value: 35, name: '中等(70-79)' },
          { value: 20, name: '及格(60-69)' },
          { value: 5, name: '不及格(<60)' }
        ],
        emphasis: {
          itemStyle: {
            shadowBlur: 10,
            shadowOffsetX: 0,
            shadowColor: 'rgba(0, 0, 0, 0.5)'
          }
        }
      }
    ]
  }
  
  chart.setOption(option)
  
  window.addEventListener('resize', () => {
    chart.resize()
  })
}

// 初始化活跃度图表
const initActivityChart = () => {
  if (!activityChart.value) return
  
  const chart = echarts.init(activityChart.value)
  const option = {
    tooltip: {
      trigger: 'axis'
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      data: ['周一', '周二', '周三', '周四', '周五', '周六', '周日']
    },
    yAxis: {
      type: 'value'
    },
    series: [
      {
        name: '活跃人数',
        type: 'line',
        smooth: true,
        data: [45, 52, 48, 61, 55, 38, 42],
        itemStyle: {
          color: '#4ECDC4'
        },
        areaStyle: {
          color: 'rgba(78, 205, 196, 0.3)'
        }
      }
    ]
  }
  
  chart.setOption(option)
  
  window.addEventListener('resize', () => {
    chart.resize()
  })
}

// 组件挂载后初始化
onMounted(async () => {
  await nextTick()
  initExamStatsChart()
  initGradeDistributionChart()
  initActivityChart()
})
</script>

<style scoped>
.teacher-dashboard {
  padding: var(--spacing-md);
}

.welcome-section {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--spacing-lg);
  margin-bottom: var(--spacing-lg);
  background: var(--gradient-success);
  color: white;
  border-radius: var(--radius-xl);
  overflow: hidden;
  position: relative;
}

.welcome-content {
  flex: 1;
  z-index: 2;
}

.welcome-title {
  font-size: 28px;
  font-weight: 700;
  margin: 0 0 var(--spacing-sm) 0;
}

.welcome-subtitle {
  font-size: 16px;
  opacity: 0.9;
  margin: 0 0 var(--spacing-lg) 0;
}

.welcome-stats {
  display: flex;
  gap: var(--spacing-lg);
}

.stat-item {
  text-align: center;
}

.stat-value {
  font-size: 24px;
  font-weight: 700;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 14px;
  opacity: 0.8;
}

.welcome-illustration {
  position: relative;
  z-index: 1;
}

.illustration-bg {
  width: 120px;
  height: 120px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  backdrop-filter: blur(10px);
}

.illustration-icon {
  color: rgba(255, 255, 255, 0.8);
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--spacing-lg);
}

.section-header h2 {
  font-size: 20px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0;
}

/* 欢迎区域和快速操作合并布局 */
.welcome-and-actions {
  display: flex;
  gap: var(--spacing-lg);
  margin-bottom: var(--spacing-lg);
  align-items: stretch;
}

.welcome-and-actions .welcome-section {
  flex: 2;
  margin-bottom: 0;
}

.welcome-and-actions .quick-actions {
  flex: 1;
  margin-bottom: 0;
}

.action-grid-2x2 {
  display: grid;
  grid-template-columns: 1fr 1fr;
  grid-template-rows: 1fr 1fr;
  gap: var(--spacing-sm);
  height: fit-content;
  max-height: 200px;
}

.action-card-compact {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: var(--spacing-sm);
  cursor: pointer;
  transition: all 0.3s ease;
  border: 2px solid transparent;
  text-align: center;
  min-height: 80px;
  max-height: 90px;
}

.action-card-compact:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow-xl);
  border-color: var(--business-blue);
}

.action-icon-compact {
  width: 32px;
  height: 32px;
  border-radius: var(--radius-md);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  margin-bottom: var(--spacing-xs);
}

.action-content-compact {
  text-align: center;
}

.action-title-compact {
  font-size: 12px;
  font-weight: 600;
  color: var(--text-primary);
  line-height: 1.2;
  text-align: center;
}

.quick-actions {
  margin-bottom: var(--spacing-lg);
}

.action-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: var(--spacing-lg);
}

.action-card {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
  padding: var(--spacing-lg);
  cursor: pointer;
  transition: all 0.3s ease;
  border: 2px solid transparent;
}

.action-card:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow-xl);
  border-color: var(--dopamine-green);
}

.action-icon {
  width: 48px;
  height: 48px;
  border-radius: var(--radius-lg);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  flex-shrink: 0;
}

.action-content {
  flex: 1;
}

.action-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 4px;
}

.action-desc {
  font-size: 14px;
  color: var(--text-secondary);
}

.action-arrow {
  color: var(--text-muted);
  transition: all 0.3s ease;
}

.action-card:hover .action-arrow {
  color: var(--dopamine-green);
  transform: translateX(4px);
}

.ongoing-exams {
  margin-bottom: var(--spacing-lg);
}

.exam-list {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(400px, 1fr));
  gap: var(--spacing-lg);
}

.exam-card {
  padding: var(--spacing-lg);
  transition: all 0.3s ease;
}

.exam-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-lg);
}

.exam-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--spacing-sm);
}

.exam-subject {
  font-size: 12px;
  font-weight: 600;
  color: var(--dopamine-green);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.exam-status {
  font-size: 12px;
  font-weight: 600;
  padding: 4px 8px;
  border-radius: var(--radius-sm);
}

.status-progress {
  background: rgba(230, 162, 60, 0.1);
  color: var(--dopamine-orange);
}

.exam-title {
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: var(--spacing-md);
}

.exam-stats {
  display: flex;
  flex-wrap: wrap;
  gap: var(--spacing-md);
  margin-bottom: var(--spacing-md);
}

.exam-stat {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 14px;
  color: var(--text-secondary);
}

.exam-progress {
  margin-bottom: var(--spacing-md);
}

.progress-label {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--spacing-sm);
  font-size: 14px;
  color: var(--text-secondary);
}

.exam-actions {
  display: flex;
  gap: var(--spacing-sm);
  justify-content: flex-end;
}

.empty-state {
  grid-column: 1 / -1;
  text-align: center;
  padding: var(--spacing-xl);
  color: var(--text-muted);
}

.empty-icon {
  margin-bottom: var(--spacing-md);
  opacity: 0.5;
}

.statistics-section {
  margin-bottom: var(--spacing-lg);
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: var(--spacing-lg);
}

.stats-card {
  padding: var(--spacing-lg);
}

.stats-card h3 {
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 var(--spacing-md) 0;
}

.chart-container {
  height: 200px;
}

.chart {
  width: 100%;
  height: 100%;
}

.question-stats {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
}

.question-type {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-sm);
}

.type-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.type-name {
  font-size: 14px;
  font-weight: 500;
  color: var(--text-primary);
}

.type-count {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-secondary);
}

.activity-stats {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
}

.activity-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--spacing-sm) 0;
  border-bottom: 1px solid var(--border-light);
}

.activity-item:last-of-type {
  border-bottom: none;
}

.activity-label {
  font-size: 14px;
  color: var(--text-secondary);
}

.activity-value {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
}

.activity-chart {
  height: 120px;
  margin-top: var(--spacing-md);
}

.todo-section {
  margin-bottom: var(--spacing-lg);
}

.todo-list {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
}

.todo-item {
  display: flex;
  align-items: flex-start;
  gap: var(--spacing-md);
  padding: var(--spacing-lg);
  transition: all 0.3s ease;
}

.todo-item.completed {
  opacity: 0.6;
}

.todo-content {
  flex: 1;
}

.todo-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 4px;
  transition: all 0.3s ease;
}

.todo-title.completed {
  text-decoration: line-through;
  color: var(--text-muted);
}

.todo-desc {
  font-size: 14px;
  color: var(--text-secondary);
  margin-bottom: var(--spacing-sm);
}

.todo-meta {
  display: flex;
  gap: var(--spacing-md);
  align-items: center;
}

.todo-priority {
  font-size: 12px;
  font-weight: 600;
  padding: 2px 6px;
  border-radius: var(--radius-sm);
}

.priority-high {
  background: rgba(245, 108, 108, 0.1);
  color: var(--dopamine-red);
}

.priority-medium {
  background: rgba(230, 162, 60, 0.1);
  color: var(--dopamine-orange);
}

.priority-low {
  background: rgba(103, 194, 58, 0.1);
  color: var(--dopamine-green);
}

.todo-date {
  font-size: 12px;
  color: var(--text-muted);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .welcome-section {
    flex-direction: column;
    text-align: center;
    gap: var(--spacing-lg);
  }
  
  .welcome-stats {
    justify-content: center;
  }
  
  .welcome-and-actions {
    flex-direction: column;
  }
  
  .welcome-and-actions .welcome-section {
    flex: none;
  }
  
  .welcome-and-actions .quick-actions {
    flex: none;
  }
  
  .action-grid-2x2 {
    grid-template-columns: 1fr 1fr;
    height: auto;
  }
  
  .action-grid {
    grid-template-columns: 1fr;
  }
  
  .exam-list {
    grid-template-columns: 1fr;
  }
  
  .stats-grid {
    grid-template-columns: 1fr;
  }
}
</style>