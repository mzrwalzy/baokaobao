import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Question } from '@/api'

export interface QuizAnswer {
  question_id: number
  answer: string
  is_correct: boolean
}

export const useQuizStore = defineStore('quiz', () => {
  const questions = ref<Question[]>([])
  const currentIndex = ref(0)
  const answers = ref<QuizAnswer[]>([])
  const startTime = ref<number>(0)
  const bankId = ref<number>(0)
  const mode = ref<'practice' | 'exam' | 'memorize'>('exam')

  const currentQuestion = computed(() => questions.value[currentIndex.value])
  const progress = computed(() => questions.value.length > 0 ? (currentIndex.value + 1) / questions.value.length : 0)
  const totalQuestions = computed(() => questions.value.length)
  const answeredCount = computed(() => answers.value.length)
  const correctCount = computed(() => answers.value.filter(a => a.is_correct).length)
  const accuracy = computed(() => answeredCount.value > 0 ? Math.round((correctCount.value / answeredCount.value) * 100) : 0)
  const duration = computed(() => Math.floor((Date.now() - startTime.value) / 1000))

  const initQuiz = (qs: Question[], bank: number, m: 'practice' | 'exam' | 'memorize' = 'exam') => {
    questions.value = qs
    currentIndex.value = 0
    answers.value = []
    bankId.value = bank
    mode.value = m
    startTime.value = Date.now()
  }

  const submitAnswer = (answer: string, isCorrect: boolean) => {
    const q = currentQuestion.value
    if (q) {
      answers.value.push({ question_id: q.id, answer, is_correct: isCorrect })
    }
  }

  const nextQuestion = () => {
    if (currentIndex.value < questions.value.length - 1) {
      currentIndex.value++
    }
  }

  const prevQuestion = () => {
    if (currentIndex.value > 0) {
      currentIndex.value--
    }
  }

  const getAnswer = (questionId: number) => answers.value.find(a => a.question_id === questionId)

  const reset = () => {
    questions.value = []
    currentIndex.value = 0
    answers.value = []
    startTime.value = 0
    bankId.value = 0
  }

  return {
    questions,
    currentIndex,
    answers,
    startTime,
    bankId,
    mode,
    currentQuestion,
    progress,
    totalQuestions,
    answeredCount,
    correctCount,
    accuracy,
    duration,
    initQuiz,
    submitAnswer,
    nextQuestion,
    prevQuestion,
    getAnswer,
    reset
  }
})
