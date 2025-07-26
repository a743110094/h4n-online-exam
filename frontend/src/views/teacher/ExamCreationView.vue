<template>
  <div class="exam-creation">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-content">
        <h1 class="page-title">
          <el-icon><DocumentAdd /></el-icon>
          {{ isEdit ? '编辑考试' : '创建考试' }}
        </h1>
        <p class="page-description">
          {{ isEdit ? '修改考试信息和题目配置' : '设置考试基本信息，选择题目并配置考试参数' }}
        </p>
      </div>
      
      <div class="header-actions">
        <el-button @click="saveDraft" :loading="saving">
          <el-icon><Document /></el-icon>
          保存草稿
        </el-button>
        <el-button type="primary" @click="publishExam" :loading="publishing">
          <el-icon><Promotion /></el-icon>
          {{ isEdit ? '更新考试' : '发布考试' }}
        </el-button>
      </div>
    </div>
    
    <!-- 考试创建步骤 -->
    <div class="creation-steps">
      <el-steps :active="currentStep" align-center>
        <el-step title="基本信息" description="设置考试基本信息" />
        <el-step title="选择题目" description="从题库中选择题目" />
        <el-step title="考试配置" description="配置考试参数" />
        <el-step title="预览发布" description="预览并发布考试" />
      </el-steps>
    </div>
    
    <!-- 步骤内容 -->
    <div class="step-content">
      <!-- 步骤1：基本信息 -->
      <div v-show="currentStep === 0" class="step-panel">
        <div class="panel-header">
          <h2>考试基本信息</h2>
          <p>请填写考试的基本信息</p>
        </div>
        
        <el-form
          ref="basicFormRef"
          :model="examForm"
          :rules="basicRules"
          label-width="120px"
          class="basic-form"
        >
          <el-row :gutter="24">
            <el-col :span="12">
              <el-form-item label="考试名称" prop="title">
                <el-input
                  v-model="examForm.title"
                  placeholder="请输入考试名称"
                  maxlength="100"
                  show-word-limit
                />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="考试科目" prop="subject">
                <el-select
                  v-model="examForm.subject"
                  placeholder="请选择考试科目"
                  style="width: 100%"
                >
                  <el-option
                    v-for="subject in subjects"
                    :key="subject.value"
                    :label="subject.label"
                    :value="subject.value"
                  />
                </el-select>
              </el-form-item>
            </el-col>
          </el-row>
          
          <el-row :gutter="24">
            <el-col :span="12">
              <el-form-item label="考试时长" prop="duration">
                <el-input-number
                  v-model="examForm.duration"
                  :min="1"
                  :max="300"
                  placeholder="分钟"
                  style="width: 100%"
                />
                <span class="form-tip">建议时长：60-120分钟</span>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="总分" prop="totalScore">
                <el-input-number
                  v-model="examForm.totalScore"
                  :min="1"
                  :max="1000"
                  placeholder="分"
                  style="width: 100%"
                />
                <span class="form-tip">当前已选题目总分：{{ selectedQuestionsScore }}</span>
              </el-form-item>
            </el-col>
          </el-row>
          
          <el-row :gutter="24">
            <el-col :span="12">
              <el-form-item label="开始时间" prop="startTime">
                <el-date-picker
                  v-model="examForm.startTime"
                  type="datetime"
                  placeholder="选择开始时间"
                  style="width: 100%"
                  :disabled-date="disabledDate"
                />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="结束时间" prop="endTime">
                <el-date-picker
                  v-model="examForm.endTime"
                  type="datetime"
                  placeholder="选择结束时间"
                  style="width: 100%"
                  :disabled-date="disabledDate"
                />
              </el-form-item>
            </el-col>
          </el-row>
          
          <el-form-item label="考试描述" prop="description">
            <el-input
              v-model="examForm.description"
              type="textarea"
              :rows="4"
              placeholder="请输入考试描述、注意事项等"
              maxlength="500"
              show-word-limit
            />
          </el-form-item>
          
          <el-form-item label="参考学生">
            <el-select
              v-model="examForm.studentIds"
              multiple
              placeholder="选择参考学生（不选择则所有学生可参加）"
              style="width: 100%"
              filterable
            >
              <el-option
                v-for="student in students"
                :key="student.id"
                :label="`${student.name} (${student.studentId})`"
                :value="student.id"
              />
            </el-select>
          </el-form-item>
        </el-form>
      </div>
      
      <!-- 步骤2：选择题目 -->
      <div v-show="currentStep === 1" class="step-panel">
        <div class="panel-header">
          <h2>选择题目</h2>
          <p>从题库中选择题目组成试卷</p>
        </div>
        
        <!-- 题目筛选 -->
        <div class="question-filters">
          <el-row :gutter="16">
            <el-col :span="6">
              <el-select
                v-model="questionFilters.subject"
                placeholder="科目"
                clearable
                @change="loadQuestions"
              >
                <el-option
                  v-for="subject in subjects"
                  :key="subject.value"
                  :label="subject.label"
                  :value="subject.value"
                />
              </el-select>
            </el-col>
            <el-col :span="6">
              <el-select
                v-model="questionFilters.type"
                placeholder="题目类型"
                clearable
                @change="loadQuestions"
              >
                <el-option label="单选题" value="single" />
                <el-option label="多选题" value="multiple" />
                <el-option label="判断题" value="judge" />
                <el-option label="填空题" value="fill" />
                <el-option label="简答题" value="essay" />
              </el-select>
            </el-col>
            <el-col :span="6">
              <el-select
                v-model="questionFilters.difficulty"
                placeholder="难度"
                clearable
                @change="loadQuestions"
              >
                <el-option label="简单" value="easy" />
                <el-option label="中等" value="medium" />
                <el-option label="困难" value="hard" />
              </el-select>
            </el-col>
            <el-col :span="6">
              <el-input
                v-model="questionFilters.keyword"
                placeholder="搜索题目内容"
                clearable
                @input="debounceSearch"
              >
                <template #prefix>
                  <el-icon><Search /></el-icon>
                </template>
              </el-input>
            </el-col>
          </el-row>
        </div>
        
        <!-- 选题统计 -->
        <div class="selection-stats">
          <div class="stats-cards">
            <div class="stat-card">
              <div class="stat-number">{{ selectedQuestions.length }}</div>
              <div class="stat-label">已选题目</div>
            </div>
            <div class="stat-card">
              <div class="stat-number">{{ selectedQuestionsScore }}</div>
              <div class="stat-label">总分</div>
            </div>
            <div class="stat-card">
              <div class="stat-number">{{ getQuestionTypeCount('single') }}</div>
              <div class="stat-label">单选题</div>
            </div>
            <div class="stat-card">
              <div class="stat-number">{{ getQuestionTypeCount('multiple') }}</div>
              <div class="stat-label">多选题</div>
            </div>
            <div class="stat-card">
              <div class="stat-number">{{ getQuestionTypeCount('judge') }}</div>
              <div class="stat-label">判断题</div>
            </div>
            <div class="stat-card">
              <div class="stat-number">{{ getQuestionTypeCount('fill') }}</div>
              <div class="stat-label">填空题</div>
            </div>
            <div class="stat-card">
              <div class="stat-number">{{ getQuestionTypeCount('essay') }}</div>
              <div class="stat-label">简答题</div>
            </div>
          </div>
          
          <div class="selection-actions">
            <el-button @click="clearSelection">
              清空选择
            </el-button>
            <el-button type="primary" @click="showSelectedQuestions = true">
              查看已选题目 ({{ selectedQuestions.length }})
            </el-button>
          </div>
        </div>
        
        <!-- 题目列表 -->
        <div class="question-list">
          <div
            v-for="question in questions"
            :key="question.id"
            class="question-item"
            :class="{ selected: isQuestionSelected(question.id) }"
          >
            <div class="question-checkbox">
              <el-checkbox
                :model-value="isQuestionSelected(question.id)"
                @change="toggleQuestionSelection(question)"
              />
            </div>
            
            <div class="question-content">
              <div class="question-header">
                <div class="question-meta">
                  <el-tag :type="getTypeTagType(question.type)" size="small">
                    {{ getTypeText(question.type) }}
                  </el-tag>
                  <el-tag :type="getDifficultyTagType(question.difficulty)" size="small">
                    {{ getDifficultyText(question.difficulty) }}
                  </el-tag>
                  <span class="question-score">{{ question.score }}分</span>
                </div>
                
                <div class="question-actions">
                  <el-button
                    type="primary"
                    link
                    @click="previewQuestion(question)"
                  >
                    预览
                  </el-button>
                </div>
              </div>
              
              <div class="question-text">
                {{ question.content }}
              </div>
              
              <div class="question-info">
                <span class="info-item">科目：{{ question.subject?.name || '未分类' }}</span>
                <span class="info-item">知识点：{{ question.knowledgePoint }}</span>
                <span class="info-item">使用次数：{{ question.usageCount || 0 }}</span>
              </div>
            </div>
          </div>
          
          <!-- 分页 -->
          <div class="pagination-wrapper">
            <el-pagination
              v-model:current-page="questionPagination.page"
              v-model:page-size="questionPagination.size"
              :total="questionPagination.total"
              :page-sizes="[10, 20, 50]"
              layout="total, sizes, prev, pager, next, jumper"
              @current-change="loadQuestions"
              @size-change="loadQuestions"
            />
          </div>
        </div>
      </div>
      
      <!-- 步骤3：考试配置 -->
      <div v-show="currentStep === 2" class="step-panel">
        <div class="panel-header">
          <h2>考试配置</h2>
          <p>配置考试的高级参数和规则</p>
        </div>
        
        <el-form
          ref="configFormRef"
          :model="examConfig"
          label-width="150px"
          class="config-form"
        >
          <div class="config-section">
            <h3>答题设置</h3>
            
            <el-form-item label="题目顺序">
              <el-radio-group v-model="examConfig.questionOrder">
                <el-radio label="fixed">固定顺序</el-radio>
                <el-radio label="random">随机顺序</el-radio>
              </el-radio-group>
            </el-form-item>
            
            <el-form-item label="选项顺序">
              <el-radio-group v-model="examConfig.optionOrder">
                <el-radio label="fixed">固定顺序</el-radio>
                <el-radio label="random">随机顺序</el-radio>
              </el-radio-group>
            </el-form-item>
            
            <el-form-item label="答题模式">
              <el-radio-group v-model="examConfig.answerMode">
                <el-radio label="sequential">顺序答题</el-radio>
                <el-radio label="free">自由答题</el-radio>
              </el-radio-group>
              <div class="form-tip">
                顺序答题：必须按顺序完成，不能跳题<br>
                自由答题：可以自由跳转题目
              </div>
            </el-form-item>
            
            <el-form-item label="允许回看">
              <el-switch v-model="examConfig.allowReview" />
              <span class="form-tip">是否允许学生回看已答题目</span>
            </el-form-item>
          </div>
          
          <div class="config-section">
            <h3>防作弊设置</h3>
            
            <el-form-item label="切屏检测">
              <el-switch v-model="examConfig.detectTabSwitch" />
              <span class="form-tip">检测学生是否切换浏览器标签页</span>
            </el-form-item>
            
            <el-form-item label="复制粘贴">
              <el-switch v-model="examConfig.allowCopyPaste" />
              <span class="form-tip">是否允许复制粘贴操作</span>
            </el-form-item>
            
            <el-form-item label="右键菜单">
              <el-switch v-model="examConfig.allowRightClick" />
              <span class="form-tip">是否允许右键菜单</span>
            </el-form-item>
            
            <el-form-item label="全屏模式">
              <el-switch v-model="examConfig.fullScreenMode" />
              <span class="form-tip">强制全屏答题</span>
            </el-form-item>
          </div>
          
          <div class="config-section">
            <h3>成绩设置</h3>
            
            <el-form-item label="立即显示成绩">
              <el-switch v-model="examConfig.showScoreImmediately" />
              <span class="form-tip">考试结束后立即显示成绩</span>
            </el-form-item>
            
            <el-form-item label="显示答案解析">
              <el-switch v-model="examConfig.showAnswerAnalysis" />
              <span class="form-tip">是否显示题目答案和解析</span>
            </el-form-item>
            
            <el-form-item label="及格分数">
              <el-input-number
                v-model="examConfig.passingScore"
                :min="0"
                :max="examForm.totalScore"
                placeholder="分"
              />
              <span class="form-tip">设置考试及格分数线</span>
            </el-form-item>
          </div>
        </el-form>
      </div>
      
      <!-- 步骤4：预览发布 -->
      <div v-show="currentStep === 3" class="step-panel">
        <div class="panel-header">
          <h2>预览发布</h2>
          <p>确认考试信息无误后发布考试</p>
        </div>
        
        <div class="exam-preview">
          <div class="preview-section">
            <h3>考试信息</h3>
            <div class="info-grid">
              <div class="info-item">
                <span class="info-label">考试名称：</span>
                <span class="info-value">{{ examForm.title }}</span>
              </div>
              <div class="info-item">
                <span class="info-label">考试科目：</span>
                <span class="info-value">{{ getSubjectLabel(examForm.subject) }}</span>
              </div>
              <div class="info-item">
                <span class="info-label">考试时长：</span>
                <span class="info-value">{{ examForm.duration }} 分钟</span>
              </div>
              <div class="info-item">
                <span class="info-label">总分：</span>
                <span class="info-value">{{ examForm.totalScore }} 分</span>
              </div>
              <div class="info-item">
                <span class="info-label">开始时间：</span>
                <span class="info-value">{{ formatDateTime(examForm.startTime) }}</span>
              </div>
              <div class="info-item">
                <span class="info-label">结束时间：</span>
                <span class="info-value">{{ formatDateTime(examForm.endTime) }}</span>
              </div>
              <div class="info-item">
                <span class="info-label">参考学生：</span>
                <span class="info-value">
                  {{ examForm.studentIds.length > 0 ? `${examForm.studentIds.length} 人` : '所有学生' }}
                </span>
              </div>
            </div>
          </div>
          
          <div class="preview-section">
            <h3>题目统计</h3>
            <div class="question-stats">
              <div class="stat-item">
                <span class="stat-label">总题数：</span>
                <span class="stat-value">{{ selectedQuestions.length }} 题</span>
              </div>
              <div class="stat-item">
                <span class="stat-label">单选题：</span>
                <span class="stat-value">{{ getQuestionTypeCount('single') }} 题</span>
              </div>
              <div class="stat-item">
                <span class="stat-label">多选题：</span>
                <span class="stat-value">{{ getQuestionTypeCount('multiple') }} 题</span>
              </div>
              <div class="stat-item">
                <span class="stat-label">判断题：</span>
                <span class="stat-value">{{ getQuestionTypeCount('judge') }} 题</span>
              </div>
              <div class="stat-item">
                <span class="stat-label">填空题：</span>
                <span class="stat-value">{{ getQuestionTypeCount('fill') }} 题</span>
              </div>
              <div class="stat-item">
                <span class="stat-label">简答题：</span>
                <span class="stat-value">{{ getQuestionTypeCount('essay') }} 题</span>
              </div>
            </div>
          </div>
          
          <div class="preview-section">
            <h3>考试配置</h3>
            <div class="config-summary">
              <div class="config-item">
                <span class="config-label">题目顺序：</span>
                <span class="config-value">{{ examConfig.questionOrder === 'fixed' ? '固定顺序' : '随机顺序' }}</span>
              </div>
              <div class="config-item">
                <span class="config-label">答题模式：</span>
                <span class="config-value">{{ examConfig.answerMode === 'sequential' ? '顺序答题' : '自由答题' }}</span>
              </div>
              <div class="config-item">
                <span class="config-label">防作弊：</span>
                <span class="config-value">
                  {{ examConfig.detectTabSwitch ? '启用' : '禁用' }}切屏检测，
                  {{ examConfig.fullScreenMode ? '启用' : '禁用' }}全屏模式
                </span>
              </div>
              <div class="config-item">
                <span class="config-label">成绩显示：</span>
                <span class="config-value">
                  {{ examConfig.showScoreImmediately ? '立即显示' : '延迟显示' }}成绩，
                  {{ examConfig.showAnswerAnalysis ? '显示' : '不显示' }}答案解析
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 步骤导航 -->
    <div class="step-navigation">
      <el-button
        :disabled="currentStep === 0"
        @click="previousStep"
      >
        上一步
      </el-button>
      
      <el-button
        v-if="currentStep < 3"
        type="primary"
        @click="nextStep"
        :disabled="!canProceedToNextStep"
      >
        下一步
      </el-button>
      
      <el-button
        v-if="currentStep === 3"
        type="primary"
        @click="publishExam"
        :loading="publishing"
      >
        {{ isEdit ? '更新考试' : '发布考试' }}
      </el-button>
    </div>
    
    <!-- 已选题目对话框 -->
    <el-dialog
      v-model="showSelectedQuestions"
      title="已选题目"
      width="80%"
      :before-close="handleSelectedQuestionsClose"
    >
      <div class="selected-questions-list">
        <div
          v-for="(question, index) in selectedQuestions"
          :key="question.id"
          class="selected-question-item"
        >
          <div class="question-index">{{ index + 1 }}</div>
          <div class="question-info">
            <div class="question-meta">
              <el-tag :type="getTypeTagType(question.type)" size="small">
                {{ getTypeText(question.type) }}
              </el-tag>
              <span class="question-score">{{ question.score }}分</span>
            </div>
            <div class="question-content">{{ question.content }}</div>
          </div>
          <div class="question-actions">
            <el-button
              type="danger"
              link
              @click="removeSelectedQuestion(question.id)"
            >
              移除
            </el-button>
          </div>
        </div>
      </div>
    </el-dialog>
    
    <!-- 题目预览对话框 -->
    <el-dialog
      v-model="showQuestionPreview"
      title="题目预览"
      width="60%"
    >
      <QuestionPreview
        v-if="previewQuestionData"
        :question="previewQuestionData"
      />
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  DocumentAdd,
  Document,
  Promotion,
  Search
} from '@element-plus/icons-vue'
import QuestionPreview from '@/components/QuestionPreview.vue'
import { debounce } from 'lodash'

