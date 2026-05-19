const app = getApp()

Page({
  data: {
    bankId: 0,
    answers: [],
    startTime: 0,
    duration: 0,
    mode: 'exam',
    accuracy: 0,
    score: 0,
    formattedDuration: '0分0秒'
  },

  onLoad(options) {
    this.setData({
      bankId: Number(options.bank_id) || 0,
      mode: options.mode || 'exam'
    })
    
    const quizData = app.globalData.quizData || {}
    const answers = quizData.answers || []
    const startTime = quizData.startTime || Date.now()
    const duration = quizData.duration || Math.floor((Date.now() - startTime) / 1000)
    const correct = answers.filter((item) => item.is_correct).length
    const accuracy = answers.length === 0 ? 0 : Math.round((correct / answers.length) * 100)
    const score = Math.round(correct * 5)
    
    this.setData({
      answers,
      startTime,
      duration,
      accuracy,
      score,
      formattedDuration: this.formatDuration(duration)
    })
  },

  formatDuration(secs) {
    const mins = Math.floor(secs / 60)
    const remainSecs = secs % 60
    return `${mins}分${remainSecs}秒`
  },

  goHome() {
    app.globalData.quizData = null
    wx.switchTab({ url: '/pages/index/index' })
  },

  redoQuiz() {
    app.globalData.quizData = null
    wx.navigateBack()
  }
})
