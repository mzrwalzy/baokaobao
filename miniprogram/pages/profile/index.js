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

  onQuickLogin() {
    wx.showLoading({ title: '登录中...' })
    wx.login({
      success: (res) => {
        if (res.code) {
          api.loginByWechat(res.code).then(result => {
            wx.hideLoading()
            const app = getApp()
            app.setUserData(result.token, result.user)
            this.setData({ 
              isLoggedIn: true,
              userInfo: result.user
            })
            this.loadData()
            wx.showToast({ title: '登录成功', icon: 'success' })
          }).catch(err => {
            wx.hideLoading()
            wx.showToast({ title: err.message || '登录失败', icon: 'none' })
          })
        } else {
          wx.hideLoading()
          wx.showToast({ title: '获取登录凭证失败', icon: 'none' })
        }
      },
      fail: () => {
        wx.hideLoading()
        wx.showToast({ title: '微信登录失败', icon: 'none' })
      }
    })
  },

  doLogin(code) {
    wx.showLoading({ title: '登录中...' })
    api.loginByWechat(code).then(result => {
      wx.hideLoading()
      const app = getApp()
      app.setUserData(result.token, result.user)
      this.setData({ 
        isLoggedIn: true,
        userInfo: result.user
      })
      this.loadData()
      wx.showToast({ title: '登录成功', icon: 'success' })
    }).catch(err => {
      wx.hideLoading()
      wx.showToast({ title: err.message || '登录失败', icon: 'none' })
    })
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
