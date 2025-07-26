<template>
  <div class="exam-management">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-content">
        <h1 class="page-title">
          <el-icon><Document /></el-icon>
          考试管理
        </h1>
        <p class="page-description">
          管理您创建的所有考试，包括编辑、发布、查看统计等操作
        </p>
      </div>
      
      <div class="header-actions">
        <el-button type="primary" @click="createExam">
          <el-icon><Plus /></el-icon>
          创建考试
        </el-button>
      </div>
    </div>
    
    <!-- 统计卡片 -->
    <div class="stats-cards">
      <div class="stat-card">
        <div class="stat-icon">
          <el-icon color="var(--dopamine-blue)"><Document /></el-icon>
        </div>
        <div class="stat-content">
          <div class="stat-number">{{ stats.totalExams }}</div>
          <div class="stat-label">总考试数</div>
        </div>
      </div>
      
      <div class="stat-card">
        <div class="stat-icon">
          <el-icon color="var(--dopamine-green)"><CircleCheck /></el-icon>
        </div>
        <div class="stat-content">
          <div class="stat-number">{{ stats.publishedExams }}</div>
          <div class="stat-label">已发布</div>
        </div>
      </div>
      
      <div class="stat-card">
        <div class="stat-icon">
          <el-icon color="var(--dopamine-orange)"><Clock /></el-icon>
        </div>
        <div class="stat-content">
          <div class="stat-number">{{ stats.ongoingExams }}</div>
          <div class="stat-label">进行中</div>
        </div>
      </div>
      
      <div class="stat-card">
        <div class="stat-icon">
          <el-icon color="var(--dopamine-purple)"><User /></el-icon>
        </div>
        <div class="stat-content">
          <div class="stat-number">{{ stats.totalParticipants }}</div>
          <div class="stat-label">参考人次</div>
        </div>
      </div>
    </div>
    
    <!-- 筛选和搜索 -->
    <div class="filters-section">
      <el-card>
        <el-row :gutter="16">
          <el-col :span="6">
            <el-select
              v-model="filters.status"
              placeholder="考试状态"
              clearable
              @change="loadExams"
            >
              <el-option label="全部" value="" />
              <el-option label="草稿" value="draft" />
              <el-option label="已发布" value="published" />
              <el-option label="进行中" value="ongoing" />
              <el-option label="已结束" value="ended" />
            </el-select>
          </el-col>
          
          <el-col :span="6">
            <el-select
              v-model="filters.subject"
              placeholder="考试科目"
              clearable
              @change="loadExams"
            >
              <el-option
                v-for="subject in subjects"
                :key="subject.value"
                :label="subject.label"
                :value="subject.value"
              />
            </el-select>
          </el-col>
          
          <el-col :span="6">
            <el-date-picker
              v-model="filters.dateRange"
              type="daterange"
              range-separator="至"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              @change="loadExams"
            />
          </el-col>
          
          <el-col :span="6">
            <el-input
              v-model="filters.keyword"
              placeholder="搜索考试名称"
              clearable
              @input="debounceSearch"
            >
              <template #prefix>
                <el-icon><Search /></el-icon>
              </template>
            </el-input>
          </el-col>
        </el-row>
      </el-card>
    </div>
    
    <!-- 考试列表 -->
    <div class="exam-list">
      <el-card>
        <template #header>
          <div class="card-header">
            <span>考试列表</span>
            <div class="header-actions">
              <el-button
                type="primary"
                link
                @click="refreshExams"
                :loading="loading"
              >
                <el-icon><Refresh /></el-icon>
                刷新
              </el-button>
            </div>
          </div>
        </template>
        
        <el-table
          v-loading="loading"
          :data="exams"
          stripe
          @selection-change="handleSelectionChange"
        >
          <el-table-column type="selection" width="55" />
          
          <el-table-column label="考试信息" min-width="300">
            <template #default="{ row }">
              <div class="exam-info">
                <div class="exam-title">{{ row.title }}</div>
                <div class="exam-meta">
                  <el-tag :type="getSubjectTagType(row.subject)" size="small">
                    {{ getSubjectLabel(row.subject) }}
                  </el-tag>
                  <span class="meta-item">
                    <el-icon><Clock /></el-icon>
                    {{ row.duration }}分钟
                  </span>
                  <span class="meta-item">
                    <el-icon><Star /></el-icon>
                    {{ row.totalScore }}分
                  </span>
                </div>
              </div>
            </template>
          </el-table-column>
          
          <el-table-column label="状态" width="120">
            <template #default="{ row }">
              <el-tag :type="getStatusTagType(row.status)">
                {{ getStatusText(row.status) }}
              </el-tag>
            </template>
          </el-table-column>
          
          <el-table-column label="时间安排" width="200">
            <template #default="{ row }">
              <div class="time-info">
                <div class="time-item">
                  <span class="time-label">开始：</span>
                  <span class="time-value">{{ formatDateTime(row.startTime) }}</span>
                </div>
                <div class="time-item">
                  <span class="time-label">结束：</span>
                  <span class="time-value">{{ formatDateTime(row.endTime) }}</span>
                </div>
              </div>
            </template>
          </el-table-column>
          
          <el-table-column label="参考情况" width="150">
            <template #default="{ row }">
              <div class="participation-info">
                <div class="participation-item">
                  <span class="participation-label">已参考：</span>
                  <span class="participation-value">{{ row.participantCount || 0 }}</span>
                </div>
                <div class="participation-item">
                  <span class="participation-label">总人数：</span>
                  <span class="participation-value">{{ row.totalStudents || 0 }}</span>
                </div>
                <div class="participation-progress">
                  <el-progress
                    :percentage="getParticipationPercentage(row)"
                    :stroke-width="4"
                    :show-text="false"
                  />
                </div>
              </div>
            </template>
          </el-table-column>
          
          <el-table-column label="操作" width="200" fixed="right">
            <template #default="{ row }">
              <div class="action-buttons">
                <el-button
                  type="primary"
                  link
                  @click="viewExam(row)"
                >
                  查看
                </el-button>
                
                <el-button
                  v-if="canEdit(row)"
                  type="primary"
                  link
                  @click="editExam(row)"
                >
                  编辑
                </el-button>
                
                <el-dropdown trigger="click">
                  <el-button type="primary" link>
                    更多
                    <el-icon><ArrowDown /></el-icon>
                  </el-button>
                  <template #dropdown>
                    <el-dropdown-menu>
                      <el-dropdown-item
                        v-if="row.status === 'draft'"
                        @click="publishExam(row)"
                      >
                        <el-icon><Promotion /></el-icon>
                        发布考试
                      </el-dropdown-item>
                      
                      <el-dropdown-item
                        v-if="row.status === 'published' && !isExamStarted(row)"
                        @click="cancelExam(row)"
                      >
                        <el-icon><Close /></el-icon>
                        取消考试
                      </el-dropdown-item>
                      
                      <el-dropdown-item
                        @click="duplicateExam(row)"
                      >
                        <el-icon><CopyDocument /></el-icon>
                        复制考试
                      </el-dropdown-item>
                      
                      <el-dropdown-item
                        @click="viewStatistics(row)"
                      >
                        <el-icon><DataAnalysis /></el-icon>
                        查看统计
                      </el-dropdown-item>
                      
                      <el-dropdown-item
                        @click="exportResults(row)"
                      >
                        <el-icon><Download /></el-icon>
                        导出成绩
                      </el-dropdown-item>
                      
                      <el-dropdown-item
                        v-if="canDelete(row)"
                        divided
                        @click="deleteExam(row)"
                      >
                        <el-icon><Delete /></el-icon>
                        删除考试
                      </el-dropdown-item>
                    </el-dropdown-menu>
                  </template>
                </el-dropdown>
              </div>
            </template>
          </el-table-column>
        </el-table>
        
        <!-- 批量操作 -->
        <div v-if="selectedExams.length > 0" class="batch-actions">
          <div class="batch-info">
            已选择 {{ selectedExams.length }} 项
          </div>
          <div class="batch-buttons">
            <el-button @click="batchPublish" :disabled="!canBatchPublish">
              批量发布
            </el-button>
            <el-button @click="batchCancel" :disabled="!canBatchCancel">
              批量取消
            </el-button>
            <el-button type="danger" @click="batchDelete" :disabled="!canBatchDelete">
              批量删除
            </el-button>
          </div>
        </div>
        
        <!-- 分页 -->
        <div class="pagination-wrapper">
          <el-pagination
            v-model:current-page="pagination.page"
            v-model:page-size="pagination.size"
            :total="pagination.total"
            :page-sizes="[10, 20, 50, 100]"
            layout="total, sizes, prev, pager, next, jumper"
            @current-change="loadExams"
            @size-change="loadExams"
          />
        </div>
      </el-card>
    </div>
    
    <!-- 考试详情对话框 -->
    <el-dialog
      v-model="showExamDetail"
      :title="selectedExam?.title"
      width="80%"
      :before-close="handleDetailDialogClose"
    >
      <div v-if="selectedExam" class="exam-detail">
        <el-tabs v-model="activeTab">
          <el-tab-pane label="基本信息" name="basic">
            <div class="detail-section">
              <div class="detail-grid">
                <div class="detail-item">
                  <span class="detail-label">考试名称：</span>
                  <span class="detail-value">{{ selectedExam.title }}</span>
                </div>
                <div class="detail-item">
                  <span class="detail-label">考试科目：</span>
                  <span class="detail-value">{{ getSubjectLabel(selectedExam.subject) }}</span>
                </div>
                <div class="detail-item">
                  <span class="detail-label">考试时长：</span>
                  <span class="detail-value">{{ selectedExam.duration }} 分钟</span>
                </div>
                <div class="detail-item">
                  <span class="detail-label">总分：</span>
                  <span class="detail-value">{{ selectedExam.totalScore }} 分</span>
                </div>
                <div class="detail-item">
                  <span class="detail-label">开始时间：</span>
                  <span class="detail-value">{{ formatDateTime(selectedExam.startTime) }}</span>
                </div>
                <div class="detail-item">
                  <span class="detail-label">结束时间：</span>
                  <span class="detail-value">{{ formatDateTime(selectedExam.endTime) }}</span>
                </div>
                <div class="detail-item">
                  <span class="detail-label">考试状态：</span>
                  <el-tag :type="getStatusTagType(selectedExam.status)">
                    {{ getStatusText(selectedExam.status) }}
                  </el-tag>
                </div>
                <div class="detail-item">
                  <span class="detail-label">创建时间：</span>
                  <span class="detail-value">{{ formatDateTime(selectedExam.createdAt) }}</span>
                </div>
              </div>
              
              <div v-if="selectedExam.description" class="detail-description">
                <h4>考试描述</h4>
                <p>{{ selectedExam.description }}</p>
              </div>
            </div>
          </el-tab-pane>
          
          <el-tab-pane label="题目列表" name="questions">
            <div class="questions-list">
              <div
                v-for="(question, index) in selectedExam.questions"
                :key="question.id"
                class="question-item"
              >
                <div class="question-header">
                  <span class="question-number">第 {{ index + 1 }} 题</span>
                  <div class="question-meta">
                    <el-tag :type="getTypeTagType(question.type)" size="small">
                      {{ getTypeText(question.type) }}
                    </el-tag>
                    <span class="question-score">{{ question.score }}分</span>
                  </div>
                </div>
                <div class="question-content">{{ question.content }}</div>
              </div>
            </div>
          </el-tab-pane>
          
          <el-tab-pane label="参考学生" name="students">
            <div class="students-list">
              <el-table :data="selectedExam.participants" stripe>
                <el-table-column label="学号" prop="studentId" width="120" />
                <el-table-column label="姓名" prop="name" width="120" />
                <el-table-column label="参考状态" width="120">
                  <template #default="{ row }">
                    <el-tag :type="row.status === 'completed' ? 'success' : row.status === 'ongoing' ? 'warning' : ''">
                      {{ getParticipantStatusText(row.status) }}
                    </el-tag>
                  </template>
                </el-table-column>
                <el-table-column label="开始时间" width="150">
                  <template #default="{ row }">
                    {{ row.startTime ? formatDateTime(row.startTime) : '-' }}
                  </template>
                </el-table-column>
                <el-table-column label="提交时间" width="150">
                  <template #default="{ row }">
                    {{ row.submitTime ? formatDateTime(row.submitTime) : '-' }}
                  </template>
                </el-table-column>
                <el-table-column label="得分" width="100">
                  <template #default="{ row }">
                    {{ row.score !== undefined ? `${row.score}/${selectedExam.totalScore}` : '-' }}
                  </template>
                </el-table-column>
              </el-table>
            </div>
          </el-tab-pane>
        </el-tabs>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Document,
  Plus,
  CircleCheck,
  Clock,
  User,
  Search,
  Refresh,
  Star,
  ArrowDown,
  Promotion,
  Close,
  CopyDocument,
  DataAnalysis,
  Download,
  Delete
} from '@element-plus/icons-vue'
import { debounce } from 'lodash'

