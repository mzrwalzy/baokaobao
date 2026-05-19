const api = require('../../api/index.js')

Page({
  data: {
    loading: false
  },

  handleWxLogin() {
    if (this.data.loading) return
    
    this.setData({ loading: true })
    
    wx.login({
      provider: 'weixin',
      success: async (res) => {
        try {
          const result = await api.loginByWechat(res.code)
          const app = getApp()
          app.setUserData(result.token, result.user)
          wx.switchTab({ url: '/pages/index/index' })
        } catch (e) {
          wx.showToast({ title: e.message || '登录失败', icon: 'none' })
        } finally {
          this.setData({ loading: false })
        }
      },
      fail: (err) => {
        wx.showToast({ title: '微信登录失败', icon: 'none' })
        this.setData({ loading: false })
      }
    })
  },

  openAgreement(e) {
    const type = e.currentTarget.dataset.type
    wx.showToast({ title: `查看${type === 'user' ? '用户协议' : '隐私政策'}`, icon: 'none' })
  }
})
