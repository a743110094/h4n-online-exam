import { get, post, put, del } from './index'

// 题目相关API接口
export interface Question {
  id?: number
  subject_id: number
  type: string
  title: string
  content: string
  options?: string
  answer: string
  explanation?: string
  difficulty: number
  score: number
  status?: string
  knowledge_point?: string
  usage_count?: number
  correct_rate?: number
  created_by?: number
  created_at?: string
  updated_at?: string
  subject?: {
    id: number
    name: string
  }
  creator?: {
    id: number
    name: string
  }
}

export interface QuestionListParams {
  page?: number
  size?: number
  subject_id?: number
  type?: string
  difficulty?: string
  search?: string
}

export interface QuestionListResponse {
  questions: Question[]
  total: number
  page: number
  size: number
}

// 获取题目列表
export const getQuestions = (params?: QuestionListParams) => {
  return get('/questions', params)
}

// 获取单个题目
export const getQuestion = (id: number) => {
  return get(`/questions/${id}`)
}

// 创建题目
export const createQuestion = (data: Question) => {
  return post('/teacher/questions', data)
}

// 更新题目
export const updateQuestion = (id: number, data: Question) => {
  return put(`/teacher/questions/${id}`, data)
}

// 删除题目
export const deleteQuestion = (id: number) => {
  return del(`/teacher/questions/${id}`)
}

// 批量导入题目
export const batchImportQuestions = (data: Question[]) => {
  return post('/teacher/questions/import', { questions: data })
}

// 获取题目统计
export const getQuestionStats = () => {
  return get('/questions/stats')
}

// 获取科目列表
export const getSubjects = () => {
  return get('/subjects')
}