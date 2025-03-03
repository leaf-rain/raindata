import type { RouteRecordRaw, RouteRecordName, RouteMeta } from 'vue-router'

const viewModules = import.meta.glob('../view/**/*.vue')
const pluginModules = import.meta.glob('../plugin/**/*.vue')



export type AsyncRouter = RouteRecordRaw & {
    name: RouteRecordName;
    meta: RouteMeta;

    children?: AsyncRouter[];
    parent?: AsyncRouter | null;


    defaultMenu?: boolean;
    btns?: any;
    hidden?: boolean;
    keepAlive?: boolean;
}

export const asyncRouterHandle = (asyncRouter: AsyncRouter[]) => {
    asyncRouter.forEach(item => {
        const componetName = item.component?.name
        console.log("componetName:",componetName);
        if (componetName) {
            if (componetName.split('/')[0] === 'view') {
                // item.component = dynamicImport(viewModules, componetName)
            } else if (componetName.split('/')[0] === 'plugin') {
                // item.component = dynamicImport(pluginModules, componetName)
            }
        }
        if (item.children) {
            asyncRouterHandle(item.children)
        }
    });
}

function dynamicImport(
    dynamicViewsModules: any,
    component: string
) {
    const keys = Object.keys(dynamicViewsModules)
    const matchKeys = keys.filter((key) => {
        const k = key.replace('../', '')
        return k === component
    })
    const matchKey = matchKeys[0]

    return dynamicViewsModules[matchKey]
}
