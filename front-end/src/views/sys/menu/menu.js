const menuTypes = [
    {
        value: "1",
        label: '目录',
    },
    {
        value: "2",
        label: '菜单',
    },
    {
        value: "3",
        label: '内链',
    },
    {
        value: "4",
        label: '外链',
    },
]

const menuStatus = [
    {
        value: "1",
        label: '显示',
    },
    {
        value: "0",
        label: '隐藏',
    },
]

export function getTypes() {
    return menuTypes
}   

export function getTypeLabel(value) {
    return menuTypes.filter(item => item.value === value)[0].label
}   

export function getStatus() {
    return menuStatus
}

export function getStatusLabel(value) {
    return menuStatus.filter(item => item.value === value)[0].label
}