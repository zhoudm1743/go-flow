<template>
    <el-dialog
         v-model="dialogVisible"
        :title="title"
        width="800"
        :before-close="handleClose"
    >
        <el-form 
            ref="ruleFormRef"
            label-position="top"
            label-width="auto"
            :rules="rules" 
            :model="formData" 
            label-suffix=":"
            status-icon
            class="demo-form-inline">
            <el-row gutter="20">
                <el-col :span="12">
                    <el-form-item label="用户名" prop="username" >
                        <el-input v-model="formData.username" placeholder="用户名" clearable />
                    </el-form-item>
                </el-col>
                <el-col :span="12">
                    <el-form-item label="昵称" prop="nick" >
                        <el-input v-model="formData.nick" placeholder="昵称" clearable />
                    </el-form-item>
                </el-col>
                <el-col :span="12" v-if="!formData.id">
                    <el-form-item label="密码" prop="password"  >
                        <el-input  type="password" v-model="formData.password" placeholder="密码" clearable />
                    </el-form-item>
                </el-col>
               
                <el-col :span="12">
                    <el-form-item label="手机号" prop="phone" >
                        <el-input v-model="formData.phone" placeholder="手机号" clearable />
                    </el-form-item>
                </el-col>
                <el-col :span="12">
                    <el-form-item label="邮箱" prop="email" >
                        <el-input v-model="formData.email" placeholder="邮箱" clearable />
                    </el-form-item>
                </el-col>
                <el-col :span="12">
                    <el-form-item label="状态">
                        <el-radio-group v-model="formData.status" > 
                            <el-radio value="1">正常</el-radio>
                            <el-radio value="0">停用</el-radio>
                        </el-radio-group>
                    </el-form-item>
                </el-col>
                <el-col :span="24" style="padding-bottom: 4px;     gap: 10px; display: flex;">
                    <span style="font-weight: bolder; font-size: 15px;">任职信息</span>
                    <el-button  type="primary" size="small" @click="positionSelectVisible = true">
                        新增任职
                    </el-button>
                </el-col>
                <el-col :span="24">
                    <el-table :data="formData.positionList" style="width: 100%">
                        <el-table-column prop="cmpName" label="公司" />
                        <el-table-column prop="orgName" label="部门" />
                        <el-table-column prop="positionName" label="职位" />
                        <el-table-column  label="操作" width="60" >
                        <template  #default="scope">
                                <el-button  link type="primary" size="small" @click="handleDelete(scope.$index)">
                                删除
                                </el-button>
                        </template>
                    </el-table-column>
                    </el-table>
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
    <el-dialog   v-model="positionSelectVisible"  title="新增任职" width="60%" >
        <PositionSelect @positionSelect="positionSelect"></PositionSelect>
    </el-dialog>
    
</template>

<script setup>
    import { ref, toRefs, reactive  } from 'vue'
    import IconSelect from  '~/components/icon/index.vue'
    import PositionSelect from  '~/views/sys/position/positionSelect.vue'
    import {addApi, editApi, detailApi } from '~/api/sys/user.js'
    // 任职信息表格
    // const tableData = ref([]);

    const iconSelectRef = ref();

    const dialogVisible = ref(false);

    const title = ref("新增");
   
    const ruleFormRef = ref()
    const formData = ref({})

    const positionSelectVisible = ref(false);

    const positionSelect = (position) => {
        console.log(position);
        if(!formData.value.positionList){
            formData.value.positionList = [];
        }
        const hasPosition =  formData.value.positionList.some(item=>{
            if(item.positionId == position.positionId){
                return true 
            } 
        })
        if(!hasPosition){
            formData.value.positionList.push(position);
        }
        
        // tableData.value.push(position);
        positionSelectVisible.value = false;
    }



    const passwordValidate = (rule, value, callback)=> {
        if(formData.value.id != null){
            if(value.length >= 6){
                callback();   
            }
            callback(new Error('密码长度不能小于6位'))
        }
        callback();
    }

    const rules = reactive({
       
        username: [{ required:true, message: "必填", trigger: 'blur' }],
        password: [{ validator: passwordValidate, trigger: 'blur' },
        { required:true, trigger: 'blur' }
        ],
        nick:  [{ required:true, message: "必填", trigger: 'blur' }],

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
            title.value = "编辑"
            detailApi({id: record.id}).then(res=> {
                formData.value = res.data;
            })
                // formData.value = record
                
           
            }else{
                formData.value = {
                    username: '',
                    nick: '',
                    password: '',
                    phone: '',
                    email: '',
                    status: "1",
                    sortNum: '',
                    positionList:[]
                }
                title.value = "新增"
            }
            dialogVisible.value = true;
 
    }


    const handleDelete = (index) => {
        formData.value.positionList.splice(index, 1);
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