<template>
  <div class="exam-taking">
    <!-- 考试头部信息 -->
    <div class="exam-header">
      <div class="exam-info">
        <h1 class="exam-title">{{ exam.title }}</h1>
        <div class="exam-meta">
          <span class="meta-item">
            <el-icon><Clock /></el-icon>
            考试时长：{{ exam.duration }} 分钟
          </span>
          <span class="meta-item">
            <el-icon><Document /></el-icon>
            总题数：{{ exam.totalQuestions }} 题
          </span>
          <span class="meta-item">
            <el-icon><Star /></el-icon>
            总分：{{ exam.totalScore }} 分
          </span>
        </div>
      </div>
      
      <div class="exam-timer">
        <div class="timer-display" :class="{ warning: timeWarning, danger: timeDanger }">
          <el-icon><Timer /></el-icon>
          <span class="time-text">{{ formatTime(remainingTime) }}</span>
        </div>
        <div class="timer-progress">
          <el-progress
            :percentage="timeProgress"
            :color="getProgressColor()"
            :show-text="false"
            :stroke-width="6"
          />
        </div>
      </div>
    </div>
    
    <!-- 答题进度 -->
    <div class="answer-progress">
      <div class="progress-info">
        <span class="progress-text">
          已答题：{{ answeredCount }}/{{ exam.totalQuestions }}
        </span>
        <span class="progress-percentage">
          {{ Math.round((answeredCount / exam.totalQuestions) * 100) }}%
        </span>
      </div>
      <el-progress
        :percentage="(answeredCount / exam.totalQuestions) * 100"
        color="var(--dopamine-blue)"
        :stroke-width="8"
      />
    </div>
    
    <!-- 题目导航 -->
    <div class="question-nav">
      <div class="nav-header">
        <h3>题目导航</h3>
        <el-button
          type="primary"
          size="small"
          @click="showQuestionNav = !showQuestionNav"
        >
          {{ showQuestionNav ? '收起' : '展开' }}
        </el-button>
      </div>
      
      <div v-show="showQuestionNav" class="nav-content">
        <div class="nav-legend">
          <div class="legend-item">
            <div class="legend-dot answered"></div>
            <span>已答题</span>
          </div>
          <div class="legend-item">
            <div class="legend-dot current"></div>
            <span>当前题</span>
          </div>
          <div class="legend-item">
            <div class="legend-dot unanswered"></div>
            <span>未答题</span>
          </div>
        </div>
        
        <div class="nav-grid">
          <div
            v-for="(question, index) in exam.questions"
            :key="question.id"
            class="nav-item"
            :class="{
              current: currentQuestionIndex === index,
              answered: answers[question.id] !== undefined,
              unanswered: answers[question.id] === undefined
            }"
            @click="goToQuestion(index)"
          >
            {{ index + 1 }}
          </div>
        </div>
      </div>
    </div>
    
    <!-- 答题区域 -->
    <div class="answer-area">
      <div class="question-container">
        <!-- 题目信息 -->
        <div class="question-header">
          <div class="question-number">
            第 {{ currentQuestionIndex + 1 }} 题
          </div>
          <div class="question-meta">
            <el-tag :type="getTypeTagType(currentQuestion.type)">
              {{ getTypeText(currentQuestion.type) }}
            </el-tag>
            <span class="question-score">{{ currentQuestion.score }} 分</span>
          </div>
        </div>
        
        <!-- 题目内容 -->
        <div class="question-content">
          <div class="question-text">
            {{ currentQuestion.content }}
          </div>
          
          <!-- 选择题选项 -->
          <div v-if="isChoiceQuestion" class="question-options">
            <div
              v-for="(option, index) in currentQuestion.options"
              :key="index"
              class="option-item"
              :class="{ selected: isOptionSelected(option.id) }"
              @click="selectOption(option.id)"
            >
              <div class="option-radio">
                <el-radio
                  v-if="currentQuestion.type === 'single'"
                  :model-value="answers[currentQuestion.id]"
                  :label="option.id"
                  @change="selectOption(option.id)"
                >
                </el-radio>
                <el-checkbox
                  v-else
                  :model-value="isOptionSelected(option.id)"
                  @change="toggleOption(option.id)"
                >
                </el-checkbox>
              </div>
              <div class="option-content">
                <div class="option-prefix">
                  {{ String.fromCharCode(65 + index) }}.
                </div>
                <div class="option-text">
                  {{ option.text }}
                </div>
              </div>
            </div>
          </div>
          
          <!-- 判断题 -->
          <div v-if="currentQuestion.type === 'judge'" class="judge-options">
            <div
              class="judge-option"
              :class="{ selected: answers[currentQuestion.id] === true }"
              @click="selectJudge(true)"
            >
              <el-radio
                :model-value="answers[currentQuestion.id]"
                :label="true"
                @change="selectJudge(true)"
              >
                正确
              </el-radio>
            </div>
            <div
              class="judge-option"
              :class="{ selected: answers[currentQuestion.id] === false }"
              @click="selectJudge(false)"
            >
              <el-radio
                :model-value="answers[currentQuestion.id]"
                :label="false"
                @change="selectJudge(false)"
              >
                错误
              </el-radio>
            </div>
          </div>
          
          <!-- 填空题 -->
          <div v-if="currentQuestion.type === 'fill'" class="fill-inputs">
            <div
              v-for="(blank, index) in currentQuestion.blanks"
              :key="index"
              class="fill-item"
            >
              <label class="fill-label">第 {{ index + 1 }} 空：</label>
              <el-input
                v-model="fillAnswers[index]"
                placeholder="请输入答案"
                @input="updateFillAnswer"
              />
            </div>
          </div>
          
          <!-- 简答题 -->
          <div v-if="currentQuestion.type === 'essay'" class="essay-input">
            <el-input
              v-model="essayAnswer"
              type="textarea"
              :rows="8"
              placeholder="请输入您的答案..."
              @input="updateEssayAnswer"
            />
            <div class="word-count">
              字数：{{ essayAnswer.length }}
            </div>
          </div>
        </div>
      </div>
      
      <!-- 操作按钮 -->
      <div class="question-actions">
        <el-button
          :disabled="currentQuestionIndex === 0"
          @click="previousQuestion"
        >
          <el-icon><ArrowLeft /></el-icon>
          上一题
        </el-button>
        
        <div class="center-actions">
          <el-button @click="clearAnswer">
            清空答案
          </el-button>
          <el-button type="warning" @click="showSubmitDialog = true">
            <el-icon><Finished /></el-icon>
            提交试卷
          </el-button>
        </div>
        
        <el-button
          :disabled="currentQuestionIndex === exam.questions.length - 1"
          @click="nextQuestion"
        >
          下一题
          <el-icon><ArrowRight /></el-icon>
        </el-button>
      </div>
    </div>
    
    <!-- 提交确认对话框 -->
    <el-dialog
      v-model="showSubmitDialog"
      title="提交试卷"
      width="500px"
      :before-close="handleSubmitDialogClose"
    >
      <div class="submit-summary">
        <div class="summary-item">
          <span class="summary-label">考试名称：</span>
          <span class="summary-value">{{ exam.title }}</span>
        </div>
        <div class="summary-item">
          <span class="summary-label">总题数：</span>
          <span class="summary-value">{{ exam.totalQuestions }} 题</span>
        </div>
        <div class="summary-item">
          <span class="summary-label">已答题：</span>
          <span class="summary-value">{{ answeredCount }} 题</span>
        </div>
        <div class="summary-item">
          <span class="summary-label">未答题：</span>
          <span class="summary-value warning">{{ exam.totalQuestions - answeredCount }} 题</span>
        </div>
        <div class="summary-item">
          <span class="summary-label">剩余时间：</span>
          <span class="summary-value" :class="{ warning: timeWarning, danger: timeDanger }">
            {{ formatTime(remainingTime) }}
          </span>
        </div>
      </div>
      
      <div v-if="exam.totalQuestions - answeredCount > 0" class="submit-warning">
        <el-alert
          title="注意"
          :description="`您还有 ${exam.totalQuestions - answeredCount} 题未答，确定要提交吗？`"
          type="warning"
          :closable="false"
        />
      </div>
      
      <template #footer>
        <el-button @click="showSubmitDialog = false">取消</el-button>
        <el-button type="primary" @click="submitExam" :loading="submitting">
          确认提交
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Clock,
  Document,
  Star,
  Timer,
  ArrowLeft,
  ArrowRight,
  Finished
} from '@element-plus/icons-vue'

