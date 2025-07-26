<template>
  <div class="grade-entry">
    <!-- 页面标题 -->
    <div class="page-header">
      <h1 class="page-title">
        <el-icon><EditPen /></el-icon>
        成绩录入
      </h1>
      <p class="page-description">
        手动录入成绩、批量导入成绩、在线阅卷评分
      </p>
    </div>

    <!-- 功能选项卡 -->
    <el-tabs v-model="activeTab" class="grade-tabs">
      <!-- 手动录入 -->
      <el-tab-pane label="手动录入" name="manual">
        <div class="manual-entry">
          <!-- 考试选择 -->
          <div class="exam-selector dopamine-card">
            <h3>选择考试</h3>
            <el-form :model="manualForm" label-width="100px">
              <el-row :gutter="20">
                <el-col :span="8">
                  <el-form-item label="考试名称">
                    <el-select v-model="manualForm.examId" placeholder="请选择考试" @change="loadStudents">
                      <el-option
                        v-for="exam in exams"
                        :key="exam.id"
                        :label="exam.title"
                        :value="exam.id"
                      />
                    </el-select>
                  </el-form-item>
                </el-col>
                <el-col :span="8">
                  <el-form-item label="科目">
                    <el-select v-model="manualForm.subject" placeholder="请选择科目">
                      <el-option
                        v-for="subject in subjects"
                        :key="subject"
                        :label="subject"
                        :value="subject"
                      />
                    </el-select>
                  </el-form-item>
                </el-col>
                <el-col :span="8">
                  <el-form-item label="班级">
                    <el-select v-model="manualForm.classId" placeholder="请选择班级">
                      <el-option
                        v-for="cls in classes"
                        :key="cls.id"
                        :label="cls.name"
                        :value="cls.id"
                      />
                    </el-select>
                  </el-form-item>
                </el-col>
              </el-row>
            </el-form>
          </div>

          <!-- 学生成绩录入表格 -->
          <div class="grade-table dopamine-card" v-if="students.length > 0">
            <div class="table-header">
              <h3>学生成绩录入</h3>
              <div class="table-actions">
                <el-button type="primary" @click="saveGrades" :loading="saving">
                  <el-icon><Check /></el-icon>
                  保存成绩
                </el-button>
                <el-button @click="clearGrades">
                  <el-icon><RefreshLeft /></el-icon>
                  清空
                </el-button>
              </div>
            </div>
            
            <el-table :data="students" stripe border>
              <el-table-column prop="studentNumber" label="学号" width="120" />
              <el-table-column prop="name" label="姓名" width="120" />
              <el-table-column prop="class" label="班级" width="120" />
              <el-table-column label="成绩" width="150">
                <template #default="{ row, $index }">
                  <el-input-number
                    v-model="row.score"
                    :min="0"
                    :max="100"
                    :precision="1"
                    placeholder="请输入成绩"
                    @change="validateScore(row, $index)"
                  />
                </template>
              </el-table-column>
              <el-table-column label="等级" width="100">
                <template #default="{ row }">
                  <el-tag :type="getGradeType(row.score)">{{ getGradeLevel(row.score) }}</el-tag>
                </template>
              </el-table-column>
              <el-table-column label="备注">
                <template #default="{ row }">
                  <el-input
                    v-model="row.remark"
                    placeholder="请输入备注"
                    maxlength="100"
                  />
                </template>
              </el-table-column>
            </el-table>
          </div>
        </div>
      </el-tab-pane>

      <!-- 批量导入 -->
      <el-tab-pane label="批量导入" name="import">
        <div class="batch-import">
          <!-- 导入说明 -->
          <div class="import-guide dopamine-card">
            <h3>批量导入说明</h3>
            <div class="guide-content">
              <el-steps :active="importStep" align-center>
                <el-step title="下载模板" description="下载Excel模板文件" />
                <el-step title="填写数据" description="按模板格式填写成绩数据" />
                <el-step title="上传文件" description="上传填写好的Excel文件" />
                <el-step title="确认导入" description="预览数据并确认导入" />
              </el-steps>
            </div>
          </div>

          <!-- 模板下载 -->
          <div class="template-download dopamine-card">
            <h3>模板下载</h3>
            <div class="download-content">
              <p>请先下载成绩导入模板，按照模板格式填写学生成绩数据。</p>
              <el-button type="primary" @click="downloadTemplate">
                <el-icon><Download /></el-icon>
                下载Excel模板
              </el-button>
            </div>
          </div>

          <!-- 文件上传 -->
          <div class="file-upload dopamine-card">
            <h3>文件上传</h3>
            <el-upload
              ref="uploadRef"
              class="upload-demo"
              drag
              :auto-upload="false"
              :on-change="handleFileChange"
              :before-upload="beforeUpload"
              accept=".xlsx,.xls"
            >
              <el-icon class="el-icon--upload"><UploadFilled /></el-icon>
              <div class="el-upload__text">
                将Excel文件拖拽到此处，或<em>点击上传</em>
              </div>
              <template #tip>
                <div class="el-upload__tip">
                  支持.xlsx、.xls格式，文件大小不超过10MB
                </div>
              </template>
            </el-upload>
          </div>

          <!-- 数据预览 -->
          <div class="data-preview dopamine-card" v-if="importData.length > 0">
            <div class="preview-header">
              <h3>数据预览</h3>
              <div class="preview-actions">
                <el-button type="success" @click="confirmImport" :loading="importing">
                  <el-icon><Check /></el-icon>
                  确认导入 ({{ importData.length }}条)
                </el-button>
                <el-button @click="clearImportData">
                  <el-icon><Close /></el-icon>
                  取消
                </el-button>
              </div>
            </div>
            
            <el-table :data="importData.slice(0, 10)" stripe border max-height="400">
              <el-table-column prop="studentNumber" label="学号" width="120" />
              <el-table-column prop="name" label="姓名" width="120" />
              <el-table-column prop="class" label="班级" width="120" />
              <el-table-column prop="subject" label="科目" width="100" />
              <el-table-column prop="score" label="成绩" width="100" />
              <el-table-column prop="remark" label="备注" />
              <el-table-column label="状态" width="100">
                <template #default="{ row }">
                  <el-tag :type="row.valid ? 'success' : 'danger'">
                    {{ row.valid ? '有效' : '无效' }}
                  </el-tag>
                </template>
              </el-table-column>
            </el-table>
            
            <div class="preview-footer" v-if="importData.length > 10">
              <p>仅显示前10条数据，共{{ importData.length }}条数据</p>
            </div>
          </div>
        </div>
      </el-tab-pane>

      <!-- 在线阅卷 -->
      <el-tab-pane label="在线阅卷" name="grading">
        <div class="online-grading">
          <!-- 考试选择 -->
          <div class="exam-selector dopamine-card">
            <h3>选择待阅卷考试</h3>
            <el-form :model="gradingForm" label-width="100px">
              <el-row :gutter="20">
                <el-col :span="12">
                  <el-form-item label="考试名称">
                    <el-select v-model="gradingForm.examId" placeholder="请选择考试" @change="loadGradingPapers">
                      <el-option
                        v-for="exam in pendingExams"
                        :key="exam.id"
                        :label="exam.title"
                        :value="exam.id"
                      >
                        <span>{{ exam.title }}</span>
                        <span style="float: right; color: #8492a6; font-size: 13px">
                          待阅卷: {{ exam.pendingCount }}
                        </span>
                      </el-option>
                    </el-select>
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="题目类型">
                    <el-select v-model="gradingForm.questionType" placeholder="请选择题目类型">
                      <el-option label="简答题" value="essay" />
                      <el-option label="论述题" value="discussion" />
                      <el-option label="计算题" value="calculation" />
                    </el-select>
                  </el-form-item>
                </el-col>
              </el-row>
            </el-form>
          </div>

          <!-- 阅卷界面 -->
          <div class="grading-interface dopamine-card" v-if="currentPaper">
            <div class="grading-header">
              <h3>在线阅卷</h3>
              <div class="grading-progress">
                <span>进度: {{ gradedCount }}/{{ totalPapers }}</span>
                <el-progress :percentage="gradingProgress" :stroke-width="8" />
              </div>
            </div>

            <div class="grading-content">
              <div class="paper-info">
                <p><strong>学生:</strong> {{ currentPaper.studentName }} ({{ currentPaper.studentNumber }})</p>
                <p><strong>题目:</strong> {{ currentPaper.questionTitle }}</p>
                <p><strong>分值:</strong> {{ currentPaper.totalScore }}分</p>
              </div>

              <div class="answer-content">
                <h4>学生答案:</h4>
                <div class="answer-text">{{ currentPaper.answer }}</div>
              </div>

              <div class="grading-panel">
                <el-form :model="gradingForm" label-width="80px">
                  <el-form-item label="得分">
                    <el-input-number
                      v-model="gradingForm.score"
                      :min="0"
                      :max="currentPaper.totalScore"
                      :precision="1"
                    />
                    <span class="score-hint">/ {{ currentPaper.totalScore }}分</span>
                  </el-form-item>
                  <el-form-item label="评语">
                    <el-input
                      v-model="gradingForm.comment"
                      type="textarea"
                      :rows="3"
                      placeholder="请输入评语"
                      maxlength="200"
                    />
                  </el-form-item>
                </el-form>

                <div class="grading-actions">
                  <el-button type="primary" @click="submitGrading" :loading="grading">
                    <el-icon><Check /></el-icon>
                    提交评分
                  </el-button>
                  <el-button @click="skipPaper">
                    <el-icon><Right /></el-icon>
                    跳过
                  </el-button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  EditPen, Check, RefreshLeft, Download, UploadFilled,
  Close, Right
} from '@element-plus/icons-vue'
import * as XLSX from 'xlsx'

