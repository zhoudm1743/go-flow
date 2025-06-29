<template>
        <div  style="display: flex; flex-direction: column; gap: 10px;">
            <el-card  :body-style="{'padding-bottom': '0px'}">
             <el-form :inline="true" :model="formInline" class="demo-form-inline" label-suffix=":"> 
                 <el-form-item label="流程编码">
                     <el-input v-model="queryData.flowCode" placeholder="流程编码" clearable />
                 </el-form-item>
                 <el-form-item label="流程名称">
                     <el-input v-model="queryData.flowName" placeholder="流程名称" clearable />
                 </el-form-item>
                 <el-form-item>
                     <el-button type="primary" @click="tableReload()" plain>查询</el-button>
                 </el-form-item>
             </el-form>
         </el-card>
        
         <el-card >
             <el-table v-direct-auto-height :data="tableData"  :border="true" style="width: 100%;"    v-if="tableData.length > 0">
                 <el-table-column prop="flowCode" label="流程编码"/>
                 <el-table-column prop="flowName" label="流程名称" />
                 <el-table-column prop="version" label="版本" >
                 </el-table-column>
                 <el-table-column prop="nodeName" label="当前节点" />
                 <el-table-column prop="createTime" label="流程开始时间" />
                 <el-table-column prop="initiatorName" label="流程发起人" >
                 </el-table-column>
                 <el-table-column prop="assigneeStr" label="正在处理人" >
                </el-table-column>
                 <el-table-column  label="流程状态"  width="100">
                    <template  #default="scope">
                        <el-tag  :color="getFlowColor(scope.row.flowStatus)"  effect="dark" style="border-color: unset !important;"> 

                            <span>{{ getFlowStatus(scope.row.flowStatus) }}</span>
                        </el-tag>
                    </template>
                </el-table-column>
                 <el-table-column  label="操作" width="80" >
                    <template  #default="scope">
                         <el-button link type="primary" size="small" @click="flowDetail(scope.row)">
                             详情
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
        <Flowdetail  ref="flowDetailRef"></Flowdetail>
 </template>
 
 <script setup>
 
     import { reactive, ref ,onMounted, defineExpose } from "vue"
     import { myApi, skipApi } from "~/api/flow/flowTask";
     import Flowdetail from "../common/detail.vue";
     import { getFlowStatus,getFlowColor } from "../index.js";

     const flowDetailRef = ref(null);
     const queryData = ref({
         flowCode: '',
         flowName: '',
         category: '',
     })
     const dialogVisible = ref(false);
     const tableData = ref([])
     
     const pageTable = (current, size) => {
        tableReload();
    }


    const pagination = ref({
        size: 20,
        current:   1,
        total:0,
    })

    onMounted(()=>{
        tableReload();
    })

    const flowParam = ref();


    const tableReload = () => {
        const param = {};
        Object.assign(param, queryData.value);
        param.current = pagination.value.current;
        param.size = pagination.value.size;
        myApi(param).then(res=> {
            pagination.value.total = res.total;
            tableData.value = res.rows;
        })
    }
    const flowDetail = (row) => {
        
        // flowParam.value = {
        //     id: row.id,
        //     flowJson: row.flowJson,
        //     formPath: row.formPath,
        //     formCustom: row.formCustom,
        //     formData: row.formData
        // }
        flowDetailRef.value.open( {
            id: row.id,
            flowJson: row.flowJson,
            formPath: row.formPath,
            formCustom: row.formCustom,
            formJson: row.formJson,
            formData: row.formData
        })
    }
    // 办理
     
 </script>
 
 <style  scoped>

 </style>