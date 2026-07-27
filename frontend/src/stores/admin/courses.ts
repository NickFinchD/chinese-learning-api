import { acceptHMRUpdate, defineStore } from 'pinia'

import { adminListCourses, createCourse, deleteCourse, updateCourse } from '@/services/admin/courses'
import type { CoursePayload } from '@/services/admin/courses'

import type { Course } from '@/types/course'

export const useAdminCoursesStore = defineStore('adminCourses', {
  state: () => ({
    items: [] as Course[],
    total: 0,
    page: 1,
    limit: 20,
    search: '',
    loading: false,
    error: '',
  }),

  getters: {
    totalPages: state => Math.max(1, Math.ceil(state.total / state.limit)),
  },

  actions: {
    async fetchPage() {
      this.loading = true
      this.error = ''

      try {
        const response = await adminListCourses({
          search: this.search || undefined,
          page: this.page,
          limit: this.limit,
        })

        this.items = response.data.items
        this.total = response.data.total
      } catch {
        this.error = 'Не удалось загрузить курсы'
      } finally {
        this.loading = false
      }
    },

    setSearch(search: string) {
      this.search = search
      this.page = 1
      this.fetchPage()
    },

    setPage(page: number) {
      this.page = page
      this.fetchPage()
    },

    async create(payload: CoursePayload) {
      await createCourse(payload)
      await this.fetchPage()
    },

    async update(id: number, payload: CoursePayload) {
      await updateCourse(id, payload)
      await this.fetchPage()
    },

    async remove(id: number) {
      await deleteCourse(id)
      await this.fetchPage()
    },
  },
})

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useAdminCoursesStore, import.meta.hot))
}
