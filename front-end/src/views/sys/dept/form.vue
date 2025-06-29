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
                <el-col :span="24" v-if="formData.pid == 0">
                    <el-form-item label="上级部门" prop="pid" >
                        <el-tree-select
                            v-model="formData.pid"
                            :data="menuTree"
                            default-expand-all
                            node-key="id"
                            :props="{
                                children: 'children', //自定义属性
                                label: 'name',
                                value: 'id'
                            }"
                        />
                        <!-- <el-input v-model="formData.parentId" placeholder="父级菜单" clearable /> -->
                    </el-form-item>
                </el-col>
              
                <el-col :span="12">
                    <el-form-item label="名称" prop="name" >
                        <el-input v-model="formData.name" placeholder="名称" clearable />
                    </el-form-item>
                </el-col>
                <el-col :span="12">
                    <el-form-item label="类型" prop="levelCode" >
                        <el-radio-group v-model="formData.levelCode">
                            <el-radio :value="item.value"  v-for="(item, index)   in getLevelCodes()" :key="index">{{item.label}}</el-radio>
                        </el-radio-group>
                    </el-form-item>
                </el-col>
                
                <el-col :span="12">
                    <el-form-item label="部门主管" prop="director" >
                        <el-select v-model="formData.director" placeholder="部门主管"  clearable>
                            <el-option :label="item.username" :value="item.id"   v-for="(item, index)   in userList" :key="index"/>
                        </el-select>
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
    import {addApi, editApi } from '~/api/sys/dept.js'
    import { getLevelCodes } from './dept';
    import { selectUsersApi } from '~/api/sys/user';
    const props = defineProps({
        menuTree: {
            type: Array,
            default:[],
        }
    })

    const userList = ref([]);

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
        pid: [{ required:true, message: "请选择上级部门", trigger: 'change' }],
        levelCode: [{ required:true, message: "请选择类别", trigger: 'change' }],

    })

    const selectIcon = (value) => {
        formData.value.icon = value;
    }

    const openIconSelect = () => {
        iconSelectRef.value.openDialog();
    }

    const categoryList = ref([])

    /**
     * @description: 打开form表单
     * @param {object} record 当前行数据
     * @return {void}
     */
    const openForm = (record) => {
        selectUsersApi().then(res=> {
            userList.value = res.rows
        })

        if(record && record.id != null){
            formData.value = record
            title.value = "编辑"
           
        }else{
            formData.value = {
                name: '',
                levelCode: 'dept',
                pid: null,
                sortNum: 0,
                director: '',
            }
            title.value = "新增"
        }
        dialogVisible.value = true;
 
    }
    
    const addChildren = (record) => {
        selectUsersApi().then(res=> {
            userList.value = res.rows
        })

        formData.value = {
                    name: '',
                    levelCode: 'dept',
                    pid: record.pid,
                    sortNum: 0,
                    director: '',
                }
        title.value = "新增"
        dialogVisible.value = true;
    }

    defineExpose({
        openForm, addChildren
    })
    
    const emits = defineEmits(['tableReload']);
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