const router = useRouter()

// 考试数据
const exam = ref({
  id: '1',
  title: '计算机基础知识测试',
  duration: 120, // 分钟
  totalQuestions: 20,
  totalScore: 100,
  questions: [
    {
      id: '1',
      type: 'single',
      content: '计算机的核心部件是什么？',
      score: 5,
      options: [
        { id: 'a', text: 'CPU' },
        { id: 'b', text: '内存' },
        { id: 'c', text: '硬盘' },
        { id: 'd', text: '显卡' }
      ]
    },
    {
      id: '2',
      type: 'multiple',
      content: '以下哪些是编程语言？',
      score: 5,
      options: [
        { id: 'a', text: 'JavaScript' },
        { id: 'b', text: 'Python' },
        { id: 'c', text: 'HTML' },
        { id: 'd', text: 'Java' }
      ]
    },
    {
      id: '3',
      type: 'judge',
      content: 'HTTP是一种安全的传输协议。',
      score: 3
    },
    {
      id: '4',
      type: 'fill',
      content: '请填写以下空白：TCP/IP协议栈包括____层和____层。',
      score: 6,
      blanks: [null, null]
    },
    {
      id: '5',
      type: 'essay',
      content: '请简述面向对象编程的三大特性。',
      score: 15
    }
  ]
})

