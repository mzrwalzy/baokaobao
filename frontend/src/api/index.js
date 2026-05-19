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

// Admin User Management
export const getAdminUsers = (params) => request.get('/admin_users', { params })
export const createAdminUser = (data) => request.post('/admin_users', data)
export const updateAdminUser = (id, data) => request.put(`/admin_users/${id}`, data)
export const resetAdminPassword = (id, data) => request.put(`/admin_users/${id}/reset_password`, data)
export const deleteAdminUser = (id) => request.delete(`/admin_users/${id}`)
export const getAdminUserBanks = (id) => request.get(`/admin_users/${id}/banks`)
export const grantAdminBankAccess = (id, data) => request.post(`/admin_users/${id}/banks`, data)
export const revokeAdminBankAccess = (id, bankId) => request.delete(`/admin_users/${id}/banks/${bankId}`)

export const getAuditLogs = (params) => request.get('/audit_logs', { params })
export const getBankStatsList = (params) => request.get('/bank_stats', { params })
export const getBankStats = (id) => request.get(`/bank_stats/${id}`)
export const getFeedback = (params) => request.get('/feedbacks', { params })
export const updateFeedbackStatus = (id, data) => request.put(`/feedbacks/${id}/status`, data)
