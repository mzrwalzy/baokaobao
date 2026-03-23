const api = require('../../api/index.js')
const app = getApp()

Page({
  data: {
    userInfo: {},
    stats: {},
    isLoggedIn: false
  },

  onLoad() {
    this.setData({ 
      isLoggedIn: !!app.globalData.token,
      userInfo: app.globalData.userInfo || {}
    })
    if (app.globalData.token) {
      this.loadData()
    }
  },

  onShow() {
    const isLoggedIn = !!app.globalData.token
    this.setData({ 
      isLoggedIn,
      userInfo: app.globalData.userInfo || {}
    })
    if (isLoggedIn) {
      this.loadData()
    }
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

  goLogin() {
    wx.navigateTo({ url: '/pages/login/index' })
  },

  goPage(e) {
    if (!app.globalData.token) {
      wx.showModal({
        title: '提示',
        content: '请先登录后操作',
        showCancel: false,
        confirmText: '去登录',
        success: () => {
          wx.navigateTo({ url: '/pages/login/index' })
        }
      })
      return
    }
    const url = e.currentTarget.dataset.url
    wx.showToast({ title: '功能开发中', icon: 'none' })
  },

  handleLogout() {
    if (!app.globalData.token) return
    
    wx.showModal({
      title: '提示',
      content: '确定要退出登录吗？',
      success: async (res) => {
        if (res.confirm) {
          try {
            await api.logout()
          } catch (e) {}
          app.clearUserData()
          this.setData({ 
            isLoggedIn: false,
            userInfo: {},
            stats: {}
          })
        }
      }
    })
  }
})
