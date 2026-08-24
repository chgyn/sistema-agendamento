<template>
  <div class="min-h-screen bg-slate-950 text-slate-100 font-sans antialiased">
    <RouterView />
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useAuthStore } from './stores/auth'

const authStore = useAuthStore()

onMounted(async () => {
  if (authStore.token) {
    await authStore.fetchMe()
  }
})
</script>
