<template>
    
    <div  style="display: flex; flex-direction: column; height: 100%; width: 100%;  ">
        <el-card  style=" min-height: 65px;">
            <el-row  justify="space-between">
                <el-row justify="space-between" style="width: 300px;" align="middle">
                    <el-button @click="$router.go(-1)">返回</el-button>
                    <span style="width: 180px; overflow: hidden;" :title="flowData.flowName">流程名称：{{ flowData.flowName  }}</span>
                </el-row>

                <el-row  justify="space-between" style="width: 350px;" align="middle" class="flow-config-step-num">
                    <template v-for="(item,index) in configSteps" :key="index">
                        <span  @click="changeStepNum(item.value)" 
                        :class="currStep == item.value ? 'step-num-active' : '' + noSelectClass(item.value)">  
                            <span class="step-num">{{index + 1}}</span> {{item.label}}
                        </span>
                    </template>
                </el-row>
                <div style="width: 300px;" >
                    <el-button style="margin-left: 8px; float: right;" @click="saveConfig">保存</el-button>
                </div>

            </el-row>
        </el-card>
        <div style="flex: 1; position: relative; overflow: auto;"  v-show="flowData.formCustom == 'N' &&  currStep === 1">
            <EDesigner ref="edDesignerRef" :hiddenHeader="true" :formMode="true" :sourceCodeReadOnly="true"/>
        </div>
       <div style="flex: 1; position: relative; overflow: auto;" v-show="currStep === 2">
            <FlowDesign :formColumns="formColumns" :data="flowData" :readonly="readonly"  ></FlowDesign> 
       </div>
       <div style="flex: 1; position: relative; overflow: auto;" v-show="currStep === 3">
            <GlobalConfig></GlobalConfig>
        </div>
        <!-- <el-drawer :key="currNodeData?.nodeId"
            v-model="drawer"
            @open="openDrawer"
            @close="closeHandler"
            :direction="direction"
            destroy-on-close
        >   
        <template #header>
            <el-text v-if="!drawerEdit" @click="changeTitle" tag="ins">{{ currNodeData?.nodeName }}</el-text>
            <el-input v-else  ref="drawerEditInput" v-model="currNodeData.nodeName" autofocus  @blur="drawerEditBlur" />
        </template>
           <template #default>
            <component ref="settingRef" :key="currNodeData.nodeId" :data="currNodeData"  :backNodeList="backNodeList" :is="currComponent"></component>
           </template>
          
        </el-drawer> -->
    </div>

      
 
</template>

