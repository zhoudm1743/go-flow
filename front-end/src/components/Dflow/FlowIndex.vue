<template>

    <!-- <div  class="my-flow" id="my-flow"> -->
    <div style="flex: 1; position: relative;   overflow: auto; background-color: var(--global-flow-background-color)" >
            
            <Maxmin :percent="percent"  @updatePercent="updatePercent"></Maxmin>
         
        <div  class="flow-container" id="flow-container"  :style="{'transform': 'scale('+(percent/100)+')'}"
            @mousemove="mousemove"
        >
            <FlowRender 
                v-model="flowData"  
                v-bind="flowData"
                @removeNodes="removeNodes" 
                @addNodes="addNodes"
                @nodeClick="$attrs.onNodeClick"
                :readonly="readonly" ></FlowRender>
        </div>
    </div>
        


    
</template>

<script setup>
    import { ref ,defineProps, toRefs, watch} from 'vue';
    import FlowRender from './nodes/render.vue'

    import Maxmin from './maxmin.vue'
    import { uuid } from '../../utils/tool';

   const props = defineProps({
        flowData:{},
        readonly: {
            type: Boolean,
            default: false
        },

    })
    const {readonly,flowData } = toRefs(props)


   


    // 百分比
    const percent = ref(100)
    
 

    const updatePercent = (num) => {
        const value = percent.value + num;
        if(value<= 200 && value >= 10){
            percent.value = value;
        }
    }
  
    // 移除节点
    const removeNodes = (nodeId) => {
        removeChild(flowData.value, nodeId)
    }

    function removeChild(parentNode,nodeId){
      
        var node = parentNode.childNode;
        

        if(node){
            if(node.nodeId == nodeId){
                parentNode.childNode = node.childNode
            }else if(node.nodeType == 'serial' || node.nodeType == 'parallel'){
                // 网关节点，判断conditionNodes的 节点
                var isFind = false;
                for(var i in node.conditionNodes){
                    if(node.conditionNodes[i].nodeId == nodeId){
                        isFind = true;
                        node.conditionNodes.splice(i,1);
                        if(node.conditionNodes.length <= 1){
                            parentNode.childNode = node.childNode
                        }
                    }
                }
                if(!isFind){
                    for(var i in node.conditionNodes){
                        removeChild(node.conditionNodes[i],nodeId);
                    }
                    removeChild(node,nodeId)
                }

            }
            else{
                removeChild(node,nodeId)
            }
           
        }
    }
  

    // 增加节点
    const addNodes = (params) => {
        var temChildNode = params.node.childNode;
        if(params.nodeType == 'between'){
            params.node.childNode= {
                                nodeId: uuid(),
                                nodeName:'审批人1',
                                nodeType: 'between',
                                value: '',
                                type: 'between',
                                properties:{},
                                childNode:temChildNode

            };
        }else if(params.nodeType == 'parallel'){
            params.node.childNode ={
                nodeId: uuid(),
                nodeType: 'parallel',
                nodeName:"并行分支",
                properties:{},
                conditionNodes:[
                        {
                            nodeId: uuid(),
                            nodeName:'审批人1',
                            nodeType: 'between',
                            type: 'between',
                            properties:{},
                            value: '',
                        },
                    
                        {
                            nodeId: uuid(),
                            nodeName:'审批人2',
                            nodeType: 'between',
                            type: 'between',
                            properties:{},
                            value: '',
                            properties:{},
                        },
                    
                ],
                childNode:temChildNode
            }
        }else if(params.nodeType == 'serial'){
            params.node.childNode =  {
                nodeId: uuid(),
                nodeType: 'serial',
                nodeName:"分支选择",
                conditionNodes:[
                
                        {
                        nodeId: uuid(),
                        nodeType: 'serial-node',
                        type: 'serial-node',
                        nodeName:"条件1",
                        sortNum: 0,
                        value: '',
                        properties:{},
                        },
                        {
                        nodeId:  uuid(),
                        nodeType: 'serial-node',
                        type: 'serial-node',
                        nodeName:"条件2",
                        value: '其他条件走此流程',
                        sortNum: 1,
                        default: true,
                        properties:{},
                        },

                ],
                childNode:temChildNode
                
            }
        }
        // addNode(flowData.value, params)
    }




    


    
   
</script>

<style  scoped>
    span{
        font-size: 12px;
    }
    div{
        font-size: 12px;
    }
    .my-flow{
        display: flex;
        width: 100%;
        height: 100%;
        flex-direction: column;
        flex-wrap: nowrap;
        background-color: var(--global-flow-background-color);
    }
    #my-flow{
        height: 100%;
    }
    /* .my-flow{
        width: 100%;
        height: 100%;
        overflow: auto;
        position: relative;
       
        background-color: var(--global-flow-background-color);
    } */
    .flow-container{
        position: relative;
        width: 100%;
        /* height: 100%; */
        padding-top: 10px;
        align-items: flex-start;
        justify-content: center;
        flex-wrap: wrap;
        min-width: -moz-min-content;
        min-width: min-content;
        transform-origin: 50% 0px 0px;
        background-color: var(--global-flow-background-color);
    }
    
    /* 滚动条整体样式 */
    .my-flow::-webkit-scrollbar {
        width: 12px; /* 适用于垂直滚动条 */
        height: 12px; /* 适用于水平滚动条 */
        }

        /* 滚动条的滑块部分 */
    .my-flow::-webkit-scrollbar-thumb {
        background-color: #c2c2c2; /* 滑块颜色 */
        border-radius: 6px; /* 滑块圆角 */
        border: 2px solid transparent; /* 滑块边框，可以设置为透明或颜色 */
        }

        /* 滚动条的轨道部分 */
    .my-flow::-webkit-scrollbar-track {
        background-color: #f0f0f0; /* 轨道颜色 */
        border-radius: 6px; /* 轨道圆角 */
        }

        /* 鼠标悬停在滚动条上时的样式 */
    .my-flow::-webkit-scrollbar-thumb:hover {
        background-color: #a6a6a6; /* 鼠标悬停时滑块颜色加深 */
        }

        /* 滚动条滑块被滚动时的样式 */
    .my-flow::-webkit-scrollbar-thumb:active {
        background-color: #8c8c8c; /* 滚动时滑块颜色更深 */
        }
</style>