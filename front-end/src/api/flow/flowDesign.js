import service from "~/utils/axios";

export function addApi(data){
    return service({
        url: '/flow/design/add',
        method: 'POST',
        data: data

    })
}

export function editApi(data){
    return service({
        url: '/flow/design/edit',
        method: 'POST',
        data: data

    })
}


export function pageApi(params){
    return service({
        url: '/flow/design/page',
        params

    })
}


export function delApi(data){
    return service({
        url: '/flow/design/del',
        method: 'POST',
        data: data

    })
}
export function saveConfigApi(data){
    return service({
        url: '/flow/design/saveConfig',
        method: 'POST',
        data: data

    })
}


export function detailApi(params){
    return service({
        url: '/flow/design/detail',
        params

    })
}
export function publishApi(data){
    return service({
        url: '/flow/design/publish',
        method: 'POST',
        data: data

    })
}

