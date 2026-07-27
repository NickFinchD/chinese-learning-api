import { api } from '@/services/client'

import type { ApiResponse } from '@/types/api'
import type { AdminListParams, PagedResult } from '@/types/admin'
import type { Course } from '@/types/course'

export interface CoursePayload {
  title: string
  description: string
  hsk_level: number
  sort_order: number
}

export async function adminGetCourse(id: number) {
  const response = await api.get<ApiResponse<Course>>(`/admin/courses/${id}`)

  return response.data
}

export async function adminListCourses(params: AdminListParams) {
  const response = await api.get<ApiResponse<PagedResult<Course>>>('/admin/courses', { params })

  return response.data
}

export async function createCourse(payload: CoursePayload) {
  const response = await api.post<ApiResponse<Course>>('/admin/courses', payload)

  return response.data
}

export async function updateCourse(id: number, payload: CoursePayload) {
  const response = await api.put<ApiResponse<Course>>(`/admin/courses/${id}`, payload)

  return response.data
}

export async function deleteCourse(id: number) {
  await api.delete(`/admin/courses/${id}`)
}
