<template>
  <div class="question-import-view">
    <!-- 页面标题 -->
    <div class="page-header">
      <h1 class="page-title">批量导入题目</h1>
      <p class="page-subtitle">支持Excel、Word等格式的题目批量导入</p>
    </div>

    <!-- 导入步骤 -->
    <div class="import-steps dopamine-card">
      <el-steps :active="currentStep" finish-status="success">
        <el-step title="选择文件" description="上传题目文件"></el-step>
        <el-step title="数据预览" description="预览导入数据"></el-step>
        <el-step title="导入设置" description="配置导入参数"></el-step>
        <el-step title="完成导入" description="确认并导入"></el-step>
      </el-steps>
    </div>

    <!-- 步骤内容 -->
    <div class="step-content dopamine-card">
      <!-- 步骤1：文件上传 -->
      <div v-if="currentStep === 0" class="upload-step">
        <div class="upload-area">
          <el-upload
            ref="uploadRef"
            class="upload-dragger"
            drag
            :auto-upload="false"
            :on-change="handleFileChange"
            :before-upload="beforeUpload"
            accept=".xlsx,.xls,.docx,.doc,.txt"
          >
            <el-icon class="el-icon--upload"><upload-filled /></el-icon>
            <div class="el-upload__text">
              将文件拖到此处，或<em>点击上传</em>
            </div>
            <template #tip>
              <div class="el-upload__tip">
                支持 Excel(.xlsx/.xls)、Word(.docx/.doc)、文本(.txt) 格式
              </div>
            </template>
          </el-upload>
        </div>

        <div class="template-download">
          <h3>下载模板</h3>
          <p>请下载标准模板，按照模板格式整理题目数据</p>
          <div class="template-buttons">
            <el-button type="primary" @click="downloadTemplate('excel')">
              <el-icon><Download /></el-icon>
              Excel模板
            </el-button>
            <el-button @click="downloadTemplate('word')">
              <el-icon><Download /></el-icon>
              Word模板
            </el-button>
          </div>
        </div>

        <div class="step-actions">
          <el-button type="primary" :disabled="!selectedFile" @click="nextStep">
            下一步
          </el-button>
        </div>
      </div>

      <!-- 步骤2：数据预览 -->
      <div v-if="currentStep === 1" class="preview-step">
        <div class="preview-info">
          <el-alert
            title="数据预览"
            :description="`共解析到 ${previewData.length} 道题目，请检查数据是否正确`"
            type="info"
            show-icon
          />
        </div>

        <div class="preview-table">
          <el-table :data="previewData" border style="width: 100%" max-height="400">
            <el-table-column prop="type" label="题型" width="80">
              <template #default="{ row }">
                <el-tag :type="getTypeTagType(row.type)">{{ row.type }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="subject" label="科目" width="100" />
            <el-table-column prop="difficulty" label="难度" width="80">
              <template #default="{ row }">
                <el-rate
                  v-model="row.difficulty"
                  disabled
                  show-score
                  text-color="#ff9900"
                  score-template="{value}"
                />
              </template>
            </el-table-column>
            <el-table-column prop="content" label="题目内容" min-width="200" show-overflow-tooltip />
            <el-table-column prop="answer" label="答案" width="120" show-overflow-tooltip />
            <el-table-column label="状态" width="80">
              <template #default="{ row }">
                <el-tag v-if="row.valid" type="success">有效</el-tag>
                <el-tag v-else type="danger">错误</el-tag>
              </template>
            </el-table-column>
          </el-table>
        </div>

        <div class="step-actions">
          <el-button @click="prevStep">上一步</el-button>
          <el-button type="primary" @click="nextStep">下一步</el-button>
        </div>
      </div>

      <!-- 步骤3：导入设置 -->
      <div v-if="currentStep === 2" class="settings-step">
        <el-form :model="importSettings" label-width="120px">
          <el-form-item label="默认科目">
            <el-select v-model="importSettings.defaultSubject" placeholder="请选择默认科目">
              <el-option
                v-for="subject in subjects"
                :key="subject.id"
                :label="subject.name"
                :value="subject.name"
              />
            </el-select>
          </el-form-item>
          
          <el-form-item label="默认难度">
            <el-rate v-model="importSettings.defaultDifficulty" show-text />
          </el-form-item>
          
          <el-form-item label="重复处理">
            <el-radio-group v-model="importSettings.duplicateHandling">
              <el-radio label="skip">跳过重复题目</el-radio>
              <el-radio label="update">更新重复题目</el-radio>
              <el-radio label="create">创建新题目</el-radio>
            </el-radio-group>
          </el-form-item>
          
          <el-form-item label="错误处理">
            <el-radio-group v-model="importSettings.errorHandling">
              <el-radio label="skip">跳过错误数据</el-radio>
              <el-radio label="stop">遇到错误停止</el-radio>
            </el-radio-group>
          </el-form-item>
        </el-form>

        <div class="step-actions">
          <el-button @click="prevStep">上一步</el-button>
          <el-button type="primary" @click="nextStep">下一步</el-button>
        </div>
      </div>

      <!-- 步骤4：完成导入 -->
      <div v-if="currentStep === 3" class="complete-step">
        <div class="import-summary">
          <el-result
            icon="success"
            title="导入完成"
            :sub-title="`成功导入 ${importResult.success} 道题目，跳过 ${importResult.skipped} 道，失败 ${importResult.failed} 道`"
          >
            <template #extra>
              <el-button type="primary" @click="resetImport">重新导入</el-button>
              <el-button @click="$router.push('/teacher/questions')">返回题库</el-button>
            </template>
          </el-result>
        </div>

        <div v-if="importResult.errors.length > 0" class="error-details">
          <h3>错误详情</h3>
          <el-table :data="importResult.errors" border>
            <el-table-column prop="row" label="行号" width="80" />
            <el-table-column prop="content" label="题目内容" min-width="200" />
            <el-table-column prop="error" label="错误信息" min-width="150" />
          </el-table>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import { UploadFilled, Download } from '@element-plus/icons-vue'

// 当前步骤
const currentStep = ref(0)

// 选中的文件
const selectedFile = ref<File | null>(null)

// 预览数据
const previewData = ref([
  {
    type: '单选题',
    subject: '数据结构',
    difficulty: 3,
    content: '以下哪种数据结构适合实现栈？',
    answer: 'A',
    valid: true
  },
  {
    type: '多选题',
    subject: '算法设计',
    difficulty: 4,
    content: '以下哪些是排序算法？',
    answer: 'ABC',
    valid: true
  }
])

// 科目列表
const subjects = ref([
  { id: 1, name: '数据结构' },
  { id: 2, name: '算法设计' },
  { id: 3, name: '操作系统' },
  { id: 4, name: '计算机网络' },
  { id: 5, name: '数据库原理' }
])

// 导入设置
const importSettings = reactive({
  defaultSubject: '',
  defaultDifficulty: 3,
  duplicateHandling: 'skip',
  errorHandling: 'skip'
})

// 导入结果
const importResult = reactive({
  success: 0,
  skipped: 0,
  failed: 0,
  errors: [] as Array<{ row: number; content: string; error: string }>
})

// 文件变化处理
const handleFileChange = (file: any) => {
  selectedFile.value = file.raw
  ElMessage.success('文件选择成功')
}

// 上传前检查
const beforeUpload = (file: File) => {
  const isValidType = ['application/vnd.openxmlformats-officedocument.spreadsheetml.sheet', 
                      'application/vnd.ms-excel',
                      'application/vnd.openxmlformats-officedocument.wordprocessingml.document',
                      'application/msword',
                      'text/plain'].includes(file.type)
  
  if (!isValidType) {
    ElMessage.error('只支持 Excel、Word、文本格式的文件')
    return false
  }
  
  const isLt10M = file.size / 1024 / 1024 < 10
  if (!isLt10M) {
    ElMessage.error('文件大小不能超过 10MB')
    return false
  }
  
  return true
}

// 下载模板
const downloadTemplate = (type: string) => {
  ElMessage.success(`正在下载${type}模板...`)
  // 这里应该实现实际的模板下载逻辑
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

// 下一步
const nextStep = () => {
  if (currentStep.value < 3) {
    currentStep.value++
    
    // 模拟导入过程
    if (currentStep.value === 3) {
      setTimeout(() => {
        importResult.success = 2
        importResult.skipped = 0
        importResult.failed = 0
      }, 1000)
    }
  }
}

// 上一步
const prevStep = () => {
  if (currentStep.value > 0) {
    currentStep.value--
  }
}

// 重置导入
const resetImport = () => {
  currentStep.value = 0
  selectedFile.value = null
  importResult.success = 0
  importResult.skipped = 0
  importResult.failed = 0
  importResult.errors = []
}
</script>

<style scoped>
.question-import-view {
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

.import-steps {
  margin-bottom: 24px;
}

.upload-step {
  text-align: center;
}

.upload-area {
  margin-bottom: 32px;
}

.upload-dragger {
  width: 100%;
}

.template-download {
  margin-bottom: 32px;
  padding: 24px;
  background: #f8f9fa;
  border-radius: 8px;
}

.template-download h3 {
  margin: 0 0 8px 0;
  color: #1a1a1a;
}

.template-download p {
  margin: 0 0 16px 0;
  color: #666;
}

.template-buttons {
  display: flex;
  gap: 12px;
  justify-content: center;
}

.step-actions {
  display: flex;
  justify-content: center;
  gap: 12px;
  margin-top: 24px;
}

.preview-step .preview-info {
  margin-bottom: 16px;
}

.preview-table {
  margin-bottom: 24px;
}

.settings-step {
  max-width: 600px;
  margin: 0 auto;
}

.complete-step .import-summary {
  margin-bottom: 24px;
}

.error-details h3 {
  margin: 0 0 16px 0;
  color: #1a1a1a;
}
</style>