const api = require('../../api/index.js')
const app = getApp()

Page({
  data: {
    tab: 'total',
    list: [],
    loading: false
  },

  onLoad() {
    if (!app.globalData.token) {
      wx.reLaunch({ url: '/pages/login/index' })
      return
    }
    this.loadRanking()
  },

  switchTab(e) {
    const tab = e.currentTarget.dataset.tab
    this.setData({ tab, list: [] })
    this.loadRanking()
  },

  async loadRanking() {
    if (this.data.loading) return
    
    this.setData({ loading: true })
    try {
      const data = await api.getRanking(this.data.tab)
      this.setData({ list: data || [] })
    } catch (e) {
      wx.showToast({ title: e.message || '加载失败', icon: 'none' })
    } finally {
      this.setData({ loading: false })
    }
  }
})
