import service from "~/utils/axios";

export function startApi(data){
    return service({
        url: '/flow/instance/start',
        method: 'POST',
        data: data

    })
}