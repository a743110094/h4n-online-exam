<template>
  <div class="student-dashboard">
    <!-- 欢迎区域 -->
    <div class="welcome-section dopamine-card">
      <div class="welcome-content">
        <div class="welcome-text">
          <h1 class="welcome-title">
            欢迎回来，{{ authStore.user?.realName || '同学' }}！
          </h1>
          <p class="welcome-subtitle">
            今天是 {{ formatDate(new Date()) }}，继续你的学习之旅吧
          </p>
        </div>
        <div class="welcome-stats">
          <div class="stat-item">
            <div class="stat-value">{{ studentStats.totalExams }}</div>
            <div class="stat-label">已参加考试</div>
          </div>
          <div class="stat-item">
            <div class="stat-value">{{ studentStats.averageScore }}</div>
            <div class="stat-label">平均分数</div>
          </div>
          <div class="stat-item">
            <div class="stat-value">{{ studentStats.practiceCount }}</div>
            <div class="stat-label">练习次数</div>
          </div>
        </div>
      </div>
      <div class="welcome-illustration">
        <div class="illustration-bg">
          <el-icon :size="80" class="illustration-icon">
            <Reading />
          </el-icon>
        </div>
      </div>
    </div>
    
    <!-- 快速操作 -->
    <div class="quick-actions">
      <div class="section-header">
        <h2>快速操作</h2>
      </div>
      <div class="action-grid">
        <div
          v-for="action in quickActions"
          :key="action.key"
          class="action-card dopamine-card"
          @click="handleQuickAction(action.key)"
        >
          <div class="action-icon" :style="{ background: action.color }">
            <el-icon :size="24">
              <component :is="action.icon" />
            </el-icon>
          </div>
          <div class="action-content">
            <div class="action-title">{{ action.title }}</div>
            <div class="action-desc">{{ action.description }}</div>
          </div>
          <el-icon class="action-arrow"><ArrowRight /></el-icon>
        </div>
      </div>
    </div>
    
    <!-- 即将到来的考试 -->
    <div class="upcoming-exams">
      <div class="section-header">
        <h2>即将到来的考试</h2>
        <el-button type="primary" @click="$router.push('/student/exams')">
          查看全部
        </el-button>
      </div>
      <div class="exam-list">
        <div
          v-for="exam in upcomingExams"
          :key="exam.id"
          class="exam-card dopamine-card"
        >
          <div class="exam-header">
            <div class="exam-subject">{{ exam.subject }}</div>
            <div class="exam-status" :class="getExamStatusClass(exam.status)">
              {{ getExamStatusText(exam.status) }}
            </div>
          </div>
          <div class="exam-title">{{ exam.title }}</div>
          <div class="exam-details">
            <div class="exam-detail">
              <el-icon><Clock /></el-icon>
              <span>{{ formatDateTime(exam.startTime) }}</span>
            </div>
            <div class="exam-detail">
              <el-icon><Timer /></el-icon>
              <span>{{ exam.duration }}分钟</span>
            </div>
            <div class="exam-detail">
              <el-icon><Document /></el-icon>
              <span>{{ exam.questionCount }}题</span>
            </div>
          </div>
          <div class="exam-actions">
            <el-button
              v-if="exam.status === 'upcoming'"
              type="primary"
              size="small"
              @click="enterExam(exam.id)"
            >
              进入考试
            </el-button>
            <el-button
              v-else-if="exam.status === 'in-progress'"
              type="warning"
              size="small"
              @click="continueExam(exam.id)"
            >
              继续考试
            </el-button>
            <el-button
              v-else
              type="info"
              size="small"
              @click="viewResult(exam.id)"
            >
              查看结果
            </el-button>
          </div>
        </div>
        
        <div v-if="upcomingExams.length === 0" class="empty-state">
          <el-icon :size="48" class="empty-icon"><Calendar /></el-icon>
          <p>暂无即将到来的考试</p>
        </div>
      </div>
    </div>
    
    <!-- 学习进度 -->
    <div class="learning-progress">
      <div class="section-header">
        <h2>学习进度</h2>
      </div>
      <div class="progress-grid">
        <!-- 学习统计 -->
        <div class="progress-card dopamine-card">
          <h3>本月学习统计</h3>
          <div class="progress-chart">
            <div ref="learningChart" class="chart"></div>
          </div>
        </div>
        
        <!-- 科目进度 -->
        <div class="progress-card dopamine-card">
          <h3>科目掌握度</h3>
          <div class="subject-progress">
            <div
              v-for="subject in subjectProgress"
              :key="subject.id"
              class="subject-item"
            >
              <div class="subject-info">
                <span class="subject-name">{{ subject.name }}</span>
                <span class="subject-score">{{ subject.progress }}%</span>
              </div>
              <el-progress
                :percentage="subject.progress"
                :color="getProgressColor(subject.progress)"
                :stroke-width="8"
              />
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 最近成绩 -->
    <div class="recent-grades">
      <div class="section-header">
        <h2>最近成绩</h2>
        <el-button type="primary" @click="$router.push('/student/grades')">
          查看全部
        </el-button>
      </div>
      <div class="grade-table">
        <el-table :data="recentGrades" style="width: 100%">
          <el-table-column prop="subject" label="科目" width="120" />
          <el-table-column prop="examTitle" label="考试名称" />
          <el-table-column prop="score" label="分数" width="80">
            <template #default="{ row }">
              <span class="grade-score" :class="getGradeClass(row.score)">
                {{ row.score }}
              </span>
            </template>
          </el-table-column>
          <el-table-column prop="totalScore" label="总分" width="80" />
          <el-table-column prop="rank" label="排名" width="80" />
          <el-table-column prop="examDate" label="考试时间" width="120">
            <template #default="{ row }">
              {{ formatDate(row.examDate) }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="100">
            <template #default="{ row }">
              <el-button
                type="primary"
                size="small"
                @click="viewGradeDetail(row.id)"
              >
                详情
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import {
  Reading, ArrowRight, Clock, Timer, Document, Calendar,
  Notebook, ChatDotRound, Medal, TrendCharts
} from '@element-plus/icons-vue'
import { useAuthStore } from '@/stores/auth'
import * as echarts from 'echarts'

const router = useRouter()
const authStore = useAuthStore()

// 学生统计数据
const studentStats = ref({
  totalExams: 15,
  averageScore: 87.5,
  practiceCount: 142
})

// 快速操作
const quickActions = ref([
  {
    key: 'practice',
    title: '开始练习',
    description: '刷题练习，巩固知识点',
    icon: 'Reading',
    color: 'var(--gradient-primary)'
  },
  {
    key: 'exams',
    title: '我的考试',
    description: '查看考试安排和历史记录',
    icon: 'Notebook',
    color: 'var(--gradient-success)'
  },
  {
    key: 'ai-assistant',
    title: 'AI助手',
    description: '智能答疑，学习辅导',
    icon: 'ChatDotRound',
    color: 'var(--gradient-warning)'
  },
  {
    key: 'grades',
    title: '成绩查询',
    description: '查看考试成绩和分析报告',
    icon: 'Medal',
    color: 'var(--gradient-danger)'
  }
])

// 即将到来的考试
const upcomingExams = ref([
  {
    id: 1,
    title: '数据结构期末考试',
    subject: '数据结构',
    startTime: new Date(Date.now() + 2 * 24 * 60 * 60 * 1000),
    duration: 120,
    questionCount: 50,
    status: 'upcoming'
  },
  {
    id: 2,
    title: '算法设计随堂测验',
    subject: '算法设计',
    startTime: new Date(Date.now() + 5 * 24 * 60 * 60 * 1000),
    duration: 60,
    questionCount: 25,
    status: 'upcoming'
  }
])

// 科目进度
const subjectProgress = ref([
  { id: 1, name: '数据结构', progress: 85 },
  { id: 2, name: '算法设计', progress: 72 },
  { id: 3, name: '操作系统', progress: 90 },
  { id: 4, name: '计算机网络', progress: 68 }
])

// 最近成绩
const recentGrades = ref([
  {
    id: 1,
    subject: '数据结构',
    examTitle: '第三次作业',
    score: 92,
    totalScore: 100,
    rank: 5,
    examDate: new Date(Date.now() - 3 * 24 * 60 * 60 * 1000)
  },
  {
    id: 2,
    subject: '算法设计',
    examTitle: '期中考试',
    score: 88,
    totalScore: 100,
    rank: 12,
    examDate: new Date(Date.now() - 7 * 24 * 60 * 60 * 1000)
  },
  {
    id: 3,
    subject: '操作系统',
    examTitle: '随堂测验',
    score: 95,
    totalScore: 100,
    rank: 2,
    examDate: new Date(Date.now() - 10 * 24 * 60 * 60 * 1000)
  }
])

// 图表引用
const learningChart = ref<HTMLElement>()

// 格式化日期
const formatDate = (date: Date): string => {
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
}

// 格式化日期时间
const formatDateTime = (date: Date): string => {
  return date.toLocaleString('zh-CN', {
    month: 'numeric',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// 获取考试状态类名
const getExamStatusClass = (status: string): string => {
  const classMap: Record<string, string> = {
    upcoming: 'status-upcoming',
    'in-progress': 'status-progress',
    completed: 'status-completed'
  }
  return classMap[status] || ''
}

// 获取考试状态文本
const getExamStatusText = (status: string): string => {
  const textMap: Record<string, string> = {
    upcoming: '即将开始',
    'in-progress': '进行中',
    completed: '已完成'
  }
  return textMap[status] || ''
}

// 获取进度颜色
const getProgressColor = (progress: number): string => {
  if (progress >= 90) return '#67C23A'
  if (progress >= 70) return '#E6A23C'
  return '#F56C6C'
}

// 获取成绩等级类名
const getGradeClass = (score: number): string => {
  if (score >= 90) return 'grade-excellent'
  if (score >= 80) return 'grade-good'
  if (score >= 70) return 'grade-fair'
  return 'grade-poor'
}

// 处理快速操作
const handleQuickAction = (key: string) => {
  switch (key) {
    case 'practice':
      router.push('/student/practice')
      break
    case 'exams':
      router.push('/student/exams')
      break
    case 'ai-assistant':
      router.push('/student/ai-assistant')
      break
    case 'grades':
      router.push('/student/grades')
      break
  }
}

// 进入考试
const enterExam = (examId: number) => {
  router.push(`/exam/${examId}`)
}

// 继续考试
const continueExam = (examId: number) => {
  router.push(`/exam/${examId}`)
}

// 查看考试结果
const viewResult = (examId: number) => {
  ElMessage.info('查看考试结果功能开发中')
}

// 查看成绩详情
const viewGradeDetail = (gradeId: number) => {
  ElMessage.info('查看成绩详情功能开发中')
}

// 初始化学习图表
const initLearningChart = () => {
  if (!learningChart.value) return
  
  const chart = echarts.init(learningChart.value)
  const option = {
    tooltip: {
      trigger: 'axis'
    },
    legend: {
      data: ['练习次数', '考试次数']
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
        name: '练习次数',
        type: 'bar',
        data: [12, 18, 15, 22],
        itemStyle: {
          color: '#4ECDC4'
        }
      },
      {
        name: '考试次数',
        type: 'bar',
        data: [2, 3, 1, 4],
        itemStyle: {
          color: '#FF6B8A'
        }
      }
    ]
  }
  
  chart.setOption(option)
  
  // 响应式
  window.addEventListener('resize', () => {
    chart.resize()
  })
}

// 组件挂载后初始化
onMounted(async () => {
  await nextTick()
  initLearningChart()
})
</script>

<style scoped>
.student-dashboard {
  padding: var(--spacing-md);
}

.welcome-section {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--spacing-md);
  margin-bottom: var(--spacing-lg);
  background: var(--gradient-primary);
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
  border-color: var(--dopamine-blue);
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
  color: var(--dopamine-blue);
  transform: translateX(4px);
}

.upcoming-exams {
  margin-bottom: var(--spacing-lg);
}

.exam-list {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(350px, 1fr));
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
  color: var(--dopamine-blue);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.exam-status {
  font-size: 12px;
  font-weight: 600;
  padding: 4px 8px;
  border-radius: var(--radius-sm);
}

.status-upcoming {
  background: rgba(64, 158, 255, 0.1);
  color: var(--dopamine-blue);
}

.status-progress {
  background: rgba(230, 162, 60, 0.1);
  color: var(--dopamine-orange);
}

.status-completed {
  background: rgba(103, 194, 58, 0.1);
  color: var(--dopamine-green);
}

.exam-title {
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: var(--spacing-md);
}

.exam-details {
  display: flex;
  flex-wrap: wrap;
  gap: var(--spacing-md);
  margin-bottom: var(--spacing-md);
}

.exam-detail {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 14px;
  color: var(--text-secondary);
}

.exam-actions {
  display: flex;
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

.learning-progress {
  margin-bottom: var(--spacing-lg);
}

.progress-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--spacing-lg);
}

.progress-card {
  padding: var(--spacing-lg);
}

.progress-card h3 {
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 var(--spacing-md) 0;
}

.progress-chart {
  height: 200px;
}

.chart {
  width: 100%;
  height: 100%;
}

.subject-progress {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
}

.subject-item {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-sm);
}

.subject-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.subject-name {
  font-size: 14px;
  font-weight: 500;
  color: var(--text-primary);
}

.subject-score {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-secondary);
}

.recent-grades {
  margin-bottom: var(--spacing-xl);
}

.grade-table {
  background: white;
  border-radius: var(--radius-lg);
  overflow: hidden;
  box-shadow: var(--shadow-sm);
}

.grade-score {
  font-weight: 600;
}

.grade-excellent {
  color: var(--dopamine-green);
}

.grade-good {
  color: var(--dopamine-blue);
}

.grade-fair {
  color: var(--dopamine-orange);
}

.grade-poor {
  color: var(--dopamine-red);
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
  
  .action-grid {
    grid-template-columns: 1fr;
  }
  
  .exam-list {
    grid-template-columns: 1fr;
  }
  
  .progress-grid {
    grid-template-columns: 1fr;
  }
}
</style>