<template>
  <div class="admin-dashboard">
    <!-- 页面标题 -->
    <div class="page-header">
      <h1 class="page-title">管理员仪表盘</h1>
      <p class="page-subtitle">系统概览与数据统计</p>
    </div>
    
    <!-- 统计卡片 -->
    <div class="stats-grid">
      <div class="stat-card dopamine-card" v-for="stat in stats" :key="stat.key">
        <div class="stat-icon" :style="{ background: stat.color }">
          <el-icon :size="24">
            <component :is="stat.icon" />
          </el-icon>
        </div>
        <div class="stat-content">
          <div class="stat-value">{{ stat.value }}</div>
          <div class="stat-label">{{ stat.label }}</div>
          <div class="stat-change" :class="stat.trend">
            <el-icon :size="12">
              <component :is="stat.trend === 'up' ? 'TrendCharts' : 'Bottom'" />
            </el-icon>
            <span>{{ stat.change }}</span>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 图表区域 -->
    <div class="charts-section">
      <div class="chart-row">
        <!-- 用户增长趋势 -->
        <div class="chart-card dopamine-card">
          <div class="chart-header">
            <h3>用户增长趋势</h3>
            <el-select v-model="userGrowthPeriod" size="small" style="width: 120px">
              <el-option label="最近7天" value="7d" />
              <el-option label="最近30天" value="30d" />
              <el-option label="最近90天" value="90d" />
            </el-select>
          </div>
          <div class="chart-container">
            <div ref="userGrowthChart" class="chart"></div>
          </div>
        </div>
        
        <!-- 考试活跃度 -->
        <div class="chart-card dopamine-card">
          <div class="chart-header">
            <h3>考试活跃度</h3>
            <el-select v-model="examActivityPeriod" size="small" style="width: 120px">
              <el-option label="今日" value="today" />
              <el-option label="本周" value="week" />
              <el-option label="本月" value="month" />
            </el-select>
          </div>
          <div class="chart-container">
            <div ref="examActivityChart" class="chart"></div>
          </div>
        </div>
      </div>
      
      <div class="chart-row">
        <!-- 系统性能监控 -->
        <div class="chart-card dopamine-card">
          <div class="chart-header">
            <h3>系统性能监控</h3>
            <div class="performance-indicators">
              <div class="indicator">
                <span class="indicator-label">CPU使用率</span>
                <el-progress :percentage="systemPerformance.cpu" :color="getPerformanceColor(systemPerformance.cpu)" />
              </div>
              <div class="indicator">
                <span class="indicator-label">内存使用率</span>
                <el-progress :percentage="systemPerformance.memory" :color="getPerformanceColor(systemPerformance.memory)" />
              </div>
              <div class="indicator">
                <span class="indicator-label">磁盘使用率</span>
                <el-progress :percentage="systemPerformance.disk" :color="getPerformanceColor(systemPerformance.disk)" />
              </div>
            </div>
          </div>
        </div>
        
        <!-- 最近活动 -->
        <div class="activity-card dopamine-card">
          <div class="chart-header">
            <h3>最近活动</h3>
            <el-button size="small" @click="refreshActivities">
              <el-icon><Refresh /></el-icon>
              刷新
            </el-button>
          </div>
          <div class="activity-list">
            <div
              v-for="activity in recentActivities"
              :key="activity.id"
              class="activity-item"
            >
              <div class="activity-icon" :style="{ background: activity.color }">
                <el-icon :size="16">
                  <component :is="activity.icon" />
                </el-icon>
              </div>
              <div class="activity-content">
                <div class="activity-title">{{ activity.title }}</div>
                <div class="activity-desc">{{ activity.description }}</div>
                <div class="activity-time">{{ formatTime(activity.time) }}</div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 快速操作 -->
    <div class="quick-actions">
      <h3>快速操作</h3>
      <div class="action-grid">
        <div
          v-for="action in quickActions"
          :key="action.key"
          class="action-card dopamine-card"
          @click="handleQuickAction(action.key)"
        >
          <div class="action-icon" :style="{ background: action.color }">
            <el-icon :size="20">
              <component :is="action.icon" />
            </el-icon>
          </div>
          <div class="action-content">
            <div class="action-title">{{ action.title }}</div>
            <div class="action-desc">{{ action.description }}</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import {
  UserFilled, School, Notebook, TrendCharts, Bottom,
  Refresh, Plus, Setting, DataBoard, View
} from '@element-plus/icons-vue'
import * as echarts from 'echarts'

