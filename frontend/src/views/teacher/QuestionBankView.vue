<template>
  <div class="question-bank-view">
    <!-- 页面标题和工具栏 -->
    <div class="page-header">
      <div class="header-content">
        <div class="title-section">
          <h1 class="page-title">
            <el-icon class="title-icon"><Document /></el-icon>
            题库管理
          </h1>
          <p class="page-subtitle">智能管理题目资源，支持批量导入和分类整理</p>
        </div>
        
        <!-- 工具栏移到头部右侧 -->
        <div class="header-toolbar">
          <div class="toolbar-actions">
            <el-button type="info" size="default" @click="loadQuestions" class="action-btn refresh-btn">
              <el-icon><Refresh /></el-icon>
              <span>刷新</span>
            </el-button>
            <el-button type="primary" size="default" @click="showAddQuestionDialog = true" class="action-btn primary-btn">
              <el-icon><Plus /></el-icon>
              <span>添加题目</span>
            </el-button>
            <el-button type="success" size="default" @click="showImportDialog = true" class="action-btn success-btn">
              <el-icon><Upload /></el-icon>
              <span>批量导入</span>
            </el-button>
            <el-button type="warning" size="default" @click="exportQuestions" class="action-btn warning-btn">
              <el-icon><Download /></el-icon>
              <span>导出题目</span>
            </el-button>
            <el-button
              type="danger"
              size="default"
              :disabled="selectedQuestions.length === 0"
              @click="batchDeleteQuestions"
              class="action-btn danger-btn"
            >
              <el-icon><Delete /></el-icon>
              <span>批量删除 ({{ selectedQuestions.length }})</span>
            </el-button>
          </div>
          
          <div class="search-container">
            <el-input
              v-model="searchKeyword"
              placeholder="搜索题目内容、知识点..."
              size="large"
              clearable
              @input="handleSearch"
              class="search-input"
            >
              <template #prefix>
                <el-icon class="search-icon"><Search /></el-icon>
              </template>
            </el-input>
          </div>
        </div>
      </div>
    </div>

        <!-- 统计信息 -->
    <!-- 统计信息 -->
    <div class="stats-section-simple">
      <div class="stats-grid-simple">
        <!-- 科目统计 -->
        <div class="stat-card-simple blue">
          <div class="stat-icon-simple">
            <el-icon><User /></el-icon>
          </div>
          <div class="stat-content-simple">
            <div class="stat-value-simple">10</div>
            <div class="stat-label-simple">科目统计</div>
          </div>
        </div>

        <!-- 题目总数 -->
        <div class="stat-card-simple green">
          <div class="stat-icon-simple">
            <el-icon><ShoppingBag /></el-icon>
          </div>
          <div class="stat-content-simple">
            <div class="stat-value-simple">21</div>
            <div class="stat-label-simple">题目总数</div>
          </div>
        </div>

        <!-- 已发布 -->
        <div class="stat-card-simple red">
          <div class="stat-icon-simple">
            <el-icon><CircleCheck /></el-icon>
          </div>
          <div class="stat-content-simple">
            <div class="stat-value-simple">20</div>
            <div class="stat-label-simple">已发布</div>
          </div>
        </div>

        <!-- 草稿 -->
        <div class="stat-card-simple orange">
          <div class="stat-icon-simple">
            <el-icon><Edit /></el-icon>
          </div>
          <div class="stat-content-simple">
            <div class="stat-value-simple">1</div>
            <div class="stat-label-simple">草稿</div>
          </div>
        </div>
      </div>
    </div>

    <!-- 题目列表 -->
    <div class="question-list-section">
      <!-- 筛选条件 -->
      <div class="filter-bar">
        <div class="filter-row">
          <div class="filter-item">
            <label class="filter-label">
              <el-icon><Folder /></el-icon>
              科目
            </label>
            <el-select v-model="filters.subject" placeholder="选择科目" clearable size="default" class="filter-select">
              <el-option
                v-for="subject in subjects"
                :key="subject.id"
                :label="subject.name"
                :value="subject.id"
              />
            </el-select>
          </div>

          <div class="filter-item">
            <label class="filter-label">
              <el-icon><List /></el-icon>
              类型
            </label>
            <el-select v-model="filters.type" placeholder="选择类型" clearable size="default" class="filter-select">
              <el-option label="单选题" value="single" />
              <el-option label="多选题" value="multiple" />
              <el-option label="判断题" value="judge" />
              <el-option label="填空题" value="fill" />
              <el-option label="简答题" value="essay" />
            </el-select>
          </div>

          <div class="filter-item">
            <label class="filter-label">
              <el-icon><Star /></el-icon>
              难度
            </label>
            <el-select v-model="filters.difficulty" placeholder="选择难度" clearable size="default" class="filter-select">
              <el-option label="简单" value="easy" />
              <el-option label="中等" value="medium" />
              <el-option label="困难" value="hard" />
            </el-select>
          </div>

          <div class="filter-item">
            <label class="filter-label">
              <el-icon><Collection /></el-icon>
              知识点
            </label>
            <el-select v-model="filters.knowledgePoint" placeholder="选择知识点" clearable size="default" class="filter-select">
              <el-option
                v-for="point in knowledgePoints"
                :key="point.id"
                :label="point.name"
                :value="point.id"
              />
            </el-select>
          </div>
          
          <div class="filter-actions">
            <el-button type="primary" @click="applyFilters" size="default">
              <el-icon><Search /></el-icon>
              筛选
            </el-button>
            <el-button @click="resetFilters" size="default">
              <el-icon><Refresh /></el-icon>
              重置
            </el-button>
          </div>
        </div>
      </div>
      
      <div class="question-list-container modern-card">
        <el-table
          :data="filteredQuestions"
          style="width: 100%"
          @selection-change="handleSelectionChange"
          row-key="id"
          v-loading="loading"
          class="modern-table"
          stripe
          :header-cell-style="{ background: '#f8fafc', color: '#374151', fontWeight: '800' }"
          :show-overflow-tooltip="false"
          height="100%"
          virtual-scrolling
          :virtual-scroll-height="400"
        >
          <el-table-column type="selection" width="55" />
          <el-table-column prop="id" label="题目ID" min-width="100" align="center">
            <template #default="{ row }">
              <el-tag type="info" size="small">#{{ row.id }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column label="题目内容" min-width="300">
            <template #default="{ row }">
              <div class="question-content-cell-compact">
                <el-popover
                  placement="top"
                  :width="400"
                  trigger="hover"
                  :disabled="!row.content || row.content.length <= 30"
                >
                  <template #default>
                    <div class="popover-content">
                      <strong>题目内容：</strong><br>
                      {{ row.content }}
                      <br><br>
                      <small>字符数：{{ row.content ? row.content.length : 0 }}</small>
                    </div>
                  </template>
                  <template #reference>
                    <p class="content-text-compact" :title="row.content">
                      {{ row.content && row.content.length > 30 ? row.content.substring(0, 30) + '...' : (row.content || '') }}
                    </p>
                  </template>
                </el-popover>
                <div v-if="row.knowledge_point" class="knowledge-point-compact">
                  <el-icon><Collection /></el-icon>
                  {{ row.knowledge_point }}
                </div>
              </div>
            </template>
          </el-table-column>
          
          <el-table-column label="难度" width="90" align="center">
            <template #default="{ row }">
              <el-tag
                v-if="getDifficultyText(row.difficulty)"
                :type="getDifficultyTagType(row.difficulty)"
                size="small"
                class="difficulty-tag-compact"
              >
                <el-icon><Star /></el-icon>
                {{ getDifficultyText(row.difficulty) }}
              </el-tag>
            </template>
          </el-table-column>
          
          <el-table-column label="题型" width="80" align="center">
            <template #default="{ row }">
              <el-tag :type="getTypeTagType(row.type)" size="small" class="type-tag-compact">
                {{ getTypeText(row.type) }}
              </el-tag>
            </template>
          </el-table-column>
          
          <el-table-column label="使用次数" width="100" align="center">
            <template #default="{ row }">
              <el-tag type="info" size="small">{{ row.usage_count || 0 }}次</el-tag>
            </template>
          </el-table-column>
          
          <el-table-column label="正确率" width="90" align="center">
            <template #default="{ row }">
              <el-tag :type="getCorrectRateTagType(row.correct_rate)" size="small">
                {{ formatCorrectRate(row.correct_rate) }}
              </el-tag>
            </template>
          </el-table-column>
          
          <el-table-column label="创建日期" width="110" align="center">
            <template #default="{ row }">
              <span class="date-text">{{ formatDate(row.created_at) }}</span>
            </template>
          </el-table-column>
          
          <el-table-column prop="subject" label="科目" min-width="120" align="center">
            <template #default="{ row }">
              <div class="subject-cell">
                <el-tag v-if="row.subject" type="primary" class="subject-tag">
                  <el-icon><Folder /></el-icon>
                  {{ row.subject.name || row.subject }}
                </el-tag>
                <span v-else class="no-data">-</span>
              </div>
            </template>
          </el-table-column>
          

          
          <el-table-column prop="status" label="状态" min-width="100" align="center">
            <template #default="{ row }">
              <div class="status-cell">
                <el-tag :type="getStatusTagType(row.status)" class="status-tag">
                  <el-icon><CircleCheck v-if="row.status === 'published'" /><Edit v-else /></el-icon>
                  {{ getStatusText(row.status) }}
                </el-tag>
              </div>
            </template>
          </el-table-column>
          
          <el-table-column label="操作" width="260" fixed="right" align="center">
            <template #default="{ row }">
              <div class="action-buttons">
                <el-button type="primary" size="small" @click="editQuestion(row)" class="action-btn-modern">
                  <el-icon><Edit /></el-icon>
                  编辑
                </el-button>
                <el-button type="success" size="small" @click="previewQuestion(row)" class="action-btn-modern">
                  <el-icon><View /></el-icon>
                  预览
                </el-button>
                <el-button type="danger" size="small" @click="deleteQuestionItem(row.id)" class="action-btn-modern">
                  <el-icon><Delete /></el-icon>
                  删除
                </el-button>
              </div>
            </template>
          </el-table-column>
        </el-table>
        
        <!-- 分页 -->
        <div class="pagination-wrapper">
          <el-pagination
            v-model:current-page="currentPage"
            v-model:page-size="pageSize"
            :page-sizes="[10, 20, 50, 100]"
            :total="totalQuestions"
            layout="total, sizes, prev, pager, next, jumper"
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
          />
        </div>
      </div>
    </div>
    </div>

    <!-- 添加题目对话框 -->
    <el-dialog
      v-model="showAddQuestionDialog"
      title="添加题目"
      width="90vw"
      :before-close="handleCloseAddDialog"
      class="question-dialog"
      top="5vh"
    >
      <QuestionForm
        :question="currentQuestion"
        @save="handleSaveQuestion"
        @cancel="showAddQuestionDialog = false"
      />
    </el-dialog>

    <!-- 批量导入对话框 -->
    <el-dialog v-model="showImportDialog" title="批量导入题目" width="600px">
      <div class="import-section">
        <div class="import-tips">
          <el-alert title="导入说明" type="info" :closable="false" show-icon>
            <template #default>
              <ul>
                <li>支持 Excel (.xlsx) 和 CSV (.csv) 格式</li>
                <li>请按照模板格式准备数据</li>
                <li>单次最多导入 1000 道题目</li>
              </ul>
            </template>
          </el-alert>
        </div>

        <div class="import-actions">
          <el-button type="primary" @click="downloadTemplate">
            <el-icon><Download /></el-icon>
            下载模板
          </el-button>
        </div>

        <el-upload
          class="upload-demo"
          drag
          :auto-upload="false"
          :on-change="handleFileChange"
          :accept="'.xlsx,.csv'"
          :limit="1"
        >
          <el-icon class="el-icon--upload"><UploadFilled /></el-icon>
          <div class="el-upload__text">将文件拖到此处，或<em>点击上传</em></div>
          <template #tip>
            <div class="el-upload__tip">只能上传 xlsx/csv 文件，且不超过 10MB</div>
          </template>
        </el-upload>

        <div class="import-progress" v-if="importProgress.show">
          <el-progress
            :percentage="importProgress.percentage"
            :status="importProgress.status"
          />
          <div class="progress-text">{{ importProgress.text }}</div>
        </div>
      </div>

      <template #footer>
        <div class="dialog-footer">
          <el-button @click="showImportDialog = false">取消</el-button>
          <el-button
            type="primary"
            :loading="importing"
            :disabled="!selectedFile"
            @click="startImport"
          >
            开始导入
          </el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 题目预览对话框 -->
    <el-dialog v-model="showPreviewDialog" title="题目预览" width="95vw" top="2.5vh" class="preview-dialog">
      <QuestionPreview v-if="previewQuestionData" :question="previewQuestionData" />
    </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Plus, Upload, Download, Delete, Search, Refresh,
  Document, CircleCheck, Edit, Star, UploadFilled,
  Filter, Folder, List, Collection, TrendCharts, Minus,
  View, Clock, Grid, Menu, PieChart
} from '@element-plus/icons-vue'
import QuestionForm from '@/components/QuestionForm.vue'
import QuestionPreview from '@/components/QuestionPreview.vue'
import {
  getQuestions,
  getQuestion,
  createQuestion,
  updateQuestion,
  deleteQuestion,
  getQuestionStats,
  getSubjects,
  type Question,
  type QuestionListParams
} from '@/api/question'

