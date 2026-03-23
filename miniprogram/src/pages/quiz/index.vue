<template>
  <view class="quiz-page">
    <view class="header">
      <view class="back-btn" @click="handleBack">←</view>
      <view class="progress-info">{{ currentIndex + 1 }}/{{ totalQuestions }}</view>
      <view class="submit-btn" @click="handleSubmit">交卷</view>
    </view>

    <view class="progress-bar">
      <view class="progress-fill" :style="{ width: progress + '%' }"></view>
    </view>

    <scroll-view scroll-y class="content">
      <view class="question-card">
        <view class="question-type">
          {{ typeText }}
        </view>
        <view class="question-content">
          {{ currentQuestion?.content || '' }}
        </view>
      </view>

      <view class="options-list">
        <view 
          v-for="(option, index) in currentQuestion?.options || []" 
          :key="option.id"
          class="option-item"
          :class="{ selected: selectedAnswer === option.option_key, correct: showResult && option.option_key === currentQuestion?.answer, wrong: showResult && selectedAnswer === option.option_key && option.option_key !== currentQuestion?.answer }"
          @click="selectOption(option.option_key)"
        >
          <view class="option-key">{{ option.option_key }}</view>
          <view class="option-value">{{ option.option_value }}</view>
        </view>
      </view>

      <view class="analysis-card" v-if="showResult && currentQuestion?.analysis">
        <view class="analysis-title">📝 答案解析</view>
        <view class="analysis-content">{{ currentQuestion.analysis }}</view>
      </view>
    </scroll-view>

    <view class="footer">
      <button 
        class="btn-nav" 
        :disabled="currentIndex === 0"
        @click="prevQuestion"
      >上一题</button>
      <button 
        v-if="currentIndex < totalQuestions - 1"
        class="btn-nav primary"
        @click="nextQuestion"
      >下一题</button>
      <button 
        v-else
        class="btn-nav primary"
        @click="handleSubmit"
      >交卷</button>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { getRandomQuestions, submitAnswer } from '@/api'
import { useQuizStore } from '@/stores/quiz'

const quizStore = useQuizStore()
const currentIndex = ref(0)
const selectedAnswer = ref('')
const showResult = ref(false)

onMounted(() => {
  const pages = getCurrentPages()
  const current = pages[pages.length - 1]
  const options = (current as any).options || {}
  
  const bankId = Number(options.bank_id) || 0
  const count = Number(options.count) || 10
  const mode = options.mode || 'exam'
  
  loadQuestions(bankId, count, mode)
})

const loadQuestions = async (bankId: number, count: number, mode: string) => {
  try {
    const questions = await getRandomQuestions(bankId, count)
    quizStore.initQuiz(questions, bankId, mode as any)
    showResult.value = mode === 'memorize'
  } catch (e: any) {
    uni.showToast({ title: e.message || '加载失败', icon: 'none' })
  }
}

const currentQuestion = computed(() => quizStore.currentQuestion)
const totalQuestions = computed(() => quizStore.totalQuestions)
const progress = computed(() => quizStore.progress * 100)

const typeText = computed(() => {
  const type = currentQuestion.value?.type
  switch (type) {
    case 'single': return '单选题'
    case 'multiple': return '多选题'
    case 'truefalse': return '判断题'
    default: return '题目'
  }
})

const selectOption = (key: string) => {
  if (showResult.value) return
  
  selectedAnswer.value = key
  
  const isCorrect = key === currentQuestion.value?.answer
  quizStore.submitAnswer(key, isCorrect)
  
  if (showResult.value) return
  
  submitAnswer({
    question_id: currentQuestion.value?.id || 0,
    answer: key,
    is_correct: isCorrect
  }).catch(() => {})
}

const prevQuestion = () => {
  if (currentIndex.value > 0) {
    currentIndex.value--
    loadCurrentAnswer()
  }
}

const nextQuestion = () => {
  if (currentIndex.value < totalQuestions.value - 1) {
    currentIndex.value++
    loadCurrentAnswer()
  }
}

const loadCurrentAnswer = () => {
  const saved = quizStore.getAnswer(currentQuestion.value?.id || 0)
  selectedAnswer.value = saved?.answer || ''
}

