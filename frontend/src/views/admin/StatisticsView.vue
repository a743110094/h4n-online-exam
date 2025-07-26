<template>
  <div class="statistics-view">
    <!-- 页面标题 -->
    <div class="page-header">
      <h1 class="page-title">数据统计</h1>
      <p class="page-subtitle">系统数据分析与统计报表</p>
    </div>

    <!-- 时间筛选 -->
    <div class="filter-section dopamine-card">
      <div class="filter-row">
        <div class="filter-item">
          <label>统计时间</label>
          <el-date-picker
            v-model="dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            format="YYYY-MM-DD"
            value-format="YYYY-MM-DD"
            @change="updateStatistics"
          />
        </div>
        <div class="filter-item">
          <label>统计维度</label>
          <el-select v-model="dimension" @change="updateStatistics">
            <el-option label="按天" value="day" />
            <el-option label="按周" value="week" />
            <el-option label="按月" value="month" />
          </el-select>
        </div>
        <el-button type="primary" :icon="Refresh" @click="refreshData">
          刷新数据
        </el-button>
        <el-button :icon="Download" @click="exportReport">
          导出报表
        </el-button>
      </div>
    </div>

    <!-- 概览统计卡片 -->
    <div class="overview-stats">
      <div class="stat-card dopamine-card" v-for="stat in overviewStats" :key="stat.key">
        <div class="stat-icon" :style="{ background: stat.color }">
          <el-icon :size="24">
            <component :is="stat.icon" />
          </el-icon>
        </div>
        <div class="stat-content">
          <div class="stat-value">{{ stat.value }}</div>
          <div class="stat-label">{{ stat.label }}</div>
          <div class="stat-change" :class="stat.trend">
            <el-icon><component :is="stat.trendIcon" /></el-icon>
            {{ stat.change }}
          </div>
        </div>
      </div>
    </div>

    <!-- 图表区域 -->
    <div class="charts-section">
      <!-- 用户增长趋势 -->
      <div class="chart-card dopamine-card">
        <div class="chart-header">
          <h3>用户增长趋势</h3>
          <el-radio-group v-model="userGrowthType" @change="updateUserGrowthChart">
            <el-radio-button label="total">总用户</el-radio-button>
            <el-radio-button label="teacher">教师</el-radio-button>
            <el-radio-button label="student">学生</el-radio-button>
          </el-radio-group>
        </div>
        <div class="chart-container" ref="userGrowthChart"></div>
      </div>

      <!-- 考试统计 -->
      <div class="chart-card dopamine-card">
        <div class="chart-header">
          <h3>考试统计</h3>
        </div>
        <div class="chart-container" ref="examStatsChart"></div>
      </div>
    </div>

    <!-- 详细数据表格 -->
    <div class="data-tables">
      <!-- 热门科目排行 -->
      <div class="table-card dopamine-card">
        <div class="table-header">
          <h3>热门科目排行</h3>
        </div>
        <el-table :data="subjectRanking" stripe>
          <el-table-column prop="rank" label="排名" width="80" />
          <el-table-column prop="subject" label="科目" />
          <el-table-column prop="examCount" label="考试次数" />
          <el-table-column prop="participantCount" label="参与人数" />
          <el-table-column prop="avgScore" label="平均分" />
          <el-table-column label="趋势" width="100">
            <template #default="{ row }">
              <el-tag :type="row.trend === 'up' ? 'success' : row.trend === 'down' ? 'danger' : ''">
                <el-icon><component :is="row.trend === 'up' ? 'ArrowUp' : row.trend === 'down' ? 'ArrowDown' : 'Minus'" /></el-icon>
                {{ row.trendText }}
              </el-tag>
            </template>
          </el-table-column>
        </el-table>
      </div>

      <!-- 活跃用户统计 -->
      <div class="table-card dopamine-card">
        <div class="table-header">
          <h3>活跃用户统计</h3>
        </div>
        <el-table :data="activeUsers" stripe>
          <el-table-column prop="date" label="日期" />
          <el-table-column prop="totalActive" label="总活跃用户" />
          <el-table-column prop="teacherActive" label="活跃教师" />
          <el-table-column prop="studentActive" label="活跃学生" />
          <el-table-column prop="newUsers" label="新增用户" />
          <el-table-column prop="loginRate" label="登录率" />
        </el-table>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import {
  Refresh,
  Download,
  User,
  UserFilled,
  Notebook,
  TrendCharts,
  ArrowUp,
  ArrowDown,
  Minus
} from '@element-plus/icons-vue'
import * as echarts from 'echarts'

