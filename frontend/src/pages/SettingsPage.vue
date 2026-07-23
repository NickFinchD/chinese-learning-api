<template>
  <div>
    <h1 class="mb-8 text-3xl font-bold">
      Настройки
    </h1>

    <BaseCard class="max-w-md">
      <h2 class="mb-4 text-xl font-semibold">
        Аккаунт
      </h2>

      <div class="mb-4">
        <div class="text-sm text-gray-500">
          Имя пользователя
        </div>
        <div class="text-gray-900">
          {{ auth.user?.username }}
        </div>
      </div>

      <div class="mb-6">
        <div class="text-sm text-gray-500">
          Эл. почта
        </div>
        <div class="text-gray-900">
          {{ auth.user?.email }}
        </div>
      </div>

      <BaseButton @click="onLogout">
        Выйти
      </BaseButton>
    </BaseCard>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'

import { useAuthStore } from '@/stores/auth'
import { logout } from '@/services'

import BaseCard from '@/components/base/BaseCard.vue'
import BaseButton from '@/components/base/BaseButton.vue'

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