// 搜索关键词
const searchKeyword = ref('')

// 筛选条件
const filters = reactive({
  subject: '',
  type: '',
  difficulty: '',
  knowledgePoint: ''
})

// 科目列表
const subjects = ref<any[]>([])

// 知识点列表
const knowledgePoints = ref([
  { id: 1, name: '线性表' },
  { id: 2, name: '栈和队列' },
  { id: 3, name: '树和二叉树' },
  { id: 4, name: '图' },
  { id: 5, name: '排序算法' },
  { id: 6, name: '查找算法' }
])

// 题目统计
const questionStats = ref({
  total: 0,
  published: 0,
  draft: 0,
  avgDifficulty: 0
})

// 选中的题目
const selectedQuestions = ref([])

// 分页
const currentPage = ref(1)
const pageSize = ref(20)
const totalQuestions = ref(0)

// 加载状态
const loading = ref(false)

// 视图模式
const viewMode = ref('table')

// 对话框状态
const showAddQuestionDialog = ref(false)
const showImportDialog = ref(false)
const showPreviewDialog = ref(false)

// 当前编辑的题目
const currentQuestion = ref(null)
const previewQuestionData = ref(null)

// 导入相关
const importing = ref(false)
const selectedFile = ref(null)
const importProgress = reactive({
  show: false,
  percentage: 0,
  status: '',
  text: ''
})

