<template>
  <div class="flex min-h-screen items-center justify-center bg-gray-100 px-4">
    <BaseCard class="w-full max-w-md">

      <h1 class="mb-2 text-4xl font-bold text-gray-900">
        Chinese Learning
      </h1>

      <p class="mb-8 text-gray-500">
        Welcome back 👋
      </p>

      <form
        class="space-y-5"
        @submit.prevent="login"
      >
        <BaseInput
          v-model="email"
          type="email"
          placeholder="Email"
        />

        <BaseInput
          v-model="password"
          type="password"
          placeholder="Password"
        />

        <BaseButton type="submit">
          Login
        </BaseButton>
      </form>

    </BaseCard>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

import BaseButton from '@/components/base/BaseButton.vue'
import BaseCard from '@/components/base/BaseCard.vue'
import BaseInput from '@/components/base/BaseInput.vue'

const email = ref('')
const password = ref('')

function login() {
  console.log({
    email: email.value,
    password: password.value,
  })
}
</script>