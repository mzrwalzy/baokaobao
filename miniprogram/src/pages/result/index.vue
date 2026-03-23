<template>
  <view class="result-page">
    <view class="header-bg">
      <view class="celebration">🎉</view>
      <view class="title">答题完成！</view>
    </view>

    <view class="score-card">
      <view class="score-grid">
        <view class="score-item">
          <view class="score-value">{{ accuracy }}%</view>
          <view class="score-label">正确率</view>
        </view>
        <view class="score-item">
          <view class="score-value">{{ formattedDuration }}</view>
          <view class="score-label">用时</view>
        </view>
        <view class="score-item">
          <view class="score-value">{{ score }}</view>
          <view class="score-label">得分</view>
        </view>
      </view>
    </view>

    <view class="detail-section">
      <view class="section-title">答题详情</view>
      <view class="answer-list">
        <view 
          v-for="(item, index) in quizStore.answers" 
          :key="index"
          class="answer-item"
        >
          <view class="answer-status" :class="{ correct: item.is_correct, wrong: !item.is_correct }">
            {{ item.is_correct ? '✓' : '✗' }}
          </view>
          <view class="answer-text">第{{ index + 1 }}题</view>
          <view class="answer-result">{{ item.is_correct ? '正确' : '错误' }}</view>
        </view>
      </view>
    </view>

    <view class="action-section">
      <button class="btn-action" @click="goHome">返回题库</button>
      <button class="btn-action primary" @click="redoQuiz">重新答题</button>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useQuizStore } from '@/stores/quiz'

const quizStore = useQuizStore()
const bankId = ref(0)

onMounted(() => {
  const pages = getCurrentPages()
  const current = pages[pages.length - 1]
  const options = (current as any).options || {}
  bankId.value = Number(options.bank_id) || 0
})

const accuracy = computed(() => quizStore.accuracy)
const score = computed(() => Math.round(quizStore.correctCount * 10))

const formattedDuration = computed(() => {
  const seconds = quizStore.duration
  const mins = Math.floor(seconds / 60)
  const secs = seconds % 60
  return `${mins}分${secs}秒`
})

const goHome = () => {
  uni.switchTab({ url: '/pages/index/index' })
}

const redoQuiz = () => {
  uni.navigateBack()
}
</script>

<style scoped>
.result-page {
  min-height: 100vh;
  background: linear-gradient(180deg, #667eea 0%, #764ba2 100%);
  padding-bottom: 200rpx;
}

.header-bg {
  padding: 80rpx 0 60rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.celebration {
  font-size: 100rpx;
  margin-bottom: 20rpx;
}

.title {
  font-size: 40rpx;
  font-weight: bold;
  color: #fff;
}

.score-card {
  margin: 0 30rpx;
  background: #fff;
  border-radius: 24rpx;
  padding: 40rpx;
  box-shadow: 0 8rpx 30rpx rgba(0, 0, 0, 0.15);
}

.score-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 20rpx;
}

.score-item {
  text-align: center;
}

.score-value {
  font-size: 48rpx;
  font-weight: bold;
  color: #333;
  margin-bottom: 8rpx;
}

.score-label {
  font-size: 26rpx;
  color: #999;
}

.detail-section {
  margin: 30rpx;
}

.section-title {
  font-size: 32rpx;
  font-weight: 500;
  color: #fff;
  margin-bottom: 20rpx;
}

.answer-list {
  background: #fff;
  border-radius: 20rpx;
  padding: 20rpx 30rpx;
  max-height: 400rpx;
  overflow-y: auto;
}

.answer-item {
  display: flex;
  align-items: center;
  padding: 20rpx 0;
  border-bottom: 1rpx solid #f0f0f0;
}

.answer-item:last-child {
  border-bottom: none;
}

.answer-status {
  width: 48rpx;
  height: 48rpx;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 26rpx;
  margin-right: 20rpx;
}

.answer-status.correct {
  background: rgba(103, 194, 58, 0.1);
  color: #67C23A;
}

.answer-status.wrong {
  background: rgba(245, 108, 108, 0.1);
  color: #F56C6C;
}

.answer-text {
  flex: 1;
  font-size: 28rpx;
  color: #333;
}

.answer-result {
  font-size: 26rpx;
  color: #999;
}

.action-section {
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

.btn-action {
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

.btn-action.primary {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: #fff;
}
</style>
