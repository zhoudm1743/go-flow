<template>
    <el-row  style=" flex-direction: column; gap: 10px;">
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
                <el-form-item>
                    <el-button type="primary" @click="tableReload()" plain>查询</el-button>
                </el-form-item>
            </el-form>
        </el-card>
        <el-row style=" gap: 10px;"  >
            <el-card  style="min-width: 280px;">
                <el-tree
                    :data="deptList"
                    :props="defaultProps"
                    @node-click="handleNodeClick"
                    default-expand-all
                    :expand-on-click-node="false"
                    node-key="id"
                    empty-text="数据为空"
                    :highlight-current="true"
                    :current-node-key="currentNodeKey"
                />
            </el-card>
            <el-card  style="flex: 1">
                <el-row  style="flex: 1; position: relative;">
                    <el-row class="row-bg" justify="space-between" style="margin-bottom: 16px;">
                        <span>
                            <el-button type="primary" @click="openForm()">新增</el-button>
                        </span>
                        <span>
                        
                        </span>
                    </el-row>
                    <el-table :data="tableData"  :border="true" style="width: 100%;"    v-if="tableData.length > 0">
                        <el-table-column prop="username" label="用户名"/>
                        <el-table-column prop="nick" label="昵称" />
                        <el-table-column prop="positionList" label="任职信息" min-width="200">
                            <template #default="scope">
                                <template v-for="(item,index) in scope.row.positionList" :key="item.id">
                                    <el-tag  type="danger"  effect="dark">{{ item.cmpName }}</el-tag> |
                                    
                                    <template v-if="item.orgName" >
                                        <el-tag  type="primary"   effect="dark">{{ item.orgName }}</el-tag>   |
                                    </template>
                                    <el-tag  type="success"  effect="dark">{{item.positionName}}</el-tag>
                                    <br  v-if="scope.row.positionList.length - 1 != index" />
                                </template>
                            
                            </template>
                        </el-table-column>
                        <el-table-column prop="phone" label="手机号" />
                        <el-table-column prop="email" label="邮箱" width="200px"/>
                        <el-table-column prop="roleStatus" label="状态" width="80" >
                            <template #default="scope">
                                <el-switch v-model="scope.row.status" 
                                    active-value="1"
                                    inactive-value="0"
                                    @change="changeUserStatus($event, scope.row)"
                                />

                            </template>
                        </el-table-column>
                    
                        <el-table-column  label="操作"  width="180">
                            <template  #default="scope">
                                <el-button link type="primary" size="small" @click="openForm(scope.row)">
                                    编辑
                                </el-button>
                                <el-popconfirm title="确定删除吗?"
                                    @confirm="handleDelete(scope.row.id)"
                                    hide-icon
                                >
                                    <template #reference>
                                        <el-button  link type="primary" size="small">
                                            删除
                                        </el-button>
                                    </template>
                                </el-popconfirm>
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
                </el-row>
            </el-card>
        </el-row>
    </el-row>
    <UserForm ref="formRef" @tableReload="tableReload"></UserForm>
    <ResourcePage ref="resourceRef" @tableReload="tableReload"></ResourcePage>

</template>

<script setup>
 import { reactive, ref ,onMounted, defineExpose } from "vue"
 import ResourcePage from '../menu/resource/index.vue'

 import UserForm from './form.vue'
 import { pageApi ,deleteApi,} from "~/api/sys/user";
 import { treeApi } from "~/api/sys/dept";
import { userType } from "../menu/resource/index.js";
 const queryData = ref({
     roleCode: '',
     roleName: '',
     roleStatus: '',
 })
 const tableData = ref([])
 const resourceRef = ref(null);
 const pagination = ref({
    size: 20,
    current:   1,
    total:0,
})

onMounted(()=>{
    treeApi().then(res=> {
        deptList.value = res.rows;
    });
    tableReload();
    
})

const openGrantResourcePage = (record) => {
    resourceRef.value.openPage({fromType: userType, fromId: record.id});
}


const tableReload = () => {
    const param = {};
    Object.assign(param, queryData.value);
    param.current = pagination.value.current;
    param.size = pagination.value.size;
    param.deptId = currentNodeKey.value;
    pageApi(param).then(res=> {
        pagination.value.total = res.total;
        tableData.value = res.rows;
    })
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

const changeUserStatus = (value, row) => {
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
 
const formRef = ref(null);
const openForm = (record) => {
    formRef.value.openForm(Object.assign({}, record));
}

const handleDelete = (id) => {
    deleteApi([{"id":id}]).then(res=> {
        tableReload();
    })
}   

// 部门配置
const deptList = ref([]);
const currentNodeKey = ref(null);
const defaultProps = ref({
    label: 'name',
})

// 点击部门，查询部门下的人员
const handleNodeClick = (curr) => {
    if(currentNodeKey.value == curr.id){
        currentNodeKey.value = null;
    }else{
        currentNodeKey.value = curr.id;
    }
    tableReload();
}




defineExpose({tableReload})

</script>

<style lang="scss" scoped>

</style>