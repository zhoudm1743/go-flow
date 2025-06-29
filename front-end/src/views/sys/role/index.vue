<template>
    <div>
        <div  style="display: flex; flex-direction: column; gap: 10px;"  v-show="showAbled">
            <el-card  :body-style="{'padding-bottom': '0px'}">
                <el-form :inline="true" :model="formInline" class="demo-form-inline" label-suffix=":"> 
                    <el-form-item label="角色编码">
                        <el-input v-model="queryData.roleCode" placeholder="角色编码" clearable />
                    </el-form-item>
                    <el-form-item label="角色名称">
                        <el-input v-model="queryData.roleName" placeholder="角色名称" clearable />
                    </el-form-item>
                    <el-form-item>
                        <el-button type="primary" @click="tableReload()" plain>查询</el-button>
                    </el-form-item>
                </el-form>
            </el-card>
        
            <el-card >
                <el-row class="row-bg" justify="space-between" style="margin-bottom: 16px;">
                    <span>
                        <el-button type="primary" @click="openForm()" >新增</el-button>
                        <el-button type="primary" @click="openGrantUserForm({id:null})">
                                授权用户
                        </el-button>
                    </span>
                    <span>
                    
                    </span>
                </el-row>
                <el-table v-direct-auto-height :data="tableData"  :border="true" style="width: 100%;"    v-if="tableData.length > 0">
                    <el-table-column prop="roleCode" label="角色编码"/>
                    <el-table-column prop="roleName" label="角色名称" />
                    <el-table-column prop="sortNum" label="排序" />
                    <el-table-column prop="roleStatus" label="状态" >
                        <template #default="scope">
                            <el-switch v-model="scope.row.roleStatus" 
                                active-value="1"
                                inactive-value="0"
                                @change="changeRoleStatus($event, scope.row)"
                            />

                        </template>
                    </el-table-column>
                    <el-table-column prop="createTime" label="创建时间" >
                    </el-table-column>
                    <el-table-column  label="操作" >
                        <template  #default="scope">
                            <el-button link type="primary" size="small" @click="openForm(scope.row)">
                                编辑
                                </el-button>
                                <el-button :disabled="scope.row.deleted == '1'" link type="primary" size="small" @click="handleDelete(scope.row.id)">
                                删除
                                </el-button>
                                <el-button  link type="primary" size="small" @click="openGrantUserForm(scope.row)">
                                授权用户
                                </el-button>
                                <el-button  link type="primary" size="small" @click="openGrantResourcePage(scope.row)">
                                授权资源
                                </el-button>
                        </template>
                    </el-table-column>
                </el-table>

                <el-empty description="无数据"   v-else></el-empty>
                <div style="padding: 5px 0;"></div>
                <el-row>
                    <el-col :span="24">
                        <el-pagination 
                            background layout="prev, pager, next" 
                            :total="pagination.total"
                            v-model:page-size="pagination.size"
                            v-model:current-page="pagination.current"
                            @change="pageTable"
                        />
                    </el-col>
                </el-row>
            </el-card>
        </div>
        <div v-show="!showAbled">
            <RoleUser ref="roleUserComponent"  v-show="!showAbled" @enableShow="enableShow" ></RoleUser>
        </div>
        
        <RoleForm ref="roleFormRef" @tableReload="tableReload" ></RoleForm>
        <ResourcePage ref="resourceRef" @tableReload="tableReload"></ResourcePage>
    </div>
</template>

<script setup>
 import { reactive, ref ,onMounted, defineExpose } from "vue"
 import RoleUser from "../roleuser/roleUser.vue";
 import RoleForm from './form.vue'
 import ResourcePage from '../menu/resource/index.vue'

 import { pageApi, addApi, editApi, enableApi, disableApi,deleteApi } from "~/api/sys/role";
 import { roleType } from "../menu/resource/index.js";

 const showAbled = ref(true);

 const queryData = ref({
     roleCode: '',
     roleName: '',
     roleStatus: '',
 })
 const tableData = ref([])
 
 const pagination = ref({
    size: 20,
    current:   1,
    total:0,
})

const roleUserComponent = ref();
const resourceRef = ref(null);
onMounted(()=>{
    tableReload();
    
})

const tableReload = () => {
    const param = {};
    Object.assign(param, queryData.value);
    param.current = pagination.value.current;
    param.size = pagination.value.size;
    pageApi(param).then(res=> {
        pagination.value.total = res.total;
        tableData.value = res.rows;
    })
}

const enableShow = (value) => {
    showAbled.value = value;
}

// 办理
const handleClick = (row) => {
    let data = {
        id: row.id,
        message: "ok",
        skipType: "PASS",
        variable: "这是一个未知数"
    }
    skipApi(data).then(res=> {
        alert(res.msg);
    })
}

const roleFormRef = ref(null);


const openForm = (record) => {
    roleFormRef.value.openForm(Object.assign({}, record));
}

const openGrantUserForm = (record) => {
    roleUserComponent.value.role(record.id);
    showAbled.value = false;
    // router.addRoute('Layout', { path: '/roleuser/roleUser', component: roleUserComponent.value })
    // router.replace({path:`/roleuser/roleUser`, query:{roleId: record.id}})
    // 
}

const openGrantResourcePage = (record) => {
    resourceRef.value.openPage({fromType: roleType, fromId: record.id});
}


const changeRoleStatus = (value, row) => {
    if(value == '1'){
        enableApi([{id: row.id}]).then(()=> {
        }).catch(e=> {
            row.roleStatus = '0';
        })
    }else{
        disableApi([{id: row.id}]).then(()=> {
        }).catch(e=> {
            row.roleStatus = '1';
        })
    }
}

defineExpose({tableReload, showAbled})

</script>

<style lang="scss" scoped>

</style>