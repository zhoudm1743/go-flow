<template>
     <el-dialog v-model="dialogTableVisible" title="图标选择" width="800" :lock-scroll="false">
    <!-- 选择图标的组件 -->
     <div  class="icon-box"  style="height: 400px; overflow-y: scroll;">
         <ul  class="icon-list">
            <li class="icon-item"  v-for="item in list" :key="item.key">
                <span class="icon-span-1" @click="selectIcon(item.key)">
                    <el-icon style="vertical-align: middle" :size="25">
                        <component :is="item.cpt"></component>
                    </el-icon>
                    <span  class="icon-title">{{ item.key }}</span>
                </span>
            </li>
            
         </ul>
        
     </div>
    </el-dialog>
</template>

<script setup>
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import { ref } from 'vue';
const dialogTableVisible = ref(false);
const list = ref([]);

for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
   list.value.push({
     key: key,
     cpt: component
   })
}

const emits =  defineEmits(['selectIcon'])
const selectIcon = (key) => {
    emits('selectIcon', key);
    dialogTableVisible.value = false
}

const openDialog = () => {
    dialogTableVisible.value = true;
}

defineExpose({openDialog})

</script>

<style  scoped>
        .icon-list{
            padding: 0px 12px;
            overflow: hidden;
            list-style: none;
            display: grid;
            grid-template-columns: repeat(6, 1fr);
        }
        .icon-item{
            text-align: center;
            color: var(--el-text-color-regular);
            height: 60px;
            font-size: 13px;
            border: 1px solid var(--el-border-color);
            transition: background-color var(--el-transition-duration);
            margin: 4px;
        }
        .icon-span-1{
            display: flex;
            flex-direction: column;
            align-items: center;
            justify-content: center;
            height: 100%;
            cursor: pointer;
        }
     
        .icon-title{
            margin-top: 4px;
        }
</style>