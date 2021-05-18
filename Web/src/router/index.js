import Vue from 'vue'
import VueRouter from 'vue-router'

import routes from './routes'

Vue.use(VueRouter)

/*
 * If not building with SSR mode, you can
 * directly export the Router instantiation;
 *
 * The function below can be async too; either use
 * async/await or return a Promise which resolves
 * with the Router instance.
 */

export default function (/* { store, ssrContext } */) {
  const Router = new VueRouter({
    scrollBehavior: () => ({ x: 0, y: 0 }),
    routes,

    // Leave these as they are and change in quasar.conf.js instead!
    // quasar.conf.js -> build -> vueRouterMode
    // quasar.conf.js -> build -> publicPath
    mode: process.env.VUE_ROUTER_MODE,
    base: process.env.VUE_ROUTER_BASE
  })
  Router.beforeEach((to, from, next) => {
    // console.log(to)
    if (to.path === '/error') {
      next()
    }
    if (to.meta.requireAuth) {
      if (Router.app.$store.state.User.auth) {
        next()
      } else {
        Router.app.$axios.get('api/status').then(function (response) {
          if (response.data.code === 1) {
            Router.app.$store.commit('User/pushFunc', { auth: response.data.code === 1, username: response.data.data })
            next()
          } else {
            next({
              path: '/login',
              query: { redirect: to.path, info: '未登录或登录时间过长令牌失效了，请重新登录！' }
            })
          }
        }).catch(() => {
          next({
            path: '/login',
            query: { redirect: to.path, info: '后端出现错误，请尝试重新登录' }
          })
        })
      }
    } else {
      next()
    }

    // Router.app.$axios.get('api/status').then(function (response) {
    //   Router.app.$store.commit('User/pushFunc', { auth: response.data.code === 1, username: response.data.data })
    //   if (to.meta.requireAuth) {
    //     if (response.data.code === 1) {
    //       next()
    //     } else {
    //       next({
    //         path: '/login',
    //         query: { redirect: to.fullPath, info: '未登录或登录时间过长令牌失效了，请重新登录！' }
    //       })
    //     }
    //   } else {
    //     next()
    //   }
    // }).catch(function (error) {
    //   console.log(error.message)
    //   let title = '错误'
    //   if (error.message === 'Network Error') {
    //     error.message = '与后端连接出现问题，网站暂时无法使用！'
    //     title = '网络错误'
    //   }
    //   next({
    //     path: '/error',
    //     query: { info: error.message, title: title, back: to.fullPath }
    //   })
    // })
  })
  return Router
}
