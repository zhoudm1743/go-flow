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
        <el-form-item label="流程编码" prop="flowCode" >
            <el-input v-model="formData.flowCode" placeholder="流程编码"  />
        </el-form-item>
        <el-form-item label="流程名称" prop="flowName" >
            <el-input v-model="formData.flowName" placeholder="流程编码"  />
        </el-form-item>
        <el-form-item label="业务分类" prop="categoryId">
            <el-select
                v-model="formData.categoryId"
                filterable
                placeholder="业务分类"
            >
                <el-option
                v-for="item in categoryList"
                :key="item.id"
                :label="item.name"
                :value="item.id"
                />
            </el-select>
        </el-form-item>
        <el-form-item label="流程管理员" prop="managerId"  required>
            <div style="display: flex; gap: 10px; align-items: center;">
            <span><a href="#" @click.prevent="openUserSelect">选择</a>
            </span>
            <el-tag closable @close="formData.managerId = null" v-if="formData.managerId != null">{{ formData.managerName }}</el-tag>
            </div>
            <span  v-if="formData.managerId == null" class="el-form-item__error">
                {{tips}}
            </span>
        </el-form-item>
        <el-form-item label="表单类型">
            <el-radio-group v-model="formData.formCustom"> 
                <el-radio value="N">设计表单</el-radio>
                <el-radio value="Y">开发表单</el-radio>
            </el-radio-group>
        </el-form-item>

        <el-form-item label="表单路径" v-show="formData.formCustom === 'Y'" prop="formPath">
            <el-input v-model="formData.formPath" placeholder="表单路径" clearable />
        </el-form-item>
        <el-form-item label="图标" prop="icon" >

         <el-input v-model="formData.icon" placeholder="图标" clearable disabled >
             <template #append>
                 <el-button @click="openIconSelect"  type="primary"  plain>选择</el-button>
             </template>
         </el-input>
     </el-form-item>
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
    <Userselect ref="userSelectRef" @roleSelectData="userSelectData" ></Userselect>
</template>

<script setup>
    import { ref, toRefs, reactive,defineEmits  } from 'vue'
    import IconSelect from  '~/components/icon/index.vue'
    import Userselect from '../../sys/user/userSelect.vue'
    import {addApi, editApi } from '~/api/flow/flowDesign.js'
    import { listApi } from '~/api/flow/flowCategory'; 
    const iconSelectRef = ref();
    const userSelectRef = ref()
    const dialogVisible = ref(false);

    const title = ref("新增");
    const formData = ref({
        flowCode: '',
        flowName: '',
        sortNum: 0,
        formCustom: 'Y',
        formPath: '',
        // version: 0,
        icon: '',
        managerId: null,
        managerName: '',
        

    })

    const openUserSelect = () => {
        userSelectRef.value.open(false);
    }
    const userSelectData = (userList) => {
        if(userList.length > 0){
            formData.value.managerId = userList[0].id;
            formData.value.managerName = userList[0].nick;
        }
    }

    const ruleFormRef = ref()

    const formPathValidate = (rule, value, callback)=> {
        if(formData.value.formCustom == 'Y'){
            if(value.endsWith(".vue") && value.length > 4){
                callback();   
            }
            callback(new Error('请配置customform下的vue文件路径， 并以.vue结尾'))
        }
        callback(); 
    }
    // 处理 自定义表单的校验 blur 无效问题
    const tips = ref('') 
    const managerValidate = (rule, value, callback)=> {
        if(value == null){
            tips.value = '管理员不能为空';
            callback(new Error(''))
        }
        callback(); 
    }
    const rules = reactive({
        formPath: [{ validator: formPathValidate, trigger: 'blur' }],
        categoryId:[{ required:true, message: "请选择分类", trigger: 'change' }],
        flowCode: [{ required:true, message: "必填", trigger: 'blur' }],
        flowName: [{ required:true, message: "必填", trigger: 'blur' }],
        managerId:[{ validator: managerValidate,  trigger: 'blur' }],
    })

    const selectIcon = (value) => {
        formData.value.icon = value;
    }

    const openIconSelect = () => {
        iconSelectRef.value.openDialog();
    }

    const categoryList = ref([])
    const openForm = (record) => {
        listApi({}).then(res => {
            categoryList.value = res.rows;
            if(record && record.id != null){
                formData.value = record
                title.value = "编辑"
           
            }else{
                formData.value = {
                    code: '',
                    name: '',
                    sortNum: 0,
                    formCustom: 'Y',
                    formPath: '',
                    icon:'',
                    managerId: null,
                }
                title.value = "新增"
            }
            dialogVisible.value = true;
        }) 
 
    }

    defineExpose({
        openForm
    })

    const emits = defineEmits(['tableReload']);
    
    const submitForm = async (formEl) => {
        if (!formEl) 
            return
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