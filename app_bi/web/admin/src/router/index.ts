import { createRouter, createWebHashHistory } from 'vue-router'
import type { AsyncRouter } from '@/utils/asyncRouter'  // 修改为仅类型导入

const homeRoute: AsyncRouter = {
  name: 'Home',
  path: '/',
  meta: {
    closeTab: false,
  },
  component: () => import('@/view/layout/index.vue'),
};

const loginRoute: AsyncRouter = {
  name: 'Login',
  path: '/login',
  meta: {
    closeTab: false,
  },
  component: () => import('@/view/login/index.vue'),
};

const catchAllRoute: AsyncRouter = {
  name: 'Error',
  path: '/:catchAll(.*)',
  meta: {
    closeTab: true,
  },
  component: () => import('@/view/error/index.vue'),
};

const routes: AsyncRouter[] = [homeRoute, loginRoute, catchAllRoute];

const router = createRouter({
  history: createWebHashHistory(),
  routes,
});

export default router;