// 当前激活的选项卡
const activeTab = ref('manual')

// 手动录入相关数据
const manualForm = reactive({
  examId: '',
  subject: '',
  classId: ''
})

const students = ref([])
const saving = ref(false)

// 批量导入相关数据
const importStep = ref(0)
const importData = ref([])
const importing = ref(false)
const uploadRef = ref()

// 在线阅卷相关数据
const gradingForm = reactive({
  examId: '',
  questionType: '',
  score: 0,
  comment: ''
})

const currentPaper = ref(null)
const gradedCount = ref(0)
const totalPapers = ref(0)
const grading = ref(false)

// 基础数据
const exams = ref([
  { id: 1, title: '数据结构期末考试', subject: '计算机科学', date: '2024-01-15' },
  { id: 2, title: '算法设计与分析', subject: '计算机科学', date: '2024-01-20' },
  { id: 3, title: '数据库原理', subject: '计算机科学', date: '2024-01-25' }
])

const subjects = ref(['计算机科学', '数学', '英语', '物理', '化学'])

const classes = ref([
  { id: 1, name: '计算机2021-1班' },
  { id: 2, name: '计算机2021-2班' },
  { id: 3, name: '软件工程2021-1班' }
])

const pendingExams = ref([
  { id: 1, title: '数据结构期末考试', pendingCount: 25 },
  { id: 2, title: '算法设计与分析', pendingCount: 18 },
  { id: 3, title: '数据库原理', pendingCount: 32 }
])

