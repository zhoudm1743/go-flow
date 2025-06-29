<template>
    <div class="openflow-login">
        <div class="openflow-login-box">
            <div class="openflow-login-header">
                <span>开源工作流平</span><span>台</span>
            </div>
            <el-form class="login-form" label-position="top" :model="formdata" :rules="rules" ref="formRef">
                <el-tabs v-model="activeName" stretch>
                <el-tab-pane label="账号登录" name="first" >
                    <el-form-item  prop="username" required class="form-item">
                        <el-input v-model="formdata.username" placeholder="账号" ></el-input>
                    </el-form-item>
                    <el-form-item  prop="password" required class="form-item">
                        <el-input v-model="formdata.password" type="password" show-password placeholder="密码"></el-input>
                    </el-form-item>
                    <el-form-item  prop="captcha" required class="form-item">
                        <el-input v-model="formdata.captcha" placeholder="验证码" >
                            <template #append>
                                <el-image :src="captcha_img" style="width: 80px; height: 32px;" @click="getCaptchaData"/>
                            </template>
                        </el-input>
                    </el-form-item>
                </el-tab-pane>
                    <el-form-item class="form-item">
                        <el-button type="primary"   :loading="loading" @click="toLogin"  style="width: 100%">
                            登 录
                        </el-button>
                    </el-form-item>
                </el-tabs>
            </el-form>
        </div>
</div>
</template>

<script setup>
    import { ref, onMounted } from 'vue'   
    import { useUserStore } from '~/store/user'
    const userStore = useUserStore();
    import { useRoute, useRouter } from 'vue-router'
    import { usePermissionStore } from "~/store/permission.js";
    import { getCaptcha } from '~/api/auth/login'
    const permissionStore = usePermissionStore();
    const router = useRouter()
    const loading = ref(false);
    const activeName = ref('first')
    const formdata = ref({
        username: "admin",
        password: "12345678",
        captcha: "",
        app_id: ""
    })

    const captcha_img = ref("");
    const toLogin = () => {
        loading.value = true;
       
        userStore.loginByAccount(formdata.value).then((result) => {
            console.log(result)
            router.push({path: "/"})
        }).catch((e)=>{
            loading.value = false;
            // 重新获取验证码
            getCaptcha();
        });
      
    }

    const getCaptchaData = async () => {
        const res = await getCaptcha();
        formdata.value.app_id = res.data.app_id;
        captcha_img.value = res.data.captcha;
    }

    onMounted(() => {
        getCaptchaData();
    })
</script>

<style  scoped>
    .form-item{
        padding-top: 10px;
        padding-bottom: 5px;
    }
  
    .openflow-login{
        display: grid;
        place-items: center;
        height: 100vh; /* 父元素的高度 */
        background-image: url('/pg.png') ;
        background-repeat: no-repeat;
        background-size: cover;
    }
    .openflow-login-header{
        text-align: justify;
        text-align-last: justify;
        /* margin: 20px 0; */
        font-size: 20px;
        font-weight: 600;
        /* letter-spacing: 28px; */
        width: 400px;
    }
    .openflow-login-box{
        padding: 40px 60px;
        background: #ffffff;
    }
</style>
