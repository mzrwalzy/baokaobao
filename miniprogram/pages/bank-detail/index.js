const api = require('../../api/index.js')
const app = getApp()

Page({
  data: {
    bankId: 0,
    bankName: '',
    bankInfo: {},
    accessStatus: null,
    actionText: '开始刷题',
    mode: 'exam'
  },

  onLoad(options) {
    this.setData({
      bankId: Number(options.id) || 0,
      bankName: decodeURIComponent(options.name || '题库')
    })
    this.loadBankInfo()
  },

  onShow() {
    if (this.data.bankId) {
      this.loadBankInfo()
    }
  },

  async loadBankInfo() {
    try {
      const tasks = [api.getQuestionBankDetail(this.data.bankId)]
      if (app.globalData.token) {
        tasks.push(api.getQuestionBankAccessStatus(this.data.bankId))
      }

      const [info, accessStatus] = await Promise.all(tasks)
      const actionText = this.getActionText(info, accessStatus)
      this.setData({
        bankInfo: info,
        accessStatus: accessStatus || null,
        actionText
      })
    } catch (e) {
      wx.showToast({ title: e.message || '加载失败', icon: 'none' })
    }
  },

  getActionText(bankInfo, accessStatus) {
    if (!app.globalData.token) {
      return '请先登录后开始刷题'
    }
    if (accessStatus && accessStatus.requires_purchase && !accessStatus.has_access) {
      return '当前账号未开通权限'
    }
    if (bankInfo && bankInfo.price > 0 && !accessStatus) {
      return '正在校验题库权限'
    }
    return '开始刷题'
  },

  selectMode(e) {
    this.setData({ mode: e.currentTarget.dataset.mode })
  },

  goBack() {
    wx.navigateBack()
  },

  startQuiz() {
    if (!app.globalData.token) {
      wx.showToast({ title: '请先登录后开始刷题', icon: 'none' })
      setTimeout(() => {
        wx.switchTab({ url: '/pages/profile/index' })
      }, 1200)
      return
    }

    if (this.data.accessStatus && this.data.accessStatus.requires_purchase && !this.data.accessStatus.has_access) {
      wx.showToast({ title: '当前题库尚未开通权限', icon: 'none' })
      return
    }

    wx.navigateTo({
      url: `/pages/quiz/index?bank_id=${this.data.bankId}&mode=${this.data.mode}&count=10&name=${encodeURIComponent(this.data.bankName)}`
    })
  },

  goWrongQuestions() {
    wx.showToast({ title: '错题功能开发中', icon: 'none' })
  }
})
