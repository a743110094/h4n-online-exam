<template>
  <div class="question-form-modern">
    <div class="form-header">
      <div class="header-content">
        <el-icon class="header-icon"><EditPen /></el-icon>
        <div class="header-text">
          <h2 class="form-title">{{ props.question ? '编辑题目' : '添加题目' }}</h2>
          <p class="form-subtitle">请填写完整的题目信息</p>
        </div>
      </div>
    </div>

    <div class="form-container-modern">
      <el-form
        ref="formRef"
        :model="formData"
        :rules="formRules"
        label-width="120px"
        @submit.prevent
        class="modern-form"
      >
        <div class="form-layout-modern">
          <!-- 左侧列 -->
          <div class="form-left-modern">
            <!-- 基本信息 -->
            <div class="form-section-modern">
              <div class="section-header">
                <el-icon class="section-icon"><Setting /></el-icon>
                <h3 class="section-title-modern">基本信息</h3>
              </div>
              <div class="section-content">
            
                <el-form-item label="题目类型" prop="type" class="modern-form-item">
                  <el-select
                    v-model="formData.type"
                    placeholder="请选择题目类型"
                    @change="handleTypeChange"
                    size="small"
                    class="modern-select"
                  >
                    <el-option label="单选题" value="single">
                      <div class="option-with-icon">
                        <el-icon><CircleCheck /></el-icon>
                        <span>单选题</span>
                      </div>
                    </el-option>
                    <el-option label="多选题" value="multiple">
                      <div class="option-with-icon">
                        <el-icon><Select /></el-icon>
                        <span>多选题</span>
                      </div>
                    </el-option>
                    <el-option label="判断题" value="judge">
                      <div class="option-with-icon">
                        <el-icon><Check /></el-icon>
                        <span>判断题</span>
                      </div>
                    </el-option>
                    <el-option label="填空题" value="fill">
                      <div class="option-with-icon">
                        <el-icon><Edit /></el-icon>
                        <span>填空题</span>
                      </div>
                    </el-option>
                    <el-option label="简答题" value="essay">
                      <div class="option-with-icon">
                        <el-icon><Document /></el-icon>
                        <span>简答题</span>
                      </div>
                    </el-option>
                  </el-select>
                </el-form-item>
                
                <el-form-item label="科目" prop="subject" class="modern-form-item">
                  <el-select v-model="formData.subject" placeholder="请选择科目" size="small" class="modern-select">
                    <el-option
                      v-for="subject in subjects"
                      :key="subject.id"
                      :label="subject.name"
                      :value="subject.name"
                    >
                      <div class="option-with-icon">
                        <el-icon><Folder /></el-icon>
                        <span>{{ subject.name }}</span>
                      </div>
                    </el-option>
                  </el-select>
                </el-form-item>
                
                <el-form-item label="知识点" prop="knowledgePoint" class="modern-form-item">
                  <el-select v-model="formData.knowledgePoint" placeholder="请选择知识点" size="small" class="modern-select">
                    <el-option
                      v-for="point in knowledgePoints"
                      :key="point.id"
                      :label="point.name"
                      :value="point.name"
                    >
                      <div class="option-with-icon">
                        <el-icon><Collection /></el-icon>
                        <span>{{ point.name }}</span>
                      </div>
                    </el-option>
                  </el-select>
                </el-form-item>
                
                <el-form-item label="难度等级" prop="difficulty" class="modern-form-item">
                  <el-select v-model="formData.difficulty" placeholder="请选择难度" size="small" class="modern-select">
                    <el-option label="简单" value="easy">
                      <div class="option-with-icon difficulty-easy">
                        <el-icon><Star /></el-icon>
                        <span>简单</span>
                      </div>
                    </el-option>
                    <el-option label="中等" value="medium">
                      <div class="option-with-icon difficulty-medium">
                        <el-icon><Star /></el-icon>
                        <el-icon><Star /></el-icon>
                        <span>中等</span>
                      </div>
                    </el-option>
                    <el-option label="困难" value="hard">
                      <div class="option-with-icon difficulty-hard">
                        <el-icon><Star /></el-icon>
                        <el-icon><Star /></el-icon>
                        <el-icon><Star /></el-icon>
                        <span>困难</span>
                      </div>
                    </el-option>
                  </el-select>
                </el-form-item>
                
                <el-form-item label="分值" prop="score" class="modern-form-item">
                  <el-input-number
                    v-model="formData.score"
                    :min="1"
                    :max="100"
                    placeholder="请输入分值"
                    size="small"
                    class="modern-input-number"
                  />
                </el-form-item>
              </div>
            </div>
          
            <!-- 题目内容 -->
            <div class="form-section-modern">
              <div class="section-header">
                <el-icon class="section-icon"><Document /></el-icon>
                <h3 class="section-title-modern">题目内容</h3>
              </div>
              <div class="section-content">
                <el-form-item label="题目描述" prop="content" class="modern-form-item">
                  <el-input
                    v-model="formData.content"
                    type="textarea"
                    :rows="5"
                    placeholder="请输入题目描述，支持富文本格式..."
                    maxlength="1000"
                    show-word-limit
                    class="modern-textarea"
                  />
                </el-form-item>
              </div>
            </div>
            
            <!-- 解析说明 -->
            <div class="form-section-modern">
              <div class="section-header">
                <el-icon class="section-icon"><ChatLineRound /></el-icon>
                <h3 class="section-title-modern">解析说明</h3>
              </div>
              <div class="section-content">
                <el-form-item label="题目解析" class="modern-form-item">
                  <el-input
                    v-model="formData.explanation"
                    type="textarea"
                    :rows="4"
                    placeholder="请输入题目解析，帮助学生理解答案..."
                    maxlength="1000"
                    show-word-limit
                    class="modern-textarea"
                  />
                </el-form-item>
              </div>
            </div>
          </div>
        
          <!-- 右侧列 -->
          <div class="form-right-modern">
            <!-- 答案设置 -->
            <div class="form-section-modern">
              <div class="section-header">
                <el-icon class="section-icon"><CircleCheck /></el-icon>
                <h3 class="section-title-modern">答案设置</h3>
              </div>
              <div class="section-content">
             
                <!-- 选择题选项 -->
                <div v-if="isChoiceQuestion" class="options-section-modern">
                  <div class="options-header">
                    <el-button
                      v-if="formData.options.length < 10"
                      type="primary"
                      size="small"
                      @click="addOption"
                      class="header-add-option-btn"
                    >
                      <el-icon><Plus /></el-icon>
                      添加选项
                    </el-button>
                  </div>
                  
                  <div class="options-container-modern">
                    <div
                      v-for="(option, index) in formData.options"
                      :key="index"
                      class="option-item-modern"
                    >
                      <div class="option-header-modern">
                        <div class="option-label-modern">
                          <span class="option-letter">{{ String.fromCharCode(65 + index) }}</span>
                          <span class="option-text">选项 {{ String.fromCharCode(65 + index) }}</span>
                        </div>
                        <el-button
                          v-if="formData.options.length > 2"
                          type="danger"
                          size="small"
                          text
                          @click="removeOption(index)"
                          class="remove-option-btn"
                        >
                          <el-icon><Delete /></el-icon>
                        </el-button>
                      </div>
                      
                      <div class="option-content-modern">
                        <div class="option-input-modern">
                          <el-input
                            v-model="option.text"
                            placeholder="请输入选项内容"
                            maxlength="200"
                            size="small"
                            class="modern-input"
                          />
                        </div>
                        
                        <div class="option-correct-modern">
                          <el-checkbox
                            v-if="formData.type === 'multiple'"
                            v-model="option.isCorrect"
                            @change="validateCorrectAnswers"
                            size="small"
                            class="correct-checkbox"
                          >
                            <span class="correct-label">正确答案</span>
                          </el-checkbox>
                          <el-radio
                            v-else
                            v-model="correctSingleAnswer"
                            :label="index"
                            @change="updateSingleCorrectAnswer"
                            size="small"
                            class="correct-radio"
                          >
                            <span class="correct-label">正确答案</span>
                          </el-radio>
                        </div>
                      </div>
                    </div>
                    
                    <div class="add-option-container">
                      <el-button
                        v-if="formData.options.length < 6"
                        type="primary"
                        @click="addOption"
                        class="add-option-btn"
                        size="small"
                      >
                        <el-icon><Plus /></el-icon>
                        添加选项
                      </el-button>
                    </div>
                  </div>
                </div>
             
                <!-- 判断题答案 -->
                <div v-if="formData.type === 'judge'" class="boolean-answer-section">
                  <h4 class="boolean-answer-title">
                    <el-icon><Check /></el-icon>
                    判断题答案
                  </h4>
                  <div class="boolean-options">
                    <div 
                      class="boolean-option" 
                      :class="{ selected: formData.judgeAnswer === true }"
                      @click="formData.judgeAnswer = true"
                    >
                      <el-icon><CircleCheck /></el-icon>
                      <span>正确</span>
                    </div>
                    <div 
                      class="boolean-option" 
                      :class="{ selected: formData.judgeAnswer === false }"
                      @click="formData.judgeAnswer = false"
                    >
                      <el-icon><Delete /></el-icon>
                      <span>错误</span>
                    </div>
                  </div>
                </div>
             
                <!-- 填空题答案 -->
                <div v-if="formData.type === 'fill'" class="fill-answer-section">
                  <h4 class="fill-answer-title">
                    <el-icon><Edit /></el-icon>
                    填空题答案
                  </h4>
                  <div class="fill-answers-modern">
                    <div
                      v-for="(answer, index) in formData.fillAnswers"
                      :key="index"
                      class="fill-answer-item-modern"
                    >
                      <div class="answer-number">{{ index + 1 }}</div>
                      <el-input
                        v-model="formData.fillAnswers[index]"
                        placeholder="请输入参考答案"
                        class="answer-input"
                      />
                      <el-button
                        v-if="formData.fillAnswers.length > 1"
                        type="danger"
                        size="small"
                        text
                        @click="removeFillAnswer(index)"
                        class="remove-answer-btn"
                      >
                        <el-icon><Delete /></el-icon>
                      </el-button>
                    </div>
                    
                    <el-button
                      type="primary"
                      text
                      @click="addFillAnswer"
                      class="add-answer-btn"
                    >
                      <el-icon><Plus /></el-icon>
                      添加答案
                    </el-button>
                  </div>
                </div>
             
                <!-- 简答题答案 -->
                <div v-if="formData.type === 'essay'" class="essay-answer-section">
                  <h4 class="essay-answer-title">
                    <el-icon><Document /></el-icon>
                    简答题答案
                  </h4>
                  <div class="essay-content">
                    <div class="essay-field">
                      <label class="field-label">参考答案</label>
                      <el-input
                        v-model="formData.essayAnswer"
                        type="textarea"
                        :rows="4"
                        placeholder="请输入参考答案"
                        maxlength="2000"
                        show-word-limit
                        class="essay-textarea"
                      />
                    </div>
                    
                    <div class="essay-field">
                      <label class="field-label">评分标准</label>
                      <el-input
                        v-model="formData.scoringCriteria"
                        type="textarea"
                        :rows="3"
                        placeholder="请输入评分标准（可选）"
                        maxlength="1000"
                        show-word-limit
                        class="essay-textarea"
                      />
                    </div>
                  </div>
                </div>
              </div>
            </div>
            
            <!-- 其他设置 -->
            <div class="form-section-modern">
              <div class="section-header">
                <el-icon class="section-icon"><Setting /></el-icon>
                <h3 class="section-title-modern">其他设置</h3>
              </div>
              <div class="section-content">
              
              <el-form-item label="状态" prop="status">
                <el-select v-model="formData.status" placeholder="请选择状态">
                  <el-option label="草稿" value="draft" />
                  <el-option label="已发布" value="published" />
                </el-select>
              </el-form-item>
              
              <el-form-item label="标签">
                <el-select
                  v-model="formData.tags"
                  multiple
                  filterable
                  allow-create
                  placeholder="请选择或输入标签"
                >
                  <el-option
                    v-for="tag in commonTags"
                    :key="tag"
                    :label="tag"
                    :value="tag"
                  />
                </el-select>
              </el-form-item>
              </div>
            </div>
          </div>
        </div>
      </el-form>
    </div>
    
    <!-- 操作按钮 -->
    <div class="form-actions">
      <el-button @click="handleCancel">取消</el-button>
      <el-button type="info" @click="handlePreview">预览</el-button>
      <el-button type="primary" @click="handleSave" :loading="saving">
        保存
      </el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { 
  Plus, Delete, EditPen, Setting, Document, ChatLineRound,
  CircleCheck, Select, Check, Edit, Folder, Collection, Star
} from '@element-plus/icons-vue'