const router = useRouter()
const route = useRoute()

// 是否为编辑模式
const isEdit = computed(() => !!route.params.id)

// 当前步骤
const currentStep = ref(0)

// 表单引用
const basicFormRef = ref()
const configFormRef = ref()

// 加载状态
const saving = ref(false)
const publishing = ref(false)

// 考试表单数据
const examForm = ref({
  title: '',
  subject: '',
  duration: 120,
  totalScore: 100,
  startTime: null,
  endTime: null,
  description: '',
  studentIds: []
})

// 考试配置
const examConfig = ref({
  questionOrder: 'fixed',
  optionOrder: 'fixed',
  answerMode: 'free',
  allowReview: true,
  detectTabSwitch: true,
  allowCopyPaste: false,
  allowRightClick: false,
  fullScreenMode: false,
  showScoreImmediately: true,
  showAnswerAnalysis: true,
  passingScore: 60
})

// 基本信息表单验证规则
const basicRules = {
  title: [
    { required: true, message: '请输入考试名称', trigger: 'blur' },
    { min: 2, max: 100, message: '长度在 2 到 100 个字符', trigger: 'blur' }
  ],
  subject: [
    { required: true, message: '请选择考试科目', trigger: 'change' }
  ],
  duration: [
    { required: true, message: '请输入考试时长', trigger: 'blur' },
    { type: 'number', min: 1, max: 300, message: '时长应在 1-300 分钟之间', trigger: 'blur' }
  ],
  totalScore: [
    { required: true, message: '请输入总分', trigger: 'blur' },
    { type: 'number', min: 1, max: 1000, message: '总分应在 1-1000 分之间', trigger: 'blur' }
  ],
  startTime: [
    { required: true, message: '请选择开始时间', trigger: 'change' }
  ],
  endTime: [
    { required: true, message: '请选择结束时间', trigger: 'change' }
  ]
}

