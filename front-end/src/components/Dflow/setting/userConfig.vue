<template>
    <el-tabs v-model="activeName" class="demo-tabs" @tab-click="handleClick" type="border-card" stretch>
       <el-form
           :model="formData"
           label-position="top"
           label-width="auto"
       >
       <el-tab-pane label="审批设置" name="first">
       
           <el-form-item label="审批人组合" >
               <el-checkbox-group v-model="formData.permissions">
                   <el-checkbox label="角色" value="role" >
                           <span  class="flow-user-select-checkbox" >
                               角色   
                               <span style="width: 20px;" @click.prevent>
                                   <el-icon style="float: right" :size="14" @click="addApproval('role')" >
                                   <CirclePlus />
                               </el-icon>
                               </span>
                           </span>
                   </el-checkbox>
                   <el-checkbox label="用户" value="user" >
                       <span  class="flow-user-select-checkbox" >
                           用户   
                               <span style="width: 20px;" @click.prevent>
                                   <el-icon style="float: right" :size="14" @click="addApproval('user')" >
                                   <CirclePlus />
                               </el-icon>
                               </span>
                           </span>
                   </el-checkbox>
                   <el-checkbox label="部门" value="dept"  >
                       <span  class="flow-user-select-checkbox" >
                           部门   
                               <span style="width: 20px;" @click.prevent>
                                   <el-icon style="float: right" :size="14" @click="addApproval('dept')" >
                                   <CirclePlus />
                               </el-icon>
                               </span>
                           </span>
                   </el-checkbox>
                   <el-checkbox label="连续多级部门" value="levelDept" disabled>
                       <span  class="flow-user-select-checkbox" >
                           连续多级部门   
                               <span style="width: 20px;" @click.prevent>
                                   <el-icon style="float: right" :size="14" @click="addApproval('levelDept')" >
                                   <CirclePlus />
                               </el-icon>
                               </span>
                           </span>
                   </el-checkbox>
                   <el-checkbox label="发起人自己" value="initiator" disabled>
                       <span  class="flow-user-select-checkbox" >
                           发起人自己   
                               <span style="width: 20px;" @click.prevent>
                                   <el-icon style="float: right" :size="14" @click="addApproval('initiator')" >
                                   <CirclePlus />
                               </el-icon>
                               </span>
                           </span>
                   </el-checkbox>
                   <el-checkbox label="发起人的部门" value="initiatorLevelDept" disabled>
                       <span  class="flow-user-select-checkbox" >
                           发起人的部门   
                               <span style="width: 20px;" @click.prevent>
                                   <el-icon style="float: right" :size="14" @click="addApproval('initiatorLevelDept')" >
                                   <CirclePlus />
                               </el-icon>
                               </span>
                           </span>
                   </el-checkbox>
                   <el-checkbox label="发起人的连续多级部门" value="initiatorMutiLevelDept" disabled>
                       <span  class="flow-user-select-checkbox" >
                           发起人的连续多级部门   
                               <span style="width: 20px;" @click.prevent>
                                   <el-icon style="float: right" :size="14" @click="addApproval('initiatorMutiLevelDept')" >
                                   <CirclePlus />
                               </el-icon>
                               </span>
                           </span>
                   </el-checkbox>
                   <el-checkbox label="表单中的人员" value="formUser" disabled>
                       <span  class="flow-user-select-checkbox" >
                           表单中的人员   
                               <span style="width: 20px;" @click.prevent>
                                   <el-icon style="float: right" :size="14" @click="addApproval('formUser')" >
                                   <CirclePlus />
                               </el-icon>
                               </span>
                           </span>
                   </el-checkbox>
               </el-checkbox-group>
           </el-form-item>
           <el-form-item label="已选择" >
                <div   div class="flow-user-selected" v-show="formData.permissions.includes('dept')">
                    <span> 部门：</span>
                    <el-tag  closable  @close="selectClose('dept', item.id)" v-for="item in formData.combination.dept" :key="item.id">
                            {{ item.name }}
                    </el-tag>
               </div>

               <div class="flow-user-selected" v-show="formData.permissions.includes('role')">
               <span> 角色：</span>
               <el-tag  closable  @close="selectClose('role', item.id)" v-for="item in formData.combination.role" :key="item.id">
                           {{ item.roleName }}
                   </el-tag>
               </div>
               <div class="flow-user-selected" v-show="formData.permissions.includes('user')">
               <span> 用户：</span>
               <el-tag  closable  @close="selectClose('user', item.id)" v-for="item in formData.combination.user" :key="item.id">
                           {{ item.name }}
                   </el-tag>
               </div>
               <div v-show="formData.permissions.length == 0">
                   无
               </div>
           </el-form-item>

           <el-form-item label="协作方式" prop="nodeRatioType">
               <el-radio-group v-model="formData.nodeRatioType">
                   <el-radio  value="1">或签</el-radio>
                   <el-radio  value="2">票签</el-radio>
                   <el-radio  value="3">会签</el-radio>
               </el-radio-group>
           </el-form-item>
           <el-form-item label="票签占比（%）" prop="nodeRatio"   v-show="formData.nodeRatioType == '2'">
               <el-input-number
                   v-model="formData.nodeRatio"
                   class="mx-4"
                   :min="0"
                   :max="100"
                   controls-position="right"
               />
           </el-form-item>

           <el-form-item label="回退时的节点" prop="backType">
               <el-radio-group v-model="formData.backType">
                   <el-radio  value="1" @click="backToStartNode()">开始节点</el-radio>
                   <el-radio  value="2">选择节点</el-radio>
                   <!-- <el-radio  value="3" disabled>动态节点</el-radio> -->
               </el-radio-group>
           </el-form-item>
           <el-form-item label="请选择节点" prop="backTypeNode" v-if="formData.backType == '2'">
               <el-select v-model="formData.backTypeNode"   >
                   <el-option   v-for="item,index in backNodeList" :key="index"  :value="item.nodeId" :label="item.nodeName" />
               </el-select>
           </el-form-item>

           <el-form-item label="审批人为空时"  prop="emptyApprove.type">
               <el-radio-group v-model="formData.emptyApprove.type">
                   <el-radio  value="AUTO">自动通过</el-radio>
                   <el-radio  value="REJECT">自动拒绝</el-radio>
                   <el-radio  value="USER">
                       <span class="flow-user-select-checkbox">
                           指定人员 
                           <span style="width: 20px; " @click.prevent>
                               <el-icon style="float: right" :size="14" @click="addDefaultUser()" >
                                   <CirclePlus />
                               </el-icon>
                           </span>
                       </span>
                   </el-radio>
                   <el-radio  value="MANAGER">流程管理员</el-radio>
               </el-radio-group>
           </el-form-item>
           <el-form-item :label="formData.emptyApprove.value.length == 0 ? '请选择用户': '已选择用户' " v-if="formData.emptyApprove.type == 'user'" >
               <div class="flow-user-selected" >
               <el-tag  closable   v-for="item in formData.emptyApprove.value" :key="item.id" @close="emptyUserRemove(item.id)" >
                       {{ item.name }}
                   </el-tag>
               </div>
           </el-form-item>
       
           
       </el-tab-pane>
       <el-tab-pane label="操作权限" name="second">
           <div class="flow-operate-setting">
          
             
                   <el-row >
                       <template  v-for="(item, index) in formData.buttons" :key="index">
                           <el-col :span="6" >
                               <el-checkbox v-model="item.checked" :disabled="operateNames[item.type].disabled" :label="operateNames[item.type].label" />
                           </el-col>
                           <el-col :span="18">
                               <el-form-item required>
                                   <el-input v-model="item.text"   placeholder=""  maxlength="10" />
                               </el-form-item>
                           </el-col>
                       </template>
                   </el-row>
               
               </div>
       </el-tab-pane>
       <el-tab-pane label="表单配置" name="third">
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

       <el-tab-pane label="监听器" name="four">
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

 <RoleSelect ref="roleSelectRef" @roleSelectData="roleSelectData"></RoleSelect>
 <UserSelect ref="userSelectRef" @roleSelectData="userSelectData"></UserSelect>
 <UserSelect ref="emptySelectRef" @roleSelectData="emptySelectData"></UserSelect>
 <DeptSelect ref="deptSelectRef" @deptSelectData="deptSelectData"> </DeptSelect>
