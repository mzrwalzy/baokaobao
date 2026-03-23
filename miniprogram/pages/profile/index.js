const api = require('../../api/index.js')
const app = getApp()

Page({
  data: {
    userInfo: {},
    stats: {}
  },

  onLoad() {
    if (!app.globalData.token) {
      wx.reLaunch({ url: '/pages/login/index' })
      return
    }
    this.setData({ userInfo: app.globalData.userInfo || {} })
    this.loadData()
  },

  async loadData() {
    try {
      const [profile, stats] = await Promise.all([
        api.getProfile(),
        api.getStats()
      ])
      this.setData({ 
        userInfo: profile,
        stats
      })
    } catch (e) {
      console.error(e)
    }
  },

  goPage(e) {
    const url = e.currentTarget.dataset.url
    wx.showToast({ title: '功能开发中', icon: 'none' })
  },

  handleLogout() {
    wx.showModal({
      title: '提示',
      content: '确定要退出登录吗？',
      success: async (res) => {
        if (res.confirm) {
          try {
            await api.logout()
          } catch (e) {}
          app.clearUserData()
          wx.reLaunch({ url: '/pages/login/index' })
        }
      }
    })
  }
})