interface Props {
  question?: any
}

interface Emits {
  save: [questionData: any]
  cancel: []
  preview: [questionData: any]
}

const props = withDefaults(defineProps<Props>(), {
  question: null
})

const emit = defineEmits<Emits>()

// 表单引用
const formRef = ref()
const saving = ref(false)

// 单选题正确答案索引
const correctSingleAnswer = ref(0)

// 表单数据
const formData = reactive({
  type: 'single',
  subject: '',
  knowledgePoint: '',
  difficulty: 'medium',
  score: 2,
  content: '',
  options: [
    { text: '', isCorrect: true },
    { text: '', isCorrect: false },
    { text: '', isCorrect: false },
    { text: '', isCorrect: false }
  ],
  judgeAnswer: true,
  fillAnswers: [''],
  essayAnswer: '',
  scoringCriteria: '',
  explanation: '',
  status: 'draft',
  tags: []
})

// 科目列表
const subjects = ref([
  { id: 1, name: '数据结构' },
  { id: 2, name: '算法设计' },
  { id: 3, name: '操作系统' },
  { id: 4, name: '计算机网络' },
  { id: 5, name: '数据库原理' }
])

// 知识点列表
const knowledgePoints = ref([
  { id: 1, name: '线性表' },
  { id: 2, name: '栈和队列' },
  { id: 3, name: '树和二叉树' },
  { id: 4, name: '图' },
  { id: 5, name: '排序算法' },
  { id: 6, name: '查找算法' }
])

