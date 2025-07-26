<template>
  <div class="exam-view">
    <!-- 页面标题 -->
    <div class="page-header">
      <h1 class="page-title">我的考试</h1>
      <p class="page-subtitle">查看考试安排，参加在线考试</p>
    </div>
    
    <!-- 考试状态筛选 -->
    <div class="filter-section dopamine-card">
      <div class="filter-tabs">
        <div
          v-for="tab in examTabs"
          :key="tab.key"
          class="filter-tab"
          :class="{ active: activeTab === tab.key }"
          @click="activeTab = tab.key"
        >
          <el-icon><component :is="tab.icon" /></el-icon>
          <span>{{ tab.label }}</span>
          <el-badge
            v-if="tab.count > 0"
            :value="tab.count"
            class="tab-badge"
          />
        </div>
      </div>
    </div>
    
    <!-- 即将开始的考试 -->
    <div v-if="activeTab === 'upcoming'" class="exam-section">
      <div class="section-header">
        <h2>即将开始的考试</h2>
        <el-button type="primary" @click="refreshExams">
          <el-icon><Refresh /></el-icon>
          刷新
        </el-button>
      </div>
      
      <div v-if="upcomingExams.length === 0" class="empty-state">
        <el-empty description="暂无即将开始的考试" />
      </div>
      
      <div v-else class="exam-grid">
        <div
          v-for="exam in upcomingExams"
          :key="exam.id"
          class="exam-card dopamine-card"
        >
          <div class="exam-header">
            <div class="exam-subject">{{ exam.subject }}</div>
            <div class="exam-status" :class="`status-${exam.status}`">
              {{ getStatusText(exam.status) }}
            </div>
          </div>
          
          <div class="exam-content">
            <div class="exam-title-row">
              <div class="exam-title">{{ exam.title }}</div>
              <div class="exam-info">
                <div class="exam-info-item">
                  <el-icon><Clock /></el-icon>
                  <span>{{ exam.duration }}分钟</span>
                </div>
                <div class="exam-info-item">
                  <el-icon><Document /></el-icon>
                  <span>{{ exam.questionCount }}题</span>
                </div>
                <div class="exam-info-item">
                  <el-icon><Star /></el-icon>
                  <span>{{ exam.totalScore }}分</span>
                </div>
              </div>
            </div>
            <div class="exam-description">{{ exam.description }}</div>
            <div class="exam-time">
              <div class="time-item">
                <span class="time-label">开始：</span>
                <span class="time-value">{{ formatDateTimeDetailed(exam.startTime) }}</span>
              </div>
              <div class="time-item">
                <span class="time-label">结束：</span>
                <span class="time-value">{{ formatDateTimeDetailed(exam.endTime) }}</span>
              </div>
            </div>
          </div>
          
          <div class="exam-countdown" v-if="exam.status === 'upcoming'">
            <div class="countdown-label">距离开始还有：</div>
            <div class="countdown-time">{{ getCountdownDetailed(exam.startTime) }}</div>
          </div>
          
          <div class="exam-actions">
            <el-button
              v-if="exam.status === 'ongoing'"
              type="primary"
              size="small"
              @click="enterExam(exam)"
            >
              <el-icon><Right /></el-icon>
              进入考试
            </el-button>
            <el-button
              v-else-if="exam.status === 'upcoming'"
              type="info"
              size="small"
              disabled
            >
              <el-icon><Clock /></el-icon>
              等待开始
            </el-button>
            <el-button
              type="default"
              size="small"
              @click="viewExamDetails(exam)"
            >
              查看详情
            </el-button>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 进行中的考试 -->
    <div v-if="activeTab === 'ongoing'" class="exam-section">
      <div class="section-header">
        <h2>进行中的考试</h2>
      </div>
      
      <div v-if="ongoingExams.length === 0" class="empty-state">
        <el-empty description="暂无进行中的考试" />
      </div>
      
      <div v-else class="exam-grid">
        <div
          v-for="exam in ongoingExams"
          :key="exam.id"
          class="exam-card dopamine-card ongoing-exam"
        >
          <div class="exam-header">
            <div class="exam-subject">{{ exam.subject }}</div>
            <div class="exam-status status-ongoing">
              <el-icon class="status-icon"><Loading /></el-icon>
              进行中
            </div>
          </div>
          
          <div class="exam-title">{{ exam.title }}</div>
          
          <div class="exam-progress">
            <div class="progress-info">
              <span>答题进度</span>
              <span>{{ exam.answeredCount }}/{{ exam.questionCount }}</span>
            </div>
            <el-progress
              :percentage="(exam.answeredCount / exam.questionCount) * 100"
              :color="getProgressColor((exam.answeredCount / exam.questionCount) * 100)"
              :stroke-width="8"
            />
          </div>
          
          <div class="exam-time-left">
            <div class="time-left-label">剩余时间：</div>
            <div class="time-left-value" :class="{ urgent: getTimeLeft(exam.endTime) < 10 }">
              {{ formatTimeLeft(exam.endTime) }}
            </div>
          </div>
          
          <div class="exam-actions">
            <el-button
              type="primary"
              size="small"
              @click="continueExam(exam)"
            >
              <el-icon><Right /></el-icon>
              继续考试
            </el-button>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 已完成的考试 -->
    <div v-if="activeTab === 'completed'" class="exam-section">
      <div class="section-header">
        <h2>已完成的考试</h2>
        <div class="header-actions">
          <el-select v-model="completedFilter" placeholder="筛选科目" clearable>
            <el-option
              v-for="subject in subjects"
              :key="subject.id"
              :label="subject.name"
              :value="subject.id"
            />
          </el-select>
        </div>
      </div>
      
      <div class="completed-table">
        <el-table :data="filteredCompletedExams" style="width: 100%">
          <el-table-column prop="subject" label="科目" min-width="120" />
          <el-table-column prop="title" label="考试名称" min-width="200" />
          <el-table-column prop="score" label="得分" min-width="100">
            <template #default="{ row }">
              <span :class="getScoreClass(row.score, row.totalScore)">
                {{ row.score }}/{{ row.totalScore }}
              </span>
            </template>
          </el-table-column>
          <el-table-column prop="percentage" label="得分率" min-width="100">
            <template #default="{ row }">
              <span :class="getPercentageClass(row.percentage)">
                {{ row.percentage }}%
              </span>
            </template>
          </el-table-column>
          <el-table-column prop="rank" label="排名" min-width="80">
            <template #default="{ row }">
              <el-tag :type="getRankTagType(row.rank)">{{ row.rank }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="duration" label="用时" min-width="100" />
          <el-table-column prop="completedAt" label="完成时间" min-width="150">
            <template #default="{ row }">
              {{ formatDateTime(row.completedAt) }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="150">
            <template #default="{ row }">
              <el-button
                type="primary"
                size="small"
                @click="viewExamResult(row.id)"
              >
                查看成绩
              </el-button>
              <el-button
                type="default"
                size="small"
                @click="reviewExam(row.id)"
              >
                复习
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Clock, Document, Star, Right, Refresh, Loading
} from '@element-plus/icons-vue'

const router = useRouter()

// 当前活跃的标签页
const activeTab = ref('upcoming')

// 标签页配置
const examTabs = ref([
  { key: 'upcoming', label: '即将开始', icon: 'Clock', count: 2 },
  { key: 'ongoing', label: '进行中', icon: 'Loading', count: 1 },
  { key: 'completed', label: '已完成', icon: 'SuccessFilled', count: 15 }
])

// 科目列表
const subjects = ref([
  { id: 1, name: '数据结构' },
  { id: 2, name: '算法设计' },
  { id: 3, name: '操作系统' },
  { id: 4, name: '计算机网络' },
  { id: 5, name: '数据库原理' }
])

// 已完成考试筛选
const completedFilter = ref('')

// 即将开始的考试
const upcomingExams = ref([
  {
    id: 1,
    title: '数据结构期中考试',
    subject: '数据结构',
    description: '包含线性表、栈、队列、树等内容',
    duration: 120,
    questionCount: 50,
    totalScore: 100,
    startTime: new Date(Date.now() + 2 * 60 * 60 * 1000),
    endTime: new Date(Date.now() + 4 * 60 * 60 * 1000),
    status: 'upcoming'
  },
  {
    id: 2,
    title: '算法设计综合测试',
    subject: '算法设计',
    description: '排序、查找、动态规划等算法',
    duration: 90,
    questionCount: 40,
    totalScore: 100,
    startTime: new Date(Date.now() + 24 * 60 * 60 * 1000),
    endTime: new Date(Date.now() + 25.5 * 60 * 60 * 1000),
    status: 'upcoming'
  }
])

// 进行中的考试
const ongoingExams = ref([
  {
    id: 3,
    title: '操作系统单元测试',
    subject: '操作系统',
    duration: 60,
    questionCount: 30,
    answeredCount: 18,
    startTime: new Date(Date.now() - 30 * 60 * 1000),
    endTime: new Date(Date.now() + 30 * 60 * 1000),
    status: 'ongoing'
  }
])

// 已完成的考试
const completedExams = ref([
  {
    id: 4,
    title: '数据库基础测试',
    subject: '数据库原理',
    score: 85,
    totalScore: 100,
    percentage: 85,
    rank: 5,
    duration: '45分钟',
    completedAt: new Date(Date.now() - 3 * 24 * 60 * 60 * 1000)
  },
  {
    id: 5,
    title: '计算机网络期末考试',
    subject: '计算机网络',
    score: 92,
    totalScore: 100,
    percentage: 92,
    rank: 2,
    duration: '118分钟',
    completedAt: new Date(Date.now() - 7 * 24 * 60 * 60 * 1000)
  },
  {
    id: 6,
    title: '数据结构章节测试',
    subject: '数据结构',
    score: 78,
    totalScore: 100,
    percentage: 78,
    rank: 12,
    duration: '52分钟',
    completedAt: new Date(Date.now() - 10 * 24 * 60 * 60 * 1000)
  }
])

// 筛选后的已完成考试
const filteredCompletedExams = computed(() => {
  if (!completedFilter.value) return completedExams.value
  const selectedSubject = subjects.value.find(s => s.id === completedFilter.value)
  if (!selectedSubject) return completedExams.value
  return completedExams.value.filter(exam => exam.subject === selectedSubject.name)
})

// 定时器
let countdownTimer: NodeJS.Timeout | null = null

// 获取状态文本
const getStatusText = (status: string): string => {
  const statusMap: Record<string, string> = {
    upcoming: '即将开始',
    ongoing: '进行中',
    completed: '已完成'
  }
  return statusMap[status] || ''
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

// 格式化详细日期时间 (yyyy-mm-dd HH:mi:ss)
const formatDateTimeDetailed = (date: Date): string => {
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  const seconds = String(date.getSeconds()).padStart(2, '0')
  return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`
}

// 获取倒计时
const getCountdown = (startTime: Date): string => {
  const now = new Date()
  const diff = startTime.getTime() - now.getTime()
  
  if (diff <= 0) return '已开始'
  
  const hours = Math.floor(diff / (1000 * 60 * 60))
  const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60))
  
  if (hours > 0) {
    return `${hours}小时${minutes}分钟`
  }
  return `${minutes}分钟`
}

// 获取精确到秒的倒计时
const getCountdownDetailed = (startTime: Date): string => {
  const now = new Date()
  const diff = startTime.getTime() - now.getTime()
  
  if (diff <= 0) return '已开始'
  
  const hours = Math.floor(diff / (1000 * 60 * 60))
  const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60))
  const seconds = Math.floor((diff % (1000 * 60)) / 1000)
  
  if (hours > 0) {
    return `${hours}小时${minutes}分${seconds}秒`
  } else if (minutes > 0) {
    return `${minutes}分${seconds}秒`
  }
  return `${seconds}秒`
}

// 获取剩余时间（分钟）
const getTimeLeft = (endTime: Date): number => {
  const now = new Date()
  const diff = endTime.getTime() - now.getTime()
  return Math.floor(diff / (1000 * 60))
}

// 格式化剩余时间
const formatTimeLeft = (endTime: Date): string => {
  const minutes = getTimeLeft(endTime)
  if (minutes <= 0) return '已结束'
  
  const hours = Math.floor(minutes / 60)
  const mins = minutes % 60
  
  if (hours > 0) {
    return `${hours}:${mins.toString().padStart(2, '0')}`
  }
  return `${mins}分钟`
}

// 获取进度颜色
const getProgressColor = (percentage: number): string => {
  if (percentage >= 80) return '#67C23A'
  if (percentage >= 50) return '#E6A23C'
  return '#F56C6C'
}

// 获取分数样式类
const getScoreClass = (score: number, totalScore: number): string => {
  const percentage = (score / totalScore) * 100
  if (percentage >= 90) return 'score-excellent'
  if (percentage >= 80) return 'score-good'
  if (percentage >= 70) return 'score-fair'
  return 'score-poor'
}

// 获取得分率样式类
const getPercentageClass = (percentage: number): string => {
  if (percentage >= 90) return 'percentage-excellent'
  if (percentage >= 80) return 'percentage-good'
  if (percentage >= 70) return 'percentage-fair'
  return 'percentage-poor'
}

// 获取排名标签类型
const getRankTagType = (rank: number): string => {
  if (rank <= 3) return 'success'
  if (rank <= 10) return 'warning'
  return 'info'
}

// 进入考试
const enterExam = async (exam: any) => {
  try {
    await ElMessageBox.confirm(
      `确定要进入考试「${exam.title}」吗？\n考试时长：${exam.duration}分钟\n题目数量：${exam.questionCount}题`,
      '进入考试',
      {
        confirmButtonText: '进入考试',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    ElMessage.success('正在进入考试...')
    // 这里可以跳转到考试页面
    router.push(`/student/exam/${exam.id}/taking`)
  } catch {
    // 用户取消
  }
}

// 继续考试
const continueExam = (exam: any) => {
  ElMessage.success('正在恢复考试...')
  router.push(`/exam/${exam.id}`)
}

// 查看考试详情
const viewExamDetails = (exam: any) => {
  ElMessage.info('查看考试详情功能开发中')
}

// 查看考试结果
const viewExamResult = (examId: number) => {
  ElMessage.info('查看考试结果功能开发中')
}

// 复习考试
const reviewExam = (examId: number) => {
  ElMessage.info('复习考试功能开发中')
}

// 刷新考试列表
const refreshExams = () => {
  ElMessage.success('考试列表已刷新')
}

// 启动倒计时定时器
const startCountdownTimer = () => {
  countdownTimer = setInterval(() => {
    // 更新倒计时显示，强制重新渲染
    upcomingExams.value = [...upcomingExams.value]
  }, 1000) // 每秒更新一次
}

// 停止倒计时定时器
const stopCountdownTimer = () => {
  if (countdownTimer) {
    clearInterval(countdownTimer)
    countdownTimer = null
  }
}

onMounted(() => {
  startCountdownTimer()
})

onUnmounted(() => {
  stopCountdownTimer()
})
</script>

<style scoped>
.exam-view {
  padding: var(--spacing-lg);
}

.page-header {
  margin-bottom: var(--spacing-xl);
}

.page-title {
  font-size: 28px;
  font-weight: 700;
  color: var(--text-primary);
  margin: 0 0 var(--spacing-sm) 0;
  background: var(--gradient-primary);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.page-subtitle {
  font-size: 16px;
  color: var(--text-secondary);
  margin: 0;
}

.filter-section {
  padding: 8px 16px;
  margin-bottom: 16px;
}

.filter-tabs {
  display: flex;
  gap: 8px;
}

.filter-tab {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 6px 12px;
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: all 0.3s ease;
  background: var(--bg-secondary);
  color: var(--text-secondary);
  position: relative;
  font-size: 14px;
}

.filter-tab:hover {
  background: var(--dopamine-blue-light);
  color: var(--dopamine-blue);
}

.filter-tab.active {
  background: var(--gradient-primary);
  color: white;
  box-shadow: var(--shadow-md);
}

.tab-badge {
  margin-left: 4px;
  font-size: 12px;
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

.header-actions {
  display: flex;
  gap: var(--spacing-md);
}

.empty-state {
  text-align: center;
  padding: var(--spacing-xl) 0;
}

.exam-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(320px, 1fr));
  gap: 12px;
}

.exam-card {
  padding: 16px;
  min-height: calc(200px + 10px);
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  background: linear-gradient(135deg, 
    rgba(255, 255, 255, 0.95) 0%, 
    rgba(248, 250, 252, 0.9) 50%, 
    rgba(241, 245, 249, 0.85) 100%);
  backdrop-filter: blur(10px);
  border-radius: 16px;
  box-shadow: 
    0 4px 20px rgba(0, 0, 0, 0.08),
    0 1px 3px rgba(0, 0, 0, 0.1),
    inset 0 1px 0 rgba(255, 255, 255, 0.6);
  position: relative;
  overflow: hidden;
}

.exam-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 3px;
  background: linear-gradient(90deg, 
    #667eea 0%, 
    #764ba2 25%, 
    #f093fb 50%, 
    #f5576c 75%, 
    #4facfe 100%);
  opacity: 0.8;
}

.exam-card:hover {
  transform: translateY(-8px) scale(1.02);
  box-shadow: 
    0 20px 40px rgba(0, 0, 0, 0.15),
    0 8px 16px rgba(0, 0, 0, 0.1),
    inset 0 1px 0 rgba(255, 255, 255, 0.8);
  border-color: rgba(102, 126, 234, 0.3);
}

.exam-card:hover::before {
  opacity: 1;
  height: 4px;
}

.ongoing-exam {
  border-color: rgba(255, 193, 7, 0.3);
  background: linear-gradient(135deg, 
    rgba(255, 193, 7, 0.1) 0%, 
    rgba(255, 235, 59, 0.05) 30%,
    rgba(255, 255, 255, 0.95) 70%, 
    rgba(248, 250, 252, 0.9) 100%);
  box-shadow: 
    0 4px 20px rgba(255, 193, 7, 0.15),
    0 1px 3px rgba(0, 0, 0, 0.1),
    inset 0 1px 0 rgba(255, 255, 255, 0.6);
}

.ongoing-exam::before {
  background: linear-gradient(90deg, 
    #ffc107 0%, 
    #ff9800 25%, 
    #ff5722 50%, 
    #e91e63 75%, 
    #9c27b0 100%);
}

.exam-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.exam-subject {
  font-size: 13px;
  font-weight: 700;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  position: relative;
}

.exam-status {
  font-size: 12px;
  font-weight: 700;
  padding: 6px 12px;
  border-radius: 20px;
  display: flex;
  align-items: center;
  gap: 4px;
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.3);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.status-upcoming {
  background: linear-gradient(135deg, 
    rgba(108, 117, 125, 0.15) 0%, 
    rgba(148, 163, 184, 0.1) 100%);
  color: #64748b;
  border-color: rgba(108, 117, 125, 0.2);
}

.status-ongoing {
  background: linear-gradient(135deg, 
    rgba(255, 193, 7, 0.2) 0%, 
    rgba(255, 152, 0, 0.15) 100%);
  color: #f59e0b;
  border-color: rgba(255, 193, 7, 0.3);
}

.status-icon {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.exam-content {
  flex: 1;
}

.exam-title-row {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 6px;
}

.exam-title {
  font-size: 18px;
  font-weight: 700;
  background: linear-gradient(135deg, #2d3748 0%, #4a5568 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  flex: 1;
  margin-right: 12px;
  line-height: 1.3;
}

.exam-description {
  font-size: 14px;
  color: var(--text-secondary);
  margin-bottom: 6px;
  line-height: 1.2;
}

.exam-info {
  display: flex;
  gap: 10px;
  flex-shrink: 0;
}

.exam-info-item {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 13px;
  color: var(--text-secondary);
  background: rgba(255, 255, 255, 0.6);
  padding: 4px 8px;
  border-radius: 8px;
  border: 1px solid rgba(255, 255, 255, 0.3);
  backdrop-filter: blur(5px);
}

.exam-info-item .el-icon {
  color: #667eea;
}

.exam-time {
  margin-bottom: 6px;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.time-item {
  display: flex;
  align-items: center;
  font-size: 13px;
  gap: 6px;
  background: rgba(255, 255, 255, 0.5);
  padding: 6px 10px;
  border-radius: 10px;
  border: 1px solid rgba(255, 255, 255, 0.3);
  backdrop-filter: blur(5px);
}

.time-label {
  color: #64748b;
  font-size: 12px;
  font-weight: 600;
}

.time-value {
  background: linear-gradient(135deg, #1e293b 0%, #475569 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  font-weight: 600;
  font-size: 12px;
}

.exam-countdown {
  background: linear-gradient(135deg, 
    rgba(102, 126, 234, 0.1) 0%, 
    rgba(118, 75, 162, 0.05) 100%);
  padding: 8px 12px;
  border-radius: 12px;
  margin-bottom: 6px;
  text-align: center;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  border: 1px solid rgba(102, 126, 234, 0.2);
  box-shadow: 0 2px 8px rgba(102, 126, 234, 0.1);
}

.countdown-label {
  font-size: 12px;
  color: var(--text-secondary);
  margin: 0;
}

.countdown-time {
  font-size: 16px;
  font-weight: 700;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.exam-progress {
  margin-bottom: 6px;
}

.progress-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 4px;
  font-size: 13px;
  color: var(--text-secondary);
}

.exam-time-left {
  background: linear-gradient(135deg, 
    rgba(245, 158, 11, 0.15) 0%, 
    rgba(251, 146, 60, 0.1) 100%);
  padding: 10px 12px;
  border-radius: 12px;
  margin-bottom: 8px;
  text-align: center;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  border: 1px solid rgba(245, 158, 11, 0.2);
  box-shadow: 0 2px 8px rgba(245, 158, 11, 0.1);
  backdrop-filter: blur(10px);
}

.time-left-label {
  font-size: 12px;
  color: var(--text-secondary);
  margin: 0;
}

.time-left-value {
  font-size: 16px;
  font-weight: 700;
  background: linear-gradient(135deg, #f59e0b 0%, #f97316 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.time-left-value.urgent {
  color: var(--dopamine-red);
  animation: pulse 1s infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.7; }
}

.exam-actions {
  display: flex;
  gap: 8px;
  justify-content: center;
  margin-top: 8px;
  padding-top: 8px;
  border-top: 1px solid rgba(255, 255, 255, 0.3);
}

.exam-actions .el-button {
  border-radius: 12px;
  font-weight: 600;
  letter-spacing: 0.3px;
  backdrop-filter: blur(10px);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.exam-actions .el-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.15);
}

.exam-actions .el-button--primary {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-color: rgba(102, 126, 234, 0.3);
}

.exam-actions .el-button--info {
  background: linear-gradient(135deg, 
    rgba(108, 117, 125, 0.1) 0%, 
    rgba(148, 163, 184, 0.05) 100%);
  color: #64748b;
  border-color: rgba(108, 117, 125, 0.2);
}

.exam-actions .el-button--default {
  background: linear-gradient(135deg, 
    rgba(255, 255, 255, 0.8) 0%, 
    rgba(248, 250, 252, 0.6) 100%);
  color: #374151;
  border-color: rgba(209, 213, 219, 0.3);
}

.completed-table {
  background: white;
  border-radius: var(--radius-lg);
  overflow: hidden;
  box-shadow: var(--shadow-sm);
}

.score-excellent, .percentage-excellent {
  color: var(--dopamine-green);
  font-weight: 600;
}

.score-good, .percentage-good {
  color: var(--dopamine-blue);
  font-weight: 600;
}

.score-fair, .percentage-fair {
  color: var(--dopamine-orange);
  font-weight: 600;
}

.score-poor, .percentage-poor {
  color: var(--dopamine-red);
  font-weight: 600;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .filter-tabs {
    flex-wrap: wrap;
  }
  
  .exam-grid {
    grid-template-columns: 1fr;
  }
  
  .exam-actions {
    flex-direction: column;
  }
}

@media (max-width: 480px) {
  .section-header {
    flex-direction: column;
    gap: var(--spacing-md);
    align-items: stretch;
  }
  
  .exam-info {
    flex-direction: column;
    gap: var(--spacing-sm);
  }
}
</style>