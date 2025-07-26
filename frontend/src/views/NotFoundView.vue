<template>
  <div class="not-found-page">
    <div class="not-found-container">
      <!-- 404 插图 -->
      <div class="illustration">
        <div class="number-404">
          <span class="number-4">4</span>
          <span class="number-0">0</span>
          <span class="number-4-2">4</span>
        </div>
        <div class="floating-elements">
          <div class="element element-1"></div>
          <div class="element element-2"></div>
          <div class="element element-3"></div>
          <div class="element element-4"></div>
        </div>
      </div>
      
      <!-- 错误信息 -->
      <div class="error-content">
        <h1 class="error-title">页面走丢了</h1>
        <p class="error-description">
          抱歉，您访问的页面不存在或已被移动。
          <br>
          请检查网址是否正确，或返回首页继续浏览。
        </p>
        
        <!-- 操作按钮 -->
        <div class="action-buttons">
          <el-button
            type="primary"
            size="large"
            class="dopamine-btn"
            @click="goHome"
          >
            <el-icon><House /></el-icon>
            返回首页
          </el-button>
          
          <el-button
            size="large"
            class="dopamine-btn-outline"
            @click="goBack"
          >
            <el-icon><ArrowLeft /></el-icon>
            返回上页
          </el-button>
        </div>
        
        <!-- 快速导航 -->
        <div class="quick-nav">
          <h3>您可能想要访问：</h3>
          <div class="nav-links">
            <router-link
              v-for="link in quickLinks"
              :key="link.path"
              :to="link.path"
              class="nav-link dopamine-card"
            >
              <div class="nav-icon" :style="{ background: link.color }">
                <el-icon :size="20">
                  <component :is="link.icon" />
                </el-icon>
              </div>
              <div class="nav-content">
                <div class="nav-title">{{ link.title }}</div>
                <div class="nav-desc">{{ link.description }}</div>
              </div>
            </router-link>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 背景装饰 -->
    <div class="background-decoration">
      <div class="decoration-circle circle-1"></div>
      <div class="decoration-circle circle-2"></div>
      <div class="decoration-circle circle-3"></div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { House, ArrowLeft, DataBoard, Reading, Notebook, School } from '@element-plus/icons-vue'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

// 快速导航链接
const quickLinks = computed(() => {
  const baseLinks = [
    {
      path: '/login',
      title: '登录',
      description: '登录到您的账户',
      icon: 'House',
      color: 'var(--gradient-primary)'
    }
  ]
  
  if (authStore.isAdmin) {
    return [
      {
        path: '/admin/dashboard',
        title: '管理后台',
        description: '系统管理和数据统计',
        icon: 'DataBoard',
        color: 'var(--gradient-primary)'
      },
      {
        path: '/admin/users/teachers',
        title: '教师管理',
        description: '管理教师账户',
        icon: 'School',
        color: 'var(--gradient-success)'
      },
      {
        path: '/admin/statistics',
        title: '数据统计',
        description: '查看系统统计数据',
        icon: 'DataBoard',
        color: 'var(--gradient-warning)'
      }
    ]
  } else if (authStore.isTeacher) {
    return [
      {
        path: '/teacher/dashboard',
        title: '教师工作台',
        description: '教学管理中心',
        icon: 'DataBoard',
        color: 'var(--gradient-primary)'
      },
      {
        path: '/teacher/questions',
        title: '题库管理',
        description: '管理考试题目',
        icon: 'Reading',
        color: 'var(--gradient-success)'
      },
      {
        path: '/teacher/exams',
        title: '考试管理',
        description: '创建和管理考试',
        icon: 'Notebook',
        color: 'var(--gradient-warning)'
      }
    ]
  } else if (authStore.isStudent) {
    return [
      {
        path: '/student/dashboard',
        title: '学习中心',
        description: '学习进度和考试安排',
        icon: 'DataBoard',
        color: 'var(--gradient-primary)'
      },
      {
        path: '/student/practice',
        title: '练习刷题',
        description: '巩固知识点',
        icon: 'Reading',
        color: 'var(--gradient-success)'
      },
      {
        path: '/student/exams',
        title: '我的考试',
        description: '查看考试安排',
        icon: 'Notebook',
        color: 'var(--gradient-warning)'
      }
    ]
  }
  
  return baseLinks
})

// 返回首页
const goHome = () => {
  if (authStore.isAdmin) {
    router.push('/admin/dashboard')
  } else if (authStore.isTeacher) {
    router.push('/teacher/dashboard')
  } else if (authStore.isStudent) {
    router.push('/student/dashboard')
  } else {
    router.push('/login')
  }
}

// 返回上一页
const goBack = () => {
  if (window.history.length > 1) {
    router.go(-1)
  } else {
    goHome()
  }
}
</script>

<style scoped>
.not-found-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--bg-secondary);
  position: relative;
  overflow: hidden;
}

.not-found-container {
  max-width: 800px;
  width: 100%;
  padding: var(--spacing-xl);
  text-align: center;
  position: relative;
  z-index: 2;
}

.illustration {
  position: relative;
  margin-bottom: var(--spacing-xl);
}