// 计算属性
const gradingProgress = computed(() => {
  if (totalPapers.value === 0) return 0
  return Math.round((gradedCount.value / totalPapers.value) * 100)
})

// 加载学生列表
const loadStudents = () => {
  // 模拟加载学生数据
  students.value = [
    { id: 1, studentNumber: '2021001', name: '张三', class: '计算机2021-1班', score: null, remark: '' },
    { id: 2, studentNumber: '2021002', name: '李四', class: '计算机2021-1班', score: null, remark: '' },
    { id: 3, studentNumber: '2021003', name: '王五', class: '计算机2021-1班', score: null, remark: '' }
  ]
}

// 验证成绩
const validateScore = (row, index) => {
  if (row.score < 0 || row.score > 100) {
    ElMessage.warning('成绩应在0-100之间')
    row.score = null
  }
}

// 获取成绩等级
const getGradeLevel = (score) => {
  if (score === null || score === undefined) return '-'
  if (score >= 90) return '优秀'
  if (score >= 80) return '良好'
  if (score >= 70) return '中等'
  if (score >= 60) return '及格'
  return '不及格'
}

// 获取成绩等级类型
const getGradeType = (score) => {
  if (score === null || score === undefined) return ''
  if (score >= 90) return 'success'
  if (score >= 80) return 'primary'
  if (score >= 70) return 'warning'
  if (score >= 60) return 'info'
  return 'danger'
}