// 常用标签
const commonTags = ref([
  '基础', '重点', '难点', '易错', '综合', '实践'
])

// 是否为选择题
const isChoiceQuestion = computed(() => {
  return ['single', 'multiple'].includes(formData.type)
})

// 表单验证规则
const formRules = {
  type: [
    { required: true, message: '请选择题目类型', trigger: 'change' }
  ],
  subject: [
    { required: true, message: '请选择科目', trigger: 'change' }
  ],
  knowledgePoint: [
    { required: true, message: '请选择知识点', trigger: 'change' }
  ],
  difficulty: [
    { required: true, message: '请选择难度等级', trigger: 'change' }
  ],
  score: [
    { required: true, message: '请输入分值', trigger: 'blur' }
  ],
  content: [
    { required: true, message: '请输入题目描述', trigger: 'blur' },
    { min: 10, message: '题目描述至少10个字符', trigger: 'blur' }
  ],
  status: [
    { required: true, message: '请选择状态', trigger: 'change' }
  ]
}

// 监听题目类型变化
watch(() => formData.type, (newType) => {
  if (newType === 'single' || newType === 'multiple') {
    // 确保至少有两个选项
    if (formData.options.length < 2) {
      formData.options = [
        { text: '', isCorrect: false },
        { text: '', isCorrect: false }
      ]
    }
  }
})

// 处理题目类型变化
const handleTypeChange = (type: string) => {
  if (type === 'single') {
    // 单选题：重置所有选项为非正确答案，设置第一个为正确答案
    formData.options.forEach((option, index) => {
      option.isCorrect = index === 0
    })
    correctSingleAnswer.value = 0
  } else if (type === 'multiple') {
    // 多选题：重置所有选项为非正确答案
    formData.options.forEach(option => {
      option.isCorrect = false
    })
  }
}

