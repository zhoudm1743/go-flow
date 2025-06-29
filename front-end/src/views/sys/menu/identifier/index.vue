<template>
    <el-drawer
        v-model="drawer"
        title="权限标识"
        direction="rtl"
        size="600"
    >
        <div  style="display: flex; flex-direction: column; gap: 20px; ">
            <el-row>
                <el-col :span="24">
                    <el-button type="primary" @click="handlerAdd('1')">新增接口权限</el-button>
                    <el-button color="#2a53fe" @click="handlerAdd('2')">新增按钮权限</el-button>
                </el-col>
            </el-row>
            <el-row>
                <el-col :span="24">
                    <el-table :data="tableData" style="width: 100%">
                        <el-table-column prop="code" label="编码" />
                        <el-table-column prop="name" label="名称"  />
                        <el-table-column prop="type" label="类别" >
                            <template #default="scope">
                                <el-tag  :type="scope.row.type === '1' ? 'primary' : 'success'"  effect="dark">{{ scope.row.type == '1' ? '接口' : '按钮' }}</el-tag>
                            </template>
                        </el-table-column>
                        <el-table-column prop="sortNum" label="排序" />
                        <el-table-column prop="operate" label="操作" width="110">
                            <template #default="scope" >
                                <el-button link type="primary"  @click="handlerEdit(scope.row)">编辑</el-button>
                                <el-button link type="primary"  @click="handlerDelete(scope.row.id)">删除</el-button>
                            </template>
                        </el-table-column>

                    </el-table>    
                </el-col>
            </el-row>
        </div>
    </el-drawer>
    <Form ref="identifierFormRef" @tableReload="tableReload"></Form>
</template>

<script setup>
import { defineExpose, onMounted, ref  } from 'vue';
import {listApi, addApi, editApi, deleteApi, detailApi} from '~/api/sys/identifier'
import Form from './form.vue';
const drawer = ref(false)

const identifierFormRef = ref(null)
const tableData = ref([])
const currMenuId = ref('')


const tableReload = () => {
    listApi({menuId: currMenuId.value}).then(result => {
        tableData.value = result.rows;
    })
}

const openDrawer = (menuId) => {
    currMenuId.value = menuId
    tableReload();
    drawer.value = true
}

const closeDrawer = () => {
    drawer.value = false
}

const handlerEdit = (row) => { 
    identifierFormRef.value.openForm(row)
}

const handlerDelete = (id) => {
    deleteApi([{id: id}]).then(result => {
        tableReload()
    })
}

const handlerAdd = (type) => {
    identifierFormRef.value.openForm({menuId: currMenuId.value, type: type});
}

defineExpose({
    openDrawer,
    closeDrawer,tableReload
})

</script>

<style lang="scss" scoped>

</style>