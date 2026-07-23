<template>
  <div class="flex h-screen flex-col">

    <AppHeader />

    <div class="flex flex-1 overflow-hidden">

      <AppSidebar />

      <main class="flex-1 overflow-y-auto p-8">
        <RouterView />
      </main>

    </div>

  </div>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted } from 'vue'

import { useGamificationStore } from '@/stores/gamification'

import AppHeader from '@/components/layout/AppHeader.vue'
import AppSidebar from '@/components/layout/AppSidebar.vue'

const HEARTBEAT_INTERVAL_MS = 30000

const gamification = useGamificationStore()

let heartbeatInterval: ReturnType<typeof setInterval> | undefined

function sendHeartbeatIfVisible() {
  if (document.visibilityState === 'visible') {
    gamification.heartbeat()
  }
}

onMounted(() => {
  gamification.loadProgress()
  sendHeartbeatIfVisible()

  heartbeatInterval = setInterval(sendHeartbeatIfVisible, HEARTBEAT_INTERVAL_MS)
})

onUnmounted(() => {
  clearInterval(heartbeatInterval)
})
</script>