// 题目数据
const questions = ref<Question[]>([])

// 加载题目列表
const loadQuestions = async () => {
  try {
    loading.value = true
    const params: QuestionListParams = {
      page: currentPage.value,
      size: pageSize.value,
      subject_id: filters.subject || undefined,
      type: filters.type || undefined,
      difficulty: filters.difficulty || undefined,
      search: searchKeyword.value || undefined
    }

    const response = await getQuestions(params)
    console.log('API响应数据:', response) // 调试日志
    questions.value = response?.questions || []
    totalQuestions.value = response?.total || 0
  } catch (error) {
    console.error('加载题目列表失败:', error)
    ElMessage.error('加载题目列表失败')
  } finally {
    loading.value = false
  }
}

// 加载科目列表
const loadSubjects = async () => {
  try {
    const response = await getSubjects()
    subjects.value = response.subjects || []
  } catch (error) {
    console.error('加载科目列表失败:', error)
  }
}

// 加载统计数据
const loadStats = async () => {
  try {
    const response = await getQuestionStats()
    if (response.data) {
      questionStats.value = response.data
    }
  } catch (error) {
    console.error('加载统计数据失败:', error)
  }
}

// 筛选后的题目列表（现在筛选在后端处理）
const filteredQuestions = computed(() => {
  return questions.value
})

