<template>
  <div class="practice-view">
    <!-- 页面标题 -->
    <div class="page-header">
      <h1 class="page-title">练习刷题</h1>
      <p class="page-subtitle">通过练习巩固知识点，提升答题技能</p>
    </div>
    
    <!-- 筛选条件 -->
    <div class="filter-section dopamine-card">
      <div class="filter-row">
        <div class="filter-item">
          <label>科目</label>
          <el-select v-model="filters.subject" placeholder="选择科目" clearable>
            <el-option
              v-for="subject in subjects"
              :key="subject.id"
              :label="subject.name"
              :value="subject.id"
            />
          </el-select>
        </div>
        
        <div class="filter-item">
          <label>题目类型</label>
          <el-select v-model="filters.type" placeholder="选择类型" clearable>
            <el-option label="单选题" value="single" />
            <el-option label="多选题" value="multiple" />
            <el-option label="判断题" value="judge" />
            <el-option label="填空题" value="fill" />
          </el-select>
        </div>
        
        <div class="filter-item">
          <label>难度等级</label>
          <el-select v-model="filters.difficulty" placeholder="选择难度" clearable>
            <el-option label="简单" value="easy" />
            <el-option label="中等" value="medium" />
            <el-option label="困难" value="hard" />
          </el-select>
        </div>
        
        <div class="filter-item">
          <label>练习模式</label>
          <el-select v-model="practiceMode" placeholder="选择模式">
            <el-option label="顺序练习" value="sequence" />
            <el-option label="随机练习" value="random" />
            <el-option label="错题重练" value="wrong" />
          </el-select>
        </div>
        
        <div class="filter-actions">
          <el-button type="primary" @click="startPractice">
            <el-icon><VideoPlay /></el-icon>
            开始练习
          </el-button>
          <el-button @click="resetFilters">
            <el-icon><Refresh /></el-icon>
            重置
          </el-button>
        </div>
      </div>
    </div>
    
    <!-- 练习统计 -->
    <div class="stats-section">
      <div class="stats-grid">
        <div class="stat-card dopamine-card">
          <div class="stat-icon" style="background: var(--gradient-primary)">
            <el-icon :size="24"><TrendCharts /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ practiceStats.total_practiced }}</div>
            <div class="stat-label">累计练习</div>
          </div>
        </div>
        
        <div class="stat-card dopamine-card">
          <div class="stat-icon" style="background: var(--gradient-success)">
            <el-icon :size="24"><SuccessFilled /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ practiceStats.correct_rate }}%</div>
            <div class="stat-label">正确率</div>
          </div>
        </div>
        
        <div class="stat-card dopamine-card">
          <div class="stat-icon" style="background: var(--gradient-warning)">
            <el-icon :size="24"><Clock /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ practiceStats.today_practiced }}</div>
            <div class="stat-label">今日练习</div>
          </div>
        </div>
        
        <div class="stat-card dopamine-card">
          <div class="stat-icon" style="background: var(--gradient-danger)">
            <el-icon :size="24"><Warning /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ practiceStats.wrong_questions }}</div>
            <div class="stat-label">错题数量</div>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 推荐练习 -->
    <div class="recommendations" style="margin-bottom: 60px;">
      <div class="section-header">
        <h2>推荐练习</h2>
        <el-button type="primary" @click="refreshRecommendations">
          <el-icon><Refresh /></el-icon>
          刷新推荐
        </el-button>
      </div>
      
      <div v-loading="loading" class="recommendation-carousel" @mouseenter="showControls = true" @mouseleave="showControls = false">
        <!-- 左箭头 -->
        <div 
          v-show="showControls && totalPages > 1" 
          class="carousel-arrow carousel-arrow-left"
          @click="previousPage"
        >
          <el-icon><ArrowLeft /></el-icon>
        </div>
        
        <!-- 右箭头 -->
        <div 
          v-show="showControls && totalPages > 1" 
          class="carousel-arrow carousel-arrow-right"
          @click="nextPage"
        >
          <el-icon><ArrowRight /></el-icon>
        </div>
        
        <!-- 卡片容器 -->
         <div class="recommendation-grid">
           <div
             v-for="recommendation in currentPageRecommendations"
             :key="recommendation.id"
             class="recommendation-card dopamine-card"
             @click="startRecommendedPractice(recommendation)"
           >
            <div class="recommendation-header">
              <div class="recommendation-subject">{{ recommendation.subject?.name || '未知科目' }}</div>
              <div class="recommendation-difficulty" :class="`difficulty-${recommendation.difficulty}`">
                {{ getDifficultyText(recommendation.difficulty) }}
              </div>
            </div>
            
            <div class="recommendation-title">{{ recommendation.title }}</div>
            <div class="recommendation-desc">{{ recommendation.description }}</div>
            
            <div class="recommendation-stats">
              <div class="recommendation-stat">
                <el-icon><Document /></el-icon>
                <span>{{ recommendation.question_count }} 题</span>
              </div>
              <div class="recommendation-stat">
                <el-icon><Timer /></el-icon>
                <span>约 {{ recommendation.estimated_time }} 分钟</span>
              </div>
              <div class="recommendation-stat">
                <el-icon><Star /></el-icon>
                <span>{{ recommendation.rating.toFixed(1) }} 分</span>
              </div>
            </div>
            
            <div class="recommendation-progress">
              <el-progress
                :percentage="recommendation.progress"
                :color="getProgressColor(recommendation.progress)"
                :stroke-width="6"
                :show-text="true"
                :format="(percentage) => `${percentage.toFixed(2)}%`"
              />
            </div>
          </div>
        </div>
        
        <!-- 底部点状指示器 -->
        <div 
          v-show="showControls && totalPages > 1" 
          class="carousel-dots"
        >
          <div
            v-for="(page, index) in totalPages"
            :key="index"
            class="carousel-dot"
            :class="{ active: currentPage === index }"
            @click="goToPage(index)"
          ></div>
        </div>
      </div>
    </div>
    
    <!-- 练习历史 -->
    <div class="practice-history">
      <div class="section-header">
        <h2>练习历史</h2>
        <el-button type="primary" @click="viewAllHistory">
          查看全部
        </el-button>
      </div>
      
      <div class="history-table">
        <el-table :data="practiceHistory" style="width: 100%">
          <el-table-column label="科目" width="120">
            <template #default="{ row }">
              {{ row.subject?.name || '未知科目' }}
            </template>
          </el-table-column>
          <el-table-column label="练习内容">
            <template #default="{ row }">
              {{ row.title || '练习记录' }}
            </template>
          </el-table-column>
          <el-table-column label="题目类型" width="100">
            <template #default="{ row }">
              <el-tag type="info">练习</el-tag>
            </template>
          </el-table-column>
          <el-table-column label="题目数" width="80">
            <template #default="{ row }">
              {{ row.total_questions }}
            </template>
          </el-table-column>
          <el-table-column label="正确数" width="80">
            <template #default="{ row }">
              {{ row.correct_count }}
            </template>
          </el-table-column>
          <el-table-column label="正确率" width="80">
            <template #default="{ row }">
              <span :class="getCorrectRateClass(getCorrectRate(row))">{{ getCorrectRate(row).toFixed(1) }}%</span>
            </template>
          </el-table-column>
          <el-table-column label="用时" width="80">
            <template #default="{ row }">
              {{ Math.floor(row.duration / 60) }}分{{ row.duration % 60 }}秒
            </template>
          </el-table-column>
          <el-table-column label="练习时间" width="120">
            <template #default="{ row }">
              {{ formatDate(row.created_at) }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="120">
            <template #default="{ row }">
              <el-button
                type="primary"
                size="small"
                @click="reviewPractice(row.id)"
              >
                复习
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import {
  VideoPlay, Refresh, TrendCharts, SuccessFilled, Clock, Warning,
  Document, Timer, Star, ArrowLeft, ArrowRight
} from '@element-plus/icons-vue'
import * as practiceApi from '@/api/practice'
import * as subjectApi from '@/api/subject'

const router = useRouter()

// 筛选条件
const filters = reactive({
  subject: '',
  type: '',
  difficulty: '',
  mode: ''
})

// 科目选项
const subjectOptions = ref([
  { label: '全部科目', value: '' }
])

// 练习统计
const practiceStats = ref<practiceApi.PracticeStats>({
  total_practiced: 0,
  correct_rate: 0,
  today_practiced: 0,
  wrong_questions: 0
})

// 推荐练习数据
const recommendations = ref<practiceApi.PracticeRecommendation[]>([])

// 加载状态
const loading = ref(false)

// 练习历史数据
const practiceHistory = ref<practiceApi.PracticeRecord[]>([])

// 轮播相关
const currentPage = ref(0)
const showControls = ref(false)
const cardsPerPage = 6 // 每页显示的卡片数量

// 计算总页数
const totalPages = computed(() => {
  return Math.ceil(recommendations.value.length / cardsPerPage)
})

// 计算当前页面的推荐卡片
const currentPageRecommendations = computed(() => {
  const start = currentPage.value * cardsPerPage
  const end = start + cardsPerPage
  return recommendations.value.slice(start, end)
})

// 轮播方法
const nextPage = () => {
  if (currentPage.value < totalPages.value - 1) {
    currentPage.value++
  } else {
    currentPage.value = 0 // 循环到第一页
  }
}

const previousPage = () => {
  if (currentPage.value > 0) {
    currentPage.value--
  } else {
    currentPage.value = totalPages.value - 1 // 循环到最后一页
  }
}

const goToPage = (pageIndex: number) => {
  currentPage.value = pageIndex
}

// 获取难度文本
const getDifficultyText = (difficulty: number) => {
  const difficultyMap: Record<number, string> = {
    1: '简单',
    2: '中等', 
    3: '困难',
    4: '专家',
    5: '大师'
  }
  return difficultyMap[difficulty] || '未知'
}

// 获取难度样式类
const getDifficultyClass = (difficulty: number) => {
  const classMap: Record<number, string> = {
    1: 'difficulty-easy',
    2: 'difficulty-medium',
    3: 'difficulty-hard', 
    4: 'difficulty-expert',
    5: 'difficulty-master'
  }
  return classMap[difficulty] || 'difficulty-medium'
}

// 获取进度颜色
const getProgressColor = (progress: number): string => {
  if (progress >= 80) return '#67C23A'
  if (progress >= 50) return '#E6A23C'
  return '#F56C6C'
}

// 计算正确率
const getCorrectRate = (row: practiceApi.PracticeRecord) => {
  if (row.total_count === 0) return 0
  return (row.correct_count / row.total_count) * 100
}

// 获取正确率样式类
const getCorrectRateClass = (rate: number) => {
  if (rate >= 80) return 'rate-high'
  if (rate >= 60) return 'rate-medium'
  return 'rate-low'
}

// 格式化日期
const formatDate = (dateString: string) => {
  const date = new Date(dateString)
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  const hours = Math.floor(diff / (1000 * 60 * 60))
  
  if (hours < 1) {
    const minutes = Math.floor(diff / (1000 * 60))
    return `${minutes}分钟前`
  } else if (hours < 24) {
    return `${hours}小时前`
  } else {
    const days = Math.floor(hours / 24)
    if (days < 7) {
      return `${days}天前`
    } else {
      return date.toLocaleDateString('zh-CN')
    }
  }
}

// 开始练习
const startPractice = () => {
  if (!filters.subject) {
    ElMessage.warning('请选择练习科目')
    return
  }
  
  ElMessage.success('正在为您准备练习题目...')
  // 这里可以跳转到练习页面
}

// 重置筛选条件
const resetFilters = () => {
  filters.subject = ''
  filters.type = ''
  filters.difficulty = ''
  practiceMode.value = 'sequence'
}

// 加载科目列表
const loadSubjects = async () => {
  try {
    const response = await subjectApi.getSubjects()
    const subjects = response.subjects.map(subject => ({
      label: subject.name,
      value: subject.id.toString()
    }))
    subjectOptions.value = [
      { label: '全部科目', value: '' },
      ...subjects
    ]
  } catch (error) {
    console.error('加载科目列表失败:', error)
  }
}

// 加载练习统计
const loadPracticeStats = async () => {
  try {
    const stats = await practiceApi.getPracticeStats()
    practiceStats.value = stats
  } catch (error) {
    console.error('加载练习统计失败:', error)
  }
}

// 加载推荐练习
const loadRecommendations = async () => {
  try {
    loading.value = true
    const recs = await practiceApi.getPracticeRecommendations()
    recommendations.value = recs
  } catch (error) {
    console.error('加载推荐练习失败:', error)
    ElMessage.error('加载推荐练习失败')
  } finally {
    loading.value = false
  }
}

// 加载练习历史
const loadPracticeHistory = async () => {
  try {
    const response = await practiceApi.getPracticeHistory()
    // response已经是提取后的data字段，包含{data: [...], total: ...}结构
    practiceHistory.value = response.data
  } catch (error) {
    console.error('加载练习历史失败:', error)
  }
}

// 刷新推荐
const refreshRecommendations = () => {
  loadRecommendations()
  ElMessage.success('推荐已刷新')
}

// 开始推荐练习
const startRecommendedPractice = async (recommendation: practiceApi.PracticeRecommendation) => {
  try {
    // 直接跳转到练习页面，传递练习参数
    router.push({
      path: `/student/practice/${recommendation.subject_id}/sequence`,
      query: {
        title: recommendation.title,
        subject: recommendation.subject.name,
        difficulty: recommendation.difficulty,
        question_count: recommendation.question_count
      }
    })
  } catch (error) {
    console.error('开始练习失败:', error)
    ElMessage.error('开始练习失败')
  }
}

// 初始化数据
const initData = async () => {
  await Promise.all([
    loadSubjects(),
    loadPracticeStats(),
    loadRecommendations(),
    loadPracticeHistory()
  ])
}

// 组件挂载时加载数据
onMounted(() => {
  initData()
})

// 查看全部历史
const viewAllHistory = () => {
  ElMessage.info('查看全部历史功能开发中')
}

// 复习练习
const reviewPractice = (practiceId: number) => {
  router.push({
    path: '/student/review',
    query: { practiceRecordId: practiceId }
  })
}
</script>

<style scoped>
.practice-view {
  padding: var(--spacing-md);
}

.page-header {
  margin-bottom: var(--spacing-lg);
}

.page-title {
  font-size: 28px;
  font-weight: 700;
  color: var(--text-primary);
  margin: 0 0 var(--spacing-sm) 0;
  background: var(--gradient-primary);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.page-subtitle {
  font-size: 16px;
  color: var(--text-secondary);
  margin: 0;
}

.filter-section {
  padding: var(--spacing-md);
  margin-bottom: var(--spacing-lg);
}

.filter-row {
  display: flex;
  flex-wrap: wrap;
  gap: var(--spacing-lg);
  align-items: end;
}

.filter-item {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-sm);
  min-width: 150px;
}

.filter-item label {
  font-size: 14px;
  font-weight: 500;
  color: var(--text-primary);
}

.filter-actions {
  display: flex;
  gap: var(--spacing-md);
}

.stats-section {
  margin-bottom: var(--spacing-lg);
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(140px, 1fr));
  gap: var(--spacing-md);
}

