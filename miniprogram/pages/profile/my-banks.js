const api = require('../../api/index.js')
const { FILE_BASE_URL } = require('../../utils/request.js')

const normalizeImage = (url) => {
  if (!url) return ''
  if (/^(https?:)?\/\//.test(url)) return url
  return `${FILE_BASE_URL}${url}`
}

Page({
  data: {
    list: [],
    loading: false
  },

  onLoad() {
    this.loadList()
  },

  async onPullDownRefresh() {
    await this.loadList()
    wx.stopPullDownRefresh()
  },

  goBack() {
    wx.navigateBack()
  },

  async loadList() {
    this.setData({ loading: true })
    try {
      const res = await api.getPurchasedBanks()
      const list = (res || []).map((item) => ({
        ...item,
        cover_image_url: normalizeImage(item.cover_image),
        price_text: Number(item.price || 0).toFixed(2)
      }))
      this.setData({ list })
    } catch (e) {
      console.error(e)
    } finally {
      this.setData({ loading: false })
    }
  },

  goBank(e) {
    const id = e.currentTarget.dataset.id
    if (!id) return
    wx.navigateTo({ url: `/pages/bank-detail/index?id=${id}` })
  }
})
