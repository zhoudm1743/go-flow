<template>
    <el-tabs v-model="activeName" class="demo-tabs" @tab-click="handleClick" type="border-card" stretch>
        <el-form
                :model="formData"
                label-position="top"
                label-width="auto"
                :rules="rules"
                >
            <el-tab-pane label="操作权限" name="first">
                
                <div class="flow-operate-setting">
                    <el-row >
                        <template  v-for="(item, index) in formData.buttons" :key="index">
                            <el-col :span="6" >
                                <el-checkbox v-model="item.checked" :disabled="operateNames[item.type].disabled" :label="operateNames[item.type].label" />
                            </el-col>
                            <el-col :span="18">
                                <el-form-item required>
                                    <el-input v-model="item.text"   placeholder=""  maxlength="10"  />
                                </el-form-item>
                            </el-col>
                        </template>
                    </el-row>
                </div>
            
            </el-tab-pane>
            <el-tab-pane label="表单配置" name="second">
                <el-form-item label="节点表单" prop="formPath" >
                    <el-input
                        v-model="formData.formPath"
                        class="mx-4"
                        placeholder="当前节点审批表单，在src/views/customform/文件夹下"                   
                        controls-position="right"
                    >
                    </el-input>
                </el-form-item>
                <el-row v-if="formData.fieldSetting && formData.fieldSetting.length > 0">
                    <el-col :span="24">
                        <span style="font-size: 14px;">配置节点表单字段后，字段配置将失效</span>
                    </el-col>
                    <el-col :span="24">
                        <p>字段配置 </p>
                    </el-col>
                </el-row>
                    <el-form-item label-position="left" label-width="auto" :label="item.label" v-for="(item,index) in formData.fieldSetting" :key="index">
                        <el-radio-group :name="item.label" v-model="item.permission">
                        <el-radio value="1" >编辑</el-radio>
                        <el-radio value="2" >只读</el-radio>
                        <el-radio value="3" >隐藏</el-radio>
                        </el-radio-group>
                    </el-form-item>
                       
            </el-tab-pane>
            <el-tab-pane label="监听器" name="third">
                <el-row >
                    <el-col :span="24">  <el-button @click="addCondition(0)" title="添加监听器"  icon="Plus" circle /> </el-col>
                </el-row>
                <el-row  :gutter="10" style="margin-top: 10px; margin-right: 10px;" v-for="(item, index) in formData.listeners" :key="index">
                    
                    <el-col :span="8">
                        <el-select v-model="item.listenerType" placeholder="条件类型" >
                            <el-option  label="任务创建"  value="create"  /> 
                            <el-option  label="任务开始办理"  value="start"  /> 
                            <el-option  label="分派监听器"  value="assignment"  /> 
                            <el-option  label="任务完成"  value="finish"  /> 
                        </el-select>
                    </el-col>
                    <el-col :span="12">
                        <el-input v-model="item.listenerPath" placeholder="监听器类路径" />
                    </el-col>
                    <el-col :span="4" style=" display: flex; align-items: center; ">
                        <el-button @click="addCondition(index+1)"   icon="Plus" circle />
                        <el-button @click="removeCondition(index)" type="danger"  icon="Minus" circle/>
                    </el-col>
                </el-row>
            </el-tab-pane>
        </el-form>
 </el-tabs>
</template>

<script setup>
   import {defineProps, defineExpose, toRefs,watch, ref, reactive} from 'vue';
   const props = defineProps({
       data:{},
   })
   const activeName = ref('first');
   const operateNames = ref({
        aggren : {label: "同意", disabled: true},
        reject : {label: "撤销", disabled: false},
        cancel : {label: "取消", disabled: true},
    })
   const formData = reactive({
        formPath: "",
        listeners:[],
        buttons:[
                {
                    type: 'aggren',
                    checked: true,
                    text:"同意",
                },
                {
                    type: 'reject',
                    checked: true,
                    text:"撤销",
                },
            ],
        fieldSetting:[],
   })

   const rules = ref({})

   const {data} = toRefs(props);
   watch(()=>data,(newv,oldv)=>{
       alert(newv);
   })



   const addCondition = (index) => {
        formData.listeners.splice(index, 0, {
            listenerType: '', //监听器类型
            listenerPath:'', // 监听器类路径
        })
    }
    
    const removeCondition = (index) => {
        formData.listeners.splice(index, 1);
    }


   Object.assign(formData, data.value.properties);
    const formConfig = () => {
        return new Promise((resolve, reject) => {
            formData.value = "所有人";
            resolve(formData);
        });
    }

    defineExpose({
        formConfig
    })

  

</script>

<style  scoped>
    .flow-operate-setting{
        display: flex;
        flex-wrap: nowrap;
        flex-direction: column;
        gap: 10px;
    }
</style>