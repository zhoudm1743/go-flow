<template>
    <div  style="display: flex; flex-direction: column; gap: 10px;  padding-bottom: 15px;">
        <el-card  :body-style="{'padding-bottom': '0px'}">
         <el-form :inline="true" :model="formInline" class="demo-form-inline" label-suffix=":"> 
             <el-form-item label="菜单名称" >
                 <el-input v-model="queryData.name" placeholder="菜单名称" clearable />
             </el-form-item>
             <el-form-item label="是否可见">
                 <el-select v-model="queryData.status" placeholder="" clearable  >
                     <el-option :label="item.label" :value="item.value" v-for="(item,index) in getStatus()" :key="index" /> 
                 </el-select>
                 <!-- <el-input v-model="queryData.roleName" placeholder="角色名称" clearable /> -->
             </el-form-item>
             <el-form-item>
                 <el-button type="primary" @click="tableReload()" plain>查询</el-button>
             </el-form-item>
         </el-form>
     </el-card>
    
     <el-card >
        <el-row class="row-bg" justify="space-between" style="margin-bottom: 16px;">
            <span>
                <el-button type="primary" @click="openForm()" >新增</el-button>
            </span>
            <span>
            
            </span>
        </el-row>
         <el-table :data="tableData"  :border="true" style="width: 100%;"  default-expand-all  row-key="id" v-if="tableData?.length > 0">
             <el-table-column prop="name" label="菜单名称"/>
             <el-table-column prop="icon" label="图标" width="60">
                 <template #default="scope">
                     <el-icon :size="25" >
                     <component :is="scope.row.icon"/>
                     </el-icon>
                 </template>
             </el-table-column>
             <el-table-column prop="type" label="类别" width="80">
                 <template #default="scope">
                     {{ getTypeLabel(scope.row.type) }}
                 </template>
             </el-table-column>
             <el-table-column prop="path" label="路由地址" />
             <el-table-column prop="component" label="组件" />
             <el-table-column prop="status" label="是否可见"  width="90">
                 <template #default="scope">
                     <el-tag  :type="scope.row.status === '1' ? 'primary' : 'info'" >{{ getStatusLabel(scope.row.status) }}</el-tag>
                 </template>
             </el-table-column>
             <el-table-column prop="sortNum" label="排序" width="60" />

             <el-table-column  label="操作"  width="180">
                 <template  #default="scope">
                     <el-button link type="primary" size="small" @click="openForm(scope.row)">
                         编辑
                         </el-button>
                         <el-button  link type="primary" size="small" @click="handleDelete(scope.row.id)">
                         删除
                         </el-button>
                         <el-button  link type="primary" size="small" @click="handleIdentifier(scope.row.id)">
                            权限标识
                         </el-button>
                 </template>
             </el-table-column>
         </el-table>

         <el-empty description="无数据"   v-else></el-empty>
   
         <div style="padding: 5px 0;"></div>
         <!-- <el-row>
             
             <el-col :span="24">
                 <el-pagination 
                     background layout="prev, pager, next" 
                     :total="pagination.total"
                     v-model:page-size="pagination.size"
                     v-model:current-page="pagination.current"
                     @change="pageTable"
                 />
             </el-col>
         </el-row> -->
     </el-card>
 
</div>
    <MenuForm ref="menuFormRef"  :menuTree="menuTree"  @tableReload="tableReload"></MenuForm>
    <Identifier ref="identifierRef"></Identifier>
</template>

<script setup>
 import { reactive, ref ,onMounted, defineExpose, watch } from "vue"

 import MenuForm from './form.vue'

 import { treeApi, addApi, editApi,deleteApi } from "~/api/sys/menu";
 import { getTypeLabel, getStatusLabel, getStatus } from './menu';
 import Identifier from "./identifier/index.vue";
 
 const identifierRef = ref(null);
 const queryData = ref({
     roleCode: '',
     roleName: '',
     roleStatus: '',
 })
 const tableData = ref([])
 
 const menuTree = ref([{
        id: "0",
        name: '根节点',
    }]);

 const pagination = ref({
    size: 20,
    current:   1,
    total:0,
})

const roleUserComponent = ref();

onMounted(()=>{
    tableReload();
})

watch(tableData, (newV, oldV) => {
    menuTree.value = [{
        id: "0",
        name: '根节点',
        children: newV
    }];
});

const tableReload = () => {
    const param = {};
    Object.assign(param, queryData.value);
    param.current = pagination.value.current;
    param.size = pagination.value.size;
    treeApi(param).then(res=> {
        tableData.value = res.rows;
    })
}

const menuFormRef = ref(null);


const openForm = (record) => {
    menuFormRef.value.openForm(Object.assign({}, record));
}


const handleDelete = (id) => {
    deleteApi([{"id":id}]).then(res => {
        tableReload();
    })
}


const handleIdentifier = (id) => {
    identifierRef.value.openDrawer(id);
}

defineExpose({tableReload})



</script>

<style lang="scss" scoped>

</style>