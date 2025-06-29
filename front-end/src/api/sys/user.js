import service from "~/utils/axios";

export function listApi(params){
    return service({
        url: '/sys/user/list',
        params

    })
}

export function pageApi(params){
    return service({
        url: '/sys/user/page',
        params

    })
}
export function selectUsersApi(params){
    return service({
        url: '/sys/user/selectUsers',
        params

    })
}


export function addApi(data){
    return service({
        url: '/sys/user/add',
        method: 'POST',
        data:data

    })
}

export function editApi(data){
    return service({
        url: '/sys/user/edit',
        method: 'POST',
        data:data

    })
}

export function deleteApi(data){
    return service({
        url: '/sys/user/delete',
        method: 'POST',
        data:data

    })
}
export function detailApi(params){
    return service({
        url: '/sys/user/detail',
        params

    })
}
