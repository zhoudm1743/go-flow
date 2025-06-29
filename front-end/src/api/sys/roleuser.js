import service from "~/utils/axios";

export function pageApi(params){
    return service({
        url: '/sys/roleUser/page',
        params

    })
}


export function canGrantPageApi(params){
    return service({
        url: '/sys/roleUser/canGrantPage',
        params

    })
}


export function revokeUserApi(data){
    return service({
        url: '/sys/roleUser/revoke',
        method: 'POST',
        data:data

    })
}