// 科目选项
const subjects = ref([
  { label: '计算机基础', value: 'computer' },
  { label: '数学', value: 'math' },
  { label: '英语', value: 'english' },
  { label: '物理', value: 'physics' },
  { label: '化学', value: 'chemistry' }
])

// 学生列表
const students = ref([
  { id: '1', name: '张三', studentId: '2021001' },
  { id: '2', name: '李四', studentId: '2021002' },
  { id: '3', name: '王五', studentId: '2021003' }
])

// 题目相关数据
const questions = ref([])
const selectedQuestions = ref([])
const questionFilters = ref({
  subject: '',
  type: '',
  difficulty: '',
  keyword: ''
})
const questionPagination = ref({
  page: 1,
  size: 10,
  total: 0
})

// 对话框状态
const showSelectedQuestions = ref(false)
const showQuestionPreview = ref(false)
const previewQuestionData = ref(null)

// 已选题目总分
const selectedQuestionsScore = computed(() => {
  return selectedQuestions.value.reduce((total, question) => total + question.score, 0)
})

// 是否可以进入下一步
const canProceedToNextStep = computed(() => {
  switch (currentStep.value) {
    case 0:
      return examForm.value.title && examForm.value.subject && examForm.value.duration && examForm.value.startTime && examForm.value.endTime
    case 1:
      return selectedQuestions.value.length > 0
    case 2:
      return true
    default:
      return false
  }
})

