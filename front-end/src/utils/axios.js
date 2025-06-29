// 导入 axios 框架
import axios from "axios";
import router from '../router/index.js'
import { useUserStore } from '~/store/user'
import {ref, h} from 'vue'
import { storeToRefs } from "pinia";
import { removeToken } from "./token.js";
const baseUrl = import.meta.env.VITE_API_URL;
const loading = ref();
const service = axios.create({
    headers: {
        'Content-Type': 'application/json'
    },
    timeout: 6000,
    baseURL: baseUrl
})

// 添加请求拦截器
service.interceptors.request.use(function (config) {
    const userStore = useUserStore();
    loading.value = ElLoading.service({
        lock: true,
        text: 'Loading',
        background: 'rgba(0, 0, 0, 0.7)',
      })
    // 在发送请求之前做些什么
    // if(!userStore.getToken){
    //     router.replace({
    //         path: '/login',
    //         query:{
    //         }
    //     })
    // }

    // 添加loading
    // 添加token
    config.headers.Authorization = 'Bearer ' + userStore.token // token 从store取
    // 默认get请求
    if(config.method == null ||  config.method.toLowerCase()  !== 'post' ){
        config.method = 'get'
    }
    return config;
  }, function (error) {
    console.log(error)
    // 对请求错误做些什么
    return Promise.reject(error);
  });

// 添加响应拦截器
service.interceptors.response.use(function (response) {
    const userStore = useUserStore();
    loading.value.close();
    // 默认返回有信息 就提示
    const msg = response.data.message;
  
    if(msg != null && msg.length > 0){
        const status = response.data?.code;
        const type = status == 200 ? 'success': 'error';
        ElNotification({
            title: '消息提示',
            message: msg,
            type:type,
            duration: 1500
          })
    }
    // 2xx 范围内的状态码都会触发该函数。
    // 对响应数据做点什么
    if(response.status  == 200){
        // 可能是void返回， 默认为成功
        const code = response.data.code || 200;
        if(code == 200){

            return Promise.resolve(response.data);
        }else if(code == 401){
            userStore.logout();
            router.replace({
                path: '/login',
                
                query:{
                    redirect:  router.currentRoute.fullPath
                }
            })
        }
    }

    //如未登录, 跳转到登录页
    if(response.status == 401) {
        userStore.logout();
        router.replace({
            path: '/login',
            query:{
                redirect:  router.currentRoute.fullPath
            }
        })
    }

   

    return Promise.reject(response);
  
  }, function (error) {
    loading.value.close();
    const errorCode =  error.code;
    if(errorCode == 'ERR_BAD_RESPONSE'){
        ElNotification({
            title: '错误提示',
            message: "服务器请求失败",
            type:'error',
            duration: 1500

          })
    }else  if(error.message && error.message.startsWith('timeout of ') ){

    }

    console.log(error);    
    // 超出 2xx 范围的状态码都会触发该函数。
    // 对响应错误做点什么
    return Promise.reject(error);
  });


  export default service;