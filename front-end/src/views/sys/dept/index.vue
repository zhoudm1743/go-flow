<template>
    <div  style="display: flex; flex-direction: column; gap: 10px;">
        <el-card  :body-style="{'padding-bottom': '0px'}">
         <el-form :inline="true" :model="formInline" class="demo-form-inline" label-suffix=":"> 
             <el-form-item label="名称" >
                 <el-input v-model="queryData.name" placeholder="名称" clearable  style="width: 200px"/>
             </el-form-item>
             <el-form-item label="类别">
                 <el-select v-model="queryData.levelCode" placeholder="" clearable  style="width: 200px;">
                     <el-option :label="item.label" :value="item.value" v-for="(item,index) in getLevelCodes()" :key="index" /> 
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
         <el-table :data="tableData"  :border="true" style="width: 100%;" default-expand-all   row-key="id" v-if="tableData?.length > 0">
             <el-table-column prop="name" label="名称"/>
             <el-table-column prop="type" label="类型" >
                 <template #default="scope">
                     {{ getLevelCodeLabel(scope.row.levelCode) }}
                 </template>
             </el-table-column>
             <el-table-column prop="directorName" label="负责人" />
             <el-table-column prop="sortNum" label="排序" />
             <el-table-column  label="操作" >
                 <template  #default="scope">
                    <el-button link type="primary" size="small" @click="openForm(scope.row)">
                         编辑 
                    </el-button>
                    <el-button link type="primary" size="small" @click="addChildren(scope.row)">
                         添加下级 
                    </el-button>
                    <el-button  link type="primary" size="small" @click="handleDelete(scope.row.id)" v-if="scope.row.parentId != 0">
                        删除
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
    <DeptForm ref="formRef"  :menuTree="menuTree" @tableReload="tableReload"></DeptForm>
   
</template>

<script setup>
 import { reactive, ref ,onMounted, defineExpose, watch } from "vue"

 import DeptForm from './form.vue'

 import { treeApi, addApi, editApi,deleteApi } from "~/api/sys/dept";
 import { getLevelCodeLabel, getLevelCodes } from './dept';
 const queryData = ref({
     name: '',
     levelCode: '',
 })
 const tableData = ref([])
 
 const menuTree = ref([]);

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
    menuTree.value = newV;
});

const tableReload = () => {
    const param = {};
    Object.assign(param, queryData.value);
    param.current = pagination.value.current;
    param.size = pagination.value.size;
    treeApi(param).then(res=> {
        pagination.value.total = res.total;
        tableData.value = res.rows;
    })
}

const formRef = ref(null);


const addChildren = (record) => {
    const params = {
        pid: record.id
    };
    formRef.value.addChildren(params);
}

const openForm = (record) => {
    formRef.value.openForm(Object.assign({}, record));
}


const handleDelete = (id) => {
    deleteApi([{"id":id}]).then(res => {
        tableReload();
    })
}

defineExpose({tableReload})

</script>

<style lang="scss" scoped>

</style>