// 禁用日期
const disabledDate = (time: Date) => {
  return time.getTime() < Date.now() - 8.64e7 // 不能选择今天之前的日期
}

// 获取科目标签
const getSubjectLabel = (value: string) => {
  const subject = subjects.value.find(s => s.value === value)
  return subject ? subject.label : value
}

// 格式化日期时间
const formatDateTime = (date: Date | null) => {
  if (!date) return ''
  return new Intl.DateTimeFormat('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  }).format(date)
}

// 获取题目类型数量
const getQuestionTypeCount = (type: string) => {
  return selectedQuestions.value.filter(q => q.type === type).length
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

// 获取难度文本
const getDifficultyText = (difficulty: string): string => {
  const textMap: Record<string, string> = {
    easy: '简单',
    medium: '中等',
    hard: '困难'
  }
  return textMap[difficulty] || ''
}

// 获取难度标签类型
const getDifficultyTagType = (difficulty: string): string => {
  const typeMap: Record<string, string> = {
    easy: 'success',
    medium: 'warning',
    hard: 'danger'
  }
  return typeMap[difficulty] || ''
}

// 加载题目列表
const loadQuestions = async () => {
  try {
    // 模拟API调用
    const mockQuestions = [
      {
        id: '1',
        type: 'single',
        content: '计算机的核心部件是什么？',
        subject: 'computer',
        knowledgePoint: '计算机硬件',
        difficulty: 'easy',
        score: 5,
        usageCount: 15
      },
      {
        id: '2',
        type: 'multiple',
        content: '以下哪些是编程语言？',
        subject: 'computer',
        knowledgePoint: '编程语言',
        difficulty: 'medium',
        score: 8,
        usageCount: 12
      },
      {
        id: '3',
        type: 'judge',
        content: 'HTTP是一种安全的传输协议。',
        subject: 'computer',
        knowledgePoint: '网络协议',
        difficulty: 'medium',
        score: 3,
        usageCount: 20
      }
    ]
    
    questions.value = mockQuestions
    questionPagination.value.total = mockQuestions.length
  } catch (error) {
    ElMessage.error('加载题目失败')
  }
}

// 防抖搜索
const debounceSearch = debounce(() => {
  loadQuestions()
}, 500)

// 判断题目是否已选择
const isQuestionSelected = (questionId: string) => {
  return selectedQuestions.value.some(q => q.id === questionId)
}

// 切换题目选择状态
const toggleQuestionSelection = (question: any) => {
  const index = selectedQuestions.value.findIndex(q => q.id === question.id)
  if (index > -1) {
    selectedQuestions.value.splice(index, 1)
  } else {
    selectedQuestions.value.push(question)
  }
}

// 清空选择
const clearSelection = () => {
  ElMessageBox.confirm('确定要清空所有已选题目吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    selectedQuestions.value = []
    ElMessage.success('已清空选择')
  })
}

// 移除已选题目
const removeSelectedQuestion = (questionId: string) => {
  const index = selectedQuestions.value.findIndex(q => q.id === questionId)
  if (index > -1) {
    selectedQuestions.value.splice(index, 1)
  }
}

// 预览题目
const previewQuestion = (question: any) => {
  previewQuestionData.value = question
  showQuestionPreview.value = true
}

// 处理已选题目对话框关闭
const handleSelectedQuestionsClose = (done: () => void) => {
  done()
}

// 上一步
const previousStep = () => {
  if (currentStep.value > 0) {
    currentStep.value--
  }
}

// 下一步
const nextStep = async () => {
  if (currentStep.value === 0) {
    // 验证基本信息
    const valid = await basicFormRef.value?.validate()
    if (!valid) return
    
    // 验证时间逻辑
    if (examForm.value.endTime <= examForm.value.startTime) {
      ElMessage.error('结束时间必须晚于开始时间')
      return
    }
  }
  
  if (currentStep.value < 3) {
    currentStep.value++
  }
}

// 保存草稿
const saveDraft = async () => {
  try {
    saving.value = true
    
    // 模拟保存
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    ElMessage.success('草稿保存成功')
  } catch (error) {
    ElMessage.error('保存失败')
  } finally {
    saving.value = false
  }
}

// 发布考试
const publishExam = async () => {
  try {
    publishing.value = true
    
    // 最终验证
    if (selectedQuestions.value.length === 0) {
      ElMessage.error('请至少选择一道题目')
      return
    }
    
    if (selectedQuestionsScore.value !== examForm.value.totalScore) {
      const result = await ElMessageBox.confirm(
        `已选题目总分(${selectedQuestionsScore.value})与设定总分(${examForm.value.totalScore})不一致，是否继续？`,
        '分数不匹配',
        {
          confirmButtonText: '继续发布',
          cancelButtonText: '取消',
          type: 'warning'
        }
      )
      if (!result) return
    }
    
    // 模拟发布
    await new Promise(resolve => setTimeout(resolve, 2000))
    
    ElMessage.success(isEdit.value ? '考试更新成功' : '考试发布成功')
    router.push('/teacher/exam')
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('发布失败')
    }
  } finally {
    publishing.value = false
  }
}