const router = useRouter()

// 统计数据
const stats = ref([
  {
    key: 'users',
    label: '总用户数',
    value: '2,847',
    change: '+12.5%',
    trend: 'up',
    icon: 'UserFilled',
    color: 'var(--gradient-primary)'
  },
  {
    key: 'teachers',
    label: '教师数量',
    value: '156',
    change: '+3.2%',
    trend: 'up',
    icon: 'School',
    color: 'var(--gradient-success)'
  },
  {
    key: 'students',
    label: '学生数量',
    value: '2,691',
    change: '+15.8%',
    trend: 'up',
    icon: 'UserFilled',
    color: 'var(--gradient-warning)'
  },
  {
    key: 'exams',
    label: '考试场次',
    value: '1,234',
    change: '+8.7%',
    trend: 'up',
    icon: 'Notebook',
    color: 'var(--gradient-danger)'
  }
])

// 系统性能数据
const systemPerformance = ref({
  cpu: 45,
  memory: 62,
  disk: 38
})

// 最近活动
const recentActivities = ref([
  {
    id: 1,
    title: '新用户注册',
    description: '张三注册了教师账号',
    time: new Date(Date.now() - 5 * 60 * 1000),
    icon: 'UserFilled',
    color: 'var(--dopamine-blue)'
  },
  {
    id: 2,
    title: '考试创建',
    description: '李老师创建了《数据结构》期末考试',
    time: new Date(Date.now() - 15 * 60 * 1000),
    icon: 'Notebook',
    color: 'var(--dopamine-green)'
  },
  {
    id: 3,
    title: '系统更新',
    description: '题库管理模块更新完成',
    time: new Date(Date.now() - 30 * 60 * 1000),
    icon: 'Setting',
    color: 'var(--dopamine-purple)'
  },
  {
    id: 4,
    title: '数据备份',
    description: '每日数据备份任务执行成功',
    time: new Date(Date.now() - 2 * 60 * 60 * 1000),
    icon: 'DataBoard',
    color: 'var(--dopamine-orange)'
  }
])

// 快速操作
const quickActions = ref([
  {
    key: 'add-teacher',
    title: '添加教师',
    description: '快速添加新的教师账号',
    icon: 'Plus',
    color: 'var(--dopamine-blue)'
  },
  {
    key: 'add-student',
    title: '添加学生',
    description: '快速添加新的学生账号',
    icon: 'Plus',
    color: 'var(--dopamine-green)'
  },
  {
    key: 'system-settings',
    title: '系统设置',
    description: '配置系统参数和选项',
    icon: 'Setting',
    color: 'var(--dopamine-purple)'
  },
  {
    key: 'view-statistics',
    title: '查看统计',
    description: '查看详细的数据统计报告',
    icon: 'TrendCharts',
    color: 'var(--dopamine-orange)'
  }
])

// 图表相关
const userGrowthChart = ref<HTMLElement>()
const examActivityChart = ref<HTMLElement>()
const userGrowthPeriod = ref('30d')
const examActivityPeriod = ref('week')

// 获取性能指标颜色
const getPerformanceColor = (value: number) => {
  if (value < 50) return '#67C23A'
  if (value < 80) return '#E6A23C'
  return '#F56C6C'
}

// 格式化时间
const formatTime = (date: Date): string => {
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  
  if (diff < 60 * 1000) {
    return '刚刚'
  } else if (diff < 60 * 60 * 1000) {
    return `${Math.floor(diff / (60 * 1000))}分钟前`
  } else if (diff < 24 * 60 * 60 * 1000) {
    return `${Math.floor(diff / (60 * 60 * 1000))}小时前`
  } else {
    return date.toLocaleDateString()
  }
}

// 初始化用户增长图表
const initUserGrowthChart = () => {
  if (!userGrowthChart.value) return
  
  const chart = echarts.init(userGrowthChart.value)
  const option = {
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'cross'
      }
    },
    legend: {
      data: ['教师', '学生']
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true
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
        name: '教师',
        type: 'line',
        smooth: true,
        data: [120, 132, 101, 134, 90, 156],
        itemStyle: {
          color: '#FF6B8A'
        }
      },
      {
        name: '学生',
        type: 'line',
        smooth: true,
        data: [2200, 2300, 2100, 2400, 2500, 2691],
        itemStyle: {
          color: '#4ECDC4'
        }
      }
    ]
  }
  
  chart.setOption(option)
  
  // 响应式
  window.addEventListener('resize', () => {
    chart.resize()
  })
}

