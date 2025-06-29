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
             <el-form-item label="名称" >
                 <el-input v-model="queryData.name" placeholder="名称" clearable  style="width: 200px"/>
             </el-form-item>
             <el-form-item label="类别">
                 <el-select v-model="queryData.levelCode" placeholder="" clearable  style="width: 200px;">
                     <el-option :label="item.label" :value="item.value" v-for="(item,index) in getLevelCodes()" :key="index" /> 
                 </el-select>
             </el-form-item>
             <el-form-item>
                 <el-button type="primary" @click="tableReload()" plain>查询</el-button>
             </el-form-item>
         </el-form>
     </el-card>
    
     <el-card >
         <el-table :data="tableData"  :border="true " style="width: 100%;" default-expand-all row-key="id"  
            @selection-change="handleSelectionChange"
            :tree-props="treeProps"
            v-if="tableData.length > 0">
            <el-table-column type="selection" width="55" />
            <el-table-column prop="name" label="名称"/>
             <el-table-column prop="type" label="类型"  width="80">
                 <template #default="scope">
                     {{ getLevelCodeLabel(scope.row.levelCode) }}
                 </template>
             </el-table-column>
             <el-table-column prop="directorName" label="负责人" />
             <el-table-column prop="sortNum" label="排序" />
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
import { getLevelCodeLabel, getLevelCodes } from './dept';

import { treeApi} from "~/api/sys/dept";

const queryData = ref({
     name: '',
     levelCode: '',
 })
const dialogVisible = ref(false);
const tableData = ref([])

// 选择的角色
const tableSelect = ref([])
const handleSelectionChange = (rows) => {
    tableSelect.value = rows;
}

const treeProps = reactive({
  checkStrictly: true,
})

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
   treeApi(param).then(res=> {
       pagination.value.total = res.total;
       tableData.value = res.rows;
   })
}


const open = (record) => {
    tableReload();
    dialogVisible.value = true;
}

const emits = defineEmits(["deptSelectData"])
const submitSelectData = () => {
    emits('deptSelectData', tableSelect.value);
    dialogVisible.value = false;
}


defineExpose({open})

  

</script>

<style  scoped>

</style>