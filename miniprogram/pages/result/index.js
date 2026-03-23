const app = getApp()

Page({
  data: {
    bankId: 0,
    answers: [],
    startTime: 0,
    duration: 0
  },

  onLoad(options) {
    this.setData({ bankId: Number(options.bank_id) || 0 })
    
    const quizData = app.globalData.quizData || {}
    const answers = quizData.answers || []
    const startTime = quizData.startTime || Date.now()
    const duration = Math.floor((Date.now() - startTime) / 1000)
    
    this.setData({ answers, startTime, duration })
  },

  get accuracy() {
    if (this.data.answers.length === 0) return 0
    const correct = this.data.answers.filter(a => a.is_correct).length
    return Math.round((correct / this.data.answers.length) * 100)
  },

  get score() {
    const correct = this.data.answers.filter(a => a.is_correct).length
    return Math.round(correct * 10)
  },

  get formattedDuration() {
    const secs = this.data.duration
    const mins = Math.floor(secs / 60)
    const remainSecs = secs % 60
    return `${mins}分${remainSecs}秒`
  },

  goHome() {
    wx.switchTab({ url: '/pages/index/index' })
  },

  redoQuiz() {
    wx.navigateBack()
  }
})
