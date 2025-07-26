<template>
  <div class="grade-view">
    <!-- 页面标题 -->
    <div class="page-header">
      <h1 class="page-title">
        <el-icon><Medal /></el-icon>
        成绩查询
      </h1>
      <p class="page-subtitle">查看我的考试成绩和学习进度</p>
    </div>

    <!-- 筛选条件 -->
    <div class="filter-section dopamine-card">
      <div class="filter-row">
        <div class="filter-item">
          <label>科目</label>
          <el-select v-model="filters.subject" placeholder="选择科目" clearable>
            <el-option label="全部" value="" />
            <el-option label="数学" value="math" />
            <el-option label="语文" value="chinese" />
            <el-option label="英语" value="english" />
            <el-option label="物理" value="physics" />
            <el-option label="化学" value="chemistry" />
          </el-select>
        </div>
        <div class="filter-item">
          <label>考试类型</label>
          <el-select v-model="filters.examType" placeholder="选择考试类型" clearable>
            <el-option label="全部" value="" />
            <el-option label="期末考试" value="final" />
            <el-option label="期中考试" value="midterm" />
            <el-option label="随堂测验" value="quiz" />
            <el-option label="模拟考试" value="mock" />
          </el-select>
        </div>
        <div class="filter-item">
          <label>时间范围</label>
          <el-date-picker
            v-model="filters.dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            format="YYYY-MM-DD"
            value-format="YYYY-MM-DD"
          />
        </div>
        <el-button type="primary" :icon="Search" @click="searchGrades">
          查询
        </el-button>
        <el-button :icon="Refresh" @click="resetFilters">
          重置
        </el-button>
      </div>
    </div>

    <!-- 成绩统计 -->
    <div class="stats-section">
      <div class="stats-grid">
        <div class="stat-card dopamine-card">
          <div class="stat-icon" style="background: var(--gradient-primary)">
            <el-icon :size="24"><TrendCharts /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ gradeStats.average }}</div>
            <div class="stat-label">平均分</div>
          </div>
        </div>
        <div class="stat-card dopamine-card">
          <div class="stat-icon" style="background: var(--gradient-success)">
            <el-icon :size="24"><Trophy /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ gradeStats.highest }}</div>
            <div class="stat-label">最高分</div>
          </div>
        </div>
        <div class="stat-card dopamine-card">
          <div class="stat-icon" style="background: var(--gradient-warning)">
            <el-icon :size="24"><Notebook /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ gradeStats.totalExams }}</div>
            <div class="stat-label">考试次数</div>
          </div>
        </div>
        <div class="stat-card dopamine-card">
          <div class="stat-icon" style="background: var(--gradient-info)">
            <el-icon :size="24"><Rank /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ gradeStats.rank }}</div>
            <div class="stat-label">班级排名</div>
          </div>
        </div>
      </div>
    </div>

    <!-- 成绩列表 -->
    <div class="grade-list dopamine-card">
      <div class="list-header">
        <h3>成绩记录</h3>
        <el-button :icon="Download" @click="exportGrades">
          导出成绩
        </el-button>
      </div>
      
      <el-table 
        :data="filteredGrades" 
        stripe 
        v-loading="loading"
        @row-click="viewGradeDetail"
        style="cursor: pointer"
      >
        <el-table-column prop="examName" label="考试名称" min-width="150" />
        <el-table-column prop="subject" label="科目" width="100" />
        <el-table-column prop="examType" label="考试类型" width="120">
          <template #default="{ row }">
            <el-tag :type="getExamTypeColor(row.examType)">{{ getExamTypeName(row.examType) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="score" label="得分" width="100">
          <template #default="{ row }">
            <span :class="getScoreClass(row.score, row.totalScore)">{{ row.score }}/{{ row.totalScore }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="percentage" label="得分率" width="100">
          <template #default="{ row }">
            <span :class="getScoreClass(row.score, row.totalScore)">{{ ((row.score / row.totalScore) * 100).toFixed(1) }}%</span>
          </template>
        </el-table-column>
        <el-table-column prop="classRank" label="班级排名" width="100" />
        <el-table-column prop="gradeRank" label="年级排名" width="100" />
        <el-table-column prop="examDate" label="考试时间" width="120" />
        <el-table-column label="操作" width="120" fixed="right">
          <template #default="{ row }">
            <el-button type="text" size="small" @click.stop="viewGradeDetail(row)">
              查看详情
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="totalGrades"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </div>

    <!-- 成绩详情对话框 -->
    <el-dialog
      v-model="showDetailDialog"
      title="成绩详情"
      width="800px"
      :before-close="closeDetailDialog"
    >
      <div v-if="selectedGrade" class="grade-detail">
        <div class="detail-header">
          <h3>{{ selectedGrade.examName }}</h3>
          <el-tag :type="getExamTypeColor(selectedGrade.examType)">{{ getExamTypeName(selectedGrade.examType) }}</el-tag>
        </div>
        
        <div class="detail-content">
          <el-row :gutter="20">
            <el-col :span="12">
              <div class="detail-item">
                <label>科目：</label>
                <span>{{ selectedGrade.subject }}</span>
              </div>
              <div class="detail-item">
                <label>得分：</label>
                <span class="score-highlight">{{ selectedGrade.score }}/{{ selectedGrade.totalScore }}</span>
              </div>
              <div class="detail-item">
                <label>得分率：</label>
                <span>{{ ((selectedGrade.score / selectedGrade.totalScore) * 100).toFixed(1) }}%</span>
              </div>
            </el-col>
            <el-col :span="12">
              <div class="detail-item">
                <label>班级排名：</label>
                <span>{{ selectedGrade.classRank }}</span>
              </div>
              <div class="detail-item">
                <label>年级排名：</label>
                <span>{{ selectedGrade.gradeRank }}</span>
              </div>
              <div class="detail-item">
                <label>考试时间：</label>
                <span>{{ selectedGrade.examDate }}</span>
              </div>
            </el-col>
          </el-row>
          
          <div class="detail-item full-width" v-if="selectedGrade.remark">
            <label>备注：</label>
            <span>{{ selectedGrade.remark }}</span>
          </div>
        </div>
      </div>
      
      <template #footer>
        <el-button @click="closeDetailDialog">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import {
  Search,
  Refresh,
  Download,
  Medal,
  TrendCharts,
  Trophy,
  Notebook,
  Rank
} from '@element-plus/icons-vue'

// 响应式数据
const loading = ref(false)
const currentPage = ref(1)
const pageSize = ref(20)
const totalGrades = ref(0)
const showDetailDialog = ref(false)
const selectedGrade = ref(null)

// 筛选条件
const filters = reactive({
  subject: '',
  examType: '',
  dateRange: []
})

// 成绩统计
const gradeStats = ref({
  average: '85.6',
  highest: '98',
  totalExams: '12',
  rank: '5'
})

// 成绩数据
const grades = ref([
  {
    id: 1,
    examName: '数学期末考试',
    subject: '数学',
    examType: 'final',
    score: 95,
    totalScore: 100,
    classRank: 3,
    gradeRank: 15,
    examDate: '2024-01-15',
    remark: '表现优秀'
  },
  {
    id: 2,
    examName: '英语期中考试',
    subject: '英语',
    examType: 'midterm',
    score: 88,
    totalScore: 100,
    classRank: 8,
    gradeRank: 45,
    examDate: '2024-01-10',
    remark: ''
  },
  {
    id: 3,
    examName: '物理随堂测验',
    subject: '物理',
    examType: 'quiz',
    score: 76,
    totalScore: 100,
    classRank: 12,
    gradeRank: 68,
    examDate: '2024-01-08',
    remark: '需要加强练习'
  }
])

// 计算属性
const filteredGrades = computed(() => {
  let result = grades.value
  
  if (filters.subject) {
    result = result.filter(grade => grade.subject === filters.subject)
  }
  
  if (filters.examType) {
    result = result.filter(grade => grade.examType === filters.examType)
  }
  
  if (filters.dateRange && filters.dateRange.length === 2) {
    const [startDate, endDate] = filters.dateRange
    result = result.filter(grade => {
      return grade.examDate >= startDate && grade.examDate <= endDate
    })
  }
  
  return result
})

// 方法
const searchGrades = () => {
  loading.value = true
  // 模拟API调用
  setTimeout(() => {
    loading.value = false
    ElMessage.success('查询完成')
  }, 1000)
}

const resetFilters = () => {
  filters.subject = ''
  filters.examType = ''
  filters.dateRange = []
  searchGrades()
}

const exportGrades = () => {
  ElMessage.info('导出功能开发中')
}

const viewGradeDetail = (row: any) => {
  selectedGrade.value = row
  showDetailDialog.value = true
}

const closeDetailDialog = () => {
  showDetailDialog.value = false
  selectedGrade.value = null
}

const handleSizeChange = (size: number) => {
  pageSize.value = size
  searchGrades()
}

const handleCurrentChange = (page: number) => {
  currentPage.value = page
  searchGrades()
}

const getExamTypeColor = (type: string) => {
  const colorMap: Record<string, string> = {
    final: 'danger',
    midterm: 'warning',
    quiz: 'info',
    mock: 'success'
  }
  return colorMap[type] || ''
}

const getExamTypeName = (type: string) => {
  const nameMap: Record<string, string> = {
    final: '期末考试',
    midterm: '期中考试',
    quiz: '随堂测验',
    mock: '模拟考试'
  }
  return nameMap[type] || type
}

const getScoreClass = (score: number, totalScore: number) => {
  const percentage = (score / totalScore) * 100
  if (percentage >= 90) return 'score-excellent'
  if (percentage >= 80) return 'score-good'
  if (percentage >= 70) return 'score-average'
  return 'score-poor'
}

onMounted(() => {
  searchGrades()
})
</script>

<style scoped>
.grade-view {
  padding: var(--spacing-lg);
}

.page-header {
  margin-bottom: var(--spacing-xl);
}

.page-title {
  font-size: 28px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 var(--spacing-sm) 0;
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
}

.page-subtitle {
  color: var(--text-secondary);
  margin: 0;
}

.filter-section {
  padding: var(--spacing-lg);
  margin-bottom: var(--spacing-lg);
}

.filter-row {
  display: flex;
  align-items: center;
  gap: var(--spacing-lg);
  flex-wrap: wrap;
}

.filter-item {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
}

.filter-item label {
  font-weight: 500;
  color: var(--text-primary);
  white-space: nowrap;
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
  display: flex;
  align-items: center;
  padding: var(--spacing-lg);
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
  flex-shrink: 0;
}

.stat-content {
  flex: 1;
}

.stat-value {
  font-size: 24px;
  font-weight: 700;
  color: var(--text-primary);
  line-height: 1;
  margin-bottom: var(--spacing-xs);
}

.stat-label {
  font-size: 14px;
  color: var(--text-secondary);
}

.grade-list {
  padding: var(--spacing-lg);
}

.list-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--spacing-lg);
}

.list-header h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
}

