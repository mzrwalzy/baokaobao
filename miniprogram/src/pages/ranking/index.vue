<template>
  <view class="ranking-page">
    <view class="header">
      <text class="header-title">排行榜</text>
    </view>

    <view class="tab-bar">
      <view 
        class="tab-item" 
        :class="{ active: tab === 'monthly' }"
        @click="tab = 'monthly'"
      >月榜</view>
      <view 
        class="tab-item"
        :class="{ active: tab === 'total' }"
        @click="tab = 'total'"
      >总榜</view>
    </view>

    <scroll-view scroll-y class="content" @scrolltolower="loadMore">
      <view class="top-three">
        <view class="rank-item second" v-if="list[1]">
          <view class="avatar">🥈</view>
          <view class="nickname">{{ list[1].nickname }}</view>
          <view class="score">{{ list[1].score }}分</view>
        </view>
        <view class="rank-item first" v-if="list[0]">
          <view class="avatar">🥇</view>
          <view class="nickname">{{ list[0].nickname }}</view>
          <view class="score">{{ list[0].score }}分</view>
        </view>
        <view class="rank-item third" v-if="list[2]">
          <view class="avatar">🥉</view>
          <view class="nickname">{{ list[2].nickname }}</view>
          <view class="score">{{ list[2].score }}分</view>
        </view>
      </view>

      <view class="rank-list">
        <view 
          v-for="(item, index) in list.slice(3)" 
          :key="index"
          class="rank-item-normal"
        >
          <view class="rank-num">{{ index + 4 }}</view>
          <view class="user-info">
            <view class="user-avatar">{{ item.nickname?.charAt(0) || '?' }}</view>
            <view class="user-name">{{ item.nickname }}</view>
          </view>
          <view class="user-score">{{ item.score }}分</view>
        </view>
      </view>

      <view class="loading-tip" v-if="loading">
        <text>加载中...</text>
      </view>
    </scroll-view>
  </view>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { getRanking, type RankingItem } from '@/api'
import { useUserStore } from '@/stores/user'

const userStore = useUserStore()
const tab = ref<'monthly' | 'total'>('total')
const list = ref<RankingItem[]>([])
const loading = ref(false)

watch(tab, () => {
  list.value = []
  loadRanking()
})

const loadRanking = async () => {
  if (loading.value) return
  
  loading.value = true
  try {
    const data = await getRanking(tab.value)
    list.value = data || []
  } catch (e: any) {
    uni.showToast({ title: e.message || '加载失败', icon: 'none' })
  } finally {
    loading.value = false
  }
}

const loadMore = () => {
  // 分页加载逻辑
}

loadRanking()
</script>

<style scoped>
.ranking-page {
  min-height: 100vh;
  background: #F5F5F5;
}

.header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 30rpx;
  padding-top: calc(30rpx + var(--status-bar-height));
  text-align: center;
}

.header-title {
  font-size: 36rpx;
  font-weight: bold;
  color: #fff;
}

.tab-bar {
  display: flex;
  background: #fff;
  padding: 0 100rpx;
}

.tab-item {
  flex: 1;
  text-align: center;
  padding: 30rpx 0;
  font-size: 30rpx;
  color: #666;
  position: relative;
}

.tab-item.active {
  color: #409EFF;
  font-weight: bold;
}

.tab-item.active::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 50%;
  transform: translateX(-50%);
  width: 60rpx;
  height: 6rpx;
  background: #409EFF;
  border-radius: 3rpx;
}

.content {
  height: calc(100vh - 250rpx);
}

.top-three {
  display: flex;
  justify-content: center;
  align-items: flex-end;
  padding: 60rpx 30rpx 30rpx;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.rank-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 0 20rpx;
}

.rank-item.first {
  order: 2;
}

.rank-item.second {
  order: 1;
  padding-bottom: 30rpx;
}

.rank-item.third {
  order: 3;
  padding-bottom: 10rpx;
}

.avatar {
  font-size: 60rpx;
  margin-bottom: 10rpx;
}

.first .avatar {
  font-size: 80rpx;
}

.nickname {
  font-size: 28rpx;
  color: #fff;
  margin-bottom: 8rpx;
  max-width: 150rpx;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.score {
  font-size: 24rpx;
  color: rgba(255, 255, 255, 0.8);
  background: rgba(255, 255, 255, 0.2);
  padding: 6rpx 20rpx;
  border-radius: 20rpx;
}

.rank-list {
  padding: 20rpx 30rpx;
}

.rank-item-normal {
  display: flex;
  align-items: center;
  background: #fff;
  border-radius: 16rpx;
  padding: 24rpx 30rpx;
  margin-bottom: 20rpx;
}

.rank-num {
  width: 60rpx;
  font-size: 32rpx;
  font-weight: bold;
  color: #999;
}

.user-info {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 20rpx;
}

.user-avatar {
  width: 72rpx;
  height: 72rpx;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-size: 32rpx;
  font-weight: bold;
}

.user-name {
  font-size: 30rpx;
  color: #333;
}

.user-score {
  font-size: 28rpx;
  color: #409EFF;
  font-weight: 500;
}

.loading-tip {
  text-align: center;
  padding: 30rpx;
  color: #999;
}
</style>
