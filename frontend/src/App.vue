<script setup lang="ts">
import { RouterView } from 'vue-router'
import { Toaster } from 'sonner'
</script>

<template>
  <RouterView />
  <Toaster position="top-center" :duration="3000" />
</template>

<style scoped>
/* App.vue 样式已清理，现在只显示路由内容 */
</style>
