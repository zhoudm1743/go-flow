import service from "~/utils/axios";


export function pageApi(params){
    return service({
        url: '/sys/position/page',
        params

    })
}

export function addApi(data){
    return service({
        url: '/sys/position/add',
        method: 'POST',
        data:data

    })
}

export function editApi(data){
    return service({
        url: '/sys/position/edit',
        method: 'POST',
        data:data

    })
}

export function deleteApi(data){
    return service({
        url: '/sys/position/delete',
        method: 'POST',
        data:data

    })
}
