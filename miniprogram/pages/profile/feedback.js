const api = require('../../api/index.js')

Page({
  data: {
    type: 'suggestion',
    content: '',
    contact: '',
    submitting: false
  },

  selectType(e) {
    this.setData({ type: e.currentTarget.dataset.type })
  },

  onContentInput(e) {
    this.setData({ content: e.detail.value })
  },

  onContactInput(e) {
    this.setData({ contact: e.detail.value })
  },

  async submitFeedback() {
    if (this.data.submitting) return

    const content = (this.data.content || '').trim()
    if (!content) {
      wx.showToast({ title: '请填写反馈内容', icon: 'none' })
      return
    }

    this.setData({ submitting: true })
    try {
      await api.createFeedback({
        type: this.data.type,
        content: content,
        contact: (this.data.contact || '').trim()
      })
      wx.showToast({ title: '提交成功', icon: 'success' })
      setTimeout(() => {
        wx.navigateBack()
      }, 1200)
    } catch (e) {
      wx.showToast({ title: e.message || '提交失败', icon: 'none' })
    } finally {
      this.setData({ submitting: false })
    }
  },

  goBack() {
    wx.navigateBack()
  }
})