// 答题状态
const currentQuestionIndex = ref(0)
const answers = ref<Record<string, any>>({})
const fillAnswers = ref<string[]>([])
const essayAnswer = ref('')
const showQuestionNav = ref(true)
const showSubmitDialog = ref(false)
const submitting = ref(false)

// 计时器
const startTime = ref(Date.now())
const remainingTime = ref(exam.value.duration * 60) // 秒
let timer: NodeJS.Timeout | null = null

// 当前题目
const currentQuestion = computed(() => {
  return exam.value.questions[currentQuestionIndex.value]
})

// 是否为选择题
const isChoiceQuestion = computed(() => {
  return ['single', 'multiple'].includes(currentQuestion.value.type)
})

// 已答题数量
const answeredCount = computed(() => {
  return Object.keys(answers.value).length
})

// 时间进度
const timeProgress = computed(() => {
  const totalTime = exam.value.duration * 60
  return ((totalTime - remainingTime.value) / totalTime) * 100
})

// 时间警告
const timeWarning = computed(() => {
  const totalTime = exam.value.duration * 60
  return remainingTime.value <= totalTime * 0.2 && remainingTime.value > totalTime * 0.1
})

const timeDanger = computed(() => {
  const totalTime = exam.value.duration * 60
  return remainingTime.value <= totalTime * 0.1
})

// 初始化填空题答案
const initFillAnswers = () => {
  if (currentQuestion.value.type === 'fill') {
    const savedAnswers = answers.value[currentQuestion.value.id]
    if (savedAnswers && Array.isArray(savedAnswers)) {
      fillAnswers.value = [...savedAnswers]
    } else {
      fillAnswers.value = new Array(currentQuestion.value.blanks.length).fill('')
    }
  }
}

// 初始化简答题答案
const initEssayAnswer = () => {
  if (currentQuestion.value.type === 'essay') {
    essayAnswer.value = answers.value[currentQuestion.value.id] || ''
  }
}

// 监听题目变化
const watchCurrentQuestion = () => {
  initFillAnswers()
  initEssayAnswer()
}

// 格式化时间
const formatTime = (seconds: number): string => {
  const hours = Math.floor(seconds / 3600)
  const minutes = Math.floor((seconds % 3600) / 60)
  const secs = seconds % 60
  
  if (hours > 0) {
    return `${hours.toString().padStart(2, '0')}:${minutes.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`
  }
  return `${minutes.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`
}

// 获取进度条颜色
const getProgressColor = () => {
  if (timeDanger.value) return 'var(--dopamine-red)'
  if (timeWarning.value) return 'var(--dopamine-orange)'
  return 'var(--dopamine-blue)'
}

