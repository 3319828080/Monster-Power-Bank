<template>
  <div class="app-shell" :class="{ 'has-tabbar': showTabBar }">
    <router-view v-slot="{ Component }">
      <keep-alive :include="['Home', 'OrderList', 'Profile']">
        <component :is="Component" />
      </keep-alive>
    </router-view>

    <van-tabbar v-model="activeTab" route placeholder v-if="showTabBar">
      <van-tabbar-item icon="home-o" to="/">首页</van-tabbar-item>
      <van-tabbar-item icon="orders-o" to="/orders">订单</van-tabbar-item>
      <van-tabbar-item icon="user-o" to="/profile">我的</van-tabbar-item>
    </van-tabbar>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const activeTab = ref(0)

const tabRoutes = ['/', '/orders', '/profile']
const showTabBar = computed(() => tabRoutes.includes(route.path))

watch(
  () => route.path,
  (path) => {
    const idx = tabRoutes.indexOf(path)
    if (idx >= 0) activeTab.value = idx
  },
  { immediate: true },
)
</script>

<style>
:root {
  --theme-color: #07c160;
  --theme-color-light: #06ad56;
}
body { margin: 0; }
.app-shell.has-tabbar { padding-bottom: 50px; }
</style>
