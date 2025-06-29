<template>
    <el-row  style=" flex-direction: column; gap: 10px;">
        <el-card  :body-style="{'padding-bottom': '0px'}">
            <el-form :inline="true" :model="formInline" class="demo-form-inline" label-suffix=":"> 
                <el-form-item label="职位名称" prop="positionName">
                    <el-input v-model="queryData.positionName" placeholder="用户名"  />
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
                    :expand-on-click-node="false"
                    @node-click="handleNodeClick"
                    default-expand-all
                    node-key="id"
                    empty-text="数据为空"
                    :highlight-current="true"
                    :current-node-key="currentNodeKey"
                />
            </el-card>
            <el-card style="flex: 1;">
                <el-row  style="flex: 1; position: relative;">
                    <el-row class="row-bg" justify="space-between" style="margin-bottom: 16px; ">
                        <span>
                            <el-button type="primary" @click="addPosition()">新增</el-button>
                        </span>
                        <span>
                        
                        </span>
                    </el-row >
                
                    <el-table  :data="tableData"  :border="true"   v-if="tableData.length > 0">
                        <el-table-column prop="positionName" label="职位名称"/>
                        <el-table-column prop="deptName" label="所属部门"/>
                        <el-table-column prop="sortNum" label="排序" />
                        <el-table-column prop="createTime" label="创建时间" >
                        </el-table-column>
                        <el-table-column  label="操作"  width="180" >
                            <template  #default="scope">
                                <el-button link type="primary" size="small" @click="openForm(scope.row)">
                                    编辑
                                    </el-button>
                                    <el-button  link type="primary" size="small" @click="handleDelete(scope.row.id)">
                                    删除
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
                </el-row>
            </el-card>
        </el-row>
    </el-row>
    <PositionForm ref="formRef" :deptList="deptList" @tableReload="tableReload"></PositionForm>
    <ResourcePage ref="resourceRef" @tableReload="tableReload"></ResourcePage>
</template>

<script setup>
 import { reactive, ref ,onMounted, defineExpose } from "vue"
 import ResourcePage from '../menu/resource/index.vue'

 import PositionForm from './form.vue'
 import { pageApi ,deleteApi } from "~/api/sys/position";
 import { treeApi } from "~/api/sys/dept";
 import { positionType } from "../menu/resource/index.js";
 const queryData = ref({
    positionName: '',
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
    resourceRef.value.openPage({fromType: positionType, fromId: record.id});
}



const tableReload = () => {
    const param = {};
    Object.assign(param, queryData.value);
    param.deptId = currentNodeKey.value;
    param.current = pagination.value.current;
    param.size = pagination.value.size;
    pageApi(param).then(res=> {
        pagination.value.total = res.total;
        tableData.value = res.rows;
    })
}
// 办理const 
const handleDelete = (id) => {
    deleteApi([{"id":id}]).then(()=> {
        tableReload();
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

const addPosition = () => {
    formRef.value.addPosition(currentNodeKey);
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