<template>
  <GlassBackground />
  <RouterView />
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import GlassBackground from '@/components/layout/GlassBackground.vue'

const auth = useAuthStore()

onMounted(() => {
  auth.loadUser()
})
</script>
