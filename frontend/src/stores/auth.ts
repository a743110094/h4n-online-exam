import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '@/api'

// 用户角色枚举
export enum UserRole {
  ADMIN = 'admin',
  TEACHER = 'teacher',
  STUDENT = 'student'
}

// 用户信息接口
export interface User {
  id: number
  username: string
  email: string
  role: UserRole
  avatar?: string
  realName?: string
  phone?: string
  createdAt: string
  updatedAt: string
}

// 登录表单接口
export interface LoginForm {
  username: string
  password: string
  remember: boolean
}

// 注册表单接口
export interface RegisterForm {
  username: string
  email: string
  password: string
  confirmPassword: string
  realName: string
  phone?: string
}

export const useAuthStore = defineStore('auth', () => {
  // 状态
  const user = ref<User | null>(null)
  const token = ref<string>('')
  const isLoading = ref(false)
  const isLoggedIn = computed(() => !!token.value && !!user.value)
  
  // 权限计算
  const isAdmin = computed(() => user.value?.role === UserRole.ADMIN)
  const isTeacher = computed(() => user.value?.role === UserRole.TEACHER)
  const isStudent = computed(() => user.value?.role === UserRole.STUDENT)
  
  // 权限检查
  const hasPermission = (requiredRole: UserRole | UserRole[]) => {
    if (!user.value) return false
    
    const roles = Array.isArray(requiredRole) ? requiredRole : [requiredRole]
    return roles.includes(user.value.role)
  }
  
  // 检查是否拥有任一角色权限
  const hasAnyRole = (roles: UserRole[]) => {
    if (!user.value) return false
    return roles.includes(user.value.role)
  }
  
  // 从localStorage恢复登录状态
  const restoreAuth = () => {
    const savedToken = localStorage.getItem('auth_token')
    const savedUser = localStorage.getItem('auth_user')
    
    if (savedToken && savedUser) {
      token.value = savedToken
      user.value = JSON.parse(savedUser)
      
      // token会通过请求拦截器自动添加到请求头中
    }
  }
  
  // 保存认证信息到localStorage
  const saveAuth = (authToken: string, userData: User) => {
    token.value = authToken
    user.value = userData
    
    localStorage.setItem('auth_token', authToken)
    localStorage.setItem('auth_user', JSON.stringify(userData))
    
    // token会通过请求拦截器自动添加到请求头中
  }
  
  // 清除认证信息
  const clearAuth = () => {
    token.value = ''
    user.value = null
    
    localStorage.removeItem('auth_token')
    localStorage.removeItem('auth_user')
    
    // token会通过响应拦截器自动清除
  }
  
  // 登录
  const login = async (loginForm: LoginForm) => {
    try {
      isLoading.value = true
      
      const response = await api.post('/auth/login', {
        username: loginForm.username,
        password: loginForm.password
      })
      
      const { token: authToken, user: userData } = response.data
      
      saveAuth(authToken, userData)
      
      return { success: true, message: '登录成功' }
    } catch (error: any) {
      console.error('登录失败:', error)
      return { 
        success: false, 
        message: error.response?.data?.message || '登录失败，请检查用户名和密码' 
      }
    } finally {
      isLoading.value = false
    }
  }
  
  // 注册
  const register = async (registerForm: RegisterForm) => {
    try {
      isLoading.value = true
      
      const response = await api.post('/auth/register', {
        username: registerForm.username,
        email: registerForm.email,
        password: registerForm.password,
        realName: registerForm.realName,
        phone: registerForm.phone
      })
      
      return { success: true, message: '注册成功，请登录' }
    } catch (error: any) {
      console.error('注册失败:', error)
      return { 
        success: false, 
        message: error.response?.data?.message || '注册失败，请稍后重试' 
      }
    } finally {
      isLoading.value = false
    }
  }
  
  // 登出
  const logout = async () => {
    try {
      await api.post('/auth/logout')
    } catch (error) {
      console.error('登出API调用失败:', error)
    } finally {
      clearAuth()
    }
  }
  
  // 获取用户信息
  const fetchUserInfo = async () => {
    try {
      if (!token.value) return
      
      const response = await api.get('/user/profile')
      user.value = response.data
      
      // 更新localStorage中的用户信息
      localStorage.setItem('auth_user', JSON.stringify(response.data))
    } catch (error) {
      console.error('刷新用户信息失败:', error)
      // 如果token无效，清除认证信息
      clearAuth()
    }
  }
  
  // 更新用户信息
  const updateProfile = async (profileData: Partial<User>) => {
    try {
      isLoading.value = true
      
      const response = await api.put('/user/profile', profileData)
      
      user.value = { ...user.value!, ...response.data }
      localStorage.setItem('auth_user', JSON.stringify(user.value))
      
      return { success: true, message: '个人信息更新成功' }
    } catch (error: any) {
      console.error('更新个人信息失败:', error)
      return { 
        success: false, 
        message: error.response?.data?.message || '更新失败，请稍后重试' 
      }
    } finally {
      isLoading.value = false
    }
  }
  
  // 修改密码
  const changePassword = async (oldPassword: string, newPassword: string) => {
    try {
      isLoading.value = true
      
      await api.put('/user/password', {
        oldPassword,
        newPassword
      })
      
      return { success: true, message: '密码修改成功' }
    } catch (error: any) {
      console.error('修改密码失败:', error)
      return { 
        success: false, 
        message: error.response?.data?.message || '密码修改失败' 
      }
    } finally {
      isLoading.value = false
    }
  }
  
  return {
    // 状态
    user,
    token,
    isLoading,
    isLoggedIn,
    isAdmin,
    isTeacher,
    isStudent,
    
    // 方法
    hasPermission,
    hasAnyRole,
    restoreAuth,
    login,
    register,
    logout,
    fetchUserInfo,
    updateProfile,
    changePassword
  }
})

// 模拟数据（开发阶段使用）
export const mockUsers: User[] = [
  {
    id: 1,
    username: 'admin',
    email: 'admin@example.com',
    role: UserRole.ADMIN,
    realName: '系统管理员',
    phone: '13800138000',
    createdAt: '2024-01-01T00:00:00Z',
    updatedAt: '2024-01-01T00:00:00Z'
  },
  {
    id: 2,
    username: 'teacher1',
    email: 'teacher1@example.com',
    role: UserRole.TEACHER,
    realName: '张老师',
    phone: '13800138001',
    createdAt: '2024-01-01T00:00:00Z',
    updatedAt: '2024-01-01T00:00:00Z'
  },
  {
    id: 3,
    username: 'student1',
    email: 'student1@example.com',
    role: UserRole.STUDENT,
    realName: '李同学',
    phone: '13800138002',
    createdAt: '2024-01-01T00:00:00Z',
    updatedAt: '2024-01-01T00:00:00Z'
  }
]