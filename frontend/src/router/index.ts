import { createRouter, createWebHistory } from 'vue-router'
import Devices from '../views/Devices.vue'
import Messages from '../views/Messages.vue'
import Send from '../views/Send.vue'
import Config from '../views/Config.vue'

export default createRouter({
  history:createWebHistory(),
  routes:[
    {path:'/devices', component:Devices},
    {path:'/messages', component:Messages},
    {path:'/send', component:Send},
    {path:'/config', component:Config}
  ]
})