// 保存成绩
const saveGrades = async () => {
  const validGrades = students.value.filter(s => s.score !== null && s.score !== undefined)
  if (validGrades.length === 0) {
    ElMessage.warning('请至少录入一个学生的成绩')
    return
  }

  saving.value = true
  try {
    // 模拟保存API调用
    await new Promise(resolve => setTimeout(resolve, 1000))
    ElMessage.success(`成功保存${validGrades.length}个学生的成绩`)
  } catch (error) {
    ElMessage.error('保存失败，请重试')
  } finally {
    saving.value = false
  }
}

// 清空成绩
const clearGrades = () => {
  students.value.forEach(student => {
    student.score = null
    student.remark = ''
  })
  ElMessage.info('已清空所有成绩')
}

// 下载模板
const downloadTemplate = () => {
  const templateData = [
    ['学号', '姓名', '班级', '科目', '成绩', '备注'],
    ['2021001', '张三', '计算机2021-1班', '数据结构', '85', ''],
    ['2021002', '李四', '计算机2021-1班', '数据结构', '92', '优秀'],
    ['2021003', '王五', '计算机2021-1班', '数据结构', '78', '']
  ]

  const ws = XLSX.utils.aoa_to_sheet(templateData)
  const wb = XLSX.utils.book_new()
  XLSX.utils.book_append_sheet(wb, ws, '成绩导入模板')
  XLSX.writeFile(wb, '成绩导入模板.xlsx')
  
  ElMessage.success('模板下载成功')
  importStep.value = 1
}

// 文件上传前验证
const beforeUpload = (file) => {
  const isExcel = file.type === 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' || 
                  file.type === 'application/vnd.ms-excel'
  const isLt10M = file.size / 1024 / 1024 < 10

  if (!isExcel) {
    ElMessage.error('只能上传Excel文件!')
    return false
  }
  if (!isLt10M) {
    ElMessage.error('文件大小不能超过10MB!')
    return false
  }
  return true
}

// 处理文件变化
const handleFileChange = (file) => {
  if (file.raw) {
    const reader = new FileReader()
    reader.onload = (e) => {
      try {
        const data = new Uint8Array(e.target.result)
        const workbook = XLSX.read(data, { type: 'array' })
        const worksheet = workbook.Sheets[workbook.SheetNames[0]]
        const jsonData = XLSX.utils.sheet_to_json(worksheet, { header: 1 })
        
        // 解析数据
        const parsedData = jsonData.slice(1).map((row, index) => ({
          id: index + 1,
          studentNumber: row[0],
          name: row[1],
          class: row[2],
          subject: row[3],
          score: row[4],
          remark: row[5] || '',
          valid: row[0] && row[1] && row[4] !== undefined && row[4] >= 0 && row[4] <= 100
        })).filter(item => item.studentNumber)
        
        importData.value = parsedData
        importStep.value = 3
        ElMessage.success(`成功解析${parsedData.length}条数据`)
      } catch (error) {
        ElMessage.error('文件解析失败，请检查文件格式')
      }
    }
    reader.readAsArrayBuffer(file.raw)
  }
}

// 确认导入
const confirmImport = async () => {
  const validData = importData.value.filter(item => item.valid)
  if (validData.length === 0) {
    ElMessage.warning('没有有效的数据可以导入')
    return
  }

  importing.value = true
  try {
    // 模拟导入API调用
    await new Promise(resolve => setTimeout(resolve, 2000))
    ElMessage.success(`成功导入${validData.length}条成绩数据`)
    clearImportData()
  } catch (error) {
    ElMessage.error('导入失败，请重试')
  } finally {
    importing.value = false
  }
}

// 清空导入数据
const clearImportData = () => {
  importData.value = []
  importStep.value = 0
  uploadRef.value?.clearFiles()
}

// 加载待阅卷试卷
const loadGradingPapers = () => {
  // 模拟加载待阅卷试卷
  currentPaper.value = {
    id: 1,
    studentName: '张三',
    studentNumber: '2021001',
    questionTitle: '请简述数据结构中栈和队列的区别',
    totalScore: 10,
    answer: '栈是后进先出(LIFO)的数据结构，只能在栈顶进行插入和删除操作。队列是先进先出(FIFO)的数据结构，在队尾插入元素，在队头删除元素。栈的应用包括函数调用、表达式求值等，队列的应用包括任务调度、广度优先搜索等。'
  }
  
  gradedCount.value = 0
  totalPapers.value = 25
  gradingForm.score = 0
  gradingForm.comment = ''
}