// 添加选项
const addOption = () => {
  if (formData.options.length < 6) {
    formData.options.push({ text: '', isCorrect: false })
  }
}

// 删除选项
const removeOption = (index: number) => {
  if (formData.options.length > 2) {
    formData.options.splice(index, 1)
    
    // 如果删除的是单选题的正确答案，重新设置第一个为正确答案
    if (formData.type === 'single') {
      const hasCorrect = formData.options.some(option => option.isCorrect)
      if (!hasCorrect && formData.options.length > 0) {
        formData.options[0].isCorrect = true
        correctSingleAnswer.value = 0
      }
    }
  }
}

// 更新单选题正确答案
const updateSingleCorrectAnswer = () => {
  formData.options.forEach((option, index) => {
    option.isCorrect = index === correctSingleAnswer.value
  })
}

// 验证多选题正确答案
const validateCorrectAnswers = () => {
  if (formData.type === 'multiple') {
    const correctCount = formData.options.filter(option => option.isCorrect).length
    if (correctCount === 0) {
      ElMessage.warning('多选题至少需要一个正确答案')
    }
  }
}

// 添加填空题答案
const addFillAnswer = () => {
  formData.fillAnswers.push('')
}

// 删除填空题答案
const removeFillAnswer = (index: number) => {
  if (formData.fillAnswers.length > 1) {
    formData.fillAnswers.splice(index, 1)
  }
}

// 处理取消
const handleCancel = () => {
  emit('cancel')
}

// 处理预览
const handlePreview = () => {
  emit('preview', { ...formData })
}

// 处理保存
const handleSave = async () => {
  try {
    await formRef.value.validate()
    
    // 验证选择题选项
    if (isChoiceQuestion.value) {
      // 过滤掉空的选项对象
      const validOptions = formData.options.filter(option => option && typeof option === 'object')
      
      const hasEmptyOption = validOptions.some(option => !option.text || !option.text.trim())
      if (hasEmptyOption) {
        ElMessage.error('请填写所有选项内容')
        return
      }
      
      const hasCorrectAnswer = validOptions.some(option => option.isCorrect)
      if (!hasCorrectAnswer) {
        ElMessage.error('请设置正确答案')
        return
      }
    }
    
    // 验证填空题答案
    if (formData.type === 'fill') {
      const hasEmptyAnswer = formData.fillAnswers.some(answer => !answer || !answer.trim())
      if (hasEmptyAnswer) {
        ElMessage.error('请填写所有参考答案')
        return
      }
    }
    
    // 验证简答题答案
    if (formData.type === 'essay' && (!formData.essayAnswer || !formData.essayAnswer.trim())) {
      ElMessage.error('请填写参考答案')
      return
    }
    
    saving.value = true
    
    // 处理数据格式，确保后端能正确解析
    const submitData = { ...formData }
    
    // 将options对象数组转换为字符串数组
    if (isChoiceQuestion.value && Array.isArray(submitData.options)) {
      submitData.options = submitData.options.map(option => option.text)
    }
    
    // 将answer字段转换为字符串格式
    if (formData.type === 'single' || formData.type === 'multiple') {
      // 选择题：将正确选项转换为JSON字符串
      const correctOptions = formData.options
        .map((option, index) => option.isCorrect ? index : -1)
        .filter(index => index !== -1)
      submitData.answer = JSON.stringify(correctOptions)
    } else if (formData.type === 'judge') {
      // 判断题：将布尔值转换为字符串
      submitData.answer = String(formData.judgeAnswer)
    } else if (formData.type === 'fill') {
      // 填空题：将答案数组转换为JSON字符串
      submitData.answer = JSON.stringify(formData.fillAnswers)
    } else if (formData.type === 'essay') {
      // 简答题：直接使用字符串
      submitData.answer = formData.essayAnswer
    }
    
    // 模拟保存过程
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    emit('save', submitData)
  } catch (error) {
    console.error('表单验证失败:', error)
  } finally {
    saving.value = false
  }
}

// 初始化表单数据
const initFormData = () => {
  if (props.question) {
    Object.assign(formData, props.question)
    
    // 解析JSON字符串字段
    if (typeof formData.options === 'string') {
      try {
        formData.options = JSON.parse(formData.options)
      } catch (e) {
        console.error('解析options失败:', e)
        formData.options = []
      }
    }
    
    // 确保options是数组且包含正确的对象结构
    if (!Array.isArray(formData.options)) {
      formData.options = []
    }
    
    // 如果options是字符串数组，转换为对象数组
    if (formData.options.length > 0 && typeof formData.options[0] === 'string') {
      formData.options = formData.options.map((text: string, index: number) => ({
        text: text,
        isCorrect: false
      }))
    }
    
    // 确保每个选项都有正确的结构
    formData.options = formData.options.map((option: any) => {
      if (typeof option === 'string') {
        return { text: option, isCorrect: false }
      }
      return {
        text: option.text || '',
        isCorrect: option.isCorrect || false
      }
    })
    
    if (typeof formData.answer === 'string' && formData.type !== 'essay') {
      try {
        formData.answer = JSON.parse(formData.answer)
      } catch (e) {
        console.error('解析answer失败:', e)
        formData.answer = []
      }
    }
    
    // 设置单选题正确答案索引
    if (formData.type === 'single' && formData.options.length > 0) {
      const correctIndex = formData.options.findIndex((option: any) => option.isCorrect)
      correctSingleAnswer.value = correctIndex >= 0 ? correctIndex : 0
    }
    
    // 处理填空题答案
    if (formData.type === 'fill') {
      if (Array.isArray(formData.answer)) {
        formData.fillAnswers = formData.answer
      } else if (typeof formData.answer === 'string') {
        formData.fillAnswers = [formData.answer]
      }
    }
    
    // 处理简答题答案
    if (formData.type === 'essay') {
      if (typeof formData.answer === 'string') {
        formData.essayAnswer = formData.answer
      }
    }
  }
}