.stat-card {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  padding: var(--spacing-md);
  transition: all 0.3s ease;
}

.stat-card:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow-xl);
}

.stat-icon {
  width: 42px;
  height: 42px;
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
  font-size: 18px;
  font-weight: 700;
  color: var(--text-primary);
  margin-bottom: 2px;
}

.stat-label {
  font-size: 12px;
  color: var(--text-secondary);
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--spacing-lg);
}

.section-header h2 {
  font-size: 20px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0;
}

.recommendations {
  margin-bottom: var(--spacing-lg);
}

/* 轮播容器 */
.recommendation-carousel {
    position: relative;
    overflow: hidden;
    border-radius: var(--radius-lg);
    padding: 0 60px; /* 为左右按钮留出空间 */
    height: 260px; /* 固定高度，不包含dots */
  }

.recommendation-grid {
  display: flex;
  gap: var(--spacing-md);
  width: 100%;
  height: 260px; /* 固定高度 */
  overflow-x: auto;
  overflow-y: hidden;
  scroll-behavior: smooth;
  padding-bottom: 8px; /* 为滚动条留出空间 */
}

/* 轮播箭头 */
.carousel-arrow {
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  width: 40px; /* 缩小约17% */
  height: 40px; /* 缩小约17% */
  background: rgba(255, 255, 255, 0.95);
  border: 1px solid rgba(226, 232, 240, 0.8);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  z-index: 10;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  backdrop-filter: blur(10px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.carousel-arrow:hover {
  background: var(--dopamine-blue);
  color: white;
  transform: translateY(-50%) scale(1.1);
  box-shadow: 0 8px 24px rgba(77, 150, 255, 0.3);
}

.carousel-arrow-left {
  left: 10px; /* 调整位置避免遮挡 */
}

.carousel-arrow-right {
  right: 10px; /* 调整位置避免遮挡 */
}

/* 点状指示器 */
.carousel-dots {
    position: absolute;
    bottom: -40px; /* 位于carousel外部 */
    left: 50%;
    transform: translateX(-50%);
    display: flex;
    justify-content: center;
    gap: 8px;
    opacity: 1; /* 保持可见 */
    height: 26px;
    align-items: center;
    z-index: 5;
  }

.carousel-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  background: rgba(226, 232, 240, 0.6);
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  border: 2px solid transparent;
}

.carousel-dot:hover {
  background: rgba(77, 150, 255, 0.5);
  transform: scale(1.2);
}

.carousel-dot.active {
  background: var(--dopamine-blue);
  transform: scale(1.3);
  box-shadow: 0 0 0 3px rgba(77, 150, 255, 0.2);
}

.recommendation-card {
  padding: var(--spacing-md);
  cursor: pointer;
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  border: 1px solid rgba(226, 232, 240, 0.8);
  position: relative;
  overflow: hidden;
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.95) 0%, rgba(248, 250, 252, 0.95) 100%);
  backdrop-filter: blur(10px);
  height: 230px; /* 调整卡片高度适应容器 */
  min-width: 280px; /* 设置最小宽度防止压缩 */
  flex-shrink: 0; /* 防止卡片被压缩 */
  display: flex;
  flex-direction: column;
}

