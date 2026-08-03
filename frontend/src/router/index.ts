import { createRouter, createWebHistory } from 'vue-router'
import Devices from '../views/Devices.vue'
import Messages from '../views/Messages.vue'
import Send from '../views/Send.vue'
import Config from '../views/Config.vue'
import Login from '../views/Login.vue'
import { authStatus } from '../api'
import { isServerAccess } from '../utils/role'

const router = createRouter({
  history:createWebHistory(),
  routes:[
    {path:'/login', component:Login, meta: { public: true }},
    {path:'/devices', component:Devices, meta: { serverOnly: true }},
    {path:'/messages', component:Messages, meta: { serverOnly: true }},
    {path:'/send', component:Send},
    {path:'/config', component:Config, meta: { serverOnly: true }},
    {path:'/', redirect:'/devices'}
  ]
})

router.beforeEach(async (to) => {
  if (to.meta.public) {
    return true
  }

  try {
    const res = await authStatus()
    const ok = !res.data.requires_password || res.data.authenticated
    if (!ok) {
      return '/login'
    }
  } catch {
    // Ignore and force login below.
    return '/login'
  }

  // 客户端（手机）只能访问「发送」页；服务端管理页一律重定向
  if (to.meta.serverOnly && !isServerAccess()) {
    return '/send'
  }

  return true
})

export default router