// 监听科目变化，自动筛选题目
watch(() => examForm.value.subject, (newSubject) => {
  if (newSubject && currentStep.value === 1) {
    questionFilters.value.subject = newSubject
    loadQuestions()
  }
})

// 监听已选题目变化，更新总分
watch(selectedQuestions, (newQuestions) => {
  const totalScore = newQuestions.reduce((total, question) => total + question.score, 0)
  if (totalScore !== examForm.value.totalScore && currentStep.value > 0) {
    examForm.value.totalScore = totalScore
  }
}, { deep: true })

onMounted(() => {
  if (isEdit.value) {
    // 加载考试数据
    // loadExamData(route.params.id)
  }
  loadQuestions()
})
</script>

<style scoped>
.exam-creation {
  padding: var(--spacing-lg);
  background: var(--bg-primary);
  min-height: 100vh;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  background: white;
  padding: var(--spacing-xl);
  border-radius: var(--radius-lg);
  margin-bottom: var(--spacing-lg);
  box-shadow: var(--shadow-sm);
}

.header-content {
  flex: 1;
}

.page-title {
  font-size: 28px;
  font-weight: 700;
  color: var(--text-primary);
  margin: 0 0 var(--spacing-sm) 0;
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
}

.page-description {
  font-size: 16px;
  color: var(--text-secondary);
  margin: 0;
  line-height: 1.5;
}