onMounted(() => {
  initFormData()
})
</script>

<style scoped>
:root {
  --modern-spacing-xs: 6px;
  --modern-spacing-sm: 6px;
  --modern-spacing-md: 6px;
  --modern-spacing-lg: 6px;
  --modern-spacing-xl: 6px;
  --modern-radius-sm: 6px;
  --modern-radius-md: 8px;
  --modern-radius-lg: 12px;
  --modern-shadow-sm: 0 2px 4px rgba(0, 0, 0, 0.05);
  --modern-shadow-md: 0 4px 12px rgba(0, 0, 0, 0.1);
  --modern-shadow-lg: 0 8px 24px rgba(0, 0, 0, 0.15);
  --primary-gradient: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  --success-gradient: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
  --warning-gradient: linear-gradient(135deg, #fa709a 0%, #fee140 100%);
}

/* 现代化表单容器 */
.question-form-modern {
  width: 100%;
  height: 100%;
  margin: 0;
  background: #f8fafc;
  display: flex;
  flex-direction: column;
}

/* 表单头部 */
.form-header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 6px;
  border-radius: 0 0 var(--modern-radius-lg) var(--modern-radius-lg);
  margin-bottom: 6px;
  box-shadow: var(--modern-shadow-lg);
}

.header-content {
  display: flex;
  align-items: center;
  gap: var(--modern-spacing-md);
  margin-bottom: var(--modern-spacing-lg);
}

.header-icon {
  font-size: 32px;
  opacity: 0.9;
}

.form-title {
  font-size: 28px;
  font-weight: 700;
  margin: 0;
  letter-spacing: -0.5px;
}

.form-subtitle {
  font-size: 16px;
  margin: 4px 0 0 0;
  opacity: 0.8;
}

/* 进度指示器 */
.progress-indicator {
  display: flex;
  justify-content: center;
}

.progress-steps {
  display: flex;
  align-items: center;
  gap: var(--modern-spacing-sm);
}

.step {
  display: flex;
  align-items: center;
  gap: var(--modern-spacing-xs);
  padding: var(--modern-spacing-sm) var(--modern-spacing-md);
  background: rgba(255, 255, 255, 0.2);
  border-radius: var(--modern-radius-md);
  backdrop-filter: blur(10px);
}

.step.active {
  background: rgba(255, 255, 255, 0.3);
}

.step-number {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.3);
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  font-size: 12px;
}

.step-label {
  font-size: 14px;
  font-weight: 500;
}

.step-divider {
  width: 20px;
  height: 2px;
  background: rgba(255, 255, 255, 0.3);
  border-radius: 1px;
}

/* 现代化表单容器 */
.form-container-modern {
  flex: 1;
  padding: 6px;
  height: calc(100% - 80px);
}

.modern-form {
  background: white;
  border-radius: var(--modern-radius-lg);
  box-shadow: var(--modern-shadow-md);
  overflow: hidden;
}

.form-layout-modern {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 6px;
  padding: 6px;
  height: 100%;
}

.form-left-modern,
.form-right-modern {
  display: flex;
  flex-direction: column;
  gap: 6px;
  min-width: 0; /* 防止内容溢出 */
  height: 100%;
}

/* 现代化表单区域 */
.form-section-modern {
  background: #ffffff;
  border-radius: var(--modern-radius-md);
  border: 1px solid #e2e8f0;
  overflow: visible; /* 改为visible防止内容被裁剪 */
  transition: all 0.3s ease;
  margin-bottom: 6px;
}

.form-section-modern:hover {
  border-color: #cbd5e0;
  box-shadow: var(--modern-shadow-sm);
}

