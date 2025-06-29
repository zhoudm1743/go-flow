import service from "~/utils/axios";

// export function listApi(params){
//     return service({
//         url: '/sys/role/page',
//         params

//     })
// }

export function addApi(data){
    return service({
        url: '/sys/menu/add',
        method: 'POST',
        data:data

    })
}

export function editApi(data){
    return service({
        url: '/sys/menu/edit',
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
        url: '/sys/menu/delete',
        method: 'POST',
        data:data

    })
}

export function pageApi(params){
    return service({
        url: '/sys/menu/page',
        params
    })
}

export function treeApi(params){
    return service({
        url: '/sys/menu/tree',
        params
    })
}

export function menuRouter(params){
    return service({
        url: '/sys/menu/menuRouter',
        params
    })
}



export function allmenus(){
    return service({
        url: '/sys/menu/allMenus',
    })
}
