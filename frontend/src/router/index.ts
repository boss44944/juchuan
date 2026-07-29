import { createRouter, createWebHistory } from 'vue-router'
import Devices from '../views/Devices.vue'
import Messages from '../views/Messages.vue'
import Send from '../views/Send.vue'
import Config from '../views/Config.vue'
import Login from '../views/Login.vue'
import { authStatus } from '../api'

const router = createRouter({
  history:createWebHistory(),
  routes:[
    {path:'/login', component:Login, meta: { public: true }},
    {path:'/devices', component:Devices},
    {path:'/messages', component:Messages},
    {path:'/send', component:Send},
    {path:'/config', component:Config},
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
    if (ok) {
      return true
    }
  } catch {
    // Ignore and force login below.
  }

  return '/login'
})

export default router
