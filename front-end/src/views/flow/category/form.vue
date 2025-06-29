<template>
        <el-dialog
         v-model="dialogVisible"
        :title="title"
        width="500"
        :before-close="handleClose"
    >
    <el-form   
        ref="ruleFormRef"
        :rules="rules" 
        label-width="auto"
        :model="formData" 
        label-suffix=":"
        status-icon
        class="demo-form-inline">
        <el-form-item label="分类编码" prop="code">
            <el-input v-model="formData.code" placeholder="分类编码" clearable />
        </el-form-item>
        <el-form-item label="分类名称" prop="name">
            <el-input v-model="formData.name" placeholder="分类名称" clearable />
        </el-form-item>
        <!-- <el-form-item label="图标" prop="icon" >
         

            <el-input v-model="formData.icon" placeholder="图标" clearable disabled >
                <template #append>
                    <el-button @click="openIconSelect"  type="primary"  plain>选择</el-button>
                </template>
            </el-input>
        </el-form-item> -->
        <el-form-item label="排序号" prop="sortNum">
            <el-input-number :min="0"  :max="9999"  controls-position="right" v-model="formData.sortNum" placeholder="排序号" clearable />
        </el-form-item>
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

    <IconSelect ref="iconSelectRef" @selectIcon="selectIcon"></IconSelect>
</template>

<script setup>
    import { ref, defineEmits, reactive  } from 'vue'
    import IconSelect from  '~/components/icon/index.vue'
    import {addApi, editApi } from '~/api/flow/flowCategory.js'
    const iconSelectRef = ref();

    const dialogVisible = ref(false);

    const title = ref("新增");
    const formData = ref({
        code: '',
        name: '',
        sortNum: 0
    })

    const selectIcon = (value) => {
        formData.value.icon = value;
    }

    const openIconSelect = () => {
        iconSelectRef.value.openDialog();
    }

    const openForm = (record) => {
        if(record && record.id != null){
            formData.value = record
            title.value = "编辑"
           
        }else{
            formData.value = {
                code: '',
                name: '',
                sortNum: 0,
                icon:'',
            }
            title.value = "新增"
        }
        dialogVisible.value = true;
    }

    defineExpose({
        openForm
    })

    const ruleFormRef = ref()
  

    const rules = reactive({
        code:[
            {required: true, trigger: 'blur', message: '请输入分类编码'}
        ],
        name:[
            {required: true, trigger: 'blur', message: '请输入分类名称'}
        ],
        icon: [
            {required: true, trigger: 'blur', message: '请选择图标'}
        ]
    })

    const emits = defineEmits(['tableReload'])

    const submitForm = async (formEl) => {
        if (!formEl) return
        await formEl.validate((valid, fields) => {
            if (valid) {
            console.log('submit!')
            if(formData.value.id != null ){
                editApi(Object.assign({}, formData.value)).then(res=> {
                emits('tableReload')
                dialogVisible.value = false;
                    console.log(res)
                }).catch(e=> {
                    console.log(e);
                }) 
            }else{
                addApi(Object.assign({}, formData.value)).then(res=> {
                emits('tableReload')
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