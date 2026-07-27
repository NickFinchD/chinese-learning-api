<template>
  <div>
    <div class="mb-6 flex items-center gap-3">
      <RouterLink
        to="/app/admin/courses"
        class="flex h-9 w-9 items-center justify-center rounded-lg text-gray-600 transition hover:bg-white/60 dark:text-gray-300 dark:hover:bg-white/10"
      >
        <AppIcon name="arrow-left" />
      </RouterLink>

      <h1 class="text-3xl font-bold text-gray-900 dark:text-white">
        {{ course?.title ?? 'Загрузка...' }}
      </h1>
    </div>

    <BaseCard
      v-if="course"
      class="mb-8 max-w-2xl"
    >
      <div class="mb-4 flex items-start justify-between gap-4">
        <div>
          <p class="text-gray-600 dark:text-gray-400">
            {{ course.description || 'Без описания' }}
          </p>
          <span class="mt-2 inline-flex rounded-full bg-[var(--color-primary)]/15 px-3 py-1 text-sm font-medium text-[var(--color-primary)] dark:bg-[var(--color-primary)]/20 dark:text-[var(--color-mint)]">
            HSK {{ course.hsk_level }}
          </span>
        </div>

        <button
          type="button"
          title="Редактировать курс"
          class="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg text-gray-600 transition hover:bg-white/60 dark:text-gray-300 dark:hover:bg-white/10"
          @click="openEditCourse"
        >
          <AppIcon
            name="pencil"
            :size="16"
          />
        </button>
      </div>
    </BaseCard>

    <div class="mb-4 flex items-center justify-between">
      <h2 class="text-xl font-semibold text-gray-900 dark:text-white">
        Уроки
      </h2>

      <button
        type="button"
        class="rounded-full bg-[var(--color-primary)] px-5 py-2.5 text-sm font-semibold text-white shadow-lg shadow-[var(--color-primary)]/30 transition hover:bg-[var(--color-primary)]/90"
        @click="openCreateLesson"
      >
        Добавить урок
      </button>
    </div>

    <div
      v-if="lessons.loading"
      class="flex justify-center py-12"
    >
      <BaseSpinner :size="28" />
    </div>

    <div
      v-else
      class="space-y-2"
    >
      <div
        v-for="lesson in lessons.items"
        :key="lesson.id"
        class="flex items-center gap-3 rounded-xl border border-white/50 bg-white/30 p-4 backdrop-blur-xl dark:border-white/10 dark:bg-white/5"
      >
        <RouterLink
          :to="{ name: 'admin-lesson', params: { id: lesson.id } }"
          class="min-w-0 flex-1 truncate font-medium text-gray-900 hover:text-[var(--color-primary)] dark:text-white"
        >
          {{ lesson.lesson_number }}. {{ lesson.title }}
        </RouterLink>

        <span class="shrink-0 text-sm text-gray-500 dark:text-gray-400">
          {{ lesson.step_count }} шаг(ов)
        </span>

        <button
          type="button"
          title="Редактировать"
          class="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg text-gray-600 transition hover:bg-white/60 dark:text-gray-300 dark:hover:bg-white/10"
          @click="openEditLesson(lesson)"
        >
          <AppIcon
            name="pencil"
            :size="16"
          />
        </button>

        <button
          type="button"
          title="Удалить"
          class="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg text-red-500 transition hover:bg-red-500/10"
          @click="deletingLesson = lesson"
        >
          <AppIcon
            name="trash"
            :size="16"
          />
        </button>
      </div>

      <p
        v-if="lessons.items.length === 0"
        class="py-10 text-center text-gray-500 dark:text-gray-400"
      >
        В этом курсе пока нет уроков
      </p>
    </div>

    <p
      v-if="lessons.error"
      class="mt-4 text-sm text-red-600 dark:text-red-400"
    >
      {{ lessons.error }}
    </p>

    <AdminFormModal
      :open="courseModalOpen"
      title="Редактировать курс"
      :saving="savingCourse"
      :error="courseFormError"
      @close="courseModalOpen = false"
      @submit="onSubmitCourse"
    >
      <BaseInput
        v-model="courseForm.title"
        placeholder="Название"
        required
      />
      <BaseTextarea
        v-model="courseForm.description"
        placeholder="Описание"
        :rows="3"
      />
      <div class="grid grid-cols-2 gap-2">
        <BaseSelect
          v-model="courseForm.hskLevel"
          :options="hskOptions"
        />
        <BaseInput
          :model-value="String(courseForm.sortOrder)"
          placeholder="Порядок"
          @update:model-value="courseForm.sortOrder = Number($event) || 0"
        />
      </div>
    </AdminFormModal>

    <AdminFormModal
      :open="lessonModalOpen"
      :title="editingLesson ? 'Редактировать урок' : 'Новый урок'"
      :saving="savingLesson"
      :error="lessonFormError"
      @close="lessonModalOpen = false"
      @submit="onSubmitLesson"
    >
      <BaseInput
        v-model="lessonForm.title"
        placeholder="Название урока"
        required
      />
      <BaseTextarea
        v-model="lessonForm.description"
        placeholder="Описание"
        :rows="2"
      />
      <BaseInput
        :model-value="String(lessonForm.lessonNumber)"
        placeholder="Номер урока"
        required
        @update:model-value="lessonForm.lessonNumber = Number($event) || 0"
      />
    </AdminFormModal>

    <ConfirmDialog
      :open="deletingLesson !== null"
      title="Удалить урок?"
      :message="deletingLesson ? `«${deletingLesson.title}» вместе со всеми шагами. Это действие нельзя отменить.` : ''"
      :loading="deletingLessonInFlight"
      @cancel="deletingLesson = null"
      @confirm="onDeleteLesson"
    />
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { RouterLink, useRoute } from 'vue-router'

