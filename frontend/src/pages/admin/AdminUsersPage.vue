<template>
  <div>
    <div class="mb-8 flex items-center gap-3">
      <RouterLink
        to="/app/admin"
        class="flex h-9 w-9 items-center justify-center rounded-lg text-gray-600 transition hover:bg-white/60 dark:text-gray-300 dark:hover:bg-white/10"
      >
        <AppIcon name="arrow-left" />
      </RouterLink>

      <h1 class="text-3xl font-bold text-gray-900 dark:text-white">
        Пользователи
      </h1>
    </div>

    <div class="mb-4 max-w-xs">
      <BaseInput
        :model-value="store.search"
        placeholder="Поиск..."
        @update:model-value="store.setSearch($event)"
      />
    </div>

    <div
      v-if="store.loading"
      class="flex justify-center py-12"
    >
      <BaseSpinner :size="28" />
    </div>

    <div
      v-else
      class="overflow-x-auto rounded-xl border border-white/50 bg-white/30 backdrop-blur-xl dark:border-white/10 dark:bg-white/5"
    >
      <table class="w-full text-left text-sm">
        <thead>
          <tr class="border-b border-white/50 dark:border-white/10">
            <th class="px-4 py-3 font-semibold text-gray-700 dark:text-gray-300">
              Имя пользователя
            </th>
            <th class="px-4 py-3 font-semibold text-gray-700 dark:text-gray-300">
              Эл. почта
            </th>
            <th class="px-4 py-3 font-semibold text-gray-700 dark:text-gray-300">
              Роль
            </th>
            <th class="px-4 py-3" />
          </tr>
        </thead>

        <tbody>
          <tr
            v-for="user in store.items"
            :key="user.id"
            class="border-b border-white/30 last:border-0 dark:border-white/5"
          >
            <td class="px-4 py-3 text-gray-800 dark:text-gray-200">
              {{ user.username }}
              <span
                v-if="user.id === auth.user?.id"
                class="ml-1 text-xs text-gray-500 dark:text-gray-400"
              >(вы)</span>
            </td>
            <td class="px-4 py-3 text-gray-800 dark:text-gray-200">
              {{ user.email }}
            </td>
            <td class="px-4 py-3">
              <span
                v-if="user.is_admin"
                class="inline-flex items-center gap-1 rounded-full bg-[var(--color-primary)]/15 px-2.5 py-1 text-xs font-medium text-[var(--color-primary)] dark:bg-[var(--color-primary)]/20 dark:text-[var(--color-mint)]"
              >
                <AppIcon
                  name="lock"
                  :size="12"
                />
                Админ
              </span>
              <span
                v-else
                class="text-xs text-gray-500 dark:text-gray-400"
              >
                Пользователь
              </span>
            </td>
            <td class="px-4 py-3 text-right">
              <button
                type="button"
                class="rounded-full border border-white/50 bg-white/40 px-4 py-1.5 text-xs font-semibold text-gray-700 transition hover:bg-white/60 disabled:cursor-not-allowed disabled:opacity-50 dark:border-white/10 dark:bg-white/5 dark:text-gray-300 dark:hover:bg-white/10"
                :disabled="togglingId === user.id"
                @click="toggle(user)"
              >
                {{ user.is_admin ? 'Забрать права' : 'Сделать админом' }}
              </button>
            </td>
          </tr>

          <tr v-if="store.items.length === 0">
            <td
              colspan="4"
              class="px-4 py-10 text-center text-gray-500 dark:text-gray-400"
            >
              Ничего не найдено
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <div
      v-if="store.totalPages > 1"
      class="mt-4 flex items-center justify-center gap-4"
    >
      <button
        type="button"
        :disabled="store.page <= 1"
        class="flex h-9 w-9 items-center justify-center rounded-lg text-gray-600 transition hover:bg-white/60 disabled:opacity-30 dark:text-gray-300 dark:hover:bg-white/10"
        @click="store.setPage(store.page - 1)"
      >
        <AppIcon name="arrow-left" />
      </button>

      <span class="text-sm text-gray-600 dark:text-gray-400">{{ store.page }} / {{ store.totalPages }}</span>

      <button
        type="button"
        :disabled="store.page >= store.totalPages"
        class="flex h-9 w-9 items-center justify-center rounded-lg text-gray-600 transition hover:bg-white/60 disabled:opacity-30 dark:text-gray-300 dark:hover:bg-white/10"
        @click="store.setPage(store.page + 1)"
      >
        <AppIcon
          name="arrow-left"
          class="rotate-180"
        />
      </button>
    </div>

    <p
      v-if="store.error"
      class="mt-4 text-sm text-red-600 dark:text-red-400"
    >
      {{ store.error }}
    </p>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { RouterLink } from 'vue-router'

import { useAdminUsersStore } from '@/stores/admin/users'
import type { AdminUser } from '@/services/admin/users'
import { useAuthStore } from '@/stores/auth'

import AppIcon from '@/components/base/AppIcon.vue'
import BaseInput from '@/components/base/BaseInput.vue'
import BaseSpinner from '@/components/base/BaseSpinner.vue'

const store = useAdminUsersStore()
const auth = useAuthStore()

const togglingId = ref<number | null>(null)

async function toggle(user: AdminUser) {
  togglingId.value = user.id

  try {
    await store.toggleAdmin(user.id, !user.is_admin)
  } finally {
    togglingId.value = null
  }
}

onMounted(() => {
  store.fetchPage()
})
</script>
