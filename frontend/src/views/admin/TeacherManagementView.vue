<template>
  <div class="teacher-management-view">
    <!-- 页面标题 -->
    <div class="page-header">
      <h1 class="page-title">教师管理</h1>
      <p class="page-subtitle">管理系统中的教师账户</p>
    </div>

    <!-- 操作栏 -->
    <div class="action-bar dopamine-card">
      <div class="action-left">
        <el-button type="primary" :icon="Plus" @click="showAddDialog = true">
          添加教师
        </el-button>
        <el-button :icon="Upload" @click="showImportDialog = true">
          批量导入
        </el-button>
        <el-button :icon="Download" @click="exportTeachers">
          导出数据
        </el-button>
      </div>
      <div class="action-right">
        <el-input
          v-model="searchKeyword"
          placeholder="搜索教师姓名、工号或邮箱"
          :prefix-icon="Search"
          clearable
          style="width: 300px"
          @input="handleSearch"
        />
      </div>
    </div>

    <!-- 筛选条件 -->
    <div class="filter-section dopamine-card">
      <div class="filter-row">
        <div class="filter-item">
          <label>状态</label>
          <el-select v-model="filters.status" placeholder="选择状态" clearable>
            <el-option label="全部" value="" />
            <el-option label="正常" value="active" />
            <el-option label="禁用" value="disabled" />
          </el-select>
        </div>
        <div class="filter-item">
          <label>学科</label>
          <el-select v-model="filters.subject" placeholder="选择学科" clearable>
            <el-option label="全部" value="" />
            <el-option label="数学" value="math" />
            <el-option label="语文" value="chinese" />
            <el-option label="英语" value="english" />
            <el-option label="物理" value="physics" />
            <el-option label="化学" value="chemistry" />
          </el-select>
        </div>
        <div class="filter-item">
          <label>入职时间</label>
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
        <el-button type="primary" :icon="Search" @click="applyFilters">
          筛选
        </el-button>
        <el-button :icon="Refresh" @click="resetFilters">
          重置
        </el-button>
      </div>
    </div>

    <!-- 教师列表 -->
    <div class="teacher-table dopamine-card">
      <el-table
        :data="filteredTeachers"
        v-loading="loading"
        stripe
        style="width: 100%"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column prop="teacherId" label="工号" min-width="120" />
        <el-table-column label="头像" width="80">
          <template #default="{ row }">
            <el-avatar :src="row.avatar" :size="40">
              <el-icon><User /></el-icon>
            </el-avatar>
          </template>
        </el-table-column>
        <el-table-column prop="realName" label="姓名" min-width="120" />
        <el-table-column prop="email" label="邮箱" min-width="200" />
        <el-table-column prop="phone" label="手机号" min-width="150" />
        <el-table-column label="学科" min-width="150">
          <template #default="{ row }">
            <el-tag v-for="subject in row.subjects" :key="subject" size="small" style="margin-right: 4px">
              {{ getSubjectText(subject) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="状态" min-width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'danger'">
              {{ row.status === 'active' ? '正常' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="joinDate" label="入职时间" min-width="140" />
        <el-table-column label="操作" width="300" fixed="right">
          <template #default="{ row }">
            <div class="action-buttons">
              <el-button 
                type="primary" 
                size="small" 
                :icon="Edit" 
                @click="editTeacher(row)"
              >
                编辑
              </el-button>
              <el-button
                :type="row.status === 'active' ? 'warning' : 'success'"
                size="small"
                :icon="row.status === 'active' ? CircleClose : CircleCheck"
                @click="toggleStatus(row)"
              >
                {{ row.status === 'active' ? '禁用' : '启用' }}
              </el-button>
              <el-button 
                type="danger" 
                size="small" 
                :icon="Delete" 
                @click="deleteTeacher(row)"
              >
                删除
              </el-button>
            </div>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="totalTeachers"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </div>

    <!-- 添加/编辑教师对话框 -->
    <el-dialog
      v-model="showAddDialog"
      :title="editingTeacher ? '编辑教师' : '添加教师'"
      width="600px"
    >
      <el-form
        ref="teacherFormRef"
        :model="teacherForm"
        :rules="teacherRules"
        label-width="80px"
      >
        <el-form-item label="工号" prop="teacherId">
          <el-input v-model="teacherForm.teacherId" placeholder="请输入教师工号" />
        </el-form-item>
        <el-form-item label="姓名" prop="realName">
          <el-input v-model="teacherForm.realName" placeholder="请输入教师姓名" />
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="teacherForm.email" placeholder="请输入邮箱地址" />
        </el-form-item>
        <el-form-item label="手机号" prop="phone">
          <el-input v-model="teacherForm.phone" placeholder="请输入手机号" />
        </el-form-item>
        <el-form-item label="学科" prop="subjects">
          <el-select v-model="teacherForm.subjects" multiple placeholder="选择任教学科">
            <el-option label="数学" value="math" />
            <el-option label="语文" value="chinese" />
            <el-option label="英语" value="english" />
            <el-option label="物理" value="physics" />
            <el-option label="化学" value="chemistry" />
          </el-select>
        </el-form-item>
        <el-form-item label="密码" prop="password" v-if="!editingTeacher">
          <el-input v-model="teacherForm.password" type="password" placeholder="请输入初始密码" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showAddDialog = false">取消</el-button>
        <el-button type="primary" @click="saveTeacher">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.action-buttons {
  display: flex;
  gap: 6px;
  justify-content: center;
  align-items: center;
  width: 100%;
}

.action-buttons .el-button {
  padding: 3px 6px;
  min-width: 50px;
  font-size: 11px;
  margin: 0;
  flex-shrink: 0;
}

.action-buttons .el-button + .el-button {
  margin-left: 0;
}
</style>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Plus,
  Upload,
  Download,
  Search,
  Refresh,
  Edit,
  Delete,
  User,
  CircleClose,
  CircleCheck
} from '@element-plus/icons-vue'


// 响应式数据
const loading = ref(false)
const showAddDialog = ref(false)
const showImportDialog = ref(false)
const editingTeacher = ref(null)
const searchKeyword = ref('')
const currentPage = ref(1)
const pageSize = ref(20)
const totalTeachers = ref(0)

// 筛选条件
const filters = reactive({
  status: '',
  subject: '',
  dateRange: []
})

// 教师表单
const teacherForm = reactive({
  teacherId: '',
  realName: '',
  email: '',
  phone: '',
  subjects: [],
  password: ''
})

// 表单验证规则
const teacherRules = {
  teacherId: [{ required: true, message: '请输入教师工号', trigger: 'blur' }],
  realName: [{ required: true, message: '请输入教师姓名', trigger: 'blur' }],
  email: [
    { required: true, message: '请输入邮箱地址', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ],
  phone: [{ required: true, message: '请输入手机号', trigger: 'blur' }],
  subjects: [{ required: true, message: '请选择任教学科', trigger: 'change' }],
  password: [{ required: true, message: '请输入初始密码', trigger: 'blur' }]
}

// 模拟教师数据
const teachers = ref([
  {
    id: 1,
    teacherId: 'T001',
    realName: '张老师',
    email: 'zhang@example.com',
    phone: '13800138001',
    subjects: ['math', 'physics'],
    status: 'active',
    joinDate: '2023-01-15',
    avatar: ''
  },
  {
    id: 2,
    teacherId: 'T002',
    realName: '李老师',
    email: 'li@example.com',
    phone: '13800138002',
    subjects: ['chinese'],
    status: 'active',
    joinDate: '2023-02-20',
    avatar: ''
  }
])

// 计算属性
const filteredTeachers = computed(() => {
  let result = teachers.value
  
  if (searchKeyword.value) {
    const keyword = searchKeyword.value.toLowerCase()
    result = result.filter(teacher => 
      teacher.realName.toLowerCase().includes(keyword) ||
      teacher.teacherId.toLowerCase().includes(keyword) ||
      teacher.email.toLowerCase().includes(keyword)
    )
  }
  
  return result
})

// 方法
const getSubjectText = (subject: string) => {
  const subjectMap = {
    math: '数学',
    chinese: '语文',
    english: '英语',
    physics: '物理',
    chemistry: '化学'
  }
  return subjectMap[subject] || subject
}

const handleSearch = () => {
  // 搜索逻辑
}

const applyFilters = () => {
  // 应用筛选
}

const resetFilters = () => {
  Object.assign(filters, {
    status: '',
    subject: '',
    dateRange: []
  })
}

const editTeacher = (teacher: any) => {
  editingTeacher.value = teacher
  Object.assign(teacherForm, teacher)
  showAddDialog.value = true
}

const toggleStatus = async (teacher: any) => {
  try {
    await ElMessageBox.confirm(
      `确定要${teacher.status === 'active' ? '禁用' : '启用'}该教师吗？`,
      '确认操作',
      { type: 'warning' }
    )
    teacher.status = teacher.status === 'active' ? 'disabled' : 'active'
    ElMessage.success('操作成功')
  } catch {
    // 用户取消
  }
}

const deleteTeacher = async (teacher: any) => {
  try {
    await ElMessageBox.confirm(
      '确定要删除该教师吗？此操作不可恢复！',
      '确认删除',
      { type: 'warning' }
    )
    const index = teachers.value.findIndex(t => t.id === teacher.id)
    if (index > -1) {
      teachers.value.splice(index, 1)
      ElMessage.success('删除成功')
    }
  } catch {
    // 用户取消
  }
}

const saveTeacher = () => {
  // 保存教师信息
  ElMessage.success('保存成功')
  showAddDialog.value = false
}

const exportTeachers = () => {
  ElMessage.success('导出功能开发中')
}

const handleSizeChange = (size: number) => {
  pageSize.value = size
}

const handleCurrentChange = (page: number) => {
  currentPage.value = page
}

onMounted(() => {
  totalTeachers.value = teachers.value.length
})
</script>

<style scoped>
.teacher-management-view {
  padding: var(--spacing-md);
}

.page-header {
  margin-bottom: var(--spacing-md);
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

.action-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--spacing-md);
  margin-bottom: var(--spacing-md);
}

.action-left {
  display: flex;
  gap: var(--spacing-sm);
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
  font-size: 14px;
}

.teacher-table {
  padding: var(--spacing-md);
}

.pagination-wrapper {
  display: flex;
  justify-content: center;
  margin-top: var(--spacing-md);
}

/* 表格紧凑化样式 */
:deep(.el-table) {
  --el-table-row-height: 44px;
}

:deep(.el-table .el-table__cell) {
  padding: 6px 10px;
}

:deep(.el-avatar) {
  width: 30px;
  height: 30px;
}

:deep(.el-tag--small) {
  padding: 0 5px;
  height: 18px;
  line-height: 16px;
  font-size: 10px;
}

/* 筛选区域紧凑化 */
:deep(.el-select) {
  width: 110px;
}

:deep(.el-date-editor) {
  width: 220px;
}
</style>