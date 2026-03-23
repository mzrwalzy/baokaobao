const api = require('../../api/index.js')

Page({
  data: {
    bankId: 0,
    bankName: '',
    bankInfo: {},
    mode: 'exam'
  },

  onLoad(options) {
    this.setData({
      bankId: Number(options.id) || 0,
      bankName: decodeURIComponent(options.name || '题库')
    })
    this.loadBankInfo()
  },

  async loadBankInfo() {
    try {
      const info = await api.getQuestionBankDetail(this.data.bankId)
      this.setData({ bankInfo: info })
    } catch (e) {
      wx.showToast({ title: e.message || '加载失败', icon: 'none' })
    }
  },

  selectMode(e) {
    this.setData({ mode: e.currentTarget.dataset.mode })
  },

  goBack() {
    wx.navigateBack()
  },

  startQuiz() {
    wx.navigateTo({
      url: `/pages/quiz/index?bank_id=${this.data.bankId}&mode=${this.data.mode}&count=10`
    })
  },

  goWrongQuestions() {
    wx.showToast({ title: '错题功能开发中', icon: 'none' })
  }
})
