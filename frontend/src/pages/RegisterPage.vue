<template>
  <div class="flex w-full items-center justify-center px-4">
    <BaseCard class="w-full max-w-md">

      <h1 class="mb-2 text-4xl font-bold text-gray-900 dark:text-white">
        Wojiao
      </h1>

      <p class="mb-8 text-gray-500 dark:text-gray-400">
        Создайте аккаунт
      </p>

      <form
        class="space-y-5"
        @submit.prevent="onRegister"
      >
        <BaseInput
          v-model="username"
          placeholder="Имя пользователя"
          required
          minlength="3"
          maxlength="50"
        />

        <BaseInput
          v-model="email"
          type="email"
          placeholder="Эл. почта"
          required
        />

        <BaseInput
          v-model="password"
          type="password"
          placeholder="Пароль (минимум 8 символов)"
          required
          minlength="8"
        />

        <p
          v-if="error"
          class="text-sm text-red-600 dark:text-red-400"
        >
          {{ error }}
        </p>

        <BaseButton
          type="submit"
          :disabled="submitting"
        >
          Зарегистрироваться
        </BaseButton>
      </form>

      <p class="mt-6 text-center text-sm text-gray-500 dark:text-gray-400">
        Уже есть аккаунт?
        <RouterLink
          to="/login"
          class="font-semibold text-[var(--color-primary)] hover:text-[var(--color-primary)]/80"
        >
          Войти
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
import { login as loginRequest, register as registerRequest } from '@/services'

import BaseButton from '@/components/base/BaseButton.vue'
import BaseCard from '@/components/base/BaseCard.vue'
import BaseInput from '@/components/base/BaseInput.vue'

const username = ref('')
const email = ref('')
const password = ref('')
const error = ref('')
const submitting = ref(false)

const router = useRouter()
const authStore = useAuthStore()

async function onRegister() {
  error.value = ''
  submitting.value = true

  try {
    await registerRequest({
      username: username.value,
      email: email.value,
      password: password.value,
    })

    await loginRequest({
      email: email.value,
      password: password.value,
    })

    await authStore.loadUser()

    await router.push('/app')
  } catch (err) {
    error.value = isAxiosError(err)
      ? err.response?.data?.message ?? 'Не удалось зарегистрироваться'
      : 'Не удалось зарегистрироваться'
  } finally {
    submitting.value = false
  }
}
</script>
