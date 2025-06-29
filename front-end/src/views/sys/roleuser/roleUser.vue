<template>

    <div  style="display: flex; flex-direction: column; gap: 10px; min-height: 400px;">
        <el-card  :body-style="{'padding-bottom': '0px'}">
         <el-form :inline="true" :model="formInline" class="demo-form-inline" label-suffix=":"> 
            <el-form-item label="用户名" prop="username">
                 <el-input v-model="queryData.username" placeholder="用户名"  />
             </el-form-item>
             <el-form-item label="昵称" prop="nick">
                 <el-input v-model="queryData.nick" placeholder="昵称"  />
             </el-form-item>
             <el-form-item label="手机号" prop="phone">
                 <el-input v-model="queryData.phone" placeholder="手机号"  />
             </el-form-item>
             <el-form-item >
                 <el-button type="primary" @click="tableReload()" plain>查询</el-button>
                 <!-- <el-button  style="margin-left: 10px;" type="primary" @click="restform()" plain>重置</el-button> -->
             </el-form-item>
             <el-form-item style="float: right;">
                        <el-button  plain  icon="ArrowLeftBold" @click="goBack()">返回</el-button>
                    </el-form-item>
         </el-form>
     </el-card>
     <el-row style="gap: 10px;  flex-wrap: nowrap;">
        <el-card  style="flex: 0 0 200px;">
            <el-tree
                style="max-width: 600px"
                :data="roleList"
                :props="defaultProps"
                @node-click="handleNodeClick"
                node-key="id"
                empty-text="数据为空"
                :highlight-current="true"
                :current-node-key="currentNodeKey"
            />
        </el-card>
        <el-card  style="flex: 1 1 auto;">
            <el-row class="row-bg" justify="space-between" style="margin-bottom: 16px;">
                <span>
                    <el-button type="primary" @click="openGrantUserForm()" plain>新增</el-button>
                </span>
                <span>
                
                </span>
            </el-row>
            <el-table 
                :data="tableData"  
                :border="true" 
                row-key="id"   
                @selection-change="handleSelectionChange"
                v-if="tableData.length > 0">
                <el-table-column type="selection" width="55" />
                <el-table-column prop="username" label="用户名"/>
                <el-table-column prop="nick" label="昵称" />
                <el-table-column prop="phone" label="手机号" />
                <el-table-column prop="email" label="邮箱" />
                <el-table-column prop="roleStatus" label="状态" >
                    <template #default="scope">
                        <el-switch v-model="scope.row.status"
                            disabled 
                            active-value="1"
                            inactive-value="0"
                        />

                    </template>
                </el-table-column>
                <el-table-column prop="createTime" label="创建时间" >
                </el-table-column>
                <el-table-column  label="操作" >
                    <template  #default="scope">
                            <el-button  link type="primary" size="small" @click="revokeUser(scope.row.roleUserId)">
                            取消授权
                            </el-button>
                    </template>
                </el-table-column>
            </el-table>

            <el-empty description="无数据"   v-else></el-empty>
    
            <div style="padding: 5px 0;"></div>
            <el-row>
                
                <el-col :span="24">
                    <el-pagination 
                        background layout="prev, pager, next,sizes,  total" 
                        :total="pagination.total"
                        :page-sizes="[20, 50, 100, 200]"
                        v-model:page-size="pagination.size"
                        v-model:current-page="pagination.current"
                        @change="pageTable"
                    />
                </el-col>
            </el-row>
        </el-card>
   </el-row>
     
    </div>
    <grantUser ref="grantUserRef" @tableReload="tableReload"></grantUser>
</template>

<script setup>
 import { reactive, ref ,onMounted, defineExpose, watch } from "vue"
 import grantUser from "./grantUser.vue";
 import { pageApi,revokeUserApi } from "~/api/sys/roleuser";
 
 import { listApi } from "~/api/sys/role";
 import router from  '~/router';

 const currentNodeKey = ref(null) //ref(router.currentRoute.value.query.roleId);
 const grantUserRef =ref(null);
 const queryData = ref({
     roleCode: '',
     roleName: '',
     roleStatus: '',
 })

const defaultProps = ref({
    label: 'roleName',
})



 const dialogVisible = ref(false)
 const tableData = ref([])
 
 const pagination = ref({
    size: 20,
    current:   1,
    total:0,
})

const revokeUser = (id) => {
    revokeUserApi([{id: id}]).then(res=> {
        tableReload();
    })
}

const roleList = ref([]);


onMounted(()=>{
    listApi({}).then(result => {
        roleList.value = result.rows;
    })
    
})

const tableSelect = ref([])
const handleSelectionChange = (rows) => {
    tableSelect.value = rows;
}

const tableReload = () => {
    const param = {};
    Object.assign(param, queryData.value);
    param.current = pagination.value.current;
    param.size = pagination.value.size;
    param.roleId = currentNodeKey.value;
    pageApi(param).then(res=> {
        pagination.value.total = res.total;
        tableData.value = res.rows;
    })
}

 

const changeRoleStatus = (value, row) => {
    // if(value == '1'){
    //     enableApi([{id: row.id}]).then(()=> {
    //     }).catch(e=> {
    //         row.roleStatus = '0';
    //     })
    // }else{
    //     disableApi([{id: row.id}]).then(()=> {
    //     }).catch(e=> {
    //         row.roleStatus = '1';
    //     })
    // }
}

//   打开页面
const openGrantUserForm = (role) => {
    if(currentNodeKey.value != null && currentNodeKey.value != undefined){
        grantUserRef.value.openGrantUserForm(currentNodeKey);
    }
  
    
}

const handleNodeClick = (curr) => {
    pagination.value.current = 1;
    currentNodeKey.value = curr.id;
    tableReload()
}

const role = (roleId) => {
    currentNodeKey.value = roleId
    tableReload();
}
/**
 * 授权
 */
const submitGrantUsre = () => {

}

const emits = defineEmits(['enableShow'])

const goBack = () => {
    // router.go(-1);
    emits('enableShow', true);
}

defineExpose({tableReload, role})

</script>

<style lang="scss" scoped>

</style>