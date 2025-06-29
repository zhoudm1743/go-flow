import service from "~/utils/axios";

// export function listApi(params){
//     return service({
//         url: '/sys/role/page',
//         params

//     })
// }

export function addApi(data){
    return service({
        url: '/sys/dept/add',
        method: 'POST',
        data:data

    })
}

export function editApi(data){
    return service({
        url: '/sys/dept/edit',
        method: 'POST',
        data:data

    })
}

/**
 * 禁用
 * @param {*} data 
 * @returns 
 */
// export function disableApi(data){
//     return service({
//         url: '/sys/role/disable',
//         method: 'POST',
//         data:data

//     })
// }

/**
 * 启用
 * @param {*} data 
 * @returns 
 */
// export function enableApi(data){
//     return service({
//         url: '/sys/role/enable',
//         method: 'POST',
//         data:data

//     })
// }

export function deleteApi(data){
    return service({
        url: '/sys/dept/delete',
        method: 'POST',
        data:data

    })
}

export function pageApi(params){
    return service({
        url: '/sys/dept/page',
        params
    })
}

export function deptListApi(params){
    return service({
        url: '/sys/dept/list',
        params
    })
}

export function treeApi(params){
    return service({
        url: '/sys/dept/tree',
        params
    })
}
