const api = require('../../api/index.js')

Page({
  data: {
    questions: [],
    currentIndex: 0,
    currentQuestion: {},
    selectedAnswer: '',
    showResult: false,
    answers: [],
    startTime: 0,
    bankId: 0,
    bankName: '',
    mode: 'exam',
    modeText: '考试模式',
    totalQuestions: 0,
    progress: 0,
    typeText: '题目',
    submitting: false,
    page: 1,
    pageSize: 1,
    hasMore: true,
    loadingNext: false,
    countdown: 0,
    countdownText: '',
    practiceSubmitted: false,
    primaryButtonText: '查看答案'
  },

  timer: null,

  onLoad(options) {
    const bankId = Number(options.bank_id) || 0
    const count = Number(options.count) || 10
    const mode = options.mode || 'exam'
    const bankName = decodeURIComponent(options.name || '题库练习')

    this.setData({
      bankId,
      bankName,
      mode,
      modeText: this.getModeText(mode),
      startTime: Date.now()
    })

    if (mode === 'practice') {
      this.initPracticeMode(bankId)
      return
    }

    this.loadExamQuestions(bankId, count, mode)
  },

  onUnload() {
    this.clearTimer()
  },

  getModeText(mode) {
    switch (mode) {
      case 'practice':
        return '练习模式'
      case 'exam':
        return '考试模式'
      case 'memorize':
        return '背题模式'
      default:
        return '答题模式'
    }
  },

  getTypeText(type) {
    switch (type) {
      case 'single':
        return '单选题'
      case 'multiple':
        return '多选题'
      case 'truefalse':
        return '判断题'
      default:
        return '题目'
    }
  },

  formatClock(totalSeconds = 0) {
    const safeSeconds = Math.max(0, totalSeconds)
    const mins = Math.floor(safeSeconds / 60)
    const secs = safeSeconds % 60
    return `${String(mins).padStart(2, '0')}:${String(secs).padStart(2, '0')}`
  },

  clearTimer() {
    if (this.timer) {
      clearInterval(this.timer)
      this.timer = null
    }
  },

  startCountdown() {
    this.clearTimer()
    this.timer = setInterval(() => {
      const nextCountdown = this.data.countdown - 1
      if (nextCountdown <= 0) {
        this.clearTimer()
        this.setData({
          countdown: 0,
          countdownText: '00:00'
        })
        wx.showToast({ title: '考试时间到，已自动交卷', icon: 'none' })
        this.submitQuizSession()
        return
      }

      this.setData({
        countdown: nextCountdown,
        countdownText: this.formatClock(nextCountdown)
      })
    }, 1000)
  },

  async initPracticeMode(bankId) {
    getApp().globalData.quizData = null
    this.setData({
      questions: [],
      currentIndex: 0,
      currentQuestion: {},
      selectedAnswer: '',
      showResult: false,
      answers: [],
      totalQuestions: 0,
      progress: 0,
      typeText: '题目',
      page: 1,
      hasMore: true,
      practiceSubmitted: false,
      primaryButtonText: '查看答案'
    })

    await this.loadPracticeQuestion(bankId, 1)
  },

  async loadPracticeQuestion(bankId, page) {
    if (this.data.loadingNext) return

    this.setData({ loadingNext: true })
    try {
      const pageData = await api.getQuestions({
        bank_id: bankId,
        page,
        page_size: 1
      })
      const list = (pageData.list || []).map((question) => this.normalizeQuestion(question))
      const questions = page === 1 ? list : [...this.data.questions, ...list]

      this.setData({
        questions,
        page,
        hasMore: page < (pageData.total_pages || page),
        totalQuestions: pageData.total || questions.length
      })

      if (page === 1 || this.data.currentIndex >= questions.length) {
        this.setData({ currentIndex: questions.length - 1 })
      }

      this.syncQuestionState()
      this.resetQuestionView()
    } catch (e) {
      wx.showToast({ title: e.message || '加载失败', icon: 'none' })
    } finally {
      this.setData({ loadingNext: false })
    }
  },

  async loadExamQuestions(bankId, count, mode) {
    try {
      const questions = (await api.getRandomQuestions(bankId, count)).map((question) => this.normalizeQuestion(question))
      const totalQuestions = questions.length
      const countdown = mode === 'exam' ? Math.max(totalQuestions * 60, 60) : 0

      getApp().globalData.quizData = null
      this.setData({
        questions,
        totalQuestions,
        countdown,
        countdownText: this.formatClock(countdown),
        showResult: mode === 'memorize'
      })
      this.syncQuestionState()

      if (mode === 'exam') {
        this.startCountdown()
      }
    } catch (e) {
      wx.showToast({ title: e.message || '加载失败', icon: 'none' })
    }
  },

  normalizeQuestion(question) {
    const normalizedType = this.normalizeQuestionType(question.type)
    return {
      ...question,
      type: normalizedType,
      imageList: this.parseImageList(question.images),
      options: (question.options || []).map((option) => ({
        ...option,
        option_image: option.option_image || '',
        isSelected: false,
        isCorrectOption: false,
        isWrongSelection: false
      }))
    }
  },

  normalizeQuestionType(type) {
    switch (type) {
      case '单选':
        return 'single'
      case '多选':
        return 'multiple'
      case '判断':
        return 'truefalse'
      default:
        return type || 'single'
    }
  },

  parseImageList(images) {
    if (!images) return []
    if (Array.isArray(images)) return images.filter(Boolean)

    const raw = String(images).trim()
    if (!raw) return []

    if (raw.startsWith('[')) {
      try {
        const parsed = JSON.parse(raw)
        if (Array.isArray(parsed)) {
          return parsed.map((item) => String(item).trim()).filter(Boolean)
        }
      } catch (e) {}
    }

    return raw
      .split(/[\n,]/)
      .map((item) => item.trim())
      .filter(Boolean)
  },

  previewImages(e) {
    const current = e.currentTarget.dataset.current
    const urls = e.currentTarget.dataset.urls || []
    if (!current || !urls.length) return

    wx.previewImage({
      current,
      urls
    })
  },

  syncQuestionState() {
    const currentQuestion = this.data.questions[this.data.currentIndex] || {}
    const totalQuestions = this.data.totalQuestions || this.data.questions.length
    const progress = totalQuestions > 0
      ? ((this.data.currentIndex + 1) / totalQuestions) * 100
      : 0

    this.setData({
      currentQuestion,
      totalQuestions,
      progress,
      typeText: this.getTypeText(currentQuestion.type)
    })
    this.updatePrimaryButtonText()
  },

  updatePrimaryButtonText() {
    let primaryButtonText = '交卷'

    if (this.data.mode === 'practice') {
      primaryButtonText = this.data.practiceSubmitted
        ? (this.data.currentIndex < this.data.totalQuestions - 1 || this.data.hasMore ? '下一题' : '完成练习')
        : '查看答案'
    } else if (this.data.currentIndex < this.data.totalQuestions - 1) {
      primaryButtonText = '下一题'
    }

    this.setData({ primaryButtonText })
  },

  resetQuestionView() {
    this.setData({
      selectedAnswer: '',
      showResult: this.data.mode === 'memorize',
      practiceSubmitted: false
    })
    this.updateOptionStates('')
    this.updatePrimaryButtonText()
  },

  getNextAnswerValue(question, key) {
    if (question.type !== 'multiple') {
      return key
    }

    const current = this.data.selectedAnswer ? this.data.selectedAnswer.split('') : []
    const next = current.includes(key)
      ? current.filter((item) => item !== key)
      : [...current, key]

    return next.sort().join('')
  },

  normalizeAnswer(answer) {
    return (answer || '').split('').sort().join('')
  },

  isAnswerCorrect(question, answerValue) {
    if (!answerValue) return false

    if (question.type === 'multiple') {
      return this.normalizeAnswer(question.answer) === this.normalizeAnswer(answerValue)
    }
    return question.answer === answerValue
  },

  buildAnswerRecord(question, answerValue) {
    return {
      question_id: question.id,
      answer: answerValue,
      is_correct: this.isAnswerCorrect(question, answerValue),
      correct_answer: question.answer,
      analysis: question.analysis || ''
    }
  },

  upsertAnswerRecord(record) {
    const answers = this.data.answers.filter((item) => item.question_id !== record.question_id)
    answers.push(record)
    this.setData({ answers })
  },

  selectOption(e) {
    if (this.data.mode === 'practice' && this.data.showResult) return

    const key = e.currentTarget.dataset.key
    const currentQuestion = this.data.currentQuestion
    const answerValue = this.getNextAnswerValue(currentQuestion, key)

    this.setData({ selectedAnswer: answerValue })
    this.updateOptionStates(answerValue)

    if (this.data.mode !== 'practice' && answerValue) {
      this.upsertAnswerRecord(this.buildAnswerRecord(currentQuestion, answerValue))
    }
  },

  updateOptionStates(answerValue = '') {
    const questions = [...this.data.questions]
    const current = questions[this.data.currentIndex]
    if (!current || !current.options) return

    const selectedSet = new Set((answerValue || '').split('').filter(Boolean))
    const correctSet = new Set(((current.answer || '') + '').split('').filter(Boolean))

    current.options = current.options.map((option) => ({
      ...option,
      isSelected: selectedSet.has(option.option_key),
      isCorrectOption: this.data.showResult && correctSet.has(option.option_key),
      isWrongSelection: this.data.showResult && selectedSet.has(option.option_key) && !correctSet.has(option.option_key)
    }))

    questions[this.data.currentIndex] = current
    this.setData({
      questions,
      currentQuestion: current
    })
  },

  loadCurrentAnswer() {
    const currentQuestion = this.data.currentQuestion
    const saved = this.data.answers.find((item) => item.question_id === currentQuestion.id)
    const selectedAnswer = saved ? saved.answer : ''

    this.setData({ selectedAnswer })
    this.updateOptionStates(selectedAnswer)
  },

  prevQuestion() {
    if (this.data.currentIndex <= 0) return

    this.setData({ currentIndex: this.data.currentIndex - 1 })
    this.syncQuestionState()
    this.loadCurrentAnswer()

    if (this.data.mode === 'practice') {
      const record = this.data.answers.find((item) => item.question_id === this.data.currentQuestion.id)
      this.setData({
        showResult: !!record,
        practiceSubmitted: !!record
      })
      this.updateOptionStates(record ? record.answer : '')
      this.updatePrimaryButtonText()
    }
  },

  async goNextQuestion() {
    if (this.data.currentIndex < this.data.questions.length - 1) {
      this.setData({ currentIndex: this.data.currentIndex + 1 })
      this.syncQuestionState()
      this.loadCurrentAnswer()

      if (this.data.mode === 'practice') {
        const record = this.data.answers.find((item) => item.question_id === this.data.currentQuestion.id)
        this.setData({
          showResult: !!record,
          practiceSubmitted: !!record
        })
        this.updateOptionStates(record ? record.answer : '')
        this.updatePrimaryButtonText()
      }
      return
    }

    if (this.data.mode === 'practice' && this.data.hasMore) {
      await this.loadPracticeQuestion(this.data.bankId, this.data.page + 1)
      this.setData({ currentIndex: this.data.questions.length - 1 })
      this.syncQuestionState()
      this.resetQuestionView()
      return
    }

    await this.submitQuizSession()
  },

  async submitPracticeQuestion() {
    if (!this.data.selectedAnswer) {
      wx.showToast({ title: '请先选择答案', icon: 'none' })
      return
    }

    wx.showLoading({ title: '校验答案中...' })
    try {
      const result = await api.submitAnswer({
        question_id: this.data.currentQuestion.id,
        my_answer: this.data.selectedAnswer
      })

      this.upsertAnswerRecord({
        question_id: this.data.currentQuestion.id,
        answer: this.data.selectedAnswer,
        is_correct: result.is_correct,
        correct_answer: result.correct_answer,
        analysis: result.analysis || ''
      })

      this.setData({
        showResult: true,
        practiceSubmitted: true
      })
      this.updateOptionStates(this.data.selectedAnswer)
      this.updatePrimaryButtonText()
    } catch (e) {
      wx.showToast({ title: e.message || '提交失败', icon: 'none' })
    } finally {
      wx.hideLoading()
    }
  },

  handleBack() {
    wx.showModal({
      title: '提示',
      content: '确定要退出答题吗？',
      success: (res) => {
        if (res.confirm) {
          this.clearTimer()
          wx.navigateBack()
        }
      }
    })
  },

  async handleSubmit() {
    if (this.data.mode === 'practice') {
      if (!this.data.practiceSubmitted) {
        await this.submitPracticeQuestion()
        return
      }
      await this.goNextQuestion()
      return
    }

    if (this.data.submitting) return

    if (this.data.currentIndex < this.data.totalQuestions - 1) {
      await this.goNextQuestion()
      return
    }

    const unanswered = this.data.totalQuestions - this.data.answers.length
    if (unanswered > 0) {
      wx.showModal({
        title: '提示',
        content: `您还有 ${unanswered} 题未答，确定要交卷吗？`,
        success: async (res) => {
          if (res.confirm) {
            await this.submitQuizSession()
          }
        }
      })
      return
    }

    await this.submitQuizSession()
  },

  async submitQuizSession() {
    if (this.data.submitting) return

    if (this.data.mode === 'memorize') {
      this.goResult()
      return
    }

    this.setData({ submitting: true })
    wx.showLoading({ title: '交卷中...' })
    try {
      const duration = Math.max(1, Math.floor((Date.now() - this.data.startTime) / 1000))

      if (this.data.mode === 'exam') {
        await api.submitExam({
          bank_id: this.data.bankId,
          duration,
          answers: this.data.answers.map((item) => ({
            question_id: item.question_id,
            my_answer: item.answer
          }))
        })
      }

      this.clearTimer()
      getApp().globalData.quizData = {
        answers: this.data.answers,
        startTime: this.data.startTime,
        duration,
        mode: this.data.mode,
        totalQuestions: this.data.totalQuestions
      }
      this.goResult()
    } catch (e) {
      wx.showToast({ title: e.message || '交卷失败', icon: 'none' })
    } finally {
      wx.hideLoading()
      this.setData({ submitting: false })
    }
  },

  goResult() {
    wx.navigateTo({
      url: `/pages/result/index?bank_id=${this.data.bankId}&mode=${this.data.mode}`
    })
  }
})
