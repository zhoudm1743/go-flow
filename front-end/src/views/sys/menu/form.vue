<template>
    <el-dialog
         v-model="dialogVisible"
        :title="title"
        width="700"
        height="600px"
        :before-close="handleClose"
    >
        <el-form 
            label-position="top"
            ref="ruleFormRef"
            label-width="auto"
            :rules="rules" 
            :model="formData" 
            label-suffix=":"
            status-icon
            class="demo-form-inline">
            <el-row gutter="20">
                <el-col :span="24">
                    <el-form-item label="上级菜单" prop="parentId" >
                        <el-tree-select
                            :check-strictly=true
                            v-model="formData.parentId"
                            :data="menuTree"
                            node-key="id"
                            :props="{
                                children: 'children', //自定义属性
                                label: 'name',
                                value: 'id'
                            }"
                            :render-after-expand="false"
                        />
                        <!-- <el-input v-model="formData.parentId" placeholder="父级菜单" clearable /> -->
                    </el-form-item>
                </el-col>
              
                <el-col :span="12">
                    <el-form-item label="菜单名称" prop="name" >
                        <el-input v-model="formData.name" placeholder="菜单名称" clearable />
                    </el-form-item>
                </el-col>
                <el-col :span="12">
                    <el-form-item label="菜单类型" prop="type" >
                        <el-radio-group v-model="formData.type">
                            <el-radio :value="item.value"  v-for="(item, index)   in getTypes()" :key="index">{{item.label}}</el-radio>
                        </el-radio-group>
                    </el-form-item>
                </el-col>
                <el-col :span="12">
                    <el-form-item label="菜单图标" prop="icon" >
                        <el-input v-model="formData.icon" placeholder="菜单图标" clearable disabled >
                            <template #append>
                                <el-button @click="openIconSelect"  type="primary"  plain>选择</el-button>
                            </template>
                        </el-input>
                    </el-form-item>
                </el-col>
                <el-col :span="12">
                    <el-form-item label="路由地址" prop="path" >
                        <el-input v-model="formData.path" placeholder="路由地址" clearable />
                    </el-form-item>
                </el-col>
                <el-col :span="24" v-if="formData.type == '2'">
                    <el-form-item label="组件路径" prop="component" >
                        <el-input v-model="formData.component" placeholder="组件路径" clearable >
                            <template #prepend>src/views/</template>
                        </el-input>
                    </el-form-item>
                </el-col>
                <el-col :span="12">
                    <el-form-item label="是否可见">
                        <el-radio-group v-model="formData.status" > 
                            <el-radio :value="item.value" v-for="(item, index)   in getStatus()" :key="index">{{ item.label }}</el-radio>
                        </el-radio-group>
                    </el-form-item>
                </el-col>
                <el-col :span="12">
                    <el-form-item label="排序" prop="sortNum">
                        <el-input-number :min="0"    controls-position="right" v-model="formData.sortNum" placeholder="排序" clearable />
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

    <IconSelect ref="iconSelectRef" @selectIcon="selectIcon"></IconSelect>
</template>

<script setup>
    import { ref, toRefs, reactive,watch   } from 'vue'
    import IconSelect from  '~/components/icon/index.vue'
    import {addApi, editApi } from '~/api/sys/menu.js'
    import { getTypes,getStatus } from './menu';
import { it } from 'element-plus/es/locales.mjs';
    const props = defineProps({
        menuTree: {
            type: Array,
            default:[],
        }
    })

    const iconSelectRef = ref();

    const dialogVisible = ref(false);

    const title = ref("新增");
    const formData = ref()
    const ruleFormRef = ref()

    watch(()=> props.menuTree,(newv,oldv)=>{
        console.log("newv:",newv)
    })


    const formTypeValidate = (rule, value, callback)=> {
        if(formData.value.type  != null){
                callback();   
        }
        callback(new Error('请选择菜单类别'))
    }
    const rules = reactive({
        name:[{ required:true, message: "请输入名称", trigger: 'blur' }],
        parentId: [{ required:true, message: "必填", trigger: 'change' }],
        type: [{ required:true, message: "请选择菜单类别", trigger: 'change' }],
        path: [{ required:true, message: "请输入路由地址", trigger: 'blur' }],
        component: [{ required:true, message: "必填", trigger: 'blur' }],

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