import { get, post, put, del } from './index'

// 科目相关API接口
export interface Subject {
  id: number
  name: string
  description?: string
  created_at?: string
  updated_at?: string
}

export interface SubjectListResponse {
  subjects: Subject[]
  total: number
  page: number
  size: number
}

// 获取科目列表
export const getSubjects = (params?: {
  page?: number
  size?: number
  keyword?: string
}): Promise<SubjectListResponse> => {
  return get('/subjects', params)
}

// 获取科目详情
export const getSubject = (id: number): Promise<Subject> => {
  return get(`/subjects/${id}`)
}

// 创建科目
export const createSubject = (data: {
  name: string
  description?: string
}): Promise<Subject> => {
  return post('/subjects', data)
}

// 更新科目
export const updateSubject = (id: number, data: {
  name: string
  description?: string
}): Promise<Subject> => {
  return put(`/subjects/${id}`, data)
}

// 删除科目
export const deleteSubject = (id: number): Promise<void> => {
  return del(`/subjects/${id}`)
}