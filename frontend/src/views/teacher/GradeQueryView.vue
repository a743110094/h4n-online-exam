<template>
  <div class="grade-query">
    <!-- 页面标题 -->
    <div class="page-header">
      <h1 class="page-title">
        <el-icon><Search /></el-icon>
        成绩查询
      </h1>
      <p class="page-description">
        查询学生成绩、统计分析、成绩对比
      </p>
    </div>

    <!-- 查询条件 -->
    <div class="query-form dopamine-card">
      <h3>查询条件</h3>
      <el-form :model="queryForm" :inline="true" label-width="80px">
        <el-form-item label="考试名称">
          <el-select v-model="queryForm.examId" placeholder="请选择考试" clearable style="width: 200px">
            <el-option
              v-for="exam in exams"
              :key="exam.id"
              :label="exam.title"
              :value="exam.id"
            />
          </el-select>
        </el-form-item>
        
        <el-form-item label="科目">
          <el-select v-model="queryForm.subject" placeholder="请选择科目" clearable style="width: 150px">
            <el-option
              v-for="subject in subjects"
              :key="subject"
              :label="subject"
              :value="subject"
            />
          </el-select>
        </el-form-item>
        
        <el-form-item label="班级">
          <el-select v-model="queryForm.classId" placeholder="请选择班级" clearable style="width: 180px">
            <el-option
              v-for="cls in classes"
              :key="cls.id"
              :label="cls.name"
              :value="cls.id"
            />
          </el-select>
        </el-form-item>
        
        <el-form-item label="学号">
          <el-input
            v-model="queryForm.studentNumber"
            placeholder="请输入学号"
            clearable
            style="width: 150px"
          />
        </el-form-item>
        
        <el-form-item label="姓名">
          <el-input
            v-model="queryForm.studentName"
            placeholder="请输入姓名"
            clearable
            style="width: 150px"
          />
        </el-form-item>
        
        <el-form-item label="成绩范围">
          <el-input-number
            v-model="queryForm.minScore"
            :min="0"
            :max="100"
            placeholder="最低分"
            style="width: 100px"
          />
          <span style="margin: 0 8px">-</span>
          <el-input-number
            v-model="queryForm.maxScore"
            :min="0"
            :max="100"
            placeholder="最高分"
            style="width: 100px"
          />
        </el-form-item>
        
        <el-form-item>
          <el-button type="primary" @click="searchGrades" :loading="searching">
            <el-icon><Search /></el-icon>
            查询
          </el-button>
          <el-button @click="resetQuery">
            <el-icon><RefreshLeft /></el-icon>
            重置
          </el-button>
        </el-form-item>
      </el-form>
    </div>

    <!-- 统计信息 -->
    <div class="statistics dopamine-card" v-if="statistics">
      <h3>统计信息</h3>
      <el-row :gutter="20">
        <el-col :span="6">
          <div class="stat-item">
            <div class="stat-value">{{ statistics.totalCount }}</div>
            <div class="stat-label">总人数</div>
          </div>
        </el-col>
        <el-col :span="6">
          <div class="stat-item">
            <div class="stat-value">{{ statistics.averageScore }}</div>
            <div class="stat-label">平均分</div>
          </div>
        </el-col>
        <el-col :span="6">
          <div class="stat-item">
            <div class="stat-value">{{ statistics.passRate }}%</div>
            <div class="stat-label">及格率</div>
          </div>
        </el-col>
        <el-col :span="6">
          <div class="stat-item">
            <div class="stat-value">{{ statistics.excellentRate }}%</div>
            <div class="stat-label">优秀率</div>
          </div>
        </el-col>
      </el-row>
      
      <el-row :gutter="20" style="margin-top: 16px">
        <el-col :span="6">
          <div class="stat-item">
            <div class="stat-value">{{ statistics.highestScore }}</div>
            <div class="stat-label">最高分</div>
          </div>
        </el-col>
        <el-col :span="6">
          <div class="stat-item">
            <div class="stat-value">{{ statistics.lowestScore }}</div>
            <div class="stat-label">最低分</div>
          </div>
        </el-col>
        <el-col :span="6">
          <div class="stat-item">
            <div class="stat-value">{{ statistics.standardDeviation }}</div>
            <div class="stat-label">标准差</div>
          </div>
        </el-col>
        <el-col :span="6">
          <div class="stat-item">
            <div class="stat-value">{{ statistics.median }}</div>
            <div class="stat-label">中位数</div>
          </div>
        </el-col>
      </el-row>
    </div>

    <!-- 成绩分布图表 -->
    <div class="charts dopamine-card" v-if="gradeData.length > 0">
      <h3>成绩分布</h3>
      <el-row :gutter="20">
        <el-col :span="12">
          <div class="chart-container">
            <h4>分数段分布</h4>
            <div ref="scoreDistributionChart" class="chart"></div>
          </div>
        </el-col>
        <el-col :span="12">
          <div class="chart-container">
            <h4>等级分布</h4>
            <div ref="gradeDistributionChart" class="chart"></div>
          </div>
        </el-col>
      </el-row>
    </div>

    <!-- 成绩列表 -->
    <div class="grade-table dopamine-card">
      <div class="table-header">
        <h3>成绩列表</h3>
        <div class="table-actions">
          <el-button @click="exportGrades" :loading="exporting">
            <el-icon><Download /></el-icon>
            导出Excel
          </el-button>
          <el-button @click="printGrades">
            <el-icon><Printer /></el-icon>
            打印
          </el-button>
        </div>
      </div>
      
      <el-table
        :data="gradeData"
        stripe
        border
        v-loading="searching"
        :default-sort="{ prop: 'score', order: 'descending' }"
        @sort-change="handleSortChange"
      >
        <el-table-column type="index" label="序号" width="60" />
        <el-table-column prop="studentNumber" label="学号" min-width="120" sortable />
        <el-table-column prop="studentName" label="姓名" min-width="120" sortable />
        <el-table-column prop="className" label="班级" min-width="150" sortable />
        <el-table-column prop="examTitle" label="考试名称" min-width="180" sortable />
        <el-table-column prop="subject" label="科目" min-width="100" sortable />
        <el-table-column prop="score" label="成绩" min-width="100" sortable>
          <template #default="{ row }">
            <span :class="getScoreClass(row.score)">{{ row.score }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="grade" label="等级" min-width="100">
          <template #default="{ row }">
            <el-tag :type="getGradeType(row.score)">{{ getGradeLevel(row.score) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="rank" label="排名" min-width="80" sortable />
        <el-table-column prop="examDate" label="考试日期" min-width="120" sortable />
        <el-table-column prop="remark" label="备注" min-width="120" />
        <el-table-column label="操作" width="120" fixed="right">
          <template #default="{ row }">
            <el-button type="text" size="small" @click="viewDetail(row)">
              <el-icon><View /></el-icon>
              详情
            </el-button>
            <el-button type="text" size="small" @click="editGrade(row)">
              <el-icon><Edit /></el-icon>
              编辑
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <!-- 分页 -->
      <div class="pagination">
        <el-pagination
          v-model:current-page="pagination.currentPage"
          v-model:page-size="pagination.pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </div>

    <!-- 成绩详情对话框 -->
    <el-dialog v-model="detailVisible" title="成绩详情" width="600px">
      <div class="grade-detail" v-if="selectedGrade">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="学号">{{ selectedGrade.studentNumber }}</el-descriptions-item>
          <el-descriptions-item label="姓名">{{ selectedGrade.studentName }}</el-descriptions-item>
          <el-descriptions-item label="班级">{{ selectedGrade.className }}</el-descriptions-item>
          <el-descriptions-item label="考试名称">{{ selectedGrade.examTitle }}</el-descriptions-item>
          <el-descriptions-item label="科目">{{ selectedGrade.subject }}</el-descriptions-item>
          <el-descriptions-item label="成绩">
            <span :class="getScoreClass(selectedGrade.score)">{{ selectedGrade.score }}</span>
          </el-descriptions-item>
          <el-descriptions-item label="等级">
            <el-tag :type="getGradeType(selectedGrade.score)">{{ getGradeLevel(selectedGrade.score) }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="排名">{{ selectedGrade.rank }}</el-descriptions-item>
          <el-descriptions-item label="考试日期">{{ selectedGrade.examDate }}</el-descriptions-item>
          <el-descriptions-item label="录入时间">{{ selectedGrade.createTime }}</el-descriptions-item>
          <el-descriptions-item label="备注" :span="2">{{ selectedGrade.remark || '无' }}</el-descriptions-item>
        </el-descriptions>
        
        <!-- 答题详情 -->
        <div class="answer-details" v-if="selectedGrade.answers">
          <h4>答题详情</h4>
          <el-table :data="selectedGrade.answers" border>
            <el-table-column prop="questionNumber" label="题号" width="80" />
            <el-table-column prop="questionType" label="题型" width="100" />
            <el-table-column prop="score" label="得分" width="80" />
            <el-table-column prop="totalScore" label="总分" width="80" />
            <el-table-column prop="isCorrect" label="正确性" width="100">
              <template #default="{ row }">
                <el-tag :type="row.isCorrect ? 'success' : 'danger'">
                  {{ row.isCorrect ? '正确' : '错误' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="answer" label="学生答案" min-width="200" />
          </el-table>
        </div>
      </div>
    </el-dialog>

    <!-- 编辑成绩对话框 -->
    <el-dialog v-model="editVisible" title="编辑成绩" width="500px">
      <el-form :model="editForm" label-width="80px">
        <el-form-item label="学生">
          <span>{{ editForm.studentName }} ({{ editForm.studentNumber }})</span>
        </el-form-item>
        <el-form-item label="考试">
          <span>{{ editForm.examTitle }}</span>
        </el-form-item>
        <el-form-item label="成绩" required>
          <el-input-number
            v-model="editForm.score"
            :min="0"
            :max="100"
            :precision="1"
            style="width: 200px"
          />
        </el-form-item>
        <el-form-item label="备注">
          <el-input
            v-model="editForm.remark"
            type="textarea"
            :rows="3"
            placeholder="请输入备注"
            maxlength="200"
          />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <el-button @click="editVisible = false">取消</el-button>
        <el-button type="primary" @click="saveEdit" :loading="saving">
          保存
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, nextTick } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Search, RefreshLeft, Download, Printer, View, Edit
} from '@element-plus/icons-vue'
import * as echarts from 'echarts'
import * as XLSX from 'xlsx'

// 查询表单
const queryForm = reactive({
  examId: '',
  subject: '',
  classId: '',
  studentNumber: '',
  studentName: '',
  minScore: null,
  maxScore: null
})

// 状态
const searching = ref(false)
const exporting = ref(false)
const saving = ref(false)
const detailVisible = ref(false)
const editVisible = ref(false)

// 数据
const gradeData = ref([])
const selectedGrade = ref(null)
const statistics = ref(null)

// 编辑表单
const editForm = reactive({
  id: '',
  studentNumber: '',
  studentName: '',
  examTitle: '',
  score: 0,
  remark: ''
})

// 分页
const pagination = reactive({
  currentPage: 1,
  pageSize: 20,
  total: 0
})

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

// 图表引用
const scoreDistributionChart = ref()
const gradeDistributionChart = ref()

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

// 获取成绩样式类
const getScoreClass = (score) => {
  if (score >= 90) return 'score-excellent'
  if (score >= 80) return 'score-good'
  if (score >= 70) return 'score-medium'
  if (score >= 60) return 'score-pass'
  return 'score-fail'
}

// 查询成绩
const searchGrades = async () => {
  searching.value = true
  try {
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    // 模拟数据
    const mockData = [
      {
        id: 1,
        studentNumber: '2021001',
        studentName: '张三',
        className: '计算机2021-1班',
        examTitle: '数据结构期末考试',
        subject: '计算机科学',
        score: 92,
        rank: 1,
        examDate: '2024-01-15',
        createTime: '2024-01-16 10:30:00',
        remark: '优秀',
        answers: [
          { questionNumber: 1, questionType: '选择题', score: 8, totalScore: 10, isCorrect: false, answer: 'A' },
          { questionNumber: 2, questionType: '填空题', score: 10, totalScore: 10, isCorrect: true, answer: '栈' },
          { questionNumber: 3, questionType: '简答题', score: 18, totalScore: 20, isCorrect: true, answer: '栈是后进先出的数据结构...' }
        ]
      },
      {
        id: 2,
        studentNumber: '2021002',
        studentName: '李四',
        className: '计算机2021-1班',
        examTitle: '数据结构期末考试',
        subject: '计算机科学',
        score: 85,
        rank: 2,
        examDate: '2024-01-15',
        createTime: '2024-01-16 10:35:00',
        remark: '良好'
      },
      {
        id: 3,
        studentNumber: '2021003',
        studentName: '王五',
        className: '计算机2021-1班',
        examTitle: '数据结构期末考试',
        subject: '计算机科学',
        score: 78,
        rank: 3,
        examDate: '2024-01-15',
        createTime: '2024-01-16 10:40:00',
        remark: ''
      }
    ]
    
    gradeData.value = mockData
    pagination.total = mockData.length
    
    // 计算统计信息
    calculateStatistics(mockData)
    
    // 绘制图表
    nextTick(() => {
      drawCharts(mockData)
    })
    
    ElMessage.success('查询成功')
  } catch (error) {
    ElMessage.error('查询失败，请重试')
  } finally {
    searching.value = false
  }
}

// 重置查询
const resetQuery = () => {
  Object.keys(queryForm).forEach(key => {
    queryForm[key] = key.includes('Score') ? null : ''
  })
  gradeData.value = []
  statistics.value = null
}

// 计算统计信息
const calculateStatistics = (data) => {
  if (data.length === 0) {
    statistics.value = null
    return
  }
  
  const scores = data.map(item => item.score)
  const totalCount = data.length
  const averageScore = (scores.reduce((sum, score) => sum + score, 0) / totalCount).toFixed(1)
  const passCount = scores.filter(score => score >= 60).length
  const excellentCount = scores.filter(score => score >= 90).length
  const passRate = ((passCount / totalCount) * 100).toFixed(1)
  const excellentRate = ((excellentCount / totalCount) * 100).toFixed(1)
  const highestScore = Math.max(...scores)
  const lowestScore = Math.min(...scores)
  
  // 计算标准差
  const variance = scores.reduce((sum, score) => sum + Math.pow(score - parseFloat(averageScore), 2), 0) / totalCount
  const standardDeviation = Math.sqrt(variance).toFixed(2)
  
  // 计算中位数
  const sortedScores = [...scores].sort((a, b) => a - b)
  const median = totalCount % 2 === 0
    ? ((sortedScores[totalCount / 2 - 1] + sortedScores[totalCount / 2]) / 2).toFixed(1)
    : sortedScores[Math.floor(totalCount / 2)].toFixed(1)
  
  statistics.value = {
    totalCount,
    averageScore,
    passRate,
    excellentRate,
    highestScore,
    lowestScore,
    standardDeviation,
    median
  }
}

// 绘制图表
const drawCharts = (data) => {
  if (data.length === 0) return
  
  // 分数段分布图
  const scoreRanges = ['0-59', '60-69', '70-79', '80-89', '90-100']
  const scoreDistribution = scoreRanges.map(range => {
    const [min, max] = range.split('-').map(Number)
    return data.filter(item => item.score >= min && item.score <= max).length
  })
  
  const scoreChart = echarts.init(scoreDistributionChart.value)
  scoreChart.setOption({
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'shadow'
      }
    },
    xAxis: {
      type: 'category',
      data: scoreRanges
    },
    yAxis: {
      type: 'value'
    },
    series: [{
      data: scoreDistribution,
      type: 'bar',
      itemStyle: {
        color: '#409EFF'
      }
    }]
  })
  
  // 等级分布饼图
  const gradeDistribution = [
    { value: data.filter(item => item.score >= 90).length, name: '优秀' },
    { value: data.filter(item => item.score >= 80 && item.score < 90).length, name: '良好' },
    { value: data.filter(item => item.score >= 70 && item.score < 80).length, name: '中等' },
    { value: data.filter(item => item.score >= 60 && item.score < 70).length, name: '及格' },
    { value: data.filter(item => item.score < 60).length, name: '不及格' }
  ]
  
  const gradeChart = echarts.init(gradeDistributionChart.value)
  gradeChart.setOption({
    tooltip: {
      trigger: 'item'
    },
    legend: {
      orient: 'vertical',
      left: 'left'
    },
    series: [{
      name: '等级分布',
      type: 'pie',
      radius: '50%',
      data: gradeDistribution,
      emphasis: {
        itemStyle: {
          shadowBlur: 10,
          shadowOffsetX: 0,
          shadowColor: 'rgba(0, 0, 0, 0.5)'
        }
      }
    }]
  })
}

// 查看详情
const viewDetail = (row) => {
  selectedGrade.value = row
  detailVisible.value = true
}

// 编辑成绩
const editGrade = (row) => {
  Object.assign(editForm, {
    id: row.id,
    studentNumber: row.studentNumber,
    studentName: row.studentName,
    examTitle: row.examTitle,
    score: row.score,
    remark: row.remark
  })
  editVisible.value = true
}

// 保存编辑
const saveEdit = async () => {
  saving.value = true
  try {
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 500))
    
    // 更新本地数据
    const index = gradeData.value.findIndex(item => item.id === editForm.id)
    if (index !== -1) {
      gradeData.value[index].score = editForm.score
      gradeData.value[index].remark = editForm.remark
    }
    
    editVisible.value = false
    ElMessage.success('修改成功')
  } catch (error) {
    ElMessage.error('修改失败，请重试')
  } finally {
    saving.value = false
  }
}

