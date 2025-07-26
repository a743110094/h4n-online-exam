<template>
  <div class="review-view">
    <!-- 页面头部 -->
    <div class="page-header">
      <h1 class="page-title">
        <Reading class="title-icon" />
        {{ practiceRecordId ? '本次练习错题复习' : '错题复习' }}
      </h1>
      <p class="page-subtitle">
        {{ practiceRecordId ? '查看和复习本次练习的错题' : '查看和复习您的错题，巩固薄弱知识点' }}
      </p>
    </div>

    <!-- 筛选区域 -->
    <el-card class="filter-section" shadow="never">
      <div class="filter-row">
        <div class="filter-item">
          <label>科目筛选</label>
          <el-select v-model="filters.subjectId" placeholder="选择科目" clearable>
            <el-option label="全部科目" :value="undefined" />
            <el-option
              v-for="subject in subjectOptions"
              :key="subject.value"
              :label="subject.label"
              :value="subject.value"
            />
          </el-select>
        </div>
        <div class="filter-actions">
          <el-button type="primary" @click="loadWrongQuestions">
            <Search class="btn-icon" />
            查询
          </el-button>
          <el-button @click="resetFilters">
            <Refresh class="btn-icon" />
            重置
          </el-button>
        </div>
      </div>
    </el-card>

    <!-- 统计信息 -->
    <div class="stats-section">
      <div class="stats-grid">
        <el-card class="stat-card" shadow="hover">
          <div class="stat-content">
            <div class="stat-icon error-bg">
              <Warning />
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ wrongQuestions.length }}</div>
              <div class="stat-label">当前错题</div>
            </div>
          </div>
        </el-card>
        <el-card class="stat-card" shadow="hover">
          <div class="stat-content">
            <div class="stat-icon primary-bg">
              <Document />
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ pagination.total }}</div>
              <div class="stat-label">总错题数</div>
            </div>
          </div>
        </el-card>
      </div>
    </div>

    <!-- 操作区域 -->
    <div class="action-section">
      <el-button
        type="primary"
        size="large"
        :disabled="wrongQuestions.length === 0"
        @click="startReview"
      >
        <VideoPlay class="btn-icon" />
        开始复习全部错题
      </el-button>
      <el-button
        v-if="selectedQuestions.length > 0"
        type="success"
        size="large"
        @click="startSelectedReview"
      >
        <Star class="btn-icon" />
        复习选中题目 ({{ selectedQuestions.length }})
      </el-button>
    </div>

    <!-- 空状态显示 -->
    <div v-if="!loading && wrongQuestions.length === 0" class="empty-container">
      <div class="modern-empty">
        <!-- 背景装饰 -->
        <div class="bg-decoration">
          <div class="floating-shape shape-1"></div>
          <div class="floating-shape shape-2"></div>
          <div class="floating-shape shape-3"></div>
          <div class="floating-shape shape-4"></div>
        </div>
        
        <!-- 主要内容 -->
        <div class="empty-main">
          <!-- 图标区域 -->
          <div class="icon-section">
            <div class="success-icon">
              <div class="icon-circle">
                <div class="checkmark">
                  <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <path d="M9 12L11 14L15 10" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                  </svg>
                </div>
              </div>
              <div class="pulse-rings">
                <div class="pulse-ring ring-1"></div>
                <div class="pulse-ring ring-2"></div>
                <div class="pulse-ring ring-3"></div>
              </div>
            </div>
            
            <!-- 粒子效果 -->
            <div class="particles">
              <div class="particle" v-for="i in 12" :key="i" :style="{animationDelay: (i * 0.1) + 's'}"></div>
            </div>
          </div>
          
          <!-- 文本内容 -->
          <div class="content-section">
            <h2 class="success-title">
              <span class="title-text">
                {{ practiceRecordId ? '完美表现！' : '学霸模式！' }}
              </span>
              <div class="title-decoration"></div>
            </h2>
            
            <p class="success-subtitle">
              {{ practiceRecordId 
                ? '本次练习全部正确，你的努力得到了回报' 
                : '目前没有错题记录，继续保持优秀状态' 
              }}
            </p>
            
            <!-- 成就徽章 -->
            <div class="achievement-badge">
              <div class="badge-icon">
                <Star class="star-icon" />
              </div>
              <span class="badge-text">{{ practiceRecordId ? '零错误' : '无错题' }}</span>
            </div>
          </div>
          
          <!-- 操作按钮 -->
          <div class="action-section">
            <button 
              @click="$router.push('/student/practice')"
              class="modern-btn primary-btn"
            >
              <div class="btn-content">
                <VideoPlay class="btn-icon" />
                <span>继续挑战</span>
              </div>
              <div class="btn-glow"></div>
            </button>
            
            <button 
              v-if="!practiceRecordId"
              @click="$router.push('/student/practice')"
              class="modern-btn secondary-btn"
            >
              <div class="btn-content">
                <Search class="btn-icon" />
                <span>浏览题库</span>
              </div>
            </button>
            
            <button 
              v-if="practiceRecordId"
              @click="$router.go(-1)"
              class="modern-btn secondary-btn"
            >
              <div class="btn-content">
                <Refresh class="btn-icon" />
                <span>返回历史</span>
              </div>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 错题列表 -->
    <el-card v-else class="questions-section" shadow="never">
      <template #header>
        <div class="section-header">
          <h2>错题列表</h2>
          <div class="header-actions">
            <el-checkbox
              v-model="selectAll"
              :indeterminate="isIndeterminate"
              @change="handleSelectAll"
            >
              全选
            </el-checkbox>
          </div>
        </div>
      </template>

      <div v-if="loading" class="loading-container">
        <el-skeleton :rows="5" animated />
      </div>

      <div v-else class="questions-list">
        <div
          v-for="question in wrongQuestions"
          :key="question.id"
          class="question-item"
        >
          <div class="question-header">
            <el-checkbox
              v-model="selectedQuestions"
              :value="question.question_id"
              class="question-checkbox"
            />
            <div class="question-meta">
              <span class="question-title">{{ question.title }}</span>
              <div class="question-tags">
                <el-tag :type="getDifficultyTagType(question.difficulty)" size="small">
                  {{ getDifficultyText(question.difficulty) }}
                </el-tag>
                <el-tag type="info" size="small">{{ question.subject_name }}</el-tag>
                <el-tag :type="getTypeTagType(question.type)" size="small">
                  {{ getTypeText(question.type) }}
                </el-tag>
              </div>
            </div>
            <div class="question-actions">
              <el-button
                type="primary"
                size="small"
                @click="reviewSingleQuestion(question.question_id)"
              >
                单独复习
              </el-button>
            </div>
          </div>

          <div class="question-content">
            <div class="question-text" v-html="question.content"></div>
            
            <!-- 选择题选项 -->
            <div v-if="question.type !== 'fill'" class="question-options">
              <div
                v-for="(option, index) in parseOptions(question.options)"
                :key="index"
                class="option-item"
                :class="{
                  'user-answer': isUserAnswer(option.key, question.user_answer),
                  'correct-answer': isCorrectAnswer(option.key, question.correct_answer),
                  'wrong-answer': isUserAnswer(option.key, question.user_answer) && !isCorrectAnswer(option.key, question.correct_answer)
                }"
              >
                <span class="option-key">{{ option.key }}.</span>
                <span class="option-text">{{ option.text }}</span>
                <span v-if="isUserAnswer(option.key, question.user_answer)" class="answer-label user-label">
                  您的答案
                </span>
                <span v-if="isCorrectAnswer(option.key, question.correct_answer)" class="answer-label correct-label">
                  正确答案
                </span>
              </div>
            </div>

            <!-- 填空题 -->
            <div v-else class="fill-answer">
              <div class="answer-comparison">
                <div class="answer-item wrong">
                  <span class="answer-label">您的答案：</span>
                  <span class="answer-text">{{ question.user_answer || '未作答' }}</span>
                </div>
                <div class="answer-item correct">
                  <span class="answer-label">正确答案：</span>
                  <span class="answer-text">{{ question.correct_answer }}</span>
                </div>
              </div>
            </div>

            <!-- 解析 -->
            <div v-if="question.explanation" class="question-explanation">
              <h4>题目解析</h4>
              <p v-html="question.explanation"></p>
            </div>

            <!-- 答题信息 -->
            <div class="question-info">
              <span class="info-item">
                <Clock class="info-icon" />
                用时: {{ formatTime(question.time_spent) }}
              </span>
              <span class="info-item">
                <Calendar class="info-icon" />
                答题时间: {{ formatDate(question.created_at) }}
              </span>
            </div>
          </div>
        </div>
      </div>

      <!-- 分页 -->
      <div v-if="pagination.total > pagination.pageSize" class="pagination-container">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :total="pagination.total"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next, jumper"
          @current-change="loadWrongQuestions"
          @size-change="loadWrongQuestions"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import {
  Reading, Search, Refresh, Warning, Document, VideoPlay, Star,
  Clock, Calendar
} from '@element-plus/icons-vue'
import * as practiceApi from '@/api/practice'
import * as subjectApi from '@/api/subject'

