<template>
    <el-dialog
         v-model="dialogVisible"
        :title="title"
        width="50%"
        :before-close="handleClose"
    >
       
    <div  style="display: flex; flex-direction: column; gap: 10px;">
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
         <el-table :data="tableData"  :border="true" style="width: 100%;"    
            @selection-change="handleSelectionChange"
            v-if="tableData.length > 0">
            <el-table-column type="selection" width="55" />
            <el-table-column prop="roleCode" label="角色编码"/>
             <el-table-column prop="roleCode" label="角色编码"/>
             <el-table-column prop="roleName" label="角色名称" />
             <el-table-column prop="sortNum" label="排序" />
             <el-table-column prop="roleStatus" label="状态" >
                <template #default="scope">
                    <el-switch v-model="scope.row.roleStatus" 
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
import { reactive, ref ,onMounted, defineExpose, defineEmits } from "vue"


import { pageApi} from "~/api/sys/role";

import router from "~/router";
const queryData = ref({
    roleCode: '',
    roleName: '',
    roleStatus: '',
})

const dialogVisible = ref(false);
const tableData = ref([])

// 选择的角色
const tableSelect = ref([])
const handleSelectionChange = (rows) => {
    tableSelect.value = rows;
}


const pagination = ref({
   size: 20,
   current:   1,
   total:0,
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


const open = (record) => {
    tableReload();
    dialogVisible.value = true;
}

const emits = defineEmits(["roleSelectData"])
const submitSelectData = () => {
    emits('roleSelectData', tableSelect.value);
    dialogVisible.value = false;
}


defineExpose({open})

  

</script>

<style  scoped>

</style>