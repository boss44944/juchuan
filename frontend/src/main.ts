import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import './styles/theme.css'
import App from './App.vue'
import i18n from './i18n'
import router from './router'

createApp(App)
  .use(createPinia())
  .use(ElementPlus)
  .use(i18n)
  .use(router)
  .mount('#app')