.recommendation-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 3px;
  background: linear-gradient(90deg, var(--dopamine-blue) 0%, var(--dopamine-purple) 50%, var(--dopamine-pink) 100%);
  opacity: 0;
  transition: opacity 0.3s ease;
}

.recommendation-card:hover {
  transform: translateY(-6px) scale(1.02);
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1), 0 8px 16px rgba(0, 0, 0, 0.06);
  border-color: rgba(77, 150, 255, 0.3);
}

.recommendation-card:hover::before {
  opacity: 1;
}

.recommendation-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--spacing-sm);
  position: relative;
}

.recommendation-subject {
  font-size: 11px;
  font-weight: 700;
  background: linear-gradient(135deg, var(--dopamine-blue), var(--dopamine-purple));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  position: relative;
}

.recommendation-difficulty {
  font-size: 10px;
  font-weight: 700;
  padding: 4px 8px;
  border-radius: 12px;
  text-transform: uppercase;
  letter-spacing: 0.3px;
  border: 1px solid rgba(255, 255, 255, 0.3);
  backdrop-filter: blur(8px);
}

.difficulty-1,
.difficulty-easy {
  background: linear-gradient(135deg, rgba(107, 207, 127, 0.15), rgba(107, 207, 127, 0.25));
  color: var(--dopamine-green);
  box-shadow: 0 2px 8px rgba(107, 207, 127, 0.2);
}

