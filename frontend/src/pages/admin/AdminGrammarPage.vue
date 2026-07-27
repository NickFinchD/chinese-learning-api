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
        Грамматика
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
    </AdminDataTable>

    <p
      v-if="store.error"
      class="mt-4 text-sm text-red-600 dark:text-red-400"
    >
      {{ store.error }}
    </p>

    <AdminFormModal
      :open="modalOpen"
      :title="editingItem ? 'Редактировать заметку' : 'Новая грамматическая заметка'"
      :saving="saving"
      :error="formError"
      @close="modalOpen = false"
      @submit="onSubmit"
    >
      <BaseInput
        v-model="form.title"
        placeholder="Заголовок"
        required
      />
      <BaseTextarea
        v-model="form.explanation"
        placeholder="Объяснение"
        :rows="4"
      />
      <BaseSelect
        v-model="form.hskLevel"
        :options="hskOptions"
      />

      <div class="grid grid-cols-3 gap-2">
        <BaseInput
          v-model="form.exampleHanzi"
          placeholder="Пример 1: ханьцзы"
        />
        <BaseInput
          v-model="form.examplePinyin"
          placeholder="Пиньинь"
        />
        <BaseInput
          v-model="form.exampleTranslation"
          placeholder="Перевод"
        />
      </div>

      <div class="grid grid-cols-3 gap-2">
        <BaseInput
          v-model="form.example2Hanzi"
          placeholder="Пример 2: ханьцзы"
        />
        <BaseInput
          v-model="form.example2Pinyin"
          placeholder="Пиньинь"
        />
        <BaseInput
          v-model="form.example2Translation"
          placeholder="Перевод"
        />
      </div>

      <div class="grid grid-cols-3 gap-2">
        <BaseInput
          v-model="form.example3Hanzi"
          placeholder="Пример 3: ханьцзы"
        />
        <BaseInput
          v-model="form.example3Pinyin"
          placeholder="Пиньинь"
        />
        <BaseInput
          v-model="form.example3Translation"
          placeholder="Перевод"
        />
      </div>
    </AdminFormModal>

    <ConfirmDialog
      :open="deletingItem !== null"
      title="Удалить заметку?"
      :message="deletingItem ? `«${deletingItem.title}». Это действие нельзя отменить.` : ''"
      :loading="deleting"
      @cancel="deletingItem = null"
      @confirm="onDelete"
    />
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { RouterLink } from 'vue-router'

import { useAdminGrammarStore } from '@/stores/admin/grammar'
import type { GrammarNotePayload } from '@/services/admin/grammar'
import type { GrammarNote } from '@/types/lesson'

import AdminDataTable from '@/components/admin/AdminDataTable.vue'
import type { AdminTableColumn } from '@/components/admin/AdminDataTable.vue'
import AdminFormModal from '@/components/admin/AdminFormModal.vue'
import ConfirmDialog from '@/components/admin/ConfirmDialog.vue'
import AppIcon from '@/components/base/AppIcon.vue'
import BaseInput from '@/components/base/BaseInput.vue'
import BaseSelect from '@/components/base/BaseSelect.vue'
import BaseTextarea from '@/components/base/BaseTextarea.vue'

const store = useAdminGrammarStore()

const columns: AdminTableColumn[] = [
  { key: 'title', label: 'Заголовок' },
  { key: 'explanation', label: 'Объяснение' },
  { key: 'hsk_level', label: 'HSK' },
]

const hskOptions = [1, 2, 3, 4, 5, 6].map(level => ({ value: level, label: `HSK ${level}` }))

const modalOpen = ref(false)
const saving = ref(false)
const formError = ref('')
const editingItem = ref<GrammarNote | null>(null)

const deletingItem = ref<GrammarNote | null>(null)
const deleting = ref(false)

const form = reactive({
  title: '',
  explanation: '',
  hskLevel: 1 as number,
  exampleHanzi: '',
  examplePinyin: '',
  exampleTranslation: '',
  example2Hanzi: '',
  example2Pinyin: '',
  example2Translation: '',
  example3Hanzi: '',
  example3Pinyin: '',
  example3Translation: '',
})

function resetForm() {
  form.title = ''
  form.explanation = ''
  form.hskLevel = 1
  form.exampleHanzi = ''
  form.examplePinyin = ''
  form.exampleTranslation = ''
  form.example2Hanzi = ''
  form.example2Pinyin = ''
  form.example2Translation = ''
  form.example3Hanzi = ''
  form.example3Pinyin = ''
  form.example3Translation = ''
}

function openCreate() {
  editingItem.value = null
  resetForm()
  formError.value = ''
  modalOpen.value = true
}

function openEdit(item: GrammarNote) {
  editingItem.value = item
  form.title = item.title
  form.explanation = item.explanation
  form.hskLevel = item.hsk_level
  form.exampleHanzi = item.example_hanzi ?? ''
  form.examplePinyin = item.example_pinyin ?? ''
  form.exampleTranslation = item.example_translation ?? ''
  form.example2Hanzi = item.example2_hanzi ?? ''
  form.example2Pinyin = item.example2_pinyin ?? ''
  form.example2Translation = item.example2_translation ?? ''
  form.example3Hanzi = item.example3_hanzi ?? ''
  form.example3Pinyin = item.example3_pinyin ?? ''
  form.example3Translation = item.example3_translation ?? ''
  formError.value = ''
  modalOpen.value = true
}

async function onSubmit() {
  saving.value = true
  formError.value = ''

  const payload: GrammarNotePayload = {
    title: form.title,
    explanation: form.explanation,
    hsk_level: form.hskLevel,
    example_hanzi: form.exampleHanzi,
    example_pinyin: form.examplePinyin,
    example_translation: form.exampleTranslation,
    example2_hanzi: form.example2Hanzi,
    example2_pinyin: form.example2Pinyin,
    example2_translation: form.example2Translation,
    example3_hanzi: form.example3Hanzi,
    example3_pinyin: form.example3Pinyin,
    example3_translation: form.example3Translation,
  }

  try {
    if (editingItem.value) {
      await store.update(editingItem.value.id, payload)
    } else {
      await store.create(payload)
    }

    modalOpen.value = false
  } catch {
    formError.value = 'Не удалось сохранить заметку'
  } finally {
    saving.value = false
  }
}

function confirmDelete(item: GrammarNote) {
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
