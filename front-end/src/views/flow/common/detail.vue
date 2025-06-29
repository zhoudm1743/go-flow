<template>
       <el-drawer
        v-model="dialogVisible"
        title="流程详情"
        @close="hanldeClose"
        size="50%"
        draggable
    > 
            <el-tabs tab-position="left" v-model="tabsModel" style="min-height: 40vh; max-height: 80vh;">
                <el-tab-pane label="表单" name="form" style="padding: 10px;">
                    <component ref="dynFormRef" :data="currFormJson" :readonly="true" :variable="variable" :is="formComponent"></component>
                </el-tab-pane>
                <el-tab-pane label="流程" name="flow" style=" display: flex; flex: 1 1 0%; position: relative; overflow: auto; min-height: 60vh; max-height: 80vh;  background-color: var(--global-flow-background-color);">
                    <FlowIndex :readonly="true" :flowData="currFlowData"></FlowIndex>
                </el-tab-pane>
                <el-tab-pane label="详情" name="detail" style="padding: 10px;">
                    <el-table :data="hisList" border style="width: 100%">
                        <el-table-column prop="nodeName" label="节点名称"  />
                        <el-table-column prop="approveName" label="审批人"  />
                        <el-table-column prop="cooperateType" label="协作方式" />
                        <el-table-column prop="message" label="审批意见" width="150px"/>
                        <el-table-column  label="流程状态"  width="100">
                            <template  #default="scope">
                                <el-tag  :color="getFlowColor(parseInt(scope.row.flowStatus))"  effect="dark" style="border-color: unset !important;"> 

                                    <span>{{ getFlowStatus(parseInt(scope.row.flowStatus)) }}</span>
                                </el-tag>
                            </template>
                        </el-table-column>
                        <el-table-column prop="createTime" label="开始时间" />
                        <el-table-column prop="updateTime" label="完成时间" />
                    </el-table>
                </el-tab-pane>
            </el-tabs>
            
            <!-- <template #footer >
            <div class="dialog-footer" v-show="tabsModel == 'form'">
                <el-button @click="dialogVisible = false">Cancel</el-button>
                <el-button type="primary" @click="handlerSubmit">
                    提交
                </el-button>
            </div>
            </template> -->
        </el-drawer>

</template>

<script setup>
    import FlowIndex from '@/components/Dflow/FlowIndex.vue'
    import {hisListApi} from '~/api/flow/flowTask.js'
    import { changeNodeStatus } from '~/utils/flowUtil';
    import { onMounted, defineEmits, ref, defineExpose, watch, toRefs  } from "vue"
    import {openCustomForm} from './common.js';
    import { getFlowColor, getFlowStatus } from '../index.js';
    const flowParam = ref();
    const dialogVisible = ref(false)
    // onMounted(()=>{
    //     openApprovalForm(hisTask.value);
    // })
    // 当前的流程定义id
    const currentDefId = ref();
    const currFormJson = ref({});
    const currFlowData = ref([]);
    const variable = ref({});
    const hisList = ref([])
  
    // watch(()=> flowParam.value.formPath, (newV) => {
    //     if(newV){
    //         openApprovalForm(newV.formPath);
    //         variable.value = newV.formData;
    //         currFlowData.value = JSON.parse(newV.flowJson)
    //         // currentDefId.value = newV.id;
    //         hisListApi({instId: newV.id}).then(res => {
    //             hisList.value  = res.rows;
    //             changeNodeStatus(currFlowData.value, res.rows);
    //             console.log(currFlowData.value);
    //         })
    //     }
        
    // },{
    //     immediate: true,
    //     deep: true
    // })
      
    // 表单组件判断
    const formComponent = ref(null);
    // 动态组件ref
    const dynFormRef = ref(null);

    const tabsModel = ref("form");


    const open = (params) => {
        flowParam.value = params;
        openApprovalForm(params);
        variable.value = params.formData;
        currFlowData.value = JSON.parse(params.flowJson)
        if(params.formJson){
            currFormJson.value = JSON.parse(params.formJson)
        }
        
        // currentDefId.value = newV.id;
        hisListApi({instId: params.id}).then(res => {
            hisList.value  = res.rows;
            // changeNodeStatus(currFlowData.value, res.rows);
            console.log(currFlowData.value);
        })

        dialogVisible.value = true
    }

    const  openApprovalForm = async (def) => {
        formComponent.value =  openCustomForm(def)
    }
    const emits = defineEmits(['close'])
    
    const hanldeClose = () => {
        emits('close')
    }
    defineExpose({open})
</script>

<style  scoped>
    :deep .el-tabs__content{
        overflow: auto;
    }
</style>