import { createMemoryHistory, createRouter, createWebHashHistory } from 'vue-router'

import Login from '~/views/login.vue'
import Layout from '~/components/layouts/index.vue';
import Home from '~/views/home/index.vue'
import FlowConfig from '~/views/flowconfig/index.vue'
import Profile from '@/views/sys/user/info/profile.vue';
export const baseRoutes = [
  { 
    path: '', 
    name:"Layout" , 
    redirect: '/home',
    component: Layout,
    children:[
      // { path: '/index', component: Home, meta:{title: "首页"} },
      // { path: '/home', component: Home, meta:{title: "首页"} },
      
      // { path: '/flowconfig', component: FlowConfig , meta:{title: "流程配置"} },

      {
        path: "/:pathMatch(.*)*",
        component: () => import('@/views/error/404.vue'),
        meta:{title: "404"} 
      },
      { path: '/profile', component: Profile, meta:{title: "个人中心"} },
     
    ] 
  },

  { path: '/login', component: Login, name: "Login"  },

]

const router = createRouter({
  history: createWebHashHistory(),
  routes:baseRoutes,
})



export default router;