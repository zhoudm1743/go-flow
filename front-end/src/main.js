import { createApp, onMounted } from 'vue'
import ElementPlus from 'element-plus'
import { setupElementPlus } from "epic-designer/dist/ui/elementPlus";
import "epic-designer/dist/style.css";
import 'element-plus/dist/index.css'
import App from './App.vue'
import router from './router/index.js'
import './style.css'
import 'element-plus/theme-chalk/dark/css-vars.css'
import zhCn from 'element-plus/es/locale/lang/zh-cn'
import direct from '@/utils/direct.js'

import './router/routerGuards.js'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'


import { createPinia } from 'pinia'
// 持久化pinia 数据
import piniaPersist from 'pinia-plugin-persist';
setupElementPlus();
const pinia = createPinia();
pinia.use(piniaPersist)

const app = createApp(App)
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
    app.component(key, component)
}
app.use(ElementPlus,{
    locale: zhCn,
  }).use(pinia).use(router)

app.directive('direct-auto-height', direct)

app.mount('#app')



