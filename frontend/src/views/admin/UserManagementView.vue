<template>
  <div class="user-management-view">
    <!-- 页面标题 -->
    <div class="page-header">
      <h1 class="page-title">用户管理</h1>
      <p class="page-subtitle">管理系统用户，包括学生、教师和管理员</p>
    </div>
    
    <!-- 操作工具栏 -->
    <div class="toolbar dopamine-card">
      <div class="toolbar-left">
        <el-button type="primary" @click="showAddUserDialog = true">
          <el-icon><Plus /></el-icon>
          添加用户
        </el-button>
        <el-button type="success" @click="showImportDialog = true">
          <el-icon><Upload /></el-icon>
          批量导入
        </el-button>
        <el-button type="warning" @click="exportUsers">
          <el-icon><Download /></el-icon>
          导出用户
        </el-button>
        <el-button
          type="danger"
          :disabled="selectedUsers.length === 0"
          @click="batchDeleteUsers"
        >
          <el-icon><Delete /></el-icon>
          批量删除
        </el-button>
      </div>
      
      <div class="toolbar-right">
        <el-input
          v-model="searchKeyword"
          placeholder="搜索用户名、姓名、邮箱..."
          style="width: 300px"
          clearable
          @input="handleSearch"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
      </div>
    </div>
    
    <!-- 筛选条件 -->
    <div class="filter-section dopamine-card">
      <div class="filter-row">
        <div class="filter-item">
          <label>用户角色</label>
          <el-select v-model="filters.role" placeholder="选择角色" clearable>
            <el-option label="管理员" value="admin" />
            <el-option label="教师" value="teacher" />
            <el-option label="学生" value="student" />
          </el-select>
        </div>
        
        <div class="filter-item">
          <label>账户状态</label>
          <el-select v-model="filters.status" placeholder="选择状态" clearable>
            <el-option label="正常" value="active" />
            <el-option label="禁用" value="disabled" />
            <el-option label="待激活" value="pending" />
          </el-select>
        </div>
        
        <div class="filter-item">
          <label>注册时间</label>
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
        
        <div class="filter-actions">
          <el-button type="primary" @click="applyFilters">
            <el-icon><Search /></el-icon>
            筛选
          </el-button>
          <el-button @click="resetFilters">
            <el-icon><Refresh /></el-icon>
            重置
          </el-button>
        </div>
      </div>
    </div>
    
    <!-- 统计信息 -->
    <div class="stats-section">
      <div class="stats-grid">
        <div class="stat-card dopamine-card">
          <div class="stat-icon" style="background: var(--gradient-primary)">
            <el-icon :size="24"><User /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ userStats.total }}</div>
            <div class="stat-label">用户总数</div>
          </div>
        </div>
        
        <div class="stat-card dopamine-card">
          <div class="stat-icon" style="background: var(--gradient-success)">
            <el-icon :size="24"><UserFilled /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ userStats.students }}</div>
            <div class="stat-label">学生</div>
          </div>
        </div>
        
        <div class="stat-card dopamine-card">
          <div class="stat-icon" style="background: var(--gradient-warning)">
            <el-icon :size="24"><Avatar /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ userStats.teachers }}</div>
            <div class="stat-label">教师</div>
          </div>
        </div>
        
        <div class="stat-card dopamine-card">
          <div class="stat-icon" style="background: var(--gradient-danger)">
            <el-icon :size="24"><Star /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ userStats.admins }}</div>
            <div class="stat-label">管理员</div>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 用户列表 -->
    <div class="user-list dopamine-card">
      <el-table
        :data="filteredUsers"
        style="width: 100%"
        @selection-change="handleSelectionChange"
        row-key="id"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column label="用户信息" min-width="200">
          <template #default="{ row }">
            <div class="user-info">
              <el-avatar
                :size="40"
                :src="row.avatar"
                :alt="row.realName"
                class="user-avatar"
              >
                {{ row.realName?.charAt(0) || row.username.charAt(0) }}
              </el-avatar>
              <div class="user-details">
                <div class="user-name">{{ row.realName || row.username }}</div>
                <div class="user-username">@{{ row.username }}</div>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="email" label="邮箱" min-width="200" />
        <el-table-column prop="phone" label="手机号" min-width="130" />
        <el-table-column prop="role" label="角色" min-width="100">
          <template #default="{ row }">
            <el-tag :type="getRoleTagType(row.role)">
              {{ getRoleText(row.role) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" min-width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusTagType(row.status)">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="lastLoginAt" label="最后登录" min-width="150">
          <template #default="{ row }">
            <span v-if="row.lastLoginAt">
              {{ formatDateTime(row.lastLoginAt) }}
            </span>
            <span v-else class="text-muted">从未登录</span>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="注册时间" min-width="120">
          <template #default="{ row }">
            {{ formatDate(row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button
              type="primary"
              size="small"
              @click="editUser(row)"
            >
              编辑
            </el-button>
            <el-button
              v-if="row.status === 'active'"
              type="warning"
              size="small"
              @click="toggleUserStatus(row, 'disabled')"
            >
              禁用
            </el-button>
            <el-button
              v-else
              type="success"
              size="small"
              @click="toggleUserStatus(row, 'active')"
            >
              启用
            </el-button>
            <el-button
              type="danger"
              size="small"
              @click="deleteUser(row.id)"
            >
              删除
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
          :total="totalUsers"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </div>
    
    <!-- 添加/编辑用户对话框 -->
    <el-dialog
      v-model="showAddUserDialog"
      :title="currentUser?.id ? '编辑用户' : '添加用户'"
      width="600px"
      :before-close="handleCloseUserDialog"
    >
      <el-form
        ref="userFormRef"
        :model="userForm"
        :rules="userFormRules"
        label-width="80px"
      >
        <el-form-item label="用户名" prop="username">
          <el-input
            v-model="userForm.username"
            placeholder="请输入用户名"
            :disabled="!!currentUser?.id"
          />
        </el-form-item>
        
        <el-form-item label="真实姓名" prop="realName">
          <el-input v-model="userForm.realName" placeholder="请输入真实姓名" />
        </el-form-item>
        
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="userForm.email" placeholder="请输入邮箱" />
        </el-form-item>
        
        <el-form-item label="手机号" prop="phone">
          <el-input v-model="userForm.phone" placeholder="请输入手机号" />
        </el-form-item>
        
        <el-form-item label="角色" prop="role">
          <el-select v-model="userForm.role" placeholder="请选择角色">
            <el-option label="管理员" value="admin" />
            <el-option label="教师" value="teacher" />
            <el-option label="学生" value="student" />
          </el-select>
        </el-form-item>
        
        <el-form-item v-if="!currentUser?.id" label="密码" prop="password">
          <el-input
            v-model="userForm.password"
            type="password"
            placeholder="请输入密码"
            show-password
          />
        </el-form-item>
        
        <el-form-item label="状态" prop="status">
          <el-select v-model="userForm.status" placeholder="请选择状态">
            <el-option label="正常" value="active" />
            <el-option label="禁用" value="disabled" />
            <el-option label="待激活" value="pending" />
          </el-select>
        </el-form-item>
      </el-form>
      
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="showAddUserDialog = false">取消</el-button>
          <el-button
            type="primary"
            :loading="saving"
            @click="saveUser"
          >
            保存
          </el-button>
        </div>
      </template>
    </el-dialog>
    
    <!-- 批量导入对话框 -->
    <el-dialog
      v-model="showImportDialog"
      title="批量导入用户"
      width="600px"
    >
      <div class="import-section">
        <div class="import-tips">
          <el-alert
            title="导入说明"
            type="info"
            :closable="false"
            show-icon
          >
            <template #default>
              <ul>
                <li>支持 Excel (.xlsx) 和 CSV (.csv) 格式</li>
                <li>请按照模板格式准备数据</li>
                <li>单次最多导入 1000 个用户</li>
                <li>导入后用户默认密码为：123456</li>
              </ul>
            </template>
          </el-alert>
        </div>
        
        <div class="import-actions">
          <el-button type="primary" @click="downloadTemplate">
            <el-icon><Download /></el-icon>
            下载模板
          </el-button>
        </div>
        
        <el-upload
          class="upload-demo"
          drag
          :auto-upload="false"
          :on-change="handleFileChange"
          :accept="'.xlsx,.csv'"
          :limit="1"
        >
          <el-icon class="el-icon--upload"><UploadFilled /></el-icon>
          <div class="el-upload__text">
            将文件拖到此处，或<em>点击上传</em>
          </div>
          <template #tip>
            <div class="el-upload__tip">
              只能上传 xlsx/csv 文件，且不超过 10MB
            </div>
          </template>
        </el-upload>
        
        <div class="import-progress" v-if="importProgress.show">
          <el-progress
            :percentage="importProgress.percentage"
            :status="importProgress.status"
          />
          <div class="progress-text">{{ importProgress.text }}</div>
        </div>
      </div>
      
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="showImportDialog = false">取消</el-button>
          <el-button
            type="primary"
            :loading="importing"
            :disabled="!selectedFile"
            @click="startImport"
          >
            开始导入
          </el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Plus, Upload, Download, Delete, Search, Refresh,
  User, UserFilled, Avatar, Star, UploadFilled
} from '@element-plus/icons-vue'

// 搜索关键词
const searchKeyword = ref('')

// 筛选条件
const filters = reactive({
  role: '',
  status: '',
  dateRange: null
})

// 用户统计
const userStats = ref({
  total: 2847,
  students: 2456,
  teachers: 378,
  admins: 13
})

// 选中的用户
const selectedUsers = ref([])

// 分页
const currentPage = ref(1)
const pageSize = ref(20)
const totalUsers = ref(2847)

// 对话框状态
const showAddUserDialog = ref(false)
const showImportDialog = ref(false)

// 当前编辑的用户
const currentUser = ref(null)
const saving = ref(false)

// 用户表单
const userForm = reactive({
  username: '',
  realName: '',
  email: '',
  phone: '',
  role: '',
  password: '',
  status: 'active'
})

// 表单验证规则
const userFormRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '用户名长度在 3 到 20 个字符', trigger: 'blur' }
  ],
  realName: [
    { required: true, message: '请输入真实姓名', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ],
  phone: [
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号', trigger: 'blur' }
  ],
  role: [
    { required: true, message: '请选择角色', trigger: 'change' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于 6 位', trigger: 'blur' }
  ]
}

const userFormRef = ref()

// 导入相关
const importing = ref(false)
const selectedFile = ref(null)
const importProgress = reactive({
  show: false,
  percentage: 0,
  status: '',
  text: ''
})

// 模拟用户数据
const users = ref([
  {
    id: 1,
    username: 'admin',
    realName: '系统管理员',
    email: 'admin@example.com',
    phone: '13800138000',
    role: 'admin',
    status: 'active',
    avatar: '',
    lastLoginAt: new Date(Date.now() - 2 * 60 * 60 * 1000),
    createdAt: new Date(Date.now() - 365 * 24 * 60 * 60 * 1000)
  },
  {
    id: 2,
    username: 'teacher001',
    realName: '张教授',
    email: 'zhang@example.com',
    phone: '13800138001',
    role: 'teacher',
    status: 'active',
    avatar: '',
    lastLoginAt: new Date(Date.now() - 30 * 60 * 1000),
    createdAt: new Date(Date.now() - 180 * 24 * 60 * 60 * 1000)
  },
  {
    id: 3,
    username: 'student001',
    realName: '李小明',
    email: 'lixiaoming@example.com',
    phone: '13800138002',
    role: 'student',
    status: 'active',
    avatar: '',
    lastLoginAt: new Date(Date.now() - 5 * 60 * 1000),
    createdAt: new Date(Date.now() - 90 * 24 * 60 * 60 * 1000)
  },
  {
    id: 4,
    username: 'student002',
    realName: '王小红',
    email: 'wangxiaohong@example.com',
    phone: '13800138003',
    role: 'student',
    status: 'disabled',
    avatar: '',
    lastLoginAt: new Date(Date.now() - 7 * 24 * 60 * 60 * 1000),
    createdAt: new Date(Date.now() - 60 * 24 * 60 * 60 * 1000)
  },
  {
    id: 5,
    username: 'student003',
    realName: '刘小强',
    email: 'liuxiaoqiang@example.com',
    phone: '',
    role: 'student',
    status: 'pending',
    avatar: '',
    lastLoginAt: null,
    createdAt: new Date(Date.now() - 1 * 24 * 60 * 60 * 1000)
  }
])

// 筛选后的用户列表
const filteredUsers = computed(() => {
  let result = users.value
  
  // 搜索关键词筛选
  if (searchKeyword.value) {
    const keyword = searchKeyword.value.toLowerCase()
    result = result.filter(user => 
      user.username.toLowerCase().includes(keyword) ||
      user.realName?.toLowerCase().includes(keyword) ||
      user.email.toLowerCase().includes(keyword)
    )
  }
  
  // 角色筛选
  if (filters.role) {
    result = result.filter(user => user.role === filters.role)
  }
  
  // 状态筛选
  if (filters.status) {
    result = result.filter(user => user.status === filters.status)
  }
  
  // 日期范围筛选
  if (filters.dateRange && filters.dateRange.length === 2) {
    const [startDate, endDate] = filters.dateRange
    result = result.filter(user => {
      const userDate = user.createdAt.toISOString().split('T')[0]
      return userDate >= startDate && userDate <= endDate
    })
  }
  
  return result
})

// 获取角色文本
const getRoleText = (role: string): string => {
  const textMap: Record<string, string> = {
    admin: '管理员',
    teacher: '教师',
    student: '学生'
  }
  return textMap[role] || ''
}

// 获取角色标签类型
const getRoleTagType = (role: string): string => {
  const typeMap: Record<string, string> = {
    admin: 'danger',
    teacher: 'warning',
    student: 'success'
  }
  return typeMap[role] || ''
}

// 获取状态文本
const getStatusText = (status: string): string => {
  const textMap: Record<string, string> = {
    active: '正常',
    disabled: '禁用',
    pending: '待激活'
  }
  return textMap[status] || ''
}

// 获取状态标签类型
const getStatusTagType = (status: string): string => {
  const typeMap: Record<string, string> = {
    active: 'success',
    disabled: 'danger',
    pending: 'warning'
  }
  return typeMap[status] || ''
}

// 格式化日期
const formatDate = (date: Date): string => {
  return date.toLocaleDateString('zh-CN')
}

// 格式化日期时间
const formatDateTime = (date: Date): string => {
  return date.toLocaleString('zh-CN', {
    month: 'numeric',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// 处理搜索
const handleSearch = () => {
  // 搜索逻辑已在计算属性中处理
}

// 应用筛选
const applyFilters = () => {
  ElMessage.success('筛选条件已应用')
}

// 重置筛选
const resetFilters = () => {
  filters.role = ''
  filters.status = ''
  filters.dateRange = null
  searchKeyword.value = ''
  ElMessage.success('筛选条件已重置')
}

// 处理选择变化
const handleSelectionChange = (selection: any[]) => {
  selectedUsers.value = selection
}

// 编辑用户
const editUser = (user: any) => {
  currentUser.value = user
  Object.assign(userForm, {
    username: user.username,
    realName: user.realName,
    email: user.email,
    phone: user.phone,
    role: user.role,
    password: '',
    status: user.status
  })
  showAddUserDialog.value = true
}

// 切换用户状态
const toggleUserStatus = async (user: any, newStatus: string) => {
  const action = newStatus === 'active' ? '启用' : '禁用'
  try {
    await ElMessageBox.confirm(
      `确定要${action}用户「${user.realName || user.username}」吗？`,
      `${action}用户`,
      {
        confirmButtonText: `确定${action}`,
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    user.status = newStatus
    ElMessage.success(`用户${action}成功`)
  } catch {
    // 用户取消
  }
}

// 删除用户
const deleteUser = async (userId: number) => {
  try {
    await ElMessageBox.confirm(
      '确定要删除这个用户吗？删除后无法恢复。',
      '删除确认',
      {
        confirmButtonText: '确定删除',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    ElMessage.success('用户删除成功')
  } catch {
    // 用户取消
  }
}

// 批量删除用户
const batchDeleteUsers = async () => {
  try {
    await ElMessageBox.confirm(
      `确定要删除选中的 ${selectedUsers.value.length} 个用户吗？删除后无法恢复。`,
      '批量删除确认',
      {
        confirmButtonText: '确定删除',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    ElMessage.success(`成功删除 ${selectedUsers.value.length} 个用户`)
    selectedUsers.value = []
  } catch {
    // 用户取消
  }
}

// 导出用户
const exportUsers = () => {
  ElMessage.success('用户导出功能开发中')
}

// 保存用户
const saveUser = async () => {
  try {
    await userFormRef.value.validate()
    
    saving.value = true
    
    // 模拟保存过程
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    if (currentUser.value?.id) {
      ElMessage.success('用户更新成功')
    } else {
      ElMessage.success('用户添加成功')
    }
    
    showAddUserDialog.value = false
    resetUserForm()
  } catch (error) {
    console.error('表单验证失败:', error)
  } finally {
    saving.value = false
  }
}

// 重置用户表单
const resetUserForm = () => {
  Object.assign(userForm, {
    username: '',
    realName: '',
    email: '',
    phone: '',
    role: '',
    password: '',
    status: 'active'
  })
  currentUser.value = null
  userFormRef.value?.clearValidate()
}

// 处理关闭用户对话框
const handleCloseUserDialog = () => {
  resetUserForm()
  showAddUserDialog.value = false
}

// 下载模板
const downloadTemplate = () => {
  ElMessage.success('模板下载功能开发中')
}

// 处理文件变化
const handleFileChange = (file: any) => {
  selectedFile.value = file
}

// 开始导入
const startImport = () => {
  if (!selectedFile.value) {
    ElMessage.warning('请选择要导入的文件')
    return
  }
  
  importing.value = true
  importProgress.show = true
  importProgress.percentage = 0
  importProgress.status = 'active'
  importProgress.text = '正在解析文件...'
  
  // 模拟导入过程
  const timer = setInterval(() => {
    importProgress.percentage += 10
    
    if (importProgress.percentage === 30) {
      importProgress.text = '正在验证数据格式...'
    } else if (importProgress.percentage === 60) {
      importProgress.text = '正在创建用户账户...'
    } else if (importProgress.percentage === 90) {
      importProgress.text = '正在发送激活邮件...'
    } else if (importProgress.percentage >= 100) {
      clearInterval(timer)
      importProgress.status = 'success'
      importProgress.text = '导入完成！'
      importing.value = false
      
      setTimeout(() => {
        showImportDialog.value = false
        importProgress.show = false
        selectedFile.value = null
        ElMessage.success('用户导入成功')
      }, 1000)
    }
  }, 300)
}

// 分页处理
const handleSizeChange = (size: number) => {
  pageSize.value = size
}

const handleCurrentChange = (page: number) => {
  currentPage.value = page
}

onMounted(() => {
  // 初始化数据
})
</script>

<style scoped>
.user-management-view {
  padding: var(--spacing-md);
  max-width: 1400px;
  margin: 0 auto;
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

.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--spacing-md);
  margin-bottom: var(--spacing-md);
}

.toolbar-left {
  display: flex;
  gap: var(--spacing-md);
}

.filter-section {
  padding: var(--spacing-md);
  margin-bottom: var(--spacing-md);
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
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: var(--spacing-lg);
}

.stat-card {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
  padding: var(--spacing-lg);
  transition: all 0.3s ease;
}

.stat-card:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow-xl);
}

.stat-icon {
  width: 60px;
  height: 60px;
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
  margin-bottom: 4px;
}

.stat-label {
  font-size: 14px;
  color: var(--text-secondary);
}

.user-list {
  padding: var(--spacing-md);
}

.user-info {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
}

.user-avatar {
  flex-shrink: 0;
}

.user-details {
  flex: 1;
}

.user-name {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 2px;
}

.user-username {
  font-size: 12px;
  color: var(--text-secondary);
}

.text-muted {
  color: var(--text-secondary);
  font-style: italic;
}

.pagination-wrapper {
  display: flex;
  justify-content: center;
  margin-top: var(--spacing-xl);
}

.import-section {
  padding: var(--spacing-lg) 0;
}

.import-tips {
  margin-bottom: var(--spacing-lg);
}

.import-tips ul {
  margin: 0;
  padding-left: var(--spacing-lg);
}

.import-tips li {
  margin-bottom: var(--spacing-sm);
}

.import-actions {
  margin-bottom: var(--spacing-lg);
  text-align: center;
}

.upload-demo {
  margin-bottom: var(--spacing-lg);
}

.import-progress {
  margin-top: var(--spacing-lg);
}

.progress-text {
  text-align: center;
  margin-top: var(--spacing-sm);
  font-size: 14px;
  color: var(--text-secondary);
}

.dialog-footer {
  text-align: right;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .toolbar {
    flex-direction: column;
    gap: var(--spacing-md);
    align-items: stretch;
  }
  
  .toolbar-left {
    flex-wrap: wrap;
  }
  
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
}

@media (max-width: 480px) {
  .stats-grid {
    grid-template-columns: 1fr;
  }
  
  .toolbar-left {
    flex-direction: column;
  }
}
</style>