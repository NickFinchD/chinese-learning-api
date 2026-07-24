<template>
  <div class="flex h-screen flex-col">

    <AppHeader @toggle-sidebar="sidebarOpen = !sidebarOpen" />

    <div class="flex flex-1 overflow-hidden">

      <AppSidebar
        :open="sidebarOpen"
        @close="sidebarOpen = false"
      />

      <main class="flex-1 overflow-y-auto p-4 sm:p-6 md:p-8">
        <RouterView v-slot="{ Component, route }">
          <Transition
            name="page"
            mode="out-in"
          >
            <component
              :is="Component"
              :key="route.fullPath"
            />
          </Transition>
        </RouterView>
      </main>

    </div>

  </div>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'

import { useGamificationStore } from '@/stores/gamification'

import AppHeader from '@/components/layout/AppHeader.vue'
import AppSidebar from '@/components/layout/AppSidebar.vue'

const HEARTBEAT_INTERVAL_MS = 30000

const sidebarOpen = ref(false)

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
