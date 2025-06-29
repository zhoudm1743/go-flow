<template>
        <el-row style="gap: 10px; ">
            <el-card  style="max-width: 400px;">
                <template #header>
                    <div class="card-header">
                        <span>个人信息</span>
                    </div>
                </template>

                <el-row :gutter="20">
                    <el-col :span="24"  style="text-align: center;">
                        <el-avatar :size="100" src="/avatar.png" />
                    </el-col>
                    <el-divider border-style="solid" />
                    <el-col :offset="2" :span="6">账号：</el-col>
                    <el-col :span="16">{{ userInfo.username }}</el-col>
                    <el-divider border-style="solid" />
                    <el-col :offset="2" :span="6">昵称：</el-col>
                    <el-col :span="16">{{ userInfo.nick }}</el-col>
                    <el-divider border-style="solid" />
                
                    <el-col :offset="2" :span="6">邮箱：</el-col>
                    <el-col :span="16">{{ userInfo.email }}</el-col>
                    <el-divider border-style="solid" />
                    <el-col>
                        <p>任职信息</p> 
                    </el-col>
                    <el-col :offset="2">
                    
                        <template v-for="(item,index) in userPositions" :key="item.id">
                                <el-tag  type="danger"  effect="dark">{{ item.cmpName }}</el-tag> |
                                <template v-if="item.orgName" >
                                    <el-tag  type="primary"   effect="dark">{{ item.orgName }}</el-tag>   |
                                </template>
                                <el-tag  type="success"  effect="dark">{{item.positionName}}</el-tag>
                                <br  v-if="userPositions.length - 1 != index" />
                            </template>
                       
                    </el-col>
                </el-row>
             
            </el-card>
       
            <el-card  style="flex: 1 1 auto; text-align: center; align-items: center;">
                <p style="padding-top: 20px; line-height: 40px;">
                    <span >千里之行始于足下</span>
                    <br />
                    <span>不忘初心,方得始终</span>
                    <br />

                    <span>你想在这里添加什么功能呢</span>
                </p>
            </el-card>
        </el-row>
       
</template>

<script setup>
    import { useUserStore } from "@/store/user";
    const userStore = useUserStore();
    const userInfo =  userStore.baseInfo;
    const userPositions = userStore.userPositions;

</script>

<style  scoped>
    :deep .el-divider--horizontal{
        margin: 10px 0;
    }
</style>