const router = useRouter()
const route = useRoute()

// 获取路由参数中的练习记录ID
const practiceRecordId = ref<number | undefined>()
if (route.query.practiceRecordId) {
  practiceRecordId.value = Number(route.query.practiceRecordId)
}

// 筛选条件
const filters = reactive({
  subjectId: undefined as number | undefined
})

// 科目选项
const subjectOptions = ref<Array<{ label: string; value: number }>>([])

// 错题数据
const wrongQuestions = ref<practiceApi.WrongQuestionDetail[]>([])
const loading = ref(false)

// 分页信息
const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

// 选择相关
const selectedQuestions = ref<number[]>([])
const selectAll = ref(false)

// 计算属性
const isIndeterminate = computed(() => {
  const selected = selectedQuestions.value.length
  const total = wrongQuestions.value.length
  return selected > 0 && selected < total
})

// 获取难度文本
const getDifficultyText = (difficulty: number) => {
  const map: Record<number, string> = {
    1: '简单',
    2: '中等',
    3: '困难',
    4: '专家',
    5: '大师'
  }
  return map[difficulty] || '未知'
}

// 获取难度标签类型
const getDifficultyTagType = (difficulty: number) => {
  const map: Record<number, string> = {
    1: 'success',
    2: 'warning',
    3: 'danger',
    4: 'danger',
    5: 'danger'
  }
  return map[difficulty] || 'info'
}

