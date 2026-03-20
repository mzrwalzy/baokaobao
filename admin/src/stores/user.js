import { defineStore } from 'pinia'
import { adminLogin, adminLogout } from '@/api'
import router from '@/router'

export const useUserStore = defineStore('user', {
  state: () => ({
    token: localStorage.getItem('admin_token') || '',
    username: localStorage.getItem('admin_username') || '',
    nickname: localStorage.getItem('admin_nickname') || '',
    role: localStorage.getItem('admin_role') || ''
  }),

  actions: {
    async login(username, password) {
      const res = await adminLogin({ username, password })
      this.token = res.token
      this.username = res.user.username
      this.nickname = res.user.nickname
      this.role = res.user.role
      
      localStorage.setItem('admin_token', this.token)
      localStorage.setItem('admin_username', this.username)
      localStorage.setItem('admin_nickname', this.nickname)
      localStorage.setItem('admin_role', this.role)
      
      return res
    },

    async logout() {
      try {
        await adminLogout()
      } catch (e) {}
      this.token = ''
      this.username = ''
      this.nickname = ''
      this.role = ''
      localStorage.removeItem('admin_token')
      localStorage.removeItem('admin_username')
      localStorage.removeItem('admin_nickname')
      localStorage.removeItem('admin_role')
      router.push('/login')
    }
  }
})
