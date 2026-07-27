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
        Тесты
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
      <template #cell-options="{ item }">
        {{ item.options.length }} вариант(ов)
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
      :title="editingItem ? 'Редактировать тест' : 'Новый тест'"
      :saving="saving"
      :error="formError"
      @close="modalOpen = false"
      @submit="onSubmit"
    >
      <BaseInput
        v-model="form.question"
        placeholder="Вопрос"
        required
      />

      <div class="grid grid-cols-2 gap-2">
        <BaseSelect
          v-model="form.hskLevel"
          :options="hskOptions"
        />
        <BaseSelect
          v-model="form.direction"
          :options="directionOptions"
        />
      </div>

      <div class="grid grid-cols-2 gap-2">
        <BaseInput
          v-model="form.hanzi"
          placeholder="Ханьцзы (если применимо)"
        />
        <BaseInput
          v-model="form.pinyin"
          placeholder="Пиньинь"
        />
      </div>

      <div>
        <div class="mb-2 text-sm font-medium text-gray-700 dark:text-gray-300">
          Варианты ответа
        </div>

        <AdminRepeatableList
          v-model="form.options"
          add-label="Добавить вариант"
          :min-items="2"
          :new-item="() => ({ text: '', pinyin: '', is_correct: false })"
        >
          <template #default="{ item, update }">
            <div class="flex items-center gap-2">
              <input
                type="checkbox"
                :checked="item.is_correct"
                title="Правильный вариант"
                class="h-5 w-5 shrink-0 accent-[var(--color-primary)]"
                @change="update({ ...item, is_correct: ($event.target as HTMLInputElement).checked })"
              >
              <BaseInput
                :model-value="item.text"
                placeholder="Текст варианта"
                @update:model-value="update({ ...item, text: $event })"
              />
              <BaseInput
                :model-value="item.pinyin"
                placeholder="Пиньинь"
                @update:model-value="update({ ...item, pinyin: $event })"
              />
            </div>
          </template>
        </AdminRepeatableList>
      </div>
    </AdminFormModal>

    <ConfirmDialog
      :open="deletingItem !== null"
      title="Удалить тест?"
      :message="deletingItem ? `«${deletingItem.question}». Это действие нельзя отменить.` : ''"
      :loading="deleting"
      @cancel="deletingItem = null"
      @confirm="onDelete"
    />
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { RouterLink } from 'vue-router'

import { useAdminQuizzesStore } from '@/stores/admin/quizzes'
import type { AdminQuiz, AdminQuizOption, QuizPayload } from '@/services/admin/quizzes'
import type { QuizDirection } from '@/types/lesson'

import AdminDataTable from '@/components/admin/AdminDataTable.vue'
import type { AdminTableColumn } from '@/components/admin/AdminDataTable.vue'
import AdminFormModal from '@/components/admin/AdminFormModal.vue'
import AdminRepeatableList from '@/components/admin/AdminRepeatableList.vue'
import ConfirmDialog from '@/components/admin/ConfirmDialog.vue'
import AppIcon from '@/components/base/AppIcon.vue'
import BaseInput from '@/components/base/BaseInput.vue'
import BaseSelect from '@/components/base/BaseSelect.vue'

const store = useAdminQuizzesStore()

const columns: AdminTableColumn[] = [
  { key: 'question', label: 'Вопрос' },
  { key: 'hsk_level', label: 'HSK' },
  { key: 'options', label: 'Варианты' },
]

const hskOptions = [1, 2, 3, 4, 5, 6].map(level => ({ value: level, label: `HSK ${level}` }))
const directionOptions = [
  { value: 'word_to_translation', label: 'Слово → перевод' },
  { value: 'translation_to_word', label: 'Перевод → слово' },
]

const modalOpen = ref(false)
const saving = ref(false)
const formError = ref('')
const editingItem = ref<AdminQuiz | null>(null)

const deletingItem = ref<AdminQuiz | null>(null)
const deleting = ref(false)

const form = reactive({
  question: '',
  hskLevel: 1 as number,
  direction: 'word_to_translation' as QuizDirection,
  hanzi: '',
  pinyin: '',
  options: [] as AdminQuizOption[],
})

function resetForm() {
  form.question = ''
  form.hskLevel = 1
  form.direction = 'word_to_translation'
  form.hanzi = ''
  form.pinyin = ''
  form.options = [
    { text: '', pinyin: '', is_correct: true },
    { text: '', pinyin: '', is_correct: false },
  ]
}

function openCreate() {
  editingItem.value = null
  resetForm()
  formError.value = ''
  modalOpen.value = true
}

function openEdit(item: AdminQuiz) {
  editingItem.value = item
  form.question = item.question
  form.hskLevel = item.hsk_level
  form.direction = item.direction
  form.hanzi = item.hanzi
  form.pinyin = item.pinyin
  form.options = item.options.map(o => ({ ...o }))
  formError.value = ''
  modalOpen.value = true
}

async function onSubmit() {
  if (!form.options.some(o => o.is_correct)) {
    formError.value = 'Отметьте хотя бы один правильный вариант'
    return
  }

  saving.value = true
  formError.value = ''

  const payload: QuizPayload = {
    question: form.question,
    hsk_level: form.hskLevel,
    direction: form.direction,
    hanzi: form.hanzi,
    pinyin: form.pinyin,
    options: form.options.map(o => ({ text: o.text, pinyin: o.pinyin, is_correct: o.is_correct })),
  }

  try {
    if (editingItem.value) {
      await store.update(editingItem.value.id, payload)
    } else {
      await store.create(payload)
    }

    modalOpen.value = false
  } catch {
    formError.value = 'Не удалось сохранить тест'
  } finally {
    saving.value = false
  }
}

function confirmDelete(item: AdminQuiz) {
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
