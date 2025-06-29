<template>
    <div>
        <el-form   ref="ruleFormRef" :model="form" label-width="auto"  :rules="rules" >
            <el-form-item label="类别"  prop="type">
                <el-radio-group v-model="form.type" class="ml-4">
                    <el-radio value="1" >病假</el-radio>
                    <el-radio value="2" >年假</el-radio>
                    <el-radio value="3" >婚假</el-radio>
                    <el-radio value="4" >其他</el-radio>
                </el-radio-group>
            </el-form-item>
            <el-form-item label="时间"  prop="times">
                <el-date-picker
                    v-model="form.times"
                    type="datetimerange"
                    format="YYYY-MM-DD HH:mm"
                    range-separator="到"
                    start-placeholder="开始时间"
                    end-placeholder="结束时间"
                    />
            </el-form-item>
            <el-form-item label="请假天数"  prop="days">
                <el-input-number v-model="form.days" placeholder="请假天数" clearable />
            </el-form-item>
            <el-form-item label="请假原因"  prop="reason" v-if="form.type === '4'">
                <el-input type="textarea" :rows="4"  v-model="form.reason" placeholder="请假原因" clearable />
            </el-form-item>
            <el-form-item label="审批页面" >
                    验证当前页面为审批页面 leaveTodo
            </el-form-item>
        </el-form>

    </div>
</template>

<script setup>
import { reactive, ref, toRefs, watch } from 'vue';
   const props = defineProps({
        variable:{
            type: Object,
        }
    })
    const form = ref({
        type: '',
        times: [],
        days:null,
        reason: ""
    })
    const { variable } = toRefs(props);
    watch(variable,(newv,oldv)=>{
        if(newv){
            form.value = newv
        }
    },{immediate: true})

    const ruleFormRef = ref(null)
    
    const rules = reactive({
        type:[
            {required: true, message: '请选择请假类型', trigger: 'blur'}
        ],
        days:[
            {required: true, message: '请输入请假天数', trigger: 'blur'}
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