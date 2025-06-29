import service from "~/utils/axios";

export function doLogin(data){
    return service({
        url: '/auth/login',
        method: 'POST',
        data: data,
    })
}

export function getUserInfo(){
    return service({
        url: '/auth/user-info',
        method: 'GET',
    })
}

export function getCaptcha(){
    return service({
        url: '/auth/captcha',
        method: 'GET',
    })
}