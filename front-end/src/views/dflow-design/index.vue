<template>
        
        <FlowIndex :flowData="data.flowJson" :readonly="readonly"  @nodeClick="nodeClick"></FlowIndex> 
        <setting ref="settingRef"  ></setting>
</template>

<script setup>
    
    import { defineProps, ref,toRefs,watch } from 'vue'
    import FlowIndex from '@/components/Dflow/FlowIndex.vue'
    import setting  from '@/components/Dflow/setting/index.vue';
    import { backNodeList } from './flow';
    const props = defineProps({
        data: Object,
        readonly: {
            type: Boolean,
            default: false
        },
        formColumns: {
            type: Array,
            default: () => []
        }
        
    })
    const settingRef = ref();
    const currNodeId = ref(null);
    const currNode = ref();
    const {data, readonly} = toRefs(props);
    const backNodes = ref([])
    
    watch(()=> data.value.flowJson, (newV) => {
        if(newV){
            backNodes.value = backNodeList(newV);
        }
    },{
        immediate: true,
        deep: true
    })
    //节点点击事件
    const nodeClick = (node) => {
        node.properties = node.properties || {};
        node.properties.formColumns = props.formColumns;
        if(readonly.value){
            return;
        }
        currNodeId.value = node.nodeId;  
        currNode.value =  node; // getNode(flowData.value, nodeId);
        //打开配置项
        settingRef.value.openHandler(node, backNodes.value);
    }
</script>

<style lang="scss" scoped>

</style>