import { asyncRouterHandle } from '@/utils/asyncRouter'
import { emitter } from '@/utils/bus.js'
import { asyncMenu } from '@/api/menu'
import { defineStore } from 'pinia'
import { ref, watchEffect, type Ref } from 'vue'
import pathInfo from "@/pathInfo.json";
import type { AsyncRouter } from '@/utils/asyncRouter'
import type { RouteRecordName } from 'vue-router'


const notLayoutRouterArr: AsyncRouter[] = []
const keepAliveRoutersArr: any[] = []
const nameMap = new Map<RouteRecordName, any>()

const formatRouter = (routes: AsyncRouter[], routeMap: { [x: RouteRecordName]: AsyncRouter }, parent: AsyncRouter|null) => {
    routes && routes.forEach(item => {
        item.parent = parent
        if (!item.meta) {
            item.meta = {};
        }
        item.meta.btns = item.btns
        item.meta.hidden = item.hidden
        if (item.meta.defaultMenu === true) {
            if (!parent) {
                item = { ...item, path: `/${item.path}` }
                notLayoutRouterArr.push(item)
            }
        }
        routeMap[item.name] = item
        if (item.children && item.children.length > 0) {
            formatRouter(item.children, routeMap, item)
        }
    })
}

const KeepAliveFilter = (routes: AsyncRouter[]) => {
    routes && routes.forEach(item => {
        // 子菜单中有 keep-alive 的，父菜单也必须 keep-alive，否则无效。这里将子菜单中有 keep-alive 的父菜单也加入。
        if ((item.children && item.children.some(ch => ch.meta?.keepAlive) || item.meta.keepAlive)) {
            const regex = /\(\) => import\("([^?]+)\??.*"\)/;
            const match = String(item.component).match(regex);
            const path = match ? match[1] : "";
            // 检查 path 是否是 pathInfo 的有效键
            if (path && (path in pathInfo)) {
                keepAliveRoutersArr.push(pathInfo[path as keyof typeof pathInfo]);
                nameMap.set(item.name, pathInfo[path as keyof typeof pathInfo]);
            }
        }
        if (item.children && item.children.length > 0) {
            KeepAliveFilter(item.children)
        }
    })
}

export const useRouterStore = defineStore('router', () :{
    topActive: Ref<string>;
    setLeftMenu: (name: string) => AsyncRouter[] | undefined;
    topMenu: Ref<AsyncRouter[]>;
    leftMenu: Ref<AsyncRouter[]>;
    asyncRouters: Ref<AsyncRouter[]>;
    keepAliveRouters: Ref<string[]>;
    asyncRouterFlag: Ref<number>;
    SetAsyncRouter: () => Promise<boolean>;
    routeMap: Record<string, AsyncRouter>;
} => {
    const keepAliveRouters = ref<string[]>([])
    const asyncRouterFlag = ref(0)

    const setKeepAliveRouters = (history: any[]) => {
        const keepArrTemp: string[] = []
        history.forEach(item => {
            if (nameMap.get(item.name)) {
                keepArrTemp.push(nameMap.get(item.name))
            }
        })
        keepAliveRouters.value = Array.from(new Set(keepArrTemp))
    }

    emitter.on('setKeepAlive', (event: any) => setKeepAliveRouters)

    const asyncRouters = ref<AsyncRouter[]>([])

    const topMenu = ref<AsyncRouter[]>([])

    const leftMenu = ref<AsyncRouter[]>([])

    const menuMap: Record<RouteRecordName, AsyncRouter> = {};

    const topActive = ref("")





    const setLeftMenu = (name:string) => {
        sessionStorage.setItem('topActive', name)
        topActive.value = name
        if (menuMap[name]?.children) {
            leftMenu.value = menuMap[name].children
        }
        return menuMap[name]?.children
    }

    watchEffect(() => {
        let topActive = sessionStorage.getItem("topActive") || ''
        let firstHasChildren:string|symbol = '';
        
        if (asyncRouters.value.length > 0 && asyncRouters.value[0]?.children) {
            asyncRouters.value[0].children.forEach((item: AsyncRouter) => {
                if (item.hidden) return;
                menuMap[item.name] = item;
                if (!firstHasChildren && item.children && item.children.length > 0) {
                    firstHasChildren = item.name;
                }
                // 确保 children 是一个空数组
                topMenu.value.push({ ...item, children: [] });
            });
        }
    
        if (!menuMap[topActive]?.children && firstHasChildren) {
            topActive = firstHasChildren;
        }
        setLeftMenu(topActive);
    })

    const routeMap = ({})
    // 从后台获取动态路由
    const SetAsyncRouter = async () => {
        asyncRouterFlag.value++
        const baseRouter: AsyncRouter[] = [{
            path: '/layout',
            name: 'layout',
            component: () => import('@/view/layout/index.vue'),
            meta: {
                title: '底层layout',
                btns: undefined,
                hidden: false,
                defaultMenu: false,
                keepAlive: false
            },
            redirect: "",

            children: [],
            parent: null,
            btns: undefined,
            hidden: false,
            defaultMenu: false,
            keepAlive: false
        }]
        const asyncRouterRes = await asyncMenu()
        const asyncRouter = asyncRouterRes.data.menus
        asyncRouter && asyncRouter.push({
            path: 'reload',
            name: 'Reload',
            hidden: true,
            meta: {
                title: '',
                closeTab: true,
            },
            component: 'view/error/reload.vue'
        })
        formatRouter(asyncRouter, routeMap, null)
        baseRouter[0].children = asyncRouter
        if (notLayoutRouterArr.length !== 0) {
            baseRouter.push(...notLayoutRouterArr)
        }
        asyncRouterHandle(baseRouter)
        KeepAliveFilter(asyncRouter)
        asyncRouters.value = baseRouter
        return true
    }

    return {
        topActive,
        setLeftMenu,
        topMenu,
        leftMenu,
        asyncRouters,
        keepAliveRouters,
        asyncRouterFlag,
        SetAsyncRouter,
        routeMap
    }
})

