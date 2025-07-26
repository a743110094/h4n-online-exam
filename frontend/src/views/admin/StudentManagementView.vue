<template>
  <div class="student-management-view">
    <!-- 页面标题 -->
    <div class="page-header">
      <h1 class="page-title">学生管理</h1>
      <p class="page-subtitle">管理系统中的学生账户</p>
    </div>

    <!-- 操作栏 -->
    <div class="action-bar dopamine-card">
      <div class="action-left">
        <el-button type="primary" :icon="Plus" @click="showAddDialog = true">
          添加学生
        </el-button>
        <el-button :icon="Upload" @click="showImportDialog = true">
          批量导入
        </el-button>
        <el-button :icon="Download" @click="exportStudents">
          导出数据
        </el-button>
      </div>
      <div class="action-right">
        <el-input
          v-model="searchKeyword"
          placeholder="搜索学生姓名、学号或邮箱"
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
          <label>年级</label>
          <el-select v-model="filters.grade" placeholder="选择年级" clearable>
            <el-option label="全部" value="" />
            <el-option label="高一" value="grade1" />
            <el-option label="高二" value="grade2" />
            <el-option label="高三" value="grade3" />
          </el-select>
        </div>
        <div class="filter-item">
          <label>班级</label>
          <el-select v-model="filters.class" placeholder="选择班级" clearable>
            <el-option label="全部" value="" />
            <el-option label="1班" value="class1" />
            <el-option label="2班" value="class2" />
            <el-option label="3班" value="class3" />
            <el-option label="4班" value="class4" />
          </el-select>
        </div>
        <div class="filter-item">
          <label>入学时间</label>
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

    <!-- 学生列表 -->
    <div class="student-table dopamine-card">
      <el-table
        :data="filteredStudents"
        v-loading="loading"
        stripe
        style="width: 100%"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column prop="studentId" label="学号" min-width="120" />
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
        <el-table-column label="年级班级" min-width="120">
          <template #default="{ row }">
            <el-tag size="small">
              {{ getGradeText(row.grade) }}{{ getClassText(row.class) }}
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
        <el-table-column prop="enrollDate" label="入学时间" min-width="140" />
        <el-table-column label="操作" width="300" fixed="right">
          <template #default="{ row }">
            <div class="action-buttons">
              <el-button 
                type="primary" 
                size="small" 
                :icon="Edit" 
                @click="editStudent(row)"
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
                @click="deleteStudent(row)"
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
          :total="totalStudents"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </div>

    <!-- 添加/编辑学生对话框 -->
    <el-dialog
      v-model="showAddDialog"
      :title="editingStudent ? '编辑学生' : '添加学生'"
      width="600px"
    >
      <el-form
        ref="studentFormRef"
        :model="studentForm"
        :rules="studentRules"
        label-width="80px"
      >
        <el-form-item label="学号" prop="studentId">
          <el-input v-model="studentForm.studentId" placeholder="请输入学生学号" />
        </el-form-item>
        <el-form-item label="姓名" prop="realName">
          <el-input v-model="studentForm.realName" placeholder="请输入学生姓名" />
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="studentForm.email" placeholder="请输入邮箱地址" />
        </el-form-item>
        <el-form-item label="手机号" prop="phone">
          <el-input v-model="studentForm.phone" placeholder="请输入手机号" />
        </el-form-item>
        <el-form-item label="年级" prop="grade">
          <el-select v-model="studentForm.grade" placeholder="选择年级">
            <el-option label="高一" value="grade1" />
            <el-option label="高二" value="grade2" />
            <el-option label="高三" value="grade3" />
          </el-select>
        </el-form-item>
        <el-form-item label="班级" prop="class">
          <el-select v-model="studentForm.class" placeholder="选择班级">
            <el-option label="1班" value="class1" />
            <el-option label="2班" value="class2" />
            <el-option label="3班" value="class3" />
            <el-option label="4班" value="class4" />
          </el-select>
        </el-form-item>
        <el-form-item label="密码" prop="password" v-if="!editingStudent">
          <el-input v-model="studentForm.password" type="password" placeholder="请输入初始密码" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showAddDialog = false">取消</el-button>
        <el-button type="primary" @click="saveStudent">确定</el-button>
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
const editingStudent = ref(null)
const searchKeyword = ref('')
const currentPage = ref(1)
const pageSize = ref(20)
const totalStudents = ref(0)

