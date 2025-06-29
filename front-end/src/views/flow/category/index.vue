<template>
        <el-card  :body-style="{'padding-bottom': '0px'}">
            <el-form :inline="true" :model="formInline" class="demo-form-inline" label-suffix=":"> 
                <el-form-item label="分类编码">
                    <el-input v-model="queryData.code" placeholder="分类编码" clearable />
                </el-form-item>
                <el-form-item label="分类名称">
                    <el-input v-model="queryData.name" placeholder="分类名称" clearable />
                </el-form-item>
                <el-form-item>
                    <el-button type="primary" @click="tableReload()" plain>查询</el-button>
                </el-form-item>
            </el-form>
        </el-card>
        <div style="padding: 5px 0;"></div>
        <el-card >
            <el-row class="row-bg" justify="space-between" style="margin-bottom: 16px;">
                <span>
                    <el-button type="primary" @click="openForm()">新增</el-button>
                </span>
                <span>
                
                </span>
            </el-row>
        
            <div style="display: flex; flex-wrap: wrap; gap: 8px;" v-if="tableData.length > 0">
                <el-card style="margin-bottom: 10px; min-width: 320px;  flex-basis: calc(25% - 8px)"  :span="6" v-for="item in tableData" :key="item.id">
                    <template #header>
                        <div class="card-header">
                            
                            <el-row class="row-bg" justify="space-between">
                                <span style="font-size: 18px; width: calc(100% - 100px);">{{ item.name }}</span>
                                <div style="font-size: 18px">
                                    <el-button  icon="Edit"  plain  @click="openForm(item)"/>
                                    <el-popconfirm title="确定删除当前数据吗?"
                                            @confirm="delRecord(item.id)"
                                            width="160"
                                            hide-icon
                                        >
                                            <template #reference>
                                                <el-button type="danger" icon="Delete"  plain style="margin-left: 4px;"/>
                                            </template>
                                        </el-popconfirm>
                                    
                                </div>
                            </el-row>
                            
                        </div>
                    </template>
                    <el-row :gutter="20">
                        <!-- <el-col :span="6" class="category-card-col" >
                            
                                <el-icon :size="80"  color="var(--el-color-primary)">
                                    <component :is="item.icon"></component>
                                </el-icon>
                            
                        </el-col> -->
                        <el-col :span="24" >
                            <el-descriptions  border   :column="1">
                                <el-descriptions-item label="分类编码" >{{item.code}}</el-descriptions-item>
                                <el-descriptions-item label="排序号" >{{ item.sortNum }}</el-descriptions-item>
                            </el-descriptions>
                        </el-col>
                    </el-row>
                </el-card>
            </div>

            <el-empty description="无数据"   v-else></el-empty>
            
            <el-row>
                
                <el-col :span="24">
                    <el-pagination 
                        background layout="prev, pager, next" 
                        :total="pagination.total"
                        v-model:page-size="pagination.size"
                        v-model:current-page="pagination.current"
                        @change="pageTable"
                    />
                </el-col>
            </el-row>

        
            <!-- <el-table :data="tableData"  :border="true" style="width: 100%;"   >
                <el-table-column prop="code" label="分类编码"/>
                <el-table-column prop="name" label="分类名称" />
                <el-table-column prop="sortNum" label="排序" />
            </el-table> -->


        </el-card>

    <CategoryForm ref="formRef" @tableReload="tableReload"></CategoryForm>
    
</template>

<script setup>

    import { onMounted, reactive, ref } from "vue"
    import CategoryForm from './form.vue'
    import { pageApi, delApi } from '~/api/flow/flowCategory.js'
   
    const tableData = ref([])

    const pagination = ref({
        size: 12,
        current:   1,
        total:0,
    })
    const queryData = ref({
        code: '',
        name: '',
        // sortNum: 1,
    })

    onMounted(()=>{
        tableReload();
    })
    
    const pageTable = (current, size) => {
        tableReload();
    }

    const delRecord = (id) => {
        delApi({id:id}).then(r=>{
            tableReload();
        })
    }

    const tableReload = () => {
        const params = queryData.value;
        pageApi(Object.assign(params, pagination.value)).then(res=> {
            pagination.value.total = res.total;
            tableData.value = res.rows;
            
        })
    }

  
   
  
    
    const formRef = ref(null);
    const openForm = (record) => {
        formRef.value.openForm(Object.assign({}, record));
    }
    
    defineExpose({tableReload})

</script>

<style  scoped>
    .category-card-col{
        padding-left: 0px  !important;
    }
    :deep .el-descriptions__label{
        width: 80px !important;
    }
</style>