.section-header {
  background: linear-gradient(135deg, #f7fafc 0%, #edf2f7 100%);
  padding: 6px;
  border-bottom: 1px solid #e2e8f0;
  display: flex;
  align-items: center;
  gap: 6px;
  position: sticky;
  top: 0;
  z-index: 10;
}

.section-icon {
  font-size: 18px;
  color: #667eea;
}

.section-title-modern {
  font-size: 16px;
  font-weight: 600;
  color: #2d3748;
  margin: 0;
}

.section-content {
  padding: 6px;
  max-height: none; /* 移除高度限制 */
}

/* 现代化表单元素 */
.modern-form-item {
  margin-bottom: 6px;
}

.modern-select,
.modern-input,
.modern-textarea,
.modern-input-number {
  border-radius: var(--modern-radius-sm);
  transition: all 0.3s ease;
}

/* 表单项标签样式 */
:deep(.el-form-item__label) {
  font-weight: 600;
  color: #374151;
  margin-bottom: 6px;
}

/* 表单项内容样式 */
:deep(.el-form-item__content) {
  margin-top: 6px;
}

/* 输入框样式优化 */
:deep(.el-input__wrapper) {
  border-radius: var(--modern-radius-sm);
  transition: all 0.3s ease;
  padding: 8px 12px;
}

:deep(.el-textarea__inner) {
  border-radius: var(--modern-radius-sm);
  transition: all 0.3s ease;
  padding: 12px;
}

:deep(.el-select .el-input__wrapper) {
  border-radius: var(--modern-radius-sm);
}

.option-with-icon {
  display: flex;
  align-items: center;
  gap: 6px;
}

.difficulty-easy {
  color: #48bb78;
}

.difficulty-medium {
  color: #ed8936;
}

.difficulty-hard {
  color: #f56565;
}

/* 现代化选项样式 */
.options-section-modern {
  background: #f7fafc;
  border-radius: var(--modern-radius-md);
  padding: 6px;
  border: 1px solid #e2e8f0;
  width: 100%;
  box-sizing: border-box;
}

.options-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 6px;
  padding-bottom: 6px;
  border-bottom: 2px solid #e2e8f0;
}

.options-title {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 16px;
  font-weight: 600;
  color: #2d3748;
  margin: 0;
}

.options-container-modern {
  display: flex;
  flex-direction: column;
  gap: 6px;
  width: 100%;
}

.option-item-modern {
  background: white;
  border-radius: var(--modern-radius-md);
  border: 2px solid #e2e8f0;
  padding: 6px;
  transition: all 0.3s ease;
  position: relative;
  width: 100%;
  box-sizing: border-box;
}

.option-item-modern:hover {
  border-color: #cbd5e0;
  box-shadow: var(--modern-shadow-sm);
  transform: translateY(-1px);
}

.option-header-modern {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 6px;
}

.option-label-modern {
  display: flex;
  align-items: center;
  gap: 6px;
}

.option-letter {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  font-size: 12px;
  flex-shrink: 0;
}

.option-text {
  font-weight: 600;
  color: #4a5568;
  font-size: 14px;
}

.remove-option-btn {
  opacity: 0.6;
  transition: opacity 0.3s ease;
  flex-shrink: 0;
}

.remove-option-btn:hover {
  opacity: 1;
}

.option-content-modern {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.option-input-modern {
  width: 100%;
}

.option-correct-modern {
  display: flex;
  align-items: center;
  justify-content: flex-start;
  margin-top: 6px;
}

.correct-label {
  font-weight: 500;
  color: #4a5568;
}

.correct-checkbox,
.correct-radio {
  font-size: 16px;
}

.add-option-container {
  display: flex;
  justify-content: center;
  padding: 6px 0;
  margin-top: 6px;
  border-top: 1px dashed #cbd5e0;
  width: 100%;
}

.add-option-btn {
  background: var(--primary-gradient);
  border: none;
  border-radius: var(--modern-radius-md);
  padding: 6px;
  font-weight: 600;
  transition: all 0.3s ease;
  color: white;
  display: flex;
  align-items: center;
  gap: 6px;
  min-height: 36px;
}

.add-option-btn:hover {
  transform: translateY(-2px);
  box-shadow: var(--modern-shadow-md);
  opacity: 0.9;
}

/* 判断题样式 */
.boolean-answer-section {
  background: #f7fafc;
  border-radius: var(--modern-radius-md);
  padding: var(--modern-spacing-lg);
  border: 1px solid #e2e8f0;
}

.boolean-answer-title {
  display: flex;
  align-items: center;
  gap: var(--modern-spacing-sm);
  font-size: 16px;
  font-weight: 600;
  color: #2d3748;
  margin: 0 0 var(--modern-spacing-lg) 0;
  padding-bottom: var(--modern-spacing-md);
  border-bottom: 2px solid #e2e8f0;
}

.boolean-options {
  display: flex;
  gap: var(--modern-spacing-lg);
}

.boolean-option {
  flex: 1;
  padding: var(--modern-spacing-lg);
  background: white;
  border: 2px solid #e2e8f0;
  border-radius: var(--modern-radius-md);
  text-align: center;
  transition: all 0.3s ease;
  cursor: pointer;
}

.boolean-option:hover {
  border-color: #cbd5e0;
  box-shadow: var(--modern-shadow-sm);
  transform: translateY(-1px);
}

.boolean-option.selected {
  border-color: #4299e1;
  background: #ebf8ff;
}

/* 填空题样式 */
.blank-answers-section {
  background: #f7fafc;
  border-radius: var(--modern-radius-md);
  padding: var(--modern-spacing-lg);
  border: 1px solid #e2e8f0;
}

.blank-answers-title {
  display: flex;
  align-items: center;
  gap: var(--modern-spacing-sm);
  font-size: 16px;
  font-weight: 600;
  color: #2d3748;
  margin: 0 0 var(--modern-spacing-lg) 0;
  padding-bottom: var(--modern-spacing-md);
  border-bottom: 2px solid #e2e8f0;
}

.blank-answers-container {
  display: flex;
  flex-direction: column;
  gap: var(--modern-spacing-md);
}

.blank-answer-item-modern {
  background: white;
  border-radius: var(--modern-radius-md);
  border: 2px solid #e2e8f0;
  padding: var(--modern-spacing-lg);
  transition: all 0.3s ease;
}

.blank-answer-item-modern:hover {
  border-color: #cbd5e0;
  box-shadow: var(--modern-shadow-sm);
}

.blank-answer-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--modern-spacing-md);
}

