<template>
    <el-dialog
         v-model="dialogVisible2"
        :title="title"
        width="1200px"
        destroy-on-close	
        :before-close="handleClose"
    >
       
    <div  style="display: flex; flex-direction: column; gap: 10px;">
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
        <el-row style="gap: 10px;  flex-wrap: nowrap;">
            <el-card  style="flex: 0 0 240px;">
                <el-tree
                    style="max-width: 600px"
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
            <el-card  style="flex: 1 1 auto;">
    
                <el-row class="row-bg" justify="space-between" style="margin-bottom: 16px;">
                    <span>
                        <el-button type="primary" @click="openForm()">新增</el-button>
                    </span>
                    <span>
                    
                    </span>
                </el-row>
                <el-table :data="tableData"  :border="true" style="width: 100%;"
                    @selection-change="handleSelectionChange"
                        v-if="tableData.length > 0">
                    <el-table-column type="selection" width="55">
                        <template #default="scope" v-if="mutiSelect === false" >
                            <el-checkbox v-model="scope.row._checked"  @click="singleSelect(scope.row)"/>
                        </template>
                    </el-table-column>
                    <el-table-column prop="username" label="用户名"/>
                    <el-table-column prop="nick" label="昵称" />
                    <el-table-column prop="positionList" label="任职信息" min-width="260">
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
                    <el-table-column prop="roleStatus" label="状态" width="80" >
                        <template #default="scope">
                            <el-switch v-model="scope.row.status" 
                                active-value="1"
                                inactive-value="0"
                                @change="changeUserStatus($event, scope.row)"
                            />

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

        <template #footer>
            <div class="dialog-footer">
                <el-space>
                    <el-button @click="dialogVisible2 = false">取消</el-button>
                    <el-button type="primary" @click="submitSelectData()">
                    确定
                    </el-button>
                </el-space>
            </div>
        </template>
    </el-dialog>
</template>

<script setup>
 import { defineEmits, ref ,onMounted, defineExpose } from "vue"

 import { pageApi ,deleteApi,} from "~/api/sys/user";
 import { treeApi } from "~/api/sys/dept";
 const queryData = ref({
     roleCode: '',
     roleName: '',
     roleStatus: '',
 })
 const tableData = ref([])
 const mutiSelect = ref(true);
 const dialogVisible2 = ref(false);

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
 
const openForm = (record) => {
    formRef.value.openForm(Object.assign({}, record));
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

const tableSelect = ref([])
const handleSelectionChange = (rows, v) => {
    if(mutiSelect.value){
        tableSelect.value = rows;
    }
   
}
const singleSelect = (row) => {
  
    tableData.value.forEach((element, index, array) => {
        if(element.id !== row.id){
            element._checked = false;
        }
  // 执行操作
    });
    tableSelect.value[0] = row;
}


const open = (type) => {
     if(typeof type === 'boolean' || typeof type === Boolean ){
        mutiSelect.value = type;
     }else{
        mutiSelect.value = true;
     }
    tableSelect.value = [];
    tableReload();
    dialogVisible2.value = true;
}

const emits = defineEmits(["roleSelectData"])
const submitSelectData = () => {
    emits('roleSelectData', tableSelect.value, mutiSelect.value);
    dialogVisible2.value = false;
}


defineExpose({open})


</script>

<style lang="scss" scoped>

</style>