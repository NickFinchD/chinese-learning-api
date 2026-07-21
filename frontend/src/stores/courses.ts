import { defineStore } from 'pinia'

import { getCourse, getCourses } from '@/services/courses'

import type { Course, CourseDetails } from '@/types/course'

export const useCoursesStore = defineStore('courses', {
  state: () => ({
    items: [] as Course[],
    current: null as CourseDetails | null,
    loading: false,
  }),

  actions: {
    async loadCourses() {
      console.log('loadCourses called')

      this.loading = true

      try {
        const response = await getCourses()

        console.log(response)

        this.items = response.data
      } finally {
        this.loading = false
      }
    },

    async loadCourse(id: number) {
      this.loading = true

      try {
        const response = await getCourse(id)

        this.current = response.data
      } finally {
        this.loading = false
      }
    },
  },
})