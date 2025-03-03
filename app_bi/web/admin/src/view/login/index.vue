<template>
    <el-row :gutter="20">
        <el-col :span="60" :offset="6">
            <el-form ref="loginForm" :model="from" :rules="loginRules" label-width="60px">
                <h3 class="title">用户登陆</h3>
                <el-form-item label="账号" prop="username">
                    <el-input v-model="from.username" placeholder="请输入用户名" />
                </el-form-item>
                <el-form-item label="密码" prop="password">
                    <el-input v-model="from.password" type="password" placeholder="请输入密码" />
                </el-form-item>
                <el-form-item label="">
                    <el-button type="primary" class="w100p" @click="Login">登录</el-button>
                </el-form-item>
            </el-form>
        </el-col>
    </el-row>
</template>

<script setup lang="ts">

defineOptions({
    name: "Login",
})

import { reactive, ref } from 'vue';
import type { LoginResponse } from '@/api/api';
import { ElForm } from 'element-plus';
import { useRouter } from 'vue-router';
import { useUserStore } from '@/pinia/modules/user'
const router = useRouter();

// 表单对象
const loginForm = ref<InstanceType<typeof ElForm> | null>(null);
// 表单数据
const from = reactive({
    username: '',
    password: ''
});
// 校验规则
const loginRules = {
    username: [
        { required: true, message: '请输入用户名', trigger: ['change', 'blur'] },
        { pattern: /^[a-zA-Z][a-zA-Z0-9_-]{3,20}$/, message: '以字母开头，可包含数字，_和-长度在 4 到 20 个字符', trigger: 'blur' },
        { min: 4, max: 20, message: '长度在 4 到 20 个字符', trigger: 'blur' }
    ],
    password: [
        { required: true, message: '请输入密码', trigger: 'blur' },
        { min: 6, max: 15, message: '长度在 6 到 15 个字符', }
    ]
}
// 使用"async"函数处理异步操作
async function Login() {
    if (!loginForm.value) {
        console.log('表单验证失败');
        return;
    }
    // 验证表单
    const valid = await loginForm.value.validate();
    if (!valid) {
        console.log('表单验证失败');
        return;
    };
    try {
        const userStore = useUserStore()
        const res: LoginResponse = await userStore.LoginIn(from)
        // 登录
        console.log(res.data);
        alert(res.msg);
        // 登录成功后处理跳转逻辑
        const previousRoute = localStorage.getItem('previousRoute');
        if (previousRoute) {
            // 如果有来源页面，则跳转到来源页面
            router.push(previousRoute);
            localStorage.removeItem('previousRoute'); // 清除记录
        } else {
            // 如果没有来源页面，则跳转到指定的布局页面（例如首页）
            router.push({ name: 'Home' });
        }
    } catch (error) {
        console.log(error);
        alert('登录失败');
    }
}
</script>

<style scoped>
.title {
    text-align: center;
    margin-bottom: 20px;
}

.w100p {
    width: 100%;
}

.txt-show {
    text-align: center;
    margin-top: 10px;
}
</style>