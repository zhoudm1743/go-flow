<template>
    <el-dialog
         v-model="dialogVisible"
        :title="title"
        :before-close="handleClose"
    >
        <el-form 
            label-position="top"
            ref="ruleFormRef"
            label-width="auto"
            :model="formData" 
            label-suffix=":"
            status-icon
            style="height: 240px;"
            class="demo-form-inline">
            <el-row gutter="20">
              
              
                <el-col :span="12">
                    <el-form-item label="权限编码" prop="code"  required   error="请输入权限编码">
                        <el-input v-model="formData.code" placeholder="权限编码" clearable  />
                    </el-form-item>
                </el-col>
                <el-col :span="12">
                    <el-form-item label="权限名称" prop="name" required   error="请输入权限名称">
                        <el-input v-model="formData.name" placeholder="权限名称" clearable />
                    </el-form-item>
                </el-col>
              
                <el-col :span="12">
                    <el-form-item label="排序" prop="sortNum">
                        <el-input-number :min="0"    controls-position="right" v-model="formData.sortNum" placeholder="排序" clearable />
                    </el-form-item>
                </el-col>
                <el-col :span="12">
                    <el-form-item label="类别" prop="type" required>
                        <el-radio-group v-model="formData.type" > 
                            <el-radio value="1" >接口权限</el-radio>
                            <el-radio value="2" >按钮权限</el-radio>
                        </el-radio-group>
                    </el-form-item>
                </el-col>
            </el-row>
        </el-form>

        
        <template #footer>
        <div class="dialog-footer">
            <el-space>
                <el-button @click="dialogVisible = false">取消</el-button>
                <el-button type="primary" @click="submitForm(ruleFormRef)">
                确定
                </el-button>
            </el-space>
        </div>
        </template>
    </el-dialog>

</template>

<script setup>
    import { ref, toRefs, reactive,watch, defineExpose   } from 'vue'
    import {addApi, editApi } from '~/api/sys/identifier.js'
   


    const dialogVisible = ref(false);

    const title = ref("新增");
    const formData = ref()
    const ruleFormRef = ref()

    const openForm = (record) => {
        if(record && record.id != null){
                formData.value = record
                title.value = "编辑"
           
            }else{
                formData.value = {
                    code: '',
                    name: '',
                    sortNum: 0,
                    type: record.type,
                    menuId: record.menuId
                }
                title.value = "新增"
            }
            dialogVisible.value = true;
 
    }
    

    defineExpose({
        openForm
    })
    
    const emits = defineEmits(['tableReload'])

    const submitForm = async (formEl) => {
        if (!formEl) return
        await formEl.validate((valid, fields) => {
            if (valid) {
            console.log('submit!')
            if(formData.value.id != null ){
                editApi(Object.assign({}, formData.value)).then(res=> {
                emits('tableReload');
                dialogVisible.value = false;
                    console.log(res)
                }).catch(e=> {
                    console.log(e);
                }) 
            }else{
                addApi(Object.assign({}, formData.value)).then(res=> {
                emits('tableReload');
                dialogVisible.value = false;
                    console.log(res)
                }).catch(e=> {
                    console.log(e);
                })
            }
            
        } else {
            console.log('error submit!', fields)
        }
    })
    }

</script>

<style  scoped>
</style>