<template>
  <div>
    <h1 class="mb-8 text-3xl font-bold text-gray-900 dark:text-white">
      Мои подборки
    </h1>

    <form
      class="mb-8 flex max-w-md gap-2"
      @submit.prevent="onCreate"
    >
      <BaseInput
        v-model="newName"
        placeholder="Название новой подборки"
      />

      <button
        type="submit"
        class="shrink-0 rounded-xl bg-[var(--color-primary)] px-4 py-3 font-semibold text-white shadow-lg shadow-[var(--color-primary)]/30 transition hover:bg-[var(--color-primary)]/90 disabled:cursor-not-allowed disabled:opacity-50"
        :disabled="!newName.trim() || creating"
      >
        <AppIcon name="plus" />
      </button>
    </form>

    <div
      v-if="collections.loading"
      class="flex items-center gap-2 text-gray-500 dark:text-gray-400"
    >
      <BaseSpinner />
      Загрузка...
    </div>

    <div
      v-else-if="collections.items.length === 0"
      class="text-gray-500 dark:text-gray-400"
    >
      У вас пока нет подборок. Создайте первую выше — потом добавляйте в неё слова из словаря.
    </div>

    <div
      v-else
      class="grid gap-4 md:grid-cols-2 xl:grid-cols-3"
    >
      <RouterLink
        v-for="(collection, index) in collections.items"
        :key="collection.id"
        :to="`/app/vocabulary/collections/${collection.id}`"
        class="animate-fade-in-up group relative rounded-xl border border-white/50 bg-white/30 p-6 shadow-sm backdrop-blur-xl transition hover:-translate-y-0.5 hover:shadow-md dark:border-white/10 dark:bg-white/5"
        :style="{ animationDelay: `${Math.min(index * 40, 300)}ms` }"
      >
        <div class="mb-3 flex items-center gap-2 pr-14">
          <AppIcon
            name="folder"
            class="text-[var(--color-primary)]"
          />

          <span class="truncate text-lg font-semibold text-gray-900 dark:text-white">
            {{ collection.name }}
          </span>
        </div>

        <div class="text-gray-500 dark:text-gray-400">
          {{ collection.word_count }} {{ wordLabel(collection.word_count) }}
        </div>

        <div class="absolute right-4 top-4 flex gap-1 opacity-0 transition group-hover:opacity-100">
          <button
            type="button"
            title="Переименовать"
            class="rounded-lg p-1.5 text-gray-400 hover:bg-white/50 hover:text-gray-700 dark:text-gray-500 dark:hover:bg-white/10 dark:hover:text-gray-200"
            @click.prevent.stop="onRename(collection.id, collection.name)"
          >
            <AppIcon
              name="pencil"
              :size="16"
            />
          </button>

          <button
            type="button"
            title="Удалить"
            class="rounded-lg p-1.5 text-gray-400 hover:bg-red-50 hover:text-red-600 dark:text-gray-500 dark:hover:bg-red-500/10 dark:hover:text-red-400"
            @click.prevent.stop="onDelete(collection.id)"
          >
            <AppIcon
              name="trash"
              :size="16"
            />
          </button>
        </div>
      </RouterLink>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { RouterLink } from 'vue-router'

import { useCollectionsStore } from '@/stores/collections'

import AppIcon from '@/components/base/AppIcon.vue'
import BaseInput from '@/components/base/BaseInput.vue'
import BaseSpinner from '@/components/base/BaseSpinner.vue'

const collections = useCollectionsStore()

const newName = ref('')
const creating = ref(false)

function wordLabel(count: number): string {
  const mod10 = count % 10
  const mod100 = count % 100

  if (mod10 === 1 && mod100 !== 11) {
    return 'слово'
  }

  if ([2, 3, 4].includes(mod10) && ![12, 13, 14].includes(mod100)) {
    return 'слова'
  }

  return 'слов'
}

async function onCreate() {
  const name = newName.value.trim()

  if (!name) {
    return
  }

  creating.value = true

  try {
    await collections.create(name)
    newName.value = ''
  } finally {
    creating.value = false
  }
}

async function onRename(id: number, currentName: string) {
  const name = window.prompt('Новое название подборки', currentName)

  if (!name || !name.trim() || name.trim() === currentName) {
    return
  }

  await collections.rename(id, name.trim())
}

async function onDelete(id: number) {
  if (!window.confirm('Удалить подборку? Слова из словаря удалены не будут.')) {
    return
  }

  await collections.remove(id)
}

onMounted(() => {
  collections.loadCollections()
})
</script>