.difficulty-2,
.difficulty-medium {
  background: linear-gradient(135deg, rgba(255, 140, 66, 0.15), rgba(255, 140, 66, 0.25));
  color: var(--dopamine-orange);
  box-shadow: 0 2px 8px rgba(255, 140, 66, 0.2);
}

.difficulty-3,
.difficulty-hard {
  background: linear-gradient(135deg, rgba(255, 87, 87, 0.15), rgba(255, 87, 87, 0.25));
  color: var(--dopamine-red);
  box-shadow: 0 2px 8px rgba(255, 87, 87, 0.2);
}

.difficulty-4,
.difficulty-expert {
  background: linear-gradient(135deg, rgba(155, 89, 182, 0.15), rgba(155, 89, 182, 0.25));
  color: var(--dopamine-purple);
  box-shadow: 0 2px 8px rgba(155, 89, 182, 0.2);
}

.recommendation-title {
  font-size: 16px;
  font-weight: 700;
  color: var(--text-primary);
  margin-bottom: var(--spacing-sm);
  line-height: 1.3;
  letter-spacing: -0.02em;
}

.recommendation-desc {
  font-size: 13px;
  color: var(--text-secondary);
  margin-bottom: var(--spacing-md);
  line-height: 1.5;
  opacity: 0.8;
  height: 40px; /* 固定高度确保对齐 */
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  flex-shrink: 0;
}