// 获取类型文本
const getTypeText = (type: string): string => {
  const textMap: Record<string, string> = {
    single: '单选题',
    multiple: '多选题',
    judge: '判断题',
    fill: '填空题',
    essay: '简答题'
  }
  return textMap[type] || ''
}

// 获取类型标签类型
const getTypeTagType = (type: string): string => {
  const typeMap: Record<string, string> = {
    single: '',
    multiple: 'success',
    judge: 'warning',
    fill: 'danger',
    essay: 'info'
  }
  return typeMap[type] || ''
}

// 跳转到指定题目
const goToQuestion = (index: number) => {
  currentQuestionIndex.value = index
  watchCurrentQuestion()
}

// 上一题
const previousQuestion = () => {
  if (currentQuestionIndex.value > 0) {
    currentQuestionIndex.value--
    watchCurrentQuestion()
  }
}

// 下一题
const nextQuestion = () => {
  if (currentQuestionIndex.value < exam.value.questions.length - 1) {
    currentQuestionIndex.value++
    watchCurrentQuestion()
  }
}

// 选择选项（单选）
const selectOption = (optionId: string) => {
  if (currentQuestion.value.type === 'single') {
    answers.value[currentQuestion.value.id] = optionId
  }
}

// 切换选项（多选）
const toggleOption = (optionId: string) => {
  if (currentQuestion.value.type === 'multiple') {
    const currentAnswers = answers.value[currentQuestion.value.id] || []
    const index = currentAnswers.indexOf(optionId)
    
    if (index > -1) {
      currentAnswers.splice(index, 1)
    } else {
      currentAnswers.push(optionId)
    }
    
    answers.value[currentQuestion.value.id] = [...currentAnswers]
  }
}

// 判断选项是否被选中
const isOptionSelected = (optionId: string): boolean => {
  const answer = answers.value[currentQuestion.value.id]
  if (currentQuestion.value.type === 'single') {
    return answer === optionId
  } else if (currentQuestion.value.type === 'multiple') {
    return Array.isArray(answer) && answer.includes(optionId)
  }
  return false
}

// 选择判断题答案
const selectJudge = (value: boolean) => {
  answers.value[currentQuestion.value.id] = value
}

// 更新填空题答案
const updateFillAnswer = () => {
  answers.value[currentQuestion.value.id] = [...fillAnswers.value]
}

// 更新简答题答案
const updateEssayAnswer = () => {
  answers.value[currentQuestion.value.id] = essayAnswer.value
}

// 清空当前题目答案
const clearAnswer = () => {
  delete answers.value[currentQuestion.value.id]
  
  if (currentQuestion.value.type === 'fill') {
    fillAnswers.value = new Array(currentQuestion.value.blanks.length).fill('')
  } else if (currentQuestion.value.type === 'essay') {
    essayAnswer.value = ''
  }
}

// 处理提交对话框关闭
const handleSubmitDialogClose = (done: () => void) => {
  ElMessageBox.confirm('确定要取消提交吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    done()
  }).catch(() => {
    // 取消关闭
  })
}

// 提交考试
const submitExam = async () => {
  try {
    submitting.value = true
    
    // 模拟提交
    await new Promise(resolve => setTimeout(resolve, 2000))
    
    ElMessage.success('试卷提交成功！')
    showSubmitDialog.value = false
    
    // 跳转到考试列表页面
    router.push('/student/exams')
  } catch (error) {
    ElMessage.error('提交失败，请重试')
  } finally {
    submitting.value = false
  }
}

// 启动计时器
const startTimer = () => {
  timer = setInterval(() => {
    remainingTime.value--
    
    if (remainingTime.value <= 0) {
      // 时间到，自动提交
      ElMessage.warning('考试时间已到，系统将自动提交试卷')
      submitExam()
    }
  }, 1000)
}

// 停止计时器
const stopTimer = () => {
  if (timer) {
    clearInterval(timer)
    timer = null
  }
}

