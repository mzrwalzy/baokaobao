const { get, post, put } = require('../utils/request.js')

const getQuestionBanks = (params) => get('/question_banks', params)
const getQuestionBankDetail = (id) => get(`/question_banks/${id}`)
const getQuestions = (params) => get('/questions', params)
const getRandomQuestions = (bank_id, count = 10) => get(`/questions/random?bank_id=${bank_id}&count=${count}`)
const submitAnswer = (data) => post('/quiz/submit', data)
const getQuizHistory = (params) => get('/quiz/history', params)
const getWrongQuestions = (params) => get('/quiz/wrong_questions', params)
const getMyScore = () => get('/score/my')
const getRanking = (type = 'total') => get(`/score/ranking?type=${type}`)
const getStats = () => get('/score/stats')
const getProfile = () => get('/user/profile')
const updateProfile = (data) => put('/user/profile', data)
const loginByWechat = (code) => post('/auth/login_by_wechat', { code })
const logout = () => post('/auth/logout')

module.exports = {
  getQuestionBanks,
  getQuestionBankDetail,
  getQuestions,
  getRandomQuestions,
  submitAnswer,
  getQuizHistory,
  getWrongQuestions,
  getMyScore,
  getRanking,
  getStats,
  getProfile,
  updateProfile,
  loginByWechat,
  logout
}
