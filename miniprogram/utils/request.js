const BASE_URL = 'http://118.196.24.214:10002/api/v1'

const request = (options) => {
  return new Promise((resolve, reject) => {
    const app = getApp()
    const token = app.globalData.token

    wx.request({
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
          if (res.data.code === 0 || res.data.code === 200) {
            resolve(res.data.data)
          } else if (res.data.code === 401) {
            app.clearUserData()
            wx.switchTab({ url: '/pages/profile/index' })
            reject(new Error(res.data.msg || '未授权'))
          } else {
            wx.showToast({ title: res.data.msg || '请求失败', icon: 'none' })
            reject(new Error(res.data.msg))
          }
        } else {
          wx.showToast({ title: '网络错误', icon: 'none' })
          reject(new Error('网络错误'))
        }
      },
      fail: (err) => {
        wx.showToast({ title: '请求失败', icon: 'none' })
        reject(err)
      }
    })
  })
}

const get = (url, data) => request({ url, method: 'GET', data })
const post = (url, data) => request({ url, method: 'POST', data })
const put = (url, data) => request({ url, method: 'PUT', data })
const del = (url, data) => request({ url, method: 'DELETE', data })

module.exports = {
  get,
  post,
  put,
  del
}
