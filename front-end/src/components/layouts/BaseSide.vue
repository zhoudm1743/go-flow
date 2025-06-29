<template>
  <el-scrollbar always>
    <el-menu
    class="el-menu-vertical-demo"
    :unique-opened="false"
    >
        <template  v-for="(item, index) in menuData" :key="item.id"  >
          <el-sub-menu :index="item.id"   v-if="item.children?.length > 0">
            <template #title>
              <el-icon>
                <component :is="item.icon"></component>
              </el-icon>
              {{item.name}}
            </template>
              <template   v-for="node in item.children" :key="node.id">
                <el-menu-item :index="node.id"   @click="handleOpen(node.path)"  v-if="node.status == 1">{{ node.name }}</el-menu-item>
              </template>
              
          </el-sub-menu>
          <el-menu-item v-else :index="item.id" :key="index"   @click="handleOpen(item.path)"> 
            <template #title>
              <el-icon>
                <component :is="item.icon"></component>
              </el-icon>
              <span>{{item.name}}</span>
            </template>  
          </el-menu-item>
        </template>
    </el-menu>
  </el-scrollbar>
</template>

<script  setup>
import { onMounted, ref, watch } from "vue";
import { usePermissionStore } from "../../store/permission";
import { storeToRefs } from "pinia";
import { useRouter } from 'vue-router';
import cloneDeep from 'lodash/cloneDeep';
const router = useRouter();

const menuData =  usePermissionStore().menuTreeList //  storeToRefs(permissionStore).menuList;
// const selectindex = ref("/home");
onMounted(() => {
  // router.push({path: '/home'})
})

const handleOpen = (key ) => {
  router.push({path:key});
  
}



const isCollapse = ref(false);

</script>
