<template>
  <div>
    <h1>layout</h1>
    <h1>{{ tmp }}</h1>
    <router-link to="/test/test1">test1</router-link>
    <router-link to="/test/test2">test2</router-link>
    <hr>
    <router-view></router-view>
  </div>
</template>

<script setup lang="ts">
import { useAppStore } from "@/pinia";
import { onMounted } from 'vue';
const userStore = useAppStore();
onMounted(() => {
  console.log('userStore:', userStore); // 调试信息
});
const tmp = import.meta.env.DEV;

import { useRouter } from 'vue-router';
const router = useRouter();
router.addRoute({
  path: '/test/test1',
  name: 'test1',
  component: () => import('./test1.vue')
},);
router.addRoute({
  path: '/test/test2',
  name: 'test2',
  component: () => import('./test2.vue')
},);

</script>