// 获取题目类型文本
const getTypeText = (type: string) => {
  const map: Record<string, string> = {
    single: '单选题',
    multiple: '多选题',
    judge: '判断题',
    fill: '填空题'
  }
  return map[type] || '未知'
}

// 获取题目类型标签类型
const getTypeTagType = (type: string) => {
  const map: Record<string, string> = {
    single: 'primary',
    multiple: 'success',
    judge: 'warning',
    fill: 'info'
  }
  return map[type] || 'info'
}

// 解析选项
const parseOptions = (optionsStr: string) => {
  try {
    const options = JSON.parse(optionsStr)
    return Object.entries(options).map(([key, text]) => ({ key, text }))
  } catch {
    return []
  }
}

// 判断是否为用户答案
const isUserAnswer = (optionKey: string, userAnswer: string) => {
  return userAnswer.includes(optionKey)
}

// 判断是否为正确答案
const isCorrectAnswer = (optionKey: string, correctAnswer: string) => {
  return correctAnswer.includes(optionKey)
}

// 格式化时间
const formatTime = (seconds: number) => {
  const mins = Math.floor(seconds / 60)
  const secs = seconds % 60
  return `${mins}分${secs}秒`
}

// 格式化日期
const formatDate = (dateString: string) => {
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN')
}

// 加载科目列表
const loadSubjects = async () => {
  try {
    const response = await subjectApi.getSubjects()
    subjectOptions.value = response.subjects.map(subject => ({
      label: subject.name,
      value: subject.id
    }))
  } catch (error) {
    console.error('加载科目列表失败:', error)
  }
}

// 加载错题列表
const loadWrongQuestions = async () => {
  try {
    loading.value = true
    console.log('开始加载错题列表，loading状态:', loading.value)
    
    const params: any = {
      page: pagination.page,
      page_size: pagination.pageSize,
      subject_id: filters.subjectId
    }
    
    // 如果有练习记录ID，添加到参数中
    if (practiceRecordId.value) {
      params.practice_record_id = practiceRecordId.value
    }
    
    console.log('请求参数:', params)
    const response = await practiceApi.getWrongQuestions(params)
    console.log('API响应:', response)
    
    wrongQuestions.value = response.data || []
    pagination.total = response.total || 0
    selectedQuestions.value = []
    selectAll.value = false
    
    console.log('错题数量:', wrongQuestions.value.length)
    console.log('总数:', pagination.total)
  } catch (error) {
    console.error('加载错题列表失败:', error)
    ElMessage.error('加载错题列表失败')
    wrongQuestions.value = []
    pagination.total = 0
  } finally {
    loading.value = false
    console.log('加载完成，loading状态:', loading.value)
    console.log('最终错题数量:', wrongQuestions.value.length)
  }
}

