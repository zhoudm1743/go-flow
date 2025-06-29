// 用户权限 store，大致看了一下若依框架的源码，自己实现一套逻辑

import { defineStore } from 'pinia'
import { menuRouter } from '~/api/sys/menu'
import Layout from '~/components/layouts/index.vue';
import router from '../router';
import FlowConfig from '~/views/flowconfig/index.vue'
const modules = import.meta.glob("../views/**/**.vue")
import { baseRoutes } from '@/router/index.js';
export const usePermissionStore = defineStore('permission', {
  state: () =>({
    // 菜单列表
    routers: [],
    menuList:[],
    menuTreeList: [],
    baseRouter:[]

  }),
  getters:{
    menuRouter : (state) => state.routers,
    menuTree : (state) => state.menuList,
    layoutRouter : (state) => state.baseRouter
  },
  // 也可以这样定义
  // state: () => ({ count: 0 })
  actions: {
    
    // 获取当前用户的菜单
    loadMenuRouter() {
      var that = this;
    
      // const uuid = userInfo.uuid
       return new Promise((resolve, reject) => {
        menuRouter().then(res => {
          that.menuTreeList = JSON.parse(JSON.stringify(res.rows));
          // setToken(res.data.access_token)
          that.menuList =   JSON.parse(JSON.stringify(res.rows));
          const data = JSON.parse(JSON.stringify(res.rows));
          // that.baseRouter.push(...firstMenu(data))
          // that.baseRouter = firstMenu(data)
          if(that.menuTreeList && that.menuTreeList.length > 0){
            that.baseRouter = baseRoutes[0].children.push(...firstMenu(data))
            const routerList  =  convertToRouter(data);
             console.log("routerList")
            console.log(routerList)
            // that.routers.push(...baseRouter)
            that.routers.push(...routerList, ...baseRoutes);
          }else{
            that.baseRouter = baseRoutes[0];
            that.routers.push(...baseRoutes);
          }
        
          resolve(that.routers)
        }).catch(error => {
          reject(error)
        })
      })
    },
    clernRoutes() {
      this.routers = [];
      this.menuList = [];
      this.menuTreeList = [];
      this.baseRouter = {};
    },
    
  },
  persist: {
    enabled: false
  }
})


const firstMenu = (data) => {
  const children = [];
  // let baseRouter = { 
  //   name:"Layout" , 
  //   component: Layout,
  //   children:[
  //   ] 
  // };
  // const children = [];
  data.forEach(item => {
    if(item.type == '2'){
      // const view =`../views/${item.component}.vue`
      item.component = loadView(item.component) //  modules[`../views/${item.component}.vue`] // () => import(view) 
      item.meta = {title: item.name}
      children.push(item);
    }
  })
  children.push({ path: '/flowconfig/:id', component: FlowConfig , name: "Flowconfig", meta:{title: "流程配置"} })
  return children;
  // return baseRouter
}

const convertToRouter = (menuList, istop = true) => {
  const menuListResult = [];

  menuList.forEach(item => {
    if(item.type == '1'){  //目录
      item.component = Layout
      menuListResult.push(item);
    }else if(item.type == '2' && !istop){ //菜单
        item.component = loadView(item.component) // () => import(view) 
        menuListResult.push(item);
     
    }
    item.meta = {title: item.name}
    if(item.children){
     const children  = convertToRouter(item.children, false);
     item.children = children;
    }
     
  })
  return menuListResult;
}


export const loadView = (view) => {
  let res;
  for (const path in modules) {
    const dir = path.split('views/')[1].split('.vue')[0];
    if (dir === view) {
      res = () => modules[path]();
    }
  }
  return res;
}