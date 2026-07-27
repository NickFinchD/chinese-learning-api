import { acceptHMRUpdate, defineStore } from 'pinia'

import { adminListLessonsByCourse, createLesson, deleteLesson, updateLesson } from '@/services/admin/lessons'
import type { AdminLessonListItem, LessonPayload } from '@/services/admin/lessons'

export const useAdminLessonsStore = defineStore('adminLessons', {
  state: () => ({
    courseId: null as number | null,
    items: [] as AdminLessonListItem[],
    loading: false,
    error: '',
  }),

  actions: {
    async fetchForCourse(courseId: number) {
      this.courseId = courseId
      this.loading = true
      this.error = ''

      try {
        const response = await adminListLessonsByCourse(courseId)

        this.items = response.data
      } catch {
        this.error = 'Не удалось загрузить уроки'
      } finally {
        this.loading = false
      }
    },

    async create(payload: LessonPayload) {
      await createLesson(payload)
      await this.fetchForCourse(payload.course_id)
    },

    async update(id: number, payload: LessonPayload) {
      await updateLesson(id, payload)
      await this.fetchForCourse(payload.course_id)
    },

    async remove(id: number) {
      if (this.courseId === null) return
      await deleteLesson(id)
      await this.fetchForCourse(this.courseId)
    },
  },
})

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useAdminLessonsStore, import.meta.hot))
}