.header-actions {
  display: flex;
  gap: var(--spacing-md);
}

.creation-steps {
  background: white;
  padding: var(--spacing-xl);
  border-radius: var(--radius-lg);
  margin-bottom: var(--spacing-lg);
  box-shadow: var(--shadow-sm);
}

.step-content {
  background: white;
  border-radius: var(--radius-lg);
  margin-bottom: var(--spacing-lg);
  box-shadow: var(--shadow-sm);
  overflow: hidden;
}

.step-panel {
  padding: var(--spacing-xl);
}

.panel-header {
  margin-bottom: var(--spacing-xl);
  padding-bottom: var(--spacing-lg);
  border-bottom: 1px solid var(--border-color);
}

.panel-header h2 {
  font-size: 20px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 var(--spacing-sm) 0;
}

.panel-header p {
  font-size: 14px;
  color: var(--text-secondary);
  margin: 0;
}

.basic-form,
.config-form {
  max-width: 800px;
}

.form-tip {
  font-size: 12px;
  color: var(--text-secondary);
  margin-left: var(--spacing-sm);
  line-height: 1.4;
}

.question-filters {
  margin-bottom: var(--spacing-lg);
  padding: var(--spacing-lg);
  background: var(--bg-secondary);
  border-radius: var(--radius-md);
}