.pagination-wrapper {
  display: flex;
  justify-content: center;
  margin-top: var(--spacing-lg);
}

.score-excellent {
  color: var(--color-success);
  font-weight: 600;
}

.score-good {
  color: var(--color-primary);
  font-weight: 600;
}

.score-average {
  color: var(--color-warning);
  font-weight: 600;
}

.score-poor {
  color: var(--color-danger);
  font-weight: 600;
}

.grade-detail {
  padding: var(--spacing-md);
}

.detail-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--spacing-lg);
  padding-bottom: var(--spacing-md);
  border-bottom: 1px solid var(--border-light);
}

.detail-header h3 {
  margin: 0;
  color: var(--text-primary);
}

.detail-content {
  margin-top: var(--spacing-lg);
}

.detail-item {
  display: flex;
  align-items: center;
  margin-bottom: var(--spacing-md);
}

.detail-item.full-width {
  flex-direction: column;
  align-items: flex-start;
}

.detail-item label {
  font-weight: 500;
  color: var(--text-secondary);
  width: 80px;
  flex-shrink: 0;
}

.detail-item.full-width label {
  width: auto;
  margin-bottom: var(--spacing-xs);
}

.score-highlight {
  font-size: 18px;
  font-weight: 600;
  color: var(--color-primary);
}
</style>