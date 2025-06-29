<template>
  <el-config-provider >
    <el-watermark :font="font" content="openflow" style="height: 100vh;">
    <router-view></router-view>
    </el-watermark>
  </el-config-provider>
  <!-- <DflowSign></DflowSign> -->
</template>

<script setup>
import { watch, reactive } from "vue";
import DflowSign from './views/dflow-design/index.vue'
import { toggleDark,isDark } from "~/composables/index.js";
const font = reactive({
  color: 'rgba(0, 0, 0, .15)',
})
watch(
  isDark,
  () => {
    font.color = isDark.value
      ? 'rgba(255, 255, 255, .15)'
      : 'rgba(0, 0, 0, .15)'
  },
  {
    immediate: true,
  }
)
</script>
<style  scoped>
#app {
  /* text-align: center; */
  /* background-color: var(--el-bg-color-page); */
  overflow: hidden;
 
  /* color: var(--el-text-color-primary); */
}


</style>
