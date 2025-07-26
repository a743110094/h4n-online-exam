<template>
  <div class="main-layout">
    <!-- 侧边栏 -->
    <el-aside
      :width="isCollapsed ? '64px' : '240px'"
      class="sidebar"
    >
      <div class="sidebar-header">
        <div class="logo">
          <el-icon class="logo-icon"><School /></el-icon>
          <span v-show="!isCollapsed" class="logo-text">考试系统</span>
        </div>
      </div>
      
      <el-scrollbar class="sidebar-menu">
        <el-menu
          :default-active="activeMenu"
          :collapse="isCollapsed"
          :unique-opened="true"
          router
          class="dopamine-menu"
        >
          <!-- 管理员菜单 -->
          <template v-if="authStore.isAdmin">
            <el-menu-item index="/admin/dashboard">
              <el-icon><DataBoard /></el-icon>
              <span>仪表盘</span>
            </el-menu-item>
            
            <el-sub-menu index="user-management">
              <template #title>
                <el-icon><UserFilled /></el-icon>
                <span>用户管理</span>
              </template>
              <el-menu-item index="/admin/users/teachers">
                <el-icon><Avatar /></el-icon>
                <span>教师管理</span>
              </el-menu-item>
              <el-menu-item index="/admin/users/students">
                <el-icon><User /></el-icon>
                <span>学生管理</span>
              </el-menu-item>
            </el-sub-menu>
            
            <el-menu-item index="/admin/statistics">
              <el-icon><TrendCharts /></el-icon>
              <span>数据统计</span>
            </el-menu-item>
            
            <el-menu-item index="/admin/settings">
              <el-icon><Setting /></el-icon>
              <span>系统设置</span>
            </el-menu-item>
          </template>
          
          <!-- 教师菜单 -->
          <template v-if="authStore.isTeacher">
            <el-menu-item index="/teacher/dashboard">
              <el-icon><DataBoard /></el-icon>
              <span>工作台</span>
            </el-menu-item>
            
            <el-sub-menu index="question-management">
              <template #title>
                <el-icon><Document /></el-icon>
                <span>题库管理</span>
              </template>
              <el-menu-item index="/teacher/questions">
                <el-icon><EditPen /></el-icon>
                <span>题目管理</span>
              </el-menu-item>
              <el-menu-item index="/teacher/questions/import">
                <el-icon><Upload /></el-icon>
                <span>批量导入</span>
              </el-menu-item>
              <el-menu-item index="/teacher/questions/collect">
                <el-icon><Connection /></el-icon>
                <span>网络采集</span>
              </el-menu-item>
              <el-menu-item index="/teacher/subjects">
                <el-icon><FolderOpened /></el-icon>
                <span>科目分类</span>
              </el-menu-item>
            </el-sub-menu>
            
            <el-sub-menu index="exam-management">
              <template #title>
                <el-icon><Notebook /></el-icon>
                <span>考试管理</span>
              </template>
              <el-menu-item index="/teacher/papers">
                <el-icon><Files /></el-icon>
                <span>组卷管理</span>
              </el-menu-item>
              <el-menu-item index="/teacher/exams">
                <el-icon><Clock /></el-icon>
                <span>考试安排</span>
              </el-menu-item>
              <el-menu-item index="/teacher/monitoring">
                <el-icon><View /></el-icon>
                <span>考试监控</span>
              </el-menu-item>
            </el-sub-menu>
            
            <el-sub-menu index="grade-management">
              <template #title>
                <el-icon><Medal /></el-icon>
                <span>成绩管理</span>
              </template>
              <el-menu-item index="/teacher/grades/entry">
                <el-icon><EditPen /></el-icon>
                <span>成绩录入</span>
              </el-menu-item>
              <el-menu-item index="/teacher/grades/query">
                <el-icon><Search /></el-icon>
                <span>成绩查询</span>
              </el-menu-item>
              <el-menu-item index="/teacher/grades/export">
                <el-icon><Download /></el-icon>
                <span>成绩导出</span>
              </el-menu-item>
            </el-sub-menu>
          </template>
          
          <!-- 学生菜单 -->
          <template v-if="authStore.isStudent">
            <el-menu-item index="/student/dashboard">
              <el-icon><Notebook /></el-icon>
              <span>学习中心</span>
            </el-menu-item>
            
            <el-menu-item index="/student/practice">
              <el-icon><Reading /></el-icon>
              <span>练习刷题</span>
            </el-menu-item>
            
            <el-menu-item index="/student/exams">
              <el-icon><Notebook /></el-icon>
              <span>我的考试</span>
            </el-menu-item>
            
            <el-menu-item index="/student/grades">
              <el-icon><Medal /></el-icon>
              <span>成绩查询</span>
            </el-menu-item>
            
            <el-menu-item index="/student/ai-assistant">
              <el-icon><ChatDotRound /></el-icon>
              <span>AI助手</span>
            </el-menu-item>
          </template>
        </el-menu>
      </el-scrollbar>
      
      <!-- 折叠按钮移到底部 -->
      <div class="sidebar-footer">
        <el-button
          class="collapse-btn"
          :icon="isCollapsed ? Expand : Fold"
          @click="toggleCollapse"
        />
      </div>
    </el-aside>
    
    <!-- 主内容区 -->
    <el-container class="main-container">
      <!-- 顶部导航栏 -->
      <el-header class="header">
        <div class="header-left">
          <el-breadcrumb class="breadcrumb">
            <el-breadcrumb-item
              v-for="item in breadcrumbs"
              :key="item.path"
              :to="item.path"
            >
              {{ item.title }}
            </el-breadcrumb-item>
          </el-breadcrumb>
        </div>
        
        <div class="header-right">
          <!-- 通知 -->
          <el-badge :value="notificationCount" class="notification-badge">
            <el-button
              class="header-btn"
              :icon="Bell"
              circle
              @click="showNotifications"
            />
          </el-badge>
          
          <!-- 用户菜单 -->
          <el-dropdown class="user-dropdown" @command="handleUserCommand">
            <div class="user-info">
              <el-avatar
                :src="authStore.user?.avatar"
                class="user-avatar"
              >
                <el-icon><User /></el-icon>
              </el-avatar>
              <span v-show="!isCollapsed" class="user-name">
                {{ authStore.user?.realName || authStore.user?.username }}
              </span>
              <el-icon v-show="!isCollapsed" class="dropdown-icon"><ArrowDown /></el-icon>
            </div>
            
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="profile">
                  <el-icon><User /></el-icon>
                  个人资料
                </el-dropdown-item>
                <el-dropdown-item command="settings">
                  <el-icon><Setting /></el-icon>
                  账户设置
                </el-dropdown-item>
                <el-dropdown-item divided command="logout">
                  <el-icon><SwitchButton /></el-icon>
                  退出登录
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>
      
      <!-- 主内容 -->
      <el-main class="main-content">
        <div class="content-wrapper">
          <router-view v-slot="{ Component }">
            <transition name="fade-transform" mode="out-in">
              <component :is="Component" />
            </transition>
          </router-view>
        </div>
      </el-main>
    </el-container>
    
    <!-- 通知抽屉 -->
    <el-drawer
      v-model="notificationDrawer"
      title="通知中心"
      direction="rtl"
      size="400px"
    >
      <div class="notification-list">
        <div
          v-for="notification in notifications"
          :key="notification.id"
          class="notification-item dopamine-card"
          :class="{ unread: !notification.read }"
        >
          <div class="notification-header">
            <span class="notification-title">{{ notification.title }}</span>
            <span class="notification-time">{{ formatTime(notification.createdAt) }}</span>
          </div>
          <p class="notification-content">{{ notification.content }}</p>
          <div class="notification-actions">
            <el-button
              v-if="!notification.read"
              size="small"
              type="primary"
              @click="markAsRead(notification.id)"
            >
              标记已读
            </el-button>
          </div>
        </div>
      </div>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  School, Expand, Fold, DataBoard, UserFilled, Avatar, User,
  TrendCharts, Setting, Document, EditPen, Upload, Connection,
  FolderOpened, Notebook, Files, Clock, View, Medal, Reading,
  ChatDotRound, Bell, ArrowDown, SwitchButton, Search, Download
} from '@element-plus/icons-vue'
import { useAuthStore } from '@/stores/auth'
import { debounce } from 'lodash'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