// 提交评分
const submitGrading = async () => {
  if (gradingForm.score === null || gradingForm.score === undefined) {
    ElMessage.warning('请输入得分')
    return
  }

  grading.value = true
  try {
    // 模拟提交评分API调用
    await new Promise(resolve => setTimeout(resolve, 500))
    
    gradedCount.value++
    ElMessage.success('评分提交成功')
    
    // 加载下一份试卷
    if (gradedCount.value < totalPapers.value) {
      loadNextPaper()
    } else {
      ElMessage.success('所有试卷已阅卷完成！')
      currentPaper.value = null
    }
  } catch (error) {
    ElMessage.error('提交失败，请重试')
  } finally {
    grading.value = false
  }
}

// 跳过试卷
const skipPaper = () => {
  if (gradedCount.value < totalPapers.value - 1) {
    loadNextPaper()
    ElMessage.info('已跳过当前试卷')
  } else {
    ElMessage.warning('这是最后一份试卷')
  }
}

// 加载下一份试卷
const loadNextPaper = () => {
  // 模拟加载下一份试卷
  const nextStudent = `学生${gradedCount.value + 1}`
  currentPaper.value = {
    ...currentPaper.value,
    studentName: nextStudent,
    studentNumber: `202100${gradedCount.value + 2}`,
    answer: `这是${nextStudent}的答案内容...`
  }
  
  gradingForm.score = 0
  gradingForm.comment = ''
}

onMounted(() => {
  // 组件挂载时的初始化操作
})
</script>

<style scoped>
.grade-entry {
  padding: var(--spacing-lg);
}

.page-header {
  margin-bottom: var(--spacing-xl);
}

.page-title {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  font-size: 24px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 var(--spacing-sm) 0;
}

.page-description {
  color: var(--text-secondary);
  margin: 0;
}

.grade-tabs {
  background: white;
  border-radius: var(--radius-lg);
  padding: var(--spacing-lg);
  box-shadow: var(--shadow-sm);
}

.dopamine-card {
  background: white;
  border-radius: var(--radius-lg);
  padding: var(--spacing-lg);
  margin-bottom: var(--spacing-lg);
  box-shadow: var(--shadow-sm);
  border: 1px solid var(--border-light);
}

.dopamine-card h3 {
  margin: 0 0 var(--spacing-md) 0;
  color: var(--text-primary);
  font-weight: 600;
}

.table-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--spacing-md);
}

.table-actions {
  display: flex;
  gap: var(--spacing-sm);
}

.guide-content {
  margin-top: var(--spacing-md);
}

.download-content {
  text-align: center;
  padding: var(--spacing-lg);
}

.download-content p {
  margin-bottom: var(--spacing-md);
  color: var(--text-secondary);
}

.upload-demo {
  margin-top: var(--spacing-md);
}

.preview-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--spacing-md);
}

.preview-actions {
  display: flex;
  gap: var(--spacing-sm);
}

.preview-footer {
  text-align: center;
  padding: var(--spacing-md);
  color: var(--text-secondary);
  border-top: 1px solid var(--border-light);
  margin-top: var(--spacing-md);
}

.grading-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--spacing-lg);
}

.grading-progress {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
  min-width: 200px;
}

.grading-content {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--spacing-lg);
}

.paper-info {
  background: var(--bg-light);
  padding: var(--spacing-md);
  border-radius: var(--radius-md);
  margin-bottom: var(--spacing-md);
}

.paper-info p {
  margin: var(--spacing-xs) 0;
}

.answer-content {
  grid-column: 1 / -1;
}

.answer-content h4 {
  margin: 0 0 var(--spacing-sm) 0;
  color: var(--text-primary);
}

.answer-text {
  background: var(--bg-light);
  padding: var(--spacing-md);
  border-radius: var(--radius-md);
  line-height: 1.6;
  min-height: 120px;
}

.grading-panel {
  background: var(--bg-light);
  padding: var(--spacing-lg);
  border-radius: var(--radius-md);
}

.score-hint {
  margin-left: var(--spacing-sm);
  color: var(--text-secondary);
}

.grading-actions {
  display: flex;
  gap: var(--spacing-sm);
  margin-top: var(--spacing-md);
}

@media (max-width: 768px) {
  .grading-content {
    grid-template-columns: 1fr;
  }
  
  .grading-header {
    flex-direction: column;
    gap: var(--spacing-md);
  }
  
  .grading-progress {
    width: 100%;
  }
}
</style>