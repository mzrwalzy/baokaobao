const { get, post, put, BASE_URL } = require('../utils/request.js')

const getQuestionBanks = (params) => get('/question_banks', params)
const getQuestionBankDetail = (id) => get(`/question_banks/${id}`)
const getQuestionBankAccessStatus = (id) => get(`/question_banks/${id}/access`)
const getQuestions = (params) => get('/questions', params)
const getRandomQuestions = (bank_id, count = 10) => get(`/questions/random?bank_id=${bank_id}&count=${count}`)
const submitAnswer = (data) => post('/quiz/submit', data)
const submitExam = (data) => post('/quiz/submit_exam', data)
const getQuizHistory = (params) => get('/quiz/history', params)
const getWrongQuestions = (params) => get('/quiz/wrong_questions', params)
const getMyScore = () => get('/score/my')
const getRanking = (type = 'total') => get(`/score/ranking?type=${type}`)
const getStats = () => get('/score/stats')
const getProfile = () => get('/user/profile')
const updateProfile = (data) => put('/user/profile', data)
const loginByWechat = (code) => post('/auth/login_by_wechat', { code })
const logout = () => post('/auth/logout')
const getExamRecords = (params) => get('/exam_records', { params })
const getPurchasedBanks = () => get('/my/banks')
const createFeedback = (data) => post('/feedbacks', data)
const uploadAvatar = (filePath) => {
  const app = getApp()
  const token = app.globalData.token

  return new Promise((resolve, reject) => {
    wx.uploadFile({
      url: `${BASE_URL}/user/avatar`,
      filePath,
      name: 'file',
      header: {
        'Authorization': token ? `Bearer ${token}` : ''
      },
      success: (res) => {
        try {
          const data = JSON.parse(res.data)
          if (res.statusCode === 200 && (data.code === 0 || data.code === 200)) {
            resolve(data.data)
            return
          }
          const message = data.msg || '上传头像失败'
          wx.showToast({ title: message, icon: 'none' })
          reject(new Error(message))
        } catch (err) {
          reject(err)
        }
      },
      fail: reject
    })
  })
}

module.exports = {
  getQuestionBanks,
  getQuestionBankDetail,
  getQuestionBankAccessStatus,
  getQuestions,
  getRandomQuestions,
  submitAnswer,
  submitExam,
  getQuizHistory,
  getWrongQuestions,
  getMyScore,
  getRanking,
  getStats,
  getProfile,
  updateProfile,
  uploadAvatar,
  loginByWechat,
  logout,
  getExamRecords,
  getPurchasedBanks,
  createFeedback
}
