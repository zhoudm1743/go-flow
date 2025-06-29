//  用户信息
// stores/counter.js
import { defineStore } from 'pinia'
import { doLogin, getUserInfo } from '~/api/auth/login'
import { getToken, setToken, removeToken } from '@/utils/token'
import router from '~/router/index'
  //login登录页放入token， 然后接口请求时取token
export const useUserStore = defineStore('user', {
  state: () =>({
    access_token: getToken(),
    roleList:[],  // 角色
    deptList:[], // 部门id
    posList:[], // 职位id
    btnList:[], // 按钮权限
    userPositions:[],  // 职位详细
    baseInfo:{
      nick:'',
    }
  }),
  getters:{
    token: (state) => state.access_token,
    hasInfo: (state) => {
      if(state.baseInfo && state.baseInfo.nick?.length > 0){
        return true
      }
      return false;
    },
  },
  // 也可以这样定义
  // state: () => ({ count: 0 })
  actions: {
    loginByAccount(userInfo) {
      var that = this;
      userInfo.username = userInfo.username.trim()
      // const uuid = userInfo.uuid
       return new Promise((resolve, reject) => {
        doLogin(userInfo).then(res => {
          console.log(res)
          // setToken(res.data.access_token)
          setToken(res.data.token)
          this.access_token = res.data.token
          resolve()
        }).catch(error => {
          reject(error)
        })
      })
    },

    loadUserInfo() {
      var _that = this;
      return new Promise((resolve, reject) => {
          getUserInfo().then(res => {
            _that.baseInfo.nick = res.data.nickname;
            _that.baseInfo.username = res.data.username;
            _that.baseInfo.email = res.data.email;
            _that.roleList = res.data.roles;
            _that.deptList = res.data.deptList;
            _that.posList = res.data.posList;
            _that.btnList = res.data.btnList;
            _that.userPositions = res.data.userPositions;

          resolve(res)
        }).catch(error => {
          reject(error)
        })
      })
    },
    logout() {
      removeToken();
      this.access_token = '';
      this.roleList = [];
      this.deptList = [];
      this.posList = [];
      this.btnList = [];
      this.userPositions = [];
      this.baseInfo = {};
      router.push({name: 'Login'})
    },
    login(data) {
      doLogin(data).then(res=> {
          // 获取token, 保存
          setToken(res.data.token)
          this.access_token = res.data.token
      });
        //  登录逻辑处理
    },
  },
  persist: {
    enabled: false
  }
})