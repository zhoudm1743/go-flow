<template>
    <div  style="display: flex; flex-direction: column; gap: 10px; min-height: 400px;">
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
        <el-row style="gap: 10px;  flex-wrap: nowrap;">
            <el-card  style="flex: 0 0 240px;">
                <el-tree
                    style="max-width: 600px"
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
            <el-card  style="flex: 1 1 auto;">
                <el-table :data="tableData"  :border="true" style="width: 100%;"    v-if="tableData.length > 0">
                    <el-table-column prop="positionName" label="职位名称"/>
                    <el-table-column prop="deptName" label="所属部门" />
                    <el-table-column prop="sortNum" label="排序" />
                    <el-table-column  label="操作" >
                        <template  #default="scope">
                            <el-button link type="primary" size="small" @click="selectPosition(scope.row)">
                                选择
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
        </el-row>
    </div>

</template>

<script setup>
 import { reactive, ref ,onMounted, defineExpose,defineEmits } from "vue"

 import PositionForm from './form.vue'
 import { pageApi ,deleteApi } from "~/api/sys/position";
 import { treeApi } from "~/api/sys/dept";
 import { getPositionInfo } from "./position";


 const emits  = defineEmits(['positionSelect']);
 const selectPosition = (position) => {
     const detpArry = getPositionInfo(position, deptList.value);
     let deptIds = null;
     let cmpId = null;
     let orgName = null;
     let cmpName = null;
     for(var i in detpArry){
        if(detpArry[i].levelCode == 'dept'){
            if(orgName && orgName.length  > 0){
                deptIds = detpArry[i].id  + "," + deptIds;
                orgName = detpArry[i].name  + " > "+  orgName;
            }else{
                deptIds = detpArry[i].id;
                orgName = detpArry[i].name ;
            }
        }else if(detpArry[i].levelCode == 'cmp'){
            cmpId = detpArry[i].id;
            cmpName = detpArry[i].name;
            break;
        }
     }
     
     const info = {
        orgName : orgName,
        cmpName: cmpName,
        deptIds:deptIds,
        cmpId:cmpId,
        positionId: position.id,
        positionName: position.positionName
     }
     emits('positionSelect',info);
 }


 const queryData = ref({
    positionName: '',
 })
 const tableData = ref([])
 
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