// 响应式数据
const dateRange = ref([])
const dimension = ref('day')
const userGrowthType = ref('total')
const userGrowthChart = ref()
const examStatsChart = ref()

// 概览统计数据
const overviewStats = ref([
  {
    key: 'totalUsers',
    label: '总用户数',
    value: '1,234',
    change: '+12.5%',
    trend: 'up',
    trendIcon: 'ArrowUp',
    icon: 'UserFilled',
    color: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)'
  },
  {
    key: 'totalExams',
    label: '总考试数',
    value: '456',
    change: '+8.3%',
    trend: 'up',
    trendIcon: 'ArrowUp',
    icon: 'Notebook',
    color: 'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)'
  },
  {
    key: 'activeUsers',
    label: '活跃用户',
    value: '892',
    change: '-2.1%',
    trend: 'down',
    trendIcon: 'ArrowDown',
    icon: 'User',
    color: 'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)'
  },
  {
    key: 'avgScore',
    label: '平均分',
    value: '78.5',
    change: '+3.2%',
    trend: 'up',
    trendIcon: 'ArrowUp',
    icon: 'TrendCharts',
    color: 'linear-gradient(135deg, #43e97b 0%, #38f9d7 100%)'
  }
])

// 科目排行数据
const subjectRanking = ref([
  {
    rank: 1,
    subject: '数学',
    examCount: 45,
    participantCount: 1200,
    avgScore: 82.5,
    trend: 'up',
    trendText: '上升'
  },
  {
    rank: 2,
    subject: '英语',
    examCount: 38,
    participantCount: 1150,
    avgScore: 78.3,
    trend: 'up',
    trendText: '上升'
  },
  {
    rank: 3,
    subject: '语文',
    examCount: 32,
    participantCount: 1100,
    avgScore: 75.8,
    trend: 'down',
    trendText: '下降'
  },
  {
    rank: 4,
    subject: '物理',
    examCount: 28,
    participantCount: 800,
    avgScore: 73.2,
    trend: 'stable',
    trendText: '持平'
  },
  {
    rank: 5,
    subject: '化学',
    examCount: 25,
    participantCount: 750,
    avgScore: 71.5,
    trend: 'up',
    trendText: '上升'
  }
])

// 活跃用户数据
const activeUsers = ref([
  {
    date: '2024-01-15',
    totalActive: 892,
    teacherActive: 45,
    studentActive: 847,
    newUsers: 12,
    loginRate: '72.3%'
  },
  {
    date: '2024-01-14',
    totalActive: 856,
    teacherActive: 42,
    studentActive: 814,
    newUsers: 8,
    loginRate: '69.5%'
  },
  {
    date: '2024-01-13',
    totalActive: 923,
    teacherActive: 48,
    studentActive: 875,
    newUsers: 15,
    loginRate: '74.8%'
  }
])

// 方法
const updateStatistics = () => {
  ElMessage.success('统计数据已更新')
}

const refreshData = () => {
  ElMessage.success('数据已刷新')
}

const exportReport = () => {
  ElMessage.success('报表导出功能开发中')
}

const updateUserGrowthChart = () => {
  initUserGrowthChart()
}

