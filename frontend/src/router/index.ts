import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { UserRole } from '@/stores/auth'
import { ElMessage } from 'element-plus'

// 导入组件
import LoginView from '@/views/LoginView.vue'
import MainLayout from '@/components/Layout/MainLayout.vue'

// 管理员页面
import AdminDashboard from '@/views/admin/Dashboard.vue'
import UserManagementView from '@/views/admin/UserManagementView.vue'

// 教师页面
import TeacherDashboard from '@/views/teacher/Dashboard.vue'
import QuestionBankView from '@/views/teacher/QuestionBankView.vue'
import QuestionImportView from '@/views/teacher/QuestionImportView.vue'
import QuestionCollectView from '@/views/teacher/QuestionCollectView.vue'
import SubjectManagementView from '@/views/teacher/SubjectManagementView.vue'
import ExamCreationView from '@/views/teacher/ExamCreationView.vue'
import ExamManagementView from '@/views/teacher/ExamManagementView.vue'
import MonitoringView from '@/views/teacher/MonitoringView.vue'
import PaperManagementView from '@/views/teacher/PaperManagementView.vue'
import GradeEntryView from '@/views/teacher/GradeEntryView.vue'
import GradeQueryView from '@/views/teacher/GradeQueryView.vue'
import GradeExportView from '@/views/teacher/GradeExportView.vue'

// 学生页面
import StudentDashboard from '@/views/student/Dashboard.vue'
import PracticeView from '@/views/student/PracticeView.vue'
import PracticeSessionView from '@/views/student/PracticeSessionView.vue'
import ReviewView from '@/views/student/ReviewView.vue'
import ExamView from '@/views/student/ExamView.vue'
import ExamTakingView from '@/views/student/ExamTakingView.vue'

// 通用页面
import NotFoundView from '@/views/NotFoundView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'Login',
      component: LoginView,
      meta: { requiresAuth: false }
    },
    {
      path: '/',
      redirect: '/dashboard'
    },
    {
      path: '/dashboard',
      redirect: (to) => {
        const authStore = useAuthStore()
        if (authStore.isAdmin) return '/admin/dashboard'
        if (authStore.isTeacher) return '/teacher/dashboard'
        if (authStore.isStudent) return '/student/dashboard'
        return '/login'
      }
    },
    {
      path: '/admin',
      component: MainLayout,
      meta: { requiresAuth: true, roles: [UserRole.ADMIN] },
      children: [
        {
          path: 'dashboard',
          name: 'AdminDashboard',
          component: AdminDashboard
        },
        {
          path: 'users',
          name: 'UserManagement',
          component: UserManagementView
        },
        {
          path: 'users/teachers',
          name: 'TeacherManagement',
          component: () => import('@/views/admin/TeacherManagementView.vue')
        },
        {
          path: 'users/students',
          name: 'StudentManagement',
          component: () => import('@/views/admin/StudentManagementView.vue')
        },
        {
          path: 'statistics',
          name: 'Statistics',
          component: () => import('@/views/admin/StatisticsView.vue')
        },
        {
          path: 'settings',
          name: 'SystemSettings',
          component: () => import('@/views/admin/SystemSettingsView.vue')
        }
      ]
    },
    {
      path: '/teacher',
      component: MainLayout,
      meta: { requiresAuth: true, roles: [UserRole.TEACHER] },
      children: [
        {
          path: 'dashboard',
          name: 'TeacherDashboard',
          component: TeacherDashboard
        },
        {
          path: 'questions',
          name: 'QuestionBank',
          component: QuestionBankView
        },
        {
          path: 'questions/import',
          name: 'QuestionImport',
          component: QuestionImportView
        },
        {
          path: 'questions/collect',
          name: 'QuestionCollect',
          component: QuestionCollectView
        },
        {
          path: 'subjects',
          name: 'SubjectManagement',
          component: SubjectManagementView
        },
        {
          path: 'exams/create',
          name: 'ExamCreation',
          component: ExamCreationView
        },
        {
          path: 'exams',
          name: 'ExamManagement',
          component: ExamManagementView
        },
        {
          path: 'monitoring',
          name: 'Monitoring',
          component: MonitoringView
        },
        {
          path: 'papers',
          name: 'PaperManagement',
          component: PaperManagementView
        },
        {
          path: 'grades/entry',
          name: 'GradeEntry',
          component: GradeEntryView
        },
        {
          path: 'grades/query',
          name: 'GradeQuery',
          component: GradeQueryView
        },
        {
          path: 'grades/export',
          name: 'GradeExport',
          component: GradeExportView
        }
      ]
    },
    {
      path: '/student',
      component: MainLayout,
      meta: { requiresAuth: true, roles: [UserRole.STUDENT] },
      children: [
        {
          path: 'dashboard',
          name: 'StudentDashboard',
          component: StudentDashboard
        },
        {
          path: 'practice',
          name: 'PracticeView',
          component: PracticeView
        },
        {
          path: 'practice/:id/:practiceType',
          name: 'PracticeSession',
          component: PracticeSessionView
        },
        {
          path: 'review',
          name: 'ReviewView',
          component: ReviewView
        },
        {
          path: 'exams',
          name: 'StudentExams',
          component: ExamView
        },
        {
          path: 'grades',
          name: 'StudentGrades',
          component: () => import('@/views/student/GradeView.vue')
        }
      ]
    },
    {
      path: '/exam/:examId',
      name: 'ExamTaking',
      component: ExamTakingView,
      meta: { requiresAuth: true, roles: [UserRole.STUDENT] }
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'NotFound',
      component: NotFoundView
    }
  ]
})

// 路由守卫
router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore()
  
  // 检查是否需要认证
  if (to.meta.requiresAuth !== false) {
    // 如果没有token，重定向到登录页
    if (!authStore.token) {
      next('/login')
      return
    }
    
    // 如果有token但没有用户信息，尝试获取用户信息
    if (!authStore.user) {
      try {
        await authStore.fetchUserInfo()
      } catch (error) {
        ElMessage.error('获取用户信息失败，请重新登录')
        authStore.logout()
        next('/login')
        return
      }
    }
    
    // 检查角色权限
    if (to.meta.roles && Array.isArray(to.meta.roles)) {
      const hasPermission = authStore.hasAnyRole(to.meta.roles as UserRole[])
      if (!hasPermission) {
        ElMessage.error('您没有权限访问该页面')
        // 根据用户角色重定向到对应的首页
        if (authStore.isAdmin) {
          next('/admin/dashboard')
        } else if (authStore.isTeacher) {
          next('/teacher/dashboard')
        } else if (authStore.isStudent) {
          next('/student/dashboard')
        } else {
          next('/login')
        }
        return
      }
    }
  }
  
  // 如果已登录用户访问登录页，重定向到对应首页
  if (to.path === '/login' && authStore.token && authStore.user) {
    if (authStore.isAdmin) {
      next('/admin/dashboard')
    } else if (authStore.isTeacher) {
      next('/teacher/dashboard')
    } else if (authStore.isStudent) {
      next('/student/dashboard')
    } else {
      next()
    }
    return
  }
  
  next()
})

export default router