<script setup>
    import { nextTick, onMounted, reactive, ref,watch } from 'vue'
    import { EDesigner } from "epic-designer";
    import UserConfig from '@/components/Dflow/setting/userConfig.vue'
    import StartConfig from '@/components/Dflow/setting/startConfig.vue'
    import SerialConfig from '@/components/Dflow/setting/serialConfig.vue'
    import FlowDesign from '@/views/dflow-design/index.vue'
    import GlobalConfig from './globalConfig.vue' 
    import { uuid, nodeList,parseFormJson } from '~/utils/flowUtil';
    import { detailApi, saveConfigApi } from '~/api/flow/flowDesign';
    import { useRouter } from 'vue-router';

    const flowData = ref({})
    const settingRef = ref(null)
    const router = useRouter();
    const id = ref(null);
    const currStep = ref(1)
    const edDesignerRef = ref(null);

    const configSteps = ref([
        {
            label: '表单配置',
            value: 1,
        },
        {
            label: '流程配置',
            value: 2,
        },
        {
            label: '全局配置',
            value: 3,
        }
    ])
    const formColumns = ref([]);

    onMounted(() => {
        formColumns.value = [];
        id.value = router.currentRoute.value.params.id;
        detailApi({id:id.value}).then((res)=>{
            flowData.value = res.data;
            if(res.data.formCustom == 'Y'){
                currStep.value = 2
            }else{
                if(flowData.value?.formJson != null){
                    edDesignerRef.value.setData(JSON.parse(flowData.value.formJson));
                }
            }
         
          
            if(flowData.value?.flowJson == null){
                flowData.value.flowJson = 
                        { 
                            nodeId: uuid(),
                            nodeType: "start",
                            nodeName:"开始节点",
                            value: "所有人",
                            properties:{},
                            childNode:{
                                nodeId: uuid(),
                                nodeType: 'between',
                                nodeName:"审批人",
                                value: '',
                                properties:{},
                                    childNode:{
                                    nodeId: uuid(),
                                    nodeType: 'end',
                                    nodeName:"结束",
                                    properties:{},
                                }
                            }
                        
                        }
                       
                    
                    
            }else{
                flowData.value.flowJson = JSON.parse(flowData.value.flowJson);
            }
        })
    })

    const changeStepNum = (num) => {
        if(!(flowData.value.formCustom == 'Y' && num == 1)){
            currStep.value = num
        }
    }

    watch(currStep, (newV, oldV) => {
        //  点击非表单配置时， 获取表单json
        if(newV != 1){
          const pageSchema =  edDesignerRef.value.getData();
          if(pageSchema.schemas){
            formColumns.value = parseFormJson(pageSchema?.schemas[0].children);
          }
        
          console.log(formColumns.value);
        }
    })

    const noSelectClass = (num) =>{
        if(flowData.value.formCustom == 'Y' && num == 1){
            return "no-select";
        }
        return '';
    }
    
    // 发布
    const saveConfig = () => {
        const formData = edDesignerRef.value.getData();
        const data = {
            id: id.value , 
            flowJson: JSON.stringify(flowData.value.flowJson),
            formJson: JSON.stringify(formData)
        }
        saveConfigApi(data).then(result => {
            console.log(result)
        })
    }
    

    const drawer = ref(false)
    const drawerEdit = ref(false);
    const drawerEditInput = ref();

    const changeTitle  = () =>{
        // 不这么写 多点几次后 input 聚焦无效
        drawerEdit.value = true;
        nextTick(()=>{
            drawerEditInput.value.focus();
        })
      
    }
    const drawerEditBlur = () => {
        if(currNodeData.value.nodeName.length == 0){
            currNodeData.value.nodeName = '未设置';
        }
        drawerEdit.value = false;   
    }

    const readonly = ref()

    const direction = ref('rtl')
   
   

    const openDrawer = (e, e1, e2) =>{
        console.log(e)
    }

    const nodeType = ref()
    const currNodeData = ref();

    const currComponent = ref();
   
    const backNodeList = ref([])

    const nodeClick = (nodeId, nodeData) => {
        const nodelistResult = [];
        nodeList(flowData.value, nodelistResult);
        backNodeList.value = nodelistResult;
        nodeType.value = nodeData.nodeType;
        currNodeData.value= nodeData;


        if(nodeData.default != null && nodeData.default){
            //跳过
            return;
        }
        else if(nodeData.nodeType == 'between' || nodeData.nodeType == 'parallel-node'){
            currComponent.value = UserConfig
        }else  if(nodeData.nodeType == 'start'){
            currComponent.value = StartConfig
        }else  if(nodeData.nodeType == 'serial-node'){
            currComponent.value = SerialConfig
        }else {
            // currComponent.value = Testc;
        }

        drawer.value = true;
    }


    const closeHandler = () => {
       settingRef.value.formConfig().then(formValue => {
            currNodeData.value.config = formValue;
            currNodeData.value.value = formValue.value;
       }).catch(error => {
           console.log(error)
       });
   
     
    }
</script>

<style  scoped>
 
    .ep-main{
        --el-main-padding: 0px;
        border-right:1px solid var(--el-border-color);
        box-sizing: border-box;
    }
    .flow-config-main{
        height: 100%;
    }
    .main-container{
        display: flex;
        width: 100%;
        height: 100%;

    }
    .left{
        width: 100px;
    }
    .left-content{
        height: 100%;
        width: calc(100% - 500px);
    }
    .right-content{
        width: 500px;
        border: 1px solid green;
    }

    .read-Model{
        position: fixed;
        top: 177px;
        right: 110px;
    }

    .flow-config-step-num > span{
        padding: 2px 10px;
        cursor: pointer;
        border: 1px solid #b5c5e7;
        background-color: aliceblue;
        border-radius: 12px;
        
    }
    .flow-config-step-num .step-num{
        display: inline-flex;
        justify-content: center;
        align-items: center;
        width: 18px;
        height: 18px;
        border: 1px solid var(--el-text-color-primary);
        border-radius: 50%;
    }
    .flow-config-step-num >span:hover{
        cursor: pointer;
        background-color: #2049fe;
        color: #ffffff;
    }
    .no-select{
        cursor: no-drop !important;
    }
    .no-select:hover{
        cursor: no-drop !important;
    }

    .flow-config-step-num > .step-num-active{
        background-color: #2049fe;
        color: #ffffff;
    }
</style>