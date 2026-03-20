<template>
  <div class="layout-container">
    <aside class="sidebar">
      <div class="logo">
        <span class="logo-icon">📚</span>
        <span class="logo-text">刷题后台</span>
      </div>
      <el-menu
        :default-active="$route.path"
        router
        background-color="transparent"
        text-color="#a8b2c4"
        active-text-color="#ffffff"
      >
        <el-menu-item v-for="item in menuList" :key="item.path" :index="item.path">
          <el-icon><component :is="item.icon" /></el-icon>
          <span>{{ item.title }}</span>
        </el-menu-item>
      </el-menu>
    </aside>
    
    <div class="main-wrapper">
      <header class="header">
        <div class="header-left">
          <h2 class="page-title">{{ $route.meta.title }}</h2>
        </div>
        <div class="header-right">
          <el-dropdown @command="handleCommand">
            <span class="user-info">
              <el-avatar :size="32">{{ userStore.nickname?.[0] || 'A' }}</el-avatar>
              <span class="username">{{ userStore.nickname || userStore.username }}</span>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="logout">退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </header>
      
      <main class="main-content">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </main>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useUserStore } from '@/stores/user'
import { useRoute } from 'vue-router'

const userStore = useUserStore()
const $route = useRoute()

const menuList = [
  { path: '/dashboard', title: '仪表盘', icon: 'Odometer' },
  { path: '/question-bank', title: '题库管理', icon: 'FolderOpened' },
  { path: '/question', title: '题目管理', icon: 'Document' },
  { path: '/user', title: '用户管理', icon: 'User' }
]

const handleCommand = (command) => {
  if (command === 'logout') {
    userStore.logout()
  }
}
</script>

<style lang="scss" scoped>
$primary-color: #f59e0b;
$bg-dark: #1a1d23;
$bg-sidebar: #13151a;
$text-secondary: #a8b2c4;

.layout-container {
  display: flex;
  height: 100vh;
  background: #f5f6f8;
}

.sidebar {
  width: 220px;
  background: $bg-sidebar;
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
}

.logo {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 20px 24px;
  border-bottom: 1px solid rgba(255,255,255,0.05);
  
  .logo-icon {
    font-size: 24px;
  }
  
  .logo-text {
    font-size: 18px;
    font-weight: 600;
    color: #fff;
    letter-spacing: 0.5px;
  }
}

:deep(.el-menu) {
  border: none;
  padding: 12px 0;
  
  .el-menu-item {
    height: 48px;
    line-height: 48px;
    margin: 4px 12px;
    border-radius: 8px;
    
    &:hover {
      background: rgba(255,255,255,0.06);
    }
    
    &.is-active {
      background: $primary-color;
      
      .el-icon {
        color: #fff;
      }
    }
  }
}

.main-wrapper {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.header {
  height: 64px;
  background: #fff;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  box-shadow: 0 1px 2px rgba(0,0,0,0.05);
  
  .page-title {
    font-size: 18px;
    font-weight: 600;
    color: #1f2937;
    margin: 0;
  }
  
  .user-info {
    display: flex;
    align-items: center;
    gap: 10px;
    cursor: pointer;
    
    .username {
      color: #374151;
      font-weight: 500;
    }
  }
}

.main-content {
  flex: 1;
  padding: 24px;
  overflow-y: auto;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
