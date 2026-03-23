<template>
  <view class="login-page">
    <view class="logo-section">
      <view class="logo">
        <text class="logo-icon">📚</text>
      </view>
      <text class="app-name">刷题宝</text>
      <text class="app-slogan">轻松刷题，高效学习</text>
    </view>

    <view class="login-section">
      <button class="btn-login" @click="handleWxLogin" :loading="loading">
        <text class="wx-icon">🍎</text>
        <text>微信一键登录</text>
      </button>
      
      <view class="agreement">
        登录即表示同意
        <text class="link" @click="openAgreement('user')">《用户协议》</text>
        和
        <text class="link" @click="openAgreement('privacy')">《隐私政策》</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { loginByWechat } from '@/api'
import { useUserStore } from '@/stores/user'

const userStore = useUserStore()
const loading = ref(false)

const handleWxLogin = async () => {
  if (loading.value) return
  
  loading.value = true
  
  try {
    const code = await new Promise<string>((resolve, reject) => {
      uni.login({
        provider: 'weixin',
        success: (res) => resolve(res.code),
        fail: (err) => reject(err)
      })
    })

    const res = await loginByWechat(code)
    userStore.login(res.token, res.user)
    
    uni.switchTab({ url: '/pages/index/index' })
  } catch (e: any) {
    uni.showToast({ title: e.message || '登录失败', icon: 'none' })
  } finally {
    loading.value = false
  }
}

const openAgreement = (type: 'user' | 'privacy') => {
  uni.showToast({ title: `查看${type === 'user' ? '用户协议' : '隐私政策'}`, icon: 'none' })
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  padding: 120rpx 60rpx 100rpx;
}

.logo-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-top: 100rpx;
}

.logo {
  width: 160rpx;
  height: 160rpx;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 40rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 40rpx;
}

.logo-icon {
  font-size: 80rpx;
}

.app-name {
  font-size: 56rpx;
  font-weight: bold;
  color: #ffffff;
  margin-bottom: 20rpx;
}

.app-slogan {
  font-size: 28rpx;
  color: rgba(255, 255, 255, 0.8);
}

.login-section {
  width: 100%;
}

.btn-login {
  width: 100%;
  height: 96rpx;
  background: #ffffff;
  border-radius: 48rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 32rpx;
  color: #333;
  font-weight: 500;
  border: none;
}

.btn-login::after {
  border: none;
}

.wx-icon {
  font-size: 40rpx;
  margin-right: 16rpx;
}

.agreement {
  text-align: center;
  color: rgba(255, 255, 255, 0.7);
  font-size: 24rpx;
  margin-top: 40rpx;
}

.link {
  color: #ffffff;
  text-decoration: underline;
}
</style>
