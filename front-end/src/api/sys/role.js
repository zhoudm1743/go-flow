import service from "~/utils/axios";

// export function listApi(params){
//     return service({
//         url: '/sys/role/page',
//         params

//     })
// }

export function addApi(data){
    return service({
        url: '/sys/role/add',
        method: 'POST',
        data:data

    })
}

export function editApi(data){
    return service({
        url: '/sys/role/edit',
        method: 'POST',
        data:data

    })
}

/**
 * 禁用
 * @param {*} data 
 * @returns 
 */
export function disableApi(data){
    return service({
        url: '/sys/role/disable',
        method: 'POST',
        data:data

    })
}

/**
 * 启用
 * @param {*} data 
 * @returns 
 */
export function enableApi(data){
    return service({
        url: '/sys/role/enable',
        method: 'POST',
        data:data

    })
}

export function deleteApi(data){
    return service({
        url: '/sys/role/delete',
        method: 'POST',
        data:data

    })
}

export function pageApi(params){
    return service({
        url: '/sys/role/page',
        params
    })
}

export function listApi(params){
    return service({
        url: '/sys/role/list',
        params
    })
}



export function grantUserApi(data){
    return service({
        url: '/sys/role/grant/user',
        method: 'POST',
        data:data

    })
}