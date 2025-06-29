import service from "~/utils/axios";


/**
 * 权限标识
 * @param {menuId} params 
 * @returns 
 */

export function listApi(params){
    return service({
        url: '/sys/resource/list',
        params

    })
}

export function addApi(data){
    return service({
        url: '/sys/resource/add',
        method: 'POST',
        data:data

    })
}

