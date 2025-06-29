import service from "~/utils/axios";
/**
 * 流程列表 分组
 * @param {*} params 
 * @returns 
 */
export function listApi(params){
    return service({
        url: '/flow/def/list',
        params

    })
}