// 重置筛选条件
const resetFilters = () => {
  filters.subjectId = undefined
  pagination.page = 1
  loadWrongQuestions()
}

// 全选处理
const handleSelectAll = (checked: boolean) => {
  if (checked) {
    selectedQuestions.value = wrongQuestions.value.map(q => q.question_id)
  } else {
    selectedQuestions.value = []
  }
}

// 开始复习全部错题
const startReview = async () => {
  try {
    const params: any = {
      subject_id: filters.subjectId,
      max_questions: 20
    }
    
    // 如果有练习记录ID，添加到参数中
    if (practiceRecordId.value) {
      params.practice_record_id = practiceRecordId.value
    }
    
    const response = await practiceApi.startWrongQuestionReview(params)
    
    // 跳转到练习页面
    router.push({
      path: `/student/practice-session/${response.practice_id}`,
      query: {
        type: 'review',
        title: practiceRecordId.value ? '本次练习错题复习' : '错题复习'
      }
    })
  } catch (error) {
    console.error('开始复习失败:', error)
    ElMessage.error('开始复习失败')
  }
}

// 开始选中题目复习
const startSelectedReview = async () => {
  try {
    const response = await practiceApi.startWrongQuestionReview({
      question_ids: selectedQuestions.value
    })
    
    // 跳转到练习页面
    router.push({
      path: `/student/practice-session/${response.practice_id}`,
      query: {
        type: 'review',
        title: `错题复习 (${selectedQuestions.value.length}题)`
      }
    })
  } catch (error) {
    console.error('开始复习失败:', error)
    ElMessage.error('开始复习失败')
  }
}

// 单独复习某个题目
const reviewSingleQuestion = async (questionId: number) => {
  try {
    const response = await practiceApi.startWrongQuestionReview({
      question_ids: [questionId]
    })
    
    // 跳转到练习页面
    router.push({
      path: `/student/practice-session/${response.practice_id}`,
      query: {
        type: 'review',
        title: '单题复习'
      }
    })
  } catch (error) {
    console.error('开始复习失败:', error)
    ElMessage.error('开始复习失败')
  }
}

// 初始化
onMounted(() => {
  loadSubjects()
  loadWrongQuestions()
})
</script>

<style scoped>
.review-view {
  padding: var(--spacing-md);
}

.page-header {
  margin-bottom: var(--spacing-lg);
}

.page-title {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  font-size: 28px;
  font-weight: 700;
  color: var(--text-primary);
  margin: 0 0 var(--spacing-sm) 0;
}

.title-icon {
  width: 32px;
  height: 32px;
  color: var(--color-primary);
}

.page-subtitle {
  font-size: 16px;
  color: var(--text-secondary);
  margin: 0;
}

.filter-section {
  margin-bottom: var(--spacing-lg);
}

.filter-row {
  display: flex;
  gap: var(--spacing-lg);
  align-items: end;
}

.filter-item {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-sm);
  min-width: 150px;
}

.filter-item label {
  font-size: 14px;
  font-weight: 500;
  color: var(--text-primary);
}

.filter-actions {
  display: flex;
  gap: var(--spacing-md);
}

.btn-icon {
  margin-right: var(--spacing-xs);
}

.stats-section {
  margin-bottom: var(--spacing-lg);
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: var(--spacing-lg);
}

.stat-card {
  transition: all 0.3s ease;
}

.stat-card:hover {
  transform: translateY(-2px);
}

.stat-content {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
}

.stat-icon {
  width: 50px;
  height: 50px;
  border-radius: var(--radius-lg);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.error-bg {
  background: var(--color-danger);
}

.primary-bg {
  background: var(--color-primary);
}

.stat-info {
  flex: 1;
}

.stat-value {
  font-size: 24px;
  font-weight: 700;
  color: var(--text-primary);
  margin-bottom: 4px;
}

.stat-label {
  font-size: 14px;
  color: var(--text-secondary);
}

.action-section {
  margin-bottom: var(--spacing-lg);
  display: flex;
  gap: var(--spacing-md);
}

.questions-section {
  margin-bottom: var(--spacing-lg);
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.section-header h2 {
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0;
}

.loading-container,
.empty-container {
  padding: var(--spacing-xl);
}

.questions-list {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-lg);
}

.question-item {
  border: 1px solid var(--border-color);
  border-radius: var(--radius-md);
  padding: var(--spacing-lg);
  background: var(--bg-color);
  transition: all 0.3s ease;
}

.question-item:hover {
  border-color: var(--color-primary);
  box-shadow: var(--shadow-md);
}

.question-header {
  display: flex;
  align-items: flex-start;
  gap: var(--spacing-md);
  margin-bottom: var(--spacing-md);
}

.question-checkbox {
  margin-top: 4px;
}

.question-meta {
  flex: 1;
}

.question-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
  display: block;
  margin-bottom: var(--spacing-sm);
}

