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
        Слова
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
      :title="editingItem ? 'Редактировать слово' : 'Новое слово'"
      :saving="saving"
      :error="formError"
      @close="modalOpen = false"
      @submit="onSubmit"
    >
      <BaseInput
        v-model="form.hanzi"
        placeholder="Ханьцзы"
        required
      />
      <BaseInput
        v-model="form.pinyin"
        placeholder="Пиньинь"
        required
      />
      <BaseInput
        v-model="form.translation"
        placeholder="Перевод"
        required
      />
      <BaseInput
        v-model="form.partOfSpeech"
        placeholder="Часть речи"
        required
      />
      <BaseSelect
        v-model="form.hskLevel"
        :options="hskOptions"
      />

      <div>
        <div class="mb-2 text-sm font-medium text-gray-700 dark:text-gray-300">
          Примеры
        </div>

        <AdminRepeatableList
          v-model="form.examples"
          add-label="Добавить пример"
          :new-item="() => ({ hanzi: '', pinyin: '', translation: '' })"
        >
          <template #default="{ item, update }">
            <div class="grid grid-cols-3 gap-2">
              <BaseInput
                :model-value="item.hanzi"
                placeholder="Ханьцзы"
                @update:model-value="update({ ...item, hanzi: $event })"
              />
              <BaseInput
                :model-value="item.pinyin"
                placeholder="Пиньинь"
                @update:model-value="update({ ...item, pinyin: $event })"
              />
              <BaseInput
                :model-value="item.translation"
                placeholder="Перевод"
                @update:model-value="update({ ...item, translation: $event })"
              />
            </div>
          </template>
        </AdminRepeatableList>
      </div>
    </AdminFormModal>

    <ConfirmDialog
      :open="deletingItem !== null"
      title="Удалить слово?"
      :message="deletingItem ? `«${deletingItem.hanzi}» — ${deletingItem.translation}. Это действие нельзя отменить.` : ''"
      :loading="deleting"
      @cancel="deletingItem = null"
      @confirm="onDelete"
    />
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { RouterLink } from 'vue-router'

import { useAdminWordsStore } from '@/stores/admin/words'
import type { ExamplePayload, WordPayload } from '@/services/admin/words'
import { getWord } from '@/services/words'
import type { Word } from '@/types/word'

import AdminDataTable from '@/components/admin/AdminDataTable.vue'
import type { AdminTableColumn } from '@/components/admin/AdminDataTable.vue'
import AdminFormModal from '@/components/admin/AdminFormModal.vue'
import AdminRepeatableList from '@/components/admin/AdminRepeatableList.vue'
import ConfirmDialog from '@/components/admin/ConfirmDialog.vue'
import AppIcon from '@/components/base/AppIcon.vue'
import BaseInput from '@/components/base/BaseInput.vue'
import BaseSelect from '@/components/base/BaseSelect.vue'

const store = useAdminWordsStore()

const columns: AdminTableColumn[] = [
  { key: 'hanzi', label: 'Ханьцзы' },
  { key: 'pinyin', label: 'Пиньинь' },
  { key: 'translation', label: 'Перевод' },
  { key: 'part_of_speech', label: 'Часть речи' },
  { key: 'hsk_level', label: 'HSK' },
]

const hskOptions = [1, 2, 3, 4, 5, 6].map(level => ({ value: level, label: `HSK ${level}` }))

const modalOpen = ref(false)
const saving = ref(false)
const formError = ref('')
const editingItem = ref<Word | null>(null)

const deletingItem = ref<Word | null>(null)
const deleting = ref(false)

const form = reactive({
  hanzi: '',
  pinyin: '',
  translation: '',
  partOfSpeech: '',
  hskLevel: 1 as number,
  examples: [] as ExamplePayload[],
})

function resetForm() {
  form.hanzi = ''
  form.pinyin = ''
  form.translation = ''
  form.partOfSpeech = ''
  form.hskLevel = 1
  form.examples = []
}

function openCreate() {
  editingItem.value = null
  resetForm()
  formError.value = ''
  modalOpen.value = true
}

async function openEdit(item: Word) {
  editingItem.value = item
  form.hanzi = item.hanzi
  form.pinyin = item.pinyin
  form.translation = item.translation
  form.partOfSpeech = item.part_of_speech
  form.hskLevel = item.hsk_level
  form.examples = []
  formError.value = ''
  modalOpen.value = true

  // The admin list doesn't carry examples (avoids an N+1 query on every
  // row) — fetch the full detail, which does, once the form is open.
  const detail = await getWord(item.id)
  form.examples = (detail.data.examples ?? []).map(e => ({
    hanzi: e.hanzi,
    pinyin: e.pinyin,
    translation: e.translation,
  }))
}

async function onSubmit() {
  saving.value = true
  formError.value = ''

  const payload: WordPayload = {
    hanzi: form.hanzi,
    pinyin: form.pinyin,
    translation: form.translation,
    part_of_speech: form.partOfSpeech,
    hsk_level: form.hskLevel,
    examples: form.examples,
  }

  try {
    if (editingItem.value) {
      await store.update(editingItem.value.id, payload)
    } else {
      await store.create(payload)
    }

    modalOpen.value = false
  } catch {
    formError.value = 'Не удалось сохранить слово'
  } finally {
    saving.value = false
  }
}

function confirmDelete(item: Word) {
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
