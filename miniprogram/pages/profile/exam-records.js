const api = require('../../api/index.js')

const formatDate = (value) => {
  if (!value) return ''
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return ''
  const pad = (n) => String(n).padStart(2, '0')
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())} ${pad(date.getHours())}:${pad(date.getMinutes())}`
}

const formatDuration = (seconds) => {
  const total = Number(seconds || 0)
  const mins = Math.floor(total / 60)
  const secs = total % 60
  if (mins <= 0) return `${secs}秒`
  return `${mins}分${secs}秒`
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
      const res = await api.getExamRecords({ page, page_size: this.data.pageSize })
      const list = (res.list || []).map((item) => {
        const totalQuestion = Number(item.total_question || 0)
        const correctRate = totalQuestion > 0 ? (Number(item.correct_count || 0) * 100) / totalQuestion : 0
        return {
          ...item,
          bank_name: item.bank && item.bank.name ? item.bank.name : `题库 ${item.bank_id}`,
          submitted_at_text: formatDate(item.submitted_at || item.created_at),
          duration_text: formatDuration(item.duration),
          correct_rate_text: correctRate.toFixed(1)
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
