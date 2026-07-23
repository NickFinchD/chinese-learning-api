<template>
  <header
    class="flex h-16 items-center justify-between border-b border-gray-200 bg-white px-6"
  >
    <div class="text-xl font-bold text-blue-600">
      Wojiao
    </div>

    <div class="flex items-center gap-4">
      <span class="text-sm font-medium text-gray-700">
        {{ auth.user?.username }}
      </span>

      <button
        class="rounded-lg bg-red-500 px-3 py-2 text-white transition hover:bg-red-600"
        @click="onLogout"
      >
        Выйти
      </button>
    </div>
  </header>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'

import { useAuthStore } from '@/stores/auth'
import { logout } from '@/services'

const auth = useAuthStore()
const router = useRouter()

async function onLogout() {
  try {
    await logout()

    auth.logout()

    await router.push('/login')
  } catch (error) {
    console.error(error)
  }
}
</script>