import service from "~/utils/axios";

export function listApi(params){
    return service({
        url: '/sys/deptUser/list',
        params

    })
}

