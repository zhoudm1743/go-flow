<template>
    <el-card >
        <el-row :gutter="20">
            <h3>全局配置正在开发中!</h3>
        </el-row>
        <el-row :gutter="20">
            <h3>发起人员配置</h3>
        </el-row>
        <el-row :gutter="20" style="flex: 1 0 auto; gap:20px; margin-bottom: 10px;">
            <el-button type="primary" round icon="Plus" size="small" @click="addApproval('org')">机构</el-button>
            <el-button type="primary" round icon="Plus" size="small" @click="addApproval('position')">职位</el-button>
            <el-button type="primary" round icon="Plus" size="small" @click="addApproval('role')">角色</el-button>
            <el-button type="primary" round icon="Plus" size="small" @click="addApproval('user')">用户</el-button>
        </el-row>
        <el-row :gutter="20" > 
            <el-form-item label="已选择" >
               <div class="flow-user-selected" v-show="formData.combination.role.length > 0">
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
               <div v-show="isNoSelected">
                <span> 无 </span>
               </div>
           </el-form-item>
        </el-row>
        
        <el-row :gutter="20">
            <h3>通知配置</h3>
        </el-row>
        <el-row :gutter="20" style="flex: 1 0 auto; gap:20px;">
            <el-form-item label="开启代办通知:">
                <el-switch
                    v-model="formData.todoNotice"
                />
            </el-form-item>

            <el-form-item label="开启回退通知:">
                <el-switch
                    v-model="formData.backNotice"
                />
            </el-form-item>
            <el-form-item label="开启完成通知:">
                <el-switch
                    v-model="formData.completedNotice"
                />
            </el-form-item>

           
        </el-row>
    </el-card>

    <RoleSelect ref="roleSelectRef" @roleSelectData="roleSelectData"></RoleSelect>
</template>

<script setup>
    import { computed, reactive, ref, watch } from 'vue'
    import RoleSelect from '~/views/sys/role/roleSelect.vue';
    import UserSelect from '~/views/sys/user/userSelect.vue';
    
    const roleSelectRef = ref(null);
    const formData = reactive({
        todoNotice: true,
        backNotice: true,
        completedNotice: true,
        permissions:[],
        combination: {
            role: [],
            user: []
        }
    })

    const isNoSelected = computed(() => {
        if(formData.combination.role.length == 0 && formData.combination.user.length == 0){
            return true
        }
        return false;
    })
   

    const addApproval = (type) => {
            if(type === 'role'){
                roleSelectRef.value.open();
            }
    }

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
       }
      
   }
</script>

<style  scoped>
   .flow-user-selected{
       display: flex;
       align-content: center;
       align-items: center;
       flex-wrap: wrap;
       gap: 10px;
       width: 100%;
   
   }
</style>