.recommendation-stats {
  display: flex;
  flex-wrap: wrap;
  gap: var(--spacing-sm);
  margin-bottom: auto; /* 推到底部 */
  padding: var(--spacing-sm) 0;
  border-top: 1px solid rgba(226, 232, 240, 0.6);
  border-bottom: 1px solid rgba(226, 232, 240, 0.6);
  margin-top: auto;
}

.recommendation-stat {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: var(--text-secondary);
  font-weight: 500;
}

.recommendation-stat .el-icon {
  color: var(--dopamine-blue);
  opacity: 0.7;
}

.recommendation-progress {
  margin-top: var(--spacing-sm);
}

.progress-label {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--spacing-sm);
  font-size: 14px;
  color: var(--text-secondary);
}

.practice-history {
  margin-bottom: var(--spacing-xl);
}

.history-table {
  background: white;
  border-radius: var(--radius-lg);
  overflow: hidden;
  box-shadow: var(--shadow-sm);
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

/* 滚动条样式 */
.recommendation-grid::-webkit-scrollbar {
  height: 6px;
}

.recommendation-grid::-webkit-scrollbar-track {
  background: rgba(0, 0, 0, 0.05);
  border-radius: 3px;
}

.recommendation-grid::-webkit-scrollbar-thumb {
  background: rgba(0, 0, 0, 0.15);
  border-radius: 3px;
  transition: background 0.3s ease;
}

.recommendation-grid::-webkit-scrollbar-thumb:hover {
  background: rgba(0, 0, 0, 0.25);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .filter-row {
    flex-direction: column;
    align-items: stretch;
  }
  
  .filter-actions {
    justify-content: center;
  }
  
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }
  
  .recommendation-carousel {
    padding: 0 20px; /* 移动端减少左右间距 */
  }
  
  .recommendation-card {
    min-width: 260px; /* 移动端稍微减小卡片宽度 */
  }
}

@media (max-width: 480px) {
  .stats-grid {
    grid-template-columns: 1fr;
  }
}
</style>