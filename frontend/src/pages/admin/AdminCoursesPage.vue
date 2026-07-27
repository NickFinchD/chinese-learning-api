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
        Курсы и уроки
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
      <template #cell-title="{ item }">
        <RouterLink
          :to="{ name: 'admin-course', params: { id: item.id } }"
          class="font-medium text-[var(--color-primary)] hover:text-[var(--color-primary)]/80"
        >
          {{ item.title }}
        </RouterLink>
      </template>
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
      :title="editingItem ? 'Редактировать курс' : 'Новый курс'"
      :saving="saving"
      :error="formError"
      @close="modalOpen = false"
      @submit="onSubmit"
    >
      <BaseInput
        v-model="form.title"
        placeholder="Название"
        required
      />
      <BaseTextarea
        v-model="form.description"
        placeholder="Описание"
        :rows="3"
      />
      <div class="grid grid-cols-2 gap-2">
        <BaseSelect
          v-model="form.hskLevel"
          :options="hskOptions"
        />
        <BaseInput
          :model-value="String(form.sortOrder)"
          placeholder="Порядок"
          @update:model-value="form.sortOrder = Number($event) || 0"
        />
      </div>
    </AdminFormModal>

    <ConfirmDialog
      :open="deletingItem !== null"
      title="Удалить курс?"
      :message="deletingItem ? `«${deletingItem.title}» вместе со всеми уроками и шагами. Это действие нельзя отменить.` : ''"
      :loading="deleting"
      @cancel="deletingItem = null"
      @confirm="onDelete"
    />
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { RouterLink } from 'vue-router'

import { useAdminCoursesStore } from '@/stores/admin/courses'
import type { CoursePayload } from '@/services/admin/courses'
import type { Course } from '@/types/course'

import AdminDataTable from '@/components/admin/AdminDataTable.vue'
import type { AdminTableColumn } from '@/components/admin/AdminDataTable.vue'
import AdminFormModal from '@/components/admin/AdminFormModal.vue'
import ConfirmDialog from '@/components/admin/ConfirmDialog.vue'
import AppIcon from '@/components/base/AppIcon.vue'
import BaseInput from '@/components/base/BaseInput.vue'
import BaseSelect from '@/components/base/BaseSelect.vue'
import BaseTextarea from '@/components/base/BaseTextarea.vue'

const store = useAdminCoursesStore()

const columns: AdminTableColumn[] = [
  { key: 'title', label: 'Название' },
  { key: 'hsk_level', label: 'HSK' },
  { key: 'sort_order', label: 'Порядок' },
]

const hskOptions = [1, 2, 3, 4, 5, 6].map(level => ({ value: level, label: `HSK ${level}` }))

const modalOpen = ref(false)
const saving = ref(false)
const formError = ref('')
const editingItem = ref<Course | null>(null)

const deletingItem = ref<Course | null>(null)
const deleting = ref(false)

const form = reactive({
  title: '',
  description: '',
  hskLevel: 1 as number,
  sortOrder: 0,
})

function resetForm() {
  form.title = ''
  form.description = ''
  form.hskLevel = 1
  form.sortOrder = 0
}

function openCreate() {
  editingItem.value = null
  resetForm()
  formError.value = ''
  modalOpen.value = true
}

function openEdit(item: Course) {
  editingItem.value = item
  form.title = item.title
  form.description = item.description
  form.hskLevel = item.hsk_level
  form.sortOrder = item.sort_order
  formError.value = ''
  modalOpen.value = true
}

async function onSubmit() {
  saving.value = true
  formError.value = ''

  const payload: CoursePayload = {
    title: form.title,
    description: form.description,
    hsk_level: form.hskLevel,
    sort_order: form.sortOrder,
  }

  try {
    if (editingItem.value) {
      await store.update(editingItem.value.id, payload)
    } else {
      await store.create(payload)
    }

    modalOpen.value = false
  } catch {
    formError.value = 'Не удалось сохранить курс'
  } finally {
    saving.value = false
  }
}

function confirmDelete(item: Course) {
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
