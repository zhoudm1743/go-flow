<template>
    <div>
        自定义页面， 测试加载能力
        <el-form   ref="ruleFormRef" :model="form" label-width="auto"  :rules="rules">
            <el-form-item label="活动名称"  prop="name">
                <el-input v-model="form.name" />
            </el-form-item>
            <el-form-item label="活动地点" prop="region">
                <el-select v-model="form.region" placeholder="please select your zone"  >
                    <el-option label="Zone one" value="shanghai" />
                    <el-option label="Zone two" value="beijing" />
                </el-select>
            </el-form-item>
        </el-form>

    </div>
</template>

<script setup>
import { reactive, ref } from 'vue';
    const ruleFormRef = ref(null)
    const form = ref({
        name: '',
        region: '',
    }) 
    const rules = reactive({
        name:[
            {required: true, message: '请输入活动名称', trigger: 'blur'}
        ],
        region:[
            {required: true, message: '请输入活动名称', trigger: 'blur'}
        ]
    })


    const submitForm = async () => {
        return new Promise((resolve, reject) => {
            ruleFormRef.value.validate((valid, fields) => {
                if (valid) {
                    resolve(form.value)
                } else {
                    reject()
                }
            })
        })
 
    }
    defineExpose({submitForm})

</script>

<style lang="scss" scoped>

</style>