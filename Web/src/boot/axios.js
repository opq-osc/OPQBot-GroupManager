import Vue from 'vue'
import axios from 'axios'
import { stringify } from 'qs'
axios.defaults.withCredentials = true
axios.defaults.baseURL = window.location.protocol + '//' + window.location.hostname + ':8888/'

axios.interceptors.request.use(
  (config) => {
    // 兼容 post 跨域问题
    if (config.method === 'post') {
      // 修改 Content-Type
      config.headers['Content-Type'] =
          'application/x-www-form-urlencoded'

      // 将对象参数转换为序列化的 URL 形式（key=val&key=val）
      config.data = stringify(config.data)
    }
    return config
  },
  (error) => {
    console.log(error)
    return Promise.reject(error)
  }
)
axios.interceptors.response.use(res => {
  if (res.data.code === 10010 || res.data.code === 10011) {
    this.$router.replace('/login')
  }
  return res
}, error => {
  return Promise.reject('出错啦', error)
})

function getCookie (cname) {
  const name = cname + '='
  const ca = document.cookie.split(';')
  for (let i = 0; i < ca.length; i++) {
    const c = ca[i].trim()
    if (c.indexOf(name) === 0) return c.substring(name.length, c.length)
  }
  return ''
}
function formatDate (value) {
  const date = new Date(value)
  const y = date.getFullYear()
  let MM = date.getMonth() + 1
  MM = MM < 10 ? ('0' + MM) : MM
  let d = date.getDate()
  d = d < 10 ? ('0' + d) : d
  let h = date.getHours()
  h = h < 10 ? ('0' + h) : h
  let m = date.getMinutes()
  m = m < 10 ? ('0' + m) : m
  let s = date.getSeconds()
  s = s < 10 ? ('0' + s) : s
  return y + '-' + MM + '-' + d + ' ' + h + ':' + m + ':' + s
}
Vue.prototype.$cookie = getCookie
Vue.prototype.$axios = axios
Vue.prototype.$timeformat = formatDate
import md5 from 'js-md5'
Vue.prototype.$md5 = md5
