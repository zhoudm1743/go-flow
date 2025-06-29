export const getFlowStatus = (num) => {
    switch(num){
        case 0:
            return '待提交';
        case 1:
            return '审批中';
        case 2:
            return '审批通过';
        case 3:
            return '自动完成';
        case 4:
            return '终止';
        case 5:
            return '作废';
        case 6:
            return '撤销';
        case 7:
            return '取回';
        case 8:
            return '已完成';
        case 9:
            return '已退回';
        case 10:
            return '失效';
        case 99:
            return '已拒绝';
        default:
            return '未知状态';
    }
}

export const getFlowColor = (num) => {
    switch(num){
        case 0:
            return 'RGB(255, 190, 152)';
        case 1:
            return 'RGB(255, 57, 82)';
        case 2:
            return 'RGB(140, 43, 50)';
        case 3:
            return 'RGB(202, 202, 202)';
        case 4:
            return 'RGB(111, 93, 68)';
        case 5:
            return 'RGB(255, 69, 76)';
        case 6:
            return 'RGB(171, 71, 60)';
        case 7:
            return 'RGB(255, 149, 0)';
        case 8:
            return 'RRGB(38, 40, 221)';
        case 9:
            return 'RGB(255, 105, 164)';
        case 10:
            return 'RGB(0, 122, 204)';
        case 99:
            return 'RGB(255, 0, 0)';
        default:
            return 'RGB(255, 255, 0)';
    }
}