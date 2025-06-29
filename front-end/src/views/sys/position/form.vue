<template>
    <el-dialog
         v-model="dialogVisible"
        :title="title"
        width="700"
        :before-close="handleClose"
    >
        <el-form 
            style="min-height: 300px;"
            ref="ruleFormRef"
            label-width="auto"
            :rules="rules" 
            :model="formData" 
            label-suffix=":"
            status-icon
            class="demo-form-inline">
            <el-form-item label="所属部门" prop="deptId" >
                <el-tree-select
                            :check-strictly=true
                            v-model="formData.deptId"
                            :data="deptList"
                            :default-expended-all="true"
                            node-key="id"
                            :props="{
                                children: 'children', //自定义属性
                                label: 'name',
                                value: 'id'
                            }"
                            :render-after-expand="false"
                        />
            </el-form-item>
            <el-form-item label="职位名称" prop="positionName" >
                <el-input v-model="formData.positionName" placeholder="职位名称" clearable />
            </el-form-item>
          
            <el-form-item label="排序号" prop="sortNum" >
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

</template>

<script setup>
    import { de } from 'element-plus/es/locales.mjs';
import { ref, toRefs, reactive  } from 'vue'
    import {addApi, editApi } from '~/api/sys/position.js'

    const props = defineProps({
        deptList: {
            type: Array,
            default:[],
        }
    })

    const iconSelectRef = ref();

    const dialogVisible = ref(false);

    const title = ref("新增");
   
    const ruleFormRef = ref()
    const formData = ref({})


    const rules = reactive({
       
        deptId: [{ required:true, message: "必填", trigger: 'blur' }],
        positionName:  [{ required:true, message: "必填", trigger: 'blur' }],

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
                    deptId: '',
                    positionName: '',
                    sortNum: 0,
                }
                title.value = "新增"
            }
            dialogVisible.value = true;
 
    }

    const addPosition = (deptId) => {
        formData.value = {
            deptId: deptId,
            positionName: '',
            sortNum: 0,
        }
        title.value = "新增"
        dialogVisible.value = true;
    }


    defineExpose({
        openForm, addPosition
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