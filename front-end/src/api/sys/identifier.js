import service from "~/utils/axios";


/**
 * 权限标识
 * @param {menuId} params 
 * @returns 
 */

export function listApi(params){
    return service({
        url: '/sys/identifier/list',
        params

    })
}

export function addApi(data){
    return service({
        url: '/sys/identifier/add',
        method: 'POST',
        data:data

    })
}

export function editApi(data){
    return service({
        url: '/sys/identifier/edit',
        method: 'POST',
        data:data

    })
}

export function deleteApi(data){
    return service({
        url: '/sys/identifier/delete',
        method: 'POST',
        data:data

    })
}

export function detailApi(params){
    return service({
        url: '/sys/identifier/detail',
        params

    })
}