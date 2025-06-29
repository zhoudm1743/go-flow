const levelCodes = [
    {
        value: "cmp",
        label: '公司',
    },
    {
        value: "dept",
        label: '部门',
    }
]


export function getLevelCodes() {
    return levelCodes
}   

export function getLevelCodeLabel(value) {
    return levelCodes.filter(item => item.value === value)[0].label
}   
