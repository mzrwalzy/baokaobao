import { get, post } from './request'

// 题库相关
export const getQuestionBanks = (params?: { page?: number; page_size?: number }) => 
  get<{ list: QuestionBank[]; total: number }>('/question_banks', params)

export const getQuestionBankDetail = (id: number) => 
  get<QuestionBank>(`/question_banks/${id}`)

// 题目相关
export const getQuestions = (params: { bank_id: number; type?: string; page?: number; page_size?: number }) =>
  get<{ list: Question[]; total: number }>('/questions', params)

export const getRandomQuestions = (bank_id: number, count: number = 10) =>
  get<Question[]>(`/questions/random?bank_id=${bank_id}&count=${count}`)

// 答题相关
export const submitAnswer = (data: { question_id: number; answer: string; is_correct: boolean }) =>
  post('/quiz/submit', data)

export const getQuizHistory = (params?: { page?: number; page_size?: number }) =>
  get('/quiz/history', params)

export const getWrongQuestions = (params?: { page?: number; page_size?: number }) =>
  get<{ list: Question[]; total: number }>('/quiz/wrong_questions', params)

// 得分相关
export const getMyScore = () => get<Score>('/score/my')
export const getRanking = (type: 'monthly' | 'total' = 'total') => 
  get<RankingItem[]>(`/score/ranking?type=${type}`)
export const getStats = () => get<UserStats>('/score/stats')

// 用户相关
export const getProfile = () => get<UserInfo>('/user/profile')
export const updateProfile = (data: Partial<UserInfo>) => put('/user/profile', data)

// 微信登录
export const loginByWechat = (code: string) => 
  post<{ token: string; user: UserInfo }>('/auth/login_by_wechat', { code })

export const logout = () => post('/auth/logout')

// 类型定义
export interface QuestionBank {
  id: number
  name: string
  description: string
  price: number
  question_count: number
  difficulty: number
  cover_image?: string
  created_at: string
}

export interface Question {
  id: number
  bank_id: number
  title?: string
  content: string
  type: 'single' | 'multiple' | 'truefalse'
  difficulty: number
  answer: string
  analysis?: string
  options: QuestionOption[]
}

export interface QuestionOption {
  id: number
  question_id: number
  option_key: string
  option_value: string
}

export interface Score {
  total_score: number
  monthly_score: number
  rank: number
}

export interface RankingItem {
  rank: number
  nickname: string
  avatar?: string
  score: number
}

export interface UserStats {
  total_questions: number
  correct_rate: number
  study_days: number
  wrong_count: number
}

export interface UserInfo {
  id: number
  nickname: string
  avatar?: string
  phone?: string
  union_id?: string
  created_at: string
}