.question-tags {
  display: flex;
  gap: var(--spacing-sm);
  flex-wrap: wrap;
}

.question-content {
  margin-left: 32px;
}

.question-text {
  font-size: 15px;
  line-height: 1.6;
  color: var(--text-primary);
  margin-bottom: var(--spacing-md);
}

.question-options {
  margin-bottom: var(--spacing-md);
}

.option-item {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  padding: var(--spacing-sm);
  margin-bottom: var(--spacing-xs);
  border-radius: var(--radius-sm);
  position: relative;
}

.option-item.user-answer {
  background: #fef0f0;
  border: 1px solid #fbc4c4;
}

.option-item.correct-answer {
  background: #f0f9ff;
  border: 1px solid #93c5fd;
}

.option-item.wrong-answer {
  background: #fef2f2;
  border: 1px solid #fca5a5;
}

.option-key {
  font-weight: 600;
  color: var(--text-primary);
  min-width: 20px;
}

.option-text {
  flex: 1;
  color: var(--text-primary);
}

.answer-label {
  font-size: 12px;
  padding: 2px 6px;
  border-radius: var(--radius-xs);
  font-weight: 500;
}

.user-label {
  background: var(--color-danger);
  color: white;
}

.correct-label {
  background: var(--color-success);
  color: white;
}

.fill-answer {
  margin-bottom: var(--spacing-md);
}

.answer-comparison {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-sm);
}

.answer-item {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  padding: var(--spacing-sm);
  border-radius: var(--radius-sm);
}

.answer-item.wrong {
  background: #fef2f2;
  border: 1px solid #fca5a5;
}

.answer-item.correct {
  background: #f0fdf4;
  border: 1px solid #86efac;
}

.answer-item .answer-label {
  font-weight: 600;
  min-width: 80px;
}

.answer-item .answer-text {
  flex: 1;
}

.question-explanation {
  background: var(--bg-secondary);
  padding: var(--spacing-md);
  border-radius: var(--radius-sm);
  margin-bottom: var(--spacing-md);
}

.question-explanation h4 {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 var(--spacing-sm) 0;
}

.question-explanation p {
  font-size: 14px;
  line-height: 1.6;
  color: var(--text-secondary);
  margin: 0;
}

.question-info {
  display: flex;
  gap: var(--spacing-lg);
  font-size: 13px;
  color: var(--text-secondary);
}

.info-item {
  display: flex;
  align-items: center;
  gap: var(--spacing-xs);
}

.info-icon {
  width: 14px;
  height: 14px;
}

.pagination-container {
  margin-top: var(--spacing-lg);
  display: flex;
  justify-content: center;
}

/* 现代化空状态样式 */
.modern-empty {
  position: relative;
  min-height: 500px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, 
    rgba(99, 102, 241, 0.05) 0%, 
    rgba(168, 85, 247, 0.05) 50%, 
    rgba(236, 72, 153, 0.05) 100%);
  border-radius: 24px;
  overflow: hidden;
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.1);
}

/* 背景装饰 */
.bg-decoration {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
  overflow: hidden;
}

.floating-shape {
  position: absolute;
  border-radius: 50%;
  background: linear-gradient(135deg, rgba(99, 102, 241, 0.1), rgba(168, 85, 247, 0.1));
  animation: float-shape 8s ease-in-out infinite;
}

.shape-1 {
  width: 120px;
  height: 120px;
  top: 10%;
  left: 10%;
  animation-delay: 0s;
}

.shape-2 {
  width: 80px;
  height: 80px;
  top: 20%;
  right: 15%;
  animation-delay: 2s;
}

.shape-3 {
  width: 60px;
  height: 60px;
  bottom: 20%;
  left: 20%;
  animation-delay: 4s;
}

