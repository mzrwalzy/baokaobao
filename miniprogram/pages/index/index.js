const api = require('../../api/index.js')
const app = getApp()

Page({
  data: {
    banks: [],
    loading: false,
    page: 1,
    pageSize: 10,
    noMore: false
  },

  onLoad() {
    if (!app.globalData.token) {
      wx.reLaunch({ url: '/pages/login/index' })
      return
    }
    this.loadBanks()
  },

  async loadBanks() {
    if (this.data.loading || this.data.noMore) return
    
    this.setData({ loading: true })
    try {
      const res = await api.getQuestionBanks({ page: this.data.page, page_size: this.data.pageSize })
      const list = res.list || []
      
      if (this.data.page === 1) {
        this.setData({ banks: list })
      } else {
        this.setData({ banks: [...this.data.banks, ...list] })
      }
      
      if (list.length < this.data.pageSize) {
        this.setData({ noMore: true })
      } else {
        this.setData({ page: this.data.page + 1 })
      }
    } catch (e) {
      wx.showToast({ title: e.message || '加载失败', icon: 'none' })
    } finally {
      this.setData({ loading: false })
    }
  },

  loadMore() {
    if (!this.data.noMore) {
      this.loadBanks()
    }
  },

  goBankDetail(e) {
    const bank = e.currentTarget.dataset.bank
    wx.navigateTo({ 
      url: `/pages/bank-detail/index?id=${bank.id}&name=${encodeURIComponent(bank.name)}` 
    })
  },

  goSearch() {
    wx.showToast({ title: '搜索功能开发中', icon: 'none' })
  },

  goNotice() {
    wx.showToast({ title: '通知功能开发中', icon: 'none' })
  }
})