const router = useRouter()

// 加载状态
const loading = ref(false)

// 统计数据
const stats = ref({
  totalExams: 0,
  publishedExams: 0,
  ongoingExams: 0,
  totalParticipants: 0
})

// 筛选条件
const filters = ref({
  status: '',
  subject: '',
  dateRange: null,
  keyword: ''
})

// 科目选项
const subjects = ref([
  { label: '计算机基础', value: 'computer' },
  { label: '数学', value: 'math' },
  { label: '英语', value: 'english' },
  { label: '物理', value: 'physics' },
  { label: '化学', value: 'chemistry' }
])

// 考试列表
const exams = ref([])
const selectedExams = ref([])
const pagination = ref({
  page: 1,
  size: 10,
  total: 0
})

// 对话框状态
const showExamDetail = ref(false)
const selectedExam = ref(null)
const activeTab = ref('basic')

// 批量操作权限
const canBatchPublish = computed(() => {
  return selectedExams.value.some(exam => exam.status === 'draft')
})

const canBatchCancel = computed(() => {
  return selectedExams.value.some(exam => exam.status === 'published' && !isExamStarted(exam))
})

const canBatchDelete = computed(() => {
  return selectedExams.value.some(exam => canDelete(exam))
})

// 获取科目标签
const getSubjectLabel = (value: string) => {
  const subject = subjects.value.find(s => s.value === value)
  return subject ? subject.label : value
}

