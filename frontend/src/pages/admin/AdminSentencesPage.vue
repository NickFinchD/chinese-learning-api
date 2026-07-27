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
        Предложения
      </h1>
    </div>

    <AdminDataTable
      :items="store.items"
      :columns="columns"
      :loading="store.loading"
      :search="store.search"
      :page="store.page"
      :total-pages="store.totalPages"
      @update:search="store.setSearch($event)"
      @update:page="store.setPage($event)"
      @create="openCreate"
      @edit="openEdit"
      @delete="confirmDelete"
    >
      <template #cell-hsk_level="{ item }">
        HSK {{ item.hsk_level }}
      </template>
      <template #cell-chunks="{ item }">
        {{ item.chunks.join(' / ') }}
      </template>
    </AdminDataTable>

    <p
      v-if="store.error"
      class="mt-4 text-sm text-red-600 dark:text-red-400"
    >
      {{ store.error }}
    </p>

    <AdminFormModal
      :open="modalOpen"
      :title="editingItem ? 'Редактировать упражнение' : 'Новое упражнение'"
      :saving="saving"
      :error="formError"
      @close="modalOpen = false"
      @submit="onSubmit"
    >
      <BaseInput
        v-model="form.translation"
        placeholder="Перевод"
        required
      />
      <BaseInput
        v-model="form.pinyin"
        placeholder="Пиньинь (полное предложение)"
      />
      <BaseSelect
        v-model="form.hskLevel"
        :options="hskOptions"
      />

      <div>
        <div class="mb-2 text-sm font-medium text-gray-700 dark:text-gray-300">
          Части предложения (по порядку)
        </div>

        <AdminRepeatableList
          v-model="form.chunks"
          add-label="Добавить часть"
          :min-items="2"
          :new-item="() => ''"
        >
          <template #default="{ item, update }">
            <BaseInput
              :model-value="item"
              placeholder="Часть предложения"
              @update:model-value="update($event)"
            />
          </template>
        </AdminRepeatableList>
      </div>
    </AdminFormModal>

    <ConfirmDialog
      :open="deletingItem !== null"
      title="Удалить упражнение?"
      :message="deletingItem ? `«${deletingItem.translation}». Это действие нельзя отменить.` : ''"
      :loading="deleting"
      @cancel="deletingItem = null"
      @confirm="onDelete"
    />
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { RouterLink } from 'vue-router'

import { useAdminSentencesStore } from '@/stores/admin/sentences'
import type { SentenceExercisePayload } from '@/services/admin/sentences'
import type { SentenceExercise } from '@/types/lesson'

import AdminDataTable from '@/components/admin/AdminDataTable.vue'
import type { AdminTableColumn } from '@/components/admin/AdminDataTable.vue'
import AdminFormModal from '@/components/admin/AdminFormModal.vue'
import AdminRepeatableList from '@/components/admin/AdminRepeatableList.vue'
import ConfirmDialog from '@/components/admin/ConfirmDialog.vue'
import AppIcon from '@/components/base/AppIcon.vue'
import BaseInput from '@/components/base/BaseInput.vue'
import BaseSelect from '@/components/base/BaseSelect.vue'

const store = useAdminSentencesStore()

const columns: AdminTableColumn[] = [
  { key: 'translation', label: 'Перевод' },
  { key: 'chunks', label: 'Части' },
  { key: 'hsk_level', label: 'HSK' },
]

const hskOptions = [1, 2, 3, 4, 5, 6].map(level => ({ value: level, label: `HSK ${level}` }))

const modalOpen = ref(false)
const saving = ref(false)
const formError = ref('')
const editingItem = ref<SentenceExercise | null>(null)

const deletingItem = ref<SentenceExercise | null>(null)
const deleting = ref(false)

const form = reactive({
  translation: '',
  pinyin: '',
  hskLevel: 1 as number,
  chunks: ['', ''] as string[],
})

function resetForm() {
  form.translation = ''
  form.pinyin = ''
  form.hskLevel = 1
  form.chunks = ['', '']
}

function openCreate() {
  editingItem.value = null
  resetForm()
  formError.value = ''
  modalOpen.value = true
}

function openEdit(item: SentenceExercise) {
  editingItem.value = item
  form.translation = item.translation
  form.pinyin = item.pinyin
  form.hskLevel = item.hsk_level
  form.chunks = [...item.chunks]
  formError.value = ''
  modalOpen.value = true
}

async function onSubmit() {
  saving.value = true
  formError.value = ''

  const payload: SentenceExercisePayload = {
    translation: form.translation,
    pinyin: form.pinyin,
    hsk_level: form.hskLevel,
    chunks: form.chunks,
  }

  try {
    if (editingItem.value) {
      await store.update(editingItem.value.id, payload)
    } else {
      await store.create(payload)
    }

    modalOpen.value = false
  } catch {
    formError.value = 'Не удалось сохранить упражнение'
  } finally {
    saving.value = false
  }
}

function confirmDelete(item: SentenceExercise) {
  deletingItem.value = item
}

async function onDelete() {
  if (!deletingItem.value) return

  deleting.value = true

  try {
    await store.remove(deletingItem.value.id)
    deletingItem.value = null
  } finally {
    deleting.value = false
  }
}

onMounted(() => {
  store.fetchPage()
})
</script>