const initUserGrowthChart = () => {
  if (!userGrowthChart.value) return
  
  const chart = echarts.init(userGrowthChart.value)
  const option = {
    title: {
      text: '用户增长趋势',
      left: 'center',
      textStyle: {
        fontSize: 16,
        fontWeight: 'normal'
      }
    },
    tooltip: {
      trigger: 'axis'
    },
    legend: {
      data: ['新增用户', '累计用户'],
      bottom: 10
    },
    xAxis: {
      type: 'category',
      data: ['1月', '2月', '3月', '4月', '5月', '6月']
    },
    yAxis: {
      type: 'value'
    },
    series: [
      {
        name: '新增用户',
        type: 'bar',
        data: [120, 132, 101, 134, 90, 230],
        itemStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: '#667eea' },
            { offset: 1, color: '#764ba2' }
          ])
        }
      },
      {
        name: '累计用户',
        type: 'line',
        data: [220, 352, 453, 587, 677, 907],
        itemStyle: {
          color: '#f093fb'
        }
      }
    ]
  }
  chart.setOption(option)
}

const initExamStatsChart = () => {
  if (!examStatsChart.value) return
  
  const chart = echarts.init(examStatsChart.value)
  const option = {
    title: {
      text: '考试统计',
      left: 'center',
      textStyle: {
        fontSize: 16,
        fontWeight: 'normal'
      }
    },
    tooltip: {
      trigger: 'item'
    },
    legend: {
      orient: 'vertical',
      left: 'left'
    },
    series: [
      {
        name: '考试类型',
        type: 'pie',
        radius: '50%',
        data: [
          { value: 1048, name: '期末考试' },
          { value: 735, name: '期中考试' },
          { value: 580, name: '随堂测验' },
          { value: 484, name: '模拟考试' },
          { value: 300, name: '补考' }
        ],
        emphasis: {
          itemStyle: {
            shadowBlur: 10,
            shadowOffsetX: 0,
            shadowColor: 'rgba(0, 0, 0, 0.5)'
          }
        }
      }
    ]
  }
  chart.setOption(option)
}

onMounted(() => {
  nextTick(() => {
    initUserGrowthChart()
    initExamStatsChart()
  })
})
</script>

<style scoped>
.statistics-view {
  padding: var(--spacing-md);
}

.page-header {
  margin-bottom: var(--spacing-lg);
}

.page-title {
  font-size: 24px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 var(--spacing-xs) 0;
}

.page-subtitle {
  color: var(--text-secondary);
  margin: 0;
  font-size: 14px;
}

.filter-section {
  padding: var(--spacing-md);
  margin-bottom: var(--spacing-md);
}

.filter-row {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
  flex-wrap: wrap;
}

.filter-item {
  display: flex;
  align-items: center;
  gap: var(--spacing-xs);
}

.filter-item label {
  font-weight: 500;
  color: var(--text-primary);
  white-space: nowrap;
  font-size: 13px;
}

.overview-stats {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: var(--spacing-md);
  margin-bottom: var(--spacing-lg);
}

.stat-card {
  display: flex;
  align-items: center;
  padding: var(--spacing-lg);
  gap: var(--spacing-md);
}

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: var(--radius-md);
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
  font-size: 26px;
  font-weight: 700;
  color: var(--text-primary);
  line-height: 1;
  margin-bottom: var(--spacing-xs);
}

.stat-label {
  font-size: 13px;
  color: var(--text-secondary);
  margin-bottom: var(--spacing-xs);
}

.stat-change {
  font-size: 11px;
  display: flex;
  align-items: center;
  gap: 3px;
}

.stat-change.up {
  color: var(--color-success);
}

.stat-change.down {
  color: var(--color-danger);
}

.charts-section {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(420px, 1fr));
  gap: var(--spacing-md);
  margin-bottom: var(--spacing-lg);
}

.chart-card {
  padding: var(--spacing-lg);
}

.chart-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--spacing-md);
}

.chart-header h3 {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
}

.chart-container {
  height: 240px;
}

.data-tables {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(500px, 1fr));
  gap: var(--spacing-md);
}

.table-card {
  padding: var(--spacing-lg);
}

.table-header {
  margin-bottom: var(--spacing-md);
}

.table-header h3 {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
}
</style>