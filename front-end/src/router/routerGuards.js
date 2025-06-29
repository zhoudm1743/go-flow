import router from "./index.js"
import { useUserStore } from "~/store/user";
import { usePermissionStore } from "~/store/permission.js";

router.beforeEach((to, from, next) => {
    const permissionStore = usePermissionStore();
    const userStore = useUserStore();
    if (userStore.token) {
        if(to.name == 'Login'){
          next({ path: '/' })
        }else if(!userStore.hasInfo ) {
            userStore.loadUserInfo().then(() => {
                // 获取菜单
              permissionStore.loadMenuRouter().then((result) => {
                result.forEach((item) => {
                    router.addRoute(item)
                })
                
                // baseRoutes.forEach((item) => {
                //   router.addRoute(item)
                // })
                // router.addRoute(permissionStore.baseRouter)
                next({ ...to, replace: true })
                return false;
              });
            })
          
       }else{
         next()
       }


    }else if(to.name == 'Login'){
        next()
    }


    else {
      next({name: 'Login'})
    }
  })
