<template>
   <el-breadcrumb :separator-icon="ArrowRight"  style="height: 35px; line-height: 30px;">
    <el-breadcrumb-item  v-for="(item,index) in pathList" :key="item.id" >
        <span :style=" {fontSize: '14px !important', fontWeight: index == 0 ? 'bold' : 'normal'}" >
            {{ item.meta.title }}
        </span>
    </el-breadcrumb-item>
    
</el-breadcrumb>
</template>

<script setup>
import {ref, onMounted, watch} from 'vue'
import { useRoute } from 'vue-router'
const route = useRoute()
const pathList = ref()

const getCurrentPath = () => {
    console.log(route.matched);
    pathList.value = route.matched.filter(item => item.meta && item.meta.title);
}
 
onMounted(() => {
    getCurrentPath();
})
 
watch(route, (to, from) => {

    pathList.value = to.matched.filter(item => item.meta && item.meta.title);
}, { immediate: true });
                        
</script>

<style lang="scss" scoped>

</style>