// 侧边栏折叠状态
const isCollapsed = ref(false)

// 通知相关
const notificationDrawer = ref(false)
const notificationCount = ref(3)
const notifications = ref([
  {
    id: 1,
    title: '考试提醒',
    content: '您有一场考试将在30分钟后开始，请做好准备。',
    read: false,
    createdAt: new Date(Date.now() - 10 * 60 * 1000)
  },
  {
    id: 2,
    title: '成绩发布',
    content: '《数据结构》考试成绩已发布，请查看。',
    read: false,
    createdAt: new Date(Date.now() - 2 * 60 * 60 * 1000)
  },
  {
    id: 3,
    title: '系统维护',
    content: '系统将于今晚22:00-24:00进行维护，请合理安排时间。',
    read: true,
    createdAt: new Date(Date.now() - 24 * 60 * 60 * 1000)
  }
])

// 当前激活的菜单
const activeMenu = computed(() => route.path)

// 面包屑导航
const breadcrumbs = computed(() => {
  const pathArray = route.path.split('/').filter(Boolean)
  const breadcrumbList = []
  
  // 根据路径生成面包屑
  let currentPath = ''
  for (const segment of pathArray) {
    currentPath += `/${segment}`
    const title = getBreadcrumbTitle(currentPath)
    if (title) {
      breadcrumbList.push({
        path: currentPath,
        title
      })
    }
  }
  
  return breadcrumbList
})

