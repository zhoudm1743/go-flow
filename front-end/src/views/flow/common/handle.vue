<template>
    <el-drawer
        v-model="dialogVisible"
        title="流程办理"
        @close="hanldeClose"
        size="50%"
        draggable
    >   
            
            

            <el-tabs tab-position="left" v-model="tabsModel" >
            <el-tab-pane label="表单" name="form" style="padding: 10px;">
                <div  style="display: flex; flex-direction: column; gap: 10px;">
                    <el-row justify="center"> 
                        <template v-for="item, index in buttonList" :key="index"> 
                            <el-button v-if="item.checked" type="primary" @click="handlerSubmit(item.type)">
                                {{ item.text  }}
                            </el-button>
                        </template >
                        <el-button @click="dialogVisible= false" >取消</el-button>
                    </el-row>
                    <div></div>
                    <component ref="dynFormRef"  :readonly="false" :data="currFormJson" :variable="variable" :is="formComponent"></component>
                </div>
            
            </el-tab-pane>
            <el-tab-pane label="流程" name="flow" style="position: relative;  display: flex;  overflow: visible; height: 100%; width: 100%;  background-color: var(--global-flow-background-color);">
                    <div style="flex: 1; position: relative; overflow: auto;" >
                        <FlowIndex :readonly="true" :flowData="currFlowData" style="justify-self: start;"></FlowIndex>
                    </div>
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

            <el-dialog
                v-model="approveShow"
                title="审批信息"
                width="500"
                :before-close="handleClose"
            >
                <el-form :model="approveFormData" label-width="auto">
                    <el-form-item  >
                        <el-input  type="textarea"  placeholder="审批意见"   :rows="6"  v-model="approveFormData.message" />
                    </el-form-item>
                    <el-form-item  label="选择用户" v-if="approveFormData.userSelectShow">
                        <el-button type="primary" :icon="Plus" circle @click="userSelectRef.open(flowParam.taskId)" />
                    </el-form-item>
                    <el-form-item  >
                        <div  class="flow-user-selected">
                            <el-tag  closable  @close="selectClose('user', item.id)" v-for="item in approveFormData.users" :key="item.id">
                                {{ item.nick }}
                        </el-tag>
                        </div>
                    </el-form-item>
                    
                </el-form>
                <template #footer>
                <div class="dialog-footer">
                    <el-button type="primary" @click="approveHandler">
                            确认
                    </el-button>
                    <el-button @click="approveShow = false">取消</el-button>
                </div>
                </template>
            </el-dialog>
            <UserSelect @selectData="selectData" ref="userSelectRef"></UserSelect>
        </el-drawer>
        

      
</template>

<script setup>
    import FlowIndex from '@/components/Dflow/FlowIndex.vue'
    import {hisListApi, skipApi, rejectApi, backApi, transferApi, deputeApi, addSignApi, reduSignApi} from '~/api/flow/flowTask.js'
    // import { changeNodeStatus } from '~/utils/flowUtil';
    import { onMounted, reactive, ref, defineExpose, watch, toRefs  } from "vue"
    import { openCustomForm } from './common';
    import UserSelect from './userSelect.vue';
    import { Plus } from '@element-plus/icons-vue';
    import { getFlowColor, getFlowStatus } from '../index.js';
    const modules = import.meta.glob('@/view/customform/**/**.vue')
    const flowParam = ref();
    const dialogVisible = ref(false)
    const approveShow = ref(false);
    // onMounted(()=>{
    //     openApprovalForm(hisTask.value);
    // })
    // 当前的流程定义id
    const currentDefId = ref();
    const currFormJson = ref({});
    const currFlowData = ref([]);
    const variable = ref({});
    const hisList = ref([])
    const buttonList = ref({})
    const userSelectRef = ref(null)
    const approveFormData = reactive({
        message: '',
    }) 
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
        buttonList.value = params.btnList;
        openApprovalForm(params);
        variable.value = params.formData;
        if(params.formJson){
            currFormJson.value = JSON.parse(params.formJson)
        }
        currFlowData.value = JSON.parse(params.flowJson)
        // currentDefId.value = newV.id;
        hisListApi({instId: params.id}).then(res => {
            hisList.value  = res.rows;
            // 有点问题，晚点再处理
            // changeNodeStatus(currFlowData.value, res.rows);
            // console.log(currFlowData.value);
        })

        dialogVisible.value = true
    }

    // 提交审批
    const approveHandler = () => {
        const skipType = approveFormData.skipType;
        if(skipType === "aggren"){
            approveFormData.skipType = "PASS";
            skipApi(approveFormData).then(res=> {
            
                dialogVisible.value = false;
            })
        }else if(skipType == "back"){
            approveFormData.skipType = "REJECT"
            backApi(approveFormData).then(res=> {

              
                dialogVisible.value = false;
            })
        }else if(skipType == "reject"){
            rejectApi(approveFormData).then(res=> {
            
                dialogVisible.value = false;
            })
        }else if(skipType == "transfer"){
            transferApi(approveFormData).then(res=> {
               
                dialogVisible.value = false;
            })
        }else if(skipType == "depute"){
            deputeApi(approveFormData).then(res=> {
               
                dialogVisible.value = false;
            })
        }else if(skipType == "signAdd"){
            addSignApi(approveFormData).then(res=> {
               
                dialogVisible.value = false;
            })
        }else if(skipType == "signRedu"){
            reduSignApi(approveFormData).then(res=> {
               
                dialogVisible.value = false;
            })
        }
        // transfer  depute  signAdd signRedu
    }

    //打开审批
    const handlerSubmit = (skipType) => {
        approveFormData.id = flowParam.value.taskId;
        approveFormData.skipType = skipType
        approveFormData.message = '';
        approveFormData.variable = flowParam.value.formData
        // const data = {
        //     id: flowParam.value.taskId,
        //     skipType: skipType,
        //     message: skipType,
        //     variable: 
        // }
        if(skipType  === 'transfer' || skipType  === 'depute' || skipType  === 'signAdd' || skipType  === 'signRedu' ){
            approveFormData.userSelectShow = true; 
        }else{
            approveFormData.userSelectShow = false;
            approveFormData.users = [];
        }
        approveShow.value = true;

      
        
       
    }

    const selectData = (selectData) => {
        approveFormData.users = selectData;
    }

    const  openApprovalForm = async (def) => {
        const component =  openCustomForm(def)
        formComponent.value = component;
        
    } 
    
  
    defineExpose({open})
</script>

<style  scoped>
        .test-class{
            display: flex; flex-direction: column; height: 100%; width: 100%; 
        }
        .flow-user-selected{
            display: flex;
            align-content: center;
            align-items: center;
            flex-wrap: wrap;
            gap: 10px;
            width: 100%;
        
        }
</style>