// 获取科目标签类型
const getSubjectTagType = (subject: string) => {
  const typeMap: Record<string, string> = {
    computer: '',
    math: 'success',
    english: 'warning',
    physics: 'danger',
    chemistry: 'info'
  }
  return typeMap[subject] || ''
}

// 获取状态文本
const getStatusText = (status: string) => {
  const textMap: Record<string, string> = {
    draft: '草稿',
    published: '已发布',
    ongoing: '进行中',
    ended: '已结束'
  }
  return textMap[status] || status
}

// 获取状态标签类型
const getStatusTagType = (status: string) => {
  const typeMap: Record<string, string> = {
    draft: 'info',
    published: 'success',
    ongoing: 'warning',
    ended: ''
  }
  return typeMap[status] || ''
}

// 获取参考状态文本
const getParticipantStatusText = (status: string) => {
  const textMap: Record<string, string> = {
    not_started: '未开始',
    ongoing: '进行中',
    completed: '已完成'
  }
  return textMap[status] || status
}

// 获取题目类型文本
const getTypeText = (type: string): string => {
  const textMap: Record<string, string> = {
    single: '单选题',
    multiple: '多选题',
    judge: '判断题',
    fill: '填空题',
    essay: '简答题'
  }
  return textMap[type] || ''
}