// 获取面包屑标题
const getBreadcrumbTitle = (path: string): string => {
  const titleMap: Record<string, string> = {
    '/admin': '管理后台',
    '/admin/dashboard': '仪表盘',
    '/admin/users': '用户管理',
    '/admin/users/teachers': '教师管理',
    '/admin/users/students': '学生管理',
    '/admin/statistics': '数据统计',
    '/admin/settings': '系统设置',
    '/teacher': '教师工作台',
    '/teacher/dashboard': '工作台',
    '/teacher/questions': '题目管理',
    '/teacher/questions/import': '批量导入',
    '/teacher/questions/collect': '网络采集',
    '/teacher/subjects': '科目分类',
    '/teacher/papers': '组卷管理',
    '/teacher/exams': '考试安排',
    '/teacher/monitoring': '考试监控',
    '/teacher/grades': '成绩管理',
    '/student': '学习中心',
    '/student/dashboard': '学习中心',
    '/student/practice': '练习刷题',
    '/student/exams': '我的考试',
    '/student/grades': '成绩查询',
    '/student/ai-assistant': 'AI助手'
  }
  
  return titleMap[path] || ''
}

// 切换侧边栏折叠状态
const toggleCollapse = debounce(() => {
  isCollapsed.value = !isCollapsed.value
  // 保存到localStorage
  localStorage.setItem('sidebar-collapsed', String(isCollapsed.value))
}, 100)

// 显示通知
const showNotifications = () => {
  notificationDrawer.value = true
}