// 导出成绩
const exportGrades = () => {
  if (gradeData.value.length === 0) {
    ElMessage.warning('没有数据可以导出')
    return
  }
  
  exporting.value = true
  
  try {
    const exportData = [
      ['学号', '姓名', '班级', '考试名称', '科目', '成绩', '等级', '排名', '考试日期', '备注']
    ]
    
    gradeData.value.forEach(item => {
      exportData.push([
        item.studentNumber,
        item.studentName,
        item.className,
        item.examTitle,
        item.subject,
        item.score,
        getGradeLevel(item.score),
        item.rank,
        item.examDate,
        item.remark || ''
      ])
    })
    
    const ws = XLSX.utils.aoa_to_sheet(exportData)
    const wb = XLSX.utils.book_new()
    XLSX.utils.book_append_sheet(wb, ws, '成绩查询结果')
    XLSX.writeFile(wb, `成绩查询结果_${new Date().toISOString().slice(0, 10)}.xlsx`)
    
    ElMessage.success('导出成功')
  } catch (error) {
    ElMessage.error('导出失败，请重试')
  } finally {
    exporting.value = false
  }
}

// 打印成绩
const printGrades = () => {
  if (gradeData.value.length === 0) {
    ElMessage.warning('没有数据可以打印')
    return
  }
  
  window.print()
}

