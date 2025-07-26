<template>
  <div class="paper-management">
    <!-- 页面标题 -->
    <div class="page-header">
      <h1 class="page-title">组卷管理</h1>
      <p class="page-description">创建和管理试卷，组织题目进行考试</p>
    </div>

    <!-- 操作栏 -->
    <div class="action-bar">
      <el-button type="primary" @click="createPaper">
        <el-icon><Plus /></el-icon>
        创建试卷
      </el-button>
      <div class="search-bar">
        <el-input
          v-model="searchKeyword"
          placeholder="搜索试卷名称或科目"
          clearable
          @input="handleSearch"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
      </div>
      <el-select v-model="filterSubject" placeholder="选择科目" clearable @change="handleFilter">
        <el-option label="全部科目" value="" />
        <el-option
          v-for="subject in subjects"
          :key="subject.id"
          :label="subject.name"
          :value="subject.id"
        />
      </el-select>
    </div>

    <!-- 试卷列表 -->
    <div class="paper-list">
      <el-table :data="filteredPapers" v-loading="loading">
        <el-table-column prop="title" label="试卷名称" min-width="200">
          <template #default="{ row }">
            <div class="paper-title">
              <span>{{ row.title }}</span>
              <el-tag v-if="row.isTemplate" type="info" size="small">模板</el-tag>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="subject" label="科目" min-width="120">
          <template #default="{ row }">
            <el-tag type="primary" size="small">{{ row.subject }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="questionCount" label="题目数量" min-width="100" align="center" />
        <el-table-column prop="totalScore" label="总分" min-width="80" align="center" />
        <el-table-column prop="duration" label="考试时长" min-width="120" align="center">
          <template #default="{ row }">
            {{ row.duration }}分钟
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" min-width="100" align="center">
          <template #default="{ row }">
            <el-tag
              :type="getStatusType(row.status)"
              size="small"
            >
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="创建时间" min-width="180">
          <template #default="{ row }">
            {{ formatDate(row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="viewPaper(row)">查看</el-button>
            <el-button size="small" type="primary" @click="editPaper(row)">编辑</el-button>
            <el-button size="small" type="success" @click="usePaper(row)">使用</el-button>
            <el-button size="small" type="danger" @click="deletePaper(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <!-- 分页 -->
    <div class="pagination">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :total="total"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </div>

    <!-- 创建/编辑试卷对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑试卷' : '创建试卷'"
      width="800px"
      @close="resetForm"
    >
      <el-form
        ref="paperFormRef"
        :model="paperForm"
        :rules="paperRules"
        label-width="100px"
      >
        <el-form-item label="试卷名称" prop="title">
          <el-input v-model="paperForm.title" placeholder="请输入试卷名称" />
        </el-form-item>
        <el-form-item label="科目" prop="subjectId">
          <el-select v-model="paperForm.subjectId" placeholder="请选择科目">
            <el-option
              v-for="subject in subjects"
              :key="subject.id"
              :label="subject.name"
              :value="subject.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="考试时长" prop="duration">
          <el-input-number
            v-model="paperForm.duration"
            :min="30"
            :max="300"
            :step="15"
            controls-position="right"
          />
          <span style="margin-left: 8px;">分钟</span>
        </el-form-item>
        <el-form-item label="试卷描述" prop="description">
          <el-input
            v-model="paperForm.description"
            type="textarea"
            :rows="3"
            placeholder="请输入试卷描述"
          />
        </el-form-item>
        <el-form-item label="设为模板">
          <el-switch v-model="paperForm.isTemplate" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="savePaper" :loading="saving">
          {{ isEdit ? '更新' : '创建' }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Search } from '@element-plus/icons-vue'
import { useRouter } from 'vue-router'

const router = useRouter()

// 响应式数据
const loading = ref(false)
const saving = ref(false)
const searchKeyword = ref('')
const filterSubject = ref('')
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)
const dialogVisible = ref(false)
const isEdit = ref(false)
const paperFormRef = ref()

// 试卷表单
const paperForm = reactive({
  id: '',
  title: '',
  subjectId: '',
  duration: 120,
  description: '',
  isTemplate: false
})

// 表单验证规则
const paperRules = {
  title: [
    { required: true, message: '请输入试卷名称', trigger: 'blur' },
    { min: 2, max: 100, message: '试卷名称长度在 2 到 100 个字符', trigger: 'blur' }
  ],
  subjectId: [
    { required: true, message: '请选择科目', trigger: 'change' }
  ],
  duration: [
    { required: true, message: '请设置考试时长', trigger: 'blur' }
  ]
}

// 模拟数据
const papers = ref([
  {
    id: '1',
    title: '高等数学期末考试',
    subject: '数学',
    subjectId: '1',
    questionCount: 25,
    totalScore: 100,
    duration: 120,
    status: 'published',
    isTemplate: false,
    description: '高等数学期末考试试卷',
    createdAt: '2024-01-15 10:30:00'
  },
  {
    id: '2',
    title: '线性代数模拟试卷',
    subject: '数学',
    subjectId: '1',
    questionCount: 20,
    totalScore: 100,
    duration: 90,
    status: 'draft',
    isTemplate: true,
    description: '线性代数模拟试卷模板',
    createdAt: '2024-01-14 14:20:00'
  },
  {
    id: '3',
    title: '大学英语四级模拟',
    subject: '英语',
    subjectId: '2',
    questionCount: 30,
    totalScore: 100,
    duration: 150,
    status: 'published',
    isTemplate: false,
    description: '大学英语四级模拟试卷',
    createdAt: '2024-01-13 09:15:00'
  }
])

const subjects = ref([
  { id: '1', name: '数学' },
  { id: '2', name: '英语' },
  { id: '3', name: '物理' },
  { id: '4', name: '化学' }
])

// 计算属性
const filteredPapers = computed(() => {
  let result = papers.value
  
  if (searchKeyword.value) {
    result = result.filter(paper => 
      paper.title.toLowerCase().includes(searchKeyword.value.toLowerCase()) ||
      paper.subject.toLowerCase().includes(searchKeyword.value.toLowerCase())
    )
  }
  
  if (filterSubject.value) {
    result = result.filter(paper => paper.subjectId === filterSubject.value)
  }
  
  return result
})

// 方法
const handleSearch = () => {
  currentPage.value = 1
}

const handleFilter = () => {
  currentPage.value = 1
}

const handleSizeChange = (size: number) => {
  pageSize.value = size
  currentPage.value = 1
}

const handleCurrentChange = (page: number) => {
  currentPage.value = page
}

const getStatusType = (status: string) => {
  const statusMap: Record<string, string> = {
    draft: 'info',
    published: 'success',
    archived: 'warning'
  }
  return statusMap[status] || 'info'
}

const getStatusText = (status: string) => {
  const statusMap: Record<string, string> = {
    draft: '草稿',
    published: '已发布',
    archived: '已归档'
  }
  return statusMap[status] || '未知'
}

const formatDate = (dateStr: string) => {
  return new Date(dateStr).toLocaleString('zh-CN')
}

const createPaper = () => {
  isEdit.value = false
  dialogVisible.value = true
}

const editPaper = (paper: any) => {
  isEdit.value = true
  Object.assign(paperForm, paper)
  dialogVisible.value = true
}

const viewPaper = (paper: any) => {
  router.push(`/teacher/papers/${paper.id}`)
}

const usePaper = (paper: any) => {
  router.push(`/teacher/exams/create?paperId=${paper.id}`)
}

const deletePaper = async (paper: any) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除试卷「${paper.title}」吗？此操作不可恢复。`,
      '删除确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    // 模拟删除操作
    const index = papers.value.findIndex(p => p.id === paper.id)
    if (index > -1) {
      papers.value.splice(index, 1)
      ElMessage.success('删除成功')
    }
  } catch {
    // 用户取消删除
  }
}

const savePaper = async () => {
  try {
    await paperFormRef.value.validate()
    saving.value = true
    
    // 模拟保存操作
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    if (isEdit.value) {
      const index = papers.value.findIndex(p => p.id === paperForm.id)
      if (index > -1) {
        papers.value[index] = { ...paperForm }
        ElMessage.success('更新成功')
      }
    } else {
      const newPaper = {
        ...paperForm,
        id: Date.now().toString(),
        questionCount: 0,
        totalScore: 0,
        status: 'draft',
        subject: subjects.value.find(s => s.id === paperForm.subjectId)?.name || '',
        createdAt: new Date().toLocaleString('zh-CN')
      }
      papers.value.unshift(newPaper)
      ElMessage.success('创建成功')
    }
    
    dialogVisible.value = false
  } catch (error) {
    console.error('保存失败:', error)
  } finally {
    saving.value = false
  }
}

const resetForm = () => {
  Object.assign(paperForm, {
    id: '',
    title: '',
    subjectId: '',
    duration: 120,
    description: '',
    isTemplate: false
  })
  paperFormRef.value?.clearValidate()
}

onMounted(() => {
  total.value = papers.value.length
})
</script>

<style scoped>
.paper-management {
  padding: var(--spacing-lg);
}

.page-header {
  margin-bottom: var(--spacing-xl);
}

.page-title {
  font-size: 24px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 var(--spacing-sm) 0;
}

.page-description {
  color: var(--text-secondary);
  margin: 0;
}

.action-bar {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
  margin-bottom: var(--spacing-lg);
  flex-wrap: wrap;
}

.search-bar {
  width: 300px;
}

.paper-list {
  background: white;
  border-radius: var(--radius-lg);
  overflow: hidden;
  box-shadow: var(--shadow-sm);
}

.paper-title {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
}

.pagination {
  display: flex;
  justify-content: center;
  margin-top: var(--spacing-lg);
}

@media (max-width: 768px) {
  .action-bar {
    flex-direction: column;
    align-items: stretch;
  }
  
  .search-bar {
    width: 100%;
  }
}
</style>