// 标记通知为已读
const markAsRead = (id: number) => {
  const notification = notifications.value.find(n => n.id === id)
  if (notification) {
    notification.read = true
    notificationCount.value = notifications.value.filter(n => !n.read).length
  }
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

// 处理用户菜单命令
const handleUserCommand = async (command: string) => {
  switch (command) {
    case 'profile':
      router.push('/profile')
      break
    case 'settings':
      router.push('/settings')
      break
    case 'logout':
      try {
        await ElMessageBox.confirm(
          '确定要退出登录吗？',
          '退出确认',
          {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning'
          }
        )
        
        await authStore.logout()
        ElMessage.success('已退出登录')
        router.push('/login')
      } catch {
        // 用户取消
      }
      break
  }
}

// 监听路由变化，更新未读通知数量
watch(
  () => route.path,
  () => {
    notificationCount.value = notifications.value.filter(n => !n.read).length
  }
)

// 组件挂载时恢复侧边栏状态
onMounted(() => {
  const savedCollapsed = localStorage.getItem('sidebar-collapsed')
  if (savedCollapsed !== null) {
    isCollapsed.value = savedCollapsed === 'true'
  }
})
</script>

<style scoped>
.main-layout {
  height: 100vh;
  display: flex;
}

.sidebar {
  background: linear-gradient(180deg, var(--business-navy) 0%, var(--business-blue) 100%);
  transition: width 0.3s ease;
  box-shadow: var(--shadow-lg);
  position: relative;
  z-index: 1000;
}

.sidebar-header {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 var(--spacing-md);
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.logo {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  color: white;
}

.logo-icon {
  font-size: 24px;
  color: var(--dopamine-yellow);
}

.logo-text {
  font-size: 18px;
  font-weight: 600;
  white-space: nowrap;
}

.collapse-btn {
  background: rgba(255, 255, 255, 0.1);
  border: none;
  color: white;
}

.collapse-btn:hover {
  background: rgba(255, 255, 255, 0.2);
}

.sidebar-menu {
  height: calc(100vh - 120px);
}

.sidebar-footer {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0 var(--spacing-md);
  border-top: 1px solid rgba(255, 255, 255, 0.1);
}

.dopamine-menu {
  background: transparent;
  border: none;
}

.dopamine-menu :deep(.el-menu-item),
.dopamine-menu :deep(.el-sub-menu__title) {
  color: rgba(255, 255, 255, 0.8);
  border-radius: var(--radius-md);
  margin: 4px var(--spacing-sm);
  transition: all 0.3s ease;
}

.dopamine-menu :deep(.el-menu-item:hover),
.dopamine-menu :deep(.el-sub-menu__title:hover) {
  background: rgba(255, 255, 255, 0.1);
  color: white;
}

.dopamine-menu :deep(.el-menu-item.is-active) {
  background: var(--gradient-warning);
  color: white;
  font-weight: 600;
}

.dopamine-menu :deep(.el-sub-menu.is-active .el-sub-menu__title) {
  color: var(--business-accent);
}

/* 子菜单容器样式 */
.dopamine-menu :deep(.el-sub-menu .el-menu) {
  background: transparent !important;
}

/* 子菜单项样式 */
.dopamine-menu :deep(.el-sub-menu .el-menu-item) {
  color: rgba(255, 255, 255, 0.7);
  background: transparent !important;
  border-radius: var(--radius-md);
  margin: 2px var(--spacing-md);
  padding-left: var(--spacing-xl);
  transition: all 0.3s ease;
}

.dopamine-menu :deep(.el-sub-menu .el-menu-item:hover) {
  background: rgba(255, 255, 255, 0.1);
  color: white;
}

.dopamine-menu :deep(.el-sub-menu .el-menu-item.is-active) {
  background: var(--gradient-warning);
  color: white;
  font-weight: 600;
}

/* 折叠状态下菜单项图标居中 */
.dopamine-menu.el-menu--collapse :deep(.el-menu-item) {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0 !important;
}

.dopamine-menu.el-menu--collapse :deep(.el-menu-item .el-icon) {
  margin-right: 0 !important;
}

.main-container {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.header {
  background: var(--bg-primary);
  box-shadow: var(--shadow-sm);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0;
  border-bottom: 1px solid var(--border-light);
}

.header-left {
  flex: 1;
  padding-left: var(--spacing-lg);
}

.breadcrumb :deep(.el-breadcrumb__item) {
  font-size: 14px;
}

.breadcrumb :deep(.el-breadcrumb__inner) {
  color: var(--text-secondary);
  font-weight: 500;
}

.breadcrumb :deep(.el-breadcrumb__inner.is-link:hover) {
  color: var(--business-blue);
}

.header-right {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
  padding-right: var(--spacing-lg);
}

.notification-badge {
  margin-right: var(--spacing-sm);
}

.header-btn {
  background: var(--bg-hover);
  border: 1px solid var(--border-light);
  color: var(--text-secondary);
}

.header-btn:hover {
  background: var(--business-blue);
  border-color: var(--business-blue);
  color: white;
}

.user-dropdown {
  cursor: pointer;
}

.user-info {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  padding: var(--spacing-sm);
  border-radius: var(--radius-md);
  transition: all 0.3s ease;
}

.user-info:hover {
  background: var(--bg-hover);
}

.user-avatar {
  background: var(--gradient-primary);
}

.user-name {
  font-size: 14px;
  font-weight: 500;
  color: var(--text-primary);
}

.dropdown-icon {
  font-size: 12px;
  color: var(--text-secondary);
}

.main-content {
  background: var(--bg-secondary);
  padding: 0;
  overflow-y: auto;
}

.content-wrapper {
  width: 100%;
  margin: 0;
}

.notification-list {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
}

.notification-item {
  padding: var(--spacing-md);
  border-left: 4px solid transparent;
  transition: all 0.3s ease;
}

.notification-item.unread {
  border-left-color: var(--dopamine-orange);
  background: rgba(255, 139, 83, 0.05);
}

.notification-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--spacing-sm);
}

.notification-title {
  font-weight: 600;
  color: var(--text-primary);
}

.notification-time {
  font-size: 12px;
  color: var(--text-muted);
}

.notification-content {
  color: var(--text-secondary);
  font-size: 14px;
  line-height: 1.5;
  margin: 0 0 var(--spacing-sm) 0;
}

.notification-actions {
  display: flex;
  justify-content: flex-end;
}

/* 页面切换动画 */
.fade-transform-enter-active,
.fade-transform-leave-active {
  transition: all 0.3s ease;
}

.fade-transform-enter-from {
  opacity: 0;
  transform: translateX(20px);
}

.fade-transform-leave-to {
  opacity: 0;
  transform: translateX(-20px);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .sidebar {
    position: fixed;
    left: 0;
    top: 0;
    height: 100vh;
    z-index: 2000;
    transform: translateX(isCollapsed ? '-100%' : '0');
  }
  
  .main-container {
    margin-left: 0;
  }
  
  .header {
    padding: 0 var(--spacing-md);
  }
  
  .user-name {
    display: none;
  }
  
  .main-content {
    padding: var(--spacing-md);
  }
}
</style>