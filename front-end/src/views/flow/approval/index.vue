<template>
    
    <el-card >
        <el-tabs type="border-card" v-if="tableData.length > 0">
            <el-tab-pane  v-for="item in tableData" :key="item.flowCategory.id" :label="item.flowCategory.name">
                
                <div  class="flow-list"  >
                    <el-card class="flow-item"  body-class="flow-item-body" @click="openApprovalForm(def)" v-for="def in item.definitionList" :key="def.id" >
                    
                                <el-icon :size="40" >
                                <component :is="def.icon"/>
                                </el-icon>
                            <span  style="font-size: 12px; font-weight: 500;">
                                {{def.flowName}}
                            </span>
                    
                    </el-card>
                </div>
            </el-tab-pane>
            
        </el-tabs>
    
        

        <el-empty description="无数据" v-else ></el-empty>
        
        <el-dialog
            v-model="dialogVisible"
            title="发起流程"
            draggable
            :before-close="handleClose"
        >   

            <el-tabs tab-position="left" v-model="tabsModel" style="min-height: 40vh; max-height: 80vh;">
                <el-tab-pane label="表单" name="form" style="padding: 10px;">
                    <component ref="dynFormRef"  :readonly="false" :data="currFormJson"  :is="formComponent"></component>
                  
                </el-tab-pane>
                <el-tab-pane label="流程" name="flow" style=" display: flex; flex: 1 1 0%; position: relative; overflow: auto; min-height: 60vh; max-height: 80vh;  background-color: var(--global-flow-background-color);">
                    <FlowIndex :readonly="true" :flowData="currFlowData"></FlowIndex>
                </el-tab-pane>
             
            </el-tabs>
            
           
            <template #footer >
            <div class="dialog-footer" v-show="tabsModel == 'form'" >
                <el-button @click="dialogVisible = false">取消</el-button>

                <el-button  v-if="userPositions && userPositions.length > 1"  ref="submitRef" type="primary" @click="popoverVisible = true"  v-click-outside="onClickOutside" >
                    提交
                </el-button>
                <el-button  v-else   type="primary"  @click="handlerSubmit">
                    提交
                </el-button>
            </div>
            </template>
        </el-dialog>

    
        <!-- <el-table :data="tableData"  :border="true" style="width: 100%;"   >
            <el-table-column prop="code" label="分类编码"/>
            <el-table-column prop="name" label="分类名称" />
            <el-table-column prop="sortNum" label="排序" />
        </el-table> -->


    </el-card>


    <el-popover
        width="min-content"
        ref="popoverRef"
        :virtual-ref="submitRef"
        trigger="click"
        title="请选择身份"
        virtual-triggering
        placement="top-end"
        :visible="popoverVisible"
        
    >
      <el-row :gutter="10" >
        <el-col :span="24" v-for="item in userPositions" :key="item.id" class="position-item"  >
           <el-button @click="poshandlerSubmit(item)"  class="position-item-button">
                <el-tag  type="danger"  effect="dark">{{ item.cmpName }}</el-tag> |
                <template v-if="item.orgName" >
                    <el-tag  type="primary"   effect="dark">{{ item.orgName }}</el-tag>   |
                </template>
                <el-tag  type="success"  effect="dark">{{item.positionName}}</el-tag>
            </el-button>
        </el-col>
      </el-row>
    </el-popover>
</template>

<script setup>
import { ClickOutside as vClickOutside } from 'element-plus'
    import FlowIndex from '@/components/Dflow/FlowIndex.vue'
    import { onMounted,  ref, defineExpose, unref  } from "vue"
    import {useUserStore} from "~/store/user";
    import { listApi } from "~/api/flow/flowDef";
    import { startApi } from "~/api/flow/flowInstance";
    import { openCustomForm } from '../common/common';
    import { EBuilder } from "epic-designer";

    const userStore = useUserStore();
    const userPositions = userStore.userPositions;
    const tableData = ref([])
    // dialog 显示 隐藏
    const dialogVisible = ref(false);
    // 当前的流程定义id
    const currentDefId = ref();
    // 表单组件判断
    const formComponent = ref(null);
    // 动态组件ref
    const dynFormRef = ref(null);

    const tabsModel = ref("form");

    const currFormJson = ref({});
    const currFlowData = ref([]);
    
    // 多身份处理
    const submitRef = ref(null);
    const popoverRef = ref(null);
    const popoverVisible = ref(false);
    // 初始化加载数据
    onMounted(()=>{
        tableReload();
    })

    const onClickOutside = () => {
        popoverVisible.value = false;
        
     
    }
    const handlerSubmit = () => {
        if(userPositions && userPositions.length == 1){
            poshandlerSubmit(userPositions[0]);
        }else{
            // 无职位信息
            poshandlerSubmit({})
        }
       ;
    }


    // 提交按钮操作  发起申请携带职位部门等信息
    const poshandlerSubmit = (positionInfo) => {
     
        // 表单验证后
        dynFormRef.value.submitForm().then(res=> {
            const data = {
                defId: currentDefId.value,
                positionInfo: positionInfo,
                formData: res
            }
            // 发起申请
            startApi(data).then(res=> {
                console.log(res)
                dialogVisible.value = false;
            })
            
         })
        
        
    }

    
    // 请求表格数据
    const tableReload = () => {
    
        listApi({}).then(res=> {
            tableData.value = res.rows;
            
        })
    }

   

    // 打开表单
    const  openApprovalForm = async (def) => {
        formComponent.value =  openCustomForm(def)
        dialogVisible.value = true;
        currFlowData.value = JSON.parse(def.flowJson)
        currFormJson.value = JSON.parse(def.formJson)
        currentDefId.value = def.id;
       
    } 

    defineExpose({tableReload})

</script>

<style  scoped>
    .category-card-col{
        padding-left: 0px  !important;
    }
    :deep .el-descriptions__label{
        width: 80px !important;
    }

    .flow-list{
        display: flex; 
        flex-wrap: wrap; 
        gap: 100px;
    }
    .flow-item{
        margin-bottom: 10px; 
        min-width:  100px;  
        max-width: 180px;
        flex-grow: 1;
       
    }
    .flow-item:hover{
        cursor: pointer;
        transform: scale(1.05)
    }
    :deep .flow-item-body{
        height: 100%;
        display: flex;
        flex-direction: column;
        flex-wrap: wrap;
        align-content: center;
        align-items: center;
        gap: 10px;
    }
    .position-item{
        margin-bottom: 5px;
    }
    .position-item-button{
        border: 1px solid #dcdfe6; 
        padding: 18px 4px;
        width: 100%;
        text-align: left !important;
    }
    .position-item-button:hover{
        background-color: var(--el-color-primary);

    }
</style>