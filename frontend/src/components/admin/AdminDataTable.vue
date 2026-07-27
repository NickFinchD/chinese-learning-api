<template>
  <div>
    <div class="mb-4 flex items-center gap-3">
      <div class="max-w-xs flex-1">
        <BaseInput
          :model-value="search"
          placeholder="Поиск..."
          @update:model-value="$emit('update:search', $event)"
        />
      </div>

      <button
        type="button"
        class="shrink-0 rounded-full bg-[var(--color-primary)] px-5 py-3 text-sm font-semibold text-white shadow-lg shadow-[var(--color-primary)]/30 transition hover:bg-[var(--color-primary)]/90"
        @click="$emit('create')"
      >
        Добавить
      </button>
    </div>

    <div
      v-if="loading"
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
            <th
              v-for="col in columns"
              :key="col.key"
              class="whitespace-nowrap px-4 py-3 font-semibold text-gray-700 dark:text-gray-300"
            >
              {{ col.label }}
            </th>
            <th class="px-4 py-3" />
          </tr>
        </thead>

        <tbody>
          <tr
            v-for="item in items"
            :key="item.id"
            class="border-b border-white/30 last:border-0 dark:border-white/5"
          >
            <td
              v-for="col in columns"
              :key="col.key"
              class="px-4 py-3 text-gray-800 dark:text-gray-200"
            >
              <slot
                :name="`cell-${col.key}`"
                :item="item"
              >
                {{ (item as Record<string, unknown>)[col.key] }}
              </slot>
            </td>

            <td class="px-4 py-3">
              <div class="flex justify-end gap-1">
                <button
                  type="button"
                  title="Редактировать"
                  class="flex h-9 w-9 items-center justify-center rounded-lg text-gray-600 transition hover:bg-white/60 dark:text-gray-300 dark:hover:bg-white/10"
                  @click="$emit('edit', item)"
                >
                  <AppIcon
                    name="pencil"
                    :size="16"
                  />
                </button>

                <button
                  type="button"
                  title="Удалить"
                  class="flex h-9 w-9 items-center justify-center rounded-lg text-red-500 transition hover:bg-red-500/10"
                  @click="$emit('delete', item)"
                >
                  <AppIcon
                    name="trash"
                    :size="16"
                  />
                </button>
              </div>
            </td>
          </tr>

          <tr v-if="items.length === 0">
            <td
              :colspan="columns.length + 1"
              class="px-4 py-10 text-center text-gray-500 dark:text-gray-400"
            >
              Ничего не найдено
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <div
      v-if="totalPages > 1"
      class="mt-4 flex items-center justify-center gap-4"
    >
      <button
        type="button"
        :disabled="page <= 1"
        class="flex h-9 w-9 items-center justify-center rounded-lg text-gray-600 transition hover:bg-white/60 disabled:opacity-30 dark:text-gray-300 dark:hover:bg-white/10"
        @click="$emit('update:page', page - 1)"
      >
        <AppIcon name="arrow-left" />
      </button>

      <span class="text-sm text-gray-600 dark:text-gray-400">{{ page }} / {{ totalPages }}</span>

      <button
        type="button"
        :disabled="page >= totalPages"
        class="flex h-9 w-9 items-center justify-center rounded-lg text-gray-600 transition hover:bg-white/60 disabled:opacity-30 dark:text-gray-300 dark:hover:bg-white/10"
        @click="$emit('update:page', page + 1)"
      >
        <AppIcon
          name="arrow-left"
          class="rotate-180"
        />
      </button>
    </div>
  </div>
</template>

<script setup lang="ts" generic="T extends { id: number }">
import BaseInput from '@/components/base/BaseInput.vue'
import BaseSpinner from '@/components/base/BaseSpinner.vue'
import AppIcon from '@/components/base/AppIcon.vue'

export interface AdminTableColumn {
  key: string
  label: string
}

defineProps<{
  items: T[]
  columns: AdminTableColumn[]
  loading: boolean
  search: string
  page: number
  totalPages: number
}>()

defineEmits<{
  (e: 'update:search', value: string): void
  (e: 'update:page', value: number): void
  (e: 'create'): void
  (e: 'edit', item: T): void
  (e: 'delete', item: T): void
}>()
</script>
