import { login, getUserInfo, setSelfInfo } from '@/api/user'
import { jsonInBlacklist } from '@/api/jwt'
import router from '@/router/index'
import { ElLoading, ElMessage } from 'element-plus'
import { defineStore } from 'pinia'
import { ref, computed, watch } from 'vue'
import { useRouterStore } from './router'
import cookie from 'js-cookie'
import type { AsyncRouter } from '@/utils/asyncRouter'
import type { LoginResponse } from "@/api/api.ts";

export interface UserInfo {
    uuid: string,
    nickName: string,
    headerImg: string,
    authority: any,
    sideMode: string | null,
    baseColor: string,

    value: any,
    authorityId: number
}

export const useUserStore = defineStore('user', () => {
    const loadingInstance = ref<any>(null)

    const userInfo = ref({
        uuid: '',
        nickName: '',
        headerImg: '',
        authority: {},
        sideMode: 'dark',
        baseColor: '#fff'
    } as UserInfo)
    const token = ref(window.localStorage.getItem('token') || cookie.get('x-token') || '')
    const setUserInfo = (val: any) => {
        userInfo.value = val
    }

    const setToken = (val: any) => {
        console.log('setToken', val)
        token.value = val
    }

    const NeedInit = () => {
        token.value = ''
        window.localStorage.removeItem('token')
        router.push({ name: 'Init', replace: true })
    }

    const ResetUserInfo = (value = {}) => {
        userInfo.value = {
            ...userInfo.value,
            ...value
        }
    }
    /* 获取用户信息*/
    const GetUserInfo = async () => {
        const res = await getUserInfo()
        if (res.data.code === 0) {
            setUserInfo(res.data.userInfo)
        }
        return res
    }
    /* 登录*/
    const LoginIn = async (loginInfo: any) => {
        loadingInstance.value = ElLoading.service({
            fullscreen: true,
            text: '登录中，请稍候...',
        })

        // const res = await login(loginInfo)

        const res = {
            data: {
                code: 0,
                msg: 'success',
                data: {
                    token: '1234567890',
                    user: {
                        uuid: 'uuid',
                        nickName: 'rain',
                        headerImg: 'test.img',
                        authority: {},
                        sideMode: 'dark',
                        baseColor: '#fff'
                    } as UserInfo
                }
            } as LoginResponse
        }
        console.log('111111111111111111111111111111111111111111')

        // 登陆失败，直接返回
        if (res.data.code !== 0) {
            loadingInstance.value.close()
            return res
        }

        // 登陆成功，设置用户信息和权限相关信息
        setUserInfo(res.data.data.user)
        setToken(res.data.data.token)
        // 初始化路由信息
        const routerStore = useRouterStore()
        await routerStore.SetAsyncRouter()
        const asyncRouters = routerStore.asyncRouters as AsyncRouter[]

        // 注册到路由表里
        asyncRouters.forEach((asyncRouter: AsyncRouter) => {
            router.addRoute(asyncRouter)
        })

        if (!router.hasRoute(userInfo.value.authority.defaultRouter)) {
            ElMessage.error('请联系管理员进行授权')
        } else {
            await router.replace({ name: userInfo.value.authority.defaultRouter })
        }

        // 全部操作均结束，关闭loading并返回
        loadingInstance.value.close()
        return res
    }
    /* 登出*/
    const LoginOut = async () => {
        const res = await jsonInBlacklist()

        // 登出失败
        if (res.data.code !== 0) {
            return
        }

        await ClearStorage()

        // 把路由定向到登录页，无需等待直接reload
        router.push({ name: 'Login', replace: true })
        window.location.reload()
    }
    /* 清理数据 */
    const ClearStorage = async () => {
        token.value = ''
        sessionStorage.clear()
        window.localStorage.removeItem('token')
        cookie.remove('x-token')
    }
    /* 设置侧边栏模式*/
    const changeSideMode = async (data: any) => {
        const res = await setSelfInfo({ sideMode: data })
        if (res.data.code === 0) {
            userInfo.value.sideMode = data
            ElMessage({
                type: 'success',
                message: '设置成功'
            })
        }
    }

    const mode = computed(() => userInfo.value.sideMode)
    const sideMode = computed(() => {
        if (userInfo.value.sideMode === 'dark') {
            return '#191a23'
        } else if (userInfo.value.sideMode === 'light') {
            return '#fff'
        } else {
            return userInfo.value.sideMode
        }
    })
    const baseColor = computed(() => {
        if (userInfo.value.sideMode === 'dark') {
            return '#fff'
        } else if (userInfo.value.sideMode === 'light') {
            return '#191a23'
        } else {
            return userInfo.value.baseColor
        }
    })

    watch(() => token.value, () => {
        window.localStorage.setItem('token', token.value)
    })

    return {
        userInfo,
        token,
        NeedInit,
        ResetUserInfo,
        GetUserInfo,
        LoginIn,
        LoginOut,
        changeSideMode,
        mode,
        sideMode,
        setToken,
        baseColor,
        loadingInstance,
        ClearStorage
    }
})
