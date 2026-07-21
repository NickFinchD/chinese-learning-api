import { api } from './client'

import type { ApiResponse } from '@/types/api'
import type { Course, CourseDetails } from '@/types/course'

export async function getCourses() {
  console.log('getCourses request')
  console.log('baseURL =', api.defaults.baseURL)

  try {
    const response = await api.get<ApiResponse<Course[]>>('/courses')

    console.log('SUCCESS', response)

    return response.data
  } catch (e) {
    console.error('AXIOS ERROR', e)

    throw e
  }
}

export async function getCourse(id: number) {
  const response = await api.get<ApiResponse<CourseDetails>>(`/courses/${id}`)

  return response.data
}