// 页面离开确认
const handleBeforeUnload = (event: BeforeUnloadEvent) => {
  event.preventDefault()
  event.returnValue = '确定要离开吗？您的答题进度可能会丢失。'
}

onMounted(() => {
  startTimer()
  watchCurrentQuestion()
  window.addEventListener('beforeunload', handleBeforeUnload)
})

onUnmounted(() => {
  stopTimer()
  window.removeEventListener('beforeunload', handleBeforeUnload)
})
</script>

<style scoped>
.exam-taking {
  min-height: 100vh;
  background: var(--bg-primary);
  padding: var(--spacing-lg);
}

.exam-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  background: white;
  padding: var(--spacing-lg);
  border-radius: var(--radius-lg);
  margin-bottom: var(--spacing-lg);
  box-shadow: var(--shadow-sm);
}

.exam-info {
  flex: 1;
}

.exam-title {
  font-size: 24px;
  font-weight: 700;
  color: var(--text-primary);
  margin: 0 0 var(--spacing-md) 0;
}

.exam-meta {
  display: flex;
  gap: var(--spacing-lg);
}

.meta-item {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  font-size: 14px;
  color: var(--text-secondary);
}

.exam-timer {
  text-align: center;
  min-width: 200px;
}

.timer-display {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--spacing-sm);
  font-size: 20px;
  font-weight: 700;
  color: var(--dopamine-blue);
  margin-bottom: var(--spacing-md);
}

.timer-display.warning {
  color: var(--dopamine-orange);
}

.timer-display.danger {
  color: var(--dopamine-red);
  animation: pulse 1s infinite;
}

.time-text {
  font-family: 'Courier New', monospace;
}

.answer-progress {
  background: white;
  padding: var(--spacing-lg);
  border-radius: var(--radius-lg);
  margin-bottom: var(--spacing-lg);
  box-shadow: var(--shadow-sm);
}

.progress-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--spacing-md);
}

.progress-text {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
}

.progress-percentage {
  font-size: 14px;
  color: var(--dopamine-blue);
  font-weight: 600;
}

.question-nav {
  background: white;
  padding: var(--spacing-lg);
  border-radius: var(--radius-lg);
  margin-bottom: var(--spacing-lg);
  box-shadow: var(--shadow-sm);
}

.nav-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--spacing-md);
}

.nav-header h3 {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
}

.nav-legend {
  display: flex;
  gap: var(--spacing-lg);
  margin-bottom: var(--spacing-md);
}

.legend-item {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  font-size: 14px;
  color: var(--text-secondary);
}

.legend-dot {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  border: 2px solid;
}

.legend-dot.answered {
  background: var(--dopamine-green);
  border-color: var(--dopamine-green);
}

.legend-dot.current {
  background: var(--dopamine-blue);
  border-color: var(--dopamine-blue);
}

.legend-dot.unanswered {
  background: white;
  border-color: var(--border-color);
}

.nav-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(40px, 1fr));
  gap: var(--spacing-sm);
}

.nav-item {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--radius-md);
  border: 2px solid var(--border-color);
  background: white;
  cursor: pointer;
  font-weight: 600;
  transition: all 0.3s ease;
}

.nav-item:hover {
  border-color: var(--dopamine-blue);
  transform: translateY(-2px);
}

.nav-item.answered {
  background: var(--dopamine-green);
  border-color: var(--dopamine-green);
  color: white;
}

.nav-item.current {
  background: var(--dopamine-blue);
  border-color: var(--dopamine-blue);
  color: white;
  transform: scale(1.1);
}

.answer-area {
  background: white;
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-sm);
  overflow: hidden;
}

.question-container {
  padding: var(--spacing-xl);
}

.question-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--spacing-lg);
  padding-bottom: var(--spacing-lg);
  border-bottom: 1px solid var(--border-color);
}

.question-number {
  font-size: 20px;
  font-weight: 700;
  color: var(--dopamine-blue);
}

.question-meta {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
}

.question-score {
  font-size: 14px;
  font-weight: 600;
  color: var(--dopamine-orange);
  background: var(--dopamine-orange-light);
  padding: 4px 8px;
  border-radius: var(--radius-sm);
}

.question-content {
  margin-bottom: var(--spacing-xl);
}

