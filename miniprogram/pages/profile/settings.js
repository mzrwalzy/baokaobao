const SETTINGS_KEY = 'app_settings'

const DEFAULT_SETTINGS = {
  soundEnabled: true,
  vibrationEnabled: true
}

Page({
  data: {
    settings: { ...DEFAULT_SETTINGS },
    cacheSize: '0 KB',
    version: '1.0.0'
  },

  onLoad() {
    this.loadSettings()
    this.calcCacheSize()
    this.loadVersion()
  },

  loadSettings() {
    try {
      const saved = wx.getStorageSync(SETTINGS_KEY)
      if (saved) {
        this.setData({ settings: { ...DEFAULT_SETTINGS, ...saved } })
      }
    } catch (e) {
      console.error(e)
    }
  },

  saveSettings(settings) {
    try {
      wx.setStorageSync(SETTINGS_KEY, settings)
    } catch (e) {
      console.error(e)
    }
  },

  toggleSound(e) {
    const settings = { ...this.data.settings, soundEnabled: e.detail.value }
    this.setData({ settings })
    this.saveSettings(settings)
  },

  toggleVibration(e) {
    const settings = { ...this.data.settings, vibrationEnabled: e.detail.value }
    this.setData({ settings })
    this.saveSettings(settings)
  },

  calcCacheSize() {
    try {
      const info = wx.getStorageInfoSync()
      const size = info.currentSize
      if (size < 1024) {
        this.setData({ cacheSize: size + ' KB' })
      } else {
        this.setData({ cacheSize: (size / 1024).toFixed(2) + ' MB' })
      }
    } catch (e) {
      console.error(e)
    }
  },

  loadVersion() {
    const info = wx.getAccountInfoSync()
    this.setData({ version: info.miniProgram.version || '1.0.0' })
  },

  clearCache() {
    wx.showModal({
      title: '提示',
      content: '确定要清除缓存吗？登录状态不会被清除',
      success: (res) => {
        if (res.confirm) {
          try {
            // 保留登录态，只清除应用缓存
            const token = wx.getStorageSync('token')
            const userInfo = wx.getStorageSync('userInfo')
            const settings = wx.getStorageSync(SETTINGS_KEY)

            wx.clearStorageSync()

            // 恢复登录态和设置
            if (token) wx.setStorageSync('token', token)
            if (userInfo) wx.setStorageSync('userInfo', userInfo)
            if (settings) wx.setStorageSync(SETTINGS_KEY, settings)

            this.setData({ settings: { ...DEFAULT_SETTINGS, ...settings }, cacheSize: '0 KB' })
            wx.showToast({ title: '缓存已清除', icon: 'success' })
          } catch (e) {
            wx.showToast({ title: '清除失败', icon: 'none' })
          }
        }
      }
    })
  },

  goBack() {
    wx.navigateBack()
  }
})