const handleBack = () => {
  uni.showModal({
    title: '提示',
    content: '确定要退出答题吗？',
    success: (res) => {
      if (res.confirm) {
        uni.navigateBack()
      }
    }
  })
}

const handleSubmit = () => {
  if (quizStore.answeredCount < totalQuestions.value) {
    uni.showModal({
      title: '提示',
      content: `您还有 ${totalQuestions.value - quizStore.answeredCount} 题未答，确定要交卷吗？`,
      success: (res) => {
        if (res.confirm) {
          submitAndGoResult()
        }
      }
    })
  } else {
    submitAndGoResult()
  }
}

const submitAndGoResult = () => {
  uni.navigateTo({
    url: `/pages/result/index?bank_id=${quizStore.bankId}`
  })
}
</script>

<style scoped>
.quiz-page {
  min-height: 100vh;
  background: #F5F5F5;
  display: flex;
  flex-direction: column;
}

.header {
  background: #409EFF;
  padding: 20rpx 30rpx;
  padding-top: calc(20rpx + var(--status-bar-height));
  display: flex;
  align-items: center;
  justify-content: space-between;
  color: #fff;
}

.back-btn {
  font-size: 40rpx;
}

.progress-info {
  font-size: 28rpx;
}

.submit-btn {
  font-size: 28rpx;
  padding: 10rpx 20rpx;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 20rpx;
}

.progress-bar {
  height: 8rpx;
  background: rgba(255, 255, 255, 0.3);
}

.progress-fill {
  height: 100%;
  background: #fff;
  transition: width 0.3s;
}

.content {
  flex: 1;
  padding: 30rpx;
  padding-bottom: 160rpx;
}

.question-card {
  background: #fff;
  border-radius: 20rpx;
  padding: 40rpx 30rpx;
  margin-bottom: 30rpx;
}

.question-type {
  display: inline-block;
  padding: 8rpx 20rpx;
  background: #409EFF;
  color: #fff;
  border-radius: 20rpx;
  font-size: 24rpx;
  margin-bottom: 24rpx;
}

.question-content {
  font-size: 32rpx;
  color: #333;
  line-height: 1.8;
}

.options-list {
  display: flex;
  flex-direction: column;
  gap: 20rpx;
}

.option-item {
  background: #fff;
  border-radius: 16rpx;
  padding: 30rpx;
  display: flex;
  align-items: center;
  gap: 24rpx;
  border: 4rpx solid transparent;
  transition: all 0.3s;
}

.option-item.selected {
  border-color: #409EFF;
  background: rgba(64, 158, 255, 0.1);
}

.option-item.correct {
  border-color: #67C23A;
  background: rgba(103, 194, 58, 0.1);
}

.option-item.wrong {
  border-color: #F56C6C;
  background: rgba(245, 108, 108, 0.1);
}

.option-key {
  width: 60rpx;
  height: 60rpx;
  background: #f5f5f5;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 28rpx;
  font-weight: bold;
  color: #666;
}

.option-item.selected .option-key {
  background: #409EFF;
  color: #fff;
}

.option-value {
  flex: 1;
  font-size: 28rpx;
  color: #333;
}

.analysis-card {
  background: #fff;
  border-radius: 16rpx;
  padding: 30rpx;
  margin-top: 30rpx;
}

.analysis-title {
  font-size: 28rpx;
  font-weight: bold;
  color: #333;
  margin-bottom: 16rpx;
}

.analysis-content {
  font-size: 28rpx;
  color: #666;
  line-height: 1.8;
}

.footer {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  padding: 30rpx;
  padding-bottom: calc(30rpx + env(safe-area-inset-bottom));
  background: #fff;
  display: flex;
  gap: 20rpx;
  box-shadow: 0 -4rpx 20rpx rgba(0, 0, 0, 0.1);
}

.btn-nav {
  flex: 1;
  height: 88rpx;
  background: #f5f5f5;
  color: #666;
  border: none;
  border-radius: 44rpx;
  font-size: 30rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.btn-nav.primary {
  background: #409EFF;
  color: #fff;
}

.btn-nav[disabled] {
  opacity: 0.5;
}
</style>
