<template>
    <el-button @click="handerData">获取数据</el-button>
    <div class="epic-designer-container">
      <EDesigner :hiddenHeader="true" :formMode="true" :sourceCodeReadOnly="true" ref="designerRef"/>
    </div>
  </template>
  <script setup >
  import { ref } from "vue";
  import { EDesigner } from "epic-designer";
  const designerRef = ref(null);

  const handerData = () => {
    const data = designerRef.value.getData();
    console.log(data);
  }
  </script>
  <style>
 .epic-designer-container {
  width: 100%;
  height: 100%;
  border: 1px solid var(--epic-border-color);
}
  </style>