import request from './request'

export const adminLogin = (data) => request.post('/login', data)
export const adminLogout = () => request.post('/logout')

export const getDashboard = () => request.get('/dashboard')
export const getUserStats = () => request.get('/stats/users')
export const getQuestionStats = () => request.get('/stats/questions')

export const getUsers = (params) => request.get('/users', { params })
export const getUserDetail = (id) => request.get(`/users/${id}`)
export const updateUserStatus = (id, data) => request.put(`/users/${id}/status`, data)

export const getQuestionBanks = (params) => request.get('/question_banks', { params })
export const createQuestionBank = (data) => request.post('/question_banks', data)
export const updateQuestionBank = (id, data) => request.put(`/question_banks/${id}`, data)
export const deleteQuestionBank = (id) => request.delete(`/question_banks/${id}`)

export const getQuestions = (params) => request.get('/questions', { params })
export const getQuestionDetail = (id) => request.get(`/questions/${id}`)
export const createQuestion = (data) => request.post('/questions', data)
export const updateQuestion = (id, data) => request.put(`/questions/${id}`, data)
export const deleteQuestion = (id) => request.delete(`/questions/${id}`)
export const importQuestions = (data) => request.post('/questions/import', data)
