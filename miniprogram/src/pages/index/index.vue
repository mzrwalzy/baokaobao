<template>
  <view class="index-page">
    <view class="header">
      <view class="search-bar" @click="goSearch">
        <text class="search-icon">🔍</text>
        <text class="search-placeholder">搜索题库...</text>
      </view>
      <view class="notice-btn" @click="goNotice">
        <text>🔔</text>
      </view>
    </view>

    <scroll-view scroll-y class="content" @scrolltolower="loadMore">
      <view class="section">
        <view class="section-title">
          <text class="title">热门题库</text>
        </view>
        
        <view class="bank-list">
          <view 
            class="bank-card" 
            v-for="bank in banks" 
            :key="bank.id"
            @click="goBankDetail(bank)"
          >
            <view class="bank-info">
              <view class="bank-icon">📚</view>
              <view class="bank-text">
                <text class="bank-name">{{ bank.name }}</text>
                <text class="bank-meta">题目: {{ bank.question_count }}+  |  难度: {{ '⭐'.repeat(bank.difficulty) }}</text>
              </view>
            </view>
            <view class="bank-action">
              <text class="start-btn">开始刷题</text>
            </view>
          </view>
        </view>

        <view class="loading-tip" v-if="loading">
          <text>加载中...</text>
        </view>
        <view class="no-more" v-if="noMore && banks.length > 0">
          <text>没有更多了</text>
        </view>
      </view>
    </scroll-view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getQuestionBanks, type QuestionBank } from '@/api'
import { useUserStore } from '@/stores/user'

const userStore = useUserStore()
const banks = ref<QuestionBank[]>([])
const loading = ref(false)
const page = ref(1)
const pageSize = 10
const noMore = ref(false)

onMounted(() => {
  if (!userStore.isLoggedIn) {
    uni.reLaunch({ url: '/pages/login/index' })
    return
  }
  loadBanks()
})

const loadBanks = async () => {
  if (loading.value || noMore.value) return
  
  loading.value = true
  try {
    const res = await getQuestionBanks({ page: page.value, page_size: pageSize })
    if (page.value === 1) {
      banks.value = res.list || []
    } else {
      banks.value = [...banks.value, ...(res.list || [])]
    }
    
    if ((res.list || []).length < pageSize) {
      noMore.value = true
    } else {
      page.value++
    }
  } catch (e: any) {
    uni.showToast({ title: e.message || '加载失败', icon: 'none' })
  } finally {
    loading.value = false
  }
}

const loadMore = () => {
  if (!noMore.value) {
    loadBanks()
  }
}

const goBankDetail = (bank: QuestionBank) => {
  uni.navigateTo({ url: `/pages/bank-detail/index?id=${bank.id}&name=${encodeURIComponent(bank.name)}` })
}

const goSearch = () => {
  uni.showToast({ title: '搜索功能开发中', icon: 'none' })
}

const goNotice = () => {
  uni.showToast({ title: '通知功能开发中', icon: 'none' })
}
</script>

<style scoped>
.index-page {
  min-height: 100vh;
  background: #F5F5F5;
}

.header {
  position: sticky;
  top: 0;
  z-index: 100;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20rpx 30rpx;
  padding-top: calc(20rpx + var(--status-bar-height));
  display: flex;
  align-items: center;
  gap: 20rpx;
}

.search-bar {
  flex: 1;
  height: 72rpx;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 36rpx;
  display: flex;
  align-items: center;
  padding: 0 30rpx;
}

.search-icon {
  font-size: 28rpx;
  margin-right: 16rpx;
}

.search-placeholder {
  color: rgba(255, 255, 255, 0.8);
  font-size: 28rpx;
}

.notice-btn {
  width: 72rpx;
  height: 72rpx;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 32rpx;
}

.content {
  height: calc(100vh - 180rpx);
}

.section {
  padding: 30rpx;
}

.section-title {
  margin-bottom: 30rpx;
}

.title {
  font-size: 36rpx;
  font-weight: bold;
  color: #333;
}

.bank-list {
  display: flex;
  flex-direction: column;
  gap: 24rpx;
}

.bank-card {
  background: #fff;
  border-radius: 20rpx;
  padding: 30rpx;
  display: flex;
  align-items: center;
  justify-content: space-between;
  box-shadow: 0 4rpx 20rpx rgba(0, 0, 0, 0.06);
}

.bank-info {
  display: flex;
  align-items: center;
  gap: 24rpx;
}

.bank-icon {
  width: 100rpx;
  height: 100rpx;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 20rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 48rpx;
}

.bank-text {
  display: flex;
  flex-direction: column;
  gap: 8rpx;
}

.bank-name {
  font-size: 32rpx;
  font-weight: 500;
  color: #333;
}

.bank-meta {
  font-size: 24rpx;
  color: #999;
}

.start-btn {
  padding: 16rpx 32rpx;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: #fff;
  border-radius: 32rpx;
  font-size: 26rpx;
}

.loading-tip,
.no-more {
  text-align: center;
  padding: 30rpx;
  color: #999;
  font-size: 26rpx;
}
</style>