.blank-answer-label {
  display: flex;
  align-items: center;
  gap: var(--modern-spacing-sm);
  font-weight: 600;
  color: #4a5568;
}

.blank-number {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  font-size: 12px;
}

.fill-answers {
  width: 100%;
}

.fill-answer-item {
  display: flex;
  gap: var(--compact-spacing-sm);
  align-items: center;
  margin-bottom: var(--compact-spacing-sm);
}

/* 判断题样式 */
.boolean-answer-section {
  background: #f7fafc;
  border-radius: var(--modern-radius-md);
  padding: var(--modern-spacing-lg);
  border: 1px solid #e2e8f0;
  margin-bottom: var(--modern-spacing-lg);
}

.boolean-answer-title {
  display: flex;
  align-items: center;
  gap: var(--modern-spacing-sm);
  font-size: 16px;
  font-weight: 600;
  color: #2d3748;
  margin: 0 0 var(--modern-spacing-lg) 0;
  padding-bottom: var(--modern-spacing-md);
  border-bottom: 2px solid #e2e8f0;
}

.boolean-options {
  display: flex;
  gap: var(--modern-spacing-lg);
}

.boolean-option {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--modern-spacing-sm);
  padding: var(--modern-spacing-lg);
  background: white;
  border: 2px solid #e2e8f0;
  border-radius: var(--modern-radius-md);
  cursor: pointer;
  transition: all 0.3s ease;
  font-weight: 600;
  color: #4a5568;
}

.boolean-option:hover {
  border-color: #cbd5e0;
  box-shadow: var(--modern-shadow-sm);
  transform: translateY(-2px);
}