.question-text {
  font-size: 16px;
  line-height: 1.6;
  color: var(--text-primary);
  margin-bottom: var(--spacing-lg);
  padding: var(--spacing-lg);
  background: var(--bg-secondary);
  border-radius: var(--radius-md);
  border-left: 4px solid var(--dopamine-blue);
}

.question-options,
.judge-options {
  margin-top: var(--spacing-lg);
}

.option-item {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
  padding: var(--spacing-md);
  margin-bottom: var(--spacing-sm);
  border-radius: var(--radius-md);
  border: 2px solid var(--border-color);
  background: white;
  cursor: pointer;
  transition: all 0.3s ease;
}

.option-item:hover {
  border-color: var(--dopamine-blue);
  background: var(--dopamine-blue-light);
}

.option-item.selected {
  border-color: var(--dopamine-blue);
  background: var(--dopamine-blue-light);
}

.option-radio {
  flex-shrink: 0;
}

.option-content {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
  flex: 1;
}

.option-prefix {
  font-weight: 600;
  color: var(--text-primary);
  min-width: 24px;
}

.option-text {
  flex: 1;
  font-size: 14px;
  line-height: 1.5;
  color: var(--text-primary);
}

.judge-options {
  display: flex;
  gap: var(--spacing-lg);
}

.judge-option {
  flex: 1;
  padding: var(--spacing-lg);
  border: 2px solid var(--border-color);
  border-radius: var(--radius-md);
  text-align: center;
  cursor: pointer;
  transition: all 0.3s ease;
}

.judge-option:hover {
  border-color: var(--dopamine-blue);
  background: var(--dopamine-blue-light);
}

.judge-option.selected {
  border-color: var(--dopamine-blue);
  background: var(--dopamine-blue-light);
}

.fill-inputs {
  margin-top: var(--spacing-lg);
}

.fill-item {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
  margin-bottom: var(--spacing-md);
}

.fill-label {
  font-weight: 600;
  color: var(--text-primary);
  min-width: 80px;
}

.essay-input {
  margin-top: var(--spacing-lg);
}

.word-count {
  text-align: right;
  font-size: 12px;
  color: var(--text-secondary);
  margin-top: var(--spacing-sm);
}

.question-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--spacing-lg) var(--spacing-xl);
  background: var(--bg-secondary);
  border-top: 1px solid var(--border-color);
}

.center-actions {
  display: flex;
  gap: var(--spacing-md);
}

.submit-summary {
  margin-bottom: var(--spacing-lg);
}

.summary-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--spacing-sm) 0;
  border-bottom: 1px solid var(--border-color);
}

.summary-item:last-child {
  border-bottom: none;
}

.summary-label {
  font-weight: 500;
  color: var(--text-secondary);
}

.summary-value {
  font-weight: 600;
  color: var(--text-primary);
}

.summary-value.warning {
  color: var(--dopamine-orange);
}

.summary-value.danger {
  color: var(--dopamine-red);
}

.submit-warning {
  margin-top: var(--spacing-lg);
}

@keyframes pulse {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0.5;
  }
}

/* 响应式设计 */
@media (max-width: 768px) {
  .exam-taking {
    padding: var(--spacing-md);
  }
  
  .exam-header {
    flex-direction: column;
    gap: var(--spacing-lg);
  }
  
  .exam-meta {
    flex-direction: column;
    gap: var(--spacing-sm);
  }
  
  .exam-timer {
    min-width: auto;
    width: 100%;
  }
  
  .nav-grid {
    grid-template-columns: repeat(auto-fill, minmax(35px, 1fr));
  }
  
  .nav-item {
    width: 35px;
    height: 35px;
  }
  
  .question-header {
    flex-direction: column;
    align-items: flex-start;
    gap: var(--spacing-md);
  }
  
  .question-actions {
    flex-direction: column;
    gap: var(--spacing-md);
  }
  
  .center-actions {
    order: -1;
  }
  
  .judge-options {
    flex-direction: column;
  }
  
  .fill-item {
    flex-direction: column;
    align-items: flex-start;
    gap: var(--spacing-sm);
  }
  
  .fill-label {
    min-width: auto;
  }
}
</style>