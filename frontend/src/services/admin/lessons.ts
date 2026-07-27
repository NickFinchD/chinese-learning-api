import { api } from '@/services/client'

import type { ApiResponse } from '@/types/api'

export interface AdminLessonListItem {
  id: number
  course_id: number
  title: string
  description: string
  lesson_number: number
  step_count: number
}

export interface LessonPayload {
  course_id: number
  title: string
  description: string
  lesson_number: number
}

export type StepType = 'word' | 'quiz' | 'grammar' | 'sentence_builder'

export interface AdminStep {
  id: number
  step_type: StepType
  entity_id: number
  sort_order: number
  label: string
}

export async function adminListLessonsByCourse(courseId: number) {
  const response = await api.get<ApiResponse<AdminLessonListItem[]>>('/admin/lessons', {
    params: { course_id: courseId },
  })

  return response.data
}

export async function createLesson(payload: LessonPayload) {
  const response = await api.post<ApiResponse<AdminLessonListItem>>('/admin/lessons', payload)

  return response.data
}

export async function updateLesson(id: number, payload: LessonPayload) {
  const response = await api.put<ApiResponse<AdminLessonListItem>>(`/admin/lessons/${id}`, payload)

  return response.data
}

export async function deleteLesson(id: number) {
  await api.delete(`/admin/lessons/${id}`)
}

export async function getLessonSteps(lessonId: number) {
  const response = await api.get<ApiResponse<AdminStep[]>>(`/admin/lessons/${lessonId}/steps`)

  return response.data
}

export async function addLessonStep(lessonId: number, stepType: StepType, entityId: number) {
  const response = await api.post<ApiResponse<AdminStep>>(`/admin/lessons/${lessonId}/steps`, {
    step_type: stepType,
    entity_id: entityId,
  })

  return response.data
}

export async function removeLessonStep(lessonId: number, stepId: number) {
  await api.delete(`/admin/lessons/${lessonId}/steps/${stepId}`)
}

export async function reorderLessonSteps(lessonId: number, stepIds: number[]) {
  await api.put(`/admin/lessons/${lessonId}/steps/reorder`, { step_ids: stepIds })
}
