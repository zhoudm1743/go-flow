<template>
   <div class="flow-operate-setting">
        <el-row >
            <el-col :span="12"><el-checkbox v-model="formData.button.aggren" disabled label="同意" /></el-col>
            <el-col :span="12">{{ formData.button.aggrenText  }}</el-col>
        </el-row>
        <el-row>
            <el-col :span="12"><el-checkbox v-model="formData.button.reject" disabled label="拒绝" /></el-col>
            <el-col :span="12">{{ formData.button.rejectText  }}</el-col>
        </el-row>
        <el-row >
            <el-col :span="12"><el-checkbox v-model="formData.button.back" disabled label="同意" /></el-col>
            <el-col :span="12">{{ formData.button.backText  }}</el-col>
        </el-row>
        <el-row>
            <el-col :span="12"><el-checkbox v-model="formData.button.skip" disabled label="跳转" /></el-col>
            <el-col :span="12">{{ formData.button.skipText  }}</el-col>
        </el-row>
        <el-row >
            <el-col :span="12"><el-checkbox v-model="formData.button.back" disabled label="回退" /></el-col>
            <el-col :span="12">{{ formData.button.backText  }}</el-col>
        </el-row>
        <el-row>
            <el-col :span="12"><el-checkbox v-model="formData.button.signature" disabled label="加签" /></el-col>
            <el-col :span="12">{{ formData.button.signatureText  }}</el-col>
        </el-row>
        <el-row>
            <el-col :span="12"><el-checkbox v-model="formData.button.transfer" disabled label="跳转" /></el-col>
            <el-col :span="12">{{ formData.button.transferText  }}</el-col>
        </el-row>
    </div>
</template>

<script setup>
    
</script>

<style lang="scss" scoped>

</style>