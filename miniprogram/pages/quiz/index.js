const api = require('../../api/index.js')

Page({
  data: {
    questions: [],
    currentIndex: 0,
    selectedAnswer: '',
    showResult: false,
    answers: [],
    startTime: 0,
    bankId: 0,
    mode: 'exam'
  },

  onLoad(options) {
    const bankId = Number(options.bank_id) || 0
    const count = Number(options.count) || 10
    const mode = options.mode || 'exam'
    
    this.setData({ bankId, mode })
    this.loadQuestions(bankId, count, mode)
  },

  async loadQuestions(bankId, count, mode) {
    try {
      const questions = await api.getRandomQuestions(bankId, count)
      this.setData({ 
        questions, 
        startTime: Date.now(),
        showResult: mode === 'memorize'
      })
    } catch (e) {
      wx.showToast({ title: e.message || '加载失败', icon: 'none' })
    }
  },

  get currentQuestion() {
    return this.data.questions[this.data.currentIndex] || {}
  },

  get totalQuestions() {
    return this.data.questions.length
  },

  get progress() {
    return this.data.questions.length > 0 
      ? ((this.data.currentIndex + 1) / this.data.questions.length) * 100 
      : 0
  },

  get typeText() {
    const type = this.currentQuestion.type
    switch (type) {
      case 'single': return '单选题'
      case 'multiple': return '多选题'
      case 'truefalse': return '判断题'
      default: return '题目'
    }
  },

  selectOption(e) {
    if (this.data.showResult) return
    
    const key = e.currentTarget.dataset.key
    const q = this.currentQuestion
    const isCorrect = key === q.answer
    
    this.setData({ selectedAnswer: key })
    
    const answers = [...this.data.answers]
    answers.push({ question_id: q.id, answer: key, is_correct: isCorrect })
    this.setData({ answers })
    
    if (!this.data.showResult) {
      api.submitAnswer({
        question_id: q.id,
        answer: key,
        is_correct: isCorrect
      }).catch(() => {})
    }
  },

  prevQuestion() {
    if (this.data.currentIndex > 0) {
      const idx = this.data.currentIndex - 1
      this.setData({ currentIndex: idx })
      this.loadCurrentAnswer()
    }
  },

  nextQuestion() {
    if (this.data.currentIndex < this.totalQuestions - 1) {
      const idx = this.data.currentIndex + 1
      this.setData({ currentIndex: idx })
      this.loadCurrentAnswer()
    }
  },

  loadCurrentAnswer() {
    const q = this.currentQuestion
    const saved = this.data.answers.find(a => a.question_id === q.id)
    this.setData({ selectedAnswer: saved ? saved.answer : '' })
  },

  handleBack() {
    wx.showModal({
      title: '提示',
      content: '确定要退出答题吗？',
      success: (res) => {
        if (res.confirm) {
          wx.navigateBack()
        }
      }
    })
  },

  handleSubmit() {
    const unanswered = this.totalQuestions - this.data.answers.length
    if (unanswered > 0) {
      wx.showModal({
        title: '提示',
        content: `您还有 ${unanswered} 题未答，确定要交卷吗？`,
        success: (res) => {
          if (res.confirm) {
            this.goResult()
          }
        }
      })
    } else {
      this.goResult()
    }
  },

  goResult() {
    wx.navigateTo({
      url: `/pages/result/index?bank_id=${this.data.bankId}`
    })
  }
})