// 初始化考试活跃度图表
const initExamActivityChart = () => {
  if (!examActivityChart.value) return
  
  const chart = echarts.init(examActivityChart.value)
  const option = {
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
  
  // 响应式
  window.addEventListener('resize', () => {
    chart.resize()
  })
}

// 刷新活动列表
const refreshActivities = () => {
  ElMessage.success('活动列表已刷新')
  // 这里可以调用API刷新数据
}

// 处理快速操作
const handleQuickAction = (key: string) => {
  switch (key) {
    case 'add-teacher':
      router.push('/admin/users/teachers')
      break
    case 'add-student':
      router.push('/admin/users/students')
      break
    case 'system-settings':
      router.push('/admin/settings')
      break
    case 'view-statistics':
      router.push('/admin/statistics')
      break
  }
}

// 组件挂载后初始化图表
onMounted(async () => {
  await nextTick()
  initUserGrowthChart()
  initExamActivityChart()
})
</script>

<style scoped>
.admin-dashboard {
  padding: var(--spacing-md);
}

.page-header {
  margin-bottom: var(--spacing-lg);
}

.page-title {
  font-size: 24px;
  font-weight: 700;
  color: var(--text-primary);
  margin: 0 0 var(--spacing-xs) 0;
  background: var(--gradient-primary);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.page-subtitle {
  font-size: 14px;
  color: var(--text-secondary);
  margin: 0;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: var(--spacing-md);
  margin-bottom: var(--spacing-lg);
}

.stat-card {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  padding: var(--spacing-md);
  transition: all 0.3s ease;
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-lg);
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
  font-size: 20px;
  font-weight: 700;
  color: var(--text-primary);
  margin-bottom: 2px;
}

.stat-label {
  font-size: 13px;
  color: var(--text-secondary);
  margin-bottom: 2px;
}

.stat-change {
  display: flex;
  align-items: center;
  gap: 2px;
  font-size: 11px;
  font-weight: 600;
}

.stat-change.up {
  color: var(--dopamine-green);
}

.stat-change.down {
  color: var(--dopamine-red);
}

.charts-section {
  margin-bottom: var(--spacing-lg);
}

.chart-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--spacing-md);
  margin-bottom: var(--spacing-md);
}

.chart-card {
  padding: var(--spacing-md);
}

.chart-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--spacing-sm);
}

.chart-header h3 {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0;
}

.chart-container {
  height: 220px;
}

.chart {
  width: 100%;
  height: 100%;
}

.performance-indicators {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-sm);
  margin-top: var(--spacing-sm);
}

.indicator {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
}

.indicator-label {
  width: 70px;
  font-size: 13px;
  color: var(--text-secondary);
}

.activity-card {
  padding: var(--spacing-md);
}

.activity-list {
  max-height: 220px;
  overflow-y: auto;
}

.activity-item {
  display: flex;
  align-items: flex-start;
  gap: var(--spacing-sm);
  padding: var(--spacing-sm) 0;
  border-bottom: 1px solid var(--border-light);
}

.activity-item:last-child {
  border-bottom: none;
}

.activity-icon {
  width: 28px;
  height: 28px;
  border-radius: var(--radius-sm);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  flex-shrink: 0;
}

.activity-content {
  flex: 1;
}

.activity-title {
  font-size: 13px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 2px;
}

.activity-desc {
  font-size: 12px;
  color: var(--text-secondary);
  margin-bottom: 2px;
}

.activity-time {
  font-size: 11px;
  color: var(--text-muted);
}

.quick-actions h3 {
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 var(--spacing-md) 0;
}

.action-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: var(--spacing-md);
}

.action-card {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  padding: var(--spacing-md);
  cursor: pointer;
  transition: all 0.3s ease;
}

.action-card:hover {
  transform: translateY(-1px);
  box-shadow: var(--shadow-md);
}

.action-icon {
  width: 40px;
  height: 40px;
  border-radius: var(--radius-md);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  flex-shrink: 0;
}

.action-content {
  flex: 1;
}

.action-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 2px;
}

.action-desc {
  font-size: 12px;
  color: var(--text-secondary);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .chart-row {
    grid-template-columns: 1fr;
  }
  
  .stats-grid {
    grid-template-columns: 1fr;
  }
  
  .action-grid {
    grid-template-columns: 1fr;
  }
}
</style>