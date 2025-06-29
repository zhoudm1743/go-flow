<template>
    <el-dialog
         v-model="dialogVisible"
        title="用户选择"
        width="900"
        style="max-height: 80%;"
        :before-close="handleClose"
    >
    <div  style="display: flex; flex-direction: column; gap: 10px; min-height: 400px;">
        <el-card  :body-style="{'padding-bottom': '0px'}">
         <el-form :inline="true" :model="formInline" class="demo-form-inline" label-suffix=":"> 
            <el-form-item label="关键词" prop="keywords" >
                 <el-input v-model="queryData.keywords" placeholder="请输入用户名或昵称"  />
             </el-form-item>
             <el-form-item label="手机号" prop="phone">
                 <el-input v-model="queryData.phone" placeholder="手机号"  />
             </el-form-item>
             <el-form-item>
                 <el-button type="primary" @click="tableReload()" plain>查询</el-button>
             </el-form-item>
         </el-form>
     </el-card>
    
     <el-card >
         <el-table 
            :data="tableData"  
            :border="true" style="width: 100%;"  
            row-key="id"   
             @selection-change="handleSelectionChange"
            v-if="tableData.length > 0">
            <el-table-column type="selection" width="55" />
             <el-table-column prop="username" label="用户名"/>
             <el-table-column prop="nick" label="昵称" />
             <el-table-column prop="phone" label="手机号" />
             <!-- <el-table-column prop="email" label="邮箱" /> -->
             <el-table-column prop="roleStatus" label="状态" >
                <template #default="scope">
                    <el-switch v-model="scope.row.status"
                        disabled 
                        active-value="1"
                        inactive-value="0"
                    />

                </template>
            </el-table-column>
             <el-table-column prop="createTime" label="创建时间" >
             </el-table-column>
         </el-table>

         <el-empty description="无数据"   v-else></el-empty>
   
         <div style="padding: 5px 0;"></div>
         <el-row>
             
             <el-col :span="24">
                 <el-pagination 
                     background layout="prev, pager, next,sizes,  total" 
                     :total="pagination.total"
                     :page-sizes="[20, 50, 100, 200]"
                     v-model:page-size="pagination.size"
                     v-model:current-page="pagination.current"
                     @change="pageTable"
                 />
             </el-col>
         </el-row>
     </el-card>
    </div>

    <template #footer>
            <el-space>
                <el-button @click="dialogVisible = false">取消</el-button>
                <el-button type="primary" @click="submitGrantUsre()">
                确定
                </el-button>
            </el-space>
        </template>
    </el-dialog>
  
</template>

<script setup>
    import { reactive, ref ,onMounted, defineExpose } from "vue"

    import { pageApi } from "~/api/sys/user";
    import { canGrantPageApi } from "~/api/sys/roleuser";
    import { grantUserApi } from "~/api/sys/role";


    const queryData = ref({
        keywords: '',
        phone: '',
    })

    const dialogVisible = ref(false)
    const tableData = ref([])
    
    const pagination = ref({
        size: 20,
        current:   1,
        total:0,
    })

    onMounted(()=>{
        
        tableReload();
        
    })

    const tableSelect = ref([])
    const handleSelectionChange = (rows) => {
        tableSelect.value = rows;
    }
    
    const tableReload = () => {
        const param = {};
        Object.assign(param, queryData.value);
        param.current = pagination.value.current;
        param.size = pagination.value.size;
        
        canGrantPageApi(param).then(res=> {
            pagination.value.total = res.total;
            tableData.value = res.rows;
        })
    }


    const changeRoleStatus = (value, row) => {
        // if(value == '1'){
        //     enableApi([{id: row.id}]).then(()=> {
        //     }).catch(e=> {
        //         row.roleStatus = '0';
        //     })
        // }else{
        //     disableApi([{id: row.id}]).then(()=> {
        //     }).catch(e=> {
        //         row.roleStatus = '1';
        //     })
        // }
    }

    //   打开页面
    const openGrantUserForm = (roleId) => {
        queryData.value.roleId = roleId;
        tableReload();
        dialogVisible.value = true;
        
    }

    const emits = defineEmits(['tableReload'])
    /**
     * 授权
     */
    const submitGrantUsre = () => {
        if(tableSelect.value.length <= 0){
            alert("没有选中数据")
        }
        const data = [];
        var roleId = queryData.value.roleId;
        for(var i in tableSelect.value){
            data.push({
                roleId:roleId,
                userId: tableSelect.value[i].id
            })
        }
        grantUserApi(data).then(res=> {
            emits('tableReload')
            dialogVisible.value = false;
        })

    }
    defineExpose({openGrantUserForm})

</script>

<style lang="scss" scoped>

</style>