// 排序变化
const handleSortChange = ({ column, prop, order }) => {
  // 实际项目中这里应该调用API重新获取排序后的数据
  console.log('排序变化:', { column, prop, order })
}

// 分页大小变化
const handleSizeChange = (size) => {
  pagination.pageSize = size
  searchGrades()
}

// 当前页变化
const handleCurrentChange = (page) => {
  pagination.currentPage = page
  searchGrades()
}

onMounted(() => {
  // 初始化时加载数据
  searchGrades()
})
</script>

<style scoped>
.grade-query {
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

.statistics {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.statistics h3 {
  color: white;
}

.stat-item {
  text-align: center;
  padding: var(--spacing-md);
  background: rgba(255, 255, 255, 0.1);
  border-radius: var(--radius-md);
}

.stat-value {
  font-size: 28px;
  font-weight: 700;
  margin-bottom: var(--spacing-xs);
}

.stat-label {
  font-size: 14px;
  opacity: 0.9;
}

.chart-container {
  text-align: center;
}

.chart-container h4 {
  margin: 0 0 var(--spacing-md) 0;
  color: var(--text-primary);
}

.chart {
  height: 300px;
  width: 100%;
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

.pagination {
  display: flex;
  justify-content: center;
  margin-top: var(--spacing-lg);
}

.score-excellent {
  color: #67c23a;
  font-weight: 600;
}

.score-good {
  color: #409eff;
  font-weight: 600;
}

.score-medium {
  color: #e6a23c;
  font-weight: 600;
}

.score-pass {
  color: #909399;
  font-weight: 600;
}

.score-fail {
  color: #f56c6c;
  font-weight: 600;
}

.grade-detail {
  padding: var(--spacing-md);
}

.answer-details {
  margin-top: var(--spacing-lg);
}

.answer-details h4 {
  margin: 0 0 var(--spacing-md) 0;
  color: var(--text-primary);
}

@media (max-width: 768px) {
  .query-form .el-form {
    display: block;
  }
  
  .query-form .el-form-item {
    display: block;
    margin-bottom: var(--spacing-md);
  }
  
  .charts .el-row {
    flex-direction: column;
  }
  
  .chart {
    height: 250px;
  }
}
</style>