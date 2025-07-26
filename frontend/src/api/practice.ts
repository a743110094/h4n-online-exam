import { get, post } from './index'

// 练习相关API接口
export interface PracticeRecommendation {
  id: number
  title: string
  description: string
  subject_id: number
  subject: {
    id: number
    name: string
  }
  difficulty: number
  question_count: number
  estimated_time: number
  rating: number
  knowledge_point: string
  question_types: string
  progress?: number
}

export interface PracticeRecord {
  id: number
  user_id: number
  subject_id: number
  subject: {
    id: number
    name: string
  }
  title: string
  description: string
  total_count: number
  correct_count: number
  wrong_count: number
  score: number
  duration: number
  difficulty: number
  practice_type: string
  is_completed: boolean
  created_at: string
  updated_at: string
}

export interface PracticeStats {
  total_practiced: number
  today_practiced: number
  correct_rate: number
  wrong_questions: number
}

export interface StartPracticeRequest {
  subject_id: number
  difficulty?: number
  question_type?: string
  practice_type?: string
  question_count?: number
}

export interface StartPracticeResponse {
  practice_id: number
  questions: any[]
}

export interface SubmitAnswerRequest {
  question_id: number
  answer: string
  time_spent?: number
}

export interface SubmitAnswerResponse {
  is_correct: boolean
  score: number
  explanation: string
}

export interface CompletePracticeResponse {
  total_questions: number
  correct_count: number
  wrong_count: number
  score: number
  accuracy: number
}

// 获取推荐练习列表
export const getPracticeRecommendations = () => {
  return get<PracticeRecommendation[]>('/practice/recommendations')
}

// 开始练习
export const startPractice = (data: StartPracticeRequest) => {
  return post<StartPracticeResponse>('/practice/start', data)
}

// 提交练习答案
export const submitPracticeAnswer = (practiceId: number, data: SubmitAnswerRequest) => {
  return post<SubmitAnswerResponse>(`/practice/${practiceId}/answer`, data)
}

// 完成练习
export const completePractice = (practiceId: number) => {
  return post<CompletePracticeResponse>(`/practice/${practiceId}/complete`)
}

// 错题相关接口
export interface WrongQuestionDetail {
  id: number
  question_id: number
  title: string
  content: string
  options: string
  type: string
  difficulty: number
  subject_name: string
  user_answer: string
  correct_answer: string
  explanation: string
  time_spent: number
  created_at: string
}

export interface WrongQuestionsResponse {
  data: WrongQuestionDetail[]
  total: number
  page: number
  page_size: number
  total_pages: number
}

export interface StartReviewRequest {
  subject_id?: number
  practice_record_id?: number
  question_ids?: number[]
  max_questions?: number
}

export interface StartReviewResponse {
  practice_id: number
  questions: Question[]
  type: string
}

// 获取错题列表
export const getWrongQuestions = (params?: {
  page?: number
  page_size?: number
  subject_id?: number
  practice_record_id?: number
}) => {
  return get<WrongQuestionsResponse>('/practice/wrong-questions', params)
}

// 开始错题复习
export const startWrongQuestionReview = (data: StartReviewRequest) => {
  return post<StartReviewResponse>('/practice/review/start', data)
}

// 获取练习历史
export const getPracticeHistory = (page = 1, pageSize = 10) => {
  return get<{
    data: PracticeRecord[]
    total: number
    page: number
    page_size: number
    total_pages: number
  }>(`/practice/history?page=${page}&page_size=${pageSize}`)
}

// 获取练习统计
export const getPracticeStats = () => {
  return get<PracticeStats>('/practice/stats')
}