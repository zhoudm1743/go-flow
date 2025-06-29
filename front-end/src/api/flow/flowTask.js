import service from "~/utils/axios";

export function todoApi(params){
    return service({
        url: '/flow/task/todo',
        params

    })
}

export function doneApi(params){
    return service({
        url: '/flow/task/done',
        params

    })
}
export function myApi(params){
    return service({
        url: '/flow/task/my',
        params

    })
}

/**
 * 办理任务
 * @param {} data 
 * @returns 
 */
export function skipApi(data){
    return service({
        url: '/flow/task/skip',
        method: 'POST',
        data: data
    })
}

export function rejectApi(data){
    return service({
        url: '/flow/task/reject',
        method: 'POST',
        data: data
    })
}

export function backApi(data){
    return service({
        url: '/flow/task/back',
        method: 'POST',
        data: data
    })
}

/**
 *  转办
 * @param {*} data 
 * @returns 
 */
export function transferApi(data){
    return service({
        url: '/flow/task/transfer',
        method: 'POST',
        data: data
    })
}


/**
 *  委派
 * @param {*} data 
 * @returns 
 */
export function deputeApi(data){
    return service({
        url: '/flow/task/depute',
        method: 'POST',
        data: data
    })
}

/**
 *  加签
 * @param {*} data 
 * @returns 
 */
export function addSignApi(data){
    return service({
        url: '/flow/task/addSign',
        method: 'POST',
        data: data
    })
}

/**
 *  减签
 * @param {*} data 
 * @returns 
 */
export function reduSignApi(data){
    return service({
        url: '/flow/task/reduSign',
        method: 'POST',
        data: data
    })
}




/**
 * 审批记录
 * @param {*} params 
 * @returns 
 */
export function hisListApi(params){
    return service({
        url: '/flow/task/hisList',
        params
    })
}

/**
 *  转办 ，委派 等， 选择用户
 * @param {*} params 
 * @returns 
 */
export function selectUsersApi(params){
    return service({
        url: '/flow/task/selectUsers',
        params
    })
}