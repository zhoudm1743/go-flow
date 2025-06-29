
<template>
  <div class="flow-menu-header">
    <el-menu class="el-menu-demo" mode="horizontal" style="width: 50vw;">
      <el-menu-item index="1" style="width: 200px;">
          <el-image style="height: 40px;" src="/logo.png"  />
      </el-menu-item>
      <el-menu-item h="full" @click="toggleDark()"  >
        
        <el-icon v-if="isDark"><Moon /></el-icon>
        <el-icon v-else><Sunny /></el-icon>
          <!-- <i inline-flex i="dark:ep-moon ep-sunny" /> -->
      </el-menu-item>
    </el-menu>

  <div style=" display: flex; gap: 8px;  align-items: center;"  >
      <span style="font-size: 14px; font-family: 'Franklin Gothic Medium', 'Arial Narrow', Arial, sans-serif;">
         {{ userStore.baseInfo.nick }}
      </span>
      <div style="width: 80px;">
        <el-dropdown @command="handleCommand">
          <span class="el-dropdown-link">
            <el-avatar
            src="/avatar.png"
          />
          <el-icon class="el-icon--right"><arrow-down /></el-icon>
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="a" @click="$router.push({path: '/profile' })">个人中心</el-dropdown-item>
              <el-dropdown-item command="b"  @click="logout">退出登录</el-dropdown-item>
            </el-dropdown-menu>
          </template>
      </el-dropdown>
      </div>
  </div>
  
  </div>
</template>

<script setup>
import {toggleDark} from '~/composables'
import { useUserStore } from '~/store/user'
const userStore = useUserStore();
const logout = () => {
  userStore.logout()
}

</script>
<style scoped>
  .flow-menu-header{
    display: flex;
    justify-content: space-between;
  }
</style>