App({
  globalData: {
    userInfo: null,
    token: ''
  },
  onLaunch() {
    const token = wx.getStorageSync('token')
    const userInfo = wx.getStorageSync('userInfo')
    if (token) {
      this.globalData.token = token
      this.globalData.userInfo = userInfo
    }
  },
  setUserData(token, userInfo) {
    this.globalData.token = token
    this.globalData.userInfo = userInfo
    wx.setStorageSync('token', token)
    wx.setStorageSync('userInfo', userInfo)
  },
  clearUserData() {
    this.globalData.token = ''
    this.globalData.userInfo = null
    wx.removeStorageSync('token')
    wx.removeStorageSync('userInfo')
  }
})
