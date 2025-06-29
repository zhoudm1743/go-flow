import service from "~/utils/axios";

export function addApi(data){
    return service({
        url: '/flow/category/add',
        method: 'POST',
        data: data

    })
}


export function editApi(data){
    return service({
        url: '/flow/category/edit',
        method: 'POST',
        data: data

    })
}

export function pageApi(params){
    return service({
        url: '/flow/category/page',
        params

    })
}


export function delApi(data){
    return service({
        url: '/flow/category/del',
        method: 'POST',
        data: data

    })
}

export function listApi(params){
    return service({
        url: '/flow/category/list',
        params
    })
}