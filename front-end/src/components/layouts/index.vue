<template>
     <el-container class="page-main">
      <el-header  class="page-header"><BaseHeader /></el-header>
      
      <el-container style="overflow: auto; ">
        <el-aside width="200px"><BaseSide /></el-aside>
        <el-main class="page-container">
          <!-- <div w="full"  class="content"> -->
        <!-- <AppIndex /> -->
        <!-- <el-breadcrumb :separator-icon="ArrowRight"  style="height: 35px; line-height: 30px;">
          <el-breadcrumb-item :to="{ path: '/' }"  ><span style=" font-size: 14px !important;">流程管理</span></el-breadcrumb-item>
          <el-breadcrumb-item style="font-size: 14px;">
            <span style=" font-size: 14px !important;">未实现路径</span>
            </el-breadcrumb-item>
        </el-breadcrumb> -->
       
          <bread-crumb></bread-crumb>
          <div style="height: calc(100% - 35px);">
            
              <router-view></router-view>
              
            </div>
         
        </el-main>
      </el-container>
    </el-container>
</template>

<script setup>
  import { onMounted, watch, ref } from 'vue';
  import BreadCrumb from './BreadCrumb.vue';
  import {useRoute}   from 'vue-router';
  const route = useRoute();
  onMounted(()=>{
    currentPath()
  })
  
  const pathList = ref([]);
  const currentPath = () => {
    pathList.value = route.matched.filter(item => item.meta && item.meta.title);
  }

  watch(route, (to, from) => {
    pathList.value  = to.matched.filter(item => item.meta && item.meta.title)
  }, {immediate: true})
</script>

<style  scoped>


.el-breadcrumb__inner{
  font-size: 14px !important;
}
.content{
  border: 4px solid var(--el-bg-color-page);
  box-sizing: border-box;
  overflow: hidden;
}

.page-main{
  height: 100%;
}
.page-container{
  border-left: 1px solid var(--el-border-color-light);
  box-sizing: border-box;
  padding: 8px 10px;
  /* padding: 8px 10px 10px 8px !important; */
  background-color: var(--el-border-color-lighter);
}
.page-header{
  box-shadow: 0 4px 6px #0003;
  z-index: 1;
}
.el-menu{
  border-right:none  !important;
}

</style>