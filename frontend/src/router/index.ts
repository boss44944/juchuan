import { createRouter, createWebHistory } from 'vue-router'
import Devices from '../views/Devices.vue'
import Messages from '../views/Messages.vue'
import Send from '../views/Send.vue'
import Config from '../views/Config.vue'
import Login from '../views/Login.vue'

export default createRouter({
  history:createWebHistory(),
  routes:[
    {path:'/login', component:Login},
    {path:'/devices', component:Devices},
    {path:'/messages', component:Messages},
    {path:'/send', component:Send},
    {path:'/config', component:Config},
    {path:'/', redirect:'/devices'}
  ]
})
