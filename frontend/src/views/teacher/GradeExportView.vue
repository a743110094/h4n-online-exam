<template>
  <div class="grade-export">
    <!-- 页面标题 -->
    <div class="page-header">
      <h1 class="page-title">
        <el-icon><Download /></el-icon>
        成绩导出
      </h1>
      <p class="page-description">
        导出成绩单、生成报告、批量打印
      </p>
    </div>

    <!-- 导出选项 -->
    <div class="export-options dopamine-card">
      <h3>导出选项</h3>
      <el-form :model="exportForm" label-width="120px">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="导出类型">
              <el-radio-group v-model="exportForm.type">
                <el-radio label="grades">成绩单</el-radio>
                <el-radio label="report">成绩报告</el-radio>
                <el-radio label="statistics">统计分析</el-radio>
                <el-radio label="transcript">成绩单据</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="导出格式">
              <el-radio-group v-model="exportForm.format">
                <el-radio label="excel">Excel</el-radio>
                <el-radio label="pdf">PDF</el-radio>
                <el-radio label="csv">CSV</el-radio>
                <el-radio label="word">Word</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
    </div>

    <!-- 筛选条件 -->
    <div class="filter-conditions dopamine-card">
      <h3>筛选条件</h3>
      <el-form :model="filterForm" :inline="true" label-width="100px">
        <el-form-item label="考试名称">
          <el-select v-model="filterForm.examId" placeholder="请选择考试" clearable multiple style="width: 250px">
            <el-option
              v-for="exam in exams"
              :key="exam.id"
              :label="exam.title"
              :value="exam.id"
            />
          </el-select>
        </el-form-item>
        
        <el-form-item label="科目">
          <el-select v-model="filterForm.subjects" placeholder="请选择科目" clearable multiple style="width: 200px">
            <el-option
              v-for="subject in subjects"
              :key="subject"
              :label="subject"
              :value="subject"
            />
          </el-select>
        </el-form-item>
        
        <el-form-item label="班级">
          <el-select v-model="filterForm.classIds" placeholder="请选择班级" clearable multiple style="width: 200px">
            <el-option
              v-for="cls in classes"
              :key="cls.id"
              :label="cls.name"
              :value="cls.id"
            />
          </el-select>
        </el-form-item>
        
        <el-form-item label="考试日期">
          <el-date-picker
            v-model="filterForm.dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            format="YYYY-MM-DD"
            value-format="YYYY-MM-DD"
            style="width: 240px"
          />
        </el-form-item>
        
        <el-form-item label="成绩范围">
          <el-input-number
            v-model="filterForm.minScore"
            :min="0"
            :max="100"
            placeholder="最低分"
            style="width: 100px"
          />
          <span style="margin: 0 8px">-</span>
          <el-input-number
            v-model="filterForm.maxScore"
            :min="0"
            :max="100"
            placeholder="最高分"
            style="width: 100px"
          />
        </el-form-item>
      </el-form>
    </div>

    <!-- 导出内容配置 -->
    <div class="export-config dopamine-card">
      <h3>导出内容配置</h3>
      <el-row :gutter="20">
        <el-col :span="12">
          <h4>基本信息</h4>
          <el-checkbox-group v-model="exportForm.basicFields">
            <el-checkbox label="studentNumber">学号</el-checkbox>
            <el-checkbox label="studentName">姓名</el-checkbox>
            <el-checkbox label="className">班级</el-checkbox>
            <el-checkbox label="examTitle">考试名称</el-checkbox>
            <el-checkbox label="subject">科目</el-checkbox>
            <el-checkbox label="examDate">考试日期</el-checkbox>
          </el-checkbox-group>
        </el-col>
        <el-col :span="12">
          <h4>成绩信息</h4>
          <el-checkbox-group v-model="exportForm.gradeFields">
            <el-checkbox label="score">成绩</el-checkbox>
            <el-checkbox label="grade">等级</el-checkbox>
            <el-checkbox label="rank">排名</el-checkbox>
            <el-checkbox label="classRank">班级排名</el-checkbox>
            <el-checkbox label="percentile">百分位</el-checkbox>
            <el-checkbox label="remark">备注</el-checkbox>
          </el-checkbox-group>
        </el-col>
      </el-row>
      
      <el-row :gutter="20" style="margin-top: 20px">
        <el-col :span="12">
          <h4>统计信息</h4>
          <el-checkbox-group v-model="exportForm.statisticsFields">
            <el-checkbox label="average">平均分</el-checkbox>
            <el-checkbox label="highest">最高分</el-checkbox>
            <el-checkbox label="lowest">最低分</el-checkbox>
            <el-checkbox label="median">中位数</el-checkbox>
            <el-checkbox label="standardDeviation">标准差</el-checkbox>
            <el-checkbox label="passRate">及格率</el-checkbox>
          </el-checkbox-group>
        </el-col>
        <el-col :span="12">
          <h4>图表选项</h4>
          <el-checkbox-group v-model="exportForm.chartOptions">
            <el-checkbox label="scoreDistribution">分数分布图</el-checkbox>
            <el-checkbox label="gradeDistribution">等级分布图</el-checkbox>
            <el-checkbox label="trendChart">成绩趋势图</el-checkbox>
            <el-checkbox label="comparisonChart">对比分析图</el-checkbox>
          </el-checkbox-group>
        </el-col>
      </el-row>
    </div>

    <!-- 模板选择 -->
    <div class="template-selection dopamine-card">
      <h3>模板选择</h3>
      <el-row :gutter="20">
        <el-col :span="6" v-for="template in templates" :key="template.id">
          <div 
            class="template-card" 
            :class="{ active: exportForm.templateId === template.id }"
            @click="selectTemplate(template.id)"
          >
            <div class="template-preview">
              <img :src="template.preview" :alt="template.name" />
            </div>
            <div class="template-info">
              <h4>{{ template.name }}</h4>
              <p>{{ template.description }}</p>
            </div>
          </div>
        </el-col>
      </el-row>
    </div>

    <!-- 预览区域 -->
    <div class="preview-area dopamine-card" v-if="previewData">
      <div class="preview-header">
        <h3>导出预览</h3>
        <div class="preview-actions">
          <el-button @click="refreshPreview" :loading="previewing">
            <el-icon><RefreshLeft /></el-icon>
            刷新预览
          </el-button>
          <el-button type="primary" @click="startExport" :loading="exporting">
            <el-icon><Download /></el-icon>
            开始导出
          </el-button>
        </div>
      </div>
      
      <div class="preview-content">
        <!-- 成绩单预览 -->
        <div v-if="exportForm.type === 'grades'" class="grades-preview">
          <div class="preview-title">{{ previewData.title }}</div>
          <el-table :data="previewData.data.slice(0, 5)" border>
            <el-table-column 
              v-for="field in selectedFields" 
              :key="field.prop"
              :prop="field.prop"
              :label="field.label"
              :width="field.width"
            />
          </el-table>
          <div class="preview-note" v-if="previewData.data.length > 5">
            仅显示前5条数据，实际导出包含全部 {{ previewData.data.length }} 条数据
          </div>
        </div>
        
        <!-- 报告预览 -->
        <div v-else-if="exportForm.type === 'report'" class="report-preview">
          <div class="report-header">
            <h2>{{ previewData.title }}</h2>
            <p>生成时间: {{ previewData.generateTime }}</p>
          </div>
          <div class="report-summary">
            <h3>考试概况</h3>
            <el-row :gutter="20">
              <el-col :span="6" v-for="stat in previewData.summary" :key="stat.label">
                <div class="summary-item">
                  <div class="summary-value">{{ stat.value }}</div>
                  <div class="summary-label">{{ stat.label }}</div>
                </div>
              </el-col>
            </el-row>
          </div>
          <div class="report-charts">
            <h3>图表分析</h3>
            <p>将包含所选的图表内容...</p>
          </div>
        </div>
      </div>
    </div>

    <!-- 导出历史 -->
    <div class="export-history dopamine-card">
      <div class="history-header">
        <h3>导出历史</h3>
        <el-button @click="clearHistory" size="small">
          <el-icon><Delete /></el-icon>
          清空历史
        </el-button>
      </div>
      
      <el-table :data="exportHistory" stripe>
        <el-table-column prop="fileName" label="文件名" min-width="200" />
        <el-table-column prop="type" label="类型" width="100">
          <template #default="{ row }">
            <el-tag>{{ getTypeLabel(row.type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="format" label="格式" width="80" />
        <el-table-column prop="size" label="大小" width="100" />
        <el-table-column prop="createTime" label="导出时间" width="160" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">{{ getStatusLabel(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150">
          <template #default="{ row }">
            <el-button 
              type="text" 
              size="small" 
              @click="downloadFile(row)"
              :disabled="row.status !== 'completed'"
            >
              <el-icon><Download /></el-icon>
              下载
            </el-button>
            <el-button type="text" size="small" @click="deleteHistory(row)">
              <el-icon><Delete /></el-icon>
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Download, RefreshLeft, Delete
} from '@element-plus/icons-vue'
import * as XLSX from 'xlsx'
import jsPDF from 'jspdf'

// 导出表单
const exportForm = reactive({
  type: 'grades',
  format: 'excel',
  templateId: 1,
  basicFields: ['studentNumber', 'studentName', 'className', 'examTitle', 'subject'],
  gradeFields: ['score', 'grade', 'rank'],
  statisticsFields: ['average', 'highest', 'lowest'],
  chartOptions: []
})

// 筛选表单
const filterForm = reactive({
  examId: [],
  subjects: [],
  classIds: [],
  dateRange: null,
  minScore: null,
  maxScore: null
})

// 状态
const previewing = ref(false)
const exporting = ref(false)
const previewData = ref(null)

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

// 模板数据
const templates = ref([
  {
    id: 1,
    name: '标准模板',
    description: '包含基本信息和成绩的标准格式',
    preview: 'https://trae-api-sg.mchost.guru/api/ide/v1/text_to_image?prompt=clean%20academic%20grade%20report%20template%20with%20student%20information%20table&image_size=square'
  },
  {
    id: 2,
    name: '详细模板',
    description: '包含详细统计和图表的完整报告',
    preview: 'https://trae-api-sg.mchost.guru/api/ide/v1/text_to_image?prompt=detailed%20academic%20report%20template%20with%20charts%20and%20statistics&image_size=square'
  },
  {
    id: 3,
    name: '简洁模板',
    description: '简洁明了的成绩单格式',
    preview: 'https://trae-api-sg.mchost.guru/api/ide/v1/text_to_image?prompt=minimalist%20grade%20sheet%20template%20clean%20design&image_size=square'
  },
  {
    id: 4,
    name: '官方模板',
    description: '符合学校官方要求的正式模板',
    preview: 'https://trae-api-sg.mchost.guru/api/ide/v1/text_to_image?prompt=official%20school%20transcript%20template%20formal%20layout&image_size=square'
  }
])

// 导出历史
const exportHistory = ref([
  {
    id: 1,
    fileName: '数据结构期末考试成绩单.xlsx',
    type: 'grades',
    format: 'excel',
    size: '156KB',
    createTime: '2024-01-16 14:30:00',
    status: 'completed'
  },
  {
    id: 2,
    fileName: '算法设计成绩报告.pdf',
    type: 'report',
    format: 'pdf',
    size: '2.3MB',
    createTime: '2024-01-15 16:45:00',
    status: 'completed'
  },
  {
    id: 3,
    fileName: '计算机2021-1班统计分析.xlsx',
    type: 'statistics',
    format: 'excel',
    size: '89KB',
    createTime: '2024-01-14 10:20:00',
    status: 'failed'
  }
])

// 计算属性
const selectedFields = computed(() => {
  const fields = []
  
  if (exportForm.basicFields.includes('studentNumber')) {
    fields.push({ prop: 'studentNumber', label: '学号', width: '120' })
  }
  if (exportForm.basicFields.includes('studentName')) {
    fields.push({ prop: 'studentName', label: '姓名', width: '120' })
  }
  if (exportForm.basicFields.includes('className')) {
    fields.push({ prop: 'className', label: '班级', width: '150' })
  }
  if (exportForm.basicFields.includes('examTitle')) {
    fields.push({ prop: 'examTitle', label: '考试名称', width: '180' })
  }
  if (exportForm.basicFields.includes('subject')) {
    fields.push({ prop: 'subject', label: '科目', width: '100' })
  }
  if (exportForm.gradeFields.includes('score')) {
    fields.push({ prop: 'score', label: '成绩', width: '100' })
  }
  if (exportForm.gradeFields.includes('grade')) {
    fields.push({ prop: 'grade', label: '等级', width: '100' })
  }
  if (exportForm.gradeFields.includes('rank')) {
    fields.push({ prop: 'rank', label: '排名', width: '80' })
  }
  
  return fields
})

// 选择模板
const selectTemplate = (templateId) => {
  exportForm.templateId = templateId
  refreshPreview()
}

// 刷新预览
const refreshPreview = async () => {
  previewing.value = true
  try {
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    if (exportForm.type === 'grades') {
      previewData.value = {
        title: '成绩单',
        data: [
          {
            studentNumber: '2021001',
            studentName: '张三',
            className: '计算机2021-1班',
            examTitle: '数据结构期末考试',
            subject: '计算机科学',
            score: 92,
            grade: '优秀',
            rank: 1
          },
          {
            studentNumber: '2021002',
            studentName: '李四',
            className: '计算机2021-1班',
            examTitle: '数据结构期末考试',
            subject: '计算机科学',
            score: 85,
            grade: '良好',
            rank: 2
          }
        ]
      }
    } else if (exportForm.type === 'report') {
      previewData.value = {
        title: '成绩分析报告',
        generateTime: new Date().toLocaleString(),
        summary: [
          { label: '参考人数', value: '45' },
          { label: '平均分', value: '82.3' },
          { label: '及格率', value: '91.1%' },
          { label: '优秀率', value: '33.3%' }
        ]
      }
    }
    
    ElMessage.success('预览生成成功')
  } catch (error) {
    ElMessage.error('预览生成失败')
  } finally {
    previewing.value = false
  }
}

// 开始导出
const startExport = async () => {
  if (!previewData.value) {
    ElMessage.warning('请先生成预览')
    return
  }
  
  exporting.value = true
  try {
    // 模拟导出过程
    await new Promise(resolve => setTimeout(resolve, 2000))
    
    if (exportForm.format === 'excel') {
      exportToExcel()
    } else if (exportForm.format === 'pdf') {
      exportToPDF()
    } else if (exportForm.format === 'csv') {
      exportToCSV()
    }
    
    // 添加到导出历史
    const fileName = `${exportForm.type}_${Date.now()}.${exportForm.format}`
    exportHistory.value.unshift({
      id: Date.now(),
      fileName,
      type: exportForm.type,
      format: exportForm.format,
      size: '128KB',
      createTime: new Date().toLocaleString(),
      status: 'completed'
    })
    
    ElMessage.success('导出成功')
  } catch (error) {
    ElMessage.error('导出失败，请重试')
  } finally {
    exporting.value = false
  }
}

// 导出到Excel
const exportToExcel = () => {
  const data = previewData.value.data || []
  const headers = selectedFields.value.map(field => field.label)
  const rows = data.map(item => 
    selectedFields.value.map(field => item[field.prop])
  )
  
  const worksheetData = [headers, ...rows]
  const ws = XLSX.utils.aoa_to_sheet(worksheetData)
  const wb = XLSX.utils.book_new()
  XLSX.utils.book_append_sheet(wb, ws, '成绩数据')
  
  const fileName = `${getTypeLabel(exportForm.type)}_${new Date().toISOString().slice(0, 10)}.xlsx`
  XLSX.writeFile(wb, fileName)
}

// 导出到PDF
const exportToPDF = () => {
  const doc = new jsPDF()
  
  // 设置字体（这里需要支持中文的字体）
  doc.setFontSize(16)
  doc.text(previewData.value.title || '成绩报告', 20, 20)
  
  doc.setFontSize(12)
  doc.text(`生成时间: ${new Date().toLocaleString()}`, 20, 35)
  
  // 这里应该添加表格和图表内容
  doc.text('成绩数据表格...', 20, 50)
  
  const fileName = `${getTypeLabel(exportForm.type)}_${new Date().toISOString().slice(0, 10)}.pdf`
  doc.save(fileName)
}

// 导出到CSV
const exportToCSV = () => {
  const data = previewData.value.data || []
  const headers = selectedFields.value.map(field => field.label)
  const rows = data.map(item => 
    selectedFields.value.map(field => item[field.prop])
  )
  
  const csvContent = [headers, ...rows]
    .map(row => row.join(','))
    .join('\n')
  
  const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' })
  const link = document.createElement('a')
  const url = URL.createObjectURL(blob)
  link.setAttribute('href', url)
  link.setAttribute('download', `${getTypeLabel(exportForm.type)}_${new Date().toISOString().slice(0, 10)}.csv`)
  link.style.visibility = 'hidden'
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}

// 获取类型标签
const getTypeLabel = (type) => {
  const labels = {
    grades: '成绩单',
    report: '成绩报告',
    statistics: '统计分析',
    transcript: '成绩单据'
  }
  return labels[type] || type
}

// 获取状态类型
const getStatusType = (status) => {
  const types = {
    completed: 'success',
    processing: 'warning',
    failed: 'danger'
  }
  return types[status] || ''
}

// 获取状态标签
const getStatusLabel = (status) => {
  const labels = {
    completed: '已完成',
    processing: '处理中',
    failed: '失败'
  }
  return labels[status] || status
}

// 下载文件
const downloadFile = (row) => {
  // 模拟下载
  ElMessage.success(`开始下载 ${row.fileName}`)
}

// 删除历史记录
const deleteHistory = async (row) => {
  try {
    await ElMessageBox.confirm('确定要删除这条导出记录吗？', '确认删除', {
      type: 'warning'
    })
    
    const index = exportHistory.value.findIndex(item => item.id === row.id)
    if (index !== -1) {
      exportHistory.value.splice(index, 1)
      ElMessage.success('删除成功')
    }
  } catch {
    // 用户取消删除
  }
}

// 清空历史
const clearHistory = async () => {
  try {
    await ElMessageBox.confirm('确定要清空所有导出历史吗？', '确认清空', {
      type: 'warning'
    })
    
    exportHistory.value = []
    ElMessage.success('清空成功')
  } catch {
    // 用户取消清空
  }
}

onMounted(() => {
  // 初始化时生成预览
  refreshPreview()
})
</script>

<style scoped>
.grade-export {
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

.dopamine-card h4 {
  margin: 0 0 var(--spacing-sm) 0;
  color: var(--text-primary);
  font-weight: 500;
  font-size: 14px;
}

.template-card {
  border: 2px solid var(--border-light);
  border-radius: var(--radius-md);
  padding: var(--spacing-md);
  cursor: pointer;
  transition: all 0.3s ease;
}

.template-card:hover {
  border-color: var(--color-primary);
  box-shadow: var(--shadow-md);
}

.template-card.active {
  border-color: var(--color-primary);
  background: var(--color-primary-light);
}

.template-preview {
  text-align: center;
  margin-bottom: var(--spacing-sm);
}

.template-preview img {
  width: 100%;
  height: 120px;
  object-fit: cover;
  border-radius: var(--radius-sm);
}

.template-info h4 {
  margin: 0 0 var(--spacing-xs) 0;
  font-size: 14px;
  font-weight: 600;
}

.template-info p {
  margin: 0;
  font-size: 12px;
  color: var(--text-secondary);
  line-height: 1.4;
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

.preview-content {
  border: 1px solid var(--border-light);
  border-radius: var(--radius-md);
  padding: var(--spacing-lg);
  background: var(--bg-light);
}

.preview-title {
  text-align: center;
  font-size: 18px;
  font-weight: 600;
  margin-bottom: var(--spacing-md);
  color: var(--text-primary);
}

.preview-note {
  text-align: center;
  margin-top: var(--spacing-md);
  color: var(--text-secondary);
  font-size: 14px;
}

.report-header {
  text-align: center;
  margin-bottom: var(--spacing-lg);
}

.report-header h2 {
  margin: 0 0 var(--spacing-sm) 0;
  color: var(--text-primary);
}

.report-header p {
  margin: 0;
  color: var(--text-secondary);
}

.report-summary {
  margin-bottom: var(--spacing-lg);
}

.summary-item {
  text-align: center;
  padding: var(--spacing-md);
  background: white;
  border-radius: var(--radius-md);
  border: 1px solid var(--border-light);
}

.summary-value {
  font-size: 24px;
  font-weight: 700;
  color: var(--color-primary);
  margin-bottom: var(--spacing-xs);
}

.summary-label {
  font-size: 14px;
  color: var(--text-secondary);
}

.history-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--spacing-md);
}

@media (max-width: 768px) {
  .filter-conditions .el-form {
    display: block;
  }
  
  .filter-conditions .el-form-item {
    display: block;
    margin-bottom: var(--spacing-md);
  }
  
  .template-card {
    margin-bottom: var(--spacing-md);
  }
  
  .preview-header {
    flex-direction: column;
    gap: var(--spacing-md);
  }
}
</style>