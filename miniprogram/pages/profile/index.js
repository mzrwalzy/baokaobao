const api = require('../../api/index.js')
const { FILE_BASE_URL } = require('../../utils/request.js')
const app = getApp()

Page({
  data: {
    userInfo: {},
    stats: {},
    isLoggedIn: false,
    editingNickname: false,
    nicknameDraft: '',
    savingProfile: false
  },

  onLoad() {
    const userInfo = this.normalizeUserInfo(app.globalData.userInfo || {})
    this.setData({ 
      isLoggedIn: !!app.globalData.token,
      userInfo
    })
    if (app.globalData.token) {
      this.loadData()
    }
  },

  onShow() {
    const isLoggedIn = !!app.globalData.token
    const userInfo = this.normalizeUserInfo(app.globalData.userInfo || {})
    this.setData({ 
      isLoggedIn,
      userInfo
    })
    if (isLoggedIn) {
      this.loadData()
    }
  },

  normalizeAvatarUrl(url) {
    if (!url) return ''
    if (/^(https?:)?\/\//.test(url) || url.startsWith('wxfile://')) {
      return url
    }
    return `${FILE_BASE_URL}${url}`
  },

  normalizeUserInfo(userInfo = {}) {
    const nickname = userInfo.nickname || ''
    return {
      ...userInfo,
      nickname,
      avatar_url: this.normalizeAvatarUrl(userInfo.avatar_url || ''),
      avatar_initial: nickname ? nickname.slice(0, 1) : '微'
    }
  },

  syncUserInfo(userInfo = {}) {
    const merged = {
      ...(app.globalData.userInfo || {}),
      ...userInfo
    }
    app.globalData.userInfo = merged
    wx.setStorageSync('userInfo', merged)
  },

  async loadData() {
    try {
      const [profile, stats] = await Promise.all([
        api.getProfile(),
        api.getStats()
      ])
      const userInfo = this.normalizeUserInfo(profile)
      this.syncUserInfo(profile)
      this.setData({ 
        userInfo,
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
            const userInfo = this.normalizeUserInfo(result.user)
            this.setData({ 
              isLoggedIn: true,
              userInfo
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

  onTapNickname() {
    if (!app.globalData.token) {
      this.onQuickLogin()
      return
    }
    this.setData({
      editingNickname: true,
      nicknameDraft: this.data.userInfo.nickname || ''
    })
  },

  onNicknameInput(e) {
    this.setData({
      nicknameDraft: e.detail.value || ''
    })
  },

  onCancelNicknameEdit() {
    this.setData({
      editingNickname: false,
      nicknameDraft: this.data.userInfo.nickname || ''
    })
  },

  async onSaveNickname() {
    const nickname = (this.data.nicknameDraft || '').trim()
    if (!nickname) {
      wx.showToast({ title: '请输入昵称', icon: 'none' })
      return
    }
    if (this.data.savingProfile) return

    this.setData({ savingProfile: true })
    wx.showLoading({ title: '保存中...' })
    try {
      await api.updateProfile({ nickname })
      const userInfo = {
        ...this.data.userInfo,
        nickname
      }
      this.syncUserInfo(userInfo)
      this.setData({
        userInfo,
        editingNickname: false
      })
      wx.showToast({ title: '昵称已更新', icon: 'success' })
    } catch (e) {
      wx.showToast({ title: e.message || '保存失败', icon: 'none' })
    } finally {
      wx.hideLoading()
      this.setData({ savingProfile: false })
    }
  },

  async onChooseAvatar(e) {
    if (!app.globalData.token) {
      this.onQuickLogin()
      return
    }

    const avatarPath = e.detail.avatarUrl
    if (!avatarPath) return

    wx.showLoading({ title: '上传头像中...' })
    try {
      const result = await api.uploadAvatar(avatarPath)
      const userInfo = {
        ...this.data.userInfo,
        avatar_url: this.normalizeAvatarUrl(result.url)
      }
      this.syncUserInfo(userInfo)
      this.setData({ userInfo })
      wx.showToast({ title: '头像已更新', icon: 'success' })
    } catch (e) {
      wx.showToast({ title: e.message || '头像上传失败', icon: 'none' })
    } finally {
      wx.hideLoading()
    }
  },

  goPage(e) {
    if (!app.globalData.token) {
      wx.showToast({ title: '请先登录', icon: 'none' })
      setTimeout(() => {
        wx.switchTab({ url: '/pages/profile/index' })
      }, 1500)
      return
    }
    const url = e.currentTarget.dataset.url
    if (url) {
      wx.navigateTo({ url })
    }
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
          } catch (e) {
            console.error('Logout failed:', e)
          }
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
