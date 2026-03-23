const api = require('../../api/index.js')
const app = getApp()

Page({
  data: {
    tab: 'total',
    list: [],
    loading: false,
    isLoggedIn: false
  },

  onLoad() {
    this.setData({ isLoggedIn: !!app.globalData.token })
    this.loadRanking()
  },

  onShow() {
    this.setData({ isLoggedIn: !!app.globalData.token })
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
      if (e.message && e.message.includes('401')) {
        wx.showModal({
          title: '提示',
          content: '登录后查看完整排行榜，是否前往登录？',
          success: (res) => {
            if (res.confirm) {
              wx.navigateTo({ url: '/pages/login/index' })
            }
          }
        })
      } else {
        wx.showToast({ title: e.message || '加载失败', icon: 'none' })
      }
    } finally {
      this.setData({ loading: false })
    }
  }
})
