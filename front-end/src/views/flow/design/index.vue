<template>
   <el-card  :body-style="{'padding-bottom': '0px'}">
            <el-form :inline="true" :model="formInline" class="demo-form-inline" label-suffix=":"> 
                <el-form-item label="分类编码">
                    <el-input v-model="queryData.flowCode" placeholder="分类编码" clearable />
                </el-form-item>
                <el-form-item label="分类名称">
                    <el-input v-model="queryData.flowName" placeholder="分类名称" clearable  />
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
                    <el-button type="primary" @click="openForm()" >新增</el-button>
                </span>
                <span>
                
                </span>
            </el-row>
            
            <el-table v-direct-auto-height :data="tableData"  :border="true" style="width: 100%;"    v-if="tableData.length > 0">
                <el-table-column prop="flowCode" label="流程编码"/>
                <el-table-column prop="flowName" label="流程名称" />
                <el-table-column prop="icon" label="图标" >
                    <template #default="scope">
                        <el-icon :size="25" >
                        <component :is="scope.row.icon"/>
                        </el-icon>
                    </template>
                </el-table-column>
                <el-table-column prop="flowVersion" label="流程版本" />
                <el-table-column prop="categoryName" label="业务类别" />
                <el-table-column prop="formCustom" label="表单类型" >
                    <template #default="scope">
                        {{scope.row.formCustom === 'Y' ? "开发表单":"设计表单" }}
                    </template>
                </el-table-column>
                <el-table-column prop="status" label="状态"  width="80">
                    <template #default="scope">
                            {{ scope.row.status == '1' ? "正常":"停用" }}
                    </template>
                </el-table-column>
                
                <el-table-column  label="操作"  width="240" >
                    <template  #default="scope">
                        <el-button link type="primary" size="small" @click="openForm(scope.row)">
                            编辑
                            </el-button>

                        <el-popconfirm title="确定删除流程吗?"
                            @confirm="delRecord(scope.row.id)"
                            hide-icon
                        >
                            <template #reference>
                                <el-button link type="primary" size="small">
                                    删除
                                </el-button>
                            </template>
                        </el-popconfirm>   
                      
                        <el-button link type="primary" size="small" @click="openFlowDesign(scope.row.id)">流程设计</el-button>
                        <el-popconfirm title="确定发布新版本吗?"
                            @confirm="publishFlow(scope.row.id)"
                            hide-icon
                        >
                            <template #reference>
                                <el-button link type="primary" size="small" >
                                    版本发布
                                </el-button>
                            </template>
                        </el-popconfirm>
                       
                    </template>
                </el-table-column>
            </el-table>

            <el-empty description="无数据"   v-else></el-empty>
      
            <div style="padding: 5px 0;"></div>
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
        </el-card>
 
    <DeSignForm  ref="formRef" @tableReload="tableReload"></DeSignForm>
</template>

<script setup>
    import { reactive, ref ,onMounted, defineExpose } from "vue"
    import DeSignForm from './form.vue'
    import {pageApi, delApi, publishApi} from '~/api/flow/flowDesign.js'
    import { useRoute, useRouter } from 'vue-router'
    const router = useRouter()
import { ElMessage, ElNotification } from "element-plus";
    const queryData = ref({
        flowCode: '',
        flowName: '',
        category: '',
    })
    const tableData = ref([])
    const pagination = ref({
        size: 20,
        current:   1,
        total:0,
    })


    onMounted(()=>{
        tableReload();
    })
    
    const pageTable = (current, size) => {
        tableReload();
    }

    const delRecord = (id) => {
        delApi([{id:id}]).then(r=>{
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

    const openFlowDesign = (id) => {
        // 使用path 传递params 无效，需要使用name
        router.push({
            path: `/flowconfig/${id}`, 
          
        })
        
    }

    const publishFlow =(id) => {
        publishApi({id:id}).then(res=> {
            ElNotification({
                title: '提示',
                message: '版本发布成功',
                type: 'success',
            })
            tableReload();
            // console.log(res);
        })
    }
   

</script>

<style lang="scss" scoped>

</style>