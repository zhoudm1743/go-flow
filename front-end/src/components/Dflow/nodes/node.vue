
<template>
    <div  class="flw-node-single">
        <Panel :class="panelClass + ' ' + statusClass" v-bind="$attrs" >
          
        </Panel>
        <Line v-bind="$attrs" ></Line>
    </div>
</template>

<script setup>
import Panel from './panel.vue'
import Line from  './line.vue'
import {defineProps,  watch, ref, useAttrs} from  'vue'
    const props =  defineProps({
        
        panelClass:{
             type: String,
            default: ''
        },
        //节点状态
        nodeStatus:{
            type: Array,
        },
       
       
       
    })
    const statusClass = ref("");


    watch(()=>props.nodeStatus,(newv,oldv)=>{

        if(newv ){
            if(newv.includes("1")){
                statusClass.value = "flow-status-progress";
            }else if(newv.includes("2")){
                statusClass.value = "flow-status-completed";
            }else if(newv.includes("9")){
                statusClass.value = "flow-status-back"
            }else if(newv.includes("99")){
                statusClass.value = "flow-status-reject";
            }else{
                statusClass.value = "";
            }

        }
    },{immediate:true, deep:true})
</script>

<style  scoped>
    .flow-status-back{
        border: 2px solid RGB(255, 105, 164);
        padding: 1px;
    }
    .flow-status-reject {
        border: 2px solid #f25643;
        padding: 1px;
    }
    .flow-status-completed{
        border: 2px solid #67B26F;
        padding: 1px;
    }
    .flow-status-progress{
        border: 2px solid #315EFB;
        padding: 1px;
    }
    .flw-node-single{
        padding-left: var(--golbal-flow-node-padding);
        padding-right: var(--golbal-flow-node-padding);
        z-index: 1;
    }
    
   
    .flw-node-left-top{
        position: absolute;
        width: 50%;
        height: 1px;
    }
</style>