.selection-stats {
  margin-bottom: var(--spacing-lg);
  padding: var(--spacing-lg);
  background: var(--bg-secondary);
  border-radius: var(--radius-md);
}

.stats-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(120px, 1fr));
  gap: var(--spacing-md);
  margin-bottom: var(--spacing-lg);
}

.stat-card {
  text-align: center;
  padding: var(--spacing-md);
  background: white;
  border-radius: var(--radius-md);
  border: 1px solid var(--border-color);
}

.stat-number {
  font-size: 24px;
  font-weight: 700;
  color: var(--dopamine-blue);
  margin-bottom: var(--spacing-sm);
}

.stat-label {
  font-size: 12px;
  color: var(--text-secondary);
}

.selection-actions {
  display: flex;
  justify-content: flex-end;
  gap: var(--spacing-md);
}

.question-list {
  margin-bottom: var(--spacing-lg);
}

.question-item {
  display: flex;
  gap: var(--spacing-md);
  padding: var(--spacing-lg);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-md);
  margin-bottom: var(--spacing-md);
  background: white;
  transition: all 0.3s ease;
}

.question-item:hover {
  border-color: var(--dopamine-blue);
  box-shadow: var(--shadow-sm);
}

.question-item.selected {
  border-color: var(--dopamine-blue);
  background: var(--dopamine-blue-light);
}

