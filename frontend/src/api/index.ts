import axios from 'axios'

// 创建axios实例
const api = axios.create({
  baseURL: 'http://localhost:8080/api/v1', // 后端API地址
  timeout: 10000, // 请求超时时间
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器
api.interceptors.request.use(
  (config) => {
    // 从localStorage获取token
    const token = localStorage.getItem('auth_token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    
    // 添加租户ID到请求头
    // 默认使用租户ID 100，实际项目中可以从用户配置或环境变量获取
    config.headers['X-Tenant-ID'] = '100'
    
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器
api.interceptors.response.use(
  (response) => {
    return response
  },
  (error) => {
    // 处理401未授权错误
    if (error.response?.status === 401) {
      // 清除本地存储的认证信息
      localStorage.removeItem('auth_token')
      localStorage.removeItem('auth_user')
      // 重定向到登录页
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

export default api

// 导出常用的HTTP方法
export const get = async <T = any>(url: string, params?: any): Promise<T> => {
  const response = await api.get(url, { params })
  // 如果响应有统一格式，提取data字段
  if (response.data && typeof response.data === 'object' && 'data' in response.data) {
    return response.data.data
  }
  return response.data
}

export const post = async <T = any>(url: string, data?: any): Promise<T> => {
  const response = await api.post(url, data)
  // 如果响应有统一格式，提取data字段
  if (response.data && typeof response.data === 'object' && 'data' in response.data) {
    return response.data.data
  }
  return response.data
}

export const put = async <T = any>(url: string, data?: any): Promise<T> => {
  const response = await api.put(url, data)
  // 如果响应有统一格式，提取data字段
  if (response.data && typeof response.data === 'object' && 'data' in response.data) {
    return response.data.data
  }
  return response.data
}

export const del = async <T = any>(url: string): Promise<T> => {
  const response = await api.delete(url)
  // 如果响应有统一格式，提取data字段
  if (response.data && typeof response.data === 'object' && 'data' in response.data) {
    return response.data.data
  }
  return response.data
}