// 筛选条件
const filters = reactive({
  status: '',
  grade: '',
  class: '',
  dateRange: []
})

// 学生表单
const studentForm = reactive({
  studentId: '',
  realName: '',
  email: '',
  phone: '',
  grade: '',
  class: '',
  password: ''
})

// 表单验证规则
const studentRules = {
  studentId: [{ required: true, message: '请输入学生学号', trigger: 'blur' }],
  realName: [{ required: true, message: '请输入学生姓名', trigger: 'blur' }],
  email: [
    { required: true, message: '请输入邮箱地址', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ],
  phone: [{ required: true, message: '请输入手机号', trigger: 'blur' }],
  grade: [{ required: true, message: '请选择年级', trigger: 'change' }],
  class: [{ required: true, message: '请选择班级', trigger: 'change' }],
  password: [{ required: true, message: '请输入初始密码', trigger: 'blur' }]
}

// 模拟学生数据
const students = ref([
  {
    id: 1,
    studentId: 'S2023001',
    realName: '王小明',
    email: 'wang@example.com',
    phone: '13800138001',
    grade: 'grade1',
    class: 'class1',
    status: 'active',
    enrollDate: '2023-09-01',
    avatar: ''
  },
  {
    id: 2,
    studentId: 'S2023002',
    realName: '李小红',
    email: 'li@example.com',
    phone: '13800138002',
    grade: 'grade2',
    class: 'class2',
    status: 'active',
    enrollDate: '2022-09-01',
    avatar: ''
  }
])

// 计算属性
const filteredStudents = computed(() => {
  let result = students.value
  
  if (searchKeyword.value) {
    const keyword = searchKeyword.value.toLowerCase()
    result = result.filter(student => 
      student.realName.toLowerCase().includes(keyword) ||
      student.studentId.toLowerCase().includes(keyword) ||
      student.email.toLowerCase().includes(keyword)
    )
  }
  
  return result
})

// 方法
const getGradeText = (grade: string) => {
  const gradeMap = {
    grade1: '高一',
    grade2: '高二',
    grade3: '高三'
  }
  return gradeMap[grade] || grade
}

const getClassText = (classValue: string) => {
  const classMap = {
    class1: '1班',
    class2: '2班',
    class3: '3班',
    class4: '4班'
  }
  return classMap[classValue] || classValue
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
    grade: '',
    class: '',
    dateRange: []
  })
}

const editStudent = (student: any) => {
  editingStudent.value = student
  Object.assign(studentForm, student)
  showAddDialog.value = true
}

const toggleStatus = async (student: any) => {
  try {
    await ElMessageBox.confirm(
      `确定要${student.status === 'active' ? '禁用' : '启用'}该学生吗？`,
      '确认操作',
      { type: 'warning' }
    )
    student.status = student.status === 'active' ? 'disabled' : 'active'
    ElMessage.success('操作成功')
  } catch {
    // 用户取消
  }
}

const deleteStudent = async (student: any) => {
  try {
    await ElMessageBox.confirm(
      '确定要删除该学生吗？此操作不可恢复！',
      '确认删除',
      { type: 'warning' }
    )
    const index = students.value.findIndex(s => s.id === student.id)
    if (index > -1) {
      students.value.splice(index, 1)
      ElMessage.success('删除成功')
    }
  } catch {
    // 用户取消
  }
}

const saveStudent = () => {
  // 保存学生信息
  ElMessage.success('保存成功')
  showAddDialog.value = false
}

const exportStudents = () => {
  ElMessage.success('导出功能开发中')
}

const handleSizeChange = (size: number) => {
  pageSize.value = size
}

const handleCurrentChange = (page: number) => {
  currentPage.value = page
}

onMounted(() => {
  totalStudents.value = students.value.length
})
</script>

<style scoped>
.student-management-view {
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

.student-table {
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