.question-checkbox {
  flex-shrink: 0;
  padding-top: 2px;
}

.question-content {
  flex: 1;
}

.question-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--spacing-md);
}

.question-meta {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
}

.question-score {
  font-size: 12px;
  font-weight: 600;
  color: var(--dopamine-orange);
  background: var(--dopamine-orange-light);
  padding: 2px 6px;
  border-radius: var(--radius-sm);
}

.question-text {
  font-size: 14px;
  line-height: 1.5;
  color: var(--text-primary);
  margin-bottom: var(--spacing-md);
}

.question-info {
  display: flex;
  gap: var(--spacing-lg);
}

.info-item {
  font-size: 12px;
  color: var(--text-secondary);
}

.pagination-wrapper {
  display: flex;
  justify-content: center;
  margin-top: var(--spacing-xl);
}

.config-section {
  margin-bottom: var(--spacing-xl);
  padding-bottom: var(--spacing-lg);
  border-bottom: 1px solid var(--border-color);
}

.config-section:last-child {
  border-bottom: none;
}

.config-section h3 {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 var(--spacing-lg) 0;
}

.exam-preview {
  max-width: 800px;
}

.preview-section {
  margin-bottom: var(--spacing-xl);
  padding: var(--spacing-lg);
  background: var(--bg-secondary);
  border-radius: var(--radius-md);
}

.preview-section h3 {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 var(--spacing-lg) 0;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: var(--spacing-md);
}

.info-item,
.stat-item,
.config-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--spacing-sm) 0;
}

.info-label,
.stat-label,
.config-label {
  font-weight: 500;
  color: var(--text-secondary);
}

.info-value,
.stat-value,
.config-value {
  font-weight: 600;
  color: var(--text-primary);
}

.question-stats,
.config-summary {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: var(--spacing-sm);
}

.step-navigation {
  display: flex;
  justify-content: center;
  gap: var(--spacing-md);
  padding: var(--spacing-lg);
  background: white;
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-sm);
}

.selected-questions-list {
  max-height: 60vh;
  overflow-y: auto;
}

.selected-question-item {
  display: flex;
  gap: var(--spacing-md);
  padding: var(--spacing-md);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-md);
  margin-bottom: var(--spacing-sm);
  background: white;
}

.question-index {
  flex-shrink: 0;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--dopamine-blue);
  color: white;
  border-radius: 50%;
  font-weight: 600;
  font-size: 14px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .exam-creation {
    padding: var(--spacing-md);
  }
  
  .page-header {
    flex-direction: column;
    gap: var(--spacing-lg);
  }
  
  .header-actions {
    width: 100%;
    justify-content: flex-end;
  }
  
  .stats-cards {
    grid-template-columns: repeat(auto-fit, minmax(100px, 1fr));
  }
  
  .question-header {
    flex-direction: column;
    align-items: flex-start;
    gap: var(--spacing-sm);
  }
  
  .question-info {
    flex-direction: column;
    gap: var(--spacing-sm);
  }
  
  .info-grid,
  .question-stats,
  .config-summary {
    grid-template-columns: 1fr;
  }
  
  .step-navigation {
    flex-direction: column;
  }
}
</style>