.number-404 {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: var(--spacing-md);
  margin-bottom: var(--spacing-lg);
}

.number-404 span {
  font-size: 120px;
  font-weight: 900;
  line-height: 1;
  background: var(--gradient-primary);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  animation: bounce 2s ease-in-out infinite;
}

.number-4 {
  animation-delay: 0s;
}

.number-0 {
  animation-delay: 0.2s;
}

.number-4-2 {
  animation-delay: 0.4s;
}

.floating-elements {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  pointer-events: none;
}

.element {
  position: absolute;
  border-radius: 50%;
  animation: float 3s ease-in-out infinite;
}

.element-1 {
  width: 20px;
  height: 20px;
  background: var(--dopamine-blue);
  top: 20%;
  left: 10%;
  animation-delay: 0s;
}

.element-2 {
  width: 30px;
  height: 30px;
  background: var(--dopamine-pink);
  top: 60%;
  right: 15%;
  animation-delay: 1s;
}

.element-3 {
  width: 15px;
  height: 15px;
  background: var(--dopamine-yellow);
  top: 40%;
  left: 80%;
  animation-delay: 2s;
}

.element-4 {
  width: 25px;
  height: 25px;
  background: var(--dopamine-green);
  top: 80%;
  left: 20%;
  animation-delay: 1.5s;
}

.error-content {
  max-width: 600px;
  margin: 0 auto;
}

.error-title {
  font-size: 36px;
  font-weight: 700;
  color: var(--text-primary);
  margin: 0 0 var(--spacing-lg) 0;
  background: var(--gradient-primary);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.error-description {
  font-size: 18px;
  color: var(--text-secondary);
  line-height: 1.6;
  margin: 0 0 var(--spacing-xl) 0;
}

.action-buttons {
  display: flex;
  gap: var(--spacing-lg);
  justify-content: center;
  margin-bottom: var(--spacing-xl);
}

.dopamine-btn {
  background: var(--gradient-primary);
  border: none;
  color: white;
  font-weight: 600;
  padding: var(--spacing-md) var(--spacing-xl);
  border-radius: var(--radius-lg);
  transition: all 0.3s ease;
  box-shadow: var(--shadow-md);
}

.dopamine-btn:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-lg);
}

.dopamine-btn-outline {
  background: transparent;
  border: 2px solid var(--dopamine-blue);
  color: var(--dopamine-blue);
  font-weight: 600;
  padding: var(--spacing-md) var(--spacing-xl);
  border-radius: var(--radius-lg);
  transition: all 0.3s ease;
}

.dopamine-btn-outline:hover {
  background: var(--dopamine-blue);
  color: white;
  transform: translateY(-2px);
}

.quick-nav {
  margin-top: var(--spacing-xl);
}

.quick-nav h3 {
  font-size: 20px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 var(--spacing-lg) 0;
}

.nav-links {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: var(--spacing-lg);
}

.nav-link {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
  padding: var(--spacing-lg);
  text-decoration: none;
  transition: all 0.3s ease;
  border: 2px solid transparent;
}

.nav-link:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow-xl);
  border-color: var(--dopamine-blue);
  text-decoration: none;
}

.nav-icon {
  width: 48px;
  height: 48px;
  border-radius: var(--radius-lg);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  flex-shrink: 0;
}

.nav-content {
  flex: 1;
  text-align: left;
}

.nav-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 4px;
}

.nav-desc {
  font-size: 14px;
  color: var(--text-secondary);
}

.background-decoration {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  pointer-events: none;
  z-index: 1;
}

.decoration-circle {
  position: absolute;
  border-radius: 50%;
  opacity: 0.1;
  animation: rotate 20s linear infinite;
}

.circle-1 {
  width: 200px;
  height: 200px;
  background: var(--gradient-primary);
  top: 10%;
  left: -100px;
}

.circle-2 {
  width: 300px;
  height: 300px;
  background: var(--gradient-success);
  bottom: 10%;
  right: -150px;
  animation-direction: reverse;
}

.circle-3 {
  width: 150px;
  height: 150px;
  background: var(--gradient-warning);
  top: 50%;
  right: 10%;
  animation-delay: 10s;
}

/* 动画 */
@keyframes bounce {
  0%, 20%, 50%, 80%, 100% {
    transform: translateY(0);
  }
  40% {
    transform: translateY(-20px);
  }
  60% {
    transform: translateY(-10px);
  }
}

@keyframes float {
  0%, 100% {
    transform: translateY(0px);
  }
  50% {
    transform: translateY(-20px);
  }
}

@keyframes rotate {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

/* 响应式设计 */
@media (max-width: 768px) {
  .number-404 span {
    font-size: 80px;
  }
  
  .error-title {
    font-size: 28px;
  }
  
  .error-description {
    font-size: 16px;
  }
  
  .action-buttons {
    flex-direction: column;
    align-items: center;
  }
  
  .nav-links {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 480px) {
  .number-404 {
    flex-direction: column;
    gap: 0;
  }
  
  .number-404 span {
    font-size: 60px;
  }
  
  .not-found-container {
    padding: var(--spacing-lg);
  }
}
</style>