</template>

<script setup>
   import {defineProps, toRefs,watch, ref, onMounted, reactive} from 'vue';
   import { listApi } from '~/api/sys/user';
   import RoleSelect from '~/views/sys/role/roleSelect.vue';
   import UserSelect from '~/views/sys/user/userSelect.vue';
   import DeptSelect from '@/views/sys/dept/deptSelect.vue';
   import { nodeList } from '~/utils/flowUtil';
   const roleSelectRef = ref(null);
   const userSelectRef = ref(null);
   const emptySelectRef = ref(null);
   const deptSelectRef = ref(null);
   const props = defineProps({
       data:{},
       backNodeList:[],
   })

   // 不变
   const operateNames = ref({
       aggren : {label: "同意", disabled: true},
       reject : {label: "拒绝", disabled: true},
       back : {label: "回退", disabled: false},
       transfer : {label: "转办", disabled: false},
       depute : {label: "委派", disabled: false},
       signAdd : {label: "加签", disabled: false},
       signRedu : {label: "减签", disabled: false},
   })
   const formData = reactive({
       nodeRatioType: "1",
       backType: '1',
       nodeRatio: 0,
       permissions: [],
       combination:{
           role: [],
           user:[],
           dept:[],
       },
       emptyApprove:{
           type: 'AUTO',
           value: [],
       },
       value:'',
       formPath:'',
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
               text:"拒绝",
           },
           {
               type: 'back',
               checked: false,
               text:"回退",
           },
           {
               type: 'transfer',
               checked: false,
               text:"转办",
           },
           {
               type: 'depute',
               checked: false,
               text:"委派",
           },
          
           {
               type: 'signAdd',
               checked: false,
               text:"加签",
           },
           {
               type: 'signRedu',
               checked: false,
               text:"减签",
           }
           
           
       ]
   })

   watch(()=> formData.nodeRatioType, (newV) => {
       if(newV == '1'){
           formData.nodeRatio = 0;
       }else if(newV == '3'){
           formData.nodeRatio = 100;
       }else{
           formData.nodeRatio = 70;
       }
   })
   
   const showValue = () => {
       let value = '';
       if(formData.permissions.includes('role')){
           const roles = formData.combination.role;
           for(const i in roles){
               value  += roles[i].roleName + ",";
           }
       }
       if(formData.permissions.includes('user')){
           const users = formData.combination.user;
           for(const i in users){
               value  += users[i].name + ",";
           }
       }
       if(formData.permissions.includes('dept')){
           const depts = formData.combination.dept;
           for(const i in depts){
               value  += depts[i].name + ",";
           }
       }
       if( value.length > 0){
           value = value.slice(0, -1);
       }
       formData.value = value;
       return true;
   }

   /**
    * 点击加号按钮
    * @param type 
    */
   const addApproval = (type) => {
       if(formData.permissions.includes(type)){
           if(type === 'role'){
               roleSelectRef.value.open();
           // 选择角色
           }else if(type === 'user'){
               // 选择用户
               userSelectRef.value.open();
           }else if(type === 'dept'){
               // 选择用户
               deptSelectRef.value.open();
           }
       }
      

   }
   const addDefaultUser = () => {
       if(formData.emptyApprove.type === 'USER'){
           emptySelectRef.value.open();
       }
   }

   // 表格中选择的角色，转为流程中的角色
   const roleSelectData = (roleList) => {
       const roles =  formData.combination.role;
       for(const i in roleList){
           const role = roleList[i];
          if( !roles.find(r => r.id == role.id )){
           roles.push({
               id: role.id,
               roleName: role.roleName
           })
          }
       }
   }
   const deptSelectData = (deptList) => {
       const depts =  formData.combination.dept = formData.combination.dept || [];
       for(const i in deptList){
           const dept = deptList[i];
          if( !depts.find(r => r.id == dept.id )){
           depts.push({
               id: dept.id,
               name: dept.name
           })
          }
       }
   }

   const emptySelectData = (userList) => {
       const users =  formData.emptyApprove.value;
       for(const i in userList){
               const user = userList[i];
           if( !users.find(r => r.id == user.id )){
               users.push({
                   id: user.id,
                   name: user.nick
               })
           }
       }
   }

   const userSelectData = (userList, mutiSelect) => {
       // 多选时为审批设置
   
       const users =  formData.combination.user;
       for(const i in userList){
               const user = userList[i];
           if( !users.find(r => r.id == user.id )){
               users.push({
                   id: user.id,
                   name: user.nick
               })
           }
       }
     
       
   }
   
   // 关闭tag时， 移除审批人
   const selectClose = (type, id) =>{
       if(type == 'role'){
           const roles =  formData.combination.role;
           let index =  roles.findIndex(item => item.id === id);
           if(index >= 0){
               roles.splice(index, 1);
           }
       }else if(type == 'user'){
           const users =  formData.combination.user;
           let index =  users.findIndex(item => item.id === id);
           if(index >= 0){
               users.splice(index, 1);
           }
       }else if(type == 'dept'){
           const depts =  formData.combination.dept;
           let index =  depts.findIndex(item => item.id === id);
           if(index >= 0){
                depts.splice(index, 1);
           }
       }
      
   }
   // 审批人为空  指定人员移除
   const emptyUserRemove = (id) => {
       const users =  formData.emptyApprove.value;
       let index =  users.findIndex(item => item.id === id);
           if(index >= 0){
               users.splice(index, 1);
           }
   }

   const backToStartNode = () =>{
       formData.backTypeNode = backNodeList.value[0].nodeId;
   }

   const addCondition = (index) => {
        formData.listeners.splice(index, 0, {
            listenerType: '', //监听器类型
            listenerPath:'', // 监听器类路径
        })
    }
    
    const removeCondition = (index) => {
        formData.listeners.splice(index, 1);
    }

   const userList = ref([])
   onMounted(()=>{
       listApi({}).then(result => {
           userList.value = result.rows;
       })
   })
   const activeName = ref('first');
   
   const {data, backNodeList} = toRefs(props);
 
   Object.assign(formData, data.value.properties);
   const formConfig = async () => {
       return new Promise((resolve, reject) => {
           const result = showValue();
           if(result){
               if(formData.backTypeNode == null){
                   formData.backTypeNode = backNodeList.value[0].nodeId;
               }

               resolve(formData);
           }else {
               reject();
           }
       })
   }

   defineExpose({
       formConfig
   })

  
</script>

<style  scoped>
   .flow-user-select-checkbox{
       display: flex;
       flex-direction: row;
       flex-wrap: nowrap;
       align-items: center;
       gap: 5px;
   }
   .flow-user-selected{
       display: flex;
       align-content: center;
       align-items: center;
       flex-wrap: wrap;
       gap: 10px;
       width: 100%;
   
   }
   .flow-operate-setting{
       display: flex;
       flex-direction: column;
       gap: 10px;
   }
</style>