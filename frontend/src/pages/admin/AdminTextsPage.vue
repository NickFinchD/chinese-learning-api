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
        Тексты
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
      :title="editingItem ? 'Редактировать текст' : 'Новый текст'"
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
        v-model="form.hanzi"
        placeholder="Текст (ханьцзы)"
        :rows="5"
      />
      <BaseTextarea
        v-model="form.pinyin"
        placeholder="Пиньинь"
        :rows="5"
      />
      <BaseTextarea
        v-model="form.translation"
        placeholder="Перевод"
        :rows="5"
      />
      <BaseSelect
        v-model="form.hskLevel"
        :options="hskOptions"
      />
    </AdminFormModal>

    <ConfirmDialog
      :open="deletingItem !== null"
      title="Удалить текст?"
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

import { useAdminTextsStore } from '@/stores/admin/texts'
import type { TextPayload } from '@/services/admin/texts'
import type { TextItem } from '@/types/text'

import AdminDataTable from '@/components/admin/AdminDataTable.vue'
import type { AdminTableColumn } from '@/components/admin/AdminDataTable.vue'
import AdminFormModal from '@/components/admin/AdminFormModal.vue'
import ConfirmDialog from '@/components/admin/ConfirmDialog.vue'
import AppIcon from '@/components/base/AppIcon.vue'
import BaseInput from '@/components/base/BaseInput.vue'
import BaseSelect from '@/components/base/BaseSelect.vue'
import BaseTextarea from '@/components/base/BaseTextarea.vue'

const store = useAdminTextsStore()

const columns: AdminTableColumn[] = [
  { key: 'title', label: 'Заголовок' },
  { key: 'hanzi', label: 'Текст' },
  { key: 'hsk_level', label: 'HSK' },
]

const hskOptions = [1, 2, 3, 4, 5, 6].map(level => ({ value: level, label: `HSK ${level}` }))

const modalOpen = ref(false)
const saving = ref(false)
const formError = ref('')
const editingItem = ref<TextItem | null>(null)

const deletingItem = ref<TextItem | null>(null)
const deleting = ref(false)

const form = reactive({
  title: '',
  hanzi: '',
  pinyin: '',
  translation: '',
  hskLevel: 1 as number,
})

function resetForm() {
  form.title = ''
  form.hanzi = ''
  form.pinyin = ''
  form.translation = ''
  form.hskLevel = 1
}

function openCreate() {
  editingItem.value = null
  resetForm()
  formError.value = ''
  modalOpen.value = true
}

function openEdit(item: TextItem) {
  editingItem.value = item
  form.title = item.title
  form.hanzi = item.hanzi
  form.pinyin = item.pinyin
  form.translation = item.translation
  form.hskLevel = item.hsk_level
  formError.value = ''
  modalOpen.value = true
}

async function onSubmit() {
  saving.value = true
  formError.value = ''

  const payload: TextPayload = {
    title: form.title,
    hanzi: form.hanzi,
    pinyin: form.pinyin,
    translation: form.translation,
    hsk_level: form.hskLevel,
  }

  try {
    if (editingItem.value) {
      await store.update(editingItem.value.id, payload)
    } else {
      await store.create(payload)
    }

    modalOpen.value = false
  } catch {
    formError.value = 'Не удалось сохранить текст'
  } finally {
    saving.value = false
  }
}

function confirmDelete(item: TextItem) {
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
