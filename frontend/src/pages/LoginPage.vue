<template>
  <div class="flex w-full items-center justify-center px-4">
    <BaseCard class="w-full max-w-md">

      <h1 class="mb-2 text-4xl font-bold text-gray-900 dark:text-white">
        Wojiao
      </h1>

      <p class="mb-8 text-gray-500 dark:text-gray-400">
        С возвращением
      </p>

      <form
        class="space-y-5"
        @submit.prevent="login"
      >
        <BaseInput
          v-model="email"
          type="email"
          placeholder="Эл. почта"
        />

        <BaseInput
          v-model="password"
          type="password"
          placeholder="Пароль"
        />

        <p
          v-if="error"
          class="text-sm text-red-600 dark:text-red-400"
        >
          {{ error }}
        </p>

        <BaseButton type="submit">
          Войти
        </BaseButton>
      </form>

      <p class="mt-6 text-center text-sm text-gray-500 dark:text-gray-400">
        Нет аккаунта?
        <RouterLink
          to="/register"
          class="font-semibold text-[var(--color-primary)] hover:text-[var(--color-primary)]/80"
        >
          Зарегистрироваться
        </RouterLink>
      </p>

    </BaseCard>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { isAxiosError } from 'axios'
import { RouterLink, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { login as loginRequest } from '@/services'

import BaseButton from '@/components/base/BaseButton.vue'
import BaseCard from '@/components/base/BaseCard.vue'
import BaseInput from '@/components/base/BaseInput.vue'


const email = ref('')
const password = ref('')
const error = ref('')

const router = useRouter()
const authStore = useAuthStore()
async function login() {
  error.value = ''

  try {
    await loginRequest({
      email: email.value,
      password: password.value,
    })

    await authStore.loadUser()

    await router.push('/app')
  } catch (err) {
    error.value = isAxiosError(err)
      ? err.response?.data?.message ?? 'Не удалось войти'
      : 'Не удалось войти'
  }
}
</script>