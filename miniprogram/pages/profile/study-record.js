const api = require('../../api/index.js')

const formatDate = (value) => {
  if (!value) return ''
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return ''
  const pad = (n) => String(n).padStart(2, '0')
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())} ${pad(date.getHours())}:${pad(date.getMinutes())}`
}

Page({
  data: {
    list: [],
    page: 1,
    pageSize: 20,
    total: 0,
    loading: false,
    finished: false
  },

  onLoad() {
    this.loadList(true)
  },

  onReachBottom() {
    if (!this.data.finished && !this.data.loading) {
      this.loadList(false)
    }
  },

  async onPullDownRefresh() {
    await this.loadList(true)
    wx.stopPullDownRefresh()
  },

  goBack() {
    wx.navigateBack()
  },

  async loadList(reset) {
    if (this.data.loading) return
    const page = reset ? 1 : this.data.page
    this.setData({ loading: true })

    try {
      const res = await api.getQuizHistory({ page, page_size: this.data.pageSize })
      const list = (res.list || []).map((item) => {
        const question = item.question || {}
        return {
          ...item,
          title: question.title || `题目 ${item.question_id}`,
          content: question.content || '',
          is_correct_text: item.is_correct ? '答对' : '答错',
          answered_at_text: formatDate(item.answered_at)
        }
      })
      const nextList = reset ? list : this.data.list.concat(list)
      this.setData({
        list: nextList,
        total: res.total || 0,
        page: page + 1,
        finished: nextList.length >= (res.total || 0)
      })
    } catch (e) {
      console.error(e)
    } finally {
      this.setData({ loading: false })
    }
  }
})
