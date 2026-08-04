import { createRouter, createWebHistory, type RouteLocationNormalized } from 'vue-router'
import Login from '../views/Login.vue'
import ServerShell from '../apps/server/ServerShell.vue'
import ClientShell from '../apps/client/ClientShell.vue'
import Devices from '../views/Devices.vue'
import Messages from '../views/Messages.vue'
import Send from '../views/Send.vue'
import Config from '../views/Config.vue'
import ClientInbox from '../views/client/Inbox.vue'
import ClientSend from '../views/client/Send.vue'
import { authStatus } from '../api'
import { currentRole, defaultRouteFor, roleFromPath, type AccessRole } from '../utils/role'

function routeRole(to: RouteLocationNormalized): AccessRole {
  return (to.meta.role as AccessRole | undefined) || roleFromPath(to.path) || currentRole(to.path)
}

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/server/login', component: Login, meta: { public: true, role: 'server' } },
    { path: '/client/login', component: Login, meta: { public: true, role: 'client' } },
    {
      path: '/server',
      component: ServerShell,
      meta: { role: 'server' },
      children: [
        { path: '', redirect: '/server/devices' },
        { path: 'devices', component: Devices },
        { path: 'messages', component: Messages },
        { path: 'send', component: Send },
        { path: 'config', component: Config },
      ],
    },
    {
      path: '/client',
      component: ClientShell,
      meta: { role: 'client' },
      children: [
        { path: '', redirect: '/client/inbox' },
        { path: 'inbox', component: ClientInbox },
        { path: 'send', component: ClientSend },
      ],
    },
    { path: '/login', redirect: () => `/${currentRole()}/login` },
    { path: '/devices', redirect: '/server/devices' },
    { path: '/messages', redirect: '/server/messages' },
    { path: '/send', redirect: () => currentRole() === 'server' ? '/server/send' : '/client/send' },
    { path: '/config', redirect: '/server/config' },
    { path: '/', redirect: () => defaultRouteFor(currentRole()) },
    { path: '/:pathMatch(.*)*', redirect: () => defaultRouteFor(currentRole()) },
  ],
})

router.beforeEach(async (to) => {
  const role = routeRole(to)
  if (to.meta.public) return true

  try {
    const response = await authStatus()
    const authenticated = !response.data.requires_password || response.data.authenticated
    if (!authenticated) {
      return { path: `/${role}/login`, query: { redirect: to.fullPath } }
    }
  } catch {
    return { path: `/${role}/login`, query: { redirect: to.fullPath } }
  }

  // A browser client needs a stable device identity so its inbox can be scoped.
  if (role === 'client' && !(localStorage.getItem('device_id') || '').trim()) {
    return { path: '/client/login', query: { redirect: to.fullPath } }
  }

  return true
})

export default router
