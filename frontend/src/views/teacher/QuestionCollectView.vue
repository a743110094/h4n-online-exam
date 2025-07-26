<template>
  <div class="question-collect-view">
    <!-- 页面标题 -->
    <div class="page-header">
      <h1 class="page-title">网络采集题目</h1>
      <p class="page-subtitle">从网络资源智能采集题目，丰富题库内容</p>
    </div>

    <!-- 采集配置 -->
    <div class="collect-config dopamine-card">
      <h3>采集配置</h3>
      <el-form :model="collectForm" label-width="120px" class="config-form">
        <el-row :gutter="24">
          <el-col :span="12">
            <el-form-item label="采集源">
              <el-select v-model="collectForm.source" placeholder="请选择采集源">
                <el-option label="题库网站A" value="source_a" />
                <el-option label="题库网站B" value="source_b" />
                <el-option label="教育平台C" value="source_c" />
                <el-option label="自定义URL" value="custom" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="科目">
              <el-select v-model="collectForm.subject" placeholder="请选择科目">
                <el-option
                  v-for="subject in subjects"
                  :key="subject.id"
                  :label="subject.name"
                  :value="subject.name"
                />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        
        <el-row :gutter="24">
          <el-col :span="12">
            <el-form-item label="题型">
              <el-select v-model="collectForm.questionType" multiple placeholder="请选择题型">
                <el-option label="单选题" value="single" />
                <el-option label="多选题" value="multiple" />
                <el-option label="判断题" value="judge" />
                <el-option label="填空题" value="fill" />
                <el-option label="简答题" value="essay" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="难度">
              <el-select v-model="collectForm.difficulty" placeholder="请选择难度">
                <el-option label="简单" value="1" />
                <el-option label="一般" value="2" />
                <el-option label="中等" value="3" />
                <el-option label="困难" value="4" />
                <el-option label="很难" value="5" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        
        <el-form-item label="采集数量">
          <el-input-number v-model="collectForm.count" :min="1" :max="100" />
        </el-form-item>
        
        <el-form-item v-if="collectForm.source === 'custom'" label="自定义URL">
          <el-input v-model="collectForm.customUrl" placeholder="请输入要采集的网页URL" />
        </el-form-item>
        
        <el-form-item label="关键词">
          <el-input v-model="collectForm.keywords" placeholder="输入关键词，多个关键词用逗号分隔" />
        </el-form-item>
        
        <el-form-item>
          <el-button type="primary" @click="startCollect" :loading="collecting">
            <el-icon><Search /></el-icon>
            开始采集
          </el-button>
          <el-button @click="resetForm">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <!-- 采集进度 -->
    <div v-if="collecting" class="collect-progress dopamine-card">
      <h3>采集进度</h3>
      <el-progress :percentage="collectProgress" :status="collectStatus" />
      <p class="progress-text">{{ progressText }}</p>
    </div>

    <!-- 采集结果 -->
    <div v-if="collectResults.length > 0" class="collect-results dopamine-card">
      <div class="results-header">
        <h3>采集结果</h3>
        <div class="results-actions">
          <el-button @click="selectAll">全选</el-button>
          <el-button @click="selectNone">取消全选</el-button>
          <el-button type="primary" @click="importSelected" :disabled="selectedQuestions.length === 0">
            导入选中 ({{ selectedQuestions.length }})
          </el-button>
        </div>
      </div>
      
      <el-table
        :data="collectResults"
        border
        @selection-change="handleSelectionChange"
        max-height="500"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column prop="type" label="题型" min-width="80">
          <template #default="{ row }">
            <el-tag :type="getTypeTagType(row.type)">{{ row.type }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="subject" label="科目" min-width="100" />
        <el-table-column prop="difficulty" label="难度" min-width="100">
          <template #default="{ row }">
            <el-rate
              v-model="row.difficulty"
              disabled
              show-score
              text-color="#ff9900"
              score-template="{value}星"
            />
          </template>
        </el-table-column>
        <el-table-column prop="content" label="题目内容" min-width="300" show-overflow-tooltip />
        <el-table-column prop="source" label="来源" min-width="120" show-overflow-tooltip />
        <el-table-column label="质量评分" min-width="100">
          <template #default="{ row }">
            <el-progress
              :percentage="row.quality"
              :color="getQualityColor(row.quality)"
              :show-text="false"
            />
            <span class="quality-text">{{ row.quality }}%</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="120">
          <template #default="{ row }">
            <el-button size="small" @click="previewQuestion(row)">预览</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <!-- 题目预览对话框 -->
    <el-dialog v-model="showPreview" title="题目预览" width="60%">
      <div v-if="previewData" class="question-preview">
        <div class="question-info">
          <el-tag :type="getTypeTagType(previewData.type)">{{ previewData.type }}</el-tag>
          <el-tag type="info">{{ previewData.subject }}</el-tag>
          <el-rate v-model="previewData.difficulty" disabled show-text />
        </div>
        
        <div class="question-content">
          <h4>题目内容：</h4>
          <p>{{ previewData.content }}</p>
        </div>
        
        <div v-if="previewData.options" class="question-options">
          <h4>选项：</h4>
          <div v-for="(option, index) in previewData.options" :key="index" class="option-item">
            <span class="option-label">{{ String.fromCharCode(65 + index) }}.</span>
            <span>{{ option }}</span>
          </div>
        </div>
        
        <div class="question-answer">
          <h4>答案：</h4>
          <p>{{ previewData.answer }}</p>
        </div>
        
        <div v-if="previewData.explanation" class="question-explanation">
          <h4>解析：</h4>
          <p>{{ previewData.explanation }}</p>
        </div>
      </div>
      
      <template #footer>
        <el-button @click="showPreview = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import { Search } from '@element-plus/icons-vue'

// 采集表单
const collectForm = reactive({
  source: '',
  subject: '',
  questionType: [],
  difficulty: '',
  count: 10,
  customUrl: '',
  keywords: ''
})

// 科目列表
const subjects = ref([
  { id: 1, name: '数据结构' },
  { id: 2, name: '算法设计' },
  { id: 3, name: '操作系统' },
  { id: 4, name: '计算机网络' },
  { id: 5, name: '数据库原理' }
])

// 采集状态
const collecting = ref(false)
const collectProgress = ref(0)
const collectStatus = ref('')
const progressText = ref('')

// 采集结果
const collectResults = ref([
  {
    id: 1,
    type: '单选题',
    subject: '数据结构',
    difficulty: 3,
    content: '以下哪种数据结构适合实现栈？',
    options: ['数组', '链表', '树', '图'],
    answer: 'A',
    explanation: '栈是后进先出的数据结构，数组可以很好地实现这一特性。',
    source: '题库网站A',
    quality: 85
  },
  {
    id: 2,
    type: '多选题',
    subject: '算法设计',
    difficulty: 4,
    content: '以下哪些是排序算法？',
    options: ['冒泡排序', '快速排序', '深度优先搜索', '归并排序'],
    answer: 'ABD',
    explanation: '冒泡排序、快速排序、归并排序都是常见的排序算法。',
    source: '教育平台C',
    quality: 92
  }
])

// 选中的题目
const selectedQuestions = ref([])

// 预览相关
const showPreview = ref(false)
const previewData = ref(null)

// 开始采集
const startCollect = () => {
  if (!collectForm.source || !collectForm.subject) {
    ElMessage.warning('请完善采集配置')
    return
  }
  
  collecting.value = true
  collectProgress.value = 0
  collectStatus.value = ''
  progressText.value = '正在连接采集源...'
  
  // 模拟采集过程
  const timer = setInterval(() => {
    collectProgress.value += 10
    
    if (collectProgress.value <= 30) {
      progressText.value = '正在分析网页结构...'
    } else if (collectProgress.value <= 60) {
      progressText.value = '正在提取题目内容...'
    } else if (collectProgress.value <= 90) {
      progressText.value = '正在进行质量评估...'
    } else {
      progressText.value = '采集完成！'
      collectStatus.value = 'success'
      collecting.value = false
      clearInterval(timer)
      ElMessage.success(`成功采集到 ${collectResults.value.length} 道题目`)
    }
  }, 300)
}

// 重置表单
const resetForm = () => {
  Object.assign(collectForm, {
    source: '',
    subject: '',
    questionType: [],
    difficulty: '',
    count: 10,
    customUrl: '',
    keywords: ''
  })
}

// 获取题型标签类型
const getTypeTagType = (type: string) => {
  const typeMap: Record<string, string> = {
    '单选题': 'primary',
    '多选题': 'success',
    '判断题': 'warning',
    '填空题': 'info',
    '简答题': 'danger'
  }
  return typeMap[type] || 'default'
}

// 获取质量颜色
const getQualityColor = (quality: number) => {
  if (quality >= 90) return '#67c23a'
  if (quality >= 70) return '#e6a23c'
  return '#f56c6c'
}

// 选择变化处理
const handleSelectionChange = (selection: any[]) => {
  selectedQuestions.value = selection
}

// 全选
const selectAll = () => {
  // 这里需要通过ref获取table组件并调用toggleAllSelection方法
  ElMessage.info('请在表格中使用全选复选框')
}

// 取消全选
const selectNone = () => {
  selectedQuestions.value = []
  ElMessage.info('已取消全选')
}

// 导入选中题目
const importSelected = () => {
  if (selectedQuestions.value.length === 0) {
    ElMessage.warning('请选择要导入的题目')
    return
  }
  
  ElMessage.success(`成功导入 ${selectedQuestions.value.length} 道题目`)
  // 这里应该实现实际的导入逻辑
}

// 预览题目
const previewQuestion = (question: any) => {
  previewData.value = question
  showPreview.value = true
}
</script>

<style scoped>
.question-collect-view {
  padding: 24px;
}

.page-header {
  margin-bottom: 24px;
}

.page-title {
  font-size: 24px;
  font-weight: 600;
  color: #1a1a1a;
  margin: 0 0 8px 0;
}

.page-subtitle {
  color: #666;
  margin: 0;
}

.dopamine-card {
  background: white;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  margin-bottom: 24px;
}

.collect-config h3 {
  margin: 0 0 20px 0;
  color: #1a1a1a;
}

.config-form {
  max-width: 800px;
}

.collect-progress h3 {
  margin: 0 0 16px 0;
  color: #1a1a1a;
}

.progress-text {
  margin: 8px 0 0 0;
  color: #666;
  text-align: center;
}

.collect-results h3 {
  margin: 0;
  color: #1a1a1a;
}

.results-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.results-actions {
  display: flex;
  gap: 8px;
}

.quality-text {
  font-size: 12px;
  color: #666;
  margin-left: 4px;
}

.question-preview {
  padding: 16px;
}

.question-info {
  display: flex;
  gap: 8px;
  align-items: center;
  margin-bottom: 16px;
}

.question-content,
.question-options,
.question-answer,
.question-explanation {
  margin-bottom: 16px;
}

.question-content h4,
.question-options h4,
.question-answer h4,
.question-explanation h4 {
  margin: 0 0 8px 0;
  color: #1a1a1a;
  font-size: 14px;
}

.question-content p,
.question-answer p,
.question-explanation p {
  margin: 0;
  color: #333;
  line-height: 1.6;
}

.option-item {
  display: flex;
  align-items: center;
  margin-bottom: 4px;
}

.option-label {
  font-weight: 600;
  margin-right: 8px;
  min-width: 20px;
}
</style>