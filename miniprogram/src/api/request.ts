const BASE_URL = 'http://localhost:10002/api/v1'

interface RequestOptions {
  url: string
  method?: 'GET' | 'POST' | 'PUT' | 'DELETE'
  data?: any
  header?: Record<string, string>
}

interface Response<T = any> {
  code: number
  msg: string
  data: T
}

const request = <T = any>(options: RequestOptions): Promise<T> => {
  return new Promise((resolve, reject) => {
    const token = uni.getStorageSync('token')
    
    uni.request({
      url: BASE_URL + options.url,
      method: options.method || 'GET',
      data: options.data,
      header: {
        'Content-Type': 'application/json',
        'Authorization': token ? `Bearer ${token}` : '',
        ...options.header
      },
      success: (res) => {
        if (res.statusCode === 200) {
          const data = res.data as Response<T>
          if (data.code === 0 || data.code === 200) {
            resolve(data.data)
          } else if (data.code === 401) {
            uni.removeStorageSync('token')
            uni.removeStorageSync('userInfo')
            uni.reLaunch({ url: '/pages/login/index' })
            reject(new Error(data.msg || '未授权'))
          } else {
            uni.showToast({ title: data.msg || '请求失败', icon: 'none' })
            reject(new Error(data.msg))
          }
        } else {
          uni.showToast({ title: '网络错误', icon: 'none' })
          reject(new Error('网络错误'))
        }
      },
      fail: (err) => {
        uni.showToast({ title: '请求失败', icon: 'none' })
        reject(err)
      }
    })
  })
}

export const get = <T = any>(url: string, data?: any) => request<T>({ url, method: 'GET', data })
export const post = <T = any>(url: string, data?: any) => request<T>({ url, method: 'POST', data })
export const put = <T = any>(url: string, data?: any) => request<T>({ url, method: 'PUT', data })
export const del = <T = any>(url: string, data?: any) => request<T>({ url, method: 'DELETE', data })