import { adminGetCourse, updateCourse } from '@/services/admin/courses'
import type { CoursePayload } from '@/services/admin/courses'
import { useAdminLessonsStore } from '@/stores/admin/lessons'
import type { AdminLessonListItem, LessonPayload } from '@/services/admin/lessons'
import type { Course } from '@/types/course'

import AdminFormModal from '@/components/admin/AdminFormModal.vue'
import ConfirmDialog from '@/components/admin/ConfirmDialog.vue'
import AppIcon from '@/components/base/AppIcon.vue'
import BaseCard from '@/components/base/BaseCard.vue'
import BaseInput from '@/components/base/BaseInput.vue'
import BaseSelect from '@/components/base/BaseSelect.vue'
import BaseSpinner from '@/components/base/BaseSpinner.vue'
import BaseTextarea from '@/components/base/BaseTextarea.vue'

const route = useRoute()
const courseId = Number(route.params.id)

const lessons = useAdminLessonsStore()

const course = ref<Course | null>(null)
const hskOptions = [1, 2, 3, 4, 5, 6].map(level => ({ value: level, label: `HSK ${level}` }))

const courseModalOpen = ref(false)
const savingCourse = ref(false)
const courseFormError = ref('')
const courseForm = reactive({
  title: '',
  description: '',
  hskLevel: 1 as number,
  sortOrder: 0,
})

function openEditCourse() {
  if (!course.value) return

  courseForm.title = course.value.title
  courseForm.description = course.value.description
  courseForm.hskLevel = course.value.hsk_level
  courseForm.sortOrder = course.value.sort_order
  courseFormError.value = ''
  courseModalOpen.value = true
}

async function onSubmitCourse() {
  savingCourse.value = true
  courseFormError.value = ''

  const payload: CoursePayload = {
    title: courseForm.title,
    description: courseForm.description,
    hsk_level: courseForm.hskLevel,
    sort_order: courseForm.sortOrder,
  }

  try {
    const response = await updateCourse(courseId, payload)
    course.value = response.data
    courseModalOpen.value = false
  } catch {
    courseFormError.value = 'Не удалось сохранить курс'
  } finally {
    savingCourse.value = false
  }
}

const lessonModalOpen = ref(false)
const savingLesson = ref(false)
const lessonFormError = ref('')
const editingLesson = ref<AdminLessonListItem | null>(null)
const lessonForm = reactive({
  title: '',
  description: '',
  lessonNumber: 1 as number,
})

const deletingLesson = ref<AdminLessonListItem | null>(null)
const deletingLessonInFlight = ref(false)

function openCreateLesson() {
  editingLesson.value = null
  lessonForm.title = ''
  lessonForm.description = ''
  lessonForm.lessonNumber = lessons.items.length + 1
  lessonFormError.value = ''
  lessonModalOpen.value = true
}

function openEditLesson(lesson: AdminLessonListItem) {
  editingLesson.value = lesson
  lessonForm.title = lesson.title
  lessonForm.description = lesson.description
  lessonForm.lessonNumber = lesson.lesson_number
  lessonFormError.value = ''
  lessonModalOpen.value = true
}

async function onSubmitLesson() {
  savingLesson.value = true
  lessonFormError.value = ''

  const payload: LessonPayload = {
    course_id: courseId,
    title: lessonForm.title,
    description: lessonForm.description,
    lesson_number: lessonForm.lessonNumber,
  }

  try {
    if (editingLesson.value) {
      await lessons.update(editingLesson.value.id, payload)
    } else {
      await lessons.create(payload)
    }

    lessonModalOpen.value = false
  } catch {
    lessonFormError.value = 'Не удалось сохранить урок'
  } finally {
    savingLesson.value = false
  }
}

async function onDeleteLesson() {
  if (!deletingLesson.value) return

  deletingLessonInFlight.value = true

  try {
    await lessons.remove(deletingLesson.value.id)
    deletingLesson.value = null
  } finally {
    deletingLessonInFlight.value = false
  }
}

onMounted(async () => {
  const response = await adminGetCourse(courseId)
  course.value = response.data

  lessons.fetchForCourse(courseId)
})
</script>
