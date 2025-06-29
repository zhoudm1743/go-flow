<template>
    <el-dialog
         v-model="dialogVisible"
        :title="title"
        width="500"
        :before-close="handleClose"
    >
        <el-form 
            ref="ruleFormRef"
            label-width="auto"
            :rules="rules" 
            :model="formData" 
            label-suffix=":"
            status-icon
            class="demo-form-inline">
            <el-form-item label="角色编码" prop="roleCode" >
                <el-input v-model="formData.roleCode" placeholder="角色编码" clearable />
            </el-form-item>
            <el-form-item label="角色名称" prop="roleName" >
                <el-input v-model="formData.roleName" placeholder="角色名称" clearable />
            </el-form-item>
            <el-form-item label="状态">
                <el-radio-group v-model="formData.roleStatus" > 
                    <el-radio value="1">正常</el-radio>
                    <el-radio value="0">停用</el-radio>
                </el-radio-group>
            </el-form-item>
            <el-form-item label="排序" prop="sortNum">
                <el-input-number :min="0"    controls-position="right" v-model="formData.sortNum" placeholder="排序" clearable />
            </el-form-item>
            <el-form-item label="描述" prop="roleDesc">
                <el-input
                    v-model="formData.roleDesc"
                    :autosize="{ minRows: 4, maxRows: 8 }"
                    type="textarea"
                    placeholder="描述信息"
                />
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
    import { ref, toRefs, reactive  } from 'vue'
    import IconSelect from  '~/components/icon/index.vue'
    import {addApi, editApi } from '~/api/sys/role.js'
  
    const iconSelectRef = ref();

    const dialogVisible = ref(false);

    const title = ref("新增");
    const formData = ref()
    const ruleFormRef = ref()

    const formPathValidate = (rule, value, callback)=> {
        if(formData.value.formCustom == 'Y'){
            if(value.endsWith(".vue") && value.length > 4){
                callback();   
            }
            callback(new Error('请配置customform下的vue文件路径， 并以.vue结尾'))
        }
    }
    const rules = reactive({
        formPath: [{ validator: formPathValidate, trigger: 'blur' }],
        categoryId:[{ required:true, message: "请选择分类", trigger: 'change' }],
        flowCode: [{ required:true, message: "必填", trigger: 'blur' }],
        flowName: [{ required:true, message: "必填", trigger: 'blur' }],
    })

    const selectIcon = (value) => {
        formData.value.icon = value;
    }

    const openIconSelect = () => {
        iconSelectRef.value.openDialog();
    }

    const categoryList = ref([])
    const openForm = (record) => {
        if(record && record.id != null){
                formData.value = record
                title.value = "编辑"
           
            }else{
                formData.value = {
                    roleCode: '',
                    roleName: '',
                    sortNum: 0,
                    roleStatus: '1',
                    roleDesc: '',
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