// 获取题目类型标签类型
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

// 格式化日期时间
const formatDateTime = (date: Date | string | null) => {
  if (!date) return '-'
  const d = new Date(date)
  return new Intl.DateTimeFormat('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  }).format(d)
}

// 获取参考率
const getParticipationPercentage = (exam: any) => {
  if (!exam.totalStudents || exam.totalStudents === 0) return 0
  return Math.round((exam.participantCount / exam.totalStudents) * 100)
}

// 判断考试是否已开始
const isExamStarted = (exam: any) => {
  return new Date(exam.startTime) <= new Date()
}

// 判断是否可以编辑
const canEdit = (exam: any) => {
  return exam.status === 'draft' || (exam.status === 'published' && !isExamStarted(exam))
}

// 判断是否可以删除
const canDelete = (exam: any) => {
  return exam.status === 'draft' || (exam.status === 'ended' && exam.participantCount === 0)
}

// 加载考试列表
const loadExams = async () => {
  try {
    loading.value = true
    
    // 模拟API调用
    const mockExams = [
      {
        id: '1',
        title: '计算机基础知识测试',
        subject: 'computer',
        duration: 120,
        totalScore: 100,
        status: 'published',
        startTime: new Date('2024-01-15 09:00:00'),
        endTime: new Date('2024-01-15 11:00:00'),
        createdAt: new Date('2024-01-10 10:00:00'),
        participantCount: 25,
        totalStudents: 30,
        description: '测试学生对计算机基础知识的掌握程度',
        questions: [
          {
            id: '1',
            type: 'single',
            content: '计算机的核心部件是什么？',
            score: 5
          },
          {
            id: '2',
            type: 'multiple',
            content: '以下哪些是编程语言？',
            score: 8
          }
        ],
        participants: [
          {
            studentId: '2021001',
            name: '张三',
            status: 'completed',
            startTime: new Date('2024-01-15 09:05:00'),
            submitTime: new Date('2024-01-15 10:30:00'),
            score: 85
          },
          {
            studentId: '2021002',
            name: '李四',
            status: 'completed',
            startTime: new Date('2024-01-15 09:02:00'),
            submitTime: new Date('2024-01-15 10:45:00'),
            score: 92
          }
        ]
      },
      {
        id: '2',
        title: '数学期中考试',
        subject: 'math',
        duration: 90,
        totalScore: 150,
        status: 'draft',
        startTime: new Date('2024-01-20 14:00:00'),
        endTime: new Date('2024-01-20 15:30:00'),
        createdAt: new Date('2024-01-12 15:00:00'),
        participantCount: 0,
        totalStudents: 28,
        description: '数学期中考试，涵盖前半学期所学内容',
        questions: [],
        participants: []
      },
      {
        id: '3',
        title: '英语听力测试',
        subject: 'english',
        duration: 60,
        totalScore: 80,
        status: 'ongoing',
        startTime: new Date('2024-01-14 10:00:00'),
        endTime: new Date('2024-01-14 11:00:00'),
        createdAt: new Date('2024-01-08 16:00:00'),
        participantCount: 15,
        totalStudents: 32,
        description: '英语听力能力测试',
        questions: [],
        participants: []
      }
    ]
    
    exams.value = mockExams
    pagination.value.total = mockExams.length
    
    // 更新统计数据
    stats.value = {
      totalExams: mockExams.length,
      publishedExams: mockExams.filter(e => e.status === 'published').length,
      ongoingExams: mockExams.filter(e => e.status === 'ongoing').length,
      totalParticipants: mockExams.reduce((total, exam) => total + exam.participantCount, 0)
    }
  } catch (error) {
    ElMessage.error('加载考试列表失败')
  } finally {
    loading.value = false
  }
}

// 防抖搜索
const debounceSearch = debounce(() => {
  loadExams()
}, 500)

// 刷新考试列表
const refreshExams = () => {
  loadExams()
}

// 处理选择变化
const handleSelectionChange = (selection: any[]) => {
  selectedExams.value = selection
}

// 创建考试
const createExam = () => {
  router.push('/teacher/exam/create')
}

// 查看考试
const viewExam = (exam: any) => {
  selectedExam.value = exam
  activeTab.value = 'basic'
  showExamDetail.value = true
}

// 编辑考试
const editExam = (exam: any) => {
  router.push(`/teacher/exam/edit/${exam.id}`)
}

// 发布考试
const publishExam = async (exam: any) => {
  try {
    await ElMessageBox.confirm(
      `确定要发布考试「${exam.title}」吗？发布后学生将可以参加考试。`,
      '发布考试',
      {
        confirmButtonText: '确定发布',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    // 模拟发布
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    exam.status = 'published'
    ElMessage.success('考试发布成功')
    loadExams()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('发布失败')
    }
  }
}

// 取消考试
const cancelExam = async (exam: any) => {
  try {
    await ElMessageBox.confirm(
      `确定要取消考试「${exam.title}」吗？取消后学生将无法参加考试。`,
      '取消考试',
      {
        confirmButtonText: '确定取消',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    // 模拟取消
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    exam.status = 'draft'
    ElMessage.success('考试已取消')
    loadExams()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('取消失败')
    }
  }
}

// 复制考试
const duplicateExam = async (exam: any) => {
  try {
    await ElMessageBox.confirm(
      `确定要复制考试「${exam.title}」吗？`,
      '复制考试',
      {
        confirmButtonText: '确定复制',
        cancelButtonText: '取消',
        type: 'info'
      }
    )
    
    // 模拟复制
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    ElMessage.success('考试复制成功')
    loadExams()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('复制失败')
    }
  }
}

// 查看统计
const viewStatistics = (exam: any) => {
  router.push(`/teacher/exam/statistics/${exam.id}`)
}

// 导出成绩
const exportResults = async (exam: any) => {
  try {
    ElMessage.info('正在导出成绩...')
    
    // 模拟导出
    await new Promise(resolve => setTimeout(resolve, 2000))
    
    ElMessage.success('成绩导出成功')
  } catch (error) {
    ElMessage.error('导出失败')
  }
}

// 删除考试
const deleteExam = async (exam: any) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除考试「${exam.title}」吗？此操作不可恢复。`,
      '删除考试',
      {
        confirmButtonText: '确定删除',
        cancelButtonText: '取消',
        type: 'error'
      }
    )
    
    // 模拟删除
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    ElMessage.success('考试删除成功')
    loadExams()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

// 批量发布
const batchPublish = async () => {
  const draftExams = selectedExams.value.filter(exam => exam.status === 'draft')
  if (draftExams.length === 0) {
    ElMessage.warning('没有可发布的考试')
    return
  }
  
  try {
    await ElMessageBox.confirm(
      `确定要发布选中的 ${draftExams.length} 个考试吗？`,
      '批量发布',
      {
        confirmButtonText: '确定发布',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    // 模拟批量发布
    await new Promise(resolve => setTimeout(resolve, 1500))
    
    ElMessage.success(`成功发布 ${draftExams.length} 个考试`)
    loadExams()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('批量发布失败')
    }
  }
}

// 批量取消
const batchCancel = async () => {
  const publishedExams = selectedExams.value.filter(exam => exam.status === 'published' && !isExamStarted(exam))
  if (publishedExams.length === 0) {
    ElMessage.warning('没有可取消的考试')
    return
  }
  
  try {
    await ElMessageBox.confirm(
      `确定要取消选中的 ${publishedExams.length} 个考试吗？`,
      '批量取消',
      {
        confirmButtonText: '确定取消',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    // 模拟批量取消
    await new Promise(resolve => setTimeout(resolve, 1500))
    
    ElMessage.success(`成功取消 ${publishedExams.length} 个考试`)
    loadExams()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('批量取消失败')
    }
  }
}

// 批量删除
const batchDelete = async () => {
  const deletableExams = selectedExams.value.filter(exam => canDelete(exam))
  if (deletableExams.length === 0) {
    ElMessage.warning('没有可删除的考试')
    return
  }
  
  try {
    await ElMessageBox.confirm(
      `确定要删除选中的 ${deletableExams.length} 个考试吗？此操作不可恢复。`,
      '批量删除',
      {
        confirmButtonText: '确定删除',
        cancelButtonText: '取消',
        type: 'error'
      }
    )
    
    // 模拟批量删除
    await new Promise(resolve => setTimeout(resolve, 1500))
    
    ElMessage.success(`成功删除 ${deletableExams.length} 个考试`)
    loadExams()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('批量删除失败')
    }
  }
}

// 处理详情对话框关闭
const handleDetailDialogClose = (done: () => void) => {
  selectedExam.value = null
  done()
}

onMounted(() => {
  loadExams()
})
</script>

<style scoped>
.exam-management {
  padding: var(--spacing-md);
  background: var(--bg-primary);
  min-height: 100vh;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  background: white;
  padding: var(--spacing-lg);
  border-radius: var(--radius-lg);
  margin-bottom: var(--spacing-md);
  box-shadow: var(--shadow-sm);
}

.header-content {
  flex: 1;
}

.page-title {
  font-size: 28px;
  font-weight: 700;
  color: var(--text-primary);
  margin: 0 0 var(--spacing-sm) 0;
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
}

.page-description {
  font-size: 16px;
  color: var(--text-secondary);
  margin: 0;
  line-height: 1.5;
}

.header-actions {
  display: flex;
  gap: var(--spacing-md);
}

.stats-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: var(--spacing-lg);
  margin-bottom: var(--spacing-lg);
}

.stat-card {
  display: flex;
  align-items: center;
  gap: var(--spacing-lg);
  background: white;
  padding: var(--spacing-lg);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-sm);
  transition: all 0.3s ease;
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
}

.stat-icon {
  width: 48px;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--bg-secondary);
  border-radius: var(--radius-md);
  font-size: 24px;
}

.stat-content {
  flex: 1;
}

.stat-number {
  font-size: 24px;
  font-weight: 700;
  color: var(--text-primary);
  margin-bottom: var(--spacing-xs);
}

.stat-label {
  font-size: 14px;
  color: var(--text-secondary);
}

.filters-section {
  margin-bottom: var(--spacing-lg);
}

.exam-list {
  margin-bottom: var(--spacing-lg);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.exam-info {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-sm);
}

.exam-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: var(--spacing-xs);
}

.exam-meta {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
}

.meta-item {
  display: flex;
  align-items: center;
  gap: var(--spacing-xs);
  font-size: 12px;
  color: var(--text-secondary);
}

.time-info {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-xs);
}

.time-item {
  display: flex;
  gap: var(--spacing-xs);
  font-size: 12px;
}

.time-label {
  color: var(--text-secondary);
  min-width: 32px;
}

.time-value {
  color: var(--text-primary);
}

.participation-info {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-xs);
}

.participation-item {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
}

.participation-label {
  color: var(--text-secondary);
}

.participation-value {
  color: var(--text-primary);
  font-weight: 500;
}

.participation-progress {
  margin-top: var(--spacing-xs);
}

.action-buttons {
  display: flex;
  gap: var(--spacing-sm);
  align-items: center;
}

.batch-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--spacing-md);
  background: var(--dopamine-blue-light);
  border-radius: var(--radius-md);
  margin-top: var(--spacing-md);
}

.batch-info {
  font-size: 14px;
  color: var(--dopamine-blue);
  font-weight: 500;
}

.batch-buttons {
  display: flex;
  gap: var(--spacing-sm);
}

.pagination-wrapper {
  display: flex;
  justify-content: center;
  margin-top: var(--spacing-lg);
}

.exam-detail {
  max-height: 70vh;
  overflow-y: auto;
}

.detail-section {
  padding: var(--spacing-lg);
}

.detail-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: var(--spacing-md);
  margin-bottom: var(--spacing-lg);
}

.detail-item {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-xs);
  padding: var(--spacing-md);
  background: var(--bg-secondary);
  border-radius: var(--radius-md);
}

.detail-label {
  font-size: 12px;
  color: var(--text-secondary);
  font-weight: 500;
}

.detail-value {
  font-size: 14px;
  color: var(--text-primary);
  font-weight: 600;
}

.detail-description {
  padding: var(--spacing-lg);
  background: var(--bg-secondary);
  border-radius: var(--radius-md);
}

.detail-description h4 {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 var(--spacing-md) 0;
}

.detail-description p {
  font-size: 14px;
  line-height: 1.6;
  color: var(--text-primary);
  margin: 0;
}

.questions-list {
  padding: var(--spacing-lg);
}

.question-item {
  padding: var(--spacing-lg);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-md);
  margin-bottom: var(--spacing-md);
  background: white;
}

.question-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--spacing-md);
}

.question-number {
  font-weight: 600;
  color: var(--dopamine-blue);
}

.question-meta {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
}

.question-score {
  font-size: 12px;
  font-weight: 600;
  color: var(--dopamine-orange);
  background: var(--dopamine-orange-light);
  padding: 2px 6px;
  border-radius: var(--radius-sm);
}

.question-content {
  font-size: 14px;
  line-height: 1.5;
  color: var(--text-primary);
}

.students-list {
  padding: var(--spacing-lg);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .exam-management {
    padding: var(--spacing-md);
  }
  
  .page-header {
    flex-direction: column;
    gap: var(--spacing-lg);
  }
  
  .header-actions {
    width: 100%;
    justify-content: flex-end;
  }
  
  .stats-cards {
    grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  }
  
  .exam-meta {
    flex-direction: column;
    align-items: flex-start;
    gap: var(--spacing-xs);
  }
  
  .action-buttons {
    flex-direction: column;
    align-items: stretch;
  }
  
  .batch-actions {
    flex-direction: column;
    gap: var(--spacing-md);
  }
  
  .detail-grid {
    grid-template-columns: 1fr;
  }
}
</style>