.boolean-option.selected {
  background: linear-gradient(135deg, #48bb78 0%, #38a169 100%);
  border-color: #38a169;
  color: white;
}

/* 填空题样式 */
.fill-answer-section {
  background: #f7fafc;
  border-radius: var(--modern-radius-md);
  padding: var(--modern-spacing-lg);
  border: 1px solid #e2e8f0;
  margin-bottom: var(--modern-spacing-lg);
}

.fill-answer-title {
  display: flex;
  align-items: center;
  gap: var(--modern-spacing-sm);
  font-size: 16px;
  font-weight: 600;
  color: #2d3748;
  margin: 0 0 var(--modern-spacing-lg) 0;
  padding-bottom: var(--modern-spacing-md);
  border-bottom: 2px solid #e2e8f0;
}

.fill-answers-modern {
  display: flex;
  flex-direction: column;
  gap: var(--modern-spacing-md);
}

.fill-answer-item-modern {
  display: flex;
  align-items: center;
  gap: var(--modern-spacing-md);
  background: white;
  border: 2px solid #e2e8f0;
  border-radius: var(--modern-radius-md);
  padding: var(--modern-spacing-md);
  transition: all 0.3s ease;
}

.fill-answer-item-modern:hover {
  border-color: #cbd5e0;
  box-shadow: var(--modern-shadow-sm);
}

.answer-number {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  font-size: 14px;
  flex-shrink: 0;
}

.answer-input {
  flex: 1;
}

.remove-answer-btn {
  flex-shrink: 0;
}

.add-answer-btn {
  align-self: flex-start;
  margin-top: var(--modern-spacing-sm);
}

/* 简答题样式 */
.essay-answer-section {
  background: #f7fafc;
  border-radius: var(--modern-radius-md);
  padding: var(--modern-spacing-lg);
  border: 1px solid #e2e8f0;
  margin-bottom: var(--modern-spacing-lg);
}

.essay-answer-title {
  display: flex;
  align-items: center;
  gap: var(--modern-spacing-sm);
  font-size: 16px;
  font-weight: 600;
  color: #2d3748;
  margin: 0 0 var(--modern-spacing-lg) 0;
  padding-bottom: var(--modern-spacing-md);
  border-bottom: 2px solid #e2e8f0;
}

.essay-content {
  display: flex;
  flex-direction: column;
  gap: var(--modern-spacing-lg);
}

.essay-field {
  display: flex;
  flex-direction: column;
  gap: var(--modern-spacing-sm);
}

.field-label {
  font-size: 14px;
  font-weight: 600;
  color: #4a5568;
}

.essay-textarea {
  background: white;
  border-radius: var(--modern-radius-md);
}

/* 其他设置样式 */
.other-settings-section {
  background: #f7fafc;
  border-radius: var(--modern-radius-md);
  padding: var(--modern-spacing-lg);
  border: 1px solid #e2e8f0;
}

.other-settings-title {
  display: flex;
  align-items: center;
  gap: var(--modern-spacing-sm);
  font-size: 16px;
  font-weight: 600;
  color: #2d3748;
  margin: 0 0 var(--modern-spacing-lg) 0;
  padding-bottom: var(--modern-spacing-md);
  border-bottom: 2px solid #e2e8f0;
}

.form-actions {
  background: white;
  padding: var(--modern-spacing-xl);
  border-top: 1px solid #e2e8f0;
  display: flex;
  justify-content: center;
  gap: var(--modern-spacing-lg);
  margin-top: var(--modern-spacing-xl);
  border-radius: var(--modern-radius-lg);
  box-shadow: var(--modern-shadow-md);
}

.form-actions .el-button {
  padding: 12px 24px;
  border-radius: var(--modern-radius-sm);
  font-weight: 600;
  min-width: 100px;
  transition: all 0.3s ease;
}

.form-actions .el-button:hover {
  transform: translateY(-1px);
  box-shadow: var(--modern-shadow-sm);
}

/* 表单操作按钮现代化 */
.form-actions-modern {
  display: flex;
  justify-content: flex-end;
  gap: var(--modern-spacing-md);
  padding: var(--modern-spacing-xl) 0 0 0;
  margin-top: var(--modern-spacing-xl);
  border-top: 2px solid #e2e8f0;
}

.action-btn-modern {
  min-width: 120px;
  height: 48px;
  border-radius: var(--modern-radius-md);
  font-weight: 600;
  font-size: 16px;
  transition: all 0.3s ease;
  border: 2px solid transparent;
}

.action-btn-modern:hover {
  transform: translateY(-2px);
  box-shadow: var(--modern-shadow-md);
}

.cancel-btn {
  background: white;
  color: #718096;
  border-color: #e2e8f0;
}

.cancel-btn:hover {
  background: #f7fafc;
  border-color: #cbd5e0;
  color: #4a5568;
}

.preview-btn {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.save-btn {
  background: linear-gradient(135deg, #48bb78 0%, #38a169 100%);
  color: white;
}

.add-blank-btn {
  background: var(--primary-gradient);
  border: none;
  border-radius: var(--modern-radius-md);
  padding: var(--modern-spacing-md) var(--modern-spacing-lg);
  font-weight: 600;
  transition: all 0.3s ease;
  margin-top: var(--modern-spacing-md);
}

.add-blank-btn:hover {
  transform: translateY(-2px);
  box-shadow: var(--modern-shadow-md);
}

/* 表单元素紧凑化 */
:deep(.el-form-item) {
  margin-bottom: var(--compact-spacing-sm);
}

:deep(.el-form-item__label) {
  font-size: 12px;
  font-weight: 500;
  color: var(--text-secondary);
  line-height: 1.3;
  padding-bottom: var(--compact-spacing-xs);
}

:deep(.el-input__inner),
:deep(.el-textarea__inner),
:deep(.el-select .el-input__inner) {
  font-size: 12px;
  padding: var(--compact-spacing-xs) var(--compact-spacing-sm);
  border-radius: var(--compact-radius-sm);
  border: 1px solid var(--border-color);
  transition: all 0.2s ease;
}

:deep(.el-input__inner:focus),
:deep(.el-textarea__inner:focus) {
  border-color: var(--dopamine-blue);
  box-shadow: 0 0 0 2px rgba(64, 158, 255, 0.1);
}

:deep(.el-button) {
  font-size: 12px;
  padding: var(--compact-spacing-xs) var(--compact-spacing-sm);
  border-radius: var(--compact-radius-sm);
}

:deep(.el-button--small) {
  padding: var(--compact-spacing-xs) var(--compact-spacing-xs);
  font-size: 11px;
}

:deep(.el-radio__label),
:deep(.el-checkbox__label) {
  font-size: 12px;
  padding-left: var(--compact-spacing-xs);
}

/* 选项样式优化 */
.option-label {
  font-size: 12px;
  font-weight: 600;
  color: var(--dopamine-blue);
  background: rgba(64, 158, 255, 0.1);
  padding: 2px 6px;
  border-radius: var(--radius-xs);
}

/* Header添加选项按钮样式 */
.options-header {
  display: flex;
  justify-content: flex-end;
  margin-bottom: var(--modern-spacing-md);
}

.header-add-option-btn {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  color: white;
  font-weight: 600;
  transition: all 0.3s ease;
}

.header-add-option-btn:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
  background: linear-gradient(135deg, #5a67d8 0%, #6b46c1 100%);
}

/* 响应式设计 */
@media (max-width: 1200px) {
  .form-layout-modern {
    grid-template-columns: 1fr;
    gap: var(--modern-spacing-lg);
    padding: var(--modern-spacing-lg);
  }
  
  .form-left-modern,
  .form-right-modern {
    max-width: 100%;
  }
}

@media (max-width: 768px) {
  .form-container-modern {
    padding: 0 var(--modern-spacing-md) var(--modern-spacing-md);
  }
  
  .form-layout-modern {
    padding: var(--modern-spacing-md);
  }
  
  .form-header {
    padding: var(--modern-spacing-lg);
    margin-bottom: var(--modern-spacing-md);
  }
  
  .form-title {
    font-size: 24px;
  }
  
  .progress-steps {
    flex-wrap: wrap;
    gap: var(--modern-spacing-xs);
  }
  
  .step {
    padding: var(--modern-spacing-xs) var(--modern-spacing-sm);
  }
  
  .step-label {
    display: none;
  }
  
  .option-content-modern {
    gap: var(--modern-spacing-xs);
  }
  
  .form-actions {
    padding: var(--modern-spacing-md);
    gap: var(--modern-spacing-md);
    flex-direction: column;
  }
  
  .form-actions .el-button {
    width: 100%;
  }
}
</style>