// 获取类型文本
const getTypeText = (type: string): string => {
  const textMap: Record<string, string> = {
    single: '单选',
    multiple: '多选',
    judge: '判断',
    fill: '填空',
    essay: '简答'
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
const getDifficultyText = (difficulty: string | number): string => {
  // 处理数字类型的难度值
  if (typeof difficulty === 'number') {
    const numMap: Record<number, string> = {
      1: '简单',
      2: '中等',
      3: '困难'
    }
    return numMap[difficulty] || ''
  }

  // 处理字符串类型的难度值
  const textMap: Record<string, string> = {
    easy: '简单',
    medium: '中等',
    hard: '困难'
  }
  return textMap[difficulty] || ''
}

// 获取难度标签类型
const getDifficultyTagType = (difficulty: string | number): string => {
  // 处理数字类型的难度值
  if (typeof difficulty === 'number') {
    const numMap: Record<number, string> = {
      1: 'success',
      2: 'warning',
      3: 'danger'
    }
    return numMap[difficulty] || 'info'
  }

  // 处理字符串类型的难度值
  const typeMap: Record<string, string> = {
    easy: 'success',
    medium: 'warning',
    hard: 'danger'
  }
  return typeMap[difficulty] || 'info'
}

// 获取状态文本
const getStatusText = (status: string): string => {
  const textMap: Record<string, string> = {
    published: '已发布',
    draft: '草稿'
  }
  return textMap[status] || ''
}

// 获取状态标签类型
const getStatusTagType = (status: string): string => {
  const typeMap: Record<string, string> = {
    published: 'success',
    draft: 'warning'
  }
  return typeMap[status] || ''
}

// 格式化正确率
const formatCorrectRate = (rate: number): string => {
  if (rate === null || rate === undefined) return '-'
  return `${rate.toFixed(1)}%`
}

// 获取正确率样式类
const getCorrectRateClass = (rate: number): string => {
  if (rate >= 90) return 'rate-excellent'
  if (rate >= 80) return 'rate-good'
  if (rate >= 70) return 'rate-fair'
  return 'rate-poor'
}

// 获取正确率标签类型
const getCorrectRateTagType = (rate: number): string => {
  if (!rate) return 'info'
  if (rate >= 80) return 'success'
  if (rate >= 60) return 'warning'
  return 'danger'
}

// 格式化日期
const formatDate = (date: Date | string | null | undefined): string => {
  if (!date) {
    return '-'
  }

  // 如果是字符串，转换为Date对象
  const dateObj = typeof date === 'string' ? new Date(date) : date

  // 检查是否是有效日期
  if (isNaN(dateObj.getTime())) {
    return '-'
  }

  return dateObj.toLocaleDateString('zh-CN')
}

// 获取题型数量
const getQuestionTypeCount = (): number => {
  const types = new Set(questions.value.map(q => q.type))
  return types.size
}

// 获取题型分布数据
const getQuestionTypeDistribution = () => {
  const typeCount: Record<string, number> = {}
  const total = questions.value.length
  
  // 统计各题型数量
  questions.value.forEach(question => {
    typeCount[question.type] = (typeCount[question.type] || 0) + 1
  })
  
  // 转换为百分比并添加颜色
  const colors = ['#409eff', '#67c23a', '#e6a23c', '#f56c6c', '#909399']
  let colorIndex = 0
  
  return Object.entries(typeCount).map(([type, count]) => ({
    type: getTypeText(type),
    count,
    percentage: total > 0 ? Math.round((count / total) * 100) : 0,
    color: colors[colorIndex++ % colors.length]
  }))
}

// 判断内容是否溢出（已移除，直接在模板中使用字符长度判断）

// 处理搜索
const handleSearch = () => {
  loadQuestions()
}

// 应用筛选
const applyFilters = () => {
  currentPage.value = 1
  loadQuestions()
}

// 重置筛选
const resetFilters = () => {
  filters.subject = ''
  filters.type = ''
  filters.difficulty = ''
  filters.knowledgePoint = ''
  searchKeyword.value = ''
  currentPage.value = 1
  loadQuestions()
}

// 处理选择变化
const handleSelectionChange = (selection: any[]) => {
  selectedQuestions.value = selection
}

// 编辑题目
const editQuestion = (question: any) => {
  currentQuestion.value = { ...question }
  showAddQuestionDialog.value = true
}

// 预览题目
const previewQuestion = (question: any) => {
  // 处理数据格式，确保选择题的options是对象数组
  const processedQuestion = { ...question }
  
  if (question.type === 'single' || question.type === 'multiple') {
    // 如果options是字符串，解析为数组
    if (typeof question.options === 'string') {
      try {
        const optionsArray = JSON.parse(question.options)
        // 将字符串数组转换为对象数组
        processedQuestion.options = optionsArray.map((text: string, index: number) => ({
          text,
          isCorrect: false // 预览时不需要显示正确答案
        }))
      } catch (error) {
        console.error('解析选项失败:', error)
        processedQuestion.options = []
      }
    } else if (Array.isArray(question.options)) {
      // 如果已经是数组，确保每个选项都有text属性
      processedQuestion.options = question.options.map((option: any, index: number) => {
        if (typeof option === 'string') {
          return { text: option, isCorrect: false }
        }
        return { ...option, isCorrect: false }
      })
    }
  }
  
  previewQuestionData.value = processedQuestion
  showPreviewDialog.value = true
}

// 删除题目
const deleteQuestionItem = async (questionId: number) => {
  try {
    await ElMessageBox.confirm(
      '确定要删除这道题目吗？删除后无法恢复。',
      '删除确认',
      {
        confirmButtonText: '确定删除',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    await deleteQuestion(questionId)
    ElMessage.success('题目删除成功')
    loadQuestions()
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('删除题目失败:', error)
      ElMessage.error('删除题目失败')
    }
  }
}

// 批量删除题目
const batchDeleteQuestions = async () => {
  try {
    await ElMessageBox.confirm(
      `确定要删除选中的 ${selectedQuestions.value.length} 道题目吗？删除后无法恢复。`,
      '批量删除确认',
      {
        confirmButtonText: '确定删除',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    ElMessage.success(`成功删除 ${selectedQuestions.value.length} 道题目`)
    selectedQuestions.value = []
  } catch {
    // 用户取消
  }
}

// 导出题目
const exportQuestions = () => {
  ElMessage.success('题目导出功能开发中')
}

// 处理保存题目
const handleSaveQuestion = async (questionData: any) => {
  try {
    console.log('保存题目数据:', questionData) // 调试日志
    
    // QuestionForm已经处理了数据格式，直接使用
    const processedData = { ...questionData }
    
    console.log('处理后的数据:', processedData) // 调试日志
    
    if (currentQuestion.value?.id) {
      // 更新题目
      console.log('更新题目 ID:', currentQuestion.value.id) // 调试日志
      await updateQuestion(currentQuestion.value.id, processedData)
      ElMessage.success('题目更新成功')
    } else {
      // 创建新题目
      console.log('创建新题目') // 调试日志
      await createQuestion(processedData)
      ElMessage.success('题目添加成功')
    }
    
    showAddQuestionDialog.value = false
    currentQuestion.value = null
    
    // 重新加载题目列表
    await loadQuestions()
  } catch (error) {
    console.error('保存题目失败:', error)
    ElMessage.error('保存题目失败')
  }
}

// 处理关闭添加对话框
const handleCloseAddDialog = () => {
  currentQuestion.value = null
  showAddQuestionDialog.value = false
}

// 下载模板
const downloadTemplate = () => {
  ElMessage.success('模板下载功能开发中')
}

// 处理文件变化
const handleFileChange = (file: any) => {
  selectedFile.value = file
}

// 开始导入
const startImport = () => {
  if (!selectedFile.value) {
    ElMessage.warning('请选择要导入的文件')
    return
  }

  importing.value = true
  importProgress.show = true
  importProgress.percentage = 0
  importProgress.status = 'active'
  importProgress.text = '正在解析文件...'

  // 模拟导入过程
  const timer = setInterval(() => {
    importProgress.percentage += 10

    if (importProgress.percentage === 30) {
      importProgress.text = '正在验证数据格式...'
    } else if (importProgress.percentage === 60) {
      importProgress.text = '正在保存题目...'
    } else if (importProgress.percentage === 90) {
      importProgress.text = '正在更新索引...'
    } else if (importProgress.percentage >= 100) {
      clearInterval(timer)
      importProgress.status = 'success'
      importProgress.text = '导入完成！'
      importing.value = false

      setTimeout(() => {
        showImportDialog.value = false
        importProgress.show = false
        selectedFile.value = null
        ElMessage.success('题目导入成功')
      }, 1000)
    }
  }, 300)
}

// 分页处理
const handleSizeChange = (size: number) => {
  pageSize.value = size
  currentPage.value = 1
  loadQuestions()
}

const handleCurrentChange = (page: number) => {
  currentPage.value = page
  loadQuestions()
}

// 初始化数据
onMounted(() => {
  loadSubjects()
  loadQuestions()
  loadStats()
})
</script>

<style scoped>
/* 页面整体布局 */
.question-bank-view {
  padding: var(--spacing-lg);
  background: var(--bg-secondary);
  min-height: 100vh;
}

/* 页面头部样式 */
.page-header {
  background: var(--gradient-primary);
  border-radius: var(--radius-xl);
  padding: 20px; /* 从 var(--spacing-xl) 减少约15% */
  margin-bottom: var(--spacing-lg);
  color: white;
  box-shadow: var(--shadow-lg);
  min-height: calc(120px - 30px); /* 压缩30px高度 */
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 24px; /* 从 var(--spacing-xl) 减少约20% */
}

.title-section {
  flex: 1;
}

.header-toolbar {
  display: flex;
  flex-direction: column;
  gap: 12px; /* 从 var(--spacing-md) 减少约20% */
  align-items: flex-end;
}

.toolbar-actions {
  display: flex;
  gap: 8px; /* 从 var(--spacing-sm) 减少约20% */
  flex-wrap: wrap;
  justify-content: flex-end;
}

.search-container {
  min-width: 300px;
}

.page-title {
  display: flex;
  align-items: center;
  gap: 12px; /* 从 var(--spacing-md) 减少约20% */
  font-size: 26px; /* 从32px减少约20% */
  font-weight: 700;
  margin: 0 0 8px 0; /* 从 var(--spacing-sm) 减少约20% */
  color: white;
}

.title-icon {
  font-size: 29px; /* 从36px减少约20% */
}

.page-subtitle {
  font-size: 14px; /* 从16px减少约15% */
  color: rgba(255, 255, 255, 0.9);
  margin: 0;
}

.header-stats {
  display: flex;
  gap: var(--spacing-xl);
}

.stat-item {
  text-align: center;
}

.stat-number {
  display: block;
  font-size: 28px;
  font-weight: 700;
  color: white;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.8);
}

/* 工具栏样式 */
.toolbar-section {
  margin-bottom: var(--spacing-lg);
}

.toolbar-card {
  background: white;
  border-radius: var(--radius-xl);
  padding: var(--spacing-lg);
  box-shadow: var(--shadow-md);
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: var(--spacing-lg);
}

.toolbar-left {
  display: flex;
  gap: var(--spacing-md);
  flex-wrap: wrap;
}

.action-btn {
  border-radius: var(--radius-lg);
  font-weight: 600;
  transition: all 0.3s ease;
  border: none;
  box-shadow: var(--shadow-sm);
}

.action-btn:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-lg);
}

.primary-btn {
  background: var(--gradient-primary);
}

.success-btn {
  background: var(--gradient-success);
}

.warning-btn {
  background: var(--gradient-warning);
}

.danger-btn {
  background: linear-gradient(135deg, var(--dopamine-red), #FF8A80);
}

.refresh-btn {
  background: linear-gradient(135deg, #17a2b8, #20c997);
}

.toolbar-right {
  flex-shrink: 0;
}

.search-container {
  position: relative;
}

.search-input {
  width: 300px;
}

.search-icon {
  color: var(--text-secondary);
}

/* 筛选条样式 - 紧凑版 */
.filter-bar {
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  padding: 8px 12px;
  margin-bottom: 12px;
}

.filter-row {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
}

.filter-item {
  display: flex;
  align-items: center;
  gap: 6px;
  min-width: 0;
}

.filter-label {
  font-weight: 500;
  color: #374151;
  font-size: 13px;
  display: flex;
  align-items: center;
  gap: 3px;
  white-space: nowrap;
}

.filter-label .el-icon {
  color: #6366f1;
  font-size: 14px;
}

.filter-select {
  width: 140px;
}

.filter-actions {
  display: flex;
  gap: 6px;
  margin-left: auto;
}

@media (max-width: 768px) {
  .filter-row {
    flex-direction: column;
    align-items: stretch;
    gap: 8px;
  }
  
  .filter-item {
    justify-content: space-between;
  }
  
  .filter-select {
    width: 180px;
  }
  
  .filter-actions {
    margin-left: 0;
    justify-content: center;
  }
}

/* 统计卡片样式 */
.stats-section {
  margin-bottom: var(--spacing-lg);
}

.stats-container {
  background: white;
  border-radius: var(--radius-xl);
  padding: var(--spacing-lg);
  box-shadow: var(--shadow-md);
}

.stats-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--spacing-lg);
  padding-bottom: var(--spacing-md);
  border-bottom: 1px solid var(--border-light);
}

.stats-title {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: var(--spacing-lg);
}

.modern-card {
  background: white;
  border-radius: var(--radius-lg);
  padding: var(--spacing-md);
  box-shadow: var(--shadow-md);
  transition: all 0.3s ease;
  border: 1px solid var(--border-light);
}
.stat-card {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
}

.stat-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.stat-icon {
  width: 60px;
  height: 60px;
  border-radius: var(--radius-lg);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  flex-shrink: 0;
  font-size: 24px;
}

.primary-gradient {
  background: var(--gradient-primary);
}

.success-gradient {
  background: var(--gradient-success);
}

.warning-gradient {
  background: var(--gradient-warning);
}

.info-gradient {
  background: linear-gradient(135deg, var(--dopamine-blue), #64B5F6);
}

.stat-trend {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  font-weight: 600;
}

.trend-up {
  color: var(--dopamine-green);
}

.trend-down {
  color: var(--dopamine-red);
}

.trend-stable {
  color: var(--text-secondary);
}

.trend-text {
  color: inherit;
}

.stat-content {
  flex: 1;
}

.stat-value {
  font-size: 32px;
  font-weight: 700;
  color: var(--text-primary);
  margin-bottom: 4px;
  line-height: 1;
}

.stat-label {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 4px;
}

.stat-desc {
  font-size: 12px;
  color: var(--text-secondary);
  line-height: 1.4;
}

/* 题目列表样式 */
.question-list-section {
  margin-bottom: var(--spacing-lg);
}

.list-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--spacing-md);
  padding: var(--spacing-lg);
  border-bottom: 1px solid var(--border-light);
  background: var(--bg-light);
}

.list-title {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0;
}

.list-actions {
  display: flex;
  gap: var(--spacing-sm);
}
/* 表格容器样式 */
.question-list-container {
  background: white;
  border-radius: var(--radius-lg);
  padding: var(--spacing-md);
  box-shadow: var(--shadow-md);
  overflow: hidden;
  height: calc(100vh - 420px);
  display: flex;
  flex-direction: column;
}

.modern-table {
  width: 100%;
  border-collapse: separate;
  border-spacing: 0;
  flex: 1;
  overflow: hidden;
}

.modern-table :deep(.el-table__body-wrapper) {
  overflow-y: auto;
  overflow-x: hidden;
}

.modern-table :deep(.el-scrollbar__wrap) {
  overflow-x: hidden;
}

.modern-table :deep(.el-table__header) {
  background: var(--bg-light);
}

.modern-table :deep(.el-table__header th) {
  background: var(--bg-light) !important;
  color: var(--text-primary);
  font-weight: 600;
  border-bottom: 2px solid var(--border-light);
  padding: 12px 8px;
  font-size: 13px;
}

.modern-table :deep(.el-table__body tr) {
  transition: none;
}
 
.modern-table :deep(.el-table__body td) {
  border-bottom: 1px solid var(--border-light);
  padding: 10px 8px;
  vertical-align: middle;
}

/* 题目内容样式 */
.question-content {
  max-width: 380px;
}

.question-content-cell {
  padding: 12px 8px;
}

.question-header-row {
  display: flex;
  gap: 8px;
  margin-bottom: 8px;
  flex-wrap: wrap;
}

.difficulty-tag,
.type-tag {
  display: flex;
  align-items: center;
  gap: 4px;
  font-weight: 600;
}

.question-text-content {
  margin-bottom: 8px;
}

.question-text-content .content-text {
  margin: 0;
  font-size: 14px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  white-space: normal;
  line-height: 1.5;
  color: var(--text-primary);
  cursor: pointer;
  transition: color 0.2s;
}

.content-text:hover {
  color: var(--primary-color);
}

.question-meta-row {
  display: flex;
  gap: 12px;
  font-size: 12px;
  color: var(--text-secondary);
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 4px;
}

.subject-cell {
  display: flex;
  justify-content: center;
}

.subject-tag {
  display: flex;
  align-items: center;
  gap: 4px;
}

.no-data {
  color: var(--text-secondary);
  font-style: italic;
}

.stats-cell {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.stat-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}

.stat-label {
  font-size: 12px;
  color: var(--text-secondary);
  white-space: nowrap;
}

.status-cell {
  display: flex;
  justify-content: center;
}

.status-tag {
  display: flex;
  align-items: center;
  gap: 4px;
}

/* 紧凑表格样式 */
.question-content-cell-compact {
  padding: 8px 0;
}

.content-text-compact {
  margin: 0;
  font-size: 14px;
  line-height: 1.4;
  color: var(--text-primary);
  cursor: pointer;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.knowledge-point-compact {
  display: flex;
  align-items: center;
  gap: 4px;
  margin-top: 6px;
  font-size: 12px;
  color: var(--text-secondary);
}

.difficulty-tag-compact,
.type-tag-compact {
  display: flex;
  align-items: center;
  gap: 2px;
  font-size: 12px;
}

.date-text {
  font-size: 12px;
  color: var(--text-secondary);
}

.popover-content {
  max-width: 350px;
  line-height: 1.6;
  word-wrap: break-word;
  white-space: pre-wrap;
}

.action-buttons {
  display: flex;
  gap: 6px;
  justify-content: center;
  flex-wrap: nowrap;
}

.action-btn-modern {
  display: flex;
  align-items: center;
  gap: 4px;
  border-radius: 4px;
  transition: all 0.2s ease;
  font-weight: 500;
  white-space: nowrap;
  min-width: 65px;
  font-size: 12px;
  padding: 8px 12px;
  cursor: pointer;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.action-btn-modern.el-button--primary {
  background: linear-gradient(135deg, #1976d2 0%, #1565c0 100%);
  color: white;
  border: 1px solid #1565c0;
}

.action-btn-modern.el-button--primary:hover {
  background: linear-gradient(135deg, #0d47a1 0%, #1565c0 100%);
  color: white;
  border-color: #0d47a1;
  transform: translateY(-1px);
  box-shadow: 0 3px 8px rgba(13, 71, 161, 0.4);
}

.action-btn-modern.el-button--success {
  background: linear-gradient(135deg, #388e3c 0%, #2e7d32 100%);
  color: white;
  border: 1px solid #2e7d32;
}

.action-btn-modern.el-button--success:hover {
  background: linear-gradient(135deg, #1b5e20 0%, #2e7d32 100%);
  color: white;
  border-color: #1b5e20;
  transform: translateY(-1px);
  box-shadow: 0 3px 8px rgba(27, 94, 32, 0.4);
}

.action-btn-modern.el-button--danger {
  background: linear-gradient(135deg, #d32f2f 0%, #c62828 100%);
  color: white;
  border: 1px solid #c62828;
}

.action-btn-modern.el-button--danger:hover {
  background: linear-gradient(135deg, #b71c1c 0%, #c62828 100%);
  color: white;
  border-color: #b71c1c;
  transform: translateY(-1px);
  box-shadow: 0 3px 8px rgba(183, 28, 28, 0.4);
}



.ml-2 {
  margin-left: var(--spacing-sm);
}

.pagination-wrapper {
  display: flex;
  justify-content: center;
  margin-top: var(--spacing-md);
  padding: var(--spacing-md);
  background: white;
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-md);
}

.rate-excellent {
  color: var(--dopamine-green);
  font-weight: 600;
}

.rate-good {
  color: var(--dopamine-blue);
  font-weight: 600;
}

.rate-fair {
  color: var(--dopamine-orange);
  font-weight: 600;
}

.rate-poor {
  color: var(--dopamine-red);
  font-weight: 600;
}

.import-section {
  padding: var(--spacing-lg) 0;
}

.import-tips {
  margin-bottom: var(--spacing-lg);
}

.import-tips ul {
  margin: 0;
  padding-left: var(--spacing-lg);
}

.import-tips li {
  margin-bottom: var(--spacing-sm);
}

.import-actions {
  margin-bottom: var(--spacing-lg);
  text-align: center;
}

.upload-demo {
  margin-bottom: var(--spacing-lg);
}

.import-progress {
  margin-top: var(--spacing-lg);
}

.progress-text {
  text-align: center;
  margin-top: var(--spacing-sm);
  font-size: 14px;
  color: var(--text-secondary);
}

.dialog-footer {
  text-align: right;
}

/* 正确率样式 */
.rate-excellent {
  color: var(--dopamine-green);
  font-weight: 600;
}

.rate-good {
  color: var(--dopamine-blue);
  font-weight: 600;
}

.rate-fair {
  color: var(--dopamine-orange);
  font-weight: 600;
}

.rate-poor {
  color: var(--dopamine-red);
  font-weight: 600;
}

/* 导入相关样式 */
.import-section {
  padding: var(--spacing-lg) 0;
}

.import-tips {
  margin-bottom: var(--spacing-lg);
}

.import-tips ul {
  margin: 0;
  padding-left: var(--spacing-lg);
}

.import-tips li {
  margin-bottom: var(--spacing-sm);
}

.import-actions {
  margin-bottom: var(--spacing-lg);
  text-align: center;
}

.upload-demo {
  margin-bottom: var(--spacing-lg);
}

.import-progress {
  margin-top: var(--spacing-lg);
}

.progress-text {
  text-align: center;
  margin-top: var(--spacing-sm);
  font-size: 14px;
  color: var(--text-secondary);
}

.dialog-footer {
  text-align: right;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .toolbar-card {
    flex-direction: column;
    gap: var(--spacing-md);
    align-items: stretch;
  }

  .toolbar-left {
    flex-wrap: wrap;
    justify-content: center;
  }

  .filter-row {
    flex-direction: column;
    align-items: stretch;
  }

  .filter-group {
    flex-direction: column;
  }

  .filter-actions {
    justify-content: center;
  }

  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }

  .header-content {
    flex-direction: column;
    text-align: center;
    gap: var(--spacing-lg);
    align-items: center;
  }

  .header-toolbar {
    width: 100%;
    align-items: center;
  }

  .toolbar-actions {
    justify-content: center;
    flex-wrap: wrap;
  }

  .search-container {
    width: 100%;
    min-width: auto;
  }

  .search-input {
    width: 100%;
  }
}

@media (max-width: 480px) {
  .stats-grid {
    grid-template-columns: 1fr;
  }

  .toolbar-left {
    flex-direction: column;
  }

  .action-buttons {
    flex-direction: column;
  }

  .page-title {
    font-size: 24px;
  }

  .title-icon {
    font-size: 28px;
  }
}

/* 题目对话框样式 */
:deep(.question-dialog) {
  .el-dialog {
    max-width: 1600px;
    height: 90vh;
    margin: 0 auto;
    border-radius: var(--radius-lg);
  }
  
  .el-dialog__header {
    padding: var(--spacing-lg) var(--spacing-xl);
    border-bottom: 1px solid var(--border-color);
    background: var(--bg-primary);
  }
  
  .el-dialog__title {
    font-size: 18px;
    font-weight: 600;
    color: var(--text-primary);
  }
  
  .el-dialog__body {
    padding: 0;
    height: calc(90vh - 120px);
    overflow: hidden;
  }
}

.question-dialog-content {
  height: 100%;
  padding: var(--spacing-lg);
  overflow-y: auto;
}

/* 优化表单在对话框中的显示 */
:deep(.question-dialog .question-form) {
  height: 100%;
}

:deep(.question-dialog .form-layout) {
  height: calc(100% - 80px);
}

:deep(.question-dialog .form-actions) {
  position: sticky;
  bottom: 0;
  background: white;
  border-top: 1px solid var(--border-color);
  margin-top: var(--spacing-lg);
  padding-top: var(--spacing-lg);
  z-index: 10;
}

/* 响应式对话框 */
@media (max-width: 1200px) {
  :deep(.question-dialog .el-dialog) {
    width: 95vw !important;
    height: 95vh;
  }
  
  .question-dialog-content {
    padding: var(--spacing-md);
  }
}

@media (max-width: 768px) {
  :deep(.question-dialog .el-dialog) {
    width: 100vw !important;
    height: 100vh;
    margin: 0;
    border-radius: 0;
  }
  
  :deep(.question-dialog .form-layout) {
    flex-direction: column;
    height: auto;
  }
  
  :deep(.question-dialog .form-left) {
    flex: none;
  }
}

/* 重新设计的统计样式 */
.stats-section-redesigned {
  margin-bottom: 24px;
}

.stats-container-redesigned {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 16px;
  padding: 20px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  height: 180px;
  overflow: hidden;
}

.stats-header-redesigned {
  margin-bottom: 16px;
}

.stats-title-redesigned {
  font-size: 16px;
  font-weight: 600;
  color: var(--el-text-color-primary);
  display: flex;
  align-items: center;
  gap: 6px;
  margin: 0;
}

.stats-grid-redesigned {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  grid-template-rows: repeat(2, 1fr);
  gap: 8px;
  height: calc(180px - 52px);
}

.stat-card-redesigned {
  background: linear-gradient(135deg, #f8fafc 0%, #e2e8f0 100%);
  border: 1px solid rgba(203, 213, 225, 0.6);
  border-radius: 8px;
  padding: 6px 10px;
  display: flex;
  align-items: flex-start;
  gap: 8px;
  transition: all 0.3s ease;
  position: relative;
  overflow: hidden;
  backdrop-filter: blur(8px);
  aspect-ratio: 16 / 9;
  min-height: 0;
}

.stat-card-redesigned::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 2px;
  background: linear-gradient(90deg, var(--card-accent-color, #409eff), transparent);
}

.stat-card-redesigned.primary {
  --card-accent-color: #409eff;
  background: linear-gradient(135deg, #e3f2fd 0%, #bbdefb 100%);
}

.stat-card-redesigned.success {
  --card-accent-color: #67c23a;
  background: linear-gradient(135deg, #e8f5e8 0%, #c8e6c9 100%);
}

.stat-card-redesigned.warning {
  --card-accent-color: #e6a23c;
  background: linear-gradient(135deg, #fff3e0 0%, #ffcc80 100%);
}

.stat-card-redesigned.info {
  --card-accent-color: #17a2b8;
  background: linear-gradient(135deg, #e0f2f1 0%, #b2dfdb 100%);
}

.stat-card-redesigned.danger {
  --card-accent-color: #f56c6c;
  background: linear-gradient(135deg, #ffebee 0%, #ffcdd2 100%);
}

.stat-card-redesigned.secondary {
  --card-accent-color: #909399;
  background: linear-gradient(135deg, #f5f5f5 0%, #e0e0e0 100%);
}

.stat-card-redesigned:hover {
  transform: translateY(-1px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
  border-color: var(--card-accent-color, #409eff);
}

.stat-icon-redesigned {
  width: 24px;
  height: 24px;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, var(--card-accent-color, #409eff), rgba(var(--card-accent-color, #409eff), 0.8));
  color: white;
  font-size: 12px;
  flex-shrink: 0;
  margin-top: 1px;
}

.stat-content-redesigned {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.stat-value-redesigned {
  font-size: 16px;
  font-weight: 700;
  color: var(--el-text-color-primary);
  line-height: 1;
  margin-bottom: 1px;
}

.stat-label-redesigned {
  font-size: 11px;
  color: var(--el-text-color-primary);
  font-weight: 600;
  line-height: 1;
  margin-bottom: 1px;
}

.stat-desc-redesigned {
  font-size: 9px;
  color: var(--el-text-color-regular);
  line-height: 1.1;
  margin-bottom: 2px;
}

.stat-chart-redesigned {
  display: flex;
  gap: 1px;
  height: 3px;
  border-radius: 1.5px;
  overflow: hidden;
  background: rgba(0, 0, 0, 0.1);
  margin-top: 1px;
}

.chart-bar {
  height: 100%;
  min-width: 1px;
  transition: all 0.3s ease;
}

/* 简洁卡片样式 - 仿照截图设计 */
.stats-section-simple {
  margin-bottom: 20px; /* 从24px减少约15% */
}

.stats-grid-simple {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 12px; /* 从16px减少约25% */
}

.stat-card-simple {
  background: white;
  border-radius: 6px; /* 从8px减少约25% */
  padding: 12px; /* 从16px减少约25% */
  display: flex;
  align-items: center;
  gap: 10px; /* 从12px减少约15% */
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
  min-height: 48px; /* 从60px减少约20% */
  cursor: pointer;
}

.stat-card-simple:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
}



.stat-icon-simple {
  width: 40px; /* 从48px减少约15% */
  height: 40px; /* 从48px减少约15% */
  border-radius: 6px; /* 从8px减少约25% */
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 16px; /* 从20px减少约20% */
  flex-shrink: 0;
}

.stat-card-simple.blue .stat-icon-simple {
  background: #2196F3;
}

.stat-card-simple.green .stat-icon-simple {
  background: #4CAF50;
}

.stat-card-simple.red .stat-icon-simple {
  background: #F44336;
}

.stat-card-simple.orange .stat-icon-simple {
  background: #FF9800;
}

.stat-content-simple {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 1px; /* 从2px减少50% */
}

.stat-value-simple {
  font-size: 20px; /* 从24px减少约15% */
  font-weight: 700;
  color: #333;
  line-height: 1;
}

.stat-label-simple {
  font-size: 11px; /* 从12px减少约10% */
  color: #666;
  font-weight: 500;
  line-height: 1;
}

/* 响应式设计 */
@media (max-width: 1200px) {
  .stats-grid-simple {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 768px) {
  .stats-grid-simple {
    grid-template-columns: 1fr;
  }
  
  .stat-card-simple {
    padding: 12px;
    gap: 10px;
  }
  
  .stat-icon-simple {
    width: 40px;
    height: 40px;
    font-size: 18px;
  }
  
  .stat-value-simple {
    font-size: 20px;
  }
}

/* 响应式调整 */
@media (max-width: 1200px) {
  .stats-grid-redesigned {
    grid-template-columns: repeat(2, 1fr);
    grid-template-rows: repeat(3, 1fr);
    gap: 10px;
  }
  
  .stats-container-redesigned {
    height: auto;
    min-height: 200px;
  }
  
  .stats-grid-redesigned {
    height: auto;
  }
}

@media (max-width: 768px) {
  .stats-grid-redesigned {
    grid-template-columns: 1fr;
    grid-template-rows: repeat(6, 1fr);
  }
  
  .stat-card-redesigned {
    padding: 8px 10px;
  }
  
  .stat-icon-redesigned {
    width: 28px;
    height: 28px;
    font-size: 12px;
  }
  
  .stat-value-redesigned {
    font-size: 18px;
  }
}

@media (max-width: 480px) {
  .stats-container-redesigned {
    padding: 16px;
  }
  
  .stat-card-redesigned {
    padding: 6px 8px;
    gap: 8px;
  }
}

/* 预览对话框样式 */
:deep(.preview-dialog .el-dialog) {
  height: 95vh;
  max-height: 95vh;
  margin: 0;
  display: flex;
  flex-direction: column;
}

:deep(.preview-dialog .el-dialog__header) {
  flex-shrink: 0;
  padding: 16px 20px;
  border-bottom: 1px solid var(--el-border-color-light);
}

:deep(.preview-dialog .el-dialog__body) {
  flex: 1;
  padding: 0;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

:deep(.preview-dialog .question-preview) {
  height: 100%;
  max-width: none;
  aspect-ratio: unset;
  display: flex;
  flex-direction: column;
  overflow-y: auto;
  padding: 20px;
}

/* 响应式预览对话框 */
@media (max-width: 768px) {
  :deep(.preview-dialog .el-dialog) {
    width: 98vw !important;
    height: 98vh;
    top: 1vh !important;
  }
  
  :deep(.preview-dialog .question-preview) {
    padding: 16px;
  }
}
</style>
