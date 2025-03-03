import { useUserStore } from '@/pinia/modules/user'
import { useRouterStore } from '@/pinia/modules/router'
import router from '@/router'
import Nprogress from 'nprogress'
import 'nprogress/nprogress.css'
import type { RouteLocationNormalized, NavigationGuardNext } from 'vue-router'
import  getPageTitle from '@/utils/page'
import  { convertToString} from '@/utils/stringFun'

Nprogress.configure({ showSpinner: false, speed: 500 })

const whiteList = ['Login', 'Home', 'Error']

const getRouter = async (userStore: any) => {
  const routerStore = useRouterStore();
  await routerStore.SetAsyncRouter();
  await userStore.GetUserInfo();
  const asyncRouters = routerStore.asyncRouters;
  
  asyncRouters.forEach(asyncRouter => {
    router.addRoute(asyncRouter);
  });
};

async function handleKeepAlive(to: RouteLocationNormalized) {
  if (to.matched.some(item => item.meta.keepAlive)) {
    if (to.matched && to.matched.length > 2) {
      for (let i = 1; i < to.matched.length; i++) {
        const element = to.matched[i - 1]
        if (element.name === 'layout') {
          to.matched.splice(i, 1)
          await handleKeepAlive(to)
        }
      }
    }
  }
}

router.beforeEach(async(to, from: RouteLocationNormalized) => {
  const routerStore = useRouterStore()
  Nprogress.start()
  const userStore = useUserStore()
  to.meta.matched = [...to.matched]
  handleKeepAlive(to)
  const token = userStore.token
  // 在白名单中的判断情况
  document.title = getPageTitle(to.meta, to)
  if(to.meta.client) {
    return true
  }
  const toName = convertToString(to.name)
  const fromName = convertToString(from.name)
  console.log("路由守卫调用: ",fromName, toName)
  console.log("token: ",token)
  if (whiteList.indexOf(convertToString(toName)) > -1) {
    if (token)  {
      if (!routerStore.asyncRouterFlag && whiteList.indexOf(fromName) < 0) {
        await getRouter(userStore)
      }
      // token 可以解析但是却是不存在的用户 id 或角色 id 会导致无限调用
      if (userStore.userInfo?.authority?.defaultRouter != null) {
        if (router.hasRoute(userStore.userInfo.authority.defaultRouter)) {
          return { name: userStore.userInfo.authority.defaultRouter }
        } else {
          return { name: 'Error' }
        }
      } else {
        // 强制退出账号
        userStore.ClearStorage()
        return {
          name: 'Login',
          query: {
            redirect: document.location.hash
          }
        }
      }
    } else {
      return true
    }
  } else {
    // 不在白名单中并且已经登录的时候
    if (token) {
      console.log(sessionStorage.getItem("needCloseAll"))
      if(sessionStorage.getItem("needToHome") === 'true') {
        sessionStorage.removeItem("needToHome")
        return { path: '/'}
      }
      // 添加flag防止多次获取动态路由和栈溢出
      if (!routerStore.asyncRouterFlag && whiteList.indexOf(fromName) < 0) {
        await getRouter(userStore)
        if (userStore.token) {
          if (router.hasRoute(userStore.userInfo.authority.defaultRouter)) {
            return { ...to, replace: true }
          } else {
            return { name: 'Error' }
          }
        } else {
          return {
            name: 'Login',
            query: { redirect: '' }
          }
        }
      } else {
        if (to.matched.length) {
          return true
        } else {
          return { path: '/layout/404' }
        }
      }
    }
    // 不在白名单中并且未登录的时候
    if (!token) {
      return {
        name: 'Login',
        query: {
          redirect: document.location.hash
        }
      }
    }
  }
})

const removeLoading = () => {
  const element = document.getElementById('loading-box');
  if (element) {
    element.remove();
  }
}

router.afterEach(() => {
  // 路由加载完成后关闭进度条
  document.getElementsByClassName('main-cont main-right')[0]?.scrollTo(0, 0)
  removeLoading()
  Nprogress.done()
})

router.onError(() => {
  // 路由发生错误后销毁进度条
  removeLoading()
  Nprogress.remove()
})

