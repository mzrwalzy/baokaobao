<template>
  <view class="bank-detail-page">
    <view class="header">
      <view class="back-btn" @click="goBack">←</view>
      <text class="header-title">{{ bankName }}</text>
      <view class="share-btn">分享</view>
    </view>

    <view class="content">
      <view class="bank-header">
        <view class="bank-icon">📚</view>
        <view class="bank-title">{{ bankName }}</view>
        <view class="bank-meta">
          <text>题目数量: {{ bankInfo?.question_count || 0 }}+</text>
          <text>难度: {{ '⭐'.repeat(bankInfo?.difficulty || 1) }}</text>
        </view>
      </view>

      <view class="card description-card">
        <view class="card-title">📖 题目简介</view>
        <text class="description">{{ bankInfo?.description || '暂无简介' }}</text>
      </view>

      <view class="card mode-section">
        <view class="card-title">🎯 刷题模式</view>
        <view class="mode-grid">
          <view 
            class="mode-item" 
            :class="{ active: mode === 'practice' }"
            @click="mode = 'practice'"
          >
            <text class="mode-icon">📝</text>
            <text class="mode-name">练习模式</text>
          </view>
          <view 
            class="mode-item"
            :class="{ active: mode === 'exam' }"
            @click="mode = 'exam'"
          >
            <text class="mode-icon">📋</text>
            <text class="mode-name">考试模式</text>
          </view>
          <view 
            class="mode-item"
            :class="{ active: mode === 'memorize' }"
            @click="mode = 'memorize'"
          >
            <text class="mode-icon">🧠</text>
            <text class="mode-name">背题模式</text>
          </view>
          <view class="mode-item" @click="goWrongQuestions">
            <text class="mode-icon">❌</text>
            <text class="mode-name">错题重做</text>
          </view>
        </view>
      </view>

      <view class="start-section">
        <button class="btn-start" @click="startQuiz">开始刷题</button>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getQuestionBankDetail } from '@/api'

const bankId = ref(0)
const bankName = ref('')
const bankInfo = ref<any>(null)
const mode = ref<'practice' | 'exam' | 'memorize'>('exam')

onMounted(() => {
  const pages = getCurrentPages()
  const current = pages[pages.length - 1]
  const options = (current as any).options || current.$page?.options || {}
  
  bankId.value = Number(options.id) || 0
  bankName.value = decodeURIComponent(options.name || '题库')
  
  loadBankInfo()
})

const loadBankInfo = async () => {
  try {
    bankInfo.value = await getQuestionBankDetail(bankId.value)
  } catch (e: any) {
    uni.showToast({ title: e.message || '加载失败', icon: 'none' })
  }
}

const goBack = () => {
  uni.navigateBack()
}

const startQuiz = () => {
  uni.navigateTo({ 
    url: `/pages/quiz/index?bank_id=${bankId.value}&mode=${mode.value}&count=10` 
  })
}

const goWrongQuestions = () => {
  uni.showToast({ title: '错题功能开发中', icon: 'none' })
}
</script>

<style scoped>
.bank-detail-page {
  min-height: 100vh;
  background: #F5F5F5;
}

.header {
  position: sticky;
  top: 0;
  z-index: 100;
  background: #fff;
  padding: 20rpx 30rpx;
  padding-top: calc(20rpx + var(--status-bar-height));
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1rpx solid #eee;
}

.back-btn {
  font-size: 40rpx;
  color: #333;
  padding: 10rpx;
}

.header-title {
  font-size: 32rpx;
  font-weight: 500;
  color: #333;
}

.share-btn {
  font-size: 28rpx;
  color: #409EFF;
}

.content {
  padding: 30rpx;
  padding-bottom: 200rpx;
}

.bank-header {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 40rpx 0;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 24rpx;
  margin-bottom: 30rpx;
}

.bank-icon {
  font-size: 80rpx;
  margin-bottom: 20rpx;
}

.bank-title {
  font-size: 40rpx;
  font-weight: bold;
  color: #fff;
  margin-bottom: 16rpx;
}

.bank-meta {
  display: flex;
  gap: 30rpx;
  color: rgba(255, 255, 255, 0.8);
  font-size: 26rpx;
}

.card {
  background: #fff;
  border-radius: 20rpx;
  padding: 30rpx;
  margin-bottom: 30rpx;
}

.card-title {
  font-size: 32rpx;
  font-weight: 500;
  color: #333;
  margin-bottom: 20rpx;
}

.description {
  font-size: 28rpx;
  color: #666;
  line-height: 1.8;
}

.mode-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 20rpx;
}

.mode-item {
  background: #f5f5f5;
  border-radius: 16rpx;
  padding: 30rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12rpx;
  transition: all 0.3s;
}

.mode-item.active {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.mode-item.active .mode-name {
  color: #fff;
}

.mode-icon {
  font-size: 48rpx;
}

.mode-name {
  font-size: 28rpx;
  color: #666;
}

.start-section {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  padding: 30rpx;
  padding-bottom: calc(30rpx + env(safe-area-inset-bottom));
  background: #fff;
  box-shadow: 0 -4rpx 20rpx rgba(0, 0, 0, 0.1);
}

.btn-start {
  width: 100%;
  height: 96rpx;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: #fff;
  border: none;
  border-radius: 48rpx;
  font-size: 34rpx;
  font-weight: 500;
  display: flex;
  align-items: center;
  justify-content: center;
}

.btn-start::after {
  border: none;
}
</style>
