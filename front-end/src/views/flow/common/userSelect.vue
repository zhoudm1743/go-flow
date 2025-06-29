<template>
    <el-dialog
         v-model="dialogVisible"
        :title="title"
        width="920"
        destroy-on-close	
        :before-close="handleClose"
    >
       
    <div  style="display: flex; flex-direction: row; gap: 2px;">
        <div  style="display: flex; flex-direction: column; gap: 10px;">
            <el-card  :body-style="{'padding-bottom': '0px'}">
                <el-form :inline="true" :model="formInline" class="demo-form-inline" label-suffix=":"> 
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
            
            <el-card >
                <el-table :data="tableData"  :border="true" style="width: 100%;"
                    @selection-change="handleSelectionChange"
                    ref="userTableRef"
                        v-if="tableData.length > 0">
                    <el-table-column type="selection" width="55">
                        <template #default="scope" v-if="mutiSelect === false" >
                            <el-checkbox v-model="scope.row._checked"  @click="singleSelect(scope.row)"/>
                        </template>
                    </el-table-column>
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
        <div style="width: 200px;">
            <el-card style="height: 100%;"  :body-style="{'padding': '0px 20px'}">
                <template #header>
                    <div class="card-header">
                        <span>已选</span>
                    </div>
                </template>
                <p v-for="item,index in tableSelect" :key="index" class="text item">
                    <el-tag closable type="primary"   @close="tagClose(item.id)">
                        {{ item.nick }}
                        </el-tag>
                </p>
            </el-card>
            
        </div>
    </div>
        
        <template #footer>
        <div class="dialog-footer">
            <el-space>
                <el-button @click="dialogVisible = false">取消</el-button>
                <el-button type="primary" @click="submitSelectData()">
                确定
                </el-button>
            </el-space>
        </div>
        </template>
    </el-dialog>

</template>

<script setup>
import { reactive, ref ,onMounted, defineExpose, defineEmits, nextTick } from "vue"
import { selectUsersApi } from "~/api/flow/flowTask";
const queryData = ref({
    taskId:'',
    roleCode: '',
    roleName: '',
    roleStatus: '',
})

const mutiSelect = ref(true);
const dialogVisible = ref(false);
const tableData = ref([])
const userTableRef = ref(null)
const title = ref("办理人选择")
// 选择的角色
const tableSelect = ref([])
const handleSelectionChange = (rows, v) => {
    if(mutiSelect.value){
        tableSelect.value = rows;
    }
   
}
const tagClose = (id)=> {
   const rows = tableSelect.value;
   for (let i = rows.length - 1; i >= 0; i--) {
        if (rows[i].id == id) { // 移除
            userTableRef.value.toggleRowSelection(rows[i],false);
        }
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

const pagination = ref({
   size: 20,
   current:   1,
   total:0,
})



const tableReload = () => {
   const param = {};
   Object.assign(param, queryData.value);
   selectUsersApi(param).then(res=> {
       pagination.value.total = res.rows.size;
       tableData.value = res.rows;
   })
}


const open = (taskId) => {
    queryData.value.taskId = taskId;
    tableSelect.value = [];
    tableReload();
    dialogVisible.value = true;
}

const emits = defineEmits(["roleSelectData"])
const submitSelectData = () => {
    emits('selectData', tableSelect.value, mutiSelect.value);
    dialogVisible.value = false;
}


defineExpose({open})

  

</script>

<style  scoped>

</style>