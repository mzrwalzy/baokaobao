<template>
  <view class="profile-page">
    <view class="header-bg">
      <view class="user-card">
        <view class="avatar">
          {{ userInfo?.nickname?.charAt(0) || '?' }}
        </view>
        <view class="nickname">{{ userInfo?.nickname || '用户' }}</view>
        <view class="vip-info" v-if="userInfo?.phone">
          会员到期: 2025-12-31
        </view>
      </view>
    </view>

    <view class="stats-section">
      <view class="stats-grid">
        <view class="stats-item">
          <view class="stats-value">{{ stats?.total_score || 0 }}</view>
          <view class="stats-label">积分</view>
        </view>
        <view class="stats-item">
          <view class="stats-value">{{ stats?.study_days || 0 }}</view>
          <view class="stats-label">学习天数</view>
        </view>
        <view class="stats-item">
          <view class="stats-value">{{ stats?.wrong_count || 0 }}</view>
          <view class="stats-label">错题数</view>
        </view>
        <view class="stats-item">
          <view class="stats-value">{{ stats?.correct_rate || 0 }}%</view>
          <view class="stats-label">正确率</view>
        </view>
      </view>
    </view>

    <view class="menu-section">
      <view class="menu-item" @click="goPage('/pages/profile/wrong-questions')">
        <text class="menu-icon">❌</text>
        <text class="menu-text">错题本</text>
        <text class="menu-badge">{{ stats?.wrong_count || 0 }}题</text>
        <text class="menu-arrow">›</text>
      </view>
      <view class="menu-item" @click="goPage('/pages/profile/study-record')">
        <text class="menu-icon">📚</text>
        <text class="menu-text">学习记录</text>
        <text class="menu-arrow">›</text>
      </view>
      <view class="menu-item" @click="goPage('/pages/profile/my-banks')">
        <text class="menu-icon">⭐</text>
        <text class="menu-text">收藏题库</text>
        <text class="menu-arrow">›</text>
      </view>
    </view>

    <view class="menu-section">
      <view class="menu-item" @click="goPage('/pages/profile/settings')">
        <text class="menu-icon">⚙️</text>
        <text class="menu-text">设置</text>
        <text class="menu-arrow">›</text>
      </view>
      <view class="menu-item" @click="goPage('/pages/profile/help')">
        <text class="menu-icon">📖</text>
        <text class="menu-text">使用帮助</text>
        <text class="menu-arrow">›</text>
      </view>
      <view class="menu-item" @click="goPage('/pages/profile/feedback')">
        <text class="menu-icon">📝</text>
        <text class="menu-text">意见反馈</text>
        <text class="menu-arrow">›</text>
      </view>
      <view class="menu-item" @click="goPage('/pages/profile/about')">
        <text class="menu-icon">ℹ️</text>
        <text class="menu-text">关于我们</text>
        <text class="menu-arrow">›</text>
      </view>
    </view>

    <view class="logout-section">
      <button class="btn-logout" @click="handleLogout">退出登录</button>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getProfile, getStats, logout } from '@/api'
import { useUserStore } from '@/stores/user'

const userStore = useUserStore()
const userInfo = ref(userStore.userInfo)
const stats = ref<any>(null)

onMounted(async () => {
  if (!userStore.isLoggedIn) {
    uni.reLaunch({ url: '/pages/login/index' })
    return
  }
  
  try {
    const [profile, userStats] = await Promise.all([
      getProfile(),
      getStats()
    ])
    userInfo.value = profile
    stats.value = userStats
  } catch (e: any) {
    console.error(e)
  }
})

const goPage = (url: string) => {
  uni.showToast({ title: '功能开发中', icon: 'none' })
}

const handleLogout = () => {
  uni.showModal({
    title: '提示',
    content: '确定要退出登录吗？',
    success: async (res) => {
      if (res.confirm) {
        try {
          await logout()
        } catch (e) {}
        userStore.logout()
        uni.reLaunch({ url: '/pages/login/index' })
      }
    }
  })
}
</script>

<style scoped>
.profile-page {
  min-height: 100vh;
  background: #F5F5F5;
  padding-bottom: 200rpx;
}

.header-bg {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 60rpx 30rpx 120rpx;
}

.user-card {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.avatar {
  width: 140rpx;
  height: 140rpx;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 60rpx;
  color: #fff;
  font-weight: bold;
  margin-bottom: 20rpx;
  border: 6rpx solid rgba(255, 255, 255, 0.3);
}

.nickname {
  font-size: 36rpx;
  font-weight: bold;
  color: #fff;
  margin-bottom: 10rpx;
}

.vip-info {
  font-size: 24rpx;
  color: rgba(255, 255, 255, 0.8);
  background: rgba(255, 255, 255, 0.15);
  padding: 8rpx 24rpx;
  border-radius: 20rpx;
}

.stats-section {
  margin: -80rpx 30rpx 30rpx;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 1rpx;
  background: #fff;
  border-radius: 20rpx;
  overflow: hidden;
  box-shadow: 0 4rpx 20rpx rgba(0, 0, 0, 0.08);
}

.stats-item {
  background: #fff;
  padding: 30rpx 10rpx;
  text-align: center;
}

.stats-value {
  font-size: 36rpx;
  font-weight: bold;
  color: #333;
  margin-bottom: 8rpx;
}

.stats-label {
  font-size: 22rpx;
  color: #999;
}

.menu-section {
  margin: 30rpx;
  background: #fff;
  border-radius: 20rpx;
  overflow: hidden;
}

.menu-item {
  display: flex;
  align-items: center;
  padding: 32rpx 30rpx;
  border-bottom: 1rpx solid #f0f0f0;
}

.menu-item:last-child {
  border-bottom: none;
}

.menu-icon {
  font-size: 36rpx;
  margin-right: 24rpx;
}

.menu-text {
  flex: 1;
  font-size: 30rpx;
  color: #333;
}

.menu-badge {
  font-size: 24rpx;
  color: #999;
  margin-right: 16rpx;
}

.menu-arrow {
  font-size: 36rpx;
  color: #ccc;
}

.logout-section {
  margin: 60rpx 30rpx 30rpx;
}

.btn-logout {
  width: 100%;
  height: 88rpx;
  background: #fff;
  color: #F56C6C;
  border: none;
  border-radius: 44rpx;
  font-size: 30rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.btn-logout::after {
  border: none;
}
</style>