.shape-4 {
  width: 100px;
  height: 100px;
  bottom: 15%;
  right: 10%;
  animation-delay: 6s;
}

/* 主要内容区域 */
.empty-main {
  position: relative;
  z-index: 2;
  text-align: center;
  max-width: 500px;
  padding: 40px;
}

/* 图标区域 */
.icon-section {
  position: relative;
  margin-bottom: 40px;
  display: flex;
  justify-content: center;
  align-items: center;
}

.success-icon {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
}

.icon-circle {
  width: 120px;
  height: 120px;
  border-radius: 50%;
  background: linear-gradient(135deg, #10b981, #059669);
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 
    0 20px 40px rgba(16, 185, 129, 0.3),
    0 0 0 1px rgba(255, 255, 255, 0.1);
  position: relative;
  z-index: 3;
  animation: icon-pulse 2s ease-in-out infinite;
}

.checkmark {
  width: 60px;
  height: 60px;
  color: white;
  animation: checkmark-draw 1s ease-in-out;
}

.checkmark svg {
  width: 100%;
  height: 100%;
}

/* 脉冲环效果 */
.pulse-rings {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
}

.pulse-ring {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  border: 2px solid rgba(16, 185, 129, 0.3);
  border-radius: 50%;
  animation: pulse-ring 2s ease-out infinite;
}

.ring-1 {
  width: 140px;
  height: 140px;
  animation-delay: 0s;
}

.ring-2 {
  width: 160px;
  height: 160px;
  animation-delay: 0.5s;
}

.ring-3 {
  width: 180px;
  height: 180px;
  animation-delay: 1s;
}

/* 粒子效果 */
.particles {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 200px;
  height: 200px;
  pointer-events: none;
}

.particle {
  position: absolute;
  width: 4px;
  height: 4px;
  background: linear-gradient(45deg, #10b981, #3b82f6);
  border-radius: 50%;
  animation: particle-float 3s ease-in-out infinite;
}

.particle:nth-child(1) { top: 0%; left: 50%; }
.particle:nth-child(2) { top: 8.3%; left: 75%; }
.particle:nth-child(3) { top: 25%; left: 93.3%; }
.particle:nth-child(4) { top: 50%; left: 100%; }
.particle:nth-child(5) { top: 75%; left: 93.3%; }
.particle:nth-child(6) { top: 91.7%; left: 75%; }
.particle:nth-child(7) { top: 100%; left: 50%; }
.particle:nth-child(8) { top: 91.7%; left: 25%; }
.particle:nth-child(9) { top: 75%; left: 6.7%; }
.particle:nth-child(10) { top: 50%; left: 0%; }
.particle:nth-child(11) { top: 25%; left: 6.7%; }
.particle:nth-child(12) { top: 8.3%; left: 25%; }

/* 内容区域 */
.content-section {
  margin-bottom: 40px;
}

.success-title {
  position: relative;
  margin-bottom: 16px;
}

.title-text {
  font-size: 32px;
  font-weight: 800;
  background: linear-gradient(135deg, #1e293b, #475569);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  display: inline-block;
  margin: 0;
  letter-spacing: -0.02em;
}

.title-decoration {
  width: 60px;
  height: 4px;
  background: linear-gradient(90deg, #10b981, #3b82f6);
  border-radius: 2px;
  margin: 12px auto 0;
  animation: decoration-glow 2s ease-in-out infinite alternate;
}

.success-subtitle {
  font-size: 16px;
  color: #64748b;
  line-height: 1.6;
  margin: 0 0 24px 0;
  max-width: 400px;
  margin-left: auto;
  margin-right: auto;
}

/* 成就徽章 */
.achievement-badge {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 12px 20px;
  background: linear-gradient(135deg, rgba(16, 185, 129, 0.1), rgba(59, 130, 246, 0.1));
  border: 1px solid rgba(16, 185, 129, 0.2);
  border-radius: 50px;
  backdrop-filter: blur(10px);
  animation: badge-glow 3s ease-in-out infinite;
}

.badge-icon {
  width: 20px;
  height: 20px;
  color: #10b981;
}

.star-icon {
  width: 100%;
  height: 100%;
  animation: star-twinkle 2s ease-in-out infinite;
}

.badge-text {
  font-size: 14px;
  font-weight: 600;
  color: #10b981;
}

/* 操作按钮区域 */
.action-section {
  display: flex;
  gap: 16px;
  justify-content: center;
  flex-wrap: wrap;
}

.modern-btn {
  position: relative;
  padding: 16px 32px;
  border: none;
  border-radius: 16px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  overflow: hidden;
  min-width: 140px;
}

.btn-content {
  position: relative;
  z-index: 2;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.btn-icon {
  width: 18px;
  height: 18px;
}

.primary-btn {
  background: linear-gradient(135deg, #3b82f6, #1d4ed8);
  color: white;
  box-shadow: 
    0 10px 25px rgba(59, 130, 246, 0.3),
    0 0 0 1px rgba(255, 255, 255, 0.1);
}

.primary-btn:hover {
  transform: translateY(-2px);
  box-shadow: 
    0 15px 35px rgba(59, 130, 246, 0.4),
    0 0 0 1px rgba(255, 255, 255, 0.2);
}

.btn-glow {
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.3), transparent);
  transition: left 0.6s;
}

.primary-btn:hover .btn-glow {
  left: 100%;
}

.secondary-btn {
  background: rgba(255, 255, 255, 0.8);
  color: #475569;
  border: 1px solid rgba(148, 163, 184, 0.3);
  backdrop-filter: blur(10px);
}

.secondary-btn:hover {
  background: rgba(255, 255, 255, 0.95);
  border-color: rgba(59, 130, 246, 0.3);
  color: #3b82f6;
  transform: translateY(-2px);
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.1);
}

/* 现代化动画效果 */
@keyframes float-shape {
  0%, 100% {
    transform: translate(0, 0) rotate(0deg);
  }
  25% {
    transform: translate(10px, -15px) rotate(90deg);
  }
  50% {
    transform: translate(-5px, -10px) rotate(180deg);
  }
  75% {
    transform: translate(-15px, 5px) rotate(270deg);
  }
}

@keyframes icon-pulse {
  0%, 100% {
    transform: scale(1);
  }
  50% {
    transform: scale(1.05);
  }
}

@keyframes checkmark-draw {
  0% {
    opacity: 0;
    transform: scale(0.5) rotate(-45deg);
  }
  50% {
    opacity: 1;
    transform: scale(1.1) rotate(0deg);
  }
  100% {
    opacity: 1;
    transform: scale(1) rotate(0deg);
  }
}

@keyframes pulse-ring {
  0% {
    transform: translate(-50%, -50%) scale(0.8);
    opacity: 1;
  }
  100% {
    transform: translate(-50%, -50%) scale(1.2);
    opacity: 0;
  }
}

@keyframes particle-float {
  0%, 100% {
    transform: translateY(0) scale(1);
    opacity: 0.7;
  }
  50% {
    transform: translateY(-20px) scale(1.2);
    opacity: 1;
  }
}

@keyframes decoration-glow {
  0% {
    box-shadow: 0 0 5px rgba(16, 185, 129, 0.3);
  }
  100% {
    box-shadow: 0 0 20px rgba(16, 185, 129, 0.6);
  }
}

@keyframes badge-glow {
  0%, 100% {
    box-shadow: 0 0 10px rgba(16, 185, 129, 0.2);
  }
  50% {
    box-shadow: 0 0 20px rgba(16, 185, 129, 0.4);
  }
}

@keyframes star-twinkle {
  0%, 100% {
    transform: scale(1) rotate(0deg);
  }
  50% {
    transform: scale(1.2) rotate(180deg);
  }
}

/* 响应式设计 */
@media (max-width: 768px) {
  .modern-empty {
    min-height: 400px;
  }
  
  .empty-main {
    padding: 24px;
  }
  
  .icon-circle {
    width: 80px;
    height: 80px;
  }
  
  .checkmark {
    width: 28px;
    height: 28px;
  }
  
  .ring-1 {
    width: 100px;
    height: 100px;
  }
  
  .ring-2 {
    width: 120px;
    height: 120px;
  }
  
  .ring-3 {
    width: 140px;
    height: 140px;
  }
  
  .particles {
    width: 150px;
    height: 150px;
  }
  
  .title-text {
    font-size: 24px;
  }
  
  .success-subtitle {
    font-size: 14px;
  }
  
  .action-section {
    flex-direction: column;
    align-items: center;
  }
  
  .modern-btn {
    width: 100%;
    max-width: 240px;
    padding: 14px 24px;
    font-size: 14px;
  }
  
  